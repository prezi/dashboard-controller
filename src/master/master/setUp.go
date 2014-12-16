package master

import (
	"time"
)

type Slave struct {
	URL                    string
	Heartbeat              time.Time
	PreviouslyDisplayedURL string
	DisplayedURL           string
}
