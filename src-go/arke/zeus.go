package arke

import (
	"encoding/binary"
	"fmt"
	"math"
)

// #include "../../include/arke.h"
import "C"

type ZeusSetPoint struct {
	Humidity    float32
	Temperature float32
	Wind        uint8
}

func (m *ZeusSetPoint) MessageClassID() uint16 {
	return ZeusSetPointMessage
}

func checkSize(buf []byte, expected int) error {
	if len(buf) < expected {
		return fmt.Errorf("Invalid buffer size %d, required %d", len(buf), expected)
	}
	return nil
}

func (m ZeusSetPoint) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 5); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], humidityFloatToBinary(m.Humidity))
	binary.LittleEndian.PutUint16(buf[2:], hih6030TemperatureFloatToBinary(m.Temperature))
	buf[4] = m.Wind
	return 5, nil
}

func (m *ZeusSetPoint) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 5); err != nil {
		return err
	}
	m.Humidity = humidityBinaryToFloat(binary.LittleEndian.Uint16(buf[0:]))
	if math.IsNaN(float64(m.Humidity)) == true {
		return fmt.Errorf("Invalid humidity value")
	}
	m.Temperature = hih6030TemperatureBinaryToFloat(binary.LittleEndian.Uint16(buf[2:]))
	if math.IsNaN(float64(m.Temperature)) == true {
		return fmt.Errorf("Invalid temperature value")
	}
	m.Wind = buf[4]
	return nil
}

type ZeusReport struct {
	Humidity    float32
	Temperature [4]float32
}

func (m *ZeusReport) MessageClassID() uint16 {
	return ZeusReportMessage
}

func (m *ZeusReport) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 8); err != nil {
		return err
	}
	packed := []uint16{
		binary.LittleEndian.Uint16(buf[0:]),
		binary.LittleEndian.Uint16(buf[2:]),
		binary.LittleEndian.Uint16(buf[4:]),
		binary.LittleEndian.Uint16(buf[6:]),
	}
	m.Humidity = humidityBinaryToFloat(packed[0] & 0x3fff)
	if math.IsNaN(float64(m.Humidity)) == true {
		return fmt.Errorf("Invalid humidity value")
	}

	m.Temperature[0] = hih6030TemperatureBinaryToFloat((packed[0] >> 14) | (packed[1]&0x0fff)<<2)
	if math.IsNaN(float64(m.Temperature[0])) == true {
		return fmt.Errorf("Invalid temperature value")
	}
	m.Temperature[1] = tmp1075BinaryToFloat((packed[1] >> 12) | (packed[2]&0x00ff)<<4)
	m.Temperature[2] = tmp1075BinaryToFloat((packed[2] >> 8) | (packed[3]&0x000f)<<8)
	m.Temperature[3] = tmp1075BinaryToFloat((packed[3] & 0xfff0) >> 4)
	return nil
}

type ZeusConfig struct {
	Humidity    PDConfig
	Temperature PDConfig
}

func (m *ZeusConfig) MessageClassID() uint16 {
	return ZeusConfigMessage
}

func (m ZeusConfig) Marshall(buf []byte) (int, error) {
	if err := checkSize(buf, 8); err != nil {
		return 0, err
	}
	if err := m.Humidity.marshall(buf[0:]); err != nil {
		return 0, err
	}
	if err := m.Temperature.marshall(buf[4:]); err != nil {
		return 4, err
	}
	return 8, nil
}

func (m *ZeusConfig) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 8); err != nil {
		return err
	}
	m.Humidity.unmarshall(buf[0:])
	m.Temperature.unmarshall(buf[4:])
	return nil
}

const (
	ZeusIdle                         uint8 = C.ARKE_ZEUS_IDLE
	ZeusActive                       uint8 = C.ARKE_ZEUS_ACTIVE
	ZeusClimateNotControlledWatchDog uint8 = C.ARKE_ZEUS_CLIMATE_UNCONTROLLED_WD
)

type ZeusStatus struct {
	Status uint8
	Fans   [2]FanStatusAndRPM
}

func (m *ZeusStatus) MessageClassID() uint16 {
	return ZeusStatusMessage
}

func (m *ZeusStatus) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 5); err != nil {
		return err
	}
	m.Status = buf[0]
	m.Fans[0] = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[1:]))
	m.Fans[1] = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[3:]))
	return nil
}

type ZeusControlPoint struct {
	Humidity    int16
	Temperature int16
}

func (m *ZeusControlPoint) MessageClassID() uint16 {
	return ZeusControlPointMessage
}

func (m *ZeusControlPoint) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 4); err != nil {
		return err
	}
	m.Humidity = int16(binary.LittleEndian.Uint16(buf[0:]))
	m.Temperature = int16(binary.LittleEndian.Uint16(buf[2:]))
	return nil
}

func init() {
	messageFactory[ZeusSetPointMessage] = func() ReceivableMessage { return &ZeusSetPoint{} }
	messageFactory[ZeusReportMessage] = func() ReceivableMessage { return &ZeusReport{} }
	messageFactory[ZeusConfigMessage] = func() ReceivableMessage { return &ZeusConfig{} }
	messageFactory[ZeusStatusMessage] = func() ReceivableMessage { return &ZeusStatus{} }
	messageFactory[ZeusControlPointMessage] = func() ReceivableMessage { return &ZeusControlPoint{} }
}
