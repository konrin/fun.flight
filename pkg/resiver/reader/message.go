package reader

import "time"

type Message struct {
	Data      string
	MlatAt    int64
	RSSI      int
	ReceiveAt time.Time
}
