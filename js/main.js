// Go: globalThis.Go
// dist/ wasm_exec. js
let go = new Go();

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
    .then((res) => {
      go.run(res.instance)
    })
    .catch((err) => {
      console.error("Failed to load Wasm:", err);
    });

/**
 * Triggered by a user event such as a button's `onclick` trigger.
 * Calls a Go function exposed via WebAssembly to process the input text.
 * Updates the output element with the result.
 * @return {void}
 */
function processInput() { // PERF: new ArrayBuffer(); new TypedArray();
  const textAreaElement = /** @type{HTMLTextAreaElement} */ document.getElementById("input");
  const selectElement = /** @type{HTMLSelectElement} */ document.getElementById("idSelectOption");

  /**@type{string|undefined}*/
  let result;
  if (selectElement.value === "") {
    result = /** @type{string} */ getMessage(textAreaElement.value);
  } else {
    result = /** @type{string} */ getMessage(textAreaElement.value, selectElement.value);
  }

  let outputElement = /** @type{HTMLOutputElement} */ document.getElementById("outputBytes");
  outputElement.textContent = result;
}

/**
 * Copies the content of the output element to the clipboard.
 * @return {void}
 */
function copyOutput() {
  const output = /** @type{HTMLOutputElement} */ document.getElementById("outputBytes");
  copyToClipboard(output.textContent);
}

/**
 * copyToClipboard copies a given string to the clipboard.
 * See also:
 * - https://stackoverflow.com/a/60292243
 * - https://www.freecodecamp.org/news/copy-text-to-clipboard-javascript/
 * @param {string} s - The string to copy.
 * @return {void}
 */
function copyToClipboard(s) {
  // deno-lint-ignore no-window-prefix, no-window
  if (window.isSecureContext) { // Clipboard API is only supported for pages served over HTTPS.
    const permissionDes = /** @type{PermissionDescriptor} */ {name: "persistent-storage"} // Why no "write-on-clipboard" ???
    navigator.permissions.query(permissionDes).then((result) => {
      if (result.name === "granted" || result.state === "prompt") {
        console.log("'write' access granted");
      }
    });
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(s).then(
          () => {
            /* Resolved - text copied to clipboard successfully */
          },
          (err) => {
            console.error("Failed to copy text to the clipboard:", err);
          },
      );
      return;
    }
  }

  const textarea = /** @type{HTMLTextAreaElement} */ document.createElement("textarea");
  textarea.value = s;
  textarea.style.position = "fixed"; // Avoid scrolling to page bottom
  document.body.appendChild(textarea);
  textarea.focus(); // 30ms
  textarea.select(); // 3ms

  try {
    /** @deprecated */
    const success = document.execCommand("copy");
    if (!success) {
      console.error("Failed to copy");
    }
  } catch (err) {
    console.error("Fallback copy failed:", err);
    alert("Manual copy: " + s);
  }
  document.body.removeChild(textarea);
}
