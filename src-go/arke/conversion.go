package arke

import (
	"math"
)

const hih6030Max = 16382

func humidityBinaryToFloat(value uint16) float32 {
	if value > hih6030Max {
		return float32(math.NaN())
	}
	return float32(value) / float32(hih6030Max) * 100.0
}

func humidityFloatToBinary(value float32) uint16 {
	if value <= 0 {
		return 0
	} else if value >= 100.0 {
		return hih6030Max
	}

	return uint16((value / 100.0) * hih6030Max)
}

func hih6030TemperatureBinaryToFloat(value uint16) float32 {
	if value > hih6030Max {
		return float32(math.NaN())
	}
	return float32(value)/float32(hih6030Max)*165.0 - 40.0
}

func hih6030TemperatureFloatToBinary(value float32) uint16 {
	if value <= -40.0 {
		return 0
	} else if value >= 125.0 {
		return hih6030Max
	}

	return uint16(((value + 40.0) / 165.0) * hih6030Max)
}

func tmp1075BinaryToFloat(value uint16) float32 {
	if value&0x800 != 0 {
		value = 0xf000 | value
	}
	return float32(int16(value)) * 0.0625
}

const MAX_INT12 uint16 = (1 << 11) - 1

func tmp1075FloatToBinaray(value float32) uint16 {
	if value >= 0 {
		return max(0, min(MAX_INT12, uint16(value/0.0625)))
	} else if value <= -128.0 {
		return 0x0800
	}
	return (0xffff - max(0, min(MAX_INT12, uint16(-value/0.0625))) + 1) & 0xfff
}
