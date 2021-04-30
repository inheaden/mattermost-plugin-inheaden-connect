import manifest from "manifest";

export const Constants = {
  actions: {
    showStartMeetingModal: "SHOW_START_MEETING_MODAL",
    closeStartMeetingModal: "CLOSE_START_MEETING_MODAL",
  },
  pluginName: manifest.id,
  pluginDisplayName: manifest.name,
};
