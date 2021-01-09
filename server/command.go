package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-plugin-api/experimental/command"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

func (p *Plugin) getCommand() (*model.Command, error) {
	iconData, err := command.GetIconData(p.API, "assets/PM_1x1_CT.svg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get icon data")
	}

	return &model.Command{
		Trigger:              "inco",
		AutoComplete:         true,
		AutoCompleteDesc:     "Available commands: meetingroom, help",
		AutoCompleteHint:     "[command] [args]",
		AutocompleteData:     p.getAutocompleteData(),
		AutocompleteIconData: iconData,
	}, nil
}

func (p *Plugin) executeCommand(c *plugin.Context, args *model.CommandArgs) (string, error) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""

	if command != "/inco" {
		return fmt.Sprintf("Command '%s' is not /inco. Please try again.", command), nil
	}

	if len(split) > 1 {
		action = split[1]
	} else {
		return "Please specify an action for /inco command.", nil
	}

	userID := args.UserId
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		return fmt.Sprintf("We could not retrieve user (userId: %v)", args.UserId), nil
	}

	switch action {
	case "meetingroom":
		return p.runMeetingRoomCommand(args, user)
	case "help":
		return p.runHelpCommand()
	default:
		return fmt.Sprintf("Unknown action %v", action), nil
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	msg, err := p.executeCommand(c, args)
	if err != nil {
		p.API.LogWarn("failed to execute command", "error", err.Error())
	}
	if msg != "" {
		p.postCommandResponse(args, msg)
	}
	return &model.CommandResponse{}, nil
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    args.UserId,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

func (p *Plugin) runMeetingRoomCommand(args *model.CommandArgs, user *model.User) (string, error) {
	split := strings.Fields(args.Command)
	if len(split) == 2 {
		p.API.KVDelete(args.ChannelId)
		return "Removed custom meeting id for this channel.", nil
	}

	if len(split) != 3 {
		return fmt.Sprintf("Please specify the meeting id to use for this channel."), nil
	}
	meetingID := split[2]

	p.API.KVSet(args.ChannelId, []byte(meetingID))

	return "Meeting id for this channel has been set.", nil
}

func (p *Plugin) runHelpCommand() (string, error) {
	return strings.ReplaceAll(`### Inheaden Connect plugin
* |/inco meetingRoom [meetingRoomId]| Set a custom meeting room id for this channel.
* |/inco help| Show this help message.
	`, "|", "`"), nil
}

func (p *Plugin) getAutocompleteData() *model.AutocompleteData {
	inco := model.NewAutocompleteData("inco", "[command]", "Starts a meeting on Inheaden Connect.")
	meetingRoom := model.NewAutocompleteData("meetingroom", "[meetingRoomId]", "Sets a meeting room id for this channel. Leave empty for default.")
	help := model.NewAutocompleteData("help", "", "Displays helpful information about this command.")

	inco.AddCommand(meetingRoom)
	inco.AddCommand(help)
	return inco
}
