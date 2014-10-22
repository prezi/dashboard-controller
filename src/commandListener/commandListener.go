package main

import(
    "bufio"
    "fmt"
    "os"
	"log"
	"net"
)

const listenAddress = "localhost:4000"

func main() {
	ListenOnPortAndReply()
}

func ReadMessage() string {
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter link: ")
    link, _ := reader.ReadString('\n')

    return link
}

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
