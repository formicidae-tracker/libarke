package arke

import (
	. "gopkg.in/check.v1"
)

type MessageSuite struct{}

var _ = Suite(&MessageSuite{})

func (s *MessageSuite) TestCANIDTIO(c *C) {
	testdata := []struct {
		IDT   uint32
		Type  MessageType
		Class MessageClass
		ID    NodeID
	}{
		{
			0x00,
			NetworkControlCommand,
			MessageClass(BroadcastClass),
			NodeID(ResetRequest),
		},
		{
			0x781,
			HeartBeat,
			MessageClass(CelaenoClass),
			1,
		},
	}

	for _, d := range testdata {
		resType, resClass, resID := extractCANIDT(d.IDT)
		c.Check(resType, Equals, d.Type)
		c.Check(resClass, Equals, d.Class)
		c.Check(resID, Equals, d.ID)
	}

	for _, d := range testdata {
		res := makeCANIDT(d.Type, d.Class, d.ID)
		c.Check(res, Equals, d.IDT)
	}
}
