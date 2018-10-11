#pragma once

#if defined(__AVR_ATmega16m1__) \
	|| defined(__AVR_ATmega32m1__) \
	|| defined(__AVR_ATmega64m1__)
#define START_TIMER0_1ms() do{	  \
		TCCR0A = _BV(WGM01); \
		TCCR0B = _BV(CS01) | _BV(CS00); \
	}while(0)
#define OCIE0A_ENABLE() TIMSK0 = _BV(OCIE0A)
#define OCIE0A_DISABLE() TIMSK0 = _BV(OCIE0A)
#define LIBARKE_TIMER0_COMPA_VECT TIMER0_COMPA_vect
#elif defined(__AVR_AT90CAN32__) \
	|| defined(__AVR_AT90CAN64__) \
	|| defined(__AVR_AT90CAN128__)
#define START_TIMER0_1ms() do{	  \
		TCCR0A = _BV(WGM01) | _BV(CS01) | _BV(CS00); \
	}while(0)
#define OCIE0A_ENABLE() TIMSK0 = _BV(OCIE0A)
#define OCIE0A_DISABLE() TIMSK0 = _BV(OCIE0A)
#define LIBARKE_TIMER0_COMPA_VECT TIMER0_COMP_vect
#else
#error "Unssuported AVR device"
#endif
