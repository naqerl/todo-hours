# todo-hours

Sum TODO hours from markdown files with section subtotals.

## Installation

### Using go install

```bash
go install github.com/yourusername/todo-hours/cmd/todo-hours@latest
```

### Using curl (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/todo-hours/main/install.sh | bash
```

This will install the latest binary to `~/.local/bin`.

### Manual Download

Download the latest binary from the [releases page](https://github.com/yourusername/todo-hours/releases).

## Usage

```bash
# Check a markdown file for TODO hour totals
todo-hours path/to/README.md

# Update the total hours line in place
todo-hours path/to/README.md --write
```

The tool looks for TODO items in markdown files matching the pattern:
```markdown
- [ ] Some task description 5h
```

And expects a total line like:
```markdown
Total planned hours from TODO items: 42h
```

## Features

- Sums hours from unchecked TODO items (`- [ ]`)
- Calculates section subtotals based on `##` headers
- Validates or updates total hours line
- Supports `--write` flag to sync totals in place

## Development

```bash
# Build
make build

# Run tests
make test

# Run vet (fmt, vet, staticcheck)
make vet

# Install locally
make install
```

## License

MIT
