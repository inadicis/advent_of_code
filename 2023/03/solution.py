import typing
from pathlib import Path

import pytest

BASE_DIR = Path(__file__).resolve().parent


def main(
    lines: list[str],
) -> list[int]:
    results = []
    for i, line in enumerate(lines):
        result = process_one_line(line)
        results.append(result)

    return results


def process_one_line(line: str) -> bool:
    data = cleanup_one_line(line)
    # calculate wanted result
    return True


MyDatastructure = typing.TypeVar("MyDatastructure")


def cleanup_one_line(line: str) -> MyDatastructure:
    """Builds an adapted data structure for the data "hidden" in the line"""
    return line


@pytest.mark.parametrize(
    ["line", "expected_result"],
    [
        [
            "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
            {"red": 4, "green": 2, "blue": 6},
        ]
    ],
)
def test_extract_minimum_colors(
    line: str,
    expected_result: MyDatastructure,
):
    v = process_one_line(
        line,
    )
    assert v == expected_result


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        results = main(
            f.readlines(),
        )
        print(results)
        print(sum(results))
