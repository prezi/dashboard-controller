package hash

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func InitializeUserAuthenticationMap(filePathToUserAuthenticationTxt string) (userAuthenticationMap map[string][16]byte) {
	content, err := ioutil.ReadFile(filePathToUserAuthenticationTxt)
	if err != nil {
		fmt.Println("User authentication data not found.\n", err)
		os.Exit(1)
	}
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

func CreateHashFromString(inputString string) (outputHash [16]byte) {
	data := []byte(inputString)
	outputHash = md5.Sum(data)
	return
}

func IsHashMatchInUserAuthenticationMap(inputUsername, inputPassword string, userAuthenticationMap map[string][16]byte) (hashMatch bool) {
	targetHash := CreateHashFromString(inputPassword)
	return reflect.DeepEqual(userAuthenticationMap[inputUsername], targetHash)
}
