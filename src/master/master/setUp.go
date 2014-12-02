package master

import "time"

type Slave struct {
	URL          string
	Heartbeat    time.Time
	DisplayedURL string // TODO: store currently displayed URL for each slave
}

func SetUp() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	return
}
