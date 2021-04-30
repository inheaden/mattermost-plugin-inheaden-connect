import { Action, Store } from "redux";

import { GlobalState } from "mattermost-redux/types/store";

import { getConfig } from "mattermost-redux/selectors/entities/general";

import React from "react";

import { openStandupModal } from "./actions";
import manifest from "./manifest";
import Client from "./client";
import reducer from "./reducers";

// eslint-disable-next-line import/no-unresolved
import { PluginRegistry } from "./types/mattermost-webapp";
import StartMeeting from "./components/startMeetingPost/StartMeetingPost";
import Icon from "./components/Icon";
import StartMeetingModal from "./components/startMeetingModal/StartMeetingModal";

export default class Plugin {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
  public async initialize(
    registry: PluginRegistry,
    store: Store<GlobalState, Action<Record<string, unknown>>>
  ) {
    registry.registerChannelHeaderButtonAction(
      <Icon />,
      (channel: { id: string }) => {
        openStandupModal(channel.id)(store.dispatch);
      },
      "Start Inheaden Connect Meeting"
    );
    Client.setServerRoute(getServerRoute(store.getState()));
    registry.registerPostTypeComponent(
      "custom_inco_start_meeting",
      StartMeeting
    );
    registry.registerRootComponent(StartMeetingModal);

    registry.registerReducer(reducer);
    // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
  }
}

declare global {
  interface Window {
    registerPlugin(id: string, plugin: Plugin): void;
  }
}

window.registerPlugin(manifest.id, new Plugin());

const getServerRoute = (state: any) => {
  const config = getConfig(state);

  let basePath = "";
  if (config && config.SiteURL) {
    basePath = new URL(config.SiteURL).pathname;

    if (basePath && basePath[basePath.length - 1] === "/") {
      basePath = basePath.substr(0, basePath.length - 1);
    }
  }

  return basePath;
};
