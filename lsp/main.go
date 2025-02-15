package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"walrus/compiler/analyzer"
	"walrus/compiler/report"
	"walrus/compiler/wio"
)

func init() {
	f, err := os.OpenFile("lsp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *LspError   `json:"error,omitempty"`
}

type LspError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var stdout = bufio.NewWriter(os.Stdout)

func main() {
	log.Println("Starting Walrus LSP...")
	for {
		msg, err := readMessage()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(msg), &req); err != nil {
			log.Printf("Invalid JSON message: %v", err)
			continue
		}

		log.Printf("Received request: %+v", req)
		handleRequest(req)
	}
}

func handleRequest(req Request) {
	switch req.Method {
	case "initialize":
		writeMessage(Response{
			Jsonrpc: "2.0",
			Id:      req.Id,
			Result: map[string]interface{}{
				"capabilities": map[string]interface{}{ "textDocumentSync": 1 },
			},
		})
		// Send initialized notification
		notification := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "initialized",
			"params":  map[string]interface{}{},
		}

		writeRawMessage(notification)

	case "textDocument/didOpen", "textDocument/didChange":
		var params struct {
			TextDocument struct { URI string `json:"uri"` } `json:"textDocument"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			log.Println("Error parsing textDocument params:", err)
			return
		}
		processDiagnostics(params.TextDocument.URI)
	case "shutdown":
		writeMessage(Response{Jsonrpc: "2.0", Id: req.Id, Result: nil})
	case "exit":
		os.Exit(0)
	}
}

func readMessage() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length: ") {
			fmt.Sscanf(line, "Content-Length: %d", &contentLength)
		}
	}
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return "", fmt.Errorf("failed to read message body: %v", err)
	}
	return string(body), nil
}

func writeMessage(resp Response) {
	writeRawMessage(resp)
}

func processDiagnostics(uri string) {
	log.Println("Processing diagnostics for:", uri)
	filePath, err := wio.UriToFilePath(uri)
	if err != nil {
		log.Println("Error converting URI to file path:", err)
		return
	}

	reports, err := analyzer.Analyze(filePath, false, false, false)
	if err != nil {
		log.Println("Error analyzing file:", err)
		return
	}

	var diagnostics []map[string]interface{}
	for _, report := range reports {
		diagnostics = append(diagnostics, map[string]interface{}{
			"range": map[string]interface{}{
				"start": map[string]interface{}{ "line": report.LineStart - 1, "character": report.ColStart - 1 },
				"end": map[string]interface{}{ "line": report.LineEnd - 1, "character": report.ColEnd - 1 },
			},
			"message":  report.Message,
			"severity": mapSeverityToLsp(report.Level),
		})
	}
	publishDiagnostics(uri, diagnostics)
}

func mapSeverityToLsp(level report.REPORT_TYPE) int {
	switch level {
	case report.CRITICAL_ERROR, report.SYNTAX_ERROR, report.NORMAL_ERROR:
		return 1 // Error
	case report.WARNING:
		return 2 // Warning
	case report.INFO:
		return 3 // Information
	default:
		return 3 // Information
	}
}

func publishDiagnostics(uri string, diagnostics []map[string]interface{}) {
	writeRawMessage(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         uri,
			"diagnostics": diagnostics,
		},
	})
}

func writeRawMessage(msg interface{}) {
    data, _ := json.Marshal(msg)
    fullMsg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
    stdout.WriteString(fullMsg)
    stdout.Flush()
}