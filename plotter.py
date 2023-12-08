import sys
import matplotlib.pyplot as plt

def read_file(filename):
    with open(filename, "r") as f:
        return list(map(int, f.read().strip().split(",")))

def avg(arr):
    return sum(arr)/len(arr)

def main():
    if len(sys.argv) < 2:
        print("Input filename")
        sys.exit(0)
    c7 = read_file(sys.argv[1])
    plt.title("Fetch times")
    x = [i for i in range(len(c7))]
    plt.scatter(x, c7, s=1)
    plt.ylabel("Time (Milliseconds)")
    plt.xlabel("Request number")
    plt.savefig("c7FetchTimes512b.png")

if __name__ == "__main__":
    main()
