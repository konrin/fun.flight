package reader

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	RawMLATLen     = 12
	RawShortMsgLen = 14
	RawLongMsgLen  = 28

	RawCharTypeManual   = '*'
	RawCharTypeWithMlat = '@'
	RawEndChar          = ';'

	DataSplitChar = '\n'
)

type RawReader struct {
	conn net.Conn

	isRunning bool

	messageCh chan *Message
	errorCh   chan error
}

func NewRawReader(conn net.Conn) *RawReader {
	return &RawReader{conn: conn}
}

func (rr *RawReader) StartReading() error {
	go func(rr *RawReader) {
		for {
			data, err := bufio.NewReader(rr.conn).ReadString(DataSplitChar)
			if err != nil {
				rr.isRunning = false

				if err == io.EOF {
					rr.errorCh <- NewConnectionError(err, "Resiver closed the connection")
				} else {
					rr.errorCh <- NewConnectionError(err, "error")
				}

				return
			}

			receiveAt := time.Now()

			message, err := rr.readMessage(data)
			if err != nil {
				rr.errorCh <- err

				continue
			}

			message.ReceiveAt = receiveAt

			rr.messageCh <- message
		}
	}(rr)

	rr.isRunning = true

	return nil
}

func (rr *RawReader) GetMessageCh() <-chan *Message {
	return rr.messageCh
}

func (rr *RawReader) GetErrorCh() <-chan error {
	return rr.errorCh
}

func (rr *RawReader) readMessage(data string) (*Message, error) {
	if len(data) < 10 {
		return nil, NewMessageFormatError(data, "Unknown message length")
	}

	if data[len(data)-1] != RawEndChar {
		return nil, NewMessageFormatError(data, "No trailing character found")
	}

	switch data[0] {
	case RawCharTypeManual:
		return rr.ReadManual(data)
	case RawCharTypeWithMlat:
		return rr.ReadWithMlat(data)
	default:
		return nil, NewMessageFormatError(data, "Start character not found")
	}
}

func (rr *RawReader) ReadWithMlat(data string) (*Message, error) {
	data = data[1 : len(data)-1]

	mlatRaw := data[0:RawMLATLen]
	msg := data[RawMLATLen:]

	mlat64, err := strconv.ParseInt(mlatRaw, 16, 64)
	if err != nil {
		return nil, err
	}

	msg = strings.ToUpper(msg)

	return &Message{
		Data:   msg,
		MlatAt: mlat64,
	}, nil
}

func (rr *RawReader) ReadManual(data string) (*Message, error) {
	data = data[1 : len(data)-1]
	data = strings.ToUpper(data)

	return &Message{
		Data: data,
	}, nil
}

type MessageFormatError struct {
	Data    string
	Message string
}

func NewMessageFormatError(data, message string) *MessageFormatError {
	return &MessageFormatError{data, message}
}

func (e *MessageFormatError) Error() string {
	return fmt.Sprintf("Message format error, %s: %s", e.Data, e.Message)
}
