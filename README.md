# AwesomeFS

A simple, lightweight HTTP file server written in Go.

## Features

- 🚀 Simple and lightweight
- 📁 Serves files from a local `files` directory
- 🔒 Secure - prevents path traversal attacks
- ⚙️ Configurable port via command-line flag

## Installation

### Build from source

```bash
git clone https://github.com/yourusername/AwesomeFS.git
cd AwesomeFS
go build -o awesomefs.exe main.go
```

### Prerequisites

- [Go](https://golang.org/) 1.21 or higher

## Usage

1. Place the files you want to share in the `files` directory

2. Run the server:

```bash
# Use default port (8081)
./awesomefs.exe

# Use custom port
./awesomefs.exe -port 3000
```

3. Open your browser and navigate to `http://localhost:8081`

## Command-line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | 8081 | Port to host the server on |
| `-h` | - | Show help |

## Security

AwesomeFS uses Go's built-in `http.FileServer` which automatically prevents path traversal attacks. Users cannot access files outside the designated `files` directory.

## License

MIT License
