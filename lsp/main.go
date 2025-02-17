package main

import "log"

func main() {
	conn, err := createServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
		return
	}
	handleConnection(conn)
}
