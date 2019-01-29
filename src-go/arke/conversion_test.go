package arke

import (
	"math"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ConversionSuite struct{}

var _ = Suite(&ConversionSuite{})

func (s *ConversionSuite) TestHumidity(c *C) {
	testData := []struct {
		FloatValue  float32
		BinaryValue uint16
	}{
		{100.00, 16382},
		{0.0, 0},
		{50.0, 16382 / 2},
		{101.0, 16382},
		{-1.0, 0},
	}

	for _, d := range testData {
		c.Check(humidityFloatToBinary(d.FloatValue), Equals, d.BinaryValue)
		if d.FloatValue < 0 || d.FloatValue > 100.0 {
			continue
		}
		c.Check(humidityBinaryToFloat(d.BinaryValue), Equals, d.FloatValue)
	}

	c.Check(math.IsNaN(float64(humidityBinaryToFloat(hih6030Max+1))), Equals, true)
}

func (s *ConversionSuite) TestHih6030Temperature(c *C) {
	testData := []struct {
		FloatValue  float32
		BinaryValue uint16
	}{
		{125.00, 16382},
		{-40.0, 0},
		{-0.003967285, 3971},
		{126.0, 16382},
		{-41.0, 0},
	}

	for _, d := range testData {
		c.Check(hih6030TemperatureFloatToBinary(d.FloatValue), Equals, d.BinaryValue, Commentf("Converting %f", d.FloatValue))
		if d.FloatValue < -40.0 || d.FloatValue > 125.0 {
			continue
		}
		c.Check(hih6030TemperatureBinaryToFloat(d.BinaryValue), Equals, d.FloatValue)
	}

	c.Check(math.IsNaN(float64(hih6030TemperatureBinaryToFloat(hih6030Max+1))), Equals, true)
}

func (s *ConversionSuite) TestTmp1075Temperature(c *C) {

	testData := []struct {
		FloatValue  float32
		BinaryValue uint16
	}{
		{127.9375, 0x7FF},
		{100, 0x640},
		{80, 0x500},
		{75, 0x4B0},
		{50, 0x320},
		{25, 0x190},
		{0.25, 0x004},
		{0.0625, 0x001},
		{0, 0x000},
		{-0.0625, 0xFFF},
		{-0.25, 0xFFC},
		{-25, 0xE70},
		{-50, 0xCE0},
		{-128, 0x800},
	}

	for _, d := range testData {
		c.Check(tmp1075BinaryToFloat(d.BinaryValue), Equals, d.FloatValue)
	}
}
