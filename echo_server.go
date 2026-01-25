package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

const defaultPort = 8080

func main() {

	port := flag.Int("port", defaultPort, "port to listen on")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("failed to create socket on port %d: %v", *port, err)
	}

	fmt.Printf("Socket created\n")

	fmt.Printf("Socket bound to port : %d\n", *port)

	fmt.Printf("Listening on port : %d\n", *port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("failed to accept incoming connection: %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		fmt.Printf("Client connected %s: %d\n", addr.IP.String(), addr.Port)
	}

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)

	if err != nil {
		log.Printf("recv failed: %v", err)
		return
	}

	if _, err := conn.Write(buffer[:n]); err != nil {
		log.Printf("send failed: %v", err)
	}
}
