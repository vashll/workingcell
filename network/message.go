package network

import (
	"github.com/golang/protobuf/proto"
)

type Message struct {
	Cmd  uint32
	Data []byte
}

func NewDataMsg(data []byte) *Message {
	return &Message{
		Data: data,
	}
}

func ParseProtoMsgData(data []byte) {
	proto.Unmarshal(data, nil)
}
