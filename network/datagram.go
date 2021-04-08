package network

import (
	"io"
	"net"
	"unsafe"
	. "workincell/log"
)

const (
	DataHeadSize = 8
)

type IDataReader interface {
	ReadData(conn net.Conn) (error, *Message)
	MsgToData(msg *Message) []byte
}

type DataHead struct {
	Len uint32
	Cmd uint32
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
type DefaultDataReader struct {
}

func (r *DefaultDataReader) ReadData(conn net.Conn) (err error, msg *Message) {
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
	msg = &Message{
		Cmd: head.Cmd,
	}
	if head.Len > 0 {
		dataBuf := make([]byte, head.Len)
		_, err = io.ReadFull(conn, dataBuf[DataHeadSize:])
		if err != nil {
			LogInfo("data handler read data error:%s", err)
			return
		}
		msg.Data = dataBuf
	}
	return
}

func (r *DefaultDataReader) MsgToData(msg *Message) []byte {
	size := len(msg.Data) + DataHeadSize
	buf := make([]byte, size)
	headBuf := buf[:DataHeadSize]
	header := (*DataHead)(unsafe.Pointer(&headBuf[0]))
	header.Cmd = msg.Cmd
	header.Len = uint32(size)
	if msg.Data != nil {
		copy(buf[:DataHeadSize], msg.Data)
	}
	return buf
}

//func DefaultDataToMsg(data []byte) *Message {
//	msg := &Message{}
//	headBuf := data[:DataHeadSize]
//	head := &DataHead{}
//	head.Unmarshal(headBuf)
//	msg.Cmd = head.Cmd
//	if head.Len > 0 {
//		msg.Data = data[DataHeadSize:]
//	}
//	return msg
//}
