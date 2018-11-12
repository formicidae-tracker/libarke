#include "arke-avr/systime.h"


#include <avr/io.h>
#include <avr/interrupt.h>

#include "config.h"

volatile ArkeSystime_t systime;

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
