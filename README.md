# errpipe

<div align="center">

![version](https://img.shields.io/badge/version-0.1.0-blue?style=flat-square)
![license](https://img.shields.io/badge/license-MIT-green?style=flat-square)
![platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square)
![built with](https://img.shields.io/badge/built%20with-go%20%E2%9D%A4-00ADD8?style=flat-square)

**Stop copy/pasting errors. Let your terminal talk to AI.**

[Install](#install) · [Usage](#usage) · [How it Works](#how-it-works) · [Contributing](#contributing)

</div>

---

## What is this

`errpipe` is an interactive shell wrapper built in Go that monitors your commands. When a command fails, it automatically captures the error output and sends it to an LLM (currently Gemini) for immediate analysis and fixes.

Instead of switching to your browser, `errpipe` brings the fix to you.

---

## Install

### From Source (Go required)

Ensure you have Go installed (1.25+ recommended).

```sh
git clone https://github.com/DiamondOsasx/errpipe.git
cd errpipe
go build -o errpipe.exe
```

Add the resulting binary to your system PATH.

---

## Usage

Simply run `errpipe` to start the interactive session.

```sh
errpipe
```

Once inside the `errpipe` shell, run your commands as usual:

```sh
[EP] C:\projects\myapp> go build
# If it fails, errpipe captures the stderr and triggers the AI
```

### Commands
- `--help`: Show help message.
- `--init`: Initialize the application.
- `exit`: Leave the `errpipe` shell.

---

## Supported LLMs

Currently, `errpipe` supports:

| Provider | Status | Integration Method |
|---|---|---|
| Google Gemini | ✅ Supported | Via Gemini CLI |
| Anthropic (Claude) | 🚧 Coming soon | Native API |
| OpenAI | 🚧 Coming soon | Native API |

*Note: For Gemini, ensure you have the Gemini CLI installed and configured on your system.*

---

## How it works

1. **REPL**: `errpipe` acts as a thin wrapper around your default shell (`cmd` on Windows, `sh` on Linux/macOS).
2. **Monitor**: It pipes `stdout` and `stderr` to your terminal while also capturing `stderr` in a buffer.
3. **Trigger**: If a command returns a non-zero exit code, `errpipe` sends the captured `stderr` to the configured AI service.
4. **Automation**: On Windows, it can automatically bring your Gemini CLI window to the front and type the error for you.

---

## Security

`errpipe` runs locally and only sends data to the LLM when a command fails. Be mindful of sensitive information in your error logs.

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
