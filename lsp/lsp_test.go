package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestReadMessageValid(t *testing.T) {
	// Prepare a valid message with header and body.
	body := "Hello, World!"
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	input := header + body

	reader := bufio.NewReader(strings.NewReader(input))
	msg, err := readMessage(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if msg != body {
		t.Fatalf("Expected message %q but got %q", body, msg)
	}
}

func TestReadMessageMissingContentLength(t *testing.T) {
	// Create a message without a valid Content-Length header.
	input := "\r\nHello, World!"
	reader := bufio.NewReader(strings.NewReader(input))
	_, err := readMessage(reader)
	if err == nil || !strings.Contains(err.Error(), "no content length header found") {
		t.Fatalf("Expected error for missing content length header, but got: %v", err)
	}
}

func TestReadMessageInvalidContentLength(t *testing.T) {
	// Provide an invalid (non-numeric) content length.
	input := "Content-Length: abcd\r\n\r\nHello, World!"
	reader := bufio.NewReader(strings.NewReader(input))
	_, err := readMessage(reader)
	if err == nil || !strings.Contains(err.Error(), "invalid Content-Length") {
		t.Fatalf("Expected error for invalid content length, but got: %v", err)
	}
}

func TestReadMessageIncompleteBody(t *testing.T) {
	// Provide a valid header but an incomplete body.
	body := "Short"
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", 10)
	input := header + body
	reader := bufio.NewReader(strings.NewReader(input))
	_, err := readMessage(reader)
	if err == nil || !strings.Contains(err.Error(), "failed to read message body") {
		t.Fatalf("Expected error for incomplete body, but got: %v", err)
	}
}

func TestWriteMessage(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	resp := Response{
		Jsonrpc: "2.0",
		Id:      1,
		Result: map[string]interface{}{
			"test": "value",
		},
	}
	writeMessage(writer, resp)

	// Ensure writer flushes the content.
	writer.Flush()
	out := buf.String()

	if !strings.HasPrefix(out, "Content-Length:") {
		t.Fatalf("Output does not start with Content-Length header: %s", out)
	}

	// Extract the JSON part after the header.
	parts := strings.SplitN(out, "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatalf("Output format invalid: %s", out)
	}

	var res Response
	if err := json.Unmarshal([]byte(parts[1]), &res); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if res.Id != resp.Id {
		t.Fatalf("Expected Id %d but got %d", resp.Id, res.Id)
	}
}

func TestWriteRawMessage(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
		"params":  map[string]interface{}{},
	}
	writeRawMessage(writer, notification)

	writer.Flush()
	out := buf.String()

	if !strings.HasPrefix(out, "Content-Length:") {
		t.Fatalf("Output does not start with Content-Length header: %s", out)
	}

	// Extract JSON part after header.
	parts := strings.SplitN(out, "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatalf("Output format invalid: %s", out)
	}

	var notif map[string]interface{}
	if err := json.Unmarshal([]byte(parts[1]), &notif); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}
	if notif["method"] != "initialized" {
		t.Fatalf("Expected method 'initialized' but got %v", notif["method"])
	}
}

func TestWriteMessageFlushError(t *testing.T) {
	// Use a writer that errors when WriteString is called.
	errorWriter := &errorWriterStub{}
	bufWriter := bufio.NewWriter(errorWriter)

	resp := Response{
		Jsonrpc: "2.0",
		Id:      1,
		Result:  "test",
	}
	// The function writes a log on error, but we cannot capture logs here.
	// Ensure it does not panic.
	writeMessage(bufWriter, resp)
}

func TestWriteRawMessageFlushError(t *testing.T) {
	// Use a writer that errors for write operations.
	errorWriter := &errorWriterStub{}
	bufWriter := bufio.NewWriter(errorWriter)

	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "testError",
	}
	// The function writes a log on error, but we do not assert on logs.
	writeRawMessage(bufWriter, notification)
}

// errorWriterStub is used to simulate write errors.
type errorWriterStub struct{}

func (e *errorWriterStub) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}
