#pragma once


#include "arke.h"
#include "yaacl.h"

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


#define ARKE_DECLARE_SENDER_FUNCTION(name) \
	yaacl_error_e ArkeSend ## name(yaacl_txn_t * txn,bool rtr, bool emergency, uint8_t subID,const Arke ## name * data)

ARKE_DECLARE_SENDER_FUNCTION(ZeusSetPoint);
ARKE_DECLARE_SENDER_FUNCTION(ZeusReport);
ARKE_DECLARE_SENDER_FUNCTION(HeliosSetPoint);
ARKE_DECLARE_SENDER_FUNCTION(CelaenoSetPoint);
ARKE_DECLARE_SENDER_FUNCTION(CelaenoStatus);

#define ArkeZeusSetPointClassValue ARKE_ZEUS_SET_POINT
#define ArkeZeusReportClassValue ARKE_ZEUS_REPORT
#define ArkeHeliosSetPointClassValue ARKE_HELIOS_SET_POINT
#define ArkeCelaenoSetPointClassValue ARKE_CELAENO_SET_POINT
#define ArkeCelaenoStatusClassValue   ARKE_CELAENO_STATUS
#define ARKE_MESSAGE_STRUCT_TO_CLASS(name) ( name ## ClassValue)
