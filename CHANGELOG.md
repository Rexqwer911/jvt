# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-12

### Added
- Initial release of JVT (Java Version Tool)
- Download and install Java versions from Adoptium/Temurin
- List available Java versions from remote repository
- List installed Java versions
- Switch between installed Java versions (session-based)
- Set default Java version (persistent via Windows registry)
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
