package arke

import (
	"encoding/binary"
	"fmt"
	"time"

	socketcan "github.com/atuleu/golang-socketcan"
)

func MakeResetRequest(c NodeClass, ID NodeID) socketcan.CanFrame {
	return socketcan.CanFrame{
		ID:       MakeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(ResetRequest)),
		Dlc:      1,
		Extended: false,
		RTR:      false,
		Data: []byte{
			byte(ID), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
	}
}

func MakePing(c NodeClass) socketcan.CanFrame {
	return socketcan.CanFrame{
		ID:       MakeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(HeartBeatRequest)),
		Dlc:      0,
		Extended: false,
		RTR:      false,
		Data:     nil,
	}
}

func MakeHeartBeatRequest(c NodeClass, t time.Duration) socketcan.CanFrame {
	period, err := castDuration(t)
	if err != nil {
		return socketcan.CanFrame{ID: 0x3ff, Dlc: 0}
	}

	f := socketcan.CanFrame{
		ID:       MakeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(HeartBeatRequest)),
		Dlc:      2,
		Extended: false,
		RTR:      false,
		Data:     make([]byte, 2),
	}
	binary.LittleEndian.PutUint16(f.Data, period)
	return f
}

func MakeIDChangeRequest(c NodeClass, original, new NodeID) socketcan.CanFrame {
	return socketcan.CanFrame{
		ID:       MakeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(IDChangeRequest)),
		Dlc:      2,
		Extended: false,
		RTR:      false,
		Data:     []byte{byte(original), byte(new)},
	}
}

type ResetRequestData struct {
	Class NodeClass
	ID    NodeID
}

func (d *ResetRequestData) MessageClassID() MessageClass {
	return ResetRequestMessage
}

func (d *ResetRequestData) String() string {
	if d.ID == 0 {
		return fmt.Sprintf("arke.ResetRequest{Class: %s, Node: All}", d.Class)
	}
	return fmt.Sprintf("arke.ResetRequest{Class: %s, Node: %d}", d.Class, d.ID)
}

func (d *ResetRequestData) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 1); err != nil {
		return err
	}
	d.ID = NodeID(buf[0])
	return nil
}

type HeartBeatRequestData struct {
	Class  NodeClass
	Period time.Duration
}

func (d *HeartBeatRequestData) MessageClassID() MessageClass {
	return HeartBeatRequestMessage
}

func (d *HeartBeatRequestData) String() string {
	periodStr := "SinglePing"
	if d.Period != 0 {
		periodStr = d.Period.String()
	}

	return fmt.Sprintf("arke.HeartBeatRequest{Class: %s, Node: All, Period: %s}", d.Class, periodStr)
}

func (d *HeartBeatRequestData) Unmarshal(buf []byte) error {
	if len(buf) == 0 {
		d.Period = 0
		return nil
	}
	if err := checkSize(buf, 2); err != nil {
		return err
	}
	periodMs := binary.LittleEndian.Uint16(buf)
	d.Period = time.Duration(periodMs) * time.Millisecond
	return nil
}

type IDChangeRequestData struct {
	Class    NodeClass
	Old, New NodeID
}

func (d *IDChangeRequestData) MessageClassID() MessageClass {
	return IDChangeRequestMessage
}

func (d *IDChangeRequestData) String() string {
	return fmt.Sprintf("arke.IDChangeRequest{Class: %s, OldID: %d, NewID: %d}", d.Class, d.Old, d.New)
}

func (d *IDChangeRequestData) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 2); err != nil {
		return err
	}
	d.Old = NodeID(buf[0])
	d.New = NodeID(buf[1])
	return nil
}

type ErrorReportData struct {
	Class     NodeClass
	ID        NodeID
	ErrorCode uint16
}

func (d *ErrorReportData) MessageClassID() MessageClass {
	return ErrorReportMessage
}

func (d *ErrorReportData) String() string {
	return fmt.Sprintf("arke.ErrorReport{Class: %s, ID: %d, ErrorCode: 0x%04x}", d.Class, d.ID, d.ErrorCode)
}

func (d *ErrorReportData) Unmarshal(buf []byte) error {
	if err := checkSize(buf, 4); err != nil {
		return err
	}
	d.Class = NodeClass(buf[0])
	d.ID = NodeID(buf[1])
	d.ErrorCode = binary.LittleEndian.Uint16(buf[2:])
	return nil
}

type HeartBeatData struct {
	Class        NodeClass
	ID           NodeID
	MajorVersion uint8
	MinorVersion uint8
	PatchVersion uint8
	TweakVersion uint8
}

func (h *HeartBeatData) Unmarshal(buf []byte) error {
	h.MajorVersion = 0
	h.MinorVersion = 0
	h.PatchVersion = 0
	h.TweakVersion = 0
	if len(buf) == 0 {
		return nil
	}

	if len(buf) == 1 {
		return fmt.Errorf("Invalid buffer size 1 (min 2 required)")
	}

	h.MajorVersion = buf[0]
	h.MinorVersion = buf[1]
	if len(buf) > 2 {
		h.PatchVersion = buf[2]
	}
	if len(buf) > 3 {
		h.TweakVersion = buf[3]
	}
	return nil
}

func (h *HeartBeatData) String() string {
	if h.MajorVersion == 0 && h.MinorVersion == 0 && h.PatchVersion == 0 && h.TweakVersion == 0 {
		return fmt.Sprintf("arke.HeartBeat{Class: %s, ID: %d}", ClassName(h.Class), h.ID)
	}

	return fmt.Sprintf("arke.HeartBeat{Class: %s, ID: %d, Version: %d.%d.%d.%d}",
		ClassName(h.Class),
		h.ID,
		h.MajorVersion,
		h.MinorVersion,
		h.PatchVersion,
		h.TweakVersion)
}

func (h *HeartBeatData) MessageClassID() MessageClass {
	return HeartBeatMessage
}

func init() {
	networkCommandFactory[NodeID(ResetRequest)] = func(c MessageClass, buffer []byte) (ReceivableMessage, NodeID, error) {
		res := &ResetRequestData{
			Class: NodeClass(c),
		}
		if err := res.Unmarshal(buffer); err != nil {
			return nil, 0, err
		}
		return res, res.ID, nil
	}

	networkCommandFactory[NodeID(IDChangeRequest)] = func(c MessageClass, buffer []byte) (ReceivableMessage, NodeID, error) {
		res := &IDChangeRequestData{
			Class: NodeClass(c),
		}
		if err := res.Unmarshal(buffer); err != nil {
			return nil, 0, err
		}
		return res, res.Old, nil
	}

	networkCommandFactory[NodeID(HeartBeatRequest)] = func(c MessageClass, buffer []byte) (ReceivableMessage, NodeID, error) {
		res := &HeartBeatRequestData{
			Class: NodeClass(c),
		}
		if err := res.Unmarshal(buffer); err != nil {
			return nil, 0, err
		}

		return res, 0, nil
	}

	networkCommandFactory[NodeID(ErrorReport)] = func(c MessageClass, buffer []byte) (ReceivableMessage, NodeID, error) {
		res := &ErrorReportData{}
		if err := res.Unmarshal(buffer); err != nil {
			return nil, 0, err
		}
		return res, res.ID, nil
	}

}
