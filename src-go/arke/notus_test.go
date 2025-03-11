package arke

import (
	"fmt"

	. "gopkg.in/check.v1"
)

type NotusSuite struct{}

var _ = Suite(&NotusSuite{})

func checkMessageEncoding(c *C, m Message, buffer []byte) {
	parsed := messageFactory[m.MessageClassID()]()

	if c.Check(parsed.Unmarshall(buffer), IsNil) == true {
		c.Check(parsed, DeepEquals, m)
	}
	resBuffer := make([]byte, len(buffer))
	n, err := m.Marshall(resBuffer)
	if c.Check(err, IsNil) == false {
		return
	}
	c.Check(n, Equals, len(buffer))
	c.Check(resBuffer, DeepEquals, buffer)

}

func checkMessageLength(c *C, m Message, size int) {
	_, err := m.Marshall([]byte{})
	c.Check(err, ErrorMatches, fmt.Sprintf("Invalid buffer size 0, required: %d", size))
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, fmt.Sprintf("Invalid buffer size 0, required: %d", size))
}

func (s *NotusSuite) TestSetPoint(c *C) {

	checkMessageEncoding(c, &NotusSetPoint{Power: 85}, []byte{0x55})
	checkMessageLength(c, &NotusSetPoint{}, 1)
}

func (s *NotusSuite) TestConfig(c *C) {

}
