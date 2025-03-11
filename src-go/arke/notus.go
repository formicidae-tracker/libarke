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
	RampDownTime time.Duration
	MinFan       uint8
	MaxHeat      uint8
}

func (m *NotusConfig) MessageClasID() MessageClass {
	return NotusConfigMessage
}

func (m NotusConfig) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 4); err != nil {
		return 0, err
	}
	binary.LittleEndian.AppendUint16(buf[0:], uint16(m.RampDownTime.Milliseconds()))
	buf[2] = m.MinFan
	buf[3] = m.MaxHeat
	return 4, nil
}

func (m *NotusConfig) Unmarshall(buf []byte) error {
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
	messagesName[NotusSetPointMessage] = "Notus.SetPoint"
	messagesName[NotusConfigMessage] = "Notus.Config"
}
