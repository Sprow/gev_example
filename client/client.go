package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, e := net.Dial("tcp", ":1833")
	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	go func() {
		buf := bufio.NewReader(conn)
		for {
			time.Sleep(100 * time.Millisecond)
			msg, err := buf.ReadBytes('\n')
			if err != nil {
				panic(err)
			}

			fmt.Printf("Message from server: %s", string(msg))
		}
	}()

	example := [][]byte{
		[]byte("111\n222\n333"),
		[]byte("444"),
		[]byte("777\n"),
		[]byte("\n"),
		[]byte("\n666\n777"),
		[]byte("111\n222"),
		[]byte("111\n222\n"),
	}

	for _, bytes := range example {
		time.Sleep(2 * time.Second)
		fmt.Println("Sent to server =>", string(bytes))
		_, err := conn.Write(bytes)
		if err != nil {
			panic(err)
		}
	}
	select {}

}
