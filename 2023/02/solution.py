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
        min_colors = extract_min_amounts(game_log)
        set_power = min_colors[0] * min_colors[1] * min_colors[2]

    return valid_games


def extract_minimum_set_powers(game_logs: list[str]) -> list[int]:
    set_powers = []
    for i, game_log in enumerate(game_logs):
        min_colors = extract_min_amounts(game_log)
        set_powers.append(min_colors[0] * min_colors[1] * min_colors[2])

    return set_powers


def is_game_possible(
    game_log: str,
    amount_red: int,
    amount_blue: int,
    amount_green: int,
) -> bool:
    game_id, pulls = game_log.split(":")
    pulls = pulls.strip().split(";")
    color_max = {"blue": amount_blue, "green": amount_green, "red": amount_red}
    for pull in pulls:
        color_pulls = pull.strip().split(",")
        for color_pull in color_pulls:
            amount, color = color_pull.strip().split(" ")
            if int(amount.strip()) > color_max[color]:  # raises if color not in mapping
                return False
            # if color == "blue" and int(amount) > amount_blue:
            #     return False
            # elif color == "red" and int(amount) > amount_red:
            #     return False
            # elif color == "green" and int(amount) > amount_green:
            #     return False
            # else:
            #     raise ValueError(f"Non recognized color {color}")
    return True


def extract_min_amounts(
    game_log: str,
) -> tuple[int, int, int]:  # r,g,b
    game_id, pulls = game_log.split(":")
    pulls = pulls.strip().split(";")
    color_maxs = {"blue": 0, "green": 0, "red": 0}
    for pull in pulls:
        color_pulls = pull.strip().split(",")
        for color_pull in color_pulls:
            amount, color = color_pull.strip().split(" ")
            amount = int(amount.strip())
            color = color.strip()
            if color_maxs.get(color) < amount:
                color_maxs[color] = amount
    return color_maxs["red"], color_maxs["green"], color_maxs["blue"]


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
def test_extract_valid_games(
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


@pytest.mark.parametrize(
    ["game_log", "expected_result"],
    [
        [
            "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
            (4, 2, 6),
        ],
        [
            "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
            (1, 3, 4),
        ],
        [
            "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
            (20, 13, 6),
        ],
        [
            "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
            (14, 3, 15),
        ],
        [
            "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
            (6, 3, 2),
        ],
        [
            "Game 6: 1 blue, 3 green; 2 blue, 2 green",
            (0, 3, 2),
        ],
    ],
)
def test_extract_minimum_colors(
    game_log: str,
    expected_result: tuple[int, int, int],
):
    v = extract_min_amounts(
        game_log,
    )
    assert v == expected_result


if __name__ == "__main__":
    with open(BASE_DIR / "input.txt", "r") as f:
        # game_ids = extract_minimum_set_powers(
        #     f.readlines(),
        #     amount_red=12,
        #     amount_green=13,
        #     amount_blue=14,
        # )
        # print(game_ids)
        # print(sum(game_ids))
        powers = extract_minimum_set_powers(
            f.readlines(),
        )
        print(powers)
        print(sum(powers))
