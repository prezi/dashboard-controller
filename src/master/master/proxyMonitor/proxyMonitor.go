package proxyMonitor

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

var IS_USING_PROXY = false
var PROXY_DEAD_TIME = 0

func ReceiveProxyHeartbeat(request *http.Request) {
	if !IS_USING_PROXY {
		fmt.Printf("Proxy detected at %v.\n", getOriginIPAddressFromRequest(request))
		IS_USING_PROXY = true
	}
}

func getOriginIPAddressFromRequest(request *http.Request) (proxyIP string) {
	proxyIP, _, _ = net.SplitHostPort(request.RemoteAddr)
	return
}

func MonitorProxy(timeInterval int) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		if IS_USING_PROXY {
			checkProxyHealth(timeInterval)
		}
	}
}

func checkProxyHealth(timeInterval int) {
	PROXY_DEAD_TIME += 1
	if PROXY_DEAD_TIME > timeInterval {
		fmt.Println("Proxy has been disconnected.")
		IS_USING_PROXY = false
	}
}
