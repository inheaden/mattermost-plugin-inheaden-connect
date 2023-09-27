// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "io.inheaden.inheaden-connect",
  "name": "Inheaden Connect plugin",
  "description": "This plugin allows you to start meetings from within mattermost.",
  "homepage_url": "https://inco.video",
  "support_url": "https://inco.video/support",
  "release_notes_url": "https://example.com/releases/v0.0.1",
  "icon_path": "assets/AppIcon.png",
  "version": "0.3.0",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "To setup this plugin, create an ApiKey on Inheaden Connect.",
    "footer": "",
    "settings": [
      {
        "key": "InheadenConnectAPIURL",
        "display_name": "Inheaden Connect URL",
        "type": "text",
        "help_text": "The URL for the Inheaden Connect Api.",
        "placeholder": "https://inco.video",
        "default": "https://inco.video"
      },
      {
        "key": "APIKey",
        "display_name": "Your api key",
        "type": "text",
        "help_text": "The api key for your account. It needs access to read all meeting rooms.",
        "placeholder": "XXX:YYY",
        "default": null
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
