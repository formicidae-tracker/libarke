package arke

import "fmt"

type PDConfig struct {
	DeadRegion              uint8
	ProportionnalMultiplier uint8
	DerivativeMultiplier    uint8
	ProportionalDivider     uint8
	DerivativeDivider       uint8
}

func (c PDConfig) marshall(buffer []byte) error {
	if c.ProportionalDivider > 15 {
		return fmt.Errorf("Maximal Proportional Divider is 15")
	}
	if c.DerivativeDivider > 15 {
		return fmt.Errorf("Maximal Derivative Divider is 15")
	}

	buffer[0] = c.DeadRegion
	buffer[1] = c.ProportionnalMultiplier
	buffer[2] = c.DerivativeMultiplier
	buffer[3] = ((c.DerivativeDivider & 0x0f) << 4) | c.ProportionalDivider&0x0f
	return nil
}

func (c PDConfig) unmarshall(buffer []byte) {
	c.DeadRegion = buffer[0]
	c.ProportionnalMultiplier = buffer[1]
	c.DerivativeMultiplier = buffer[2]
	c.ProportionalDivider = buffer[3] & 0x0f
	c.DerivativeDivider = (buffer[3] & 0xf0) >> 4
}
