from collections.abc import Sequence
from pathlib import Path
from typing import Iterable
import pytest

BASE_DIR = Path(__file__).resolve().parent


def extract_calibration_values(input_data: Iterable[str]) -> list[int]:
    values = []
    for s in input_data:
        digits = []
        for char in s:
            if char.isdigit():
                digits.append(char)
                break

        for char in reversed(s):
            if char.isdigit():
                digits.append(char)
                break
        if not digits:
            raise ValueError("No digits found on this line")

                # if len(digits) < 2:
                #     digits.append(char)
                # else:
                #     digits[1] = char
        values.append(int(digits[0] + digits[1]))
    return values



test_data = """1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet""".split("\n")


@pytest.mark.parametrize(
    ["input_data", "expected_result"],
    [[test_data, [12, 38, 15, 77]], ]
)
def test_extract_calibration(input_data: list[str], expected_result: list[int]):

    v = extract_calibration_values(input_data)
    assert v == expected_result



if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        calibration_values = extract_calibration_values(f.readlines())
        print(calibration_values)
        print(sum(calibration_values))
