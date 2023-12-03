import pprint
import typing
from pathlib import Path
import re

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


def extract_gears(lines: list[str]) -> list[int]:
    """Returns a number, its line index and its last character char index"""
    gears = []
    wrapped_lines = wrap_lines_with_points(lines)
    for line_index, line in enumerate(wrapped_lines):
        for char_index, char in enumerate(line):
            if char == "*":
                numbers = find_adjacent_numbers(
                    wrapped_lines, line_index=line_index, char_index=char_index
                )
                if len(numbers) == 2:
                    gears.append(numbers[0] * numbers[1])
    return gears


def find_adjacent_numbers(lines: list[str], line_index: int, char_index: int) -> list[int]:
    # find adjacent digits then looks left and right of these digits to complete the number
    numbers = []
    for i, line in enumerate(lines[line_index - 1 : line_index + 2]):
        skip_amount = 0
        for j, char in enumerate(line[char_index - 1 : char_index + 2]):
            if skip_amount:
                skip_amount -= 1
                continue
            if not char.isdigit():
                continue

            number = char
            if j == 0:  # then we have to look left of it
                for previous_char in reversed(line[: char_index - 1]):
                    if previous_char.isdigit():
                        number = previous_char + number
                    else:
                        break
            for next_char in line[char_index + j :]:  # look right of it
                if next_char.isdigit():
                    number = number + next_char
                    skip_amount += 1
                else:
                    # line ends with a point, so this will happen eventually
                    numbers.append(int(number))
                    break
    return numbers


def wrap_lines_with_points(lines: list[str]):
    # deletes pending newlines as well
    wrapped_lines = ["." * len(lines[0]), *lines, "." * len(lines[0])]
    return ["." + line.strip() + "." for line in wrapped_lines]


@pytest.mark.parametrize(
    ["lines", "expected_result", "description"],
    [
        [
            [
                ".1*2*3*4",
            ],
            [2, 6, 12],
            "one line full of chained multiplications",
        ],
        [
            [
                "467..114.+\n",
                "...*......\n",
                "..35..633.\n",
                "......#...\n",
                "617*......\n",
                ".....+.58.\n",
                "..592.....\n",
                "......755.\n",
                "...$.*....\n",
                ".664.598..\n",
                ".........3\n",
                "3*........\n",
                "3.........\n",
            ],
            [467 * 35, 598 * 755, 3 * 3],
            "given test data",
        ],
        [
            [
                "........",
                ".24*4...",
                "......*.",
            ],
            [24 * 4],
            "one same-line multiplication in a triple line, with one isolated star",
        ],
        [
            [
                "........",
                ".24$-4..",
                "......*.",
            ],
            [],
            "One star adjacent to only one number in triple line data",
        ],
        [
            [
                "11....12",
                "..*..*..",
                "11....12",
            ],
            [11 * 11, 12 * 12],
            "two multi-line multiplications with numbers going further in both directions",
        ],
        [
            [
                ".1....5.",
                "*2*3*4..",
            ],
            [1 * 2, 3 * 4],
            "one multi-line multiplication, a start adjacent to 3 numbers, and a same line multiplication",
        ],
        [
            [
                ".1*2*3*4",
                ".1*2*3*4",
                ".1*2*3*4",
            ],
            [],
            "Only stars adjacent to more than 2 numbers",
        ],
        [
            [
                "...3..",
                ".1*...",
                "...2..",
            ],
            [],
            "One star adjacent to 3 numbers (in 3 lines)",
        ],
        [
            [
                "*...",
                ".113",
                ".113",
                "..*.",
            ],
            [],
            "Stars only adjacent to 1 number",
        ],
    ],
)
def test_extract_gears(lines: list[str], expected_result: list[int], description: str):
    numbers = extract_gears(lines)
    assert numbers == expected_result


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        results = extract_gears(
            f.readlines(),
        )
        print(results)
        print(sum(results))

    # with open(BASE_DIR / "input.txt", "r") as f:
    #     qa = wrap_lines_with_points(f.readlines(), results)
