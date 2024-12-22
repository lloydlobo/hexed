package main

import (
	"syscall/js"

	"hextxt/internal"
)

const (
	inputMaxSize = 1024 * 4
)

func main() {
	js.Global().Set("getMessage", js.FuncOf(getMessage)) // Register functions to be called from JavaScript
	select {}                                            // Block forever
}

// getMessage returns converted args[0] result to Javascript.
func getMessage(_ js.Value, args []js.Value) interface{} {
	if len(args) == 0 || args[0].Type() != js.TypeString {
		return "Invalid input"
	}
	s := args[0].String()
	if len(s) > (inputMaxSize) {
		return "Exceeded input size limit"
	}
	return internal.Hexdump(s, []string{"-C"})
}
