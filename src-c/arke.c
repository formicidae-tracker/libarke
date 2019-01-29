#include "arke-c.h"



#include "arke-private-conversion.h"


void ArkeZeusSetTargetHumidity(ArkeZeusSetPoint * sp,float value) {
	sp->Humidity = float_to_humidity(value);
}

void ArkeZeusSetTargetTemperature(ArkeZeusSetPoint * sp,float value) {
	sp->Temperature = float_to_hih6030_temperature(value);
}

void ArkeZeusSetTargetWind(ArkeZeusSetPoint * sp,uint8_t power) {
	sp->Power = power;
}


float ArkeZeusGetTargetHumidity(const ArkeZeusSetPoint * sp) {
	return humidity_to_float(sp->Humidity);
}

float ArkeZeusGetTargetTemperature(const ArkeZeusSetPoint * sp) {
	return hih6030_temperature_to_float(sp->Temperature);
}

uint8_t ArkeZeusGetTargetWind(const ArkeZeusSetPoint * sp) {
	return sp->Wind;
}


float ArkeZeusGetHumidity(const ArkeZeusReport * r) {
	return humidity_to_float(r->Humidity);
}

float ArkeZeusGetTemperature1(const ArkeZeusReport * r) {
	return hih6030_temperature_to_float(r->Temperature1);
}

float ArkeZeusGetTemperature2(const ArkeZeusReport * r) {

}

float ArkeZeusGetTemperature3(const ArkeZeusReport * r) {

}

float ArkeZeusGetTemperature4(const ArkeZeusReport * r) {

}
