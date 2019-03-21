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

func (c CelaenoSetPoint) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 1); err != nil {
		return 0, err
	}
	buf[0] = c.Power
	return 1, nil
}

func (c *CelaenoSetPoint) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 1); err != nil {
		return err
	}
	c.Power = buf[0]
	return nil
}

func (c *CelaenoSetPoint) String() string {
	return fmt.Sprintf("Celaeno.SetPoint{Power: %d}", c.Power)
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

func (c *CelaenoStatus) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 3); err != nil {
		return err
	}
	if buf[0]&0x04 != 0 {
		c.WaterLevel = CelaenoWaterReadError
	} else if buf[0]&0x02 != 0 {
		c.WaterLevel = CelaenoWaterCritical
	} else {
		c.WaterLevel = WaterLevelStatus(buf[0])
	}

	c.Fan = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[1:]))
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

func (c *CelaenoStatus) String() string {
	return fmt.Sprintf("Celaeno.Status{WaterLevel: %s, Fan:%s}", c.WaterLevel, c.Fan)
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

func (c CelaenoConfig) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 8); err != nil {
		return 0, err
	}
	for i, t := range []time.Duration{c.RampUpTime, c.RampDownTime, c.MinimumOnTime, c.DebounceTime} {
		data, err := castDuration(t)
		if err != nil {
			return 2 * i, err
		}
		binary.LittleEndian.PutUint16(buf[(2*i):], data)
	}
	return 8, nil
}

func (c *CelaenoConfig) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 8); err != nil {
		return err
	}
	c.RampUpTime = time.Duration(binary.LittleEndian.Uint16(buf[0:])) * time.Millisecond
	c.RampDownTime = time.Duration(binary.LittleEndian.Uint16(buf[2:])) * time.Millisecond
	c.MinimumOnTime = time.Duration(binary.LittleEndian.Uint16(buf[4:])) * time.Millisecond
	c.DebounceTime = time.Duration(binary.LittleEndian.Uint16(buf[6:])) * time.Millisecond
	return nil
}

func (c *CelaenoConfig) String() string {
	return fmt.Sprintf("Celaeno.Config{RampUp: %s, RampDown: %s, MinimumOn: %s, Debounce: %s}",
		c.RampUpTime,
		c.RampDownTime,
		c.MinimumOnTime,
		c.DebounceTime)
}

func init() {
	messageFactory[CelaenoSetPointMessage] = func() ReceivableMessage { return &CelaenoSetPoint{} }
	messagesName[CelaenoSetPointMessage] = "Celaeno.SetPoint"
	messageFactory[CelaenoStatusMessage] = func() ReceivableMessage { return &CelaenoStatus{} }
	messagesName[CelaenoStatusMessage] = "Celaeno.Status"
	messageFactory[CelaenoConfigMessage] = func() ReceivableMessage { return &CelaenoConfig{} }
	messagesName[CelaenoConfigMessage] = "Celaeno.Config"
}
