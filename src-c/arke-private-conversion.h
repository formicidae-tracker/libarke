#pragma once

#include <math.h>

#include <inttypes.h>

#define MAX_HIH6030_VALUE 16382

#ifdef __cplusplus
extern "C" {
#endif



inline uint16_t float_to_humidity(float value) {
	if (value <= 0.0 ) {
		return 0;
	} else if (value >= 100.0 ) {
		return MAX_HIH6030_VALUE;
	} else {
		return value / 100.0 * (float)MAX_HIH6030_VALUE;
	}
}

inline float humidity_to_float(uint16_t value) {
	if (value > MAX_HIH6030_VALUE ) {
		return NAN;
	}
	return (float)value / ((float)MAX_HIH6030_VALUE) * 100.0;
}

inline uint16_t float_to_hih6030_temperature(float value) {
	if ( value <= -40.0 ) {
		return 0;
	} else if ( value >= 125.0 ) {
		return MAX_HIH6030_VALUE;
	} else {
		return (value + 40.0) / 165.0 * MAX_HIH6030_VALUE;
	}
}

inline float hih6030_temperature_to_float(uint16_t value) {
	if (value > MAX_HIH6030_VALUE ) {
		return NAN;
	}
	return (float)value / ((float)MAX_HIH6030_VALUE) * 165.0 - 40.0;
}

inline float tmp1075_to_float(uint16_t value) {
	// conversion from 12-bit ninary complement to 16-bit binary complement
	int16_t sValue = value;
	if ( (value & 0x0800) != 0 ) {
		sValue |= 0xf000;
	}
	return sValue * 0.0625;
}

#ifdef __cplusplus
}
#endif
