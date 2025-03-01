package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	CONN_PORT      = ":3334"
	CONN_TYPE      = "tcp"
	MSG_DISCONNECT = "Disconnected from the server.\n"
)

var wg sync.WaitGroup

func Read(conn net.Conn) {
	defer wg.Done()
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(MSG_DISCONNECT)
			return
		}
		fmt.Print(str)
	}
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		fmt.Print("> ")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		str = strings.TrimSpace(str)
		if str != "" {
			notifyTyping(conn)
		}

		_, err = writer.WriteString(str + "\n")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func notifyTyping(conn net.Conn) {
	conn.Write([]byte("User is typing...\n"))
}

func main() {
	conn, err := net.Dial(CONN_TYPE, CONN_PORT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	wg.Add(1)
	go Read(conn)
	go Write(conn)

	wg.Wait()
}
