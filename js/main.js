const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then((res) => {
    go.run(res.instance);
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
function processInput() {
  // PERF
  // new ArrayBuffer();
  // new TypedArray();
  const textarea = document.getElementById("input");
  const options = document.getElementById("options");
  let hexdump;
  if (options.value === "") {
    // Default for Option "hex"
    hexdump = getMessage(textarea.value);
  } else {
    hexdump = getMessage(textarea.value, options.value);
  }
  document.getElementById("outputBytes").textContent = hexdump;
}

/**
 * Copies the content of the output element to the clipboard.
 * @return {void}
 */
function copyOutput() {
  const output = document.getElementById("outputBytes");
  const text = output.textContent;
  copyToClipboard(text);
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
  if (window.isSecureContext) {
    // Clipboard API is only supported for pages served over HTTPS.
    navigator.permissions
      .query({ name: "write-on-clipboard" })
      .then((result) => {
        if (result.name == "granted" || result.state == "prompt") {
          console.log("'write' access granted");
        }
      });
    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(s).then(
        () => {
          /* Resolved - text copied to clipboard successfully */
        },
        (err) => {
          console.error("Failed to copy text:", err);
          /* Rejected - text failed to copy to the clipboard */
        },
      );
      return;
    }
  }

  const textarea = document.createElement("textarea");
  textarea.value = s;
  textarea.style.position = "fixed"; // Avoid scrolling to page bottom
  document.body.appendChild(textarea);

  textarea.focus(); // 30ms
  textarea.select();// 3ms

  try {
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
