package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	fmt.Println("Starting Load Balancer on PORT", Config.Port)

	ln, err := net.Listen("tcp", Config.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}

}

func handleConn(src net.Conn) {
	defer src.Close()
	fmt.Printf("Received request from %s", src.RemoteAddr())
	dst, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()
	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatal(err)
	}
}
