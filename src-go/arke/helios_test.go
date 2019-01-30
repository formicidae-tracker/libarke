package arke

import (
	. "gopkg.in/check.v1"
)

type HeliosSuite struct{}

var _ = Suite(&HeliosSuite{})

func (s *HeliosSuite) TestSetPoint(c *C) {
	testData := []struct {
		Message HeliosSetPoint
		Buffer  []byte
	}{
		{
			HeliosSetPoint{Visible: 123, UV: 231},
			[]byte{123, 231},
		},
	}

	for _, d := range testData {
		res := HeliosSetPoint{}
		if c.Check(res.Unmarshall(d.Buffer), IsNil) == false {
			continue
		}
		c.Check(res, DeepEquals, d.Message)
	}

	for _, d := range testData {
		res := make([]byte, 2)
		res[0] = 0xff
		n, err := d.Message.Marshall(res)
		if c.Check(err, IsNil) == false {
			continue
		}
		c.Check(n, Equals, 2)
		c.Check(res, DeepEquals, d.Buffer)
	}
	m := HeliosSetPoint{}
	_, err := m.Marshall([]byte{})
	c.Check(err, ErrorMatches, "Invalid buffer size .*")
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, "Invalid buffer size .*")
}
