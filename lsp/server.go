package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func createServer() (net.Conn, error) {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds) // Add microseconds to log timestamps

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return nil, err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("PORT:%d\n", port)
	os.Stdout.Sync() // Force flush the port number

	log.Printf("LSP Server listening on port %d", port)

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Failed to accept connection: %v", err)
		return nil, err
	}
	defer conn.Close()

	log.Printf("Client connected from: %s", conn.RemoteAddr())

	return conn, nil
}
