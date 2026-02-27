// Package todohours provides functionality for summing TODO hours from markdown files.
package todohours

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	// lineRe matches TODO lines with hours: "- [ ] Some task 5h"
	lineRe = regexp.MustCompile(`^- \[ \]\s.*\d+h$`)
	// hoursRe extracts hours from end of line: "5h"
	hoursRe = regexp.MustCompile(`(\d+)h$`)
	// totalLineRe matches the total hours line
	totalLineRe = regexp.MustCompile(`(?m)^Total planned hours from TODO items:\s*\d+h$`)
	// h2Re matches H2 markdown headers
	h2Re = regexp.MustCompile(`^##\s+(.*\S)\s*$`)
)

// Result holds the parsed TODO hours data.
type Result struct {
	Count      int
	Total      int
	Subtotals  map[string]int
	TotalLine  string
	TotalStart int
	TotalEnd   int
}

// SumHours calculates the total hours from TODO items in the text.
// Returns the count of matched lines and total hours.
func SumHours(text string) (count, total int) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if lineRe.MatchString(line) {
			if match := hoursRe.FindStringSubmatch(line); match != nil {
				var hours int
				fmt.Sscanf(match[1], "%d", &hours)
				total += hours
				count++
			}
		}
	}
	return count, total
}

// FindTotalLine finds the total line in the text and returns its position.
func FindTotalLine(text string) (start, end int, found bool) {
	loc := totalLineRe.FindStringIndex(text)
	if loc == nil {
		return 0, 0, false
	}
	return loc[0], loc[1], true
}

// SectionSubtotals calculates subtotals for each section (H2 header) in the text.
func SectionSubtotals(text string) map[string]int {
	subtotals := make(map[string]int)
	currentSection := "Unsectioned"
	sectionHasHours := make(map[string]bool)

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()

		if match := h2Re.FindStringSubmatch(line); match != nil {
			currentSection = match[1]
			continue
		}

		if lineRe.MatchString(line) {
			if match := hoursRe.FindStringSubmatch(line); match != nil {
				var hours int
				fmt.Sscanf(match[1], "%d", &hours)
				subtotals[currentSection] += hours
				sectionHasHours[currentSection] = true
			}
		}
	}

	// Filter out sections without TODO-hour lines
	result := make(map[string]int)
	for section, total := range subtotals {
		if sectionHasHours[section] {
			result[section] = total
		}
	}

	return result
}

// Parse analyzes the markdown text and returns the parsed results.
func Parse(text string) (*Result, error) {
	count, total := SumHours(text)
	subtotals := SectionSubtotals(text)
	start, end, found := FindTotalLine(text)

	if !found {
		return nil, fmt.Errorf("total line not found")
	}

	return &Result{
		Count:      count,
		Total:      total,
		Subtotals:  subtotals,
		TotalLine:  text[start:end],
		TotalStart: start,
		TotalEnd:   end,
	}, nil
}

// ExpectedTotalLine returns the expected total line format.
func ExpectedTotalLine(total int) string {
	return fmt.Sprintf("Total planned hours from TODO items: %dh", total)
}
