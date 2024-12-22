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

/*
	echo hello world | hexdump -C
	00000000  68 65 6c 6c 6f 20 77 6f  72 6c 64 0a              |hello world.|
	0000000c

	echo hello world | hexdump
	0000000 6568 6c6c 206f 6f77 6c72 0a64
	000000c

	"reversal"

	- The default hexdump format groups bytes into 16-bit words and uses
	  little-endian order.
	- hexdump -C shows bytes individually in natural order, avoiding the
	  confusion.

		got := hex.EncodeToString([]byte("hello world\n"))
		want := "68656c6c6f20776f726c640a"
		isMatch := reflect.DeepEqual(got, want)
		slog.Info("hex.EncodeToString", "got", got, "want", want, "isMatch", isMatch)
*/

/*
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title></title>
        <link href="css/style.css" rel="stylesheet">
    </head>
    <body>

    </body>
</html>

00000000   21 44 4f 43 54 59 50 45  20 68 74 6d 6c 0a 68 74   |!DOCTYPE html.ht|
00000010   6d 6c 20 6c 61 6e 67 3d  22 65 6e 22 0a 20 20 20   |ml lang="en".   |
00000020   20 68 65 61 64 0a 20 20  20 20 20 20 20 20 6d 65   | head.        me|
00000030   74 61 20 63 68 61 72 73  65 74 3d 22 55 54 46 2d   |ta charset="UTF-|
00000040   38 22 0a 20 20 20 20 20  20 20 20 6d 65 74 61 20   |8".        meta |
00000050   6e 61 6d 65 3d 22 76 69  65 77 70 6f 72 74 22 20   |name="viewport" |
00000060   63 6f 6e 74 65 6e 74 3d  22 77 69 64 74 68 3d 64   |content="width=d|
00000070   65 76 69 63 65 2d 77 69  64 74 68 2c 20 69 6e 69   |evice-width, ini|
00000080   74 69 61 6c 2d 73 63 61  6c 65 3d 31 22 0a 20 20   |tial-scale=1".  |
00000090   20 20 20 20 20 20 74 69  74 6c 65 2f 74 69 74 6c   |      title/titl|
000000a0   65 0a 20 20 20 20 20 20  20 20 6c 69 6e 6b 20 68   |e.        link h|
000000b0   72 65 66 3d 22 63 73 73  2f 73 74 79 6c 65 2e 63   |ref="css/style.c|
000000c0   73 73 22 20 72 65 6c 3d  22 73 74 79 6c 65 73 68   |ss" rel="stylesh|
000000d0   65 65 74 22 0a 20 20 20  20 2f 68 65 61 64 0a 20   |eet".    /head. |
000000e0   20 20 20 62 6f 64 79 0a  20 20 20 20 0a 20 20 20   |   body.    .   |
000000f0   20 2f 62 6f 64 79 0a 2f  68 74 6d 6c               | /body./html|
*/
