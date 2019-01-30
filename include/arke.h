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
	ARKE_ZEUS_VIBRATION_REPORT = 0x40,
	ARKE_ZEUS_CONFIG = 0x41,
	ARKE_ZEUS_STATUS = 0x42,
	ARKE_ZEUS_CONTROL_POINT = 0x43,
	ARKE_HELIOS_SET_POINT = 0x34,
	ARKE_HELIOS_PULSE_MODE = 0x35,
	ARKE_CELAENO_SET_POINT = 0x30,
	ARKE_CELAENO_STATUS = 0x31,
	ARKE_CELAENO_CONFIG = 0x32
} ArkeMessageClass;


struct ArkeZeusSetPoint_t {
	uint16_t Humidity;
	uint16_t Temperature;
	uint8_t  Wind;
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


struct ArkePDConfig_t {
	uint8_t DeadRegion;
	uint8_t ProportionalMult;
	uint8_t DerivativeMult;
	uint8_t ProportionalDivPower2:4;
	uint8_t DerivativeDivPower2:4;
} __attribute__((packed));

typedef struct ArkePDConfig_t ArkePDConfig;

struct ArkeZeusConfig_t {
	ArkePDConfig Humidity;
	ArkePDConfig Temperature;
} __attribute__((packed));

typedef struct ArkeZeusConfig_t ArkeZeusConfig;

#define ARKE_FAN_AGING_ALERT (0x4000)
#define ARKE_FAN_STALL_ALERT (0x8000)
#define ARKE_FAN_RPM_MASK (0x3fff)

#define ArkeFanAging(status) (((status).fanStatus & ARKE_FAN_AGING_ALERT) != 0x0000)
#define ArkeFanStall(status) (((status).fanStatus & ARKE_FAN_STALL_ALERT) != 0x0000)
#define ArkeFanRPM(status) ( (status).fanStatus & ARKE_FAN_RPM_MASK )

typedef uint16_t ArkeFanStatus;

struct ArkeZeusStatus_t {
	uint8_t       Status;
	ArkeFanStatus Fan[2];
} __attribute__((packed));

#define ARKE_ZEUS_IDLE 0x00
#define ARKE_ZEUS_ACTIVE 0x01
#define ARKE_ZEUS_CLIMATE_UNCONTROLLED_WD 0x02

typedef struct ArkeZeusStatus_t ArkeZeusStatus;

struct ArkeZeusControlPoint_t {
	int16_t Humidity;
	int16_t Temperature;
} __attribute__((packed));

typedef struct ArkeZeusControlPoint_t ArkeZeusControlPoint;

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
	uint8_t  waterLevel;
	ArkeFanStatus fanStatus;
} __attribute__((packed));
typedef struct ArkeCelaenoStatus_t  ArkeCelaenoStatus;


typedef enum ArkeCelaenoWaterLevel_e {
	ARKE_CELAENO_NOMINAL = 0,
	ARKE_CELAENO_WARNING = (1 << 0),
	ARKE_CELAENO_CRITICAL = (1 << 1),
	ARKE_CELAENO_RO_ERROR = (1 << 2)
} ArkeCelaenoWaterLevel;

#define ArkeCelaenoWaterNominal(status) ( (status).waterLevel == 0 )
#define ArkeCelaenoWaterWarning(status) ( (status).waterLevel == ARKE_CELAENO_WARNING )
#define ArkeCelaenoWaterCritical(status) ( ((status).waterLevel & ~(ARKE_CELAENO_CRITICAL | ARKE_CELAENO_RO_ERROR) ) == ARKE_CELAENO_WARNING )
#define ArkeCelaenoWaterHasRoError(status) (((status).waterLevel & ARKE_CELAENO_RO_ERROR) != 0x00)


struct ArkeCelaenoConfig_t {
	uint16_t RampUpTimeMS;
	uint16_t RampDownTimeMS;
	uint16_t MinOnTimeMS;
	uint16_t DebounceTimeMS;
} __attribute__((packed));

typedef struct ArkeCelaenoConfig_t ArkeCelaenoConfig;



#ifdef __cplusplus
}
#endif //__cplusplus
