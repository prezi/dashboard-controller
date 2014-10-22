package main

import(
	"fmt"
	"log"
	"net"
)

const listenAddress = "localhost:4000"

func ListenOnPortAndReply() {
	l, err := net.Listen("tcp", listenAddress)

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
