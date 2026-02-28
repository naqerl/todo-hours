package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/todo-hours/internal/todohours"
)

const defaultPath = "delivery/README.md"

func main() {
	var writeFlag bool

	// Define both short (-w) and long (-write) flags
	flag.BoolVar(&writeFlag, "write", false, "Update an existing but incorrect total line")
	flag.BoolVar(&writeFlag, "w", false, "Update an existing but incorrect total line (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [PATH]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Sum TODO hours from markdown files with section subtotals.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  PATH    Path to file to parse (default: %s)\n\n", defaultPath)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nThe tool looks for TODO items matching: - [ ] ... <N>h\n")
		fmt.Fprintf(os.Stderr, "And expects a total line: Total planned hours from TODO items: <N>h\n\n")
		fmt.Fprintf(os.Stderr, "Behavior:\n")
		fmt.Fprintf(os.Stderr, "  - If total line is missing: it will be auto-created\n")
		fmt.Fprintf(os.Stderr, "  - If total line is correct: validation passes\n")
		fmt.Fprintf(os.Stderr, "  - If total line is wrong: use -write to update it\n")
	}
	flag.Parse()

	// Get file path from arguments or use default
	path := defaultPath
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: file not found: %s\n", path)
		os.Exit(1)
	}

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: reading file: %v\n", err)
		os.Exit(1)
	}
	text := string(content)

	// Calculate hours and subtotals
	count, total := todohours.SumHours(text)
	subtotals := todohours.SectionSubtotals(text)
	expectedLine := todohours.ExpectedTotalLine(total)

	// Find total line position
	start, end, found := todohours.FindTotalLine(text)

	// Check for multiple total lines
	countMatches := 0
	if found {
		countMatches = 1
	}
	if strings.Count(text, "Total planned hours from TODO items:") > 1 {
		fmt.Fprintf(os.Stderr, "error: expected exactly one total line matching 'Total planned hours from TODO items: <N>h', found %d\n",
			strings.Count(text, "Total planned hours from TODO items:"))
		os.Exit(1)
	}

	// Handle write flag
	updated := false
	if writeFlag {
		if found {
			// Update existing line
			currentLine := text[start:end]
			if currentLine != expectedLine {
				newText := text[:start] + expectedLine + text[end:]
				err := os.WriteFile(path, []byte(newText), 0o644)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: writing file: %v\n", err)
					os.Exit(1)
				}
				updated = true
			}
		} else {
			// Add new line at end of file
			newText := text
			if !strings.HasSuffix(text, "\n") {
				newText += "\n"
			}
			newText += "\n" + expectedLine + "\n"
			err := os.WriteFile(path, []byte(newText), 0o644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: writing file: %v\n", err)
				os.Exit(1)
			}
			updated = true
		}
	} else {
		// Without --write flag, still try to add/fix the total line
		if !found {
			// Add new line at end of file
			newText := text
			if !strings.HasSuffix(text, "\n") {
				newText += "\n"
			}
			newText += "\n" + expectedLine + "\n"
			err := os.WriteFile(path, []byte(newText), 0o644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: writing file: %v\n", err)
				os.Exit(1)
			}
			updated = true
		} else {
			// Validate existing line
			currentLine := text[start:end]
			if currentLine != expectedLine {
				fmt.Fprintf(os.Stderr, "error: total line is out of sync; expected '%s' but found '%s'\n", expectedLine, currentLine)
				os.Exit(1)
			}
		}
	}

	// Output results
	fmt.Printf("matched_lines=%d\n", count)
	fmt.Printf("total_hours=%d\n", total)
	fmt.Printf("total_line_matches=%d\n", countMatches)
	for section, subtotal := range subtotals {
		fmt.Printf("subtotal[%s]=%d\n", section, subtotal)
	}
	if updated {
		fmt.Println("updated=yes")
	} else {
		fmt.Println("updated=no")
	}
}
