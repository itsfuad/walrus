package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

var filename = "file.dummy"

// TestHandleMessageInitialize tests the "initialize" message response.
func TestHandleMessageInitialize(t *testing.T) {
	id := json.RawMessage(`1`)
	msg := Message{
		JsonRPC: "2.0",
		ID:      &id,
		Method:  "initialize",
	}

	resp := handleMessage(&msg)
	if resp == nil {
		t.Fatalf("Expected non-nil response for initialize")
	}

	// Marshal and unmarshal Result to verify its structure.
	marshalled, err := json.Marshal(resp.Result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	var result InitializeResult
	if err := json.Unmarshal(marshalled, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if !result.Capabilities.TextDocumentSync.OpenClose {
		t.Errorf("Expected OpenClose to be true, got false")
	}
	if result.Capabilities.TextDocumentSync.Change != 1 {
		t.Errorf("Expected Change to be 1, got %d", result.Capabilities.TextDocumentSync.Change)
	}
}

// TestWriteAndReadMessage tests that writeMessage and readMessage work together.
func TestWriteAndReadMessage(t *testing.T) {
	origMsg := Message{
		JsonRPC: "2.0",
		Method:  "test",
		Params:  json.RawMessage(`{"key":"value"}`),
	}

	// Write the message to an in-memory buffer.
	var buf bytes.Buffer
	if err := writeMessage(&buf, &origMsg); err != nil {
		t.Fatalf("writeMessage failed: %v", err)
	}

	// Read back the message using a bufio.Reader.
	reader := bufio.NewReader(bytes.NewReader(buf.Bytes()))
	var readMsg Message
	if err := readMessage(reader, &readMsg); err != nil {
		t.Fatalf("readMessage failed: %v", err)
	}

	if readMsg.Method != origMsg.Method {
		t.Errorf("Expected method %s, got %s", origMsg.Method, readMsg.Method)
	}
	if string(readMsg.Params) != string(origMsg.Params) {
		t.Errorf("Expected params %s, got %s", string(origMsg.Params), string(readMsg.Params))
	}
}

// TestHandleMessageDidOpen tests the "textDocument/didOpen" message handling.
func TestHandleMessageDidOpen(t *testing.T) {

	//create a dummy file
	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	defer func () {
		f.Close()
		os.Remove(filename)
	}()

	content := "let a := 3;\n";

	//write some content to the file
	_, err = f.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to dummy file: %v", err)
	}

	params := struct {
		TextDocument struct {
			URI  string `json:"uri"`
			Text string `json:"text"`
		} `json:"textDocument"`
	}{
		TextDocument: struct {
			URI  string `json:"uri"`
			Text string `json:"text"`
		}{
			URI:  filename,
			Text: content,
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	msg := Message{
		JsonRPC: "2.0",
		Method:  "textDocument/didOpen",
		Params:  data,
	}

	// didOpen notifications do not return a response.
	resp := handleMessage(&msg)
	if resp != nil {
		t.Errorf("Expected nil response for didOpen, got non-nil")
	}
}

// TestHandleMessageDidChange tests the "textDocument/didChange" message handling.
func TestHandleMessageDidChange(t *testing.T) {

	//create a dummy file
	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	defer func () {
		f.Close()
		os.Remove(filename)
	}()

	content := "let a := 3;\n";

	//write some content to the file
	_, err = f.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to dummy file: %v", err)
	}

	params := struct {
		TextDocument struct {
			URI  string `json:"uri"`
			Text string `json:"text"`
		} `json:"textDocument"`
	}{
		TextDocument: struct {
			URI  string `json:"uri"`
			Text string `json:"text"`
		}{
			URI:  filename,
			Text: content,
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	msg := Message{
		JsonRPC: "2.0",
		Method:  "textDocument/didChange",
		Params:  data,
	}

	// didChange notifications do not return a response.
	resp := handleMessage(&msg)
	if resp != nil {
		t.Errorf("Expected nil response for didChange, got non-nil")
	}
}

// TestReadMessageEOF tests that readMessage returns io.EOF when no data is available.
func TestReadMessageEOF(t *testing.T) {
	emptyReader := bufio.NewReader(bytes.NewReader([]byte{}))
	var msg Message
	err := readMessage(emptyReader, &msg)
	if err != io.EOF {
		t.Errorf("Expected io.EOF error, got %v", err)
	}
}