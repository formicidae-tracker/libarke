package arke

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"
)

type NotusSuite struct{}

var _ = Suite(&NotusSuite{})

func checkMessageEncoding(c *C, m Message, buffer []byte) {
	builder, ok := messageFactory[m.MessageClassID()]
	c.Assert(ok, Equals, true, Commentf("missing factory"))
	parsed := builder()

	if c.Check(parsed.Unmarshal(buffer), IsNil) == true {
		c.Check(parsed, DeepEquals, m)
	}
	resBuffer := make([]byte, len(buffer))
	n, err := m.Marshal(resBuffer)
	if c.Check(err, IsNil) == false {
		return
	}
	c.Check(n, Equals, len(buffer))
	c.Check(resBuffer, DeepEquals, buffer)

}

func checkMessageLength(c *C, m Message, size int) {
	_, err := m.Marshal([]byte{})
	c.Check(err, ErrorMatches, fmt.Sprintf("Invalid buffer size 0, required: %d", size))
	c.Check(m.Unmarshal([]byte{}), ErrorMatches, fmt.Sprintf("Invalid buffer size 0, required: %d", size))
}

func (s *NotusSuite) TestSetPoint(c *C) {

	checkMessageEncoding(c, &NotusSetPoint{Power: 85}, []byte{0x55})
	checkMessageLength(c, &NotusSetPoint{}, 1)
}

func (s *NotusSuite) TestConfig(c *C) {
	checkMessageEncoding(c, &NotusConfig{
		RampDownTime: 2 * time.Second,
		MinFan:       33,
		MaxHeat:      211,
	}, []byte{0xd0, 0x7, 0x21, 0xd3})
	checkMessageLength(c, &NotusConfig{}, 4)
}
