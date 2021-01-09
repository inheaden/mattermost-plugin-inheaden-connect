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

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/api/v1/meetings":
		p.handleStartMeeting(w, r)
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.API.LogError("error when trying to read response", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var startMeetingRequest StartMeetingRequest
	err = json.Unmarshal(body, &startMeetingRequest)
	if err != nil {
		p.API.LogError("error when trying to pares request", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.API.LogDebug(fmt.Sprintf("request body: %s", startMeetingRequest))

	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		http.Error(w, appErr.Error(), appErr.StatusCode)
		return
	}

	var response JoinResponse

	fullName := user.GetFullName()
	if len(fullName) == 0 {
		fullName = user.Username
	}

	meetingID := p.configuration.DefaultMeetingRoomID
	if len(startMeetingRequest.ChannelID) != 0 {
		id, _ := p.API.KVGet(startMeetingRequest.ChannelID)
		if id != nil {
			meetingID = string(id)
		}
	}

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

func (p *Plugin) createMeeting(w http.ResponseWriter, r *http.Request, joinRequest JoinRequest, response *JoinResponse, meetingID string) {

	apiURL := p.configuration.InheadenConnectAPIURL
	apiKey := p.configuration.APIKey

	requestBody, err := json.Marshal(joinRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
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

// See https://developers.mattermost.com/extend/plugins/server/reference/
