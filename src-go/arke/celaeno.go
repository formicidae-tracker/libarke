package arke

import (
	"encoding/binary"
	"fmt"
	"time"
)

type CelaenoSetPoint struct {
	Power uint8
}

func (m *CelaenoSetPoint) MessageClassID() MessageClass {
	return CelaenoSetPointMessage
}

func (m CelaenoSetPoint) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 1); err != nil {
		return 0, err
	}
	buf[0] = m.Power
	return 1, nil
}

func (m *CelaenoSetPoint) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 1); err != nil {
		return err
	}
	m.Power = buf[0]
	return nil
}

func (m *CelaenoSetPoint) String() string {
	return fmt.Sprintf("Celaeno.SetPoint{Power: %d}", m.Power)
}

type WaterLevelStatus uint8

const (
	CelaenoWaterNominal   WaterLevelStatus = 0x00
	CelaenoWaterWarning   WaterLevelStatus = 0x01
	CelaenoWaterCritical  WaterLevelStatus = 0x02
	CelaenoWaterReadError WaterLevelStatus = 0x04
)

type CelaenoStatus struct {
	WaterLevel WaterLevelStatus
	Fan        FanStatusAndRPM
}

func (m *CelaenoStatus) MessageClassID() MessageClass {
	return CelaenoStatusMessage
}

func (m *CelaenoStatus) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 3); err != nil {
		return 0, err
	}

	buf[0] = byte(m.WaterLevel)
	binary.LittleEndian.PutUint16(buf[1:], uint16(m.Fan))

	return 3, nil
}

func (m *CelaenoStatus) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 3); err != nil {
		return err
	}
	if buf[0]&0x04 != 0 {
		m.WaterLevel = CelaenoWaterReadError
	} else if buf[0]&0x02 != 0 {
		m.WaterLevel = CelaenoWaterCritical
	} else {
		m.WaterLevel = WaterLevelStatus(buf[0])
	}

	m.Fan = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[1:]))
	return nil
}

func (s WaterLevelStatus) String() string {
	if s == CelaenoWaterReadError {
		return "readout-error"
	}
	prefix := ""
	if s&CelaenoWaterReadError != 0 {
		prefix = "readout-error|"
	}
	if s&CelaenoWaterCritical != 0 {
		return prefix + "critical"
	}
	if s&CelaenoWaterWarning != 0 {
		return prefix + "warning"
	}
	return prefix + "nominal"
}

func (m *CelaenoStatus) String() string {
	return fmt.Sprintf("Celaeno.Status{WaterLevel: %s, Fan:%s}", m.WaterLevel, m.Fan)
}

type CelaenoConfig struct {
	RampUpTime    time.Duration
	RampDownTime  time.Duration
	MinimumOnTime time.Duration
	DebounceTime  time.Duration
}

func (m *CelaenoConfig) MessageClassID() MessageClass {
	return CelaenoConfigMessage
}

const MaxUint16 = ^uint16(0)

func castDuration(t time.Duration) (uint16, error) {
	res := t.Nanoseconds() / 1000000
	if res > int64(MaxUint16) {
		return 0xffff, fmt.Errorf("Time constant overflow")
	}
	return uint16(res), nil
}

func (m CelaenoConfig) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 8); err != nil {
		return 0, err
	}
	for i, t := range []time.Duration{m.RampUpTime, m.RampDownTime, m.MinimumOnTime, m.DebounceTime} {
		data, err := castDuration(t)
		if err != nil {
			return 2 * i, err
		}
		binary.LittleEndian.PutUint16(buf[(2*i):], data)
	}
	return 8, nil
}

func (m *CelaenoConfig) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 8); err != nil {
		return err
	}
	m.RampUpTime = time.Duration(binary.LittleEndian.Uint16(buf[0:])) * time.Millisecond
	m.RampDownTime = time.Duration(binary.LittleEndian.Uint16(buf[2:])) * time.Millisecond
	m.MinimumOnTime = time.Duration(binary.LittleEndian.Uint16(buf[4:])) * time.Millisecond
	m.DebounceTime = time.Duration(binary.LittleEndian.Uint16(buf[6:])) * time.Millisecond
	return nil
}

func (m *CelaenoConfig) String() string {
	return fmt.Sprintf("Celaeno.Config{RampUp: %s, RampDown: %s, MinimumOn: %s, Debounce: %s}",
		m.RampUpTime,
		m.RampDownTime,
		m.MinimumOnTime,
		m.DebounceTime)
}

func init() {
	messageFactory[CelaenoSetPointMessage] = func() Message { return &CelaenoSetPoint{} }
	messagesName[CelaenoSetPointMessage] = "Celaeno.SetPoint"
	messageFactory[CelaenoStatusMessage] = func() Message { return &CelaenoStatus{} }
	messagesName[CelaenoStatusMessage] = "Celaeno.Status"
	messageFactory[CelaenoConfigMessage] = func() Message { return &CelaenoConfig{} }
	messagesName[CelaenoConfigMessage] = "Celaeno.Config"
}
