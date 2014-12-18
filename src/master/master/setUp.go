package master

import (
	"path"
	"runtime"
	"time"
	"sort"
	"net/http"
	"fmt"
	"flag"
)

const (
	DEFAULT_PROXY_IP_ADDRESS = "localhost"
	DEFAULT_PROXY_PORT       = "8080"
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

func checkStatusCode(urlToDisplay string) int {
	if (len(urlToDisplay) <= 6) {
		urlToDisplay = "http://" + urlToDisplay
	} else if (string(urlToDisplay[0:6]) != "http:/" && string(urlToDisplay[0:6]) != "https:") {
		urlToDisplay = "http://" + urlToDisplay
	}

	response, err := http.Head(urlToDisplay)
	if err != nil {
		return 0
	} else {
		return response.StatusCode
	}
}

func IsURLValid(url string) bool {
	if 400 <= checkStatusCode(url) || checkStatusCode(url) == 0 { return false }
	return true
}

func SetUpProxy() (proxyAddress, proxyPort string) {
	proxyAddress, proxyPort = configFlags()

	fmt.Printf("Registered proxy at %v", proxyAddress)
	fmt.Printf(" on port: %v\n", proxyPort)
	return
}

func configFlags() (proxyIP, proxyPort string) {
	flag.StringVar(&proxyIP, "proxyIP", DEFAULT_PROXY_IP_ADDRESS, "proxy IP")
	flag.StringVar(&proxyPort, "proxyPort", DEFAULT_PROXY_PORT, "proxy port")
	flag.Parse()

	return
}
