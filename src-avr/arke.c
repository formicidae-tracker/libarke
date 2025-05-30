#include <avr/eeprom.h>
#include <avr/interrupt.h>

#include <arke-avr.h>
#include <arke-avr/systime.h>



#include <yaacl.h>


#if defined(ARKE_MY_FW_TWEAK_VERSION)
#define VERSION_LENGTH 4
#elif defined(ARKE_MY_FW_PATCH_VERSION)
#define VERSION_LENGTH 3
#else
#define VERSION_LENGTH 2
#endif

#ifndef ARKE_MY_SIZE
#error "Please define ARKE_MY_SIZE"
#else
#define ARKE_RX_IDT(ID_) ( (ARKE_MESSAGE << 9) | (ARKE_MY_CLASS << 3) | (ID_))
#define ARKE_RX_IDT_MASK (ARKE_MESSAGE_TYPE_MASK | ((~(ARKE_MY_SIZE - 1) & 0x3f) << 3) | 0x07 )
#define ARKE_RX_BROAD_MASK (ARKE_MESSAGE_TYPE_MASK | ((~(ARKE_MY_SIZE - 1) & 0x3f) << 3))
#endif

typedef enum ArkeHeartbeatStatus_e {
	NO_HEARTBEAT = 0,
	HEARTBEAT_ONCE = 1,
	HEARTBEAT_REPEAT = 2
} ArkeHeartbeatStatus_e;


#define ARKE_ERROR_FIFO_SIZE 16
#define ARKE_ERROR_FIFO_MASK 15

struct ArkeErrorFIFO_t {
	ArkeError_t data[ARKE_ERROR_FIFO_SIZE];
	uint8_t size,head,tail;
};

#define FIFO_INIT(f) do {	  \
		(f).size = 0; \
		(f).head = 0; \
		(f).tail = 0; \
	}while(0)
#define FIFO_FULL(f) ((f).size >= ARKE_ERROR_FIFO_SIZE)
#define FIFO_EMPTY(f) ((f).size == 0)
#define FIFO_HEAD(f) ((f).data[(f).head])
#define FIFO_TAIL(f) ((f).data[(f).tail])
#define FIFO_INCR_PTR(p) (p) = ((p) + 1) & ARKE_ERROR_FIFO_MASK
#define FIFO_INCREMENT_HEAD(f) do{	  \
		FIFO_INCR_PTR((f).head); \
		(f).size -= 1; \
	}while(0)
#define FIFO_INCREMENT_TAIL(f) do{	  \
		FIFO_INCR_PTR((f).tail); \
		(f).size += 1; \
	}while(0)



typedef struct ArkeData_t {
	//we limit ourselfs to 6 MOb
	yaacl_txn_t control;
	yaacl_txn_t heartbeat,error;
	yaacl_txn_t rx,broadcast;
	uint8_t rxLength;
	uint8_t * rxBuffer;
	uint8_t controlData[8];
	uint8_t heartbeatData[VERSION_LENGTH];
	uint8_t errorData[4];
	ArkeHeartbeatStatus_e heartbeatStatus;
	ArkeSystime_t heartbeatPeriod;
	ArkeSystime_t lastHeartbeat;
	uint8_t ID;
	struct ArkeErrorFIFO_t errors;
} ArkeData_t;

#define arke_prepare_rx_txn(ID_) do {	  \
		yaacl_make_std_idt(arke.rx.ID,ARKE_RX_IDT(ID_),0); \
		yaacl_make_std_mask(arke.rx.mask,ARKE_RX_IDT_MASK,0,1); \
		arke.rx.length = arke.rxLength; \
		arke.rx.data = arke.rxBuffer; \
	}while(0);

#define arke_prepare_broadcast_txn(ID_) do {	  \
		yaacl_make_std_idt(arke.broadcast.ID,ARKE_RX_IDT(ID_),0); \
		yaacl_make_std_mask(arke.broadcast.mask,ARKE_RX_BROAD_MASK,0,1); \
		arke.broadcast.length = arke.rxLength; \
		arke.broadcast.data = arke.rxBuffer; \
	}while(0);


ArkeData_t arke;

uint8_t EEMEM arkeID = 1;


void InitArke(uint8_t * rxBuffer, uint8_t length) {
	uint8_t sreg = SREG;
	cli();
	arke.ID = eeprom_read_byte(&arkeID);
	if ( arke.ID == 0 || arke.ID > 7 ) {
		arke.ID = 1;
	}
	SREG = sreg;

	arke.heartbeatData[0] = ARKE_MY_FW_MAJOR_VERSION;
	arke.heartbeatData[1] = ARKE_MY_FW_MINOR_VERSION;
#if VERSION_LENGTH > 2
	arke.heartbeatData[2] = ARKE_MY_FW_PATCH_VERSION;
#endif

#if VERSION_LENGTH > 3
	arke.heartbeatData[3] = ARKE_MY_FW_TWEAK_VERSION;
#endif
	arke.heartbeatStatus = NO_HEARTBEAT;
	arke.heartbeatPeriod = 0;

	ArkeInitSystime();

	yaacl_config_t config;
	config.baudrate = YAACL_BR_250;


	yaacl_init(&config);

	yaacl_init_txn(&(arke.control));
	yaacl_init_txn(&(arke.error));
	yaacl_init_txn(&(arke.heartbeat));
	yaacl_init_txn(&(arke.rx));
	//reserve immediatly the highest priority Mob for Network Control Handling
	yaacl_make_std_idt(arke.control.ID,0,0);
	yaacl_make_std_mask(arke.control.mask,ARKE_MESSAGE_TYPE_MASK,1,1);
	arke.control.length = 8;
	arke.control.data = &(arke.controlData[0]);

	arke.error.ID = ARKE_ERROR_REPORT;
	arke.error.length = 4;
	arke.error.data = &(arke.errorData[0]);
	arke.error.data[0] = ARKE_MY_CLASS;
	arke.error.data[1] = arke.ID;
	arke.error.data[2] = 0;
	arke.error.data[3] = 0;


	arke.heartbeat.ID = (ARKE_HEARTBEAT << 9 ) | (ARKE_MY_CLASS << 3) | arke.ID;
	arke.heartbeat.data = &(arke.heartbeatData[0]);

	arke.rxLength = length;
	arke.rxBuffer = rxBuffer;
	arke_prepare_rx_txn(arke.ID);
	arke_prepare_broadcast_txn(arke.ID);

	yaacl_listen(&arke.control);
	yaacl_listen(&arke.broadcast);
	yaacl_listen(&arke.rx);

	FIFO_INIT(arke.errors);

}

void ProcessControl() {
	uint8_t command = arke.control.ID & ARKE_SUBID_MASK;
	uint8_t dataLength = arke.control.length;
	// we are forbidden to access any yaacl function here

	if ( command == ARKE_RESET_REQUEST && dataLength == 1 ) {
		uint8_t targetID = arke.controlData[0];
		if ( targetID == 0 || targetID == arke.ID) {
			ArkeSoftwareReset();
		}
		return;
	}

	if ( command == ARKE_HEARTBEAT_REQUEST ) {
		if (dataLength == 0
		    || (dataLength == 2
		        && arke.controlData[0] == 0
		        && arke.controlData[1] == 0 ) ) {
			// single heartbeat request
			arke.heartbeatPeriod = 0;
			arke.lastHeartbeat = ArkeGetSystime();
			arke.heartbeatStatus = HEARTBEAT_ONCE;
		} else if ( dataLength == 2 ) {

			*(((uint8_t*)(&arke.heartbeatPeriod)) + 0) = arke.controlData[0];
			*(((uint8_t*)(&arke.heartbeatPeriod)) + 1) = arke.controlData[1];
			arke.lastHeartbeat = ArkeGetSystime();
			arke.heartbeatStatus = HEARTBEAT_REPEAT;
		}
		return;
	}

	if ( command == ARKE_ID_CHANGE_REQUEST
	     && dataLength == 2
	     && arke.controlData[0] == arke.ID
	     && arke.controlData[1] > 0
	     && arke.controlData[1] < 8 ) {
		arke.ID = arke.controlData[1];
		uint8_t sreg = SREG;
		cli();
		eeprom_update_byte(&arkeID,arke.ID);
		SREG = sreg;


		// needed, but normally you won't change ID on the fly.
		ArkeSoftwareReset();
		return;
	}
}

yaacl_idt_t ArkeProcessNodeMessage(uint8_t * length) {
	yaacl_txn_status_e s = yaacl_txn_status(&arke.rx);
	yaacl_idt_t ret = ARKE_NO_MESSAGE;
	if ( s == YAACL_TXN_COMPLETED ) {
		*length = arke.rx.length;
		ret = arke.rx.ID;
	}

	if ( s != YAACL_TXN_PENDING ) {
		arke_prepare_rx_txn(arke.ID);
		yaacl_listen(&arke.rx);
	}
	if ( ret != ARKE_NO_MESSAGE) {
		return ret;
	}

	s = yaacl_txn_status(&arke.broadcast);

	if ( s == YAACL_TXN_COMPLETED ) {
		*length = arke.broadcast.length;
		ret = arke.broadcast.ID;
	}

	if ( s != YAACL_TXN_PENDING ) {
		arke_prepare_broadcast_txn(arke.ID);
		yaacl_listen(&arke.broadcast);
	}

	return ret;

}

void arkeSendErrors() {
	if (FIFO_EMPTY(arke.errors) || yaacl_txn_status(&(arke.error)) != YAACL_TXN_UNSUBMITTED ) {
		return;
	}

	ArkeError_t e = FIFO_HEAD(arke.errors);
	arke.errorData[2] = e & 0x00ff;
	arke.errorData[3] = (e & 0xff00) >> 8;

	yaacl_error_e err = yaacl_send(&(arke.error));
	if (err == YAACL_ERR_MOB_OVERFLOW) {
		return;
	}

	FIFO_INCREMENT_HEAD(arke.errors);
}

yaacl_idt_t ArkeProcess(uint8_t * length) {
	yaacl_txn_status_e s = yaacl_txn_status(&(arke.control));
	if ( s == YAACL_TXN_COMPLETED ) {
		// ID and data will not change unless yaacl_txn_status is
		// called again.
		uint8_t targetID = (arke.control.ID & ARKE_NODE_CLASS_MASK) >> 3;
		if ( targetID == 0x00 || targetID == ARKE_MY_CLASS ) {
			ProcessControl();
		}
	}
	if ( s != YAACL_TXN_PENDING ) {
		// we received or had an error, we re-listen
		arke.control.length = 8;
		yaacl_listen(&(arke.control));
	}

	yaacl_idt_t toReturn = ArkeProcessNodeMessage(length);

	arkeSendErrors();

	// heartbeat
	if ( arke.heartbeatStatus == NO_HEARTBEAT
	     || yaacl_txn_status(&(arke.heartbeat)) != YAACL_TXN_UNSUBMITTED ) {
		//we are sending or just have sent an heartbeat, or we do not
		//need a periodic heartbeat
		return toReturn;
	}

	ArkeSystime_t now = ArkeGetSystime();
	if ( (now - arke.lastHeartbeat) < arke.heartbeatPeriod ) {
		// wait before sending a new heartbeat
		return toReturn;
	}

	arke.heartbeat.length = (arke.heartbeatStatus == HEARTBEAT_ONCE) ? VERSION_LENGTH : 0;

	yaacl_error_e err =  yaacl_send(&(arke.heartbeat));
	if ( err == YAACL_ERR_MOB_OVERFLOW ) {
		// no free Mob, no worries, we fo it later.
		return toReturn;
	}
	//heartbeat succesfully submitted
	if (arke.heartbeatStatus == HEARTBEAT_ONCE ) {
		arke.heartbeatStatus = NO_HEARTBEAT;
	} else {
		arke.lastHeartbeat = now;
	}

	return toReturn;
}

#define implement_sender_function(name)                                        \
	ARKE_DECLARE_SENDER_FUNCTION(name) {                                       \
		txn->ID =                                                              \
		    arke.ID | (Arke##name##ClassValue << 3) |                          \
		    ((emergency ? ARKE_HIGH_PRIORITY_MESSAGE : ARKE_MESSAGE) << 9);    \
		txn->length = sizeof(Arke##name);                                      \
		txn->data   = (uint8_t *)data;                                         \
		return yaacl_send(txn);                                                \
	}

// clang-format off

implement_sender_function(ZeusSetPoint)
implement_sender_function(ZeusReport)
implement_sender_function(ZeusStatus)
implement_sender_function(ZeusConfig)
implement_sender_function(ZeusControlPoint)
implement_sender_function(ZeusDeltaTemperature)
implement_sender_function(HeliosSetPoint)
implement_sender_function(CelaenoSetPoint)
implement_sender_function(CelaenoStatus)
implement_sender_function(CelaenoConfig)
implement_sender_function(NotusSetPoint)

    // clang-format on

    void ArkeReportError(ArkeError_t error) {
	if (FIFO_FULL(arke.errors)) {
		return;
	}
	FIFO_TAIL(arke.errors) = error;
	FIFO_INCREMENT_TAIL(arke.errors);
}
