import pprint
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


def extract_safe_numbers(lines: list[str]) -> list[tuple[int, int, int]]:
    """Returns a number, its line index and its last character char index"""
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
                    safe_numbers.append((int(current_number), line_index, char_index))
                current_number = ""
                is_current_number_safe = False

            # handle numbers that end the line as well
        if current_number and is_current_number_safe:
            safe_numbers.append((int(current_number), line_index, len(line) - 1))

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
            if not adjacent_char.isdigit() and not adjacent_char in [".", "\n"]:
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
            [467, 35, 633, 617, 592, 755, 664, 598, 3, 3],
        ],
        [
            [
                "........",
                ".24..4..",
                "......*.",
            ],
            [4],
        ],
        [
            [
                "........",
                ".24$-4..",
                "......*.",
            ],
            [24, 4],
        ],
        [
            [
                "11....11",
                "..$..$..",
                "11....11",
            ],
            [11, 11, 11, 11],
        ],
        [
            [
                "$......$",
                ".1....1.",
                ".1....1.",
                "$......$",
            ],
            [1, 1, 1, 1],
        ],
        [
            [
                "$......$",
                ".11..11.",
                ".11..11.",
                "$......$",
            ],
            [11, 11, 11, 11],
        ],
        [
            [
                "$11",
                "...",
                "11$",
                "...",
            ],
            [11, 11],
        ],
        [
            [
                "$..",
                ".11",
                ".11",
                "$..",
                "..$",
                "11.",
                "11.",
                "..$",
            ],
            [11, 11, 11, 11],
        ],
        [
            [
                "11.$.",
            ],
            [],
        ],
    ],
)
def test_extract_safe_numbers(
        lines: list[str],
        expected_result: list[int],
):
    numbers = [number for number, line, char in extract_safe_numbers(lines)]
    assert numbers == expected_result


def qa_results(lines: list[str], safe_numbers: list[int, int, int]):
    # # add points around to avoid out of bounds when looking around digits
    # # pprint.pp(lines)
    # wrapped_lines = ["." + line.strip() + "." for line in lines]
    # line_length = len(wrapped_lines[0])
    # text = "." * line_length + "\n" + "\n".join(wrapped_lines) + "\n" + "." * line_length
    # # pprint.pp(wrapped_lines)
    # # pprint.pp(text)
    # current_index = 0
    # for number_id, (number, line_index, char_index) in enumerate(safe_numbers):
    #     sub = str(number)
    #     # i = text.find(sub, current_index, len(text))
    #     i = line_length * (line_index + 1) + char_index + 1  # taking into account wrapped points
    #     min_col = i - 1
    #     max_col = i + len(sub) + 1  # excluded
    #     # print(f"{number_id=}, {number=}, {i=}, {min_col=}, {max_col=}")
    #     offset = line_length + 1  # we have to count newline char now
    #     chars_line_before = text[min_col - offset: max_col - offset]
    #     chars_same_line = text[min_col:max_col]
    #     chars_line_after = text[min_col + offset: max_col + offset]
    #     chars = {*chars_line_before, *chars_same_line, *chars_line_after}
    #     symbols = {char for char in chars if char != "." and not char.isdigit()}
    #     print(line_index, char_index)
    #     print(f"{number=}, {i=},  {number_id=}, {chars=}")
    #     print(chars_line_before)
    #     print(chars_same_line)
    #     print(chars_line_after)
    #     if not symbols:
    #         break
    #     current_index = i
    #     # if number_id > 500:
    #     #     break
    wrapped_lines = ["." * len(lines[0]), *lines, "." * len(lines[0])]
    wrapped_lines = ["." + line + "." for line in wrapped_lines]
    for number, line_index, last_char_index in safe_numbers:
        line_before = wrapped_lines[line_index][
                      last_char_index - len(str(number)): last_char_index + 2
                      ]
        line_of_number = wrapped_lines[line_index + 1][
                         last_char_index - len(str(number)): last_char_index + 2
                         ]
        line_after = wrapped_lines[line_index + 2][
                     last_char_index - len(str(number)): last_char_index + 2
                     ]
        # print(line_before)
        # print(line_of_number)
        # print(line_after)
        # if line_index > 1:
        #     return
        chars = {*line_before, *line_after, *line_of_number}
        symbols = {char for char in chars if char != "." and not char.isdigit()}
        if not symbols:
            print(line_before)
            print(line_of_number)
            print(line_after)
            break


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        results = extract_safe_numbers(
            f.readlines(),
        )
        print(results)
        print(sorted(results))
        print(sum(numbers for numbers, line, last_char in results))

    with open(BASE_DIR / "input.txt", "r") as f:
        qa = qa_results(f.readlines(), results)
