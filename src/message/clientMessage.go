package message

type ClientMessage struct {
	Type    MessageType
	Content []byte
}

func (c ClientMessage) ToBytes() []byte {
	return append([]byte{byte(c.Type)}, c.Content...)
}

func ClientMessageFromBytes(b []byte) ClientMessage {
	return ClientMessage{
		Type:    MessageType(b[0]),
		Content: b[1:],
	}
}

func NewErrorClientMessage(err error) ClientMessage {
	return ClientMessage{
		Type:    ErrorMessageType,
		Content: []byte(err.Error()),
	}
}
