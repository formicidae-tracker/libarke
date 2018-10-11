#pragma once


#include "arke.h"


#include "inttypes.h"

typedef uint16_t ArkeSystime_t;

ArkeSystime_t ArkeGetSystime();

void InitArke();

void ArkeProcess();

void ArkeSoftwareReset();

#include <avr/wdt.h>

#define implements_ArkeSoftwareReset()	  \
	void ArkeSoftwareReset() { \
		wdt_enable(WDTO_15MS); \
		for(;;){} \
	} \
	void wdt_init() __attribute__((naked)) __attribute__((section(".init3"))); \
	void wdt_init() { \
		MCUSR = 0; \
		wdt_disable(); \
	}
