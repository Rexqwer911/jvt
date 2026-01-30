# JVT - Java Version Tool

[![GitHub release](https://img.shields.io/github/v/release/rexqwer911/jvt)](https://github.com/rexqwer911/jvt/releases/latest)
[![Chocolatey](https://img.shields.io/chocolatey/v/jvt)](https://community.chocolatey.org/packages/jvt)
[![License](https://img.shields.io/github/license/rexqwer911/jvt)](LICENSE)

A cross-platform command-line utility that simplifies downloading, installing, and switching between different Java versions - similar to nvm for Node.js.

## Downloads

**Latest Release:**
- **[Download jvt-installer.exe](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-installer.exe)** - Installer executable (Windows)
- **[Download jvt-windows-amd64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-windows-amd64.zip)** - ZIP archive (Windows x64)
- **[Download jvt-macos-amd64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-macos-amd64.zip)** - ZIP archive (macOS x64)
- **[Download jvt-macos-arm64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-macos-arm64.zip)** - ZIP archive (macOS ARM64)
- **[Download jvt-linux-amd64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-linux-amd64.zip)** - ZIP archive (Linux x64) 

Or view [all releases](https://github.com/rexqwer911/jvt/releases)

## Features

- Download and install multiple Java versions
- Easily switch between installed Java versions
- Remove unused Java installations
- List available and installed Java versions
- Support for major Java distributions (Oracle JDK, OpenJDK, Temurin, etc.)
- Cross-platform support (Windows, macOS, Linux)

## Installation

### Windows

#### Option 1: Direct Download (Easiest)
1. **[Download jvt-installer.exe](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-installer.exe)**
2. Run the installer

#### Option 2: Via Chocolatey
```bash
choco install jvt
```

#### Option 3: Manual Installation from ZIP
1. Download [jvt-windows-amd64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-windows-amd64.zip)
2. Extract `jvt.exe` to a directory (e.g., `C:\Program Files\jvt`)
3. Add the directory to your PATH environment variable

### macOS / Linux

1. Download the archive for your OS and architecture (e.g., `jvt-macos-arm64.zip` or `jvt-linux-amd64.zip`).
2. Extract the archive.
3. Move the `jvt` binary to a location in your PATH (e.g., `/usr/local/bin`).
   ```bash
   unzip jvt-macos-arm64.zip
   sudo mv jvt /usr/local/bin/
   chmod +x /usr/local/bin/jvt
   ```
4. **Important:** Add the JVT shell initialization to your shell config (e.g., `~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):
   ```bash
   # Add this line to your shell config file
   eval "$(jvt env)"
   ```
5. Restart your terminal or source your config file.

## Usage

```bash
# List available Java versions for download
jvt list-remote

# Install a specific Java version
jvt install 21

# List installed versions
jvt list

# Switch to a specific version (persists across sessions)
jvt use 21

# Uninstall a version
jvt uninstall 11

# Upgrade Java to the latest version
jvt upgrade 17                 # Upgrade Java 17 to latest
jvt upgrade --all              # Upgrade all installed versions
jvt upgrade --all --dry-run    # Check for updates without installing
jvt upgrade --all --keep-old   # Upgrade but keep old versions

# Show current active version
jvt current
# or
java -version
```

## Development

### Prerequisites
- Go 1.21 or higher

### Building from Source

```bash
# Windows
make build-windows

# macOS (Intel & Apple Silicon)
make build-macos

# Linux
make build-linux
```

## Project Structure

```
jvt/
├── cmd/
│   └── jvt/
│       └── main.go          # Entry point
├── internal/
│   ├── cli/                 # Command-line interface
│   ├── config/              # Configuration management
│   ├── download/            # Download logic
│   ├── install/             # Installation logic
│   ├── registry/            # Java distribution registry
│   └── version/             # Version management
├── installer/               # Installer files
├── chocolatey/              # Chocolatey package files
├── go.mod
├── go.sum
├── README.md
└── CHANGELOG.md
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
