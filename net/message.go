package net

import (
	"github.com/golang/protobuf/proto"
)

type Message struct {
	Cmd uint32
	Data []byte
}

func NewDataMsg(data []byte) *Message {
	return &Message{
		Data: data,
	}
}

func ParseProroMsgData(data []byte) {

}