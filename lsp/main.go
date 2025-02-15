package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"walrus/compiler/analyzer"
	"walrus/compiler/report"
)

// UriToFilePath converts a file:// URI to a platform-specific file path.
func UriToFilePath(uri string) (string, error) {
	// Parse the URI
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	// Handle Windows paths (e.g., file:///C:/path/to/file)
	if parsed.Path[0] == '/' && len(parsed.Path) > 3 && parsed.Path[2] == ':' {
		parsed.Path = parsed.Path[1:]
	}

	// Decode URL-encoded characters (e.g., %3A -> :)
	filePath, err := url.PathUnescape(parsed.Path)
	if err != nil {
		return "", err
	}

	// Convert to platform-specific file path
	return filepath.FromSlash(filePath), nil
}

//log to file
func init() {
	f, err := os.OpenFile("lsp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

// LSP Request and Response structures
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

		switch req.Method {
		case "initialize":
			response := Response{
				Jsonrpc: "2.0",
				Id:      req.Id,
				Result: map[string]interface{}{
					"capabilities": map[string]interface{}{
						"textDocumentSync": 1, // Full text sync
					},
				},
			}
			writeMessage(response)
		
			// Send initialized notification
			notification := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "initialized",
				"params":  map[string]interface{}{},
			}

			writeRawMessage(notification)

		case "textDocument/didOpen", "textDocument/didChange":
			var params struct {
				TextDocument struct {
					URI string `json:"uri"`
				} `json:"textDocument"`
			}
			if err := json.Unmarshal(req.Params, &params); err != nil {
				log.Println("Error parsing textDocument params:", err)
				continue
			}
			processDiagnostics(params.TextDocument.URI)

		case "shutdown":
			writeMessage(Response{Jsonrpc: "2.0", Id: req.Id, Result: nil})

		case "exit":
			os.Exit(0)
		}
	}
}


func readMessage() (string, error) {

    reader := bufio.NewReader(os.Stdin)
    
    // Read headers
    contentLength := 0
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return "", err
        }
        line = strings.TrimSpace(line)
        
        if line == "" { // End of headers
            break
        }

        if strings.HasPrefix(line, "Content-Length: ") {
            fmt.Sscanf(line, "Content-Length: %d", &contentLength)
        }
    }

    // Read body
    body := make([]byte, contentLength)
    _, err := io.ReadFull(reader, body)
    if err != nil {
        return "", fmt.Errorf("failed to read message body: %v", err)
    }

    return string(body), nil
}


// Write an LSP message
func writeMessage(resp Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	msg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
	if _, err := stdout.WriteString(msg); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
	stdout.Flush()
}

func processDiagnostics(uri string) {
	log.Println("Processing diagnostics for:", uri)

	// Convert URI to file path
	filePath, err := UriToFilePath(uri)
	if err != nil {
		log.Println("Error converting URI to file path:", err)
		// Handle URI to file path conversion error
		writeMessage(Response{
			Jsonrpc: "2.0",
			Id:      0,
			Error: &LspError{
				Code:    -32603, // Internal error code
				Message: fmt.Sprintf("Error converting URI to file path: %v", err),
			},
		})
		return
	}

	// Analyze the file for diagnostics
	reports, err := analyzer.Analyze(filePath, false, false, false)
	if err != nil {
		log.Println("Error analyzing file:", err)
		// Handle analysis error
		writeMessage(Response{
			Jsonrpc: "2.0",
			Id:      0,
			Error: &LspError{
				Code:    -32603, // Internal error code
				Message: fmt.Sprintf("Error analyzing file: %v", err),
			},
		})
		return
	}

	// Convert the report to diagnostics
	diagnostics := []map[string]interface{}{}
	for _, report := range reports {
		diagnostic := map[string]interface{}{
			"range": map[string]interface{}{
				"start": map[string]interface{}{
					"line": report.LineStart - 1, // LSP uses 0-based indices
					"character": report.ColStart - 1,
				},
				"end": map[string]interface{}{
					"line": report.LineEnd - 1,
					"character": report.ColEnd - 1,
				},
			},
			"message": report.Message,
			"severity": mapSeverityToLsp(report.Level),
		}
		diagnostics = append(diagnostics, diagnostic)
	}

	// Send diagnostics to LSP client
	response := Response{
		Jsonrpc: "2.0",
		Id:      0,
		Result: map[string]interface{}{
			"uri":        uri,
			"diagnostics": diagnostics,
		},
	}
	
	writeMessage(response)
}

// Convert report severity to LSP severity
func mapSeverityToLsp(level report.REPORT_TYPE) int {
	switch level {
	case report.CRITICAL_ERROR, report.SYNTAX_ERROR:
		return 1 // Error
	case report.NORMAL_ERROR:
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
    notification := map[string]interface{}{
        "jsonrpc": "2.0",
        "method":  "textDocument/publishDiagnostics",
        "params": map[string]interface{}{
            "uri":         uri,
            "diagnostics": diagnostics,
        },
    }
    writeRawMessage(notification)
}

func writeRawMessage(msg interface{}) {
    data, _ := json.Marshal(msg)
    fullMsg := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
    stdout.WriteString(fullMsg)
    stdout.Flush()
}