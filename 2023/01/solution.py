from pathlib import Path
from typing import Iterable, Mapping

import pytest

BASE_DIR = Path(__file__).resolve().parent

DIGITS = ["one", "two", "three", "four", "five", "six", "seven", "eight", "nine"]
REVERSED_DIGITS = [s[::-1] for s in DIGITS]

DIGITS_MAPPING = {str(i + 1): str(i + 1) for i in range(9)}
BASE_MAPPING: dict[str, str] = {
    **DIGITS_MAPPING,
    **{digit: str(i + 1) for i, digit in enumerate(DIGITS)},
}
BASE_MAPPING_REVERSE: dict[str, str] = {
    **DIGITS_MAPPING,
    **{digit: str(i + 1) for i, digit in enumerate(REVERSED_DIGITS)},
}


def extract_calibration_values(input_data: Iterable[str]) -> list[int]:
    values = []
    for s in input_data:
        digits = []

        digits.append(
            extract_first_digit(s, base_choices=BASE_MAPPING, current_choices=BASE_MAPPING)
        )
        digits.append(
            extract_first_digit(
                s[::-1], base_choices=BASE_MAPPING_REVERSE, current_choices=BASE_MAPPING_REVERSE
            )
        )
        if not digits:
            raise ValueError("No digits found on this line")

        values.append(int(f"{digits[0]}{digits[1]}"))
    return values


def extract_first_digit(
    s: str, base_choices: Mapping[str, str], current_choices: Mapping[str, str]
) -> str:
    if s[0] in current_choices:
        return current_choices[s[0]]
    new_choices = {
        **base_choices,
        **{key[1:]: value for key, value in base_choices.items() if key.startswith(s[0])},
        **{key[1:]: value for key, value in current_choices.items() if key.startswith(s[0])},
    }
    return extract_first_digit(s[1:], current_choices=new_choices, base_choices=base_choices)


@pytest.mark.parametrize(
    ["input_data", "expected_result"],
    [
        [["1abc2", "pqr3stu8vwx", "a1b2c3d4e5f", "treb7uchet"], [12, 38, 15, 77]],
        [
            [
                "two1nine",
                "eightwothree",
                "abcone2threexyz",
                "xtwone3four",
                "4nineeightseven2",
                "zoneight234",
                "7pqrstsixteen",
            ],
            [29, 83, 13, 24, 42, 14, 76],
        ],
    ],
)
def test_extract_calibration(input_data: list[str], expected_result: list[int]):
    v = extract_calibration_values(input_data)
    assert v == expected_result


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        calibration_values = extract_calibration_values(f.readlines())
        print(calibration_values)
        print(sum(calibration_values))
