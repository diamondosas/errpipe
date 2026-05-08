# errpipe

<!-- PLACE YOUR DEMO GIF/IMAGE HERE -->
<!-- Example: ![demo](./assets/demo.gif) -->

<div align="center">

![version](https://img.shields.io/badge/version-0.1.0-blue?style=flat-square)
![license](https://img.shields.io/badge/license-MIT-green?style=flat-square)
![platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square)
![built with](https://img.shields.io/badge/built%20with-terminal%20%E2%9D%A4-black?style=flat-square)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat-square)

**It's 2026. Stop copy/pasting errors into ChatGPT.**

[Install](#install) · [Usage](#usage) · [Config](#config) · [Contributing](#contributing)

</div>

---

## What is this

Your terminal throws an error. You copy it. You open the browser. You paste it. You wait.

That is 4 steps too many.

`errpipe` catches the error from your terminal and sends it to any LLM you want. One pipe. Done.

```sh
npm run build 2>&1 | errpipe
```

That is it.

---

## Why

- Errors are inevitable. That is just engineering.
- Copy/paste is not a workflow. It is a tax on your time.
- LLMs are good at this. Use them properly.

> Full-time vibe coders will not enjoy this app. Errors are a normal part of building software. This tool is for engineers who know that.

---

## Install

**From release (recommended)**

Download the binary for your platform from the [Releases](../../releases) folder.

```sh
# macOS / Linux
chmod +x errpipe
sudo mv errpipe /usr/local/bin/
```

```sh
# Windows
# Add the .exe to your PATH
```

**From source**

```sh
git clone https://github.com/yourusername/errpipe.git
cd errpipe
# build instructions here
```

---

## Usage

Pipe any stderr output directly into `errpipe`.

```sh
# Basic
python app.py 2>&1 | errpipe

# Specific command
cargo build 2>&1 | errpipe

# Save the fix to a file
npm test 2>&1 | errpipe --output fix.md
```

---

## Config

Set your LLM of choice. `errpipe` does not pick one for you.

```sh
errpipe config --llm openai --key YOUR_API_KEY
errpipe config --llm anthropic --key YOUR_API_KEY
errpipe config --llm groq --key YOUR_API_KEY
```

Config is stored at `~/.errpipe/config.json`.

```json
{
  "llm": "openai",
  "model": "gpt-4o",
  "api_key": "sk-..."
}
```

---

## Supported LLMs

| Provider | Status |
|---|---|
| OpenAI (GPT-4o, GPT-4) | ✅ Supported |
| Anthropic (Claude) | ✅ Supported |
| Groq | ✅ Supported |
| Ollama (local) | 🚧 Coming soon |
| Google Gemini | 🚧 Coming soon |

---

## How it works

```
your terminal
     │
     ▼
  stderr
     │
     ▼
  errpipe         ← strips sensitive paths (optional)
     │
     ▼
 LLM API          ← your provider, your key
     │
     ▼
  stdout          ← explanation + fix, right in your terminal
```

No middleman. No cloud storage. Your error goes to your LLM.

---

## Security

Stack traces can contain secrets. File paths. Env variable names. Token prefixes.

`errpipe` has a scrub mode. Use it.

```sh
errpipe config --scrub on
```

This strips common patterns before the error leaves your machine. It is not perfect. Be careful with what you pipe.

---

## Contributing

PRs are welcome. Keep it simple.

1. Fork the repo
2. Make your change
3. Open a pull request

If you find a bug, open an issue. If you want a feature, open an issue first.

---

## License

MIT. Use it. Break it. Fix it.

---

<div align="center">

Built for engineers who read their error messages.

<!-- PLACE YOUR FOOTER LOGO/EMOJI/BANNER HERE -->
<!-- Example: ![footer](./assets/footer.png) -->

</div>