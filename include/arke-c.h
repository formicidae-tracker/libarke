#pragma once

#include "arke.h"

#ifdef __cplusplus
extern "C" {
#endif


	void ArkeZeusSetTargetHumidity(ArkeZeusSetPoint * sp,float humidity);
	void ArkeZeusSetTargetTemperature(ArkeZeusSetPoint * sp,float humidity);
	void ArkeZeusSetTargetWind(const ArkeZeusSetPoint * sp,uint8_t power);

	float ArkeZeusGetTargetHumidity(const ArkeZeusSetPoint * sp);
	float ArkeZeusGetTargetTemperature(const ArkeZeusSetPoint * sp);
	uint8_t ArkeZeusGetTargetWind(const ArkeZeusSetPoint * sp);



	float ArkeZeusGetHumidity(const ArkeZeusReport * r);
	float ArkeZeusGetTemperature1(const ArkeZeusReport * r);
	float ArkeZeusGetTemperature2(const ArkeZeusReport * r);
	float ArkeZeusGetTemperature3(const ArkeZeusReport * r);
	float ArkeZeusGetTemperature4(const ArkeZeusReport * r);







#ifdef __cplusplus
}
#endif
