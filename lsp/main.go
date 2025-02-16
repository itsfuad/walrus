package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"walrus/compiler/report"
	"walrus/compiler/wio"

	"walrus/compiler/analyzer"
)

// LSP structures remain the same
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

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds) // Add microseconds to log timestamps

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("PORT:%d\n", port)
	os.Stdout.Sync() // Force flush the port number

	log.Printf("LSP Server listening on port %d", port)

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Failed to accept connection: %v", err)
	}
	defer conn.Close()

	log.Printf("Client connected from: %s", conn.RemoteAddr())
	handleConnection(conn)
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg, err := readMessage(reader)
		if err == io.EOF {
			log.Printf("Client disconnected")
			break
		}
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		if msg == "" {
			log.Printf("Empty message received, skipping")
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(msg), &req); err != nil {
			log.Printf("Invalid JSON message %q: %v", msg, err)
			continue
		}

		switch req.Method {
		case "initialize":
			handleInitialize(writer, req)
		case "textDocument/didOpen", "textDocument/didChange", "textDocument/didSave":
			handleTextDocumentChange(writer, req)
		case "shutdown":
			handleShutdown(writer, req)
		case "exit":
			handleExit(conn)
		default:
			handleUnknownMethod(req)
		}
	}
}

func handleInitialize(writer *bufio.Writer, req Request) {
	resp := Response{
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result: map[string]interface{}{
			"capabilities": map[string]interface{}{
				"textDocumentSync": 1,
			},
		},
	}
	writeMessage(writer, resp)
	
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
		"params":  map[string]interface{}{},
	}
	writeRawMessage(writer, notification)
}

func handleTextDocumentChange(writer *bufio.Writer, req Request) {
	var params map[string]interface{}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		log.Printf("Invalid params: %v", err)
		return
	}

	uri, ok := params["textDocument"].(map[string]interface{})["uri"].(string)
	if !ok {
		log.Printf("Invalid uri: %v", params)
		return
	}

	processDiagnostics(writer, uri)
}

func handleShutdown(writer *bufio.Writer, req Request) {
	resp := Response{
		Jsonrpc: "2.0",
		Id:      req.Id,
		Result:  nil,
	}
	writeMessage(writer, resp)
}

func handleExit(conn net.Conn) {
	log.Printf("Client requested exit")
	conn.Close()
}

func handleUnknownMethod(req Request) {
	log.Printf("Unknown method: %v", req.Method)
}

func readMessage(reader *bufio.Reader) (string, error) {
	contentLength := 0
	
	// Read headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		
		// Trim both \r and \n
		line = strings.TrimRight(line, "\r\n")
		
		if line == "" { // End of headers
			break
		}

		if strings.HasPrefix(line, "Content-Length: ") {
			lengthStr := strings.TrimPrefix(line, "Content-Length: ")
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return "", fmt.Errorf("invalid Content-Length: %v", err)
			}
			log.Printf("Content length: %d", contentLength)
		}
	}

	if contentLength == 0 {
		return "", fmt.Errorf("no content length header found")
	}

	// Read body
	body := make([]byte, contentLength)
	n, err := io.ReadFull(reader, body)
	if err != nil {
		return "", fmt.Errorf("failed to read message body (read %d of %d bytes): %v", n, contentLength, err)
	}

	bodyStr := string(body)
	log.Printf("Received message body: %q", bodyStr)
	return bodyStr, nil
}


func writeMessage(writer *bufio.Writer, resp Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	msg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	if _, err := writer.WriteString(msg); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
	writer.Flush()
}

func writeRawMessage(writer *bufio.Writer, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}
	
	fullMsg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	if _, err := writer.WriteString(fullMsg); err != nil {
		log.Printf("Failed to write message: %v", err)
	}
	writer.Flush()
}

func processDiagnostics(writer *bufio.Writer, uri string) {
	log.Println("Processing diagnostics for:", uri)

	filePath, err := wio.UriToFilePath(uri)
	if err != nil {
		log.Println("Error converting URI to file path:", err)
		return
	}

	log.Println("File path:", filePath)

	report, err := analyzer.Analyze(filePath, false, false, false)
	if err != nil {
		log.Println("Error analyzing file:", err)
	}

	diagnostics := make([]map[string]interface{}, 0)

	log.Printf("Found %d problems\n", len(report))
	
	for _, r := range report {
		diagnostics = append(diagnostics, map[string]interface{}{
			"range": map[string]interface{}{
				"start": map[string]int{"line": r.LineStart - 1, "character": r.ColStart - 1},
				"end":   map[string]int{"line": r.LineEnd - 1, "character": r.ColEnd - 1},
			},
			"message":  r.Message,
			"severity": getSeverity(r.Level),
		})
	}

	publishDiagnostics(writer, uri, diagnostics)
}

func getSeverity(level report.REPORT_TYPE) int {
	switch level {
	case report.CRITICAL_ERROR, report.SYNTAX_ERROR, report.NORMAL_ERROR:
		return 1
	case report.WARNING:
		return 2
	case report.INFO:
		return 3
	default:
		return 4
	}
}

func publishDiagnostics(writer *bufio.Writer, uri string, diagnostics []map[string]interface{}) {
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         uri,
			"diagnostics": diagnostics,
		},
	}
	writeRawMessage(writer, notification)
}