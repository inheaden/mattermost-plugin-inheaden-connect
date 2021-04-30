import { Constants } from "./constants";
import { combineReducers } from "redux";

let prevState = false;

export const startMeetingModalVisible = (state = false, action) => {
  switch (action.type) {
    case Constants.actions.showStartMeetingModal:
      prevState = true;
      return true;
    case Constants.actions.closeStartMeetingModal:
      prevState = false;
      return false;
    default:
      return prevState || false;
  }
};

export default combineReducers({
  startMeetingModalVisible,
});
