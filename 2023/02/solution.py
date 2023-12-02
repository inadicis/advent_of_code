from collections.abc import Sequence
from pathlib import Path
from typing import Iterable, Mapping
import pytest

BASE_DIR = Path(__file__).resolve().parent


def extract_possible_games(
    game_logs: list[str],
    amount_per_color: Mapping[str, int],
) -> list[int]:
    valid_games = []
    for i, game_log in enumerate(game_logs):
        if is_game_possible(game_log, amount_per_color):
            valid_games.append(i + 1)

    return valid_games


def extract_minimum_set_powers(game_logs: list[str]) -> list[int]:
    set_powers = []
    for i, game_log in enumerate(game_logs):
        min_colors = extract_min_amounts(game_log)
        power = 1
        # now works with any amount of colors,
        # but if one color is not needed, it will not multiply by 0
        for amount in min_colors.values():
            power *= amount
        set_powers.append(power)

    return set_powers


def get_clean_game_log(game_log: str) -> list[dict[str, int]]:
    """Eliminates game id, splits the different pulls, and for each
    one builds a mapping from colors to their amount"""
    clean_pulls = []
    game_id, pulls = game_log.split(":")
    pulls = pulls.split(";")
    for pull in pulls:
        color_pulls = pull.split(",")
        amount_per_color = {}
        for color_pull in color_pulls:
            amount, color = color_pull.strip().split(" ")
            amount_per_color[color] = int(amount)
        clean_pulls.append(amount_per_color)
    return clean_pulls


def is_game_possible(
    game_log: str,
    amount_per_color: Mapping[str, int],
) -> bool:
    for pull in get_clean_game_log(game_log):
        for color, amount in pull.items():
            if amount > amount_per_color.get(color, 0):
                # if a color was not provided in the mapping,
                # we interpret it as not being available (amount 0)
                return False

    return True


def extract_min_amounts(game_log: str) -> Mapping[str, int]:  # works with any color
    color_maxs = {}
    for pull in get_clean_game_log(game_log):
        for color, amount in pull.items():
            if color_maxs.get(color, 0) < amount:
                color_maxs[color] = amount

    return color_maxs


@pytest.mark.parametrize(
    ["game_logs", "amount_per_color", "expected_result"],
    [
        [
            [
                "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
                "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
                "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
                "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
                "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
            ],
            {"blue": 12, "green": 13, "red": 14},
            [1, 2, 5],
        ],
    ],
)
def test_extract_valid_games(
    game_logs: list[str],
    amount_per_color: Mapping[str, int],
    expected_result: list[int],
):
    v = extract_possible_games(game_logs, amount_per_color=amount_per_color)
    assert v == expected_result


@pytest.mark.parametrize(
    ["game_log", "expected_result"],
    [
        [
            "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
            {"red": 4, "green": 2, "blue": 6},
        ],
        [
            "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
            {"red": 1, "green": 3, "blue": 4},
        ],
        [
            "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
            {"red": 20, "green": 13, "blue": 6},
        ],
        [
            "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
            {"red": 14, "green": 3, "blue": 15},
        ],
        [
            "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
            {"red": 6, "green": 3, "blue": 2},
        ],
        [
            "Game 6: 1 blue, 3 green; 2 blue, 2 green",
            {"green": 3, "blue": 2},
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
        game_ids = extract_possible_games(
            f.readlines(), amount_per_color={"blue": 12, "green": 13, "red": 14}
        )
        print(game_ids)
        print(sum(game_ids))
        # powers = extract_minimum_set_powers(
        #     f.readlines(),
        # )
        # print(powers)
        # print(sum(powers))
