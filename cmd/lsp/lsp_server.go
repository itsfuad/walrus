package main

import (
	"encoding/json"
	"log"
	"net"
)

type JsonRPC2 struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Server struct {
	conn  net.Conn
	files map[string]string
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("LSP server listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection failed: %v", err)
			continue
		}
		server := &Server{
			conn:  conn,
			files: make(map[string]string),
		}
		go server.handleConnection()
	}
}

func (s *Server) handleConnection() {
	defer s.conn.Close()

	decoder := json.NewDecoder(s.conn)
	encoder := json.NewEncoder(s.conn)

	for {
		var msg JsonRPC2
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Failed to decode message: %v", err)
			return
		}

		switch msg.Method {
		case "initialize":
			response := JsonRPC2{
				JsonRPC: "2.0",
				ID:      msg.ID,
				Result: map[string]interface{}{
					"capabilities": map[string]interface{}{
						"textDocumentSync": 1,
						"completionProvider": map[string]interface{}{
							"triggerCharacters": []string{".", ":"},
						},
					},
				},
			}
			encoder.Encode(response)

		case "textDocument/didOpen":
			// Handle file open
			s.handleNotification(msg)

		case "textDocument/didChange":
			// Handle file change
			s.handleNotification(msg)

		case "textDocument/didSave":
			// Handle file save
			s.handleNotification(msg)

		case "shutdown":
			response := JsonRPC2{
				JsonRPC: "2.0",
				ID:      msg.ID,
				Result:  nil,
			}
			encoder.Encode(response)
			return
		}
	}
}

func (s *Server) handleNotification(msg JsonRPC2) {
	// Here you can add your language-specific logic
	// For example, triggering the lexer, parser, etc.
	log.Printf("Received notification: %s", msg.Method)
}
