package internal

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	base           int    = 16
	hexChunkSep    string = " "
	hexChunkSepLen int    = len(hexChunkSep)
	hexChunkSize   int    = 2
)

// Hexdump converts a string to its hexadecimal representation.
func Hexdump(s string, opts []string) string {
	out := chunkString(hex.EncodeToString([]byte(s)), hexChunkSep, hexChunkSize)
	if len(opts) == 0 {
		return out
	}
	switch cmd := opts[0]; cmd {
	case "-C":
		size := 16
		h := chunkHexBytes([]byte(out), size*(hexChunkSize+hexChunkSepLen))
		c := Chunks{
			offset: chunkHexOffset(len(h), size),
			hex:    h,
			ascii:  chunkHexASCII(s, size),
		}
		return encodeChunks(c)
	default:
		panic("unimplemented command argument: " + cmd)
	}
}

func chunkString(s, sep string, size int) string {
	var a []string
	n := len(s)
	for i := 0; i < n; i += size {
		j := i + size
		if j > n {
			j = n
		}
		a = append(a, s[i:j])
	}
	return strings.Join(a, sep)
}

func chunkHexOffset(n int, size int) []string {
	a := make([]string, 0, n)
	for i := range n {
		offset := strconv.FormatInt(int64(i*size), base)         // Increment by chunkSize
		a = append(a, strings.Repeat("0", 8-len(offset))+offset) // Pad with leading zeroes
	}
	return a
}

func chunkHexBytes(b []byte, size int) [][]byte {
	var a [][]byte
	n := len(b)
	for i := 0; i < n; i += size {
		j := i + size
		if j > n {
			j = n
		}
		a = append(a, bytes.Map(toValidASCIIRune, b[i:j]))
	}
	return a
}

// UnescapeString unescapes entities like "&lt;" to become "\<". It unescapes a
// larger range of entities than EscapeString escapes. substr = html.UnescapeString(substr)
func chunkHexASCII(s string, size int) []string {
	var a []string
	n := len(s)
	for i := 0; i < n; i += size {
		j := i + size
		if j > n {
			j = n
		}
		substr := string(s[i:j])
		substr = strings.Map(toValidASCIIRune, substr)
		a = append(a, substr)
	}
	return a
}

// toValidASCIIRune replaces non-printable ASCII characters with '.'.
func toValidASCIIRune(r rune) rune {
	if r < 32 || r > 126 {
		return '.'
	}
	return r
}

type Chunks struct { // size=72 (0x48)
	offset []string
	hex    [][]byte
	ascii  []string
}

// encodeChunks combines the offset, hex, and ASCII chunks into a string.
func encodeChunks(c Chunks) string {
	var sb strings.Builder
	for i := range max(len(c.offset), len(c.hex), len(c.ascii)) {
		if len(c.offset) > 0 && i < len(c.offset) {
			sb.WriteString(safeIndex(c.offset, i, "invalid"))
		}
		if len(c.hex) > 0 && i < len(c.hex) {
			hex := safeIndex(c.hex, i, []byte("invalid"))
			if delta := len(c.hex[0]) - len(hex); delta > 0 {
				for range delta {
					hex = append(hex, ' ')
				}
			}
			sb.WriteString("   ")
			sb.Write(hex[:len(hex)/2])
			sb.WriteString(" ")
			sb.Write(hex[len(hex)/2:])
			sb.WriteString("  ")
		}
		if len(c.ascii) > 0 && i < len(c.ascii) {
			sb.WriteString("|")
			sb.WriteString(safeIndex(c.ascii, i, "invalid"))
			sb.WriteString("|")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// safeIndex returns the element at the given index or the default value
// if the index is out of bounds.
func safeIndex[T any](a []T, index int, defaultValue T) T {
	if index < len(a) {
		return a[index]
	}
	return defaultValue
}
