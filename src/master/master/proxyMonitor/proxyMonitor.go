package proxyMonitor

import (
	"fmt"
	"net"
	"net/http"
	"network"
	"time"
)

var (
	IS_USING_PROXY   = false
	PROXY_DEAD_TIME  = 0
	PROXY_PORT       string
	PROXY_IP_ADDRESS string
	PROXY_URL        string
)

func ReceiveProxyHeartbeat(request *http.Request) {
	if !IS_USING_PROXY {
		PROXY_IP_ADDRESS = getOriginIPAddressFromRequest(request)
		PROXY_PORT = request.PostFormValue("ProxyHTTPServerPort")
		PROXY_URL = "http://" + PROXY_IP_ADDRESS + ":" + PROXY_PORT
		fmt.Printf("Proxy detected at %v.\n", PROXY_IP_ADDRESS)
		IS_USING_PROXY = true
	}
	PROXY_DEAD_TIME = 0
}

func getOriginIPAddressFromRequest(request *http.Request) (proxyIP string) {
	proxyIP, _, _ = net.SplitHostPort(request.RemoteAddr)
	return
}

func MonitorProxy(timeInterval int) {
	timer := time.Tick(time.Duration(timeInterval) * time.Second)
	for _ = range timer {
		if IS_USING_PROXY {
			PROXY_DEAD_TIME += 1
			checkProxyHealth(timeInterval)
		}
	}
}

func checkProxyHealth(timeInterval int) {
	if PROXY_DEAD_TIME > timeInterval {
		fmt.Println("Proxy has been disconnected.")
		IS_USING_PROXY = false
	}
}

func RequestProxyToAddNewSlaveToIPTables(proxyURL, IPAddressToAdd string) (err error) {
	proxyURLToUpdateIPTables := proxyURL + "/update_iptables"

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"IPAddressToAdd": IPAddressToAdd})

	_, err = client.PostForm(proxyURLToUpdateIPTables, form)
	network.ErrorHandler(err, "Error communicating with proxy: %v\n")
	return
}

func RequestProxyToRemoveDeadSlaveFromIPTables(proxyURL, IPAddressToDelete string) (err error) {
	proxyURLToUpdateIPTables := proxyURL + "/update_iptables"

	client := &http.Client{}
	form := network.CreateFormWithInitialValues(map[string]string{"IPAddressToDelete": IPAddressToDelete})

	_, err = client.PostForm(proxyURLToUpdateIPTables, form)
	network.ErrorHandler(err, "Error communicating with proxy: %v\n")
	return
}
