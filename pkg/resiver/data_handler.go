package resiver

import (
	"fmt"
	"net"
	"time"

	"fun.flight/pkg/resiver/reader"
)

const (
	MessageChQueueLen = 10
)

type DataHandler struct {
	resiver *Resiver

	isRunning bool

	conn net.Conn

	reader *reader.Reader

	messageCh chan *reader.Message
	errCh     chan error
	stopCh    chan bool

	// startAt time.Time
	// stopAt  time.Time
}

func NewDataHandler(r *Resiver) *DataHandler {
	return &DataHandler{
		resiver:   r,
		messageCh: make(chan *reader.Message, MessageChQueueLen),
		errCh:     make(chan error, 1),
	}
}

func (r *DataHandler) IsRunning() bool {
	return r.isRunning
}

func (r *DataHandler) GetMessageCh() <-chan *reader.Message {
	return r.messageCh
}

func (r *DataHandler) GetErrorCh() <-chan error {
	return r.errCh
}

func (r *DataHandler) Start() error {
	go func() {
		for {
			err := r.run()
			if err != nil {
				time.Sleep(5 * time.Second)

				continue
			}

			return
		}
	}()

	return nil
}

func (r *DataHandler) run() error {
	r.stopCh = make(chan bool)

	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", r.resiver.cfg.Host, r.resiver.cfg.Port),
	)
	if err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	r.isRunning = true
	r.resiver.state = ResiverStateConnecting

	err = r.startReading(conn)
	if err != nil {
		return err
	}

	r.isRunning = false
	r.resiver.state = ResiverStateStopped

	return nil
}

func (r *DataHandler) Stop() error {
	if !r.isRunning {
		return nil
	}

	r.stopCh <- true

	if err := r.conn.Close(); err != nil {
		return err
	}

	r.reader = nil

	return nil
}

func (r *DataHandler) startReading(conn net.Conn) error {
	var _reader reader.Reader

	switch r.resiver.Type() {
	case ResiverTypeRaw:
		_reader = reader.NewRawReader(conn)
		break
		// case ResiverTypeBeast:
		// 	return
	}

	if err := _reader.StartReading(); err != nil {
		return err
	}

	r.resiver.state = ResiverStateRunning

	for {
		select {
		case <-r.stopCh:
			return nil
		case message := <-_reader.GetMessageCh():
			r.messageCh <- message
			break
		case err := <-_reader.GetErrorCh():
			r.errCh <- err
			break
		}
	}
}
