#include "arke-avr.h"

#include <avr/io.h>
#include <avr/interrupt.h>

volatile ArkeSystime_t systime;

#if defined(__AVR_ATmega16m1__)
#define START_TIMER0_1ms() do{	  \
		TCCR0A = _BV(WGM01); \
		TCCR0B = _BV(CS01) | _BV(CS00); \
	}while(0)
#define OCIE0A_ENABLE() TIMSK0 = _BV(OCIE0A)
#define OCIE0A_DISABLE() TIMSK0 = _BV(OCIE0A)
#define LIBARKE_TIMER0_COMPA_VECT TIMER0_COMPA_vect
#elif defined(__AVR_AT90CAN128__)
#define START_TIMER0_1ms() do{	  \
		TCCR0A = _BV(WGM01) | _BV(CS01) | _BV(CS00); \
	}while(0)
#define OCIE0A_ENABLE() TIMSK0 = _BV(OCIE0A)
#define OCIE0A_DISABLE() TIMSK0 = _BV(OCIE0A)
#define LIBARKE_TIMER0_COMPA_VECT TIMER0_COMP_vect
#else
#error "Unssuported AVR device"
#endif

void ArkeInitSystime() {
	systime = 0;

	//frequency 16MHz/ ( N * (1+OCRnA))
	OCR0A = 249;
	OCIE0A_ENABLE();
	START_TIMER0_1ms();
	sei();
}

ISR(LIBARKE_TIMER0_COMPA_VECT) {
	++systime;
}


ArkeSystime_t ArkeGetSystime() {
	OCIE0A_DISABLE();
	ArkeSystime_t res = systime;
	OCIE0A_ENABLE();
	return res;
}
