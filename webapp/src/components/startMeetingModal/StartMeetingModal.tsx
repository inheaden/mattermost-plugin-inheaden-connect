import React, { useEffect, useState } from "react";
import { isStartMeetingModalVisible } from "../../selectors";
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import {
  closeStandupModal,
  showMeetingMessage,
  startMeeting,
} from "../../actions";
import {
  Alert,
  Button,
  FormControl,
  FormGroup,
  InputGroup,
  Modal,
  OverlayTrigger,
  Tooltip,
} from "react-bootstrap";
import { Constants } from "../../constants";
import Client from "../../client";
import "./style.css";

const mapStateToProps = (state: any) => ({
  currentUserId: state.entities.users.currentUserId,
  channelId: state.entities.channels.currentChannelId,
  visible: isStartMeetingModalVisible(state),
});

type mapProps = ReturnType<typeof mapStateToProps> &
  ReturnType<typeof mapDispatchToProps>;

const mapDispatchToProps = (dispatch: any) =>
  bindActionCreators(
    {
      close: closeStandupModal,
      showMeeting: showMeetingMessage,
      startMeeting: startMeeting,
    },
    dispatch
  );
export interface Props extends mapProps {
  theme: any;
}

const StartMeetingModal = ({
  close,
  currentUserId,
  channelId,
  visible,
  showMeeting,
  startMeeting,
}: Props) => {
  const [meetingRooms, setMeetingRooms] = useState([] as any[]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (meetingRooms.length === 0 && visible) {
      setIsLoading(true);
      Client.getAllMeetingRooms()
        .then((res) => {
          console.log(res);
          setMeetingRooms(res);
        })
        .catch(console.error)
        .finally(() => setIsLoading(false));
    }
  }, [visible]);

  return (
    <Modal show={visible} onHide={close}>
      <Modal.Header closeButton={true}>
        <Modal.Title>{Constants.pluginDisplayName}</Modal.Title>
      </Modal.Header>
      <div
        style={{
          padding: "10px",
        }}
      >
        <p>Select a meeting room from the list below:</p>
        {isLoading && "loading ..."}
        {!isLoading && (
          <div
            style={{
              display: "grid",
              gap: "10px",
            }}
          >
            {meetingRooms.length ? (
              meetingRooms.map((room, idx) => (
                <MeetingRoomItem
                  room={room}
                  key={idx}
                  onClick={() => {
                    close(channelId);
                    showMeeting(channelId, room);
                    startMeeting(channelId, room.id);
                  }}
                />
              ))
            ) : (
              <div style={{ textAlign: "center" }}>No meeting rooms found</div>
            )}
          </div>
        )}
      </div>
    </Modal>
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(StartMeetingModal);

const MeetingRoomItem = ({
  room,
  onClick,
}: {
  room: any;
  onClick: () => void;
}) => {
  return (
    <div
      style={{
        padding: "10px",
        border: "1px solid grey",
        borderRadius: "5px",
        cursor: "pointer",
      }}
      className="meeting-room-item"
      onClick={onClick}
    >
      <strong style={{ display: "block" }}>{room.name}</strong>
      <small>Max participants: {room.type?.maxParticipants}</small>
    </div>
  );
};
