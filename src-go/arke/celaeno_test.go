package arke

import (
	"time"

	. "gopkg.in/check.v1"
)

type CelaenoSuite struct{}

var _ = Suite(&CelaenoSuite{})

func (s *CelaenoSuite) TestSetPoint(c *C) {

	checkMessageEncoding(c, &CelaenoSetPoint{Power: 127}, []byte{0x7f})
	checkMessageLength(c, &CelaenoSetPoint{}, 1)
}

func (s *CelaenoSuite) TestStatus(c *C) {
	checkMessageEncoding(c, &CelaenoStatus{
		WaterLevel: CelaenoWaterWarning,
		Fan:        1200,
	},
		[]byte{0x01, 0xb0, 0x04},
	)
	checkMessageLength(c, &CelaenoStatus{}, 3)
}

func (s *CelaenoSuite) TestConfig(c *C) {

	checkMessageEncoding(c, &CelaenoConfig{
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
		})

	checkMessageLength(c, &CelaenoConfig{}, 8)
	m := CelaenoConfig{
		RampDownTime: (1 << 16) * time.Millisecond,
	}
	_, err := m.Marshall(make([]byte, 8))
	c.Check(err, ErrorMatches, "Time constant overflow")
}
