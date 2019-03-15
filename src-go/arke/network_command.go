package arke

import (
	"encoding/binary"
	"fmt"
	"time"

	socketcan "github.com/atuleu/golang-socketcan"
)

func SendResetRequest(itf socketcan.RawInterface, c NodeClass, ID NodeID) error {
	return itf.Send(socketcan.CanFrame{
		ID:       makeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(ResetRequest)),
		Dlc:      1,
		Extended: false,
		RTR:      false,
		Data: []byte{
			byte(ID), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
	})
}

func Ping(itf socketcan.RawInterface, c NodeClass) error {
	return itf.Send(socketcan.CanFrame{
		ID:       makeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(HeartBeatRequest)),
		Dlc:      0,
		Extended: false,
		RTR:      false,
		Data:     nil,
	})
}

func SendHeartBeatRequest(itf socketcan.RawInterface, c NodeClass, t time.Duration) error {
	period, err := castDuration(t)
	if err != nil {
		return err
	}

	f := socketcan.CanFrame{
		ID:       makeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(HeartBeatRequest)),
		Dlc:      2,
		Extended: false,
		RTR:      false,
		Data:     make([]byte, 2),
	}
	binary.LittleEndian.PutUint16(f.Data, period)
	return itf.Send(f)
}

func SendIDChangeRequest(itf socketcan.RawInterface, c NodeClass, original, new NodeID) error {
	if original == 0 || original > 7 {
		return fmt.Errorf("Invalid Original ID %d", original)
	}
	if new == 0 || new > 7 {
		return fmt.Errorf("Invalid Target ID %d", new)
	}

	return itf.Send(socketcan.CanFrame{
		ID:       makeCANIDT(NetworkControlCommand, MessageClass(c), NodeID(IDChangeRequest)),
		Dlc:      2,
		Extended: false,
		RTR:      false,
		Data:     []byte{byte(original), byte(new)},
	})
}

type HeartBeatData struct {
	Class        NodeClass
	ID           NodeID
	MajorVersion uint8
	MinorVersion uint8
	PatchVersion uint8
	TweakVersion uint8
}

func (h *HeartBeatData) Unmarshall(buf []byte) error {
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
