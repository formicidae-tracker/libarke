package arke

import (
	"encoding/binary"
	"fmt"
	"time"
)

type NotusSetPoint struct {
	Power uint8
}

func (m *NotusSetPoint) MessageClassID() MessageClass {
	return NotusSetPointMessage
}

func (m NotusSetPoint) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 1); err != nil {
		return 0, err
	}
	buf[0] = m.Power
	return 1, nil
}

func (m *NotusSetPoint) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 1); err != nil {
		return err
	}
	m.Power = buf[0]
	return nil
}

func (m *NotusSetPoint) String() string {
	return fmt.Sprintf("Notus.SetPoint{Power: %d}", m.Power)
}

type NotusConfig struct {
	RampDownTime time.Duration `long:"ramp-down" description:"time to keep fan off on poweroff" default:"2s"`
	MinFan       uint8         `long:"min-fan" description:"minimum fan power (0-255)" default:"50"`
	MaxHeat      uint8         `long:"max-fan" description:"maximum heat power (0-255)" default:"200"`
}

func (m *NotusConfig) MessageClassID() MessageClass {
	return NotusConfigMessage
}

func (m NotusConfig) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 4); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], uint16(m.RampDownTime.Milliseconds()))
	buf[2] = m.MinFan
	buf[3] = m.MaxHeat
	return 4, nil
}

func (m *NotusConfig) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 4); err != nil {
		return err
	}
	m.RampDownTime = time.Duration(binary.LittleEndian.Uint16(buf[0:])) * time.Millisecond
	m.MinFan = buf[2]
	m.MaxHeat = buf[3]
	return nil
}

func (m *NotusConfig) String() string {
	return fmt.Sprintf("Notus.Config{RampDownTime: %s, MinFan: %d, MaxHeat: %d}",
		m.RampDownTime, m.MinFan, m.MaxHeat)
}

func init() {
	messageFactory[NotusSetPointMessage] = func() Message { return &NotusSetPoint{} }
	messageFactory[NotusConfigMessage] = func() Message { return &NotusConfig{} }
	messagesName[NotusSetPointMessage] = "Notus.SetPoint"
	messagesName[NotusConfigMessage] = "Notus.Config"
}
