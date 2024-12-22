// Fix editor/lsp import issues:
//
//	sudo apt install direnv			# Install once
//	echo "$(direnv hook zsh)		# Enable direnv once
//
//	echo "export GOOS=js\nexport GOARCH=wasm" > .envrc && direnv allow
//
// Manually:
//
//	GOOS=js GOARCH=wasm nvim .
package main

import (
	"syscall/js"

	"hextxt/internal"
)

const (
	inputMaxSize = (2 << 8) * (2 << 4) // 512*32=>16384
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

	switch {
	case len(s) > (inputMaxSize):
		return "Exceeded input size limit"
	case len(args) == 1:
		return internal.Hexdump(s, []string{})
	default:
		opts := make([]string, 0, len(args)-1)
		for _, v := range args[1:] {
			if v.Type() != js.TypeString {
				return "Invalid non-string option"
			}
			opts = append(opts, v.String())
		}
		return internal.Hexdump(s, opts)
	}
}

// ValueOf returns x as a JavaScript value:
//
//	| Go                     | JavaScript             |
//	| ---------------------- | ---------------------- |
//	| js.Value               | [its value]            |
//	| js.Func                | function               |
//	| nil                    | null                   |
//	| bool                   | boolean                |
//	| integers and floats    | number                 |
//	| string                 | string                 |
//	| []interface{}          | new array              |
//	| map[string]interface{} | new object             |
