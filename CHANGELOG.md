# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [1.0.0] - 2026-02-14

### Added
- Initial stable release of `gyml`
- **SetValue** function: Set values at any path within a YAML document with automatic creation of intermediate nodes
- **GetValue** function: Retrieve values from YAML documents using generic type parameters
- **DeleteValue** function: Remove nodes or entire sub-trees from YAML documents
- Path-based access using simple key notation (e.g., `["database", "connections", "[0]", "port"]`)
- Support for dynamic type handling (maps, sequences, scalars)
- Comprehensive error handling for invalid paths and operations
- Full test coverage for all public APIs
- Complete documentation with usage examples

[1.0.0]: https://github.com/matus-u/gyml/releases/tag/v1.0.0
