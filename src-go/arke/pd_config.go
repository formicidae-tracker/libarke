package arke

type PDConfig struct {
	ProportionnalMultiplier uint8
	DerivativeMultiplier    uint8
	IntegralMultiplier      uint8
	DividerPower            uint8
}

func (c PDConfig) marshall(buffer []byte) error {
	buffer[0] = c.ProportionnalMultiplier
	buffer[1] = c.DerivativeMultiplier
	buffer[2] = c.IntegralMultiplier
	buffer[3] = c.DividerPower
	return nil
}

func (c *PDConfig) unmarshall(buffer []byte) {
	c.ProportionnalMultiplier = buffer[0]
	c.DerivativeMultiplier = buffer[1]
	c.IntegralMultiplier = buf[2]
	c.DividerPower = buf[2]
}
