# FORT Communication Protocol Specification

## Overview

The Communication protocol is based on CAN Bus 2.0 A. According to
the OSI model it can be described according to the folowing layers:

### Physical and Data Link Layer

[CAN Bus 2.0 A (11-bit identifier)](https://en.wikipedia.org/wiki/CAN_bus#Base_frame_format)
specification running at 250 kbit/s.

### Network Layer

The network consits of a host, usually a desktop computer, running Linux and a set of nodes (here: AVR microcontrollers).
The host acknowledges any message on the CAN bus. Each node has an unique 9 bit identifier. Each physical node is configured to respond to a certain set of CAN identifier (IDT) alone on the bus, i.e. for any message identifier, only one single node should listen that type of message.

The 11 bits of the CAN base identifier are subdivided as follows to specify any
message on the bus:

<table>
	<tr>
		<th></th>
		<th colspan="2">Message Class</th>
		<th colspan="6">Message Category</th>
		<th colspan="3">ID</th>
	</tr>
	<tr>
		<th>Bit</th>
		<td>10</td>
		<td>9</td>
		<td>8</td>
		<td>7</td>
		<td>6</td>
		<td>5</td>
		<td>4</td>
		<td>3</td>
		<td>2</td>
		<td>1</td>
		<td>0</td>
	</tr>
</table>


#### Message Class

The first two bits are used to specify the class of the message according to the following table :

|        Bits Value | Message Class                           | Node Permission        |
|------------------:|:----------------------------------------|:-----------------------|
|              0b00 | Network Control Command                 |         Reception Only |
|              0b01 | High Priority Message (error, emergency)|          Emission Only |
|              0b10 | Standard Priority Messages              | Emission and Reception |
|              0b11 | Heartbeat message                       |          Emission Only |


Except for Network Control Command and Heartbeat messages, the 9 remaning bits are used to specify a single node on the bus using message class and ID field


#### Message Category

The message class identifies the kind of message and payload according to this specifications. A physical device on the bus can receive several different classes of messages.

The message class is encoded MSB first in bits number 8...3 of the CAN IDT.

Each device on the bus has an assigned a class which corresponds to the lowest message class value that this node accepts.


#### ID field

The ID field ensures that several physically identical nodes can co-exist on the bus and be addressed independently. This field consists of the lowest 3 bits of the CAN IDT. The value 0 is reserved for broadcasting.

Each device of a given class must have a unique ID. When emmitting message, any device uses its own ID in the IDT. If two devices are required by the application to communicate independently of the host, they need to share the same ID (but obviously different classes). TODO: The last sentence does not make sense, must be: If two devices are required by the application to communicate independently of the host, they need to share the same ID (but obviously different categories).

Only the host can use the ID-value 0 to broadcast a message.


### Session Layer

The session is managed using heartbeat messages. The host uses a Heartbeat request message to ask all nodes of given class to provide a heartbeat every X milliseconds. This heartbeat is used to monitor which nodes are online. The same command is used to
request a single heartbeat to detect all nodes present on the bus.

## FORT Standard Message Category Table and Node Class ID

This table lists all possible message categories of the bus:

| Message Category | Name                                    | Node    | Node Class ID (without subID) |
|-----------------:|:----------------------------------------|:--------|:------------------------------|
|             0x3f | Reserved (Zeus)                         | Zeus    | 0x38                          |
|             0x3e | Zeus Delta Temperature Configuration    | Zeus    | 0x38                          |
|             0x3d | Zeus Control Point Report               | Zeus    | 0x38                          |
|             0x3c | Zeus Status Report                      | Zeus    | 0x38                          |
|             0x3b | Zeus Configuration                      | Zeus    | 0x38                          |
|             0x3a | Zeus Vibration Report                   | Zeus    | 0x38                          |
|             0x39 | Zeus Temperature and Humidity Report    | Zeus    | 0x38                          |
|             0x38 | Zeus Temperature and Humidity Set Point | Zeus    | 0x38                          |
|        0x36-0x37 | Reserved (Helios)                       | Helios  | 0x34                          |
|             0x35 | Helios Pulse Mode                       | Helios  | 0x34                          |
|             0x34 | Helios Set Point                        | Helios  | 0x34                          |
|             0x33 | Reserved (Celaeno)                      | Celaeno | 0x30                          |
|             0x31 | Celaeno Configuration                   | Celaeno | 0x30                          |
|             0x31 | Celaeno Status                          | Celaeno | 0x30                          |
|             0x30 | Celaeno Humidity Set Point              | Celaeno | 0x30                          |
|        0x01-0x29 | Reserved for Future Use                 | n.a     | n.a                           |
|             0x00 | Reserved for broadcast                  | all     | n.a.                          |

Any of these messages can be sent using low (0b10) or high (0b01) priority.

The following table lists the class IDs for the FORT nodes currently defined:

| Node Name | Device Class ID |
|----------:|:----------------|
|      Zeus | 0x38            |
|    Helios | 0x34            |
|   Celaeno | 0x30            |

For all of these messages, the host can use a Remote Transmission Request (RTR) (with a Data length code (DLC) field of strictly zero) to actively fetch the data if required.


## Message Specifications.


#### 0x30 Celaeno Humidity Set Point

* Host Access: Read/Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 1
  * Data fields:
	* Byte 0: Amount of humidity to be produced

#### 0x31 Celaeno Status

* Host Access: Read
* Periodically emitted by node: yes on any exceptional situation
* Payload:
  * Data Length: 3
  * Data fields:
	* Byte 0: Water level status
	  * 0x00: Functionning normally
	  * 0x01: Warning level reached
	  * 0x02: Critical level reached, humidity production disabled
	  * 0x04: Sensor readout error
	* Bytes 1-2: Fan Status, little endian
	  * B0 - B13 : Current fan RPM
	  * B14 : If set, specifies a fan aging alert
	  * B15 : If set, specifies a fan stall alert (Fan should spin but is currently not)

#### 0x32 Celaeno Configuration

* Host Access: Read/Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 8
  * Data fields:
	* Bytes 0-1: Ramp-up time [ms], little endian
	* Bytes 2-3: Ramp-down time [ms], little endian
	* Bytes 4-5: Minimum On time [ms], little endian
	* Bytes 6-7: Floating sensor debounce time [ms], little endian

#### 0x34 Helios Set Point

* Host Access: Read/Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 2
  * Data fields:
	* Byte 0: Visible light [TODO: amount?]
	* Byte 1: UV light [TODO: amount?]

#### 0x35 Helios Pulse Mode Toggle

Toggles a pulse mode where light output is a triangle wave with a few seconds period. Mainly for debug purpose.

* Host Access:  Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 0

#### 0x38 Zeus Set Point

* Host Access:  Read/Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 4
  * Data fields:
  	* Bytes 0-1: Target relative humidity, result of (relative_humidity/100.0) * 16382, little endian
	* Bytes 2-3: Target temperature, results of ( (temp+40.0) / 165.0 ) * 16382, little endian

#### 0x39 Zeus Climate Report

* Host Access:  Read
* Periodically emitted by node: yes
* Payload:
  * Data Length: 8
  * Data fields:
    	* Bits 0-13: Current relative humidity, little endian, conversion: x -> (x/16382)*100.0 in %
	* Bits 14-27: Ant temperature,  little endian, conversion: x -> (x/16382)*165 - 40.0 in °C
    	* Bits 28-39: Aux temperature 1,  2-complement on 12 bits, little endian, conversion: x -> x * 0.0625 in °C
	* Bits 40-51: Aux temperature 2,  2-complement on 12 bits, little endian, conversion: x -> x * 0.0625 in °C
	* Bits 52-63: Aux temperature 3,  2-complement on 12 bits, little endian, conversion: x -> x * 0.0625 in °C

#### 0x3a Zeus Vibration Report

This message is reserved for future use

* Host Access:  Read
* Periodically emitted by node: yes
* Payload: Undefined

#### 0x3b Zeus Configuration

* Host Access:  Read/Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 8
  * Data fields:
	* Bytes 0-3: Humidity PID Control configuration
	  * Bits 0-7: Proportional gain (P)
	  * Bits 8-15: Derivative gain (D)
	  * Bits 16-24: Integral gain (I)
	  * Bits 25-28: Proportional and Derivative divider (DIV), in power of 2. if P is bits 0..7 and DIV is bits 25-28, the final proportional gain is P/(2^DIV)
	  * Bits 29-31: Integral divider in power of 2
	  * Bytes 4-7: Temperature PID Control Configuration, same structure as above

#### 0x3c Zeus Status

* Host Access:  Read
* Periodically emitted by node: on exceptional situation
* Payload:
  * Data Length: 7
  * Data fields:
	* Byte 0: Generas Zeus Status
	  * Bit 0: Climate control loop is running (set point received)
	  * Bit 1: Climate uncontrolled for too long flag
	  * Bit 2: Target humidity cannot be reached flag
	  * Bit 3: Target temperature cannot be reached flag
  	* Bytes 1-2: Wind fan status and RPM, same structure as Celaeno Fan status
    	* Bytes 3-4: Right extraction fan status and RPM, same structure as Celaeno Fan status
    	* Bytes 3-4: Left extraction fan status and RPM, same structure as Celaeno Fan status

#### 0x3d Zeus Control Point

* Host Access:  Read
* Periodically emitted by node: yes
* Payload:
  * Data Length: 4
  * Data fields:
	* Bytes 0-1: Humidity PID Control command output, signed word little endian
	* Bytes 2-3: Temperature PID Control command output, signed word little endian


#### 0x3e Zeus Delta Temperature

* Host Access:  Read/Write
* Periodically emitted by node: never
* Payload:
  * Data Length: 8
  * Data fields:
	* Bytes 0-1: Ant Temperature delta, signed word little endian, (x) -> x*16382/165 °C
	* Bytes 2-3: Aux1 Temperature delta, signed word little endian, (x) -> x*0.0625 °C
	* Bytes 4-5: Aux2 Temperature delta, signed word little endian, (x) -> x*0.0625 °C
	* Bytes 6-7: Aux2 Temperature delta, signed word little endian, (x) -> x*0.0625 °C


## FORT Network Control Command Specification

Network Control Command IDTs are formatted differently than other messages:
* The message category field is used to target specific node classes, or 0x00
  t for broadcasts
* The ID fields is used as command specification, that's to be broadcast to all nodes of the class specified.

The following table lists the commands that are specified.
Note that for some, the implementation is not stricly required.

| Code  | Command                   | Implementation | Payload      |
|-------|---------------------------|----------------|--------------|
| 0b000 | Software Reset Request    | Required       | 1 byte       |
| 0b001 | Timestamp Synchronization | Optional       | n.a.         |
| 0b010 | Node ID Change            | Required       | 2 bytes      |
| 0b011 | Device Error Report       | Required       | 4 bytes      |
| 0b111 | Heartbeat Request         | Required       | 0 or 2 bytes |


Network commands cannot use the RTR flag.

### 0b000 Software Reset Request

Any node on the bus must implement this feature in order to reset itself after acknowledging the command.
The payload of this command is the target ID of the node that needs to perform a reset. A value of 0 resets all nodes of the
specified class(es).

* Payload:
  * Data Length: 1
  * Data fields:
	* Byte 0: Target node ID

### 0b001 Timestamp Synchronization

Not yet specified.

### 0b010 Node ID Change Request

Request a node to change its ID to a new one. This triggers a subsequent
software reset. The target ID cannot be 0.

* Payload:
  * Data Length: 2
  * Data Fields:
	* Byte 0: Old target ID, cannot be 0 (broadcasts not possible)
	* Byte 1: New target ID, should be in 1 to 7 range

### 0b011 Node Internal Error Report

These special messages are use by nodes to report important internal errors. Their main purpose is for development debugging. Any production applications should not rely on this kind of error reporting.

* Payload:
  * Data Length: 4
  * Data Fields:
	* Byte 0: Class of the device issuing the error
	* Byte 1: ID of the device issuing the error
	* Bytes 2-3: Error code, little endian word

### 0b111 Heartbeat Request

This command is used to request for the targeted nodes to transmit a heartbeat.
The Host specifies an Heartbeat period in ms, and the targeted nodes are expected to transmit a heartbeat periodically after the acknowledgment. If the period is 0 or simply omitted, a single heartbeat requested from the targeted nodes.

* Payload:
  * Data Length: 0 or 2
  * Data Fields:
	* Bytes 0-1: Heartbeat period [ms], little endian

## Heartbeat Message Specification

Heartbeat messages are issued by nodes to monitor their online status.
The CAN IDT consists of 0b11 followed by the node's unique ID (Node Class + ID).
In case the heartbeat is emitted following a single heartbeat request (node pinging/enumeration), the heartbeat must contain the firmware version information.
In case of periodically sent heartbeats, the nodes must not transmit this information.

* Payload:
  * Data Length 0, 2, 3 or 4.
  * Data Field :
  	* Byte 0 : Major Version number (required if any versionning is transmitted)
  	* Byte 1 : Minor Version number (required if any versionning is transmitted)
  	* Byte 2 : Patch Version number (optional)
  	* Byte 3 : Tweak Version number (optional)
