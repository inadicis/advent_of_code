def find_conflicts(numbers: list[int]) -> list[int]:
    """
    returns the indexes of the numbers that conflict with the previous one,
    assuming the first two numbers define the order
    """

    first_diff = numbers[1] - numbers[0]
    first_is_descending = first_diff < 0
    previous = numbers[0]
    conflicts = []
    for index, number in enumerate(numbers[1:]):
        current_diff = number - previous
        current_is_descending = current_diff < 0
        if (
            current_is_descending != first_is_descending
            or abs(current_diff) > 3
            or current_diff == 0
        ):
            conflicts.append(index)
        previous = number
    return conflicts


def main(file_name: str) -> int:
    valid_lines = []
    invalid_lines = []

    with open(file_name) as f:
        for line in f.readlines():
            numbers = [int(i) for i in line.strip().split(" ")]
            conflicts = find_conflicts(numbers)

            # if len(conflicts) > 3:
            #     # cannot fix with only one removal for sure
            #     # (at least 2 independent pairs of involved numbers,
            #     # removing one would fix maximum 1 pair)
            #     invalid_lines.append(line)
            # elif len(conflicts) == 0:
            #     # with 0 conflict, obviously valid
            #     valid_lines.append(line)
            # else:
            if find_conflicts(numbers) == []:
                valid_lines.append(line)
                continue

            for index, _ in enumerate(numbers):
                if find_conflicts(numbers[:index] + numbers[index + 1 :]) == []:
                    valid_lines.append(line)
                    break
            else:
                invalid_lines.append(line)

    return len(valid_lines)


if __name__ == "__main__":
    # main("test_input.txt")
    print(main("input.txt"))
