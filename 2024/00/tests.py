import pytest
from solution import main


@pytest.mark.parametrize(
    ["description", "input", "expected"],
    [
        [
            "Minimal request",
            "test_input.txt",
            True,
        ],
    ],
)
def test_cases(input: str, description: str, expected: bool):
    assert main(input) == expected, description
