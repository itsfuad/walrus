package builtins

import (
	"testing"
)

func TestGetBitSize(t *testing.T) {
	tests := []struct {
		kind     string
		expected uint8
	}{
		{INT8, 8},
		{UINT8, 8},
		{BYTE, 8},
		{INT16, 16},
		{UINT16, 16},
		{INT32, 32},
		{UINT32, 32},
		{FLOAT32, 32},
		{INT64, 64},
		{UINT64, 64},
		{FLOAT64, 64},
		{STRING, 0},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			if got := GetBitSize(tt.kind); got != tt.expected {
				t.Errorf("GetBitSize(%s) = %d; want %d", tt.kind, got, tt.expected)
			}
		})
	}
}

func TestIsSigned(t *testing.T) {
	tests := []struct {
		kind     string
		expected bool
	}{
		{INT8, true},
		{INT16, true},
		{INT32, true},
		{INT64, true},
		{UINT8, false},
		{UINT16, false},
		{UINT32, false},
		{UINT64, false},
		{BYTE, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			if got := IsSigned(tt.kind); got != tt.expected {
				t.Errorf("IsSigned(%s) = %v; want %v", tt.kind, got, tt.expected)
			}
		})
	}
}

func TestIsUnsigned(t *testing.T) {
	tests := []struct {
		kind     string
		expected bool
	}{
		{INT8, false},
		{INT16, false},
		{INT32, false},
		{INT64, false},
		{UINT8, true},
		{UINT16, true},
		{UINT32, true},
		{UINT64, true},
		{BYTE, true},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			if got := IsUnsigned(tt.kind); got != tt.expected {
				t.Errorf("IsUnsigned(%s) = %v; want %v", tt.kind, got, tt.expected)
			}
		})
	}
}
