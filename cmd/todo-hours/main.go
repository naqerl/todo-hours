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
	flag.BoolVar(&writeFlag, "write", false, "Replace the total-hours line in place with the computed sum")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [PATH]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Sum TODO hours from markdown files with section subtotals.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  PATH    Path to file to parse (default: %s)\n\n", defaultPath)
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nThe tool looks for TODO items matching: - [ ] ... <N>h\n")
		fmt.Fprintf(os.Stderr, "And expects a total line: Total planned hours from TODO items: <N>h\n")
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

	// Parse the file
	result, err := todohours.Parse(text)
	if err != nil {
		if err.Error() == "total line not found" {
			fmt.Fprintf(os.Stderr, "error: expected exactly one total line matching 'Total planned hours from TODO items: <N>h', found 0\n")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "error: parsing file: %v\n", err)
		os.Exit(1)
	}

	// Check for multiple total lines
	totalLine := todohours.ExpectedTotalLine(result.Total)
	countMatches := strings.Count(text, "Total planned hours from TODO items:")
	if countMatches != 1 {
		fmt.Fprintf(os.Stderr, "error: expected exactly one total line matching 'Total planned hours from TODO items: <N>h', found %d\n", countMatches)
		os.Exit(1)
	}

	// Validate or update the total line
	updated := false
	if writeFlag {
		if result.TotalLine != totalLine {
			newText := text[:result.TotalStart] + totalLine + text[result.TotalEnd:]
			err := os.WriteFile(path, []byte(newText), 0o644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: writing file: %v\n", err)
				os.Exit(1)
			}
			updated = true
		}
	} else if result.TotalLine != totalLine {
		fmt.Fprintf(os.Stderr, "error: total line is out of sync; expected '%s' but found '%s'\n", totalLine, result.TotalLine)
		os.Exit(1)
	}

	// Output results
	fmt.Printf("matched_lines=%d\n", result.Count)
	fmt.Printf("total_hours=%d\n", result.Total)
	fmt.Printf("total_line_matches=1\n")
	for section, subtotal := range result.Subtotals {
		fmt.Printf("subtotal[%s]=%d\n", section, subtotal)
	}
	if writeFlag {
		if updated {
			fmt.Println("updated=yes")
		} else {
			fmt.Println("updated=no")
		}
	}
}
