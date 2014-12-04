package master

import (
	"path"
	"runtime"
	"time"
)

type Slave struct {
	URL          string
	Heartbeat    time.Time
	DisplayedURL string // TODO: store currently displayed URL for each slave
}

func SetUp() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	return
}

func GetRelativeFilePath(relativeFileName string) (filePath string) {
	_, filename, _, _ := runtime.Caller(1)
	filePath = path.Join(path.Dir(filename), relativeFileName)
	return
}
