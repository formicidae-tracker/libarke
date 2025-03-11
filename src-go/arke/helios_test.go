package arke

import (
	. "gopkg.in/check.v1"
)

type HeliosSuite struct{}

var _ = Suite(&HeliosSuite{})

func (s *HeliosSuite) TestSetPoint(c *C) {
	checkMessageEncoding(c, &HeliosSetPoint{Visible: 123, UV: 231}, []byte{123, 231})
	checkMessageLength(c, &HeliosSetPoint{}, 2)
}
