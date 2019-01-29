package arke

// #include "../../include/arke.h"
import "C"

type FanStatus uint8

const (
	FanOK      FanStatus = 0x00
	FanAging   FanStatus = 0x01
	FanStalled FanStatus = 0x02
)

type FanStatusAndRPM uint16

func (s FanStatusAndRPM) RPM() uint16 {
	return uint16(s & C.ARKE_FAN_RPM_MASK)
}

func (s FanStatusAndRPM) Status() FanStatus {
	if s&C.ARKE_FAN_STALL_ALERT != 0 {
		return FanStalled
	} else if s&C.ARKE_FAN_AGING_ALERT != 0 {
		return FanAging
	} else {
		return FanOK
	}
}
