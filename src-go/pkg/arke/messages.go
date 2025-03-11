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
	IDChangeRequest        MessageClass = 0x02
	ErrorReport            MessageClass = 0x03
	HeartBeatRequest       MessageClass = 0x07
	IDMask                 uint16       = 0x07

	BroadcastID NodeID = 0x00

	RTRRequestMessage             MessageClass = MessageClass(1 << 10)
	HeartBeatMessage              MessageClass = MessageClass(HeartBeat << 9)
	ResetRequestMessage           MessageClass = MessageClass(0x7f8 | ResetRequest)
	SynchronisationRequestMessage MessageClass = MessageClass(0x7f8 | SynchronisationRequest)
	IDChangeRequestMessage        MessageClass = MessageClass(0x7f8 | IDChangeRequest)
	ErrorReportMessage            MessageClass = MessageClass(0x7f8 | ErrorReport)
	HeartBeatRequestMessage       MessageClass = MessageClass(0x7f8 | HeartBeatRequest)
	ZeusSetPointMessage           MessageClass = 0x38
	ZeusReportMessage             MessageClass = 0x39
	ZeusVibrationReportMessage    MessageClass = 0x3a
	ZeusConfigMessage             MessageClass = 0x3b
	ZeusStatusMessage             MessageClass = 0x3c
	ZeusControlPointMessage       MessageClass = 0x3d
	ZeusDeltaTemperatureMessage   MessageClass = 0x3e
	HeliosSetPointMessage         MessageClass = 0x34
	HeliosPulseModeMessage        MessageClass = 0x35
	CelaenoSetPointMessage        MessageClass = 0x30
	CelaenoStatusMessage          MessageClass = 0x31
	CelaenoConfigMessage          MessageClass = 0x32
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
	String() string
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

func SendMessage(itf socketcan.RawInterface, m SendableMessage, highPriority bool, ID NodeID) error {
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

func RequestMessage(itf socketcan.RawInterface, m ReceivableMessage, ID NodeID) error {
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

type networkCommandParser func(c MessageClass, buffer []byte) (ReceivableMessage, NodeID, error)

var networkCommandFactory = make(map[NodeID]networkCommandParser)

type MessageRequestData struct {
	Class MessageClass
	ID    NodeID
}

func (d *MessageRequestData) MessageClassID() MessageClass {
	return RTRRequestMessage
}

var messagesName = make(map[MessageClass]string)

func (c MessageClass) String() string {
	if n, ok := messagesName[c]; ok == true {
		return n
	}
	return "<unknown>"
}

func (d *MessageRequestData) String() string {
	if d.ID == NodeID(0) {
		return fmt.Sprintf("arke.MessageRequest{Message:%s, Node: all}", d.Class)
	}
	return fmt.Sprintf("arke.MessageRequest{Message:%s, Node: %d}", d.Class, d.ID)
}

func (d *MessageRequestData) Unmarshall(buf []byte) error {
	return nil
}

func parseRTR(f *socketcan.CanFrame) (ReceivableMessage, NodeID, error) {
	mType, mClass, mID := extractCANIDT(f.ID)

	if f.Dlc > 0 {
		return nil, 0, fmt.Errorf("RTR frame with a payload")
	}
	if mType != StandardMessage && mType != HighPriorityMessage {
		return nil, 0, fmt.Errorf("Unauthorized network command RTR frame")
	}

	_, ok := messageFactory[mClass]
	if ok == false {
		return nil, mID, fmt.Errorf("Unknown message type 0x%02x", int(mClass))
	}

	return &MessageRequestData{
		Class: mClass,
		ID:    mID,
	}, mID, nil

}

func ParseMessage(f *socketcan.CanFrame) (ReceivableMessage, NodeID, error) {
	if f.Extended == true {
		return nil, 0, fmt.Errorf("Arke does not support extended IDT")
	}

	if f.RTR == true {
		return parseRTR(f)
	}

	mType, mClass, mID := extractCANIDT(f.ID)
	if mType == NetworkControlCommand {
		parser, ok := networkCommandFactory[mID]
		if ok == false {
			return nil, 0, fmt.Errorf("Unknown network command 0x%02x", mID)
		}
		return parser(mClass, f.Data[0:f.Dlc])
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
		return nil, mID, fmt.Errorf("Unknown message type 0x%02x", int(mClass))
	}

	m := creator()
	err := m.Unmarshall(f.Data[0:f.Dlc])
	if err != nil {
		err = fmt.Errorf("Could not parse message data: %s", err)
	}

	return m, mID, err
}
