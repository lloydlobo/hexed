package internal

import (
	"reflect"
	"testing"
)

func TestTxtToHex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "hello world",
			input: "hello world",
			want:  "68 65 6c 6c 6f 20 77 6f 72 6c 64",
		},
		{
			name:  "Alphabets",
			input: "abcdefghijklmnopqrstuvwxyz",
			want:  "61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a",
		},
		{
			name:  "newline",
			input: "\n",
			want:  "0a",
		},
		{
			name:  "package internal",
			input: "package internal",
			want:  "70 61 63 6b 61 67 65 20 69 6e 74 65 72 6e 61 6c",
		},
		{
			name: `package internal with newline comment`,
			input: `package internal
\\`,
			want: "70 61 63 6b 61 67 65 20 69 6e 74 65 72 6e 61 6c 0a 5c 5c",
		},
		{
			name: `package internal with newline + newline comment`,
			input: `package internal

\\`,
			want: "70 61 63 6b 61 67 65 20 69 6e 74 65 72 6e 61 6c 0a 0a 5c 5c",
		},
		{
			name: `package internal with newline + newline imports`,
			input: `package internal

import "fmt"`,
			want: "70 61 63 6b 61 67 65 20 69 6e 74 65 72 6e 61 6c 0a 0a 69 6d 70 6f 72 74 20 22 66 6d 74 22",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{}
			if got := Hexdump(tt.input, args); got != tt.want {
				t.Errorf("expected '%v', got '%v'", tt.want, got)
			}
		})
	}
}

func TestGenerateHexOffsets(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		output []string
	}{
		{
			name:   "Zero rows",
			input:  0,
			output: []string{},
		},
		{
			name:   "Single row",
			input:  1,
			output: []string{"00000000"},
		},
		{
			name:   "Multiple rows",
			input:  5,
			output: []string{"00000000", "00000010", "00000020", "00000030", "00000040"},
		},
		{
			name:   "10 rows",
			input:  10,
			output: []string{"00000000", "00000010", "00000020", "00000030", "00000040", "00000050", "00000060", "00000070", "00000080", "00000090"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunkSize := 16
			got := chunkHexOffset(tt.input, chunkSize)
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("input = %v, got = %v; want %v", tt.input, got, tt.output)
			}
		})
	}
}
