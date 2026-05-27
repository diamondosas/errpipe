<!--<div align="center">
  <div style="display: flex; align-items: center; justify-content: center;">
    <!--<img src="assets/logo.png" alt="errpipe" width="100" style="margin-right: 20px;">-->
    <h1 style="border-bottom: none;">errpipe</h1>
  </div>-->

  <br>

  ![version](https://img.shields.io/badge/version-0.1.0-blue?style=flat-square)
  ![license](https://img.shields.io/badge/license-MIT-green?style=flat-square)
  ![platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square)
  ![built with](https://img.shields.io/badge/built%20with-go%20%E2%9D%A4-00ADD8?style=flat-square)

  <br>

  **Stop copy-pasting errors. Get instant fixes right in your terminal..**
  <br>
  <br>
  🚀 **No API Key? No Problem!** Includes a completely free mode out of the box.

  <br>

  [Install](#install) · [Usage](#usage) · [How it Works](#how-it-works) · [Contributing](#contributing)

  <br>
  <br>

  <img src="assets/demo.gif" alt="errpipe demo" width="600">
</div>

---

## What is this

`errpipe` is an interactive shell wrapper built in Go that monitors your commands. When a command fails, it automatically captures the error output and sends it to your preferred LLM for immediate analysis and fixes.

Instead of switching to your browser manually, `errpipe` brings the fix to you, either directly in your terminal or by automating your browser.

---

## Install

### Quick Install (Recommended)

**Windows (PowerShell)**
```powershell
irm https://github.com/diamondosas/errpipe/releases/download/vlatest/install.ps1 | iex
```

**macOS / Linux**
```bash
curl -fsSL https://github.com/diamondosas/errpipe/releases/download/vlatest/install.sh | sh
```

<details>
<summary><b>From Source (Go required)</b></summary>

Ensure you have Go installed (1.21+ recommended).

```sh
git clone https://github.com/diamondosas/errpipe
cd errpipe
go mod tidy
go build -o errpipe
```

Add the resulting binary to your system PATH.
</details>

<!--
>[!Note]
>This project is for people who are learning programming or a concept very quickly and dont want to go through the hassle of debugging  .
-->
---

## Usage

Run `errpipe` to start the interactive session. On your first run, it will automatically guide you through a quick setup.

```sh
errpipe
```

Once inside the `errpipe` shell, run your commands as usual. It works with any language or tool:

```sh
# Go
[EP] C:\projects\myapp> go build

# Python
[EP] ~/projects/myapp$ python main.py

# Node.js
[EP] ~/projects/myapp$ node app.js

# Rust
[EP] ~/projects/myapp$ cargo build

# If any command fails, errpipe automatically captures the stderr and triggers the AI
```

### Commands
- `--help`: Show help message.
- `--init`: Reconfigure the application.
- `exit`: Leave the `errpipe` shell.

---

## Supported LLMs & Modes

`errpipe` supports multiple ways to interact with AI:

| Provider | Status | Inline (Streaming) | CLI Mode | Web Mode |
|---|---|---|---|---|
| **Free Mode** | ✅ Supported | ✅ Yes | ❌ No | ❌ No |
| Google Gemini | ✅ Supported | ✅ Yes | ✅ Yes | ✅ Yes |
| Anthropic Claude | ✅ Supported | ✅ Yes | ✅ Yes | ✅ Yes |
| OpenAI ChatGPT | ✅ Supported | ✅ Yes | ✅ Yes | ✅ Yes |

### Integration Modes
1.  **Free Mode**: Completely free, built-in AI mode. **No API key required!**
2.  **Inline CLI Mode**: Streams AI responses directly into your terminal. Requires an **API Key** from the provider.
3.  **CLI Mode**: Interacts with the official CLI tools of the providers installed on your system.
4.  **Web Mode**: Automatically detects your browser, opens the provider's chat page, and types the error message for you using browser automation.

---

## How it works

1.  **REPL**: `errpipe` acts as a thin wrapper around your default shell (`cmd/powershell` on Windows, `sh/bash` on Linux/macOS).
2.  **Monitor**: It pipes `stdout` and `stderr` to your terminal while also capturing `stderr` in a buffer.
3.  **Trigger**: If a command returns a non-zero exit code, `errpipe` sends the captured `stderr` to the configured AI service.
4.  **Analysis**: Depending on your mode, it will either stream the fix directly or automate your environment to get you the answer.

---

## Security

`errpipe` runs locally and only sends data to the LLM when a command fails. Be mindful of sensitive information in your error logs before sending them to public AI models.

---

## Contributing

PRs are welcome.

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

---

## License

MIT.
