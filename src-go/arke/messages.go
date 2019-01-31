package arke

// #include "../../include/arke.h"
import "C"
import (
	"fmt"

	socketcan "github.com/atuleu/golang-socketcan"
)

type MessageType uint16
type MessageClass uint16
type NodeClass uint16
type NodeID uint8

const (
	NetworkControlCommand MessageType = C.ARKE_NETWORK_CONTROL_COMMAND
	HighPriorityMessage   MessageType = C.ARKE_HIGH_PRIORITY_MESSAGE
	StandardMessage       MessageType = C.ARKE_MESSAGE
	HeartBeat             MessageType = C.ARKE_HEARTBEAT
	MessageTypeMask       uint16      = C.ARKE_MESSAGE_TYPE_MASK

	BroadcastClass NodeClass = C.ARKE_BROADCAST
	ZeusClass      NodeClass = C.ARKE_ZEUS
	HeliosClass    NodeClass = C.ARKE_HELIOS
	CelaenoClass   NodeClass = C.ARKE_CELAENO
	NodeClassMask  uint16    = C.ARKE_NODE_CLASS_MASK

	ResetRequest           MessageClass = C.ARKE_RESET_REQUEST
	SynchronisationRequest MessageClass = C.ARKE_SYNCHRONISATION
	HeartBeatRequest       MessageClass = C.ARKE_HEARTBEAT_REQUEST
	IDMask                 uint16       = C.ARKE_SUBID_MASK

	BroadcastID NodeID = 0x00

	ZeusSetPointMessage        MessageClass = C.ARKE_ZEUS_SET_POINT
	ZeusReportMessage          MessageClass = C.ARKE_ZEUS_REPORT
	ZeusVibrationReportMessage MessageClass = C.ARKE_ZEUS_VIBRATION_REPORT
	ZeusConfigMessage          MessageClass = C.ARKE_ZEUS_CONFIG
	ZeusStatusMessage          MessageClass = C.ARKE_ZEUS_STATUS
	ZeusControlPointMessage    MessageClass = C.ARKE_ZEUS_CONTROL_POINT
	HeliosSetPointMessage      MessageClass = C.ARKE_HELIOS_SET_POINT
	HeliosPulseModeMessage     MessageClass = C.ARKE_HELIOS_PULSE_MODE
	CelaenoSetPointMessage     MessageClass = C.ARKE_CELAENO_SET_POINT
	CelaenoStatusMessage       MessageClass = C.ARKE_CELAENO_STATUS
	CelaenoConfigNessage       MessageClass = C.ARKE_CELAENO_CONFIG
)

func makeCANIDT(t MessageType, c MessageClass, n NodeID) uint32 {
	return uint32((uint32(t) << 9) | (uint32(c) << 3) | uint32(n))
}

func extractCANIDT(idt uint32) (t MessageType, c MessageClass, n NodeID) {
	n = NodeID(idt & 0x7)
	c = MessageClass((idt & 0x1f8) >> 3)
	t = MessageType((idt & 0x600) >> 9)
	return
}

type Marshaller interface {
	Marshall([]byte) (int, error)
}

type Unmarshaller interface {
	Unmarshall([]byte) error
}

type identifiable interface {
	MessageClassID() MessageClass
}

type SendableMessage interface {
	Marshaller
	identifiable
}

type ReceivableMessage interface {
	Unmarshaller
	identifiable
}

type Message interface {
	Marshaller
	Unmarshaller
	identifiable
}

func checkID(ID NodeID) error {
	if ID > 7 {
		return fmt.Errorf("Invalid device ID %d (max is 7)", ID)
	}
	return nil
}

func SendMessage(itf *socketcan.RawInterface, m SendableMessage, highPriority bool, ID NodeID) error {
	if err := checkID(ID); err != nil {
		return err
	}
	mType := StandardMessage
	if highPriority == true {
		mType = HighPriorityMessage
	}

	f := socketcan.CanFrame{
		ID:       makeCANIDT(mType, m.MessageClassID(), ID),
		Extended: false,
		RTR:      false,
		Data:     make([]byte, 8),
	}
	dlc, err := m.Marshall(f.Data)
	if err != nil {
		return fmt.Errorf("Could not marshall %v: %s", m, err)
	}
	f.Dlc = uint8(dlc)
	return itf.Send(f)
}

func RequestMessage(itf *socketcan.RawInterface, m ReceivableMessage, ID NodeID) error {
	if err := checkID(ID); err != nil {
		return err
	}
	return itf.Send(socketcan.CanFrame{
		ID:       makeCANIDT(StandardMessage, m.MessageClassID(), ID),
		Extended: false,
		RTR:      false,
		Data:     make([]byte, 0),
		Dlc:      0,
	})
}

type messageCreator func() ReceivableMessage

var messageFactory = make(map[MessageClass]messageCreator)

func ParseMessage(f *socketcan.CanFrame) (ReceivableMessage, NodeID, error) {
	if f.Extended == true {
		return nil, 0, fmt.Errorf("Arke does not support extended can ID")
	}

	if f.RTR == true {
		return nil, 0, fmt.Errorf("RTR frame")
	}

	mType, mClass, mID := extractCANIDT(f.ID)
	if mType == NetworkControlCommand {
		return nil, 0, fmt.Errorf("Network Command")
	}

	if mType == HeartBeat {
		return nil, 0, fmt.Errorf("Heartbeating not implemented yet")
	}

	creator, ok := messageFactory[mClass]
	if ok == false {
		return nil, mID, fmt.Errorf("Unknown message type 0x%05x", mClass)
	}

	m := creator()
	err := m.Unmarshall(f.Data)
	if err != nil {
		err = fmt.Errorf("Could not parse message data: %s", err)
	}

	return m, mID, err
}
