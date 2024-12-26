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
	cmd := Command(strings.TrimSpace(os.Args[1]))
	fun, ok := commandToFunc[cmd]
	if !ok {
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsageAndExit(1)
	}
	fun()
}

type Command string

const (
	Clean Command = "clean"
	Build Command = "build"
	Serve Command = "serve"
	Test  Command = "test"
)

var (
	progn         = "main.wasm"
	commandToFunc = map[Command]func(){
		Clean: func() {
			run("go", "clean")
			run("rm", "-f", progn)
		},
		Build: func() {
            //run("cp", "$(go env GOROOT)/misc/wasm/wasm_exec.js", "dist/")
            run("go","env","GOROOT")
            run("cp",strings.Fields("cp /usr/local/go/misc/wasm/wasm_exec.js .")[1:]...) // cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
			run("go", "build", "-o", progn, "main.go")
		},
		Serve: func() {
			run("go", "run", "-v", "cmd/serve/serve.go")
		},
		Test: func() {
			run("go", "test", "-v", "./internal/...") // Avoid env error: `fork/exec /tmp/go-build1319822464/b001/internal.test: exec format error`
			_ = os.Setenv("GOOS", "js")
			_ = os.Setenv("GOARCH", "wasm")
			run("go", "test", "-v", "./cmd/...")
		},
		// Deploy: func() { run("deploy", "my-deploy-script") },
	}
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if name == "go" {
		switch Command(args[0]) {
		case Build, Serve:
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
	for k := range commandToFunc {
		a = append(a, string(k))
	}
	sort.Strings(a) // Sort commands for predictable listing order
	return a
}
