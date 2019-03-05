# FORT Communication Protocol Specification

## Overview

The Communication protocol is built upon CAN Bus 2.0 A. According to
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
single node could emit and should listen that type of message.

The CAN base identifier 11 bit are subdicided as follow to specify any
message on the bus:

<table>
	<tr>
		<th></th>
		<th colspan="2">Message Class</th>
		<th colspan="6">Message Category</th>
		<th colspan="3">subID</th>
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


The first two bits are used to specify the class of the
message according to the following table :

|        Bits Value | Message Class                           | Node Permission        |
|------------------:|:----------------------------------------|:-----------------------|
|              0b00 | Network Control Command                 |           Receive Only |
|              0b01 | High Priority Message (error, emergency)|          Emission Only |
|              0b10 | Standard Priority Messages              | Emission and Reception |
|              0b11 | Heartbeat message                       |          Emission Only |

Except for control messages the following 9 bit are subdivided in two
subset, a 6 bit message category (MSB) and a 3 bit subID (LSB). As
mentionned before only one single node over the bus should be able to
handle an unique IDT, therefore two physically identicate board should
be used on the bus, they could use the subID field to specify their
specific message. The special subID 0b000 is used for broadcasting,
and therefore up to 7 identical device can work on the bus at the same
time.

Any node on the bus can manage several message category, as it will be
describe furtherly. The unique 9-bit identifier of a node consist of
its lowest Message Category he handles and its subID.

Network control Message are a bit more perticular. They are always
broadcasted the message category is used to specify all or a specifc
type, and the subID is used to convey the actual command.




### Session Layer

The session could be managed using heartbeat message. The host can use
a Heartbeat request message to ask for all or for a specific node to
provide a heartbeat every X milliseconds. This heartbeat could be used
to monitor which node are available. The same command could be used to
specify a single heartbeat to enumerate all present node on the bus.



## FORT Standard Message Category Table and Node Class ID

This table gives all message categories that could be found on the bus.

| Message Category | Name                                          | Node         | Node Class ID (without subID) |
|-----------------:|:----------------------------------------------|:-------------|:------------------------------|
|        0x39-0x3f | Reserved (Zeus)                               | Zeus         | 0x38                          |
|             0x39 | Zeus Temperature and Humidity Report          | Zeus         | 0x38                          |
|             0x38 | Zeus Temperature and Humidity Set Poiny       | Zeus         | 0x38                          |
|        0x35-0x37 | Reserved (Helios)                             | Helios       | 0x34                          |
|             0x34 | Helios Illumination Set Point                 | Helios       | 0x34                          |
|        0x32-0x33 | Reserved (Celaeno)                            | Celaeno      | 0x30                          |
|             0x31 | Celaeno Set Point                             | Celaeno      | 0x30                          |
|             0x30 | Celaeno Status                                | Celaeno      | 0x30                          |
|        0x04-0x29 | For future use                                | n.a          | n.a                           |
|             0x00 | Reserved for broadcast                        | all          | n.a.                          |
Any of this messages can be sent using low (0b10) and high (0b01)
priority.

This gives the following class ID for the FORT nodes currently
available :

| Node Name | Device Class ID |
|----------:|:----------------|
|      Zeus | 0x38            |
|    Helios | 0x34            |
|   Celaeno | 0x30            |

For all messages, the host can use a RTR request on the device to
actively fetch the data if he requires it and this data is not sent
periodically.

## Message Specifications.

#### 0x39 Zeus Temperature and Humidity Report

* Access :  Read Only
* Periodically emitted by node : yes (frequency not yet determined)
* Write : n.a.
* Read :
  * Data Length: to be determined.
  * Data Fields: to be determined.


#### 0x38 Zeus Temperature and Humidity Set Point

* Access :  Read/Write
* Periodically emitted by node : never
* Read/Write :
  * Data Length: to be determined.
  * Data fields: to be determined.


#### 0x34 Helios Illumination Set Point

* Access :  Read/Write
* Periodically emitted by node : never
* Read/Write :
  * Data Length: 2
  * Data fields:
	* Byte 0: Visible Light Amount level from none to max (255).
	* Byte 1: UV Light Level

Note, unstar the previous system, the infrared light pulse duration is
solely managed by the framegrabber.


#### 0x31 Celaeno Humidity Production Set Point

* Access :  Read/Write
* Periodically emitted by node : never
* Read/Write :
  * Data Length: 1
  * Data fields:
	* Byte 0: Amount of humidity currently produced.


#### 0x30 Celaeno Status

* Access :  Read
* Periodically emitted by node : yes on exceptional situation
  (warning / critical level reached).
* Write: n.a.
* Read :
  * Data Length: 3
  * Data fields:
	* Byte 0: Water Level Status
	  * 0x00: Functionning Normally
	  * 0x01: Warning Level Reached
	  * 0x02: Critical Level Reached, Humidity Production Disabled.
	  * 0x04: Incoherent sensor readout.
  * Byte 1-2: Fan Status
      * B0 - B13 : Fan Current RPM
	  * B14 : If set, specifies a fan aging alert
	  * B15 : If set, specifies a fan stall alert (Fan should spin but is not)


## FORT Network Control Command Specification

As mentionned previously, Network Control Command are specified differently:
* The category field is used to target specific node class, or 0x00
  to broadcast to all classes
* The subID field is used as a command specification, all node of
  the same class are broadcasted.

The following command are specified, note that for some their
implementation is not necesarly required

| Code  | Command                   | Implementation |
|-------|---------------------------|----------------|
| 0b000 | Software Reset Request    | Required       |
| 0b001 | Timestamp Synchrnoization | Optional       |
| 0b010 | Node ID Change            | Required       |
| 0b111 | Heartbeat Request         | Required       |

As these message are broadcasted, the RTR beat should always be set to
zero.


### 0b000 Software Reset Request

Any node on the bus should implement this feature and reset itself
after acknowledgement. The expected data length is 0 byte.


### 0b001 Timestamp Synchronisation

Not specified yet, be if timestamped data would be required on the
bus, this high priority frame should be used for re-synchronizing
node clocks to the Host one.

### 0b111 Heartbeat Request

This command is used to request for the targeted nodes to transmit
an heartbeat. The Host specifies an Heartbeat period in ms, and the
targeted nodes are expected to transmit an heartbeat periodically
after acknowledgment. If the period is 0 or simply omitted, a single
heartbeat request should be sent by the targeted nodes.

* Data Length: 0 or 2
* Data Field:
  * Byte 0: Heartbeat period LSB
  * Byte 1: Heartbeat period MSB

## Heartbeat Message Specification

Heartbeat message are issued by nodes to monitor all the nodes
status. The CAN IDT would consist of 0b11 followed by the node
unique ID (Node Class + subID). In the case the heartbeat is emitted
following a single heartbeat request (node enumeration), the
heartbeat should contain a firmware version information. In case of
heartbeat periodically sent, the nodes could avoid to transmit this
information (but are allowed to transmit). Firmware version consist of
2 to 4 number, corresponding to Major,Minor,Patch and Tweak numbering

* Data Length 0, 2, 3 or 4.
* Data Field :
  * Byte 0 : Major Version number (required if any versionning is transmitted)
  * Byte 1 : Minor Version number (required if any versionning is transmitted)
  * Byte 2 : Patch Version number (optional)
  * Byte 3 : Tweak Version number (optional)
