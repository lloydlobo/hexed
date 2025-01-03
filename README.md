:q

# hexed

[//]: # (hexadecimal editor)

[![Deploy Pages](https://github.com/lloydlobo/hextxt/actions/workflows/static.yml/badge.svg)](https://github.com/lloydlobo/hextxt/actions/workflows/static.yml)

![preview](./static/preview.png)

## Build

- First, compile Go code to WebAssembly:

  ```sh
  GOOS=js GOARCH=wasm go build -o main.wasm main.go
  ```

- Copy the WebAssembly JavaScript support file from local Go installation:

  ```sh
  cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
  # or
  cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
  ```

  See <https://go.dev/wiki/WebAssembly> guide.

- Final directory structure should look like:

  ```txt
  app/
  ├── index.html
  ├── main.wasm
  └── wasm_exec.js
  ```

- Push these files to GitHub/GitLab and enable Pages.

## Develop

```sh
find -name '*go' -or -name '*.html' -or -name '*.js' | entr -crs "GOOS=js GOARCH=wasm go build -o main.wasm main.go; python -m http.server"
```

## Thanks

- <https://js.hyperui.dev/>
- <https://www.penguinui.com/>

## Related

- <https://binspec.org/>
