import sys

def avg(arr: list[int]) -> float:
    return sum(arr)/len(arr)

def main():
    filename = sys.argv[1]
    with open(filename, "r") as f:
        vals = list(map(int, f.read().strip().split(",")))
    print(f"Total values {len(vals)}")
    print(f"Min value is {min(vals)}")
    print(f"Max value is {max(vals)}")
    print(f"Avg value is {avg(vals)}")
    sorted_vals = sorted(vals)
    nine_five = round(len(sorted_vals)*0.95)
    print(f"95th percentile is {sorted_vals[nine_five]}")
    nine_nine = round(len(sorted_vals)*0.99)
    print(f"99th percentile is {sorted_vals[nine_nine]}")


if __name__ == "__main__":
    main()
