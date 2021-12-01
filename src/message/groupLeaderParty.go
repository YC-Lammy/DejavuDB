package message

import "unsafe"

type GroupLeaderPartyMessage interface {
	ToBytes() []byte
	isGroupLeaderPartyMessage()
}

type GroupLeaderPartyMessageType byte

const (
	PleaseGiveMeThatValue_Type GroupLeaderPartyMessageType = iota
)

func UnMarshalGroupLeaderPartyMessage(b []byte) GroupLeaderPartyMessage {
	switch GroupLeaderPartyMessageType(b[0]) {
	case PleaseGiveMeThatValue_Type:
		p := PleaseGiveMeThatValue_Message{}
		p.FromBytes(b[1:])
		return p
	}
	return nil
}

type PleaseGiveMeThatValue_Message struct {
	SenderId uint32
	Key      string
}

func (PleaseGiveMeThatValue_Message) isGroupLeaderPartyMessage() {}

func (p *PleaseGiveMeThatValue_Message) FromBytes(b []byte) {
	id := [4]byte{b[0], b[1], b[2], b[3]}
	p.SenderId = *(*uint32)(unsafe.Pointer(&id))
	p.Key = string(b[4:])
}

func (p PleaseGiveMeThatValue_Message) ToBytes() []byte {
	re := []byte{}
	id := *(*[4]byte)(unsafe.Pointer(&p.SenderId))
	re = append(re, id[0], id[1], id[2], id[3])
	re = append(re, p.Key...)
	return re
}

type IamGivingYouThatValue_Message struct {
	SenderId  uint32
	Timestamp uint64
	Key       string
	Value     string
}

func (IamGivingYouThatValue_Message) isGroupLeaderPartyMessage() {}
