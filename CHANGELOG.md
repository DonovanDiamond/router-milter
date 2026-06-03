# Changelog

## [v0.3.1] - 2026-06-03

- Added the `sha256` function to the script handler that generates hashes.

## [v0.3.0] - 2026-06-03

- **BREAKING CHANGE:** The milter's actions are now controlled by a JavaScript
script instead of configuration fields. This allows for highly flexible email
filtering by modifying the script as needed, while significantly simplifying
the underlying Go code.

## [v0.2.2] - 2026-06-01

- Support for `reject_to_sha256` in configuration file.
- Support for `-reject-to-sha256` command-line flag.
- Unit tests for SHA256-based recipient rejection.

## [v0.2.1] - 2026-06-01

- Made recipient and sender rejection case-insensitive.
- Made regex-based rejection case-insensitive.

## [v0.2.0] - 2026-06-01

- Support for `reject_to_regex` in configuration file.
- Support for `-reject-to-regex` command-line flag.
- Unit tests for regex-based recipient rejection.
- Configuration loading tests.
