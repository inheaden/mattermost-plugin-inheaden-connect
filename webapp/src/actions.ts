import { PostTypes } from "mattermost-redux/action_types";

import Client from "./client";

export function showMeetingMessage(channelId) {
  return async (dispatch, getState) => {
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
      message: "",
      type: "custom_inco_start_meeting",
      props: {},
      hashtags: "",
      pending_post_id: "",
    };

    dispatch({
      type: PostTypes.RECEIVED_NEW_POST,
      data: post,
      channelId,
    });

    return { error: null };
  };
}

export function startMeeting(channelId) {
  return async (dispatch, getState) => {
    try {
      const startFunction = Client.startMeeting;
      const meetingURL = await startFunction(channelId, true);
      if (meetingURL) {
        window.open(meetingURL);
      }
    } catch (error) {
      let m = error.message;
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

    return { data: true };
  };
}
