package lexer

type Position struct {
	Line	int
	Column	int
	Index	int
}

func (p *Position) Advance(toSkip string) *Position {
	for _, char := range toSkip { // a rune is an alias for int32
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