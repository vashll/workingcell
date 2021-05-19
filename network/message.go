package network

type Message struct {
	Cmd  uint32
	Data []byte
}

func NewDataMsg(data []byte) *Message {
	return &Message{
		Data: data,
	}
}

type CellMessage struct {
	Msg     *Message
	MsgCell *TcpMsgCell
	User    interface{} //user data bind to message cell.
}
