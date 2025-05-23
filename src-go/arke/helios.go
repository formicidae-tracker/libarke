package arke

import (
	"encoding/binary"
	"fmt"
	"time"
)

type HeliosSetPoint struct {
	Visible uint8 `positional-arg-name:"visible" required:"yes"`
	UV      uint8 `positional-arg-name:"UV" required:"yes"`
}

func (m HeliosSetPoint) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 2); err != nil {
		return 0, err
	}
	buf[0] = m.Visible
	buf[1] = m.UV
	return 2, nil
}

func (m *HeliosSetPoint) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 2); err != nil {
		return err
	}
	m.Visible = buf[0]
	m.UV = buf[1]
	return nil
}

func (m *HeliosSetPoint) MessageClassID() MessageClass {
	return HeliosSetPointMessage
}

func (m *HeliosSetPoint) String() string {
	return fmt.Sprintf("Helios.SetPoint{Visible: %d, UV: %d}", m.Visible, m.UV)
}

type HeliosPulseMode struct {
	Period time.Duration `positional-arg-name:"period" required:"yes"`
}

func (m HeliosPulseMode) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 2); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], uint16(m.Period.Milliseconds()))
	return 2, nil
}

func (m *HeliosPulseMode) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 2); err != nil {
		return err
	}
	m.Period = time.Duration(binary.LittleEndian.Uint16(buf)) * time.Millisecond
	return nil
}

func (m *HeliosPulseMode) MessageClassID() MessageClass {
	return HeliosPulseModeMessage
}

func (m *HeliosPulseMode) String() string {
	return fmt.Sprintf("Helios.PulseMode{Period: %s}", m.Period)
}

type HeliosTriggerMode struct {
	Period      time.Duration `positional-arg-name:"period" required:"yes"`
	PulseLength time.Duration `positional-arg-name:"length" required:"yes"`
	CameraDelay time.Duration `positional-arg-name:"delay" required:"no" default:"0s"`
}

func (m HeliosTriggerMode) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 6); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], uint16(m.Period.Microseconds()/100))
	binary.LittleEndian.PutUint16(buf[2:], uint16(m.PulseLength.Microseconds()))
	binary.LittleEndian.PutUint16(buf[4:], uint16(m.CameraDelay.Microseconds()))
	return 6, nil
}

func (m *HeliosTriggerMode) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 6); err != nil {
		return err
	}
	m.Period = time.Duration(binary.LittleEndian.Uint16(buf[0:])) * 100 * time.Microsecond
	m.PulseLength = time.Duration(binary.LittleEndian.Uint16(buf[2:])) * time.Microsecond
	m.CameraDelay = time.Duration(int16(binary.LittleEndian.Uint16(buf[4:]))) * time.Microsecond
	return nil
}

func (m *HeliosTriggerMode) MessageClassID() MessageClass {
	return HeliosTriggerModeMessage
}

func (m *HeliosTriggerMode) String() string {
	return fmt.Sprintf("Helios.TriggerMode{Period: %s, PulseLength: %s, CameraDelay: %s}", m.Period, m.PulseLength, m.CameraDelay)
}

func init() {
	messageFactory[HeliosSetPointMessage] = func() Message { return &HeliosSetPoint{} }
	messagesName[HeliosSetPointMessage] = "Helios.SetPoint"
	messageFactory[HeliosPulseModeMessage] = func() Message { return &HeliosPulseMode{} }
	messagesName[HeliosPulseModeMessage] = "Helios.PulseMode"
	messageFactory[HeliosTriggerModeMessage] = func() Message { return &HeliosTriggerMode{} }
	messagesName[HeliosTriggerModeMessage] = "Helios.TriggerMode"
}
