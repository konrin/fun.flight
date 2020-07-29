package reader

type Reader interface {
	StartReading() error
	GetMessageCh() <-chan *Message
	GetErrorCh() <-chan error
}
