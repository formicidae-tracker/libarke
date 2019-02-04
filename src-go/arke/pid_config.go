package arke

import "fmt"

type PDConfig struct {
	ProportionnalMultiplier uint8
	DerivativeMultiplier    uint8
	IntegralMultiplier      uint8
	DividerPower            uint8
	DividerPowerIntegral    uint8
}

func (c PDConfig) marshall(buffer []byte) error {
	if c.DividerPower > 15 {
		return fmt.Errorf("Maximal Proportional&Derivative Divider is 15")
	}
	if c.DividerPowerIntegral > 15 {
		return fmt.Errorf("Maximal Integral Divider is 15")
	}

	buffer[0] = c.ProportionnalMultiplier
	buffer[1] = c.DerivativeMultiplier
	buffer[2] = c.IntegralMultiplier
	buffer[3] = ((c.DividerPowerIntegral & 0x0f) << 4) | c.DividerPower&0x0f
	return nil
}

func (c *PDConfig) unmarshall(buffer []byte) {
	c.ProportionnalMultiplier = buffer[0]
	c.DerivativeMultiplier = buffer[1]
	c.IntegralMultiplier = buffer[2]
	c.DividerPower = buffer[3] & 0x0f
	c.DividerPowerIntegral = (buffer[3] & 0xf0) >> 4
}
