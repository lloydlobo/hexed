package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit(1)
	}
	fn, ok := commandToFunc[Command(strings.TrimSpace(os.Args[1]))]
	if !ok {
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsageAndExit(1)
	}
	fn()
}

type Command string

const (
	clean Command = "clean"
	build Command = "build"
	serve Command = "serve"
	test  Command = "test"
)

var (
	progn         = "main.wasm"
	commandToFunc = map[Command]func(){
		clean: func() { runCmd("go", "clean"); runCmd("rm", "-f", progn) },
		build: func() { runCmd("go", "build", "-o", progn, "main.go") },
		serve: func() { runCmd("go", "run", "-v", "cmd/serve/serve.go") },
		test: func() {
			runCmd("go", "test", "-v", "./internal/...") // Avoid env error: `fork/exec /tmp/go-build1319822464/b001/internal.test: exec format error`
			_ = os.Setenv("GOOS", "js")
			_ = os.Setenv("GOARCH", "wasm")
			runCmd("go", "test", "-v", "./cmd/...")
		},
	}
)

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if name == "go" {
		switch Command(args[0]) {
		case build, serve:
			cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
		}
	}
	fmt.Printf("Running command: %s %v\n", name, args)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Command '%s' failed with error: %v\n", name, err)
		os.Exit(1)
	}
}

func printUsageAndExit(code int) {
	fmt.Println("Usage: go run cmd/make/make.go <command>")
	fmt.Printf("Available commands: %s\n", strings.Join(commandKeys(), ", "))
	os.Exit(code)
}

func commandKeys() []string {
	a := make([]string, 0, len(commandToFunc))
	for s := range commandToFunc {
		a = append(a, string(s))
	}
	sort.Strings(a) // Sort commands for predictable listing order
	return a
}
