# AwesomeFS

When you just want to serve a file over HTTP on whatever port you specify. Nothing fancy.

## Features

- 🚀 Simple and lightweight
- 📁 Serves files from a local `files` directory
- 🔒 Secure - prevents path traversal attacks
- ⚙️ Configurable port via command-line flag
- 🖥️ **Windows System Tray Support** - Run in background with system tray icon (Windows only)
- ⚡ Auto-start on Windows boot
- 🎯 Quick access to server from system tray

## Installation

### Build from source

```bash
git clone https://github.com/yourusername/AwesomeFS.git
cd AwesomeFS
go build -o bin/awesomefs.exe .
```

### Prerequisites

- [Go](https://golang.org/) 1.21 or higher
- Windows OS (for system tray mode)

## Usage

### CLI Mode (Default)

1. Place the files you want to share in the `files` directory

2. Run the server:

```bash
# Use default port (8081)
./bin/awesomefs.exe

# Use custom port
./bin/awesomefs.exe -port 3000
```

3. Open your browser and navigate to `http://localhost:8081`

### System Tray Mode (Windows Only)

Run AwesomeFS in the background with a system tray icon:

```bash
# Run in system tray mode
./bin/awesomefs.exe -tray

# Run with custom port in tray mode
./bin/awesomefs.exe -tray -port 3000
```

**System Tray Features:**
- Start/Stop server from tray menu
- Open server in browser with one click
- Configure settings (host and port)
- Enable/disable auto-start on Windows boot
- Status indicator showing server state

## Command-line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-tray` | false | Run in system tray mode (Windows only) |
| `-host` | 0.0.0.0 | Interface to bind to |
| `-port` | 8081 | Port to host the server on |
| `-h` | - | Show help |

## Security

AwesomeFS uses Go's built-in `http.FileServer` which automatically prevents path traversal attacks. Users cannot access files outside the designated `files` directory.

## License

MIT License
