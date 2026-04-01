# 🐸 gh-quoi 🐸

[![Go Report Card](https://goreportcard.com/badge/github.com/BenbenIO/gh-quoi)](https://goreportcard.com/report/github.com/BenbenIO/gh-quoi)
[![License](https://img.shields.io/github/license/BenbenIO/gh-quoi)](LICENSE)
[![Release](https://img.shields.io/github/v/release/BenbenIO/gh-quoi)](https://github.com/BenbenIO/gh-quoi/releases/latest)
[![CI](https://github.com/BenbenIO/gh-quoi/actions/workflows/ci.yml/badge.svg)](https://github.com/BenbenIO/gh-quoi/actions/workflows/ci.yml)

Filtered GitHub notifications in your terminal.

“Quoi?” (kwa 🐸🎵) — that moment when your GitHub notifications explode into chaos and
you are asking yourself: "What should I focus on right now?"

> `Quoi` means `What` in French.

`gh-quoi` helps you answer that question by giving you a TUI for managing notifications,
with custom filtering to show only what you actually care about.

## Installation

### Requirements

- Go 1.24 or later (if building from source)
- A GitHub Personal Access Token with `notifications` scope
- :turtle: On `wls`, please install [wlsu](https://github.com/WhitewaterFoundry/wslu-documentation/blob/master/install.md) to open the browser from WSL:

  ```shell
  sudo apt install wslu
  ```

### Quick Install (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/BenbenIO/gh-quoi/main/install.sh | sh
```

Or download and run the script manually:

```bash
wget https://raw.githubusercontent.com/BenbenIO/gh-quoi/main/install.sh
chmod +x install.sh
./install.sh
```

The script will:

- Detect your OS and architecture
- Download the latest release
- Install to `/usr/local/bin` (customizable via `INSTALL_DIR` env var)

### Binary Releases

Download pre-built binaries for your OS from the [releases page](https://github.com/BenbenIO/gh-quoi/releases).

### Via Go Install

```bash
go install github.com/BenbenIO/gh-quoi/cmd/gh-quoi@latest
```

### From Source

```bash
git clone https://github.com/BenbenIO/gh-quoi.git
cd gh-quoi
go build -o gh-quoi ./cmd/gh-quoi
```

### Get a GitHub Personal Access Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Generate a new token (classic)
3. Select the `notifications` scope
4. Copy the generated token

> **Note:** Only the `notifications` scope is required for gh-quoi to function.

## Usage

### Quick Start

1. Initialize configuration:

   ```bash
   # Create a default config file in the expected location to start
   gh-quoi init
   ```

2. Export your GitHub Personal Access Token as an environment variable (or source a `.env` file):

   ```bash
   export GH_TOKEN=your_personal_access_token
   ```

3. Start the TUI:

   ```bash
   gh-quoi
   ```

### Commands

```bash
gh-quoi                # Launch the TUI
gh-quoi init           # Initialize configuration
gh-quoi version        # Show version information
gh-quoi help           # Show help
```

### Keyboard Shortcuts

| Key     | Action                       |
| ------- | ---------------------------- |
| `↑/k`   | Move up                      |
| `↓/j`   | Move down                    |
| `Enter` | Open notification in browser |
| `r`     | Refresh notifications        |
| `f`     | Toggle filter menu           |
| `?`     | Show help                    |
| `q/Esc` | Quit                         |

### Create Filters

Filters are defined in your configuration file (`~/.config/gh-quoi/config.yaml`):

```yaml
---
since_day_ago: 7

filters:
  - name: "Pull Requests"
    subject_type: "PullRequest"

  - name: "Mentions"
    reason: "mention"

  - name: "My Repositories"
    repository: "myorg/myrepo"
```

Available filter fields:

- `name`: Display name for the filter
- `subject_type`: Filter by notification type (`Issue`, `PullRequest`, `Release`, etc.)
- `reason`: Filter by reason (`mention`, `assign`, `review_requested`, `subscribed`, etc.)
- `repository`: Filter by repository (format: `owner/repo`)
- `unread`: Filter by read status (`true` or `false`)

### Filter logics

#### Rule Structure

> Each rule = AND logic

Inside a rule, all fields must match to be included.

Example:

```yaml
filters:
  - name: "PR Mentions from myorg"
    types: ["PullRequest"]
    reasons: ["mention"]
    repos: ["myorg/*"]
```

Means: PR AND mention AND from org myorg.

#### Combining Rules

> Multiple rules = OR logic

If a notification matches any rule, it is included.

This is extremely flexible:

```yaml
filters:
  - name: "PRs"
    types: ["PullRequest"]

  - name: "Mentions"
    reasons: ["mention"]
```

Means: (PR) OR (mention) are included.

#### Supported fields

Field Meaning:

- `types`: Issue / PullRequest / Release / Commit / etc.
- `reasons`: mention / assigned / review_requested / …
- `repos`: glob patterns over repo names
- `title_regex`: filter by regex on PR/issue title

- <https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28>

#### GitHub API limitations

GitHub only allows server-side filtering by:

- since=
- all=true (read)
- participating=true
- Everything else — repo filtering, type filtering, reason filtering — is done client-side by `gh-quoi`.

> GitHub UI does allows filtering with `author` which is practical but not supported by the API and
> not returned in notification objects so it cannot (easily) be implemented in `gh-quoi` 😟.

## Dev

### Feature ideas

- Mark WebUI notification as `Read`: to clean up the Web UI notifications board (`until` config).
- Configurable refresh interval.
- In the status bar show number of filtered vs total notifications.
- Auto sizing for columns.
- Add sorting by reason, for example: `mentioned` first
- Integration to `vim`.

### Building

```bash
go build -o gh-quoi ./cmd/gh-quoi
```

### Testing

```bash
go test ./...
```

### Linting and formatting

```bash
pre-commit install --config .pre-commit-config.yaml
pre-commit run --all-files
```

### Releasing

This project uses [GoReleaser](https://goreleaser.com/) for automated releases via GitHub Actions.

#### Creating a New Release

1. **Ensure all changes are merged to main**

   ```bash
   git checkout main
   git pull origin main
   ```

2. **Create and push a version tag** (following [semantic versioning](https://semver.org/))

   ```bash
   git tag v0.0.0
   git push origin v0.0.0
   ```

3. **The release will be available at:**
   - Release page: `https://github.com/BenbenIO/gh-quoi/releases/latest`
   - Users can install via: `go install github.com/BenbenIO/gh-quoi/cmd/gh-quoi@v1.2.0`

#### Local Release Testing

To test the release process locally (requires GoReleaser installed):

```bash
go install github.com/goreleaser/goreleaser/v2@latest
goreleaser release --snapshot --clean
ls -la dist/
```

### Resources

- <https://docs.github.com/en/rest/activity/notifications?apiVersion=2022-11-28>
- <https://pkg.go.dev/github.com/google/go-github/github>
- <https://github.com/charmbracelet/bubbletea>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
