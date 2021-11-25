package message

type ClientMessage struct {
	Type    MessageType
	Content []byte
}

func NewErrorClientMessage(err error) ClientMessage {
	return ClientMessage{
		Type:    ErrorMessageType,
		Content: []byte(err.Error()),
	}
}
