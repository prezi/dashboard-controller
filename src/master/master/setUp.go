package master

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"time"
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

func checkStatusCode(urlToDisplay string) int {
	if len(urlToDisplay) <= 6 {
		urlToDisplay = "http://" + urlToDisplay
	} else if string(urlToDisplay[0:6]) != "http:/" && string(urlToDisplay[0:6]) != "https:" {
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
	if 400 <= checkStatusCode(url) || checkStatusCode(url) == 0 {
		return false
	}
	return true
}

func GetProxyPort() (proxyPort string) {
	proxyPort = configFlags()
	fmt.Printf("Registered proxy on port: %v\n", proxyPort)
	return
}

func configFlags() (proxyPort string) {
	flag.StringVar(&proxyPort, "proxyPort", "8080", "proxy port")
	flag.Parse()
	return
}
