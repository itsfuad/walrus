package lexer

// Position represents a specific location in the source code with line, column, and index information.
type Position struct {
	Line   int // Line number in the source code.
	Column int // Column number in the source code.
	Index  int // Index in the source code.
}

// Advance updates the Position by advancing it based on the bytes in the provided string.
// It increments the line number for newline bytes and the column number for other bytes.
// The index is incremented for each byte in the string.
//
// Parameters:
// - toSkip: A string containing bytes to advance the position by.
//
// Returns:
// - A pointer to the updated Position.
func (p *Position) Advance(toSkip string) *Position {
	for _, char := range toSkip {
		if char == '\n' {
			p.Line++
			p.Column = 1
		} else {
			p.Column++
		}
		p.Index++
	}
	return p
}
