# Inheaden Connect Plugin — Overview

This document describes the architecture, setup, and usage of the Mattermost plugin for [Inheaden Connect](https://inco.video).

For general Mattermost plugin concepts, refer to the [Mattermost Plugin Developer Guide](https://developers.mattermost.com/extend/plugins/).

---

## What It Does

The plugin integrates Inheaden Connect video meetings into Mattermost. Users can:

- Open a meeting room picker modal from the channel header.
- Search and select a meeting room from their Inheaden Connect account.
- Post a join link into the current channel so all members can join with one click.
- Assign a persistent default meeting room to a channel via slash command.

---

## Architecture

The plugin follows the standard Mattermost plugin structure with a **server** component (Go) and a **webapp** component (TypeScript/React).

### Server (Go)

The server component registers the `/inco` slash command on activation and exposes three internal HTTP endpoints that the webapp calls:

| Endpoint | Method | Description |
|---|---|---|
| `/api/v1/startMeeting` | `POST` | Calls the Inheaden Connect API to create a join URL for the given room |
| `/api/v1/showMeetingPost` | `POST` | Posts a meeting message with the join link into the channel |
| `/api/v1/getAllMeetingRooms` | `POST` | Returns a filtered/paginated list of meeting rooms from Inheaden Connect |

All server endpoints require an authenticated Mattermost session (`Mattermost-User-Id` header). Calls to the Inheaden Connect API are authenticated with HTTP Basic Auth using the configured API key.

The plugin stores per-channel meeting room assignments in the Mattermost key–value store (`p.API.KVSet`).

### Webapp (TypeScript / React)

The webapp component registers three UI elements:

- **Channel header button** — opens the meeting room picker modal.
- **`StartMeetingModal`** — displays a searchable list of meeting rooms. On selection it calls `showMeetingPost` (posts to channel) and `startMeeting` (opens the join URL).
- **`StartMeetingPost`** — custom post renderer for posts of type `custom_inco_start_meeting`, showing the meeting name and a join button inline in the channel feed.

For more on plugin webapp extension points, see [Webapp Plugin Reference](https://developers.mattermost.com/extend/plugins/webapp/reference/).

---

## Installation

1. **Build or download** the plugin bundle (`.tar.gz`).  
   To build from source: `make dist`
2. In Mattermost, go to **System Console → Plugins → Plugin Management**.
3. Upload the bundle and enable the plugin.
4. Navigate to the plugin's settings page (see [Configuration](#configuration)) and save your credentials.

For general plugin installation guidance see [Installing Plugins](https://developers.mattermost.com/extend/plugins/using-and-managing-plugins/).

---

## Configuration

Configure the plugin under **System Console → Plugins → Inheaden Connect**.

| Setting | Key | Description | Default |
|---|---|---|---|
| Inheaden Connect URL | `InheadenConnectAPIURL` | Base URL of the Inheaden Connect instance | `https://inco.video` |
| API Key | `APIKey` | API key in `XXX:YYY` format with read access to meeting rooms | *(required)* |

The API key can be created in the Inheaden Connect admin panel. It must have permission to list and join meeting rooms.

---

## Usage

### Starting a Meeting

1. Open a channel.
2. Click the **Inheaden Connect icon** (🎥) in the channel header.
3. Search for and select a meeting room from the list.
4. A message with the join link is posted in the channel. All members can join by clicking **Join Meeting**.

### Slash Commands

| Command | Description |
|---|---|
| `/inco meetingroom <roomId>` | Set a default meeting room for the current channel. The room ID can be found in Inheaden Connect. |
| `/inco meetingroom` | Remove the custom meeting room assignment for the channel (reverts to manual selection). |
| `/inco help` | Display available commands. |

---

## Development

### Prerequisites

- Go 1.16+
- Node.js / npm
- A running Mattermost server for deployment testing

### Build

```bash
# Install webapp dependencies and build everything
make dist

# Run server and webapp tests
make test

# Lint (requires golangci-lint)
make check-style
```

### Deploy to a Local Server

```bash
make deploy
```

This builds the bundle and uploads it to the server configured via `MM_SERVICESETTINGS_SITEURL` and credentials in `env.sh`. See [Deploying Plugins](https://developers.mattermost.com/extend/plugins/deployment/) for details.

### Debugging

Set `MM_DEBUG=1` to build the server and webapp in debug mode:

```bash
MM_DEBUG=1 make dist
```

Use `make attach-headless` to attach `dlv` to a running plugin process on port `2346`.

---

## References

- [Mattermost Plugin Developer Guide](https://developers.mattermost.com/extend/plugins/)
- [Server Plugin Reference](https://developers.mattermost.com/extend/plugins/server/reference/)
- [Webapp Plugin Reference](https://developers.mattermost.com/extend/plugins/webapp/reference/)
- [Inheaden Connect](https://inco.video)
- [Support](https://inco.video/support)
