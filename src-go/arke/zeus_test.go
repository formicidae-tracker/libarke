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
	testData := []struct {
		Message ZeusSetPoint
		Buffer  []byte
	}{
		{
			Message: ZeusSetPoint{
				Humidity:    42,
				Temperature: 25,
				Wind:        127,
			},
			Buffer: []byte{
				0xe0, 0x1a, 0x35, 0x19, 0x7f,
			},
		},
	}

	for _, d := range testData {
		res := make([]byte, 5)
		written, err := d.Message.Marshall(res)
		if c.Check(err, IsNil) == false {
			continue
		}
		c.Check(written, Equals, 5)
		c.Check(res, DeepEquals, d.Buffer)
	}

	for _, d := range testData {
		res := ZeusSetPoint{}
		if c.Check(res.Unmarshall(d.Buffer), IsNil) == false {
			continue
		}
		c.Check(res.Wind, Equals, d.Message.Wind)
		c.Check(res.Humidity, AlmostChecker, d.Message.Humidity, 0.01)
		c.Check(res.Temperature, AlmostChecker, d.Message.Temperature, 0.01)
	}

	m := ZeusSetPoint{}
	written, err := m.Marshall([]byte{})
	c.Check(err, ErrorMatches, "Invalid buffer size .*")
	c.Check(written, Equals, 0)
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, "Invalid buffer size .*")
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
	testData := []struct {
		Message ZeusReport
		Buffer  []byte
	}{
		{
			Message: ZeusReport{
				Humidity: 40,
				Temperature: [4]float32{
					25, 26, 27, 28,
				},
			},
			Buffer: []byte{
				0x99, 0x99,
				0x4d, 0x06,
				0x1a, 0xb0,
				0x01, 0x1c,
			},
		},
	}

	for _, d := range testData {
		res := ZeusReport{}
		if c.Check(res.Unmarshall(d.Buffer), IsNil) == false {
			continue
		}
		c.Check(res.Humidity, AlmostChecker, d.Message.Humidity, 0.01)
		for i, _ := range res.Temperature {
			c.Check(res.Temperature[i], AlmostChecker, d.Message.Temperature[i], 0.01)
		}

	}
	m := ZeusReport{}
	c.Check(m.Unmarshall([]byte{}), ErrorMatches, "Invalid buffer size .*")
}
