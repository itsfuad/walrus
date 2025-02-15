package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"walrus/compiler/analyzer"
	"walrus/compiler/report"
	"walrus/compiler/wio"
)

// Message represents a JSON-RPC message.
type Message struct {
	JsonRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method,omitempty"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

// InitializeResult is the response to the initialize request.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

// ServerCapabilities defines the server's capabilities.
type ServerCapabilities struct {
	TextDocumentSync TextDocumentSyncOptions `json:"textDocumentSync"`
}

// TextDocumentSyncOptions defines how text document changes are handled.
type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose"`
	Change    int  `json:"change"` // 1 = full
}

// Position represents a position in a text document.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range represents a range in a text document.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

func handleMessage(msg *Message) *Message {
	log.Printf("Received message: %s", msg.Method)

	switch msg.Method {
	case "initialize":
		log.Println("Handling initialize request")
		// Respond to the initialize request with server capabilities.
		response := &Message{
			JsonRPC: "2.0",
			ID:      msg.ID,
			Result: InitializeResult{
				Capabilities: ServerCapabilities{
					TextDocumentSync: TextDocumentSyncOptions{
						OpenClose: true,
						Change:    1, // Full text synchronization
					},
				},
			},
		}
		log.Printf("Sending initialize response: %+v", response)
		return response

	case "textDocument/didOpen":
		log.Println("Handling textDocument/didOpen notification")
		var params struct {
			TextDocument struct {
				URI  string `json:"uri"`
				Text string `json:"text"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(msg.Params, &params); err != nil {
			log.Printf("Error parsing didOpen params: %v", err)
			return nil
		}
		processDiagnostics(params.TextDocument.URI)

	case "textDocument/didChange":
		log.Println("Handling textDocument/didChange notification")
		var params struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			ContentChanges []struct {
				Text string `json:"text"`
			} `json:"contentChanges"`
		}
		if err := json.Unmarshal(msg.Params, &params); err != nil {
			log.Printf("Error parsing didChange params: %v", err)
			return nil
		}
		if len(params.ContentChanges) > 0 {
			// Process diagnostics for the full text update.
			processDiagnostics(params.TextDocument.URI)
		}
	}

	return nil
}

// processDiagnostics analyzes the file and publishes diagnostics.
func processDiagnostics(uri string) {
	filePath, err := wio.UriToFilePath(uri)
	if err != nil {
		log.Printf("Error converting URI to file path: %v", err)
		return
	}
	log.Printf("Processing diagnostics for: %s", filePath)
	reports, err := analyzer.Analyze(filePath, false, false, false)
	if err != nil {
		log.Printf("Error analyzing file: %v", err)
		return
	}
	publishDiagnostics(filePath, reports)
}

// publishDiagnostics sends diagnostics to the client.
func publishDiagnostics(uri string, diagnostics report.IReport) {
	params := struct {
		URI         string              `json:"uri"`
		Diagnostics report.IReport 		`json:"diagnostics"`
	}{
		URI:         uri,
		Diagnostics: diagnostics,
	}
	data, err := json.Marshal(params)
	if err != nil {
		log.Printf("Error marshalling diagnostics: %v", err)
		return
	}
	notification := Message{
		JsonRPC: "2.0",
		Method:  "textDocument/publishDiagnostics",
		Params:  data,
	}
	if err := writeMessage(os.Stdout, &notification); err != nil {
		log.Printf("Error writing diagnostics: %v", err)
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
			log.Printf("Ignoring header: %s", line)
			continue
		}
	}

	log.Printf("Content-Length: %d", contentLength)

	// Read content
	content := make([]byte, contentLength)
	if _, err := io.ReadFull(r, content); err != nil {
		log.Printf("Error reading content: %v", err)
		return err
	}

	log.Printf("Received content: %s", content)

	return json.Unmarshal(content, msg)
}

// writeMessage writes a JSON-RPC message to the output stream.
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

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		var msg Message
		if err := readMessage(reader, &msg); err != nil {
			if err == io.EOF {
				log.Println("EOF received, exiting")
				return
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		if response := handleMessage(&msg); response != nil {
			if err := writeMessage(os.Stdout, response); err != nil {
				log.Printf("Error writing response: %v", err)
			}
		} else {
			log.Printf("No response for message: %+v", msg)
		}
	}
}