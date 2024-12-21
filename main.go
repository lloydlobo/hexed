package main

import (
	"syscall/js"

	"hextxt/internal"
)

const (
	bytesMaxSize = 1024 * 4
)

func main() {
	// Register functions to be called from JavaScript
	js.Global().Set("getMessage", js.FuncOf(getMessage))
	// Block forever
	select {}
}

// getMessage returns converted args[0] result to Javascript.
func getMessage(_ js.Value, args []js.Value) interface{} {
	if len(args) == 0 || args[0].Type() != js.TypeString {
		return "Invalid input"
	}
	s := args[0].String()
	if len(s) > (bytesMaxSize) {
		return "Exceeded input size limit"
	}
	return internal.Hexdump(s, []string{"-C"})
}
