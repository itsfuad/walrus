package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestCreateServer(t *testing.T) {
	// Capture stdout to get the printed port.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Run createServer in a goroutine.
	connCh := make(chan net.Conn, 1)
	errCh := make(chan error, 1)
	go func() {
		conn, err := createServer()
		if err != nil {
			errCh <- err
			return
		}
		connCh <- conn
	}()

	// Allow some time for createServer to start and write the port.
	time.Sleep(100 * time.Millisecond)
	// Close the write end so we can read all output.
	w.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	output := buf.String()
	t.Logf("Captured stdout: %s", output)

	// Extract port from output.
	re := regexp.MustCompile(`PORT:(\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatal("failed to extract port from output")
	}
	port := matches[1]

	// Dial the server using the extracted port.
	c, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	c.Close()

	// Wait for createServer to return.
	select {
	case srvConn := <-connCh:
		if srvConn == nil {
			t.Fatal("received nil connection")
		}
		// The connection is closed by createServer via defer.
	case err := <-errCh:
		t.Fatalf("createServer error: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for createServer to return")
	}
}
