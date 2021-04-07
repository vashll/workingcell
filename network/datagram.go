package net

import (
	"errors"
	"io"
	"net"
	"unsafe"
	. "workincell/log"
)

const (
	DataHeadSize = 8
)

type IDataHandler interface {
	ReadData(conn net.Conn) (error, []byte)
	SendData(conn net.Conn, data []byte) error

}

type DataHead struct {
	Len   uint32
	Cmd   uint32
}

func (r *DataHead) Marshal() []byte {
	buf := make([]byte, DataHeadSize)
	head := (*DataHead)(unsafe.Pointer(&buf[0]))
	head.Len = r.Len
	head.Cmd = r.Cmd
	return buf
}

func (r *DataHead) Unmarshal(data []byte) {
	head := (*DataHead)(unsafe.Pointer(&data[0]))
	r.Len = head.Len
	r.Cmd = head.Cmd
}

//默认数据处理程序
type DefaultDataHandler struct {
}

func (r *DefaultDataHandler) ReadData(conn net.Conn) (err error, dataBuf []byte) {
	headBuf := make([]byte, DataHeadSize)
	var head *DataHead
	size, err := io.ReadFull(conn, headBuf)
	if err != nil {
		//if err != io.EOF {
		//	LogError("connection receive data error:%s", err.Error())
		//}
		return
	}
	if size != DataHeadSize {
		LogError("read data head fail data len:%v", size)
		return
	}
	head = (*DataHead)(unsafe.Pointer(&headBuf[0]))
	dataBuf = make([]byte, head.Len + DataHeadSize)
	copy(dataBuf[:DataHeadSize], headBuf)
	if head.Len > 0 {
		_, err = io.ReadFull(conn, dataBuf[DataHeadSize:])
		if err != nil {
			LogInfo("data handler read data error:%s", err)
			return
		}
	}
	return
}

func (r *DefaultDataHandler) SendData(conn net.Conn, data []byte) error {
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	if n < len(data) {
		return errors.New("size of write data is less than all")
	}
	return nil
}

func DefaultDataToMsg(data []byte) *Message {
	msg := &Message{}
	headBuf := data[:DataHeadSize]
	head := &DataHead{}
	head.Unmarshal(headBuf)
	msg.Cmd = head.Cmd
	if head.Len > 0 {
		msg.Data = data[DataHeadSize:]
	}
	return msg
}