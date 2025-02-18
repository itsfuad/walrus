package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"testing"
)

const (
	CONTENT_LENGTH_HEADER = "Content-Length: "
	SEP                   = "\r\n\r\n"
)

func makeHeader(length int) string {
	return CONTENT_LENGTH_HEADER + strconv.Itoa(length) + SEP
}

// Test readMessage by simulating an input containing a Content-Length header and JSON body.
func TestReadMessage(t *testing.T) {
	body := "Hello"
	header := makeHeader(len(body))
	input := header + body
	reader := bufio.NewReader(strings.NewReader(input))

	msg, err := readMessage(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if msg != body {
		t.Errorf("Expected %q, got %q", body, msg)
	}
}

// Test writeMessage: verifies that the output contains the proper header and JSON response.
func TestWriteMessage(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	resp := Response{
		Jsonrpc: "2.0",
		Id:      1,
		Result:  map[string]string{"test": "value"},
	}
	writeMessage(writer, resp)
	writer.Flush()

	outStr := buf.String()
	if !strings.HasPrefix(outStr, CONTENT_LENGTH_HEADER) {
		t.Errorf("Output missing Content-Length header: %s", outStr)
	}

	parts := strings.SplitN(outStr, SEP, 2)
	if len(parts) != 2 {
		t.Fatalf("Output format invalid: %s", outStr)
	}

	var parsedResp Response
	if err := json.Unmarshal([]byte(parts[1]), &parsedResp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if parsedResp.Id != resp.Id || parsedResp.Jsonrpc != resp.Jsonrpc {
		t.Errorf("Expected %+v, got %+v", resp, parsedResp)
	}
}

// Test writeRawMessage: checks that raw messages are marshaled correctly.
func TestWriteRawMessage(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "testNotification",
		"params":  map[string]string{"key": "value"},
	}
	writeRawMessage(writer, notification)
	writer.Flush()

	outStr := buf.String()
	if !strings.Contains(outStr, `"method":"testNotification"`) {
		t.Errorf("Notification not found in output: %s", outStr)
	}
}

// Test handleInitialize uses a dummy initialize request and checks for the "initialized" notification.
func TestHandleInitialize(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	req := Request{
		Jsonrpc: "2.0",
		Id:      42,
		Method:  "initialize",
		Params:  nil,
	}
	handleInitialize(writer, req)
	writer.Flush()

	outStr := buf.String()
	if !strings.Contains(outStr, `"method":"initialized"`) {
		t.Errorf("Initialized notification missing: %s", outStr)
	}
}

// Test handleShutdown: ensures that a shutdown response is correctly written.
func TestHandleShutdown(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	req := Request{
		Jsonrpc: "2.0",
		Id:      100,
		Method:  "shutdown",
		Params:  nil,
	}
	handleShutdown(writer, req)
	writer.Flush()

	outStr := buf.String()
	parts := strings.SplitN(outStr, SEP, 2)
	if len(parts) != 2 {
		t.Fatalf("Invalid message format: %s", outStr)
	}

	var parsedResp Response
	if err := json.Unmarshal([]byte(parts[1]), &parsedResp); err != nil {
		t.Fatalf("Failed to unmarshal shutdown response: %v", err)
	}
	if parsedResp.Id != req.Id {
		t.Errorf("Expected response id %d, got %d", req.Id, parsedResp.Id)
	}
	if parsedResp.Result != nil {
		t.Errorf("Expected shutdown result to be nil, got %+v", parsedResp.Result)
	}
}

// Test handleConnection using a net.Pipe to simulate a client-server connection.
// This sends an 'initialize' request and expects both a response and a notification.
func TestHandleConnection(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go handleConnection(serverConn)

	// Create an initialize request.
	req := Request{
		Jsonrpc: "2.0",
		Id:      7,
		Method:  "initialize",
		Params:  nil,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	header := makeHeader(len(data))
	message := header + string(data)

	// Write the request to the connection.
	if _, err := clientConn.Write([]byte(message)); err != nil {
		t.Fatalf("Failed to write request: %v", err)
	}

	reader := bufio.NewReader(clientConn)
	var responses []map[string]interface{}

	// Expect two messages: one response and one notification.
	for i := 0; i < 2; i++ {
		responses = append(responses, readResponseFromReader(t, reader))
	}

	// Check the initialize response.
	if id, ok := responses[0]["id"]; !ok || int(id.(float64)) != 7 {
		t.Errorf("Initialize response id mismatch, got: %+v", responses[0])
	}
	// Check the 'initialized' notification.
	if method, ok := responses[1]["method"]; !ok || method != "initialized" {
		t.Errorf("Expected initialized notification, got: %+v", responses[1])
	}
}

// readResponseFromReader reads a single response from the reader and returns it as a map.
func readResponseFromReader(t *testing.T, reader *bufio.Reader) map[string]interface{} {
	msg, err := readMessage(reader)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	return resp
}
