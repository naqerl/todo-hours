#!/usr/bin/env python3
"""Sum TODO hours, print section subtotals, and optionally sync final total."""

from __future__ import annotations

import argparse
import re
import sys
from collections import OrderedDict
from pathlib import Path

LINE_RE = re.compile(r"^- \[ \]\s.*\d{1,}h$", re.MULTILINE)
HOURS_RE = re.compile(r"(\d{1,})h$")
TOTAL_LINE_RE = re.compile(r"^Total planned hours from TODO items:\s*\d+h$", re.MULTILINE)
H2_RE = re.compile(r"^##\s+(.*\S)\s*$")


def sum_hours(text: str) -> tuple[int, int]:
    total = 0
    count = 0
    for line in LINE_RE.findall(text):
        match = HOURS_RE.search(line)
        if not match:
            continue
        total += int(match.group(1))
        count += 1
    return count, total


def find_total_line_matches(text: str) -> list[re.Match[str]]:
    return list(TOTAL_LINE_RE.finditer(text))


def section_subtotals(text: str) -> OrderedDict[str, int]:
    subtotals: OrderedDict[str, int] = OrderedDict()
    current_section = "Unsectioned"

    for raw_line in text.splitlines():
        section_match = H2_RE.match(raw_line)
        if section_match:
            current_section = section_match.group(1)
            if current_section not in subtotals:
                subtotals[current_section] = 0
            continue

        if LINE_RE.match(raw_line):
            hours_match = HOURS_RE.search(raw_line)
            if not hours_match:
                continue
            if current_section not in subtotals:
                subtotals[current_section] = 0
            subtotals[current_section] += int(hours_match.group(1))

    # Hide sections without TODO-hour lines from output.
    return OrderedDict((k, v) for k, v in subtotals.items() if v > 0)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Sum TODO hours from lines matching: ^- [ ]\\s.*\\d{1,}h$"
    )
    parser.add_argument(
        "path",
        nargs="?",
        default="delivery/README.md",
        help="Path to file to parse (default: delivery/README.md)",
    )
    parser.add_argument(
        "--write",
        action="store_true",
        help="Replace the total-hours line in place with the computed sum",
    )
    args = parser.parse_args()

    path = Path(args.path)
    if not path.exists():
        print(f"error: file not found: {path}", file=sys.stderr)
        return 1

    text = path.read_text(encoding="utf-8")
    count, total = sum_hours(text)
    subtotals = section_subtotals(text)
    matches = find_total_line_matches(text)

    if len(matches) != 1:
        print(
            (
                "error: expected exactly one total line matching "
                "'Total planned hours from TODO items: <N>h', "
                f"found {len(matches)}"
            ),
            file=sys.stderr,
        )
        return 1

    expected_line = f"Total planned hours from TODO items: {total}h"
    current_line = matches[0].group(0)

    updated = False
    if args.write:
        new_text = text[: matches[0].start()] + expected_line + text[matches[0].end() :]
        if new_text != text:
            path.write_text(new_text, encoding="utf-8")
            updated = True
    elif current_line != expected_line:
        print(
            (
                "error: total line is out of sync; "
                f"expected '{expected_line}' but found '{current_line}'"
            ),
            file=sys.stderr,
        )
        return 1

    print(f"matched_lines={count}")
    print(f"total_hours={total}")
    print(f"total_line_matches={len(matches)}")
    for section, subtotal in subtotals.items():
        print(f"subtotal[{section}]={subtotal}")
    if args.write:
        print(f"updated={'yes' if updated else 'no'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
