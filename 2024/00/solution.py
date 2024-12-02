import argparse


def cleanup_line(line: str) -> list[int]:
    return [int(i) for i in line.strip().split(" ")]

def process_values(values: list[int]) -> list[int]:
    """

    """
    result = []
    for index, number in enumerate(values):
        ...
    return result


def main(file_name: str) -> int:

    with open(file_name, "r") as f:
        total = []
        for line in f:
            values = cleanup_line(line)
            result = process_values(values)
            total.append(result)

    return len(total)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Advent of Code - 2024 - Day 0 - Part 1")
    parser.add_argument("filename", type=str, help="The input file.")

    args = parser.parse_args()

    print(main(args.filename))
