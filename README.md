# gh-dependabot

A [GitHub CLI](https://cli.github.com/) extension to interact with pull requests opened by Dependabot. It allows you to approve, merge, and rebase Dependabot PRs in bulk.

## Requirements

- [Go](https://go.dev/) **1.25+**
- [GitHub CLI](https://cli.github.com/) (`gh`) **2.0+**, authenticated (`gh auth login`)

## Installation

```sh
gh extension install clementvtrd/gh-dependabot
```

After installation the extension is available as a `gh` subcommand:

```sh
gh dependabot --help
```

## Usage

| Command | Description |
|---------|-------------|
| `gh dependabot approve [--dry-run]` | Approve all open Dependabot PRs with passing CI checks |
| `gh dependabot merge [--dry-run] [--rebase]` | Merge all open Dependabot PRs with passing CI checks |
| `gh dependabot rebase [--dry-run]` | Ask Dependabot to rebase all its open PRs |

### Flags

- `--dry-run` — Print what would be done without executing any action.
- `--rebase` _(merge only)_ — Use the rebase merge strategy instead of a merge commit.

## Development workflow

### Clone and build

```sh
git clone https://github.com/clementvtrd/gh-dependabot.git
cd gh-dependabot
go build
```

### Run locally

```sh
# Run the built binary directly
./gh-dependabot approve --dry-run

# Or install it as a local extension for integration testing
gh extension install .
gh dependabot approve --dry-run
```

### Run tests

```sh
go test ./...
```

### Lint (optional)

If you have [`golangci-lint`](https://golangci-lint.run/) installed:

```sh
golangci-lint run
```

## Publish workflow

1. Ensure all tests pass:

   ```sh
   go test ./...
   ```

2. Tag a new release following [semantic versioning](https://semver.org/):

   ```sh
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. Create a GitHub release for the tag (manually or via `gh`):

   ```sh
   gh release create v1.0.0 --generate-notes
   ```

   Optionally attach pre-built binaries for different platforms:

   ```sh
   GOOS=linux   GOARCH=amd64 go build -o dist/linux-amd64/gh-dependabot
   GOOS=darwin  GOARCH=arm64 go build -o dist/darwin-arm64/gh-dependabot
   GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/gh-dependabot.exe

   gh release upload v1.0.0 dist/**/*
   ```

4. Users can then install or upgrade with:

   ```sh
   gh extension install clementvtrd/gh-dependabot
   gh extension upgrade dependabot
   ```
