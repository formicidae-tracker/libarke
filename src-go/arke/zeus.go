package arke

import (
	"encoding/binary"
	"fmt"
	"math"
)

type ZeusSetPoint struct {
	Humidity    float32
	Temperature float32
	Wind        uint8
}

func (m *ZeusSetPoint) MessageClassID() MessageClass {
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

func (sp *ZeusSetPoint) String() string {
	return fmt.Sprintf("Zeus.SetPoint{Humidity: %.2f%%, Temperature: %.2f°C, Wind: %d}",
		sp.Humidity, sp.Temperature, sp.Wind)
}

type ZeusReport struct {
	Humidity    float32
	Temperature [4]float32
}

func (m *ZeusReport) MessageClassID() MessageClass {
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

func (sp *ZeusReport) String() string {
	return fmt.Sprintf("Zeus.Report{Humidity: %.2f%%, Ant: %.2f°C, Aux1: %.2f°C, Aux2: %.2f°C, Aux3: %.2f°C}",
		sp.Humidity, sp.Temperature[0], sp.Temperature[1], sp.Temperature[2], sp.Temperature[3])
}

type ZeusConfig struct {
	Humidity    PDConfig
	Temperature PDConfig
}

func (c *ZeusConfig) String() string {
	return fmt.Sprintf("Zeus.Config{Humidity:%s, Temperature:%s}",
		c.Humidity, c.Temperature)
}

func (m *ZeusConfig) MessageClassID() MessageClass {
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

type ZeusStatusValue uint8

const (
	ZeusIdle                         ZeusStatusValue = 0x00
	ZeusActive                       ZeusStatusValue = 1 << 0
	ZeusClimateNotControlledWatchDog ZeusStatusValue = 1 << 1
	ZeusHumidityUnreachable          ZeusStatusValue = 1 << 2
	ZeusTemperatureUnreachable       ZeusStatusValue = 1 << 3
)

type ZeusStatus struct {
	Status ZeusStatusValue
	Fans   [3]FanStatusAndRPM
}

func (s ZeusStatusValue) String() string {
	prefix := ""
	if s&ZeusTemperatureUnreachable != 0 {
		prefix += "temperature-unreachable|"
	}
	if s&ZeusHumidityUnreachable != 0 {
		prefix += "humidity-unreachable|"
	}
	if s&ZeusClimateNotControlledWatchDog != 0 {
		if s&ZeusActive != 0 {
			return prefix + "sensor-issue"
		}
		prefix += "climate-uncontrolled|"
	}
	if s&ZeusActive != 0 {
		return prefix + "active"
	}
	return prefix + "idle"
}

func (s *ZeusStatus) String() string {
	return fmt.Sprintf("Zeus.Status{General: %s, WindFan: %s, LeftFan: %s, RightFan: %s}",
		s.Status,
		s.Fans[0],
		s.Fans[2],
		s.Fans[1],
	)
}

func (m *ZeusStatus) MessageClassID() MessageClass {
	return ZeusStatusMessage
}

func (m *ZeusStatus) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 7); err != nil {
		return err
	}
	m.Status = ZeusStatusValue(buf[0])
	m.Fans[0] = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[1:]))
	m.Fans[1] = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[3:]))
	m.Fans[2] = FanStatusAndRPM(binary.LittleEndian.Uint16(buf[5:]))
	return nil
}

type ZeusControlPoint struct {
	Humidity    int16
	Temperature int16
}

func (cp *ZeusControlPoint) String() string {
	return fmt.Sprintf("Zeus.ControlPoint{Humidity: %d, Temperature: %d}",
		cp.Humidity,
		cp.Temperature,
	)
}

func (m *ZeusControlPoint) MessageClassID() MessageClass {
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

type ZeusDeltaTemperature struct {
	Delta [4]float32
}

func (d *ZeusDeltaTemperature) String() string {
	return fmt.Sprintf("Zeus.DeltaTemperature{Ants: %.4f°C, Aux1: %.4f°C, Aux2: %.4f°C, Aux3: %.4f°C}",
		d.Delta[0],
		d.Delta[1],
		d.Delta[2],
		d.Delta[3],
	)
}

func (m *ZeusDeltaTemperature) MessageClassID() MessageClass {
	return ZeusDeltaTemperatureMessage
}

func (m *ZeusDeltaTemperature) Marshall(buf []byte) (int, error) {
	binary.LittleEndian.PutUint16(buf[0:], uint16(int16(m.Delta[0]*float32(hih6030Max)/165.0)))
	for i := 1; i < 4; i++ {
		binary.LittleEndian.PutUint16(buf[(2*i):], uint16(int16(m.Delta[i]/0.0625)))
	}
	return 8, nil
}

func (m *ZeusDeltaTemperature) Unmarshall(buf []byte) error {
	if err := checkSize(buf, 8); err != nil {
		return err
	}

	m.Delta[0] = float32(int16(binary.LittleEndian.Uint16(buf[0:]))) * 165.0 / float32(hih6030Max)
	for i := 1; i < 4; i++ {
		m.Delta[i] = float32(int16(binary.LittleEndian.Uint16(buf[(2*i):]))) * 0.0625
	}

	return nil
}

func init() {
	messageFactory[ZeusSetPointMessage] = func() ReceivableMessage { return &ZeusSetPoint{} }
	messagesName[ZeusSetPointMessage] = "Zeus.SetPoint"
	messageFactory[ZeusReportMessage] = func() ReceivableMessage { return &ZeusReport{} }
	messagesName[ZeusReportMessage] = "Zeus.Report"
	messageFactory[ZeusConfigMessage] = func() ReceivableMessage { return &ZeusConfig{} }
	messagesName[ZeusConfigMessage] = "Zeus.Config"
	messageFactory[ZeusStatusMessage] = func() ReceivableMessage { return &ZeusStatus{} }
	messagesName[ZeusStatusMessage] = "Zeus.Status"
	messageFactory[ZeusControlPointMessage] = func() ReceivableMessage { return &ZeusControlPoint{} }
	messagesName[ZeusControlPointMessage] = "Zeus.ControlPoint"
	messageFactory[ZeusDeltaTemperatureMessage] = func() ReceivableMessage { return &ZeusDeltaTemperature{} }
	messagesName[ZeusDeltaTemperatureMessage] = "Zeus.DeltaTemperature"
	messagesName[ZeusVibrationReportMessage] = "Zeus.VibrationReport"
}
