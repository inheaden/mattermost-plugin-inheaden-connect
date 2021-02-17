import {Action, Store} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import React from 'react';

import {showMeetingMessage} from './actions';
import manifest from './manifest';
import Client from './client';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';
import StartMeeting from './components/StartMeeting';
import Icon from './components/Icon';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(
        registry: PluginRegistry,
        store: Store<GlobalState, Action<Record<string, unknown>>>,
    ) {
        registry.registerChannelHeaderButtonAction(
            <Icon/>,
            (channel: { id: string; }) => {
                showMeetingMessage(channel.id)(store.dispatch, store.getState);
            },
            'Start Inheaden Connect Meeting',
        );
        Client.setServerRoute(getServerRoute(store.getState()));
        registry.registerPostTypeComponent(
            'custom_inco_start_meeting',
            StartMeeting,
        );

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

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath;
};
