# JVT - Java Version Tool

[![GitHub release](https://img.shields.io/github/v/release/rexqwer911/jvt)](https://github.com/rexqwer911/jvt/releases/latest)
[![Chocolatey](https://img.shields.io/chocolatey/v/jvt)](https://community.chocolatey.org/packages/jvt)
[![License](https://img.shields.io/github/license/rexqwer911/jvt)](LICENSE)

A command-line utility for Windows that simplifies downloading, installing, and switching between different Java versions - similar to nvm for Node.js.

## Downloads

**Latest Release:**
- **[Download jvt-installer.exe](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-installer.exe)** - Installer executable
- **[Download jvt-windows-amd64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-windows-amd64.zip)** - ZIP archive

Or view [all releases](https://github.com/rexqwer911/jvt/releases)

## Features

- Download and install multiple Java versions
- Easily switch between installed Java versions
- Remove unused Java installations
- List available and installed Java versions
- Support for major Java distributions (Oracle JDK, OpenJDK, Temurin, etc.)

## Installation

### Option 1: Direct Download (Easiest)
1. **[Download jvt-installer.exe](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-installer.exe)**
2. Run the installer

### Option 2: Via Chocolatey
```bash
choco install jvt
```

### Option 3: Manual Installation from ZIP
1. Download [jvt-windows-amd64.zip](https://github.com/rexqwer911/jvt/releases/latest/download/jvt-windows-amd64.zip)
2. Extract `jvt.exe` to a directory (e.g., `C:\Program Files\jvt`)
3. Add the directory to your PATH environment variable

## Usage

```bash
# List available Java versions
jvt list-remote

# Install a specific Java version
jvt install 21

# List installed versions
jvt list

# Use a specific version
jvt use 21

# Uninstall a version
jvt uninstall 11

# Show current version
jvt current
# or
java -version
```

## Development

### Prerequisites
- Go 1.21 or higher
- Windows 10/11

### Building from Source
```bash
go build -o jvt.exe cmd/jvt/main.go
```
or
```bash
make build
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
