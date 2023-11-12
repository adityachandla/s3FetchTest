import matplotlib.pyplot as plt

def read_files():
    with open("./t3.nano_512b_1000.csv", "r") as f:
        t3 = f.read().strip().split(",")
    t3 = [int(i) for i in t3]
    with open("./m6in.xlarge_512b_1000.csv", "r") as f:
        m6 = f.read().strip().split(",")
    m6 = [int(i) for i in m6]
    return (t3, m6)

def avg(arr):
    s = 0
    for i in arr:
        s += i
    return s/len(arr)

def main():
    t3,m6 = read_files()
    assert len(t3) == len(m6)
    x = [2*(i+1) for i in range(len(t3))]
    plt.title("Fetch times")
    ax = plt.gca()
    ax.set_xticks([])
    plt.scatter(x, t3, s=1, label="t3.micro")
    plt.ylabel("Time (Milliseconds)")
    plt.scatter(x, m6, s=1, label="m6in.xlarge")
    plt.legend(loc="upper right")
    plt.savefig("comparision.png")
    print(f"T3 running times: min={min(t3)}, max={max(t3)}, avg={avg(t3)}")
    print(f"M6In running times: min={min(m6)}, max={max(m6)}, avg={avg(m6)}")

if __name__ == "__main__":
    main()
