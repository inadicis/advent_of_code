import typing
from pathlib import Path

import pytest

BASE_DIR = Path(__file__).resolve().parent


# symbols / digits / points
# - symbols create a "safe area"
# - digits are chained to each other to build a number -> one digit must be part of a safe area
# - points are neutral


# vertical dependency: one char depends on only 3 lines (before, same, after),
# horizontal: because of digits chaining together, potentially the whole line is necessary


# ideas
# - boolean filter eliminating all non-safe zones -> find remaining numbers -> look around them
# in original input
# - go through numbers, find out if there is a symbol adjacent

def extract_safe_numbers(lines: list[str]) -> list[int]:
    safe_numbers = []
    # previous_line = None
    for line_index, line in enumerate(lines):
        # line: str
        # if len(lines) > line_index + 1:
        #     next_line = lines[line_index + 1]
        # else:
        #     next_line = None

        # previous_char = None
        current_number = ""
        is_current_number_safe = False
        for char_index, char in enumerate(line):
            char: str
            # if len(line) > char_index + 1:
            #     next_char = line[char_index + 1]
            # else:
            #     next_char = None

            if char.isdigit():
                current_number += char
                is_current_number_safe = is_current_number_safe or check_adjacent_symbols(
                    lines, line_index, char_index
                )

            else:
                if current_number and is_current_number_safe:
                    safe_numbers.append(int(current_number))
                current_number = ""
                is_current_number_safe = False

            # handle numbers that end the line as well
        if current_number and is_current_number_safe:
            safe_numbers.append(int(current_number))

            # previous_char = char

        # previous_line = line

    # data = cleanup_one_line(line)
    # calculate wanted result
    return safe_numbers


def check_adjacent_symbols(lines: list[str], line_index: int, char_index: int) -> True:
    rows_min = max(0, line_index - 1)
    rows_max = min(len(lines), line_index + 1)

    # assuming the 3 lines have the same length
    col_min = max(0, char_index - 1)
    col_max = min(len(lines[line_index]), char_index + 1)

    for i in range(rows_min, rows_max + 1):
        if len(lines) <= i:
            continue
        for j in range(col_min, col_max + 1):
            # we do check for positions that are waste of processor time (e.g. middle,
            # or adjacent numbers), or adjacent chars of previous numbers
            if len(lines[i]) <= j:
                continue
            adjacent_char = lines[i][j]
            if not adjacent_char.isdigit() and not adjacent_char == ".":
                return True
    return False


#
# def cleanup_one_line(line: str) -> list[int]:
#     """Builds an adapted data structure for the data "hidden" in the line"""
#     return line


@pytest.mark.parametrize(
    ["lines", "expected_result"],
    [
        [
            [
                "467..114..",
                "...*......",
                "..35..633.",
                "......#...",
                "617*......",
                ".....+.58.",
                "..592.....",
                "......755.",
                "...$.*....",
                ".664.598..",
            ],
            [467, 35, 633, 617, 592, 755, 664, 598],
        ]
    ],
)
def test_extract_safe_numbers(
        lines: list[str],
        expected_result: list[int],
):
    v = extract_safe_numbers(
        lines,
    )
    assert v == expected_result


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        results = extract_safe_numbers(
            f.readlines(),
        )
        print(results)
        print(sum(results))
