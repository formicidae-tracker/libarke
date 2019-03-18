#pragma once


#include "arke-avr/systime.h"
#include <arke.h>
#include "yaacl.h"

#include "inttypes.h"

void InitArke(uint8_t * rxBuffer, uint8_t length);

#define ARKE_NO_MESSAGE 0

yaacl_idt_t ArkeProcess(uint8_t * length);


uint8_t ArkeMyID();

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


typedef uint16_t ArkeError_t;

void ArkeReportError(ArkeError_t error);

#define ARKE_DECLARE_SENDER_FUNCTION(name) \
	yaacl_error_e ArkeSend ## name(yaacl_txn_t * txn, bool emergency,const Arke ## name * data)



ARKE_DECLARE_SENDER_FUNCTION(ZeusSetPoint);
ARKE_DECLARE_SENDER_FUNCTION(ZeusReport);
ARKE_DECLARE_SENDER_FUNCTION(ZeusConfig);
ARKE_DECLARE_SENDER_FUNCTION(ZeusStatus);
ARKE_DECLARE_SENDER_FUNCTION(ZeusControlPoint);
ARKE_DECLARE_SENDER_FUNCTION(ZeusDeltaTemperature);
ARKE_DECLARE_SENDER_FUNCTION(HeliosSetPoint);
ARKE_DECLARE_SENDER_FUNCTION(CelaenoSetPoint);
ARKE_DECLARE_SENDER_FUNCTION(CelaenoStatus);
ARKE_DECLARE_SENDER_FUNCTION(CelaenoConfig);

#define ArkeZeusSetPointClassValue ARKE_ZEUS_SET_POINT
#define ArkeZeusReportClassValue ARKE_ZEUS_REPORT
#define ArkeZeusStatusClassValue ARKE_ZEUS_STATUS
#define ArkeZeusConfigClassValue ARKE_ZEUS_CONFIG
#define ArkeZeusControlPointClassValue ARKE_ZEUS_CONTROL_POINT
#define ArkeZeusDeltaTemperatureClassValue ARKE_ZEUS_DELTA_TEMPERATURE
#define ArkeHeliosSetPointClassValue ARKE_HELIOS_SET_POINT
#define ArkeCelaenoSetPointClassValue ARKE_CELAENO_SET_POINT
#define ArkeCelaenoStatusClassValue   ARKE_CELAENO_STATUS
#define ArkeCelaenoConfigClassValue   ARKE_CELAENO_CONFIG
#define ARKE_MESSAGE_STRUCT_TO_CLASS(name) ( name ## ClassValue)
