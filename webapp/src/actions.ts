import { PostTypes } from "mattermost-redux/action_types";

import Client from "./client";

import { Constants } from "./constants";

function handleError(
  error: any,
  getState: () => any,
  channelId: string,
  dispatch: any
) {
  let m = error.message;
  // eslint-disable-next-line no-console
  console.error(error);
  if (error.message && error.message[0] === "{") {
    const e = JSON.parse(error.message);

    // Error is from Zoom API
    if (e && e.message) {
      m = "Inheaden Connect error: " + e.message;
    }
  }

  const post = {
    id: "incoPlugin" + Date.now(),
    create_at: Date.now(),
    update_at: 0,
    edit_at: 0,
    delete_at: 0,
    is_pinned: false,
    user_id: getState().entities.users.currentUserId,
    channel_id: channelId,
    root_id: "",
    parent_id: "",
    original_id: "",
    message: m,
    type: "system_ephemeral",
    props: {},
    hashtags: "",
    pending_post_id: "",
  };

  dispatch({
    type: PostTypes.RECEIVED_NEW_POST,
    data: post,
    channelId,
  });

  return { error };
}

export function showMeetingMessage(channelId: string, room: any) {
  return async (dispatch: any, getState: () => any) => {
    try {
      const startFunction = Client.showMeetingPost;
      await startFunction(channelId, room);
    } catch (error) {
      return handleError(error, getState, channelId, dispatch);
    }

    return { data: true };
  };
}

export function startMeeting(channelId: string, roomId: string) {
  return async (dispatch: any, getState: () => any) => {
    try {
      const startFunction = Client.startMeeting;
      const meetingURL = await startFunction(channelId, roomId);
      if (meetingURL) {
        window.open(meetingURL);
      }
    } catch (error) {
      handleError(error, getState, channelId, dispatch);
    }

    return { data: true };
  };
}

export const openStandupModal = (channelId: string) => (dispatch) => {
  dispatch({
    type: Constants.actions.showStartMeetingModal,
    channelId,
  });
};

export const closeStandupModal = (channelId: string) => (dispatch) => {
  dispatch({
    type: Constants.actions.closeStartMeetingModal,
    channelId,
  });
};
