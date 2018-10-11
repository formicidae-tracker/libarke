#pragma once

#include "inttypes.h"

#ifdef __cplusplus
extern "C"{
#endif //__cplusplus


typedef enum ArkeMessageType_e {
	ARKE_NETWORK_CONTROL_COMMAND = 0x00,
	ARKE_HIGH_PRIORITY_MESSAGE = 0x01,
	ARKE_MESSAGE = 0x02,
	ARKE_HEARTBEAT = 0x03,
	ARKE_MESSAGE_TYPE_MASK = 0x03 << 9
} ArkeMessageType;

typedef enum ArkeNodeClass_e {
	ARKE_BROADCAST = 0x0,
	ARKE_ZEUS = 0x38,
	ARKE_HELIOS = 0x34,
	ARKE_CELAENO = 0x30,
	ARKE_NODE_CLASS_MASK = 0x3f << 3
} ArkeNodeClass;


typedef enum ArkeNetworkCommand_e {
	ARKE_RESET_REQUEST = 0x00,
	ARKE_SYNCHRONISATION = 0x01,
	ARKE_HEARTBEAT_REQUEST = 0x07,
	ARKE_SUBID_MASK = 0x07
} ArkeNetworkCommand;

//#define ARKE_SUBID_MASK ARKE_HEARTBEAT_REQUEST

typedef enum ArkeMessageClass_e {
	ARKE_ZEUS_SET_POINT = 0x38,
	ARKE_ZEUS_REPORT = 0x39,
	ARKE_HELIOS_SET_POINT = 0x34,
	ARKE_CELAENO_SET_POINT = 0x30,
	ARKE_CELAENO_STATUS = 0x31
} ArkeMessageClass;


struct ArkeZeusSetPoint_t {
	uint16_t Humidity;
	uint16_t Temperature;
} __attribute__((packed));
typedef struct ArkeZeusSetPoint_t ArkeZeusSetPoint;

struct ArkeZeusReport_t {
	uint16_t Humidity:14;
	uint16_t Temperature1:14;
	uint16_t Temperature2:12;
	uint16_t Temperature3:12;
	uint16_t Temperature4:12;
} __attribute__((packed));
typedef struct ArkeZeusReport_t ArkeZeusReport;

struct ArkeHeliosSetPoint_t {
	uint8_t Visible;
	uint8_t UV;
} __attribute__((packed));
typedef struct ArkeHeliosSetPoint_t ArkeHeliosSetPoint;


struct ArkeCelaenoSetPoint_t {
	uint8_t Power;
} __attribute__((packed));
typedef struct ArkeCelaenoSetPoint_t  ArkeCelaenoSetPoint;

struct ArkeCelaenoStatus_t {
	uint8_t  Level;
	uint16_t FanSpeed;
} __attribute__((packed));
typedef struct ArkeCelaenoStatus_t  ArkeCelaenoStatus;




#ifdef __cplusplus
}
#endif //__cplusplus