package network

import (
	"fmt"
	"net/url"
	"path"
	"runtime"
)

const (
	DEFAULT_SLAVE_NAME        = "SLAVE NAME UNSPECIFIED"
	DEFAULT_MASTER_IP_ADDRESS = "localhost"
	DEFAULT_MASTER_PORT       = "5000"
)

func ErrorHandler(err error, message string) (errorOccurred bool) {
	if err != nil {
		fmt.Printf(message, err)
		return true
	}
	return false
}

func CreateFormWithInitialValues(formEntries map[string]string) (form url.Values) {
	form = url.Values{}
	for key, value := range formEntries {
		form.Set(key, value)
	}
	return
}

func GetRelativeFilePath(relativeFileName string) (filePath string) {
	_, filename, _, _ := runtime.Caller(1)
	filePath = path.Join(path.Dir(filename), relativeFileName)
	return
}
