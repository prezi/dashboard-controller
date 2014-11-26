package master

import (
	"fmt"
	"net/http"
)

func MonitorWebserverRegistration(_ http.ResponseWriter, request *http.Request) {
	webserverUrl := request.PostFormValue("webserverUrl")
	fmt.Println("############## WebserverURL :", webserverUrl)
	//	webserverAddress = webserverUrl
}
