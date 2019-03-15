package arke

import "fmt"

type HeliosSetPoint struct {
	Visible uint8
	UV      uint8
}

func (c HeliosSetPoint) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 2); err != nil {
		return 0, err
	}
	buf[0] = c.Visible
	buf[1] = c.UV
	return 2, nil
}

func (c *HeliosSetPoint) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 2); err != nil {
		return err
	}
	c.Visible = buf[0]
	c.UV = buf[1]
	return nil
}

func (m *HeliosSetPoint) MessageClassID() MessageClass {
	return HeliosSetPointMessage
}

func (c *HeliosSetPoint) String() string {
	return fmt.Sprintf("Helios.SetPoint{Visible: %d, UV: %d}", c.Visible, c.UV)
}

func init() {
	messageFactory[HeliosSetPointMessage] = func() ReceivableMessage { return &HeliosSetPoint{} }
}
