package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// JoinResponse will be returned by the Inheaden Connect backend.
type JoinResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	JoinURL string `json:"joinUrl"`
}

// JoinRequest is used when trying to join a meeting.
type JoinRequest struct {
	FullName string `json:"fullName"`
}

type StartMeetingRequest struct {
	ChannelID string `json:"channel_id"`
	RoomID    string `json:"room_id"`
	RoomName  string `json:"room_name"`
}

type GetMeetingRoomsRequest struct {
	Name string `json:"name"`
}

type MeetingRoomType struct {
	MaxParticipants int `json:"maxParticipants"`
}

type MeetingRoomResponse struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Type      MeetingRoomType `json:"type"`
	RoomToken string          `json:"roomToken"`
	Password  string          `json:"password"`
}

func (p *Plugin) OnActivate() error {
	command, err := p.getCommand()
	if err != nil {
		return errors.Wrap(err, "failed to get command")
	}

	err = p.API.RegisterCommand(command)
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/api/v1/startMeeting":
		p.handleStartMeeting(w, r)
	case "/api/v1/showMeetingPost":
		p.handleShowMeetingPost(w, r)
	case "/api/v1/getAllMeetingRooms":
		p.handleGetAllMeetingRooms(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleStartMeeting(w http.ResponseWriter, r *http.Request) {
	p.API.LogInfo("handleStartMeeting")

	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return
	}

	startMeetingRequest := p.getStartMeetingRequest(w, r)
	if startMeetingRequest == nil {
		return
	}

	var response JoinResponse

	fullName := user.GetFullName()
	if len(fullName) == 0 {
		fullName = user.Username
	}

	meetingID := startMeetingRequest.RoomID

	p.createMeeting(w, r, JoinRequest{
		FullName: fullName,
	}, &response, meetingID)

	result, err := json.Marshal(map[string]string{
		"joinUrl": response.JoinURL,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (p *Plugin) getStartMeetingRequest(w http.ResponseWriter, r *http.Request) *StartMeetingRequest {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.API.LogError("error when trying to read response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	var startMeetingRequest StartMeetingRequest
	err = json.Unmarshal(body, &startMeetingRequest)
	if err != nil {
		p.API.LogError("error when trying to pares request", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	p.API.LogDebug(fmt.Sprintf("request body: %s", startMeetingRequest))

	return &startMeetingRequest
}

func (p *Plugin) createMeeting(w http.ResponseWriter, r *http.Request, joinRequest JoinRequest, response *JoinResponse, meetingID string) {

	apiURL := p.configuration.InheadenConnectAPIURL
	apiKey := p.configuration.APIKey

	requestBody, err := json.Marshal(joinRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/connect/v1/meetingRoom/%s/join", apiURL, meetingID), bytes.NewBuffer(requestBody))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(apiKey))))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.API.LogDebug("starting request")
	resp, err := client.Do(request)
	if err != nil {
		p.API.LogError("error when trying to create meeting", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		message := fmt.Sprintf("error when trying to create meeting: %d", resp.StatusCode)
		p.API.LogError(message)
		http.Error(w, message, resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("error when trying to read response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		message := fmt.Sprintf("error when trying to create meeting: %s", response.Message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
}

func (p *Plugin) handleShowMeetingPost(w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return
	}

	fullName := user.GetFullName()
	if len(fullName) == 0 {
		fullName = user.Username
	}

	startMeetingRequest := p.getStartMeetingRequest(w, r)
	if startMeetingRequest == nil {
		return
	}

	meetingRoom, err := p.getMeetingRoomById(w, r, startMeetingRequest.RoomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	joinURL := p.makeJoinURL(meetingRoom)
	textPost := &model.Post{
		UserId:    userID,
		ChannelId: startMeetingRequest.ChannelID,
		Message: fmt.Sprintf(`Inheaden Connect Meeting in **%s**

[Join Meeting](%s)`, startMeetingRequest.RoomName, joinURL),
		Type: "custom_inco_start_meeting",
	}

	textPost.Props = model.StringInterface{
		"from_webhook":      "true",
		"creator_name":      fullName,
		"override_username": "Inheaden Connect",
		"override_icon_url": "https://cdn.inheaden.cloud/inco/brand/App%20Icons/AppIcon__512x512.png",
		"room_id":           startMeetingRequest.RoomID,
		"room_name":         startMeetingRequest.RoomName,
		"join_url":          joinURL,
	}

	_, appErr = p.API.CreatePost(textPost)
	if appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return
	}

	response, err := json.Marshal(map[string]string{
		"success": "ok",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (p *Plugin) makeJoinURL(meetingRoom *MeetingRoomResponse) string {
	apiURL := p.configuration.InheadenConnectAPIURL

	return fmt.Sprintf("%s/app/join?token=%s&password=%s", apiURL, meetingRoom.RoomToken, meetingRoom.Password)
}

func (p *Plugin) handleGetAllMeetingRooms(w http.ResponseWriter, r *http.Request) {
	p.API.LogInfo("handleStartMeeting")

	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	userID := r.Header.Get("Mattermost-User-Id")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	request := p.getMeetingRoomsRequest(w, r)

	response, err := p.getAllMeetingRomms(request, w, r)
	if err != nil {
		return
	}

	res, _ := json.Marshal(response)
	w.Write(res)
}

func (p *Plugin) getMeetingRoomsRequest(w http.ResponseWriter, r *http.Request) *GetMeetingRoomsRequest {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.API.LogError("error when trying to read response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	var request GetMeetingRoomsRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		p.API.LogError("error when trying to parse request", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	p.API.LogDebug(fmt.Sprintf("request body: %s", request))

	return &request
}

func (p *Plugin) getMeetingRoomById(w http.ResponseWriter, r *http.Request, meetingID string) (*MeetingRoomResponse, error) {
	apiURL := p.configuration.InheadenConnectAPIURL
	apiKey := p.configuration.APIKey

	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/connect/v1/meetingRoom/%s", apiURL, meetingID), http.NoBody)
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(apiKey))))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	p.API.LogDebug("starting request")
	resp, err := client.Do(request)
	if err != nil {
		p.API.LogError("error when trying to get meeting room", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		message := fmt.Sprintf("error when trying to get meeting room: %d", resp.StatusCode)
		p.API.LogError(message)
		http.Error(w, message, resp.StatusCode)
		return nil, errors.New(message)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("error when trying to read response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	var response MeetingRoomResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &response, nil
}

type FilterResponse struct {
	Elements []MeetingRoomResponse `json:"elements"`
}

func (p *Plugin) getAllMeetingRomms(meetingRoomsRequest *GetMeetingRoomsRequest, w http.ResponseWriter, r *http.Request) ([]MeetingRoomResponse, error) {
	apiURL := p.configuration.InheadenConnectAPIURL
	apiKey := p.configuration.APIKey

	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}

	requestBodyMap := map[string]interface{}{
		"paging": map[string]interface{}{
			"pageSize":   10,
			"pageNumber": 0,
			"sorting": []map[string]interface{}{
				{
					"sortBy":  "name",
					"sortDir": "asc",
				},
			},
		},
	}

	if meetingRoomsRequest.Name != "" {
		requestBodyMap["filter"] = map[string]interface{}{
			"comparator": map[string]interface{}{
				"attribute":  "name",
				"comparator": "isLike",
				"value":      meetingRoomsRequest.Name,
			},
		}
	}

	requestBody, err := json.Marshal(requestBodyMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/connect/v1/meetingRoom/filter", apiURL), bytes.NewBuffer(requestBody))
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(apiKey))))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	p.API.LogDebug("starting request")
	resp, err := client.Do(request)
	if err != nil {
		p.API.LogError("error when trying to get meeting rooms", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		message := fmt.Sprintf("error when trying to get meeting rooms: %d", resp.StatusCode)
		p.API.LogError(message)
		http.Error(w, message, resp.StatusCode)
		return nil, errors.New(message)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("error when trying to read response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	var response FilterResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return response.Elements, nil
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
