# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2025-12-16

### Added
- Cross-platform support (macOS and Linux).
- Shell integration (`jvt env`) for simplified PATH management on Unix-like systems.
- README update

## [1.1.0] - 2025-12-13

### Changed
- Updated version to 1.1.0
- Added automated installer build via GitHub Actions
- Changed: Merged jvt default functionality into jvt use.
- Changed: jvt list now highlights the currently active version.
- Removed: jvt completion command (no longer needed).
- Fixed: Prevented duplicate PATH environment variable entries.
- Fixed: Installer no longer creates a desktop shortcut or prompts to launch on exit.

## [1.0.0] - 2025-12-12

### Added
- Initial release of JVT (Java Version Tool)
- Download and install Java versions from Adoptium/Temurin
- List available Java versions from remote repository
- List installed Java versions
- Switch between installed Java versions (persistent via Windows registry)
- Uninstall Java versions
- Check current active Java version
- Automatic JAVA_HOME and PATH management
- SHA256 checksum verification for downloads
- Detection of system-level Java installations

### Features
- Windows x64 support
- User-friendly CLI with Cobra framework
- Installation to `~/.jvt/versions/`
- Download caching to `~/.jvt/cache/`
- Chocolatey package for easy installation

[1.0.0]: https://github.com/rexqwer911/jvt/releases/tag/v1.0.0
