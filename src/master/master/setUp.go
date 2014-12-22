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

const (
	DEFAULT_PROXY_IP_ADDRESS = "localhost"
	DEFAULT_PROXY_PORT       = "8080"
)

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

func GetProxyURL() (proxyURL string) {
	proxyIP, proxyPort := configFlags()
	proxyURL = "http://" + proxyIP + ":" + proxyPort
	fmt.Printf("Proxy registered at %v", proxyURL)
	return proxyURL
}

func configFlags() (proxyIP, proxyPort string) {
	flag.StringVar(&proxyIP, "proxyIP", DEFAULT_PROXY_IP_ADDRESS, "proxy address")
	flag.StringVar(&proxyPort, "proxyPort", DEFAULT_PROXY_PORT, "proxy port")
	flag.Parse()
	return
}
