package parser

import (
	"testing"
	"walrus/frontend/ast"
	"walrus/frontend/builtins"
	"walrus/frontend/lexer"
)

func TestGetBP(t *testing.T) {
	// Setup
	BPLookup[lexer.PLUS_TOKEN] = ADDITIVE_BP

	tests := []struct {
		token builtins.TOKEN_KIND
		want  BINDING_POWER
	}{
		{lexer.PLUS_TOKEN, ADDITIVE_BP},
		{lexer.MINUS_TOKEN, DEFAULT_BP}, // Not set in BPLookup, should return DEFAULT_BP
	}

	for _, tt := range tests {
		t.Run(string(tt.token), func(t *testing.T) {
			if got := GetBP(tt.token); got != tt.want {
				t.Errorf("GetBP(%v) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestLed(t *testing.T) {
	// Setup
	handler := func(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node { return nil }
	led(lexer.PLUS_TOKEN, ADDITIVE_BP, handler)

	if got, ok := LEDLookup[lexer.PLUS_TOKEN]; !ok || got == nil {
		t.Errorf("LEDLookup[%v] not set correctly", lexer.PLUS_TOKEN)
	}

	if got := BPLookup[lexer.PLUS_TOKEN]; got != ADDITIVE_BP {
		t.Errorf("BPLookup[%v] = %v, want %v", lexer.PLUS_TOKEN, got, ADDITIVE_BP)
	}
}

func TestNud(t *testing.T) {
	// Setup
	handler := func(p *Parser) ast.Node { return nil }
	nud(lexer.MINUS_TOKEN, handler)

	if got, ok := NUDLookup[lexer.MINUS_TOKEN]; !ok || got == nil {
		t.Errorf("NUDLookup[%v] not set correctly", lexer.MINUS_TOKEN)
	}
}

func TestStmt(t *testing.T) {
	// Setup
	handler := func(p *Parser) ast.Node { return nil }
	stmt(lexer.LET_TOKEN, handler)

	if got, ok := STMTLookup[lexer.LET_TOKEN]; !ok || got == nil {
		t.Errorf("STMTLookup[%v] not set correctly", lexer.LET_TOKEN)
	}
}

func TestBindLookupHandlers(t *testing.T) {
	// Call the function to bind handlers
	bindLookupHandlers()

	// Check a few handlers to ensure they are set correctly
	if got, ok := LEDLookup[lexer.EQUALS_TOKEN]; !ok || got == nil {
		t.Errorf("LEDLookup[%v] not set correctly", lexer.EQUALS_TOKEN)
	}

	if got, ok := NUDLookup[lexer.IDENTIFIER_TOKEN]; !ok || got == nil {
		t.Errorf("NUDLookup[%v] not set correctly", lexer.IDENTIFIER_TOKEN)
	}

	if got, ok := STMTLookup[lexer.LET_TOKEN]; !ok || got == nil {
		t.Errorf("STMTLookup[%v] not set correctly", lexer.LET_TOKEN)
	}
}
