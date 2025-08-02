

# docker-wizard

![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=flat-square)

> 🚀 **docker-wizard** is a Go CLI tool that detects Python, Node.js, and Go projects in the current directory and generates a Dockerfile with the correct base image, build/run instructions, and multistage formatting. It also supports React, Angular, and Vue frontend frameworks.

---

## ✨ Features

- **Automatic project detection**: Scans for `requirements.txt`, `package.json`, or `go.mod`.
- **Dynamic base image selection**: Uses the correct image version for your project.
- **Multistage Dockerfile generation**: Optimized for production.
- **Frontend support**: Adds build steps for React, Angular, and Vue projects.
- **Colorful CLI output**: Easy to read and visually appealing.

---

## ⚡ Usage

```bash
# Download dependencies
make deps

# Build the CLI
make build

# Install the CLI (system-wide)
make install

# Run the CLI
docker-wizard --output Dockerfile
```

---

## 📁 Project Structure

- `cmd/docker-wizard/main.go` — CLI entry point
- `internal/generate/generate.go` — Dockerfile generation logic
- `internal/generate/generate_test.go` — Unit tests
- `.github/copilot-instructions.md` — Copilot custom instructions
- `.gitignore` — Ignore build artifacts and IDE files
- `Makefile` — Build and install commands

---

## 📝 License

MIT
