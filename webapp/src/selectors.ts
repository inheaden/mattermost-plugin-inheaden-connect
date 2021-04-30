import { Constants } from "./constants";

const getPluginState = (state: any) =>
  state[`plugins-${Constants.pluginName}`] || {};

export const isStartMeetingModalVisible = (state: any) =>
  getPluginState(state).startMeetingModalVisible || false;
