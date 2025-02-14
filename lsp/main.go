package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Message struct {
	JsonRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method,omitempty"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	TextDocumentSync TextDocumentSyncOptions `json:"textDocumentSync"`
}

type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose"`
	Change    int  `json:"change"` // 1 = full
}

func handleMessage(msg *Message) *Message {
	if msg.Method == "initialize" {
		return &Message{
			JsonRPC: "2.0",
			ID:      msg.ID,
			Result: InitializeResult{
				Capabilities: ServerCapabilities{
					TextDocumentSync: TextDocumentSyncOptions{
						OpenClose: true,
						Change:    1,
					},
				},
			},
		}
	}
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var msg Message
		if err := readMessage(reader, &msg); err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		if response := handleMessage(&msg); response != nil {
			if err := writeMessage(os.Stdout, response); err != nil {
				log.Printf("Error writing response: %v", err)
			}
		}
	}
}

func readMessage(r *bufio.Reader, msg *Message) error {
	// Read headers
	contentLength := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		if line == "\r\n" {
			break
		}
		if _, err := fmt.Sscanf(line, "Content-Length: %d\r\n", &contentLength); err != nil {
			continue
		}
	}

	// Read content
	content := make([]byte, contentLength)
	if _, err := io.ReadFull(r, content); err != nil {
		return err
	}

	return json.Unmarshal(content, msg)
}

func writeMessage(w io.Writer, msg *Message) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	if _, err := w.Write([]byte(header)); err != nil {
		return err
	}

	_, err = w.Write(content)
	return err
}
