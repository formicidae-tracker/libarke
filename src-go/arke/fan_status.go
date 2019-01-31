package arke

type FanStatus uint8

const (
	FanOK      FanStatus = 0x00
	FanAging   FanStatus = 0x01
	FanStalled FanStatus = 0x02
)

type FanStatusAndRPM uint16

func (s FanStatusAndRPM) RPM() uint16 {
	return uint16(s & 0xc000)
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
