# Changelog
Tracking changes for Soar (using [SemVer 2](http://semver.org/)).

## [0.2.0] - 16-09-2022

### Added
- Replace config `--local` flag with `--global` flag
- Client `subusers` commands
- Client `startup` commands
- Client `--id` flag to `servers:get` command
- Client `settings` commands
- `--page` and `--per-page` flags for application API commands
- Metadata field in data response objects when present
- `config.http.parse_errors` option to parse errors or return as JSON

### Removed
- `config.http.max_body` option (unused)

## [0.1.0-alpha] - 18-08-2022
Initial commit.
