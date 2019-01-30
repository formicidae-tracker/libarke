package arke

import (
	. "gopkg.in/check.v1"
)

type FanStatusSuite struct{}

var _ = Suite(&FanStatusSuite{})

func (s *FanStatusSuite) TestFanStatus(c *C) {
	testData := []struct {
		Binary uint16
		RPM    uint16
		Status FanStatus
	}{
		{0, 0, FanOK},
		{1200, 1200, FanOK},
		{800 | (uint16(FanAging) << 14), 800, FanAging},
		{0 | (uint16(FanStalled) << 14), 0, FanStalled},
	}

	for _, d := range testData {
		fan := FanStatusAndRPM(d.Binary)
		c.Check(fan.RPM(), Equals, d.RPM)
		c.Check(fan.Status(), Equals, d.Status)
	}

}
