def main(file_name: str) -> int:
    valid_lines = []
    invalid_lines = []

    with open(file_name) as f:
        for line in f.readlines():
            print(line.strip(), end=" - ")
            numbers = [int(i) for i in line.strip().split(" ")]
            # print(numbers)
            # is_descending = numbers[1] < numbers[0]
            first_diff = numbers[1] - numbers[0]
            first_is_descending = first_diff < 0
            print(f"{first_is_descending=}", end="")
            previous = numbers[0]
            # print(f"{descending=}, {previous=}")
            for number in numbers[1:]:
                current_diff = number - previous
                current_is_descending = current_diff < 0
                if (
                    current_is_descending != first_is_descending
                    or abs(current_diff) > 3
                    or current_diff == 0
                ):
                    invalid_lines.append(line)
                    print(" -- INVALID", f"conflict: {previous}, {number}")
                    break
                previous = number
            else:
                valid_lines.append(line)
                print(" -- is valid")

    # print("Valid lines:")
    # for line in invalid_lines:
    #     print(line)
    print(len(valid_lines))
    return len(valid_lines)


if __name__ == "__main__":
    # main("test_input.txt")
    main("input.txt")
