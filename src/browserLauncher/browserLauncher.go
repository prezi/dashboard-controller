// open the browser with http://www.google.com, hardcoded
// 
package main

import (
    "os/exec"
)

func open(url string) *exec.Cmd {
	return exec.Command("open", url)
}

func main() {
	url := "http://www.google.com"
	open(url).Start()
}
