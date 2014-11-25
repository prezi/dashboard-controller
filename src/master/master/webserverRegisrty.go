package master

import (
	"net/http"
	"fmt"
)

func MonitorWebserverRegistration(_ http.ResponseWriter, request *http.Request) {
	webserverUrl := request.PostFormValue("webserverUrl")
	fmt.Println("############## WebserverURL :", webserverUrl)
//	webserverAddress = webserverUrl
}
