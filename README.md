# JVT - Java Version Tool

A command-line utility for Windows that simplifies downloading, installing, and switching between different Java versions - similar to nvm for Node.js.

## Features

- Download and install multiple Java versions
- Easily switch between installed Java versions
- Remove unused Java installations
- List available and installed Java versions
- Set default Java version
- Support for major Java distributions (Oracle JDK, OpenJDK, Temurin, etc.)

## Installation

### Via Chocolatey (Recommended)
```bash
choco install jvt
```

### Manual Installation
1. Download the latest release from [Releases](https://github.com/rexqwer911/jvt/releases)
2. Extract to a directory (e.g., `C:\Program Files\jvt`)
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

# Set default version
jvt default 21

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
