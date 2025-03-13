package arke

import (
	"time"

	. "gopkg.in/check.v1"
)

type HeliosSuite struct{}

var _ = Suite(&HeliosSuite{})

func (s *HeliosSuite) TestSetPoint(c *C) {
	checkMessageEncoding(c, &HeliosSetPoint{Visible: 123, UV: 231}, []byte{123, 231})
	checkMessageLength(c, &HeliosSetPoint{}, 2)
}

func (s *HeliosSuite) TestPulseMode(c *C) {
	checkMessageEncoding(c, &HeliosPulseMode{Period: 2 * time.Second},
		[]byte{0xd0, 0x07})
	checkMessageLength(c, &HeliosPulseMode{}, 2)
}

func (s *HeliosSuite) TestTriggerMode(c *C) {
	checkMessageEncoding(c, &HeliosTriggerMode{
		Period:      100 * time.Millisecond,
		PulseLength: 3200 * time.Microsecond,
	}, []byte{0xe8, 0x03, 0x80, 0x0c})
	checkMessageLength(c, &HeliosTriggerMode{}, 4)
}
