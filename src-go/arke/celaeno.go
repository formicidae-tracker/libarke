package arke

// #include "../../include/arke.h"
import "C"

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

type WaterLevelStatus uint8

const (
	CelaenoWaterNominal   WaterLevelStatus = C.ARKE_CELAENO_NOMINAL
	CelaenoWaterWarning   WaterLevelStatus = C.ARKE_CELAENO_WARNING
	CelaenoWaterCritical  WaterLevelStatus = C.ARKE_CELAENO_CRITICAL
	CelaenoWaterReadError WaterLevelStatus = C.ARKE_CELAENO_RO_ERROR
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
	if (buf[0] & C.ARKE_CELAENO_RO_ERROR) != 0 {
		c.WaterLevel = CelaenoWaterReadError
	} else if buf[0]&C.ARKE_CELAENO_CRITICAL != 0 {
		c.WaterLevel = CelaenoWaterCritical
	} else {
		c.WaterLevel = WaterLevelStatus(buf[0])
	}

	c.Fan = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[1:]))
	return nil
}

type CelaenoConfig struct {
	RampUpTime    time.Duration
	RampDownTime  time.Duration
	MinimumOnTime time.Duration
	DebounceTime  time.Duration
}

func (m *CelaenoConfig) MessageClassID() MessageClass {
	return CelaenoConfigNessage
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

func init() {
	messageFactory[CelaenoSetPointMessage] = func() ReceivableMessage { return &CelaenoSetPoint{} }
	messageFactory[CelaenoStatusMessage] = func() ReceivableMessage { return &CelaenoStatus{} }
	messageFactory[CelaenoConfigNessage] = func() ReceivableMessage { return &CelaenoConfig{} }
}
