package hash

import (
	"crypto/md5"
	// "fmt"
	"io/ioutil"
	"master/master"
	"network"
	"reflect"
	"strings"
)

var FILE_PATH_TO_USER_AUTHENTICATION_DATA = master.GetRelativeFilePath("./user_authentication_data.txt")

func CreateHashFromString(inputString string) (outputHash [16]byte) {
	data := []byte(inputString)
	outputHash = md5.Sum(data)
	return outputHash
}

func IsHashMatchInUserAuthenticationMap(inputUsername, inputPassword string, userAuthenticationMap map[string][16]byte) (hashMatch bool) {
	targetHash := CreateHashFromString(inputPassword)
	return reflect.DeepEqual(userAuthenticationMap[inputUsername], targetHash)
}

func InitializeUserAuthenticationMap() (userAuthenticationMap map[string][16]byte) {
	content, err := ioutil.ReadFile(FILE_PATH_TO_USER_AUTHENTICATION_DATA)
	network.ErrorHandler(err, "Error encountered while parsing user authentication data: %v")

	username, password := GetUserNameAndPasswordFromFile(string(content))

	userAuthenticationMap = make(map[string][16]byte)
	userAuthenticationMap[username] = CreateHashFromString(password)
	return
}

func GetUserNameAndPasswordFromFile(fileContent string) (username, password string) {
	lines := strings.Split(string(fileContent), "\n")
	username = lines[0]
	password = lines[1]
	return
}
