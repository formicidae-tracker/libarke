# FORT Communication Protocol Specification

## Overview

The Communication protocol is based upon CAN Bus 2.0 A. According to
the OSI model it could be described has the following

### Physical and Data Link Layer

[CAN Bus 2.0 A (11-bit identifier)](https://en.wikipedia.org/wiki/CAN_bus#Base_frame_format)
specification running at 250 kbit/s.

### Network Layer

The network consits of an host, usually a desktop computer running
Linux, and a set of nodes. The host listen and acknowledge any message
on the CAN bus. Each node has an unique 9 bit identifier. Each
physical node is configured to respond to a certain set of CAN
identifier alone on the bus, i.e. for any message identifier, only one
single node should listen that type of message.

The CAN base identifier 11 bit are subdivided as follow to specify any
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


### Message Category

The first two bits are used to specify the class of the
message according to the following table :

|        Bits Value | Message Class                           | Node Permission        |
|------------------:|:----------------------------------------|:-----------------------|
|              0b00 | Network Control Command                 |           Receive Only |
|              0b01 | High Priority Message (error, emergency)|          Emission Only |
|              0b10 | Standard Priority Messages              | Emission and Reception |
|              0b11 | Heartbeat message                       |          Emission Only |


At the exception of Network Control Command and Heartbeat messages,
the 9 remaning bits are used to specify a single node on the bus,
using message class and ID field

### Message Class

The message class identify the kind of message and payload accordingly
to this specifications. A physical device on the bus can receive
several different class of message.

The message class is encoded MSB first in bits 8...3 of the CAN IDT.

Each device on the bus has an assigned class which correspond to the
lowest message class value this node accepts.

### ID field

The ID field is here to ensure that several physically identical node
could co-exist on the bus, and that we could address all of them
independantly. It consist of the lowest 3 bits of the CAN IDT. The
value of 0 is reserved for broadcasting.

Each device of a given class must have a unique ID. When emmitting
message, any device should use its own ID in the IDT. If the
application needs two device to communicate indepandtly of the host,
they have to share the same ID (but obviously different class).

Only host can use the 0 value for broadcasting a message.


### Session Layer

The session could be managed using heartbeat message. The host can use
a Heartbeat request message to ask for all nodes of given class to
provide a heartbeat every X milliseconds. This heartbeat could be used
to monitor which nodes are online. The same command could be used to
specify a single heartbeat to enumerate all present node on the bus.

## FORT Standard Message Category Table and Node Class ID

This table gives all message categories that could be found on the bus.

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
|             0x30 | Celaeno Set Point                       | Celaeno | 0x30                          |
|        0x01-0x29 | Reserved for Future Use                 | n.a     | n.a                           |
|             0x00 | Reserved for broadcast                  | all     | n.a.                          |

Any of this messages can be sent using low (0b10) and high (0b01)
priority.

This gives the following class ID for the FORT node currently defined :

| Node Name | Device Class ID |
|----------:|:----------------|
|      Zeus | 0x38            |
|    Helios | 0x34            |
|   Celaeno | 0x30            |

For all of this messages, the host can use a RTR request (with a DLC
of strictly zero) to actively fetch the data if required.

## Message Specifications.


#### 0x30 Celaeno Humidity Set Point

* Host Access :  Read/Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 1
  * Data fields:
	* Byte 0: Amount of humidity to be produced.

#### 0x31 Celaeno Status

* Host Access :  Read
* Periodically emitted by node : yes on any exceptional situation.
* Payload :
  * Data Length: 3
  * Data fields:
	* Byte 0: Water Level Status
	  * 0x00: Functionning Normally
	  * 0x01: Warning Level Reached
	  * 0x02: Critical Level Reached, Humidity Production Disabled.
	  * 0x04: Sensor readout error.
	* Bytes 1-2: Fan Status, Little Endian
      * B0 - B13 : Fan Current RPM
	  * B14 : If set, specifies a fan aging alert
	  * B15 : If set, specifies a fan stall alert (Fan should spin but it is currently)

#### 0x32 Celaeno Configuration

* Host Access :  Read/Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 8
  * Data fields:
	* Bytes 0-1: Ramp Up time, in ms, Little Endian
	* Bytes 2-3: Ramp Down time, in ms, Little Endian
	* Bytes 4-5: Minimum On time, in ms, Little Endian
	* Bytes 6-7: Floating Sensor Debounce Time, in ms, Little Endian

#### 0x34 Helios Set Point

* Host Access :  Read/Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 2
  * Data fields:
	* Byte 0: Visible Light Amount
	* Byte 1: UV Light Amount

#### 0x35 Helios Pulse Mode Toggle

Toggles a pulse mode where light output is a triangle wave with a few
seconds period. Mainly here for debug purpose.

* Host Access :  Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 0


#### 0x38 Zeus Set Point

* Host Access :  Read/Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 4
  * Data fields:
    * Bytes 0-1: Target Relative Humidity, result of (relative_humidity/100.0) * 16382, Little Endian
	* Bytes 2-3: Target Temperature, results of ( (temp+40.0) / 165.0 ) * 16382, Little Endian

#### 0x39 Zeus Climate Report

* Host Access :  Read
* Periodically emitted by node : yes
* Payload :
  * Data Length: 8
  * Data fields:
    * Bits 0-13: Current Relative Humidity, Little Endian, conversion: x -> (x/16382)*100.0 in %
	* Bits 14-27: Ant Temperature,  Little Endian, conversion: x -> (x/16382)*165 - 40.0 in °C
    * Bits 28-39: Aux Temperature 1,  2-complement on 12 bits, Little Endian, conversion: x -> x * 0.0625 in °C
	* Bits 40-51: Aux Temperature 2,  2-complement on 12 bits, Little Endian, conversion: x -> x * 0.0625 in °C
	* Bits 52-63: Aux Temperature 3,  2-complement on 12 bits, Little Endian, conversion: x -> x * 0.0625 in °C

#### 0x3a Zeus Vibration Report

This message is reserved for future use

* Host Access :  Read
* Periodically emitted by node : yes
* Payload : Undefined

#### 0x3b Zeus Configuration

* Host Access :  Read/Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 8
  * Data fields:
	* Bytes 0-3: Humidity PID Control Configuration
	  * Bits 0-7: Proportional Gain
	  * Bits 8-15: Derivative Gain
	  * Bits 16-24: Integral Gain
	  * Bits 25-28: Proportional and Derivative divider, in power of 2. if p is bits 0..7 and d is bits 25-28, final gain is p/(2^d)
	  * Bits 29-31: Integral divider in power of 2
    * Bytes 4-7: Temperature PID Control Configuration, same structure than above


#### 0x3c Zeus Status

* Host Access :  Read
* Periodically emitted by node : on exceptional situation
* Payload :
  * Data Length: 7
  * Data fields:
	* Byte 0: Generas Zeus Status
	  * Bit 0: Climate Control Loop is running (set point received)
	  * Bit 1: Climate Uncontrolled for too long flag
	  * Bit 2: Target Humidity cannot be reached flag
	  * Bit 3: Target Temperature cannot be reached flag
    * Bytes 1-2: Wind fan status and RPM, same structure than Celaeno Fan status
    * Bytes 3-4: Right Extraction fan status and RPM, same structure than Celaeno Fan status
    * Bytes 3-4: Left Extraction fan status and RPM, same structure than Celaeno Fan status


#### 0x3d Zeus Control Point

* Host Access :  Read
* Periodically emitted by node : yes
* Payload :
  * Data Length: 4
  * Data fields:
	* Bytes 0-1: Humidity PID Control command output, signed word Little Endian
	* Bytes 2-3: Temperature PID Control command output, signed word Little Endian


#### 0x3e Zeus Delta Temperature

* Host Access :  Read/Write
* Periodically emitted by node : never
* Payload :
  * Data Length: 8
  * Data fields:
	* Bytes 0-1: Ant Temperature delta, signed word little endian, (x) -> x*16382/165 °C
	* Bytes 2-3: Aux1 Temperature delta, signed word little endian, (x) -> x*0.0625 °C
	* Bytes 4-5: Aux2 Temperature delta, signed word little endian, (x) -> x*0.0625 °C
	* Bytes 6-7: Aux2 Temperature delta, signed word little endian, (x) -> x*0.0625 °C


## FORT Network Control Command Specification

Network Control Command IDT are formatted differently than other messages:
* The category field is used to target specific node class, or 0x00
  to broadcast to all classes
* The ID fields is used as a command specification, all node of
  the same class are broadcasted.

The following command are specified, note that for some their
implementation is not necesarly required

| Code  | Command                   | Implementation | Payload      |
|-------|---------------------------|----------------|--------------|
| 0b000 | Software Reset Request    | Required       | 1 byte       |
| 0b001 | Timestamp Synchronization | Optional       | n.a.         |
| 0b010 | Node ID Change            | Required       | 2 bytes      |
| 0b011 | Device Error Report       | Required       | 4 bytes      |
| 0b111 | Heartbeat Request         | Required       | 0 or 2 bytes |


Network command cannot use the RTR flag

### 0b000 Software Reset Request

Any node on the bus should implement this feature and reset itself
after acknowledgement. The payload of this command is the target ID
that need to perform the reset. A value of 0 resets all nodes of the
choosen class(es)

* Payload:
  * Data Length: 1
  * Data fields:
	* Byte 0: Target Node ID

### 0b001 Timestamp Synchronization

Not specified yet.

### 0b010 Node ID Change Request

Request a node to change its ID to a new one. It will trigger a
software reset immediatly after. The target ID cannot be 0.

* Payload:
  * Data Length: 2
  * Data Fields:
	* Byte 0: Target old ID, 0 does not broadcast
	* Byte 1: Target new ID, should be in 1 to 7 range

### 0b011 Node Internal Error Report

These special messages are use by nodes to report important internal
errors. There main use is for development debug and any production 
applications should not rely on these kind of error reporting.

* Payload:
  * Data Length: 4
  * Data Fields:
	* Byte 0: Class of the device issuing the error
	* Byte 1: ID of the device issuing the error
	* Bytes 2-3: Error Code, little endian word

### 0b111 Heartbeat Request

This command is used to request for the targeted nodes to transmit an
heartbeat. The Host specifies an Heartbeat period in ms, and the
targeted nodes are expected to transmit an heartbeat periodically
after acknowledgment. If the period is 0 or simply omitted, a single
heartbeat request should be sent by the targeted nodes.

* Payload
  * Data Length: 0 or 2
  * Data Fields:
	* Bytes 0-1: Heartbeat period, in ms, little endian

## Heartbeat Message Specification

Heartbeat message are issued by nodes to monitor all the nodes live
status. The CAN IDT would consist of 0b11 followed by the node unique
ID (Node Class + ID). In the case the heartbeat is emitted following a
single heartbeat request (node pinging/enumeration), the heartbeat
should contain a firmware version information. In case of heartbeat
periodically sent, the nodes should not transmit this information

* Data Length 0, 2, 3 or 4.
* Data Field :
  * Byte 0 : Major Version number (required if any versionning is transmitted)
  * Byte 1 : Minor Version number (required if any versionning is transmitted)
  * Byte 2 : Patch Version number (optional)
  * Byte 3 : Tweak Version number (optional)
