package builtins

// common types

type TC_TYPE string
type TOKEN_KIND string
type PARSER_TYPE string

const (
	INT8      = "i8"
	INT16     = "i16"
	INT32     = "i32"
	INT64     = "i64"
	UINT8     = "u8"
	UINT16    = "u16"
	UINT32    = "u32"
	UINT64    = "u64"
	FLOAT32   = "f32"
	FLOAT64   = "f64"
	STRING    = "str"
	BYTE      = "byte"
	BOOL      = "bool"
	NULL      = "null"
	FUNCTION  = "fn"
	STRUCT    = "struct"
	INTERFACE = "interface"
	ARRAY     = "array"
	MAP       = "map"
	VOID      = "void"
)

type Searchable interface {
	TC_TYPE | PARSER_TYPE
}

func GetBitSize[T Searchable](kind T) uint8 {
	switch kind {
	case INT8, UINT8, BYTE:
		return 8
	case INT16, UINT16:
		return 16
	case INT32, UINT32, FLOAT32:
		return 32
	case INT64, UINT64, FLOAT64:
		return 64
	default:
		return 0
	}
}

func IsSigned[T Searchable](kind T) bool {
	switch kind {
	case INT8, INT16, INT32, INT64:
		return true
	default:
		return false
	}
}

func IsUnsigned[T Searchable](kind T) bool {
	switch kind {
	case UINT8, UINT16, UINT32, UINT64, BYTE:
		return true
	default:
		return false
	}
}
