#pragma once

// #ifdef __cplusplus
// extern "C"{
// #endif //__cplusplus


typedef enum ArkeMessageClass_e {
	ARKE_NETWORK_CONTROL_COMMAND = 0b00,
	ARKE_HIGH_PRIORITY_MESSAGE = 0b01,
	ARKE_MESSAGE = 0b10,
	ARKE_HEARTBEAT = 0b11,
	ARKE_MESSAGE_CLASS_MASK = 0x3 << 8
} ArkeMessage;

typedef enum ArkeNodeClass_e {
	ARKE_BROADCAST = 0x0,
	ARKE_ZEUS = 0x01,
	ARKE_HELIOS = 0x09,
	ARKE_CELAENO = 0x0d,
	ARKE_NODE_CLASS_MASK = 0x3f << 3
} ArkeNodeClass;


typedef enum ArkeNetworkCommand_e {
	ARKE_RESET_REQUEST = 0b000,
	ARKE_SYNCHRONISATION = 0b001,
	ARKE_HEARTBEAT = 0b111,
	ARKE_SUBID_MASK = 0b111
} ArkeNetworkCommand;


typedef enum ArkeMessageClass_e {
	ARKE_ZEUS_REPORT = 0x01,
	ARKE_ZEUS_SET_POINT = 0x02,
	ARKE_HELIOS_SET_POINT = 0x09,
	ARKE_CELAENO_SET_POINT = 0x0d,
	ARKE_CELAENO_WATER_LEVEL = 0x0e,
} ArkeMessageClass;


#ifdef __cplusplus
}
#endif //__cplusplus
