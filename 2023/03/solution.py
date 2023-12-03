import typing
from pathlib import Path

import pytest

BASE_DIR = Path(__file__).resolve().parent


# symbols / digits / points
# - symbols create a "safe area"
# - digits are chained to each other to build a number -> one digit must be part of a safe area
# - points are neutral


def extract_safe_numbers(lines: list[str]) -> list[int]:
    # data = cleanup_one_line(line)
    # calculate wanted result
    return []


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
