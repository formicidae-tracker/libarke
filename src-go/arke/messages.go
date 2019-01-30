package arke

// #include "../../include/arke.h"
import "C"
import (
	"fmt"

	socketcan "github.com/atuleu/golang-socketcan"
)

const (
	NetworkControlCommand uint16 = C.ARKE_NETWORK_CONTROL_COMMAND
	HighPriorityMessage   uint16 = C.ARKE_HIGH_PRIORITY_MESSAGE
	StandardMessage       uint16 = C.ARKE_MESSAGE
	HeartBeat             uint16 = C.ARKE_HEARTBEAT
	MessageTypeMask       uint16 = C.ARKE_MESSAGE_TYPE_MASK

	BroadcastClass uint16 = C.ARKE_BROADCAST
	ZeusClass      uint16 = C.ARKE_ZEUS
	HeliosClass    uint16 = C.ARKE_HELIOS
	CelaenoClass   uint16 = C.ARKE_CELAENO
	NodeClassMask  uint16 = C.ARKE_NODE_CLASS_MASK

	ResetRequest           uint16 = C.ARKE_RESET_REQUEST
	SynchronisationRequest uint16 = C.ARKE_SYNCHRONISATION
	HeartBeatRequest       uint16 = C.ARKE_HEARTBEAT_REQUEST
	IDMask                 uint16 = C.ARKE_SUBID_MASK

	ZeusSetPointMessage        uint16 = C.ARKE_ZEUS_SET_POINT
	ZeusReportMessage          uint16 = C.ARKE_ZEUS_REPORT
	ZeusVibrationReportMessage uint16 = C.ARKE_ZEUS_VIBRATION_REPORT
	ZeusConfigMessage          uint16 = C.ARKE_ZEUS_CONFIG
	ZeusStatusMessage          uint16 = C.ARKE_ZEUS_STATUS
	ZeusControlPointMessage    uint16 = C.ARKE_ZEUS_CONTROL_POINT
	HeliosSetPointMessage      uint16 = C.ARKE_HELIOS_SET_POINT
	HeliosPulseModeMessage     uint16 = C.ARKE_HELIOS_PULSE_MODE
	CelaenoSetPointMessage     uint16 = C.ARKE_CELAENO_SET_POINT
	CelaenoStatusMessage       uint16 = C.ARKE_CELAENO_STATUS
	CelaenoConfigNessage       uint16 = C.ARKE_CELAENO_CONFIG
)

type Marshaller interface {
	Marshall([]byte) (int, error)
}

type Unmarshaller interface {
	Unmarshall([]byte) error
}

type identifiable interface {
	MessageClassID() uint16
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

func checkID(ID uint8) error {
	if ID > 7 {
		return fmt.Errorf("Invalid device ID %d (max is 7)", ID)
	}
	return nil
}

func SendMessage(itf *socketcan.RawInterface, m SendableMessage, highPriority bool, ID uint8) error {
	if err := checkID(ID); err != nil {
		return err
	}
	canID := (m.MessageClassID() << 3) | uint16(ID)
	if highPriority == true {
		canID |= HighPriorityMessage << 9
	} else {
		canID |= StandardMessage << 9
	}

	f := socketcan.CanFrame{
		ID:       uint32(canID),
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

func RequestMessage(itf *socketcan.RawInterface, m ReceivableMessage, ID uint8) error {
	if err := checkID(ID); err != nil {
		return err
	}
	return itf.Send(socketcan.CanFrame{
		ID:       uint32((StandardMessage << 9) | (m.MessageClassID() << 3) | uint16(ID)),
		Extended: false,
		RTR:      false,
		Data:     make([]byte, 0),
		Dlc:      0,
	})
}

type messageCreator func() ReceivableMessage

var messageFactory = make(map[uint16]messageCreator)

func ParseMessage(f *socketcan.CanFrame) (ReceivableMessage, uint8, error) {
	if f.Extended == true {
		return nil, 0, fmt.Errorf("Arke does not support extended can ID")
	}

	if f.RTR == true {
		return nil, 0, fmt.Errorf("RTR frame")
	}

	mType := uint16((f.ID & 0x7ff) >> 9)
	if mType == NetworkControlCommand {
		return nil, 0, fmt.Errorf("Network Command")
	}

	if mType == HeartBeat {
		return nil, 0, fmt.Errorf("Heartbeating not implemented yet")
	}

	mID := uint8(f.ID & 0x7)
	mClass := uint16((f.ID & 0x1f8) >> 3)
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
