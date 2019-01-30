package arke

import (
	"time"

	. "gopkg.in/check.v1"
)

type CelaenoSuite struct{}

var _ = Suite(&CelaenoSuite{})

func (s *CelaenoSuite) TestSetPoint(c *C) {
	testData := []struct {
		Message CelaenoSetPoint
		Buffer  []byte
	}{
		{
			CelaenoSetPoint{Power: 127},
			[]byte{0x7f},
		},
	}

	for _, d := range testData {
		res := CelaenoSetPoint{}
		if c.Check(res.Unmarshall(d.Buffer), IsNil) == false {
			continue
		}
		c.Check(res, DeepEquals, d.Message)
	}

	for _, d := range testData {
		res := make([]byte, 1)
		res[0] = 0xff
		n, err := d.Message.Marshall(res)
		if c.Check(err, IsNil) == false {
			continue
		}
		c.Check(n, Equals, 1)
		c.Check(res, DeepEquals, d.Buffer)
	}
	m := CelaenoSetPoint{}
	_, err := m.Marshall([]byte{})
	c.Check(err, ErrorMatches, "Invalid buffer size .*")
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, "Invalid buffer size .*")
}

func (s *CelaenoSuite) TestStatus(c *C) {
	testData := []struct {
		Message CelaenoStatus
		Buffer  []byte
	}{
		{
			CelaenoStatus{
				WaterLevel: 0x01,
				Fan:        1200,
			},
			[]byte{0x01, 0xb0, 0x04},
		},
	}

	for _, d := range testData {
		res := CelaenoStatus{}
		if c.Check(res.Unmarshall(d.Buffer), IsNil) == false {
			continue
		}
		c.Check(res, DeepEquals, d.Message)
	}

	m := CelaenoStatus{}
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, "Invalid buffer size .*")
}

func (s *CelaenoSuite) TestConfig(c *C) {
	testData := []struct {
		Message CelaenoConfig
		Buffer  []byte
	}{
		{
			CelaenoConfig{
				200 * time.Millisecond,
				300 * time.Millisecond,
				400 * time.Millisecond,
				500 * time.Millisecond,
			},
			[]byte{
				0xc8, 0x00,
				0x2c, 0x01,
				0x90, 0x01,
				0xf4, 0x01,
			},
		},
	}

	for _, d := range testData {
		res := CelaenoConfig{}
		if c.Check(res.Unmarshall(d.Buffer), IsNil) == false {
			continue
		}
		c.Check(res, DeepEquals, d.Message)
	}

	for _, d := range testData {
		res := make([]byte, 8)
		res[0] = 0xff
		n, err := d.Message.Marshall(res)
		if c.Check(err, IsNil) == false {
			continue
		}
		c.Check(n, Equals, 8)
		c.Check(res, DeepEquals, d.Buffer)
	}
	m := CelaenoConfig{}
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, "Invalid buffer size .*")

	errorData := []struct {
		Message CelaenoConfig
		Buffer  []byte
		Ematch  string
	}{
		{
			CelaenoConfig{},
			make([]byte, 0),
			"Invalid buffer size .*",
		},

		{
			CelaenoConfig{
				RampDownTime: (1 << 16) * time.Millisecond,
			},
			make([]byte, 8),
			"Time constant overflow",
		},
	}

	for _, d := range errorData {
		_, err := d.Message.Marshall(d.Buffer)
		c.Check(err, ErrorMatches, d.Ematch)
	}

}
