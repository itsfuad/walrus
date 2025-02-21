package constants

type MESSAGE string

const (
	// Error and warning messages
	UNEXPECTED_KEYWORD MESSAGE = "unexpected keyword"
	UNEXPECTED_TOKEN   MESSAGE = "unexpected token"

	UNDEFINED_IDENTIFIER MESSAGE = "undefined identifier"
	UNDEFINED_TYPE       MESSAGE = "undefined type"

	UNDEFINED_PROP   MESSAGE = "undefined property"
	UNDEFINED_METHOD MESSAGE = "undefined method"
)
