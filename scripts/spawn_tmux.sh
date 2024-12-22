#!/usr/bin/env bash

# chmod +x scripts/spawn_tmux.sh
# ./scripts/spawn_tmux.sh

SESSION_NAME="hextxt"

if tmux has-session -t $SESSION_NAME 2>/dev/null; then
	echo "Session '$SESSION_NAME' already exists. Attaching..."
	tmux attach-session -t $SESSION_NAME
	exit 0
fi

tmux new-session -d -s $SESSION_NAME
tmux send-keys -t $SESSION_NAME "nvim ." C-m

# Watch for changes, rebuild and serve (with python)
tmux split-window -h -t $SESSION_NAME
tmux send-keys -t $SESSION_NAME "find -name '*go' -or -name '*.html' | entr -crs 'GOOS=js GOARCH=wasm go build -o main.wasm main.go; python -m http.server'" C-m

# Watch for changes and test
tmux split-window -v -t $SESSION_NAME
tmux send-keys -t $SESSION_NAME "find -name '*go' -or -name '*.html' | entr -crs 'go test -v internal/hexdump.go internal/hexdump_test.go'" C-m

# Watch for changes and serve (with Go cmd)
tmux split-window -v -t $SESSION_NAME
tmux send-keys -t $SESSION_NAME "make serve-watch" C-m

# Select the first pane
tmux select-pane -t $SESSION_NAME:0.0

tmux attach -t $SESSION_NAME
