package arke

import (
	"time"

	socketcan "github.com/atuleu/golang-socketcan"
	. "gopkg.in/check.v1"
)

type MessageSuite struct{}

var _ = Suite(&MessageSuite{})

func (s *MessageSuite) TestCANIDTIO(c *C) {
	testdata := []struct {
		IDT   uint32
		Type  MessageType
		Class MessageClass
		ID    NodeID
	}{
		{
			0x00,
			NetworkControlCommand,
			MessageClass(BroadcastClass),
			NodeID(ResetRequest),
		},
		{
			0x781,
			HeartBeat,
			MessageClass(CelaenoClass),
			1,
		},
	}

	for _, d := range testdata {
		resType, resClass, resID := ExtractCANIDT(d.IDT)
		c.Check(resType, Equals, d.Type)
		c.Check(resClass, Equals, d.Class)
		c.Check(resID, Equals, d.ID)
	}

	for _, d := range testdata {
		res := MakeCANIDT(d.Type, d.Class, d.ID)
		c.Check(res, Equals, d.IDT)
	}
}

func (s *MessageSuite) TestMessageParsing(c *C) {
	testdata := []struct {
		F  socketcan.CanFrame
		ID NodeID
		M  ReceivableMessage
	}{
		{
			socketcan.CanFrame{ID: MakeCANIDT(HeartBeat, MessageClass(ZeusClass), 2)},
			2,
			&HeartBeatData{ZeusClass, 2, 0, 0, 0, 0},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(HeartBeat, MessageClass(CelaenoClass), 4), Dlc: 4, Data: []byte{1, 2, 3, 4}},
			4,
			&HeartBeatData{CelaenoClass, 4, 1, 2, 3, 4},
		},

		{
			socketcan.CanFrame{ID: MakeCANIDT(NetworkControlCommand, MessageClass(ZeusClass), NodeID(HeartBeatRequest))},
			0,
			&HeartBeatRequestData{ZeusClass, 0},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(NetworkControlCommand, MessageClass(ZeusClass), NodeID(HeartBeatRequest)), Dlc: 2, Data: []byte{0xe8, 0x03}},
			0,
			&HeartBeatRequestData{ZeusClass, 1 * time.Second},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(NetworkControlCommand, MessageClass(0), NodeID(ResetRequest)), Dlc: 1, Data: []byte{0x00}},
			0,
			&ResetRequestData{BroadcastClass, BroadcastID},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(NetworkControlCommand, MessageClass(HeliosClass), NodeID(IDChangeRequest)), Dlc: 2, Data: []byte{0x01, 0x02}},
			1,
			&IDChangeRequestData{HeliosClass, 1, 2},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(NetworkControlCommand, MessageClass(0), NodeID(ErrorReport)), Dlc: 4, Data: []byte{byte(ZeusClass), 3, 0x42, 0}},
			3,
			&ErrorReportData{ZeusClass, 3, 0x0042},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusSetPointMessage, 2), Dlc: 5, Data: []byte{0, 0, 0, 0, 0}},
			2,
			&ZeusSetPoint{0, -40, 0},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusReportMessage, 3), Dlc: 8, Data: []byte{0, 0, 0, 0, 0, 0, 0, 0}},
			3,
			&ZeusReport{0, [4]float32{-40, 0, 0, 0}},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusConfigMessage, 4), Dlc: 8, Data: []byte{0, 0, 0, 0, 0, 0, 0, 0}},
			4,
			&ZeusConfig{PDConfig{0, 0, 0, 0, 0}, PDConfig{0, 0, 0, 0, 0}},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusStatusMessage, 5), Dlc: 7, Data: []byte{0, 0, 0, 0, 0, 0, 0}},
			5,
			&ZeusStatus{0, [3]FanStatusAndRPM{0, 0, 0}},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusControlPointMessage, 6), Dlc: 4, Data: []byte{2, 0, 3, 0}},
			6,
			&ZeusControlPoint{2, 3},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusDeltaTemperatureMessage, 7), Dlc: 8, Data: []byte{0, 0, 0, 0, 0, 0, 0, 0}},
			7,
			&ZeusDeltaTemperature{[4]float32{0, 0, 0, 0}},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, HeliosSetPointMessage, 1), Dlc: 2, Data: []byte{0x7f, 0xff}},
			1,
			&HeliosSetPoint{127, 255},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, CelaenoSetPointMessage, 1), Dlc: 1, Data: []byte{0x7f}},
			1,
			&CelaenoSetPoint{127},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, CelaenoStatusMessage, 2), Dlc: 3, Data: []byte{0x06, 0x00, 0x00}},
			2,
			&CelaenoStatus{WaterLevel: CelaenoWaterReadError, Fan: 0},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, CelaenoStatusMessage, 3), Dlc: 3, Data: []byte{0x02, 0x00, 0x00}},
			3,
			&CelaenoStatus{WaterLevel: CelaenoWaterCritical, Fan: 0},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, CelaenoConfigMessage, 4), Dlc: 8, Data: []byte{0xe8, 0x03, 0xe8, 0x03, 0xe8, 0x03, 0xe8, 0x03}},
			4,
			&CelaenoConfig{time.Second, time.Second, time.Second, time.Second},
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, CelaenoConfigMessage, 5), Dlc: 0, Data: []byte{}, RTR: true},
			5,
			&MessageRequestData{Class: CelaenoConfigMessage, ID: 5},
		},
	}

	for _, d := range testdata {
		m, ID, err := ParseMessage(&d.F)
		if c.Check(err, IsNil) == false {
			continue
		}
		c.Check(ID, Equals, d.ID)
		c.Check(m, DeepEquals, d.M)
		c.Check(m.MessageClassID(), Equals, d.M.MessageClassID())
	}

	errorData := []struct {
		F socketcan.CanFrame
		E string
	}{
		{
			socketcan.CanFrame{Extended: true},
			"Arke does not support extended IDT",
		},
		{
			socketcan.CanFrame{RTR: true, Dlc: 1},
			"RTR frame with a payload",
		},
		{
			socketcan.CanFrame{RTR: true, Dlc: 0, ID: MakeCANIDT(NetworkControlCommand, 0, 0)},
			"Unauthorized network command RTR frame",
		},

		{
			socketcan.CanFrame{ID: MakeCANIDT(NetworkControlCommand, 0, 6)},
			"Unknown network command 0x06",
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(HeartBeat, MessageClass(ZeusClass), 1), Dlc: 1, Data: []byte{0}},
			"Invalid buffer size 1 .*",
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, 0, 1), Dlc: 1, Data: []byte{0}},
			"Unknown message type 0x00",
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, 0, 1), Dlc: 0, Data: []byte{}, RTR: true},
			"Unknown message type 0x00",
		},
		{
			socketcan.CanFrame{ID: MakeCANIDT(StandardMessage, ZeusReportMessage, 1), Dlc: 1, Data: []byte{0}},
			"Could not parse message data: .*",
		},
	}

	for _, d := range errorData {
		_, _, err := ParseMessage(&d.F)
		c.Check(err, ErrorMatches, d.E)
	}

	// does nothing
	c.Check((&MessageRequestData{}).Unmarshal(nil), Equals, nil)
}

func (s *MessageSuite) TestMessagesName(c *C) {
	testdata := []struct {
		C MessageClass
		E string
	}{
		{ZeusSetPointMessage, "Zeus.SetPoint"},
		{ZeusReportMessage, "Zeus.Report"},
		{ZeusVibrationReportMessage, "Zeus.VibrationReport"},
		{ZeusConfigMessage, "Zeus.Config"},
		{ZeusStatusMessage, "Zeus.Status"},
		{ZeusControlPointMessage, "Zeus.ControlPoint"},
		{ZeusDeltaTemperatureMessage, "Zeus.DeltaTemperature"},
		{HeliosSetPointMessage, "Helios.SetPoint"},
		{HeliosPulseModeMessage, "Helios.PulseMode"},
		{HeliosTriggerModeMessage, "Helios.TriggerMode"},
		{CelaenoSetPointMessage, "Celaeno.SetPoint"},
		{CelaenoStatusMessage, "Celaeno.Status"},
		{CelaenoConfigMessage, "Celaeno.Config"},
		{NotusSetPointMessage, "Notus.SetPoint"},
		{NotusConfigMessage, "Notus.Config"},
		{0, "<unknown>"},
	}

	for _, d := range testdata {
		c.Check(d.C.String(), Equals, d.E)
	}

}
