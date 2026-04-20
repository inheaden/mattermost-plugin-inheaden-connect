# Mattermost Plugin for [Inheaden Connect](https://inco.video)

> ⚠️ **This plugin is currently under active development and is not fully production ready.**

Bring **Inheaden Connect** video meetings directly into your Mattermost workspace. Start a meeting from any channel with one click, post a join link for your teammates, and manage per-channel meeting rooms — all without leaving Mattermost.

## Features

- Start Inheaden Connect meetings from the channel header button
- Select a meeting room from a searchable list
- Post a meeting join link into the channel for teammates
- Configure a dedicated meeting room per channel (`/inco meetingroom`)
- Admin-configurable API endpoint and credentials

## Quick Start

1. Download the latest release bundle (`.tar.gz`) from the [releases page](https://github.com/inheaden/mattermost-plugin-inheaden-connect/releases).
2. Upload it in **System Console → Plugins → Plugin Management**.
3. Configure the plugin with your Inheaden Connect URL and API key (see [Configuration](#configuration)).
4. Click the Inheaden Connect icon in any channel header to start a meeting.

For full setup instructions see [docs/overview.md](docs/overview.md).

## Configuration

| Setting | Description | Default |
|---|---|---|
| **Inheaden Connect URL** | Base URL of your Inheaden Connect instance | `https://inco.video` |
| **API Key** | API key (`XXX:YYY`) with permission to read meeting rooms | *(required)* |

## Slash Commands

| Command | Description |
|---|---|
| `/inco meetingroom <id>` | Assign a specific meeting room to the current channel |
| `/inco meetingroom` | Remove the custom meeting room assignment for the current channel |
| `/inco help` | Show available commands |

## Development

```bash
# Build and bundle the plugin
make dist

# Run tests
make test

# Deploy to a local Mattermost server
make deploy
```

See the [Mattermost Plugin Developer Guide](https://developers.mattermost.com/extend/plugins/) for environment setup and further guidance.

## Documentation

- [Overview & Architecture](docs/overview.md)
- [Mattermost Plugin Developer Guide](https://developers.mattermost.com/extend/plugins/)
- [Inheaden Connect Support](https://inco.video/support)

## License

[MIT](LICENSE)
