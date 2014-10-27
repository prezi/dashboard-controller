package main

import(
	"fmt"
	"log"
	"net"
)

const (
	DEFAULT_ADDRESS = "localhost:4000"
	DEFAULT_LOG = "/log/commandListener.log"
)

func ListenOnPortAndReply() {
	l, err := net.Listen("tcp", DEFAULT_ADDRESS)

	CheckError(err)

	for {
		c, err :=l.Accept()

		CheckError(err)

		fmt.Fprintln(c, "Hello!")
		c.Close()
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
