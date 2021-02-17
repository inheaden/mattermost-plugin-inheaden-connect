import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';

import {id} from './manifest';

class Client {
    private url: string | undefined;

    setServerRoute(url: string) {
        this.url = url + '/plugins/' + id;
    }

    startMeeting = async (
        channelId: any,
        personal = true,
        topic = '',
        meetingId = 0,
        force = false,
    ) => {
        const res = await doPost(`${this.url}/api/v1/startMeeting`, {
            channel_id: channelId,
            personal,
            topic,
            meeting_id: meetingId,
        });
        return res.joinUrl;
    };

    showMeetingPost = async (channelId: any) => {
        await doPost(`${this.url}/api/v1/showMeetingPost`, {
            channel_id: channelId,
        });
    };
}

export const doPost = async (url: string, body: any, headers = {}) => {
    const options = {
        method: 'post',
        body: JSON.stringify(body),
        headers,
    };

    const response = await fetch(url, Client4.getOptions(options));

    if (response.ok) {
        return response.json();
    }

    const text = await response.text();

    throw new ClientError(Client4.url, {
        message: text || '',
        status_code: response.status,
        url,
    });
};

const ClientInstance = new Client();
export default ClientInstance;
