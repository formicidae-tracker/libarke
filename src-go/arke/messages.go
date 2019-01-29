package arke

// #include "../../include/arke.h"
import "C"

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
	HeliosSetPointMessage      uint16 = C.ARKE_HELIOS_SET_POINT
	HeliosPulseModeMessage     uint16 = C.ARKE_HELIOS_PULSE_MODE
	CelaenoSetPointMessage     uint16 = C.ARKE_CELAENO_SET_POINT
	CelaenoStatusMessage       uint16 = C.ARKE_CELAENO_STATUS
	CelaenoConfigNessage       uint16 = C.ARKE_CELAENO_CONFIG
)

type Marshaller interface {
	Marshall([]byte) uint8
}

type Unmarshaller interface {
	Unmarshall([]byte) error
}
