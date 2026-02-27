"""Tests for todo-hours."""

from todo_hours.main import sum_hours, section_subtotals, find_total_line_matches


def test_sum_hours_basic():
    text = """
- [ ] Task one 5h
- [ ] Task two 3h
- [x] Done task 2h
"""
    count, total = sum_hours(text)
    assert count == 2
    assert total == 8


def test_sum_hours_no_matches():
    text = "No todo items here"
    count, total = sum_hours(text)
    assert count == 0
    assert total == 0


def test_section_subtotals():
    text = """
## Section A
- [ ] Task one 5h
- [ ] Task two 3h

## Section B
- [ ] Task three 2h
"""
    subtotals = section_subtotals(text)
    assert subtotals["Section A"] == 8
    assert subtotals["Section B"] == 2


def test_find_total_line():
    text = "Total planned hours from TODO items: 10h\n"
    matches = find_total_line_matches(text)
    assert len(matches) == 1
