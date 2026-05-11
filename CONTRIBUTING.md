# Contributing to errpipe

First off, thank you for considering contributing to `errpipe`! It's people like you that make `errpipe` a great tool for the developer community.

## How Can I Contribute?

### Reporting Bugs
* **Check the issue tracker** to see if the bug has already been reported.
* If you find a new bug, **open a new issue**. Be sure to include:
    * Your operating system.
    * The shell you are using (bash, zsh, cmd, powershell).
    * Steps to reproduce the bug.
    * The expected behavior vs. what actually happened.

### Suggesting Enhancements
* Open an issue with the tag `enhancement`.
* Explain why the feature would be useful and how you imagine it working.

### Pull Requests
1. **Fork the repository** and create your branch from `main`.
2. **Setup your environment** (see below).
3. **Make your changes**. If you're adding a new AI provider, follow the existing pattern in `internal/ai/`.
4. **Follow Go idioms**. Run `go fmt ./...` before committing.
5. **Update documentation** if you change user-facing behavior.
6. **Submit a Pull Request** with a clear description of your changes.

---

## Development Setup

### Prerequisites
- **Go 1.21 or higher**.
- Access to AI API keys (optional, but needed for testing "Inline CLI Mode").
- System-specific dependencies (e.g., `libx11` headers on Linux for `robotgo` / `rod`).

### Building from source
```bash
git clone https://github.com/DiamondOsasx/errpipe.git
cd errpipe
go mod download
go build -o errpipe
```

### Running in Development
You can run the project directly using:
```bash
go run main.go ai.go
```

---

## Project Structure

- `main.go`: Entry point and REPL loop.
- `ai.go`: Orchestrates sending errors to the correct AI provider/mode.
- `internal/ai/`: Contains provider implementations.
    - `gemini/`, `claude/`, `chatgpt/`: Each contains logic for `Web`, `CLI`, and `Inline` modes.
- `internal/cli/`: Handles configuration (`--init`) and state.
- `internal/utils/`:
    - `browser.go`: Cross-platform browser automation using `rod`.
    - `window/`: OS-specific window focus logic.

---

## Coding Standards

- **Go Standards**: We follow standard Go formatting and idioms. Always run `go fmt`.
- **Commit Messages**:
    - Use the present tense ("Add feature" not "Added feature").
    - Use the imperative mood ("Move cursor to..." not "Moves cursor to...").
    - Keep the first line under 50 characters if possible.

## License
By contributing, you agree that your contributions will be licensed under its [MIT License](./LICENSE).
