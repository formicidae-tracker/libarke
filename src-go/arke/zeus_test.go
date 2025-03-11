package arke

import (
	"fmt"
	"math"

	. "gopkg.in/check.v1"
)

type ZeusSuite struct{}

var _ = Suite(&ZeusSuite{})

type almostEqualChecker struct {
	*CheckerInfo
}

var AlmostChecker = &almostEqualChecker{
	&CheckerInfo{Name: "AlmostEqual", Params: []string{"obtained", "expected", "within"}},
}

func (c *almostEqualChecker) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()
	a := params[0].(float32)
	b := params[1].(float32)
	bound := params[2].(float64)
	result = math.Abs(float64(a-b)) <= bound
	return
}

func (s *ZeusSuite) TestSetPointIO(c *C) {
	checkMessageEncoding(c, &ZeusSetPoint{
		Humidity:    41.997314,
		Temperature: 24.994812,
		Wind:        127,
	}, []byte{
		0xe0, 0x1a, 0x35, 0x19, 0x7f,
	})

	checkMessageLength(c, &ZeusSetPoint{}, 5)
	errorData := []struct {
		Buffer []byte
		EMatch string
	}{
		{
			[]byte{
				0xff, 0xff, 0x00, 0x00, 0x00,
			},
			"Invalid humidity value",
		},
		{
			[]byte{
				0x00, 0x00, 0xff, 0xff, 0x00,
			},
			"Invalid temperature value",
		},
	}

	for _, d := range errorData {
		m := ZeusSetPoint{}
		c.Check(m.Unmarshall(d.Buffer), ErrorMatches, d.EMatch)
	}

}

func (s *ZeusSuite) TestReportIO(c *C) {

	checkMessageEncoding(c, &ZeusReport{
		Humidity: 40.0012207,
		Temperature: [4]float32{
			25.0048828, 26, 27, 28,
		},
	}, []byte{
		0x99, 0x99,
		0x4d, 0x06,
		0x1a, 0xb0,
		0x01, 0x1c,
	})

	errorData := []struct {
		Buffer []byte
		Ematch string
	}{
		{
			[]byte{0xff, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			"Invalid humidity value",
		},
		{
			[]byte{0x00, 0xc0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00},
			"Invalid temperature value",
		},
	}

	for _, d := range errorData {
		m := ZeusReport{}
		c.Check(m.Unmarshall(d.Buffer), ErrorMatches, d.Ematch)
	}
}

func (s *ZeusSuite) TestConfigIO(c *C) {

	checkMessageEncoding(c, &ZeusConfig{
		Humidity: PDConfig{
			ProportionnalMultiplier: 100,
			DerivativeMultiplier:    50,
			IntegralMultiplier:      1,
			DividerPower:            6,
		},
		Temperature: PDConfig{
			ProportionnalMultiplier: 103,
			DerivativeMultiplier:    102,
			IntegralMultiplier:      0,
			DividerPower:            4,
		},
	}, []byte{
		100, 50, 1, 6, 103, 102, 0, 4,
	})

	checkMessageLength(c, &ZeusConfig{}, 8)

	errorData := []struct {
		Message ZeusConfig
		Buffer  []byte
		Ematch  string
	}{
		{
			ZeusConfig{},
			make([]byte, 0),
			"Invalid buffer size .*",
		},
		{
			ZeusConfig{PDConfig{0, 0, 0, 16, 0}, PDConfig{}},
			make([]byte, 8),
			"Maximal Proportional&Derivative Divider is 15",
		},
		{
			ZeusConfig{PDConfig{}, PDConfig{0, 0, 0, 0, 16}},
			make([]byte, 8),
			"Maximal Integral Divider is 15",
		},
	}

	for _, d := range errorData {
		_, err := d.Message.Marshall(d.Buffer)
		c.Check(err, ErrorMatches, d.Ematch)
	}

}

func (s *ZeusSuite) TestStatusIO(c *C) {
	checkMessageEncoding(c, &ZeusStatus{
		Status: ZeusIdle,
		Fans: [3]FanStatusAndRPM{
			1200,
			FanStatusAndRPM(0 | (uint16(FanStalled) << 14)),
			FanStatusAndRPM(400 | (uint16(FanAging) << 14)),
		},
	}, []byte{
		0, 0xb0, 0x04, 0x00, 0x80, 0x90, 0x41,
	})
	checkMessageLength(c, &ZeusStatus{}, 7)
}

func (s *ZeusSuite) TestControlPointIO(c *C) {
	checkMessageEncoding(c, &ZeusControlPoint{
		Humidity:    1234,
		Temperature: -275,
	}, []byte{
		0xd2, 0x04, 0xed, 0xfe,
	})

	checkMessageLength(c, &ZeusControlPoint{}, 4)
}

func (s *ZeusSuite) TestTemperatureDelta(c *C) {
	checkMessageEncoding(c, &ZeusDeltaTemperature{
		Delta: [4]float32{0, 0, 0, 0},
	}, []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
	checkMessageEncoding(c, &ZeusDeltaTemperature{
		Delta: [4]float32{-0.75540227078, 2.625, -1, 0},
	}, []byte{
		0xb5, 0xff, 42, 0x00, 0xf0, 0xff, 0x00, 0x00,
	})

	checkMessageLength(c, &ZeusDeltaTemperature{}, 8)
}
