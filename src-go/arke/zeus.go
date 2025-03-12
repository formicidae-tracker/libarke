package arke

import (
	"encoding/binary"
	"fmt"
	"math"
)

type ZeusSetPoint struct {
	Humidity    float32 `positional-arg-name:"humidity" required:"yes"`
	Temperature float32 `positional-arg-name:"temperature" required:"yes"`
	Wind        uint8   `positional-arg-name:"wind" required:"yes"`
}

func (m *ZeusSetPoint) MessageClassID() MessageClass {
	return ZeusSetPointMessage
}

func checkSize(buf []byte, expected int) error {
	if len(buf) < expected {
		return fmt.Errorf("Invalid buffer size %d, required: %d", len(buf), expected)
	}
	return nil
}

func (m ZeusSetPoint) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 5); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], humidityFloatToBinary(m.Humidity))
	binary.LittleEndian.PutUint16(buf[2:], hih6030TemperatureFloatToBinary(m.Temperature))
	buf[4] = m.Wind
	return 5, nil
}

func (m *ZeusSetPoint) Unmarshal(buf []byte) error {
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

func (m *ZeusSetPoint) String() string {
	return fmt.Sprintf("Zeus.SetPoint{Humidity: %.2f%%, Temperature: %.2f°C, Wind: %d}",
		m.Humidity, m.Temperature, m.Wind)
}

type ZeusReport struct {
	Humidity    float32
	Temperature [4]float32
}

func (m *ZeusReport) MessageClassID() MessageClass {
	return ZeusReportMessage
}

func (m ZeusReport) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 8); err != nil {
		return 0, err
	}
	packed := make([]uint16, 4)
	binTemp := hih6030TemperatureFloatToBinary(m.Temperature[0])
	auxs := []uint16{
		tmp1075FloatToBinaray(m.Temperature[1]),
		tmp1075FloatToBinaray(m.Temperature[2]),
		tmp1075FloatToBinaray(m.Temperature[3]),
	}
	packed[0] = humidityFloatToBinary(m.Humidity) | (binTemp << 14)
	packed[1] = binTemp>>2 | (auxs[0]&0xf)<<12
	packed[2] = auxs[0]>>4 | (auxs[1]&0xff)<<8
	packed[3] = auxs[1]>>8 | (auxs[2]&0xfff)<<4

	for i, word := range packed {
		binary.LittleEndian.PutUint16(buf[2*i:], word)
	}

	return 8, nil
}

func (m *ZeusReport) Unmarshal(buf []byte) error {
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

func (m *ZeusReport) String() string {
	return fmt.Sprintf("Zeus.Report{Humidity: %.2f%%, Ant: %.2f°C, Aux1: %.2f°C, Aux2: %.2f°C, Aux3: %.2f°C}",
		m.Humidity, m.Temperature[0], m.Temperature[1], m.Temperature[2], m.Temperature[3])
}

type ZeusConfig struct {
	Humidity    PDConfig
	Temperature PDConfig
}

func (m *ZeusConfig) String() string {
	return fmt.Sprintf("Zeus.Config{Humidity:%s, Temperature:%s}",
		m.Humidity, m.Temperature)
}

func (m *ZeusConfig) MessageClassID() MessageClass {
	return ZeusConfigMessage
}

func (m ZeusConfig) Marshal(buf []byte) (int, error) {
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

func (m *ZeusConfig) Unmarshal(buf []byte) error {
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

func (m ZeusStatus) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 7); err != nil {
		return 0, err
	}
	buf[0] = byte(m.Status)
	binary.LittleEndian.PutUint16(buf[1:], uint16(m.Fans[0]))
	binary.LittleEndian.PutUint16(buf[3:], uint16(m.Fans[1]))
	binary.LittleEndian.PutUint16(buf[5:], uint16(m.Fans[2]))
	return 7, nil
}

func (m *ZeusStatus) Unmarshal(buf []byte) error {
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

func (m *ZeusControlPoint) String() string {
	return fmt.Sprintf("Zeus.ControlPoint{Humidity: %d, Temperature: %d}",
		m.Humidity,
		m.Temperature,
	)
}

func (m *ZeusControlPoint) MessageClassID() MessageClass {
	return ZeusControlPointMessage
}

func (m *ZeusControlPoint) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 4); err != nil {
		return 0, err
	}
	binary.LittleEndian.PutUint16(buf[0:], uint16(m.Humidity))
	binary.LittleEndian.PutUint16(buf[2:], uint16(m.Temperature))
	return 4, nil
}

func (m *ZeusControlPoint) Unmarshal(buf []byte) error {
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

func (m *ZeusDeltaTemperature) String() string {
	return fmt.Sprintf("Zeus.DeltaTemperature{Ants: %.4f°C, Aux1: %.4f°C, Aux2: %.4f°C, Aux3: %.4f°C}",
		m.Delta[0],
		m.Delta[1],
		m.Delta[2],
		m.Delta[3],
	)
}

func (m *ZeusDeltaTemperature) MessageClassID() MessageClass {
	return ZeusDeltaTemperatureMessage
}

func (m *ZeusDeltaTemperature) Marshal(buf []byte) (int, error) {
	if err := checkSize(buf, 8); err != nil {
		return 0, err
	}

	binary.LittleEndian.PutUint16(buf[0:], uint16(int16(m.Delta[0]*float32(hih6030Max)/165.0)))
	for i := 1; i < 4; i++ {
		binary.LittleEndian.PutUint16(buf[(2*i):], uint16(int16(m.Delta[i]/0.0625)))
	}
	return 8, nil
}

func (m *ZeusDeltaTemperature) Unmarshal(buf []byte) error {
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
	messageFactory[ZeusSetPointMessage] = func() Message { return &ZeusSetPoint{} }
	messagesName[ZeusSetPointMessage] = "Zeus.SetPoint"
	messageFactory[ZeusReportMessage] = func() Message { return &ZeusReport{} }
	messagesName[ZeusReportMessage] = "Zeus.Report"
	messageFactory[ZeusConfigMessage] = func() Message { return &ZeusConfig{} }
	messagesName[ZeusConfigMessage] = "Zeus.Config"
	messageFactory[ZeusStatusMessage] = func() Message { return &ZeusStatus{} }
	messagesName[ZeusStatusMessage] = "Zeus.Status"
	messageFactory[ZeusControlPointMessage] = func() Message { return &ZeusControlPoint{} }
	messagesName[ZeusControlPointMessage] = "Zeus.ControlPoint"
	messageFactory[ZeusDeltaTemperatureMessage] = func() Message { return &ZeusDeltaTemperature{} }
	messagesName[ZeusDeltaTemperatureMessage] = "Zeus.DeltaTemperature"
	messagesName[ZeusVibrationReportMessage] = "Zeus.VibrationReport"
}
