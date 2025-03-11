package arke

import (
	"encoding/binary"
	"fmt"
	"time"
)

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

type HeliosPulseMode struct {
	Period time.Duration
}

func (m HeliosPulseMode) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 2); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], uint16(m.Period.Milliseconds()))
	return 2, nil
}

func (m *HeliosPulseMode) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 2); err != nil {
		return err
	}
	m.Period = time.Duration(binary.LittleEndian.Uint16(buf)) * time.Millisecond
	return nil
}

func (m *HeliosPulseMode) MessageClassID() MessageClass {
	return HeliosPulseModeMessage
}

type HeliosTriggerMode struct {
	Period      time.Duration
	PulseLength time.Duration
}

func (m HeliosTriggerMode) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 4); err != nil {
		return 0, err
	}
	binary.LittleEndian.AppendUint16(buf[0:], uint16(m.Period.Microseconds()/100))
	binary.LittleEndian.AppendUint16(buf[2:], uint16(m.PulseLength.Microseconds()))
	return 4, nil
}

func (m *HeliosTriggerMode) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 4); err != nil {
		return err
	}
	m.Period = time.Duration(binary.LittleEndian.Uint16(buf[0:])*100) * time.Microsecond
	m.PulseLength = time.Duration(binary.LittleEndian.Uint16(buf[2:])) * time.Microsecond
	return nil
}

func (m *HeliosTriggerMode) MessageClassID() MessageClass {
	return HeliosTriggerModeMessage
}

func init() {
	messageFactory[HeliosSetPointMessage] = func() ReceivableMessage { return &HeliosSetPoint{} }
	messagesName[HeliosSetPointMessage] = "Helios.SetPoint"
	messagesName[HeliosPulseModeMessage] = "Helios.PulseMode"
	messagesName[HeliosTriggerModeMessage] = "Helios.TriggerMode"
}
