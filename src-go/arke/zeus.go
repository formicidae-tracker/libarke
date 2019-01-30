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
		return fmt.Errorf("Invalid Temperature[0] value")
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

func (m ZeusConfig) Marshall(buf []byte) (int, error) {
	if len(buf) < 8 {
		return 0, fmt.Errorf("Invalid buffer size %d, required 8", len(buf))
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
	m.Humidity.unmarshall(buf[0:])
	m.Temperature.unmarshall(buf[4:])
	return nil
}

type ZeusStatus struct {
	status      uint8
	fans        [2]FanStatus
	humidity    int16
	temperature int16
}

type ZeusControlStatus uint8

const (
	ZeusIdle                         ZeusControlStatus = C.ARKE_ZEUS_IDLE
	ZeusActive                       ZeusControlStatus = C.ARKE_ZEUS_ACTIVE
	ZeusClimateNotControlledWatchDog ZeusControlStatus = C.ARKE_ZEUS_CLIMATE_UNCONTROLLED_WD
)

func (s ZeusStatus) ControlStatus() ZeusControlStatus {
	return ZeusControlStatus(s.status & 0x03)
}

func (s ZeusStatus) IsFanStatus() bool {
	return (s.status)&C.ARKE_ZEUS_STATUS_IS_COMMAND_DATA == 0
}

func (s ZeusStatus) IsCommandData() bool {
	return (s.status)&C.ARKE_ZEUS_STATUS_IS_COMMAND_DATA == 0
}

func (s ZeusStatus) FanStatus() ([2]FanStatus, error) {
	if s.IsCommandData() {
		return [2]FanStatus{}, fmt.Errorf("Packet contains command data")
	}
	return s.fans, nil
}

func (s ZeusStatus) HumidityCommand() (int16, error) {
	if s.IsFanStatus() {
		return 0, fmt.Errorf("Packet contains fan status")
	}
	return s.humidity, nil
}

func (s ZeusStatus) TemperatureCommand() (int16, error) {
	if s.IsFanStatus() {
		return 0, fmt.Errorf("Packet contains fan status")
	}
	return s.temperature, nil
}

func (m ZeusStatus) Unmarshall(buf []byte) error {
	m.status = buf[0]
	if m.IsFanStatus() {
		m.fans[0] = FanStatus(binary.LittleEndian.Uint16(buf[1:]))
		m.fans[1] = FanStatus(binary.LittleEndian.Uint16(buf[3:]))
	} else {
		m.humidity = int16(binary.LittleEndian.Uint16(buf[1:]))
		m.temperature = int16(binary.LittleEndian.Uint16(buf[3:]))
	}
	return nil
}
