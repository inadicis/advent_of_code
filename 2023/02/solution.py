from collections.abc import Sequence
from pathlib import Path
from typing import Iterable, Mapping
import pytest

BASE_DIR = Path(__file__).resolve().parent


def extract_possible_games(
    game_logs: list[str],
    amount_red: int,
    amount_blue: int,
    amount_green: int,
) -> list[int]:
    valid_games = []
    for i, game_log in enumerate(game_logs):
        if is_game_possible(
            game_log, amount_red=amount_red, amount_green=amount_green, amount_blue=amount_blue
        ):
            valid_games.append(i + 1)
    return valid_games


def is_game_possible(
    game_log: str,
    amount_red: int,
    amount_blue: int,
    amount_green: int,
) -> bool:
    return True


@pytest.mark.parametrize(
    ["game_logs", "amount_red", "amount_blue", "amount_green", "expected_result"],
    [
        [
            [
                "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
                "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
                "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
                "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
                "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
            ],
            12,
            13,
            14,
            [1, 2, 5],
        ],
    ],
)
def test_extract_calibration(
    game_logs: list[str],
    amount_red: int,
    amount_blue: int,
    amount_green: int,
    expected_result: list[int],
):
    v = extract_possible_games(
        game_logs, amount_blue=amount_blue, amount_red=amount_red, amount_green=amount_green
    )
    assert v == expected_result


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        game_ids = extract_possible_games(
            f.readlines(),
            amount_red=12,
            amount_green=13,
            amount_blue=14,
        )
        print(game_ids)
        print(sum(game_ids))
