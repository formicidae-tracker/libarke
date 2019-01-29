package arke

import (
	"encoding/binary"
	"fmt"
	"math"
)

// #include "../../include/arke.h"
import "C"

type ZeusSetPoint struct {
	Humidity   float32
	Temprature float32
	Wind       uint8
}

func (m ZeusSetPoint) Marshall(buf []byte) uint8 {
	binary.LittleEndian.PutUint16(buf[0:], humidityFloatToBinary(m.Humidity))
	binary.LittleEndian.PutUint16(buf[2:], hih6030TemperatureFloatToBinary(m.Temprature))
	buf[4] = m.Wind
	return 5
}

func (m ZeusSetPoint) Unmarshall(buf []byte) error {
	m.Humidity = humidityBinaryToFloat(binary.LittleEndian.Uint16(buf[0:]))
	if math.IsNaN(float64(m.Humidity)) == true {
		return fmt.Errorf("Invalid humidity value")
	}
	m.Temprature = hih6030TemperatureBinaryToFloat(binary.LittleEndian.Uint16(buf[2:]))
	if math.IsNaN(float64(m.Temprature)) == true {
		return fmt.Errorf("Invalid temperature value")
	}
	m.Wind = buf[4]
	return nil
}

type ZeusReport struct {
	Humidity    float32
	Temperature [4]float32
}

func (m ZeusReport) Unmarshall(buf []byte) error {
	packed := []uint16{
		binary.LittleEndian.Uint16(buf[0:]),
		binary.LittleEndian.Uint16(buf[2:]),
		binary.LittleEndian.Uint16(buf[4:]),
		binary.LittleEndian.Uint16(buf[6:]),
	}
	m.Humidity = humidityBinaryToFloat(packed[0] >> 2)
	if math.IsNaN(float64(m.Humidity)) == true {
		return fmt.Errorf("Invalid humidity value")
	}

	m.Temperature[0] = hih6030TemperatureBinaryToFloat(((packed[0] & 0x03) << 12) | packed[1]>>4)
	if math.IsNaN(float64(m.Temperature[0])) == true {
		return fmt.Errorf("Invalid Temperature[0] value")
	}
	m.Temperature[1] = tmp1075BinaryToFloat(((packed[1] & 0x0f) << 8) | packed[2]>>8)
	m.Temperature[2] = tmp1075BinaryToFloat(((packed[2] & 0xff) << 4) | packed[3]>>12)
	m.Temperature[3] = tmp1075BinaryToFloat(packed[3] & 0x0fff)
	return nil
}

type ZeusConfig struct {
	Humidity    PDConfig
	Temperature PDConfig
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
