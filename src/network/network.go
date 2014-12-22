package network

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
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

func GetRelativeFilePath(relativeFileName string) (filePath string) {
	_, filename, _, _ := runtime.Caller(1)
	filePath = path.Join(path.Dir(filename), relativeFileName)
	return
}

func CreateFormWithInitialValues(formEntries map[string]string) (form url.Values) {
	form = url.Values{}
	for key, value := range formEntries {
		form.Set(key, value)
	}
	return
}

func GetOS() (OS string) {
	operatingSystemBytes, err := exec.Command("uname", "-a").Output() // display operating system name...why do we need the -a?
	operatingSystemName := string(operatingSystemBytes)

	var kernel string

	if ErrorHandler(err, "Error encountered while reading kernel: %v\n") {
		kernel = "Unknown"
	} else {
		kernel = strings.Split(operatingSystemName, " ")[0]
	}
	fmt.Println("Kernel detected: ", kernel)

	switch kernel {
	case "Linux":
		OS = "Linux"
	case "Darwin":
		OS = "OS X"
	default:
		OS = "Unknown"
	}

	if OS == "Unknown" {
		fmt.Println("ERROR: Failed to detect operating system.")
		fmt.Println("Aborting program.")
		os.Exit(1)
	}

	fmt.Printf("Operating system detected: %v\n", OS)
	return OS
}
