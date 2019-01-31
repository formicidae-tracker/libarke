package arke

import (
	"fmt"

	socketcan "github.com/atuleu/golang-socketcan"
)

type MessageType uint16
type MessageClass uint16
type NodeClass uint16
type NodeID uint8

const (
	NetworkControlCommand MessageType = 0x00
	HighPriorityMessage   MessageType = 0x01
	StandardMessage       MessageType = 0x02
	HeartBeat             MessageType = 0x03
	MessageTypeMask       uint16      = 0x03 << 9

	BroadcastClass NodeClass = 0x0
	ZeusClass      NodeClass = 0x38
	HeliosClass    NodeClass = 0x34
	CelaenoClass   NodeClass = 0x30
	NodeClassMask  uint16    = 0x3f << 3

	ResetRequest           MessageClass = 0x00
	SynchronisationRequest MessageClass = 0x01
	HeartBeatRequest       MessageClass = 0x07
	IDMask                 uint16       = 0x07

	BroadcastID NodeID = 0x00

	HeartBeatMessage           MessageClass = MessageClass(HeartBeat << 9)
	ZeusSetPointMessage        MessageClass = 0x38
	ZeusReportMessage          MessageClass = 0x39
	ZeusVibrationReportMessage MessageClass = 0x3a
	ZeusConfigMessage          MessageClass = 0x3b
	ZeusStatusMessage          MessageClass = 0x3c
	ZeusControlPointMessage    MessageClass = 0x3d
	HeliosSetPointMessage      MessageClass = 0x34
	HeliosPulseModeMessage     MessageClass = 0x35
	CelaenoSetPointMessage     MessageClass = 0x30
	CelaenoStatusMessage       MessageClass = 0x31
	CelaenoConfigNessage       MessageClass = 0x32
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
		RTR:      true,
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
		return nil, 0, fmt.Errorf("Network Command Received")
	}

	if mType == HeartBeat {
		res := &HeartBeatData{}
		if err := res.Unmarshall(f.Data[0:f.Dlc]); err != nil {
			return nil, mID, err
		}
		res.Class = NodeClass(mClass)
		res.ID = mID
		return res, mID, nil
	}

	creator, ok := messageFactory[mClass]
	if ok == false {
		return nil, mID, fmt.Errorf("Unknown message type 0x%05x", mClass)
	}

	m := creator()
	err := m.Unmarshall(f.Data[0:f.Dlc])
	if err != nil {
		err = fmt.Errorf("Could not parse message data: %s", err)
	}

	return m, mID, err
}
