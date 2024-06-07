package main

import (
	"bufio"
	"fmt"
	"net"
)

type LoadBalancer struct {
	Port    string
	servers []Server
	algo    RoutingAlgorithm
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Port: Config.Port,
		servers: []Server{
			{url: "http://localhost:3001", active: true},
			{url: "http://localhost:3002", active: true},
		},
		algo: &RoundRobin{current_index: -1},
	}
}

func (lb *LoadBalancer) getServer() Server {
	return lb.algo.getNextServer(lb.servers)

}

func (lb *LoadBalancer) ListenAndServe() {

	ln, err := net.Listen("tcp", lb.Port)
	if err != nil {
		fmt.Println("error Listen", err)
	}
	fmt.Println("Starting Load Balancer on PORT", lb.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error Accept", err)
		}
		go handleConn(conn)
	}

}

func handleConn(src net.Conn) {
	// buffer := make([]byte, 1024)
	// _, err := src.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		t := scanner.Text()
		fmt.Println(t)
		if t == "" {
			// Empty line encountered, indicating end of headers
			break
		}
	}

	fmt.Println("---------------")
	src.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nDate: Sat, 01 May 2024 12:00:00 GMT\r\nServer: MyServer/1.0\r\n\r\ndataaaaa\n"))

	// src.Write([]byte("dataaaaa"))
	src.Close()
}
