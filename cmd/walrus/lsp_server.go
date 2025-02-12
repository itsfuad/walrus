package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/tools/gopls/lsp"
	"golang.org/x/tools/gopls/lsp/protocol"
	"golang.org/x/tools/gopls/lsp/server"
	"golang.org/x/tools/gopls/span"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to listen on localhost:8080: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	ctx := context.Background()
	stream := lsp.NewStream(conn, conn)
	protocolServer := server.NewServer(ctx, stream)

	if err := protocolServer.Run(ctx); err != nil {
		log.Printf("Failed to run LSP server: %v", err)
	}
}

func initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.Full,
				Save:      &protocol.SaveOptions{IncludeText: true},
			},
		},
	}, nil
}

func didOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	uri := span.URI(params.TextDocument.URI)
	fmt.Printf("File opened: %s\n", uri)
	return nil
}

func didChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := span.URI(params.TextDocument.URI)
	fmt.Printf("File changed: %s\n", uri)
	return nil
}

func didSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	uri := span.URI(params.TextDocument.URI)
	fmt.Printf("File saved: %s\n", uri)
	return nil
}
