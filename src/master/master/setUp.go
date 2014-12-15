package master

import (
	"path"
	"runtime"
	"time"
	"sort"
)

type Slave struct {
	URL                    string
	Heartbeat              time.Time
	PreviouslyDisplayedURL string
	DisplayedURL           string
}

func GetSlaveMap() (slaveMap map[string]Slave) {
	slaveMap = make(map[string]Slave)
	return
}

func GetSlaveNamesFromMap(slaveMap map[string]Slave) (slaveNames []string) {
	for index := range slaveMap {
		slaveNames = append(slaveNames, index)
	}
	sort.Strings(slaveNames)
	return
}

func GetRelativeFilePath(relativeFileName string) (filePath string) {
	_, filename, _, _ := runtime.Caller(1)
	filePath = path.Join(path.Dir(filename), relativeFileName)
	return
}
