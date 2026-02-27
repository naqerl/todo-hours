# todo-hours

Sum TODO hours from markdown files with section subtotals.

## Installation

```bash
uv tool install todo-hours
```

Or with pip:

```bash
pip install todo-hours
```

## Usage

```bash
# Check a markdown file for TODO hour totals
todo-hours path/to/README.md

# Update the total hours line in place
todo-hours path/to/README.md --write
```

The script looks for TODO items in markdown files matching the pattern:
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
# Setup
uv sync

# Run tests
uv run pytest

# Lint
uv run ruff check .
uv run ruff format .

# Type check
uv run mypy src/todo_hours
```

## License

MIT
