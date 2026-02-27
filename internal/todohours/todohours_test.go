package todohours

import (
	"testing"
)

func TestSumHoursBasic(t *testing.T) {
	text := `
- [ ] Task one 5h
- [ ] Task two 3h
- [x] Done task 2h
`
	count, total := SumHours(text)
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
	if total != 8 {
		t.Errorf("expected total 8, got %d", total)
	}
}

func TestSumHoursNoMatches(t *testing.T) {
	text := "No todo items here"
	count, total := SumHours(text)
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}

func TestSectionSubtotals(t *testing.T) {
	text := `
## Section A
- [ ] Task one 5h
- [ ] Task two 3h

## Section B
- [ ] Task three 2h
`
	subtotals := SectionSubtotals(text)

	if subtotals["Section A"] != 8 {
		t.Errorf("expected Section A subtotal 8, got %d", subtotals["Section A"])
	}
	if subtotals["Section B"] != 2 {
		t.Errorf("expected Section B subtotal 2, got %d", subtotals["Section B"])
	}
}

func TestFindTotalLine(t *testing.T) {
	text := "Total planned hours from TODO items: 10h\n"
	start, end, found := FindTotalLine(text)
	if !found {
		t.Error("expected to find total line")
	}
	if start == 0 && end == 0 {
		t.Error("expected non-zero positions")
	}
}

func TestFindTotalLineNotFound(t *testing.T) {
	text := "Some random text without total line"
	_, _, found := FindTotalLine(text)
	if found {
		t.Error("expected not to find total line")
	}
}

func TestExpectedTotalLine(t *testing.T) {
	expected := "Total planned hours from TODO items: 42h"
	result := ExpectedTotalLine(42)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestParse(t *testing.T) {
	text := `## Section A
- [ ] Task one 5h
- [ ] Task two 3h

Total planned hours from TODO items: 8h
`
	result, err := Parse(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}
	if result.Total != 8 {
		t.Errorf("expected total 8, got %d", result.Total)
	}
	if result.Subtotals["Section A"] != 8 {
		t.Errorf("expected Section A subtotal 8, got %d", result.Subtotals["Section A"])
	}
}

func TestParseNoTotalLine(t *testing.T) {
	text := `## Section A
- [ ] Task one 5h
`
	_, err := Parse(text)
	if err == nil {
		t.Error("expected error for missing total line")
	}
}
