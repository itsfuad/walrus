package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"walrus/compiler/parser"
	"walrus/compiler/report"
	"walrus/compiler/typechecker"
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

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// handleMessage processes requests and notifications.
func handleMessage(msg *Message) *Message {
	// Removed output to stdout to avoid protocol corruption.
	// log to stderr if needed:
	log.Printf("Received message: %s", msg.Method)

	switch msg.Method {
	case "initialize":
		return &Message{
			JsonRPC: "2.0",
			ID:      msg.ID,
			Result: InitializeResult{
				Capabilities: ServerCapabilities{
					TextDocumentSync: TextDocumentSyncOptions{
						OpenClose: true,
						Change:    1, // using full text on change
					},
				},
			},
		}
	case "textDocument/didOpen":
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
			// Use full text update.
			processDiagnostics(params.TextDocument.URI)
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

// processDiagnostics now uses the compiler's lexer, parser, and typechecker.
func processDiagnostics(filePath string) {

	tree := parser.NewParser(filePath, false).Parse(false)
	env := typechecker.ProgramEnv(filePath)
	typechecker.CheckAST(tree, env)

	// Fetch diagnostics produced during typechecking.
	diagnostics := report.GetDiagnostics()
	publishDiagnostics(filePath, diagnostics)
}

func publishDiagnostics(uri string, diagnostics []report.Diagnostic) {
	params := struct {
		URI         string              `json:"uri"`
		Diagnostics []report.Diagnostic `json:"diagnostics"`
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
