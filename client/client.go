package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
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

	f, err := os.Open("book.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	buf := make([]byte, 0, 2000)
	for {
		time.Sleep(2 * time.Second)
		randLen := rand.Intn(500) // 1-n
		fmt.Printf("take %d bytes from file\n", randLen)
		n, err := io.ReadFull(r, buf[:randLen])
		buf = buf[:n]
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.ErrUnexpectedEOF {
				log.Println(err)
				break
			}
		}
		_, err = conn.Write(buf)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	// -=Old=-
	//	example := [][]byte{
	//	[]byte("111\n222\n333"),
	//	[]byte("444"),
	//	[]byte("777\n"),
	//	[]byte("\n"),
	//	[]byte("\n666\n777"),
	//	[]byte("111\n222"),
	//	[]byte("111\n222\n"),
	//}
	//
	//for _, bytes := range example {
	//	time.Sleep(2 * time.Second)
	//	fmt.Println("Sent to server =>", string(bytes))
	//	_, err := conn.Write(bytes)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	select {}

}
