package arke

import "fmt"

type FanStatus uint8

const (
	FanOK      FanStatus = 0x00
	FanAging   FanStatus = 0x01
	FanStalled FanStatus = 0x02
)

type FanStatusAndRPM uint16

func (s FanStatusAndRPM) RPM() uint16 {
	return uint16(s & 0x3fff)
}

func (s FanStatusAndRPM) Status() FanStatus {
	if s&0x8000 != 0 {
		return FanStalled
	} else if s&0x4000 != 0 {
		return FanAging
	} else {
		return FanOK
	}
}

func (s FanStatus) String() string {
	if s == FanOK {
		return "OK"
	}
	if s == FanAging {
		return "Aging"
	}
	return "Stalled"
}

func (s FanStatusAndRPM) String() string {
	return fmt.Sprintf("{Status: %s, RPM: %d}", s.Status(), s.RPM())
}
