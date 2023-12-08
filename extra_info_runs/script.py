import sys
import re
from dataclasses import dataclass

@dataclass
class RequestInfo:
    operation: float
    http: float

def main():
    with open(sys.argv[1], "r") as f:
        lines = f.read().strip().split("\n")
    lines = [l for l in lines if "=" not in l] # Filter request id lines
    matcher = re.compile("\[(.*)\] ([a-z]*): ([0-9.]*)(m?)s")
    values = []
    for l in lines:
        values.append(matcher.findall(l)[0])
    request_info = dict()
    for v in values:
        if v[0] in request_info:
            curr = request_info[v[0]]
        else:
            curr = RequestInfo(0,0)
        val = float(v[2])
        if v[3] != 'm':
            val *= 1000 # Value was in seconds
        if v[1] == "http":
            curr.http = float(v[2])
        else:
            curr.operation = float(v[2])
        request_info[v[0]] = curr
    for k,v in request_info.items():
        delay = v.operation-v.http
        if delay > 1:
            print(v)


if __name__ == "__main__":
    main()
