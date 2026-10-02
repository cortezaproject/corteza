#!/usr/bin/env python3
"""Shape a JSON body for reading.

Usage: jsonout.py [--pick EXPR] [--limit N] [--compact|--pretty] < body

EXPR is a path into the body:
  response.set[].email                   each element's email, one per line
  response.set[0].name                   one element
  response.set[].{recordID,values.name}  keep these keys of each element
A missing key yields null rather than an error, so a typo is visible in the
output instead of fatal.

Output: a scalar prints raw; a list of scalars prints one per line; anything
else prints as JSON, pretty when stdout is a TTY or --pretty, compact otherwise.

--limit N (default 20): every list longer than N is cut to N items and a final
"… K more (--all)" string marks the cut; 0 disables. The marker is a string
inside the list, so it is unmistakable to a reader and never silent.

Exit 1 and print nothing when stdin is not JSON, so the caller can tell a
non-JSON body from an empty one.
"""
import json
import sys


def parse(expr):
    """Turn an EXPR into steps: ('key', name) | ('each',) | ('index', n) | ('keys', [expr…])."""
    steps, i, n = [], 0, len(expr)
    buf = ""

    def flush():
        nonlocal buf
        if buf:
            steps.append(("key", buf))
            buf = ""

    while i < n:
        c = expr[i]
        if c == ".":
            flush()
            i += 1
        elif c == "[":
            flush()
            j = expr.index("]", i)
            inner = expr[i + 1 : j]
            steps.append(("each",) if inner == "" else ("index", int(inner)))
            i = j + 1
        elif c == "{":
            flush()
            depth, j = 0, i
            while True:
                if expr[j] == "{":
                    depth += 1
                elif expr[j] == "}":
                    depth -= 1
                    if depth == 0:
                        break
                j += 1
            parts, depth, start = [], 0, i + 1
            for k in range(i + 1, j):
                ch = expr[k]
                if ch == "{":
                    depth += 1
                elif ch == "}":
                    depth -= 1
                elif ch == "," and depth == 0:
                    parts.append(expr[start:k].strip())
                    start = k + 1
            parts.append(expr[start:j].strip())
            steps.append(("keys", parts))
            i = j + 1
        else:
            buf += c
            i += 1
    flush()
    return steps


def apply(value, steps):
    for pos, step in enumerate(steps):
        kind = step[0]
        if kind == "each":
            rest = steps[pos + 1 :]
            if not isinstance(value, list):
                return None
            return [apply(v, rest) for v in value]
        if kind == "key":
            value = value.get(step[1]) if isinstance(value, dict) else None
        elif kind == "index":
            value = value[step[1]] if isinstance(value, list) and -len(value) <= step[1] < len(value) else None
        elif kind == "keys":
            value = {sub: apply(value, parse(sub)) for sub in step[1]}
    return value


def truncate(value, limit):
    if limit <= 0:
        return value
    if isinstance(value, list):
        out = [truncate(v, limit) for v in value[:limit]]
        if len(value) > limit:
            out.append(f"… {len(value) - limit} more (--all)")
        return out
    if isinstance(value, dict):
        return {k: truncate(v, limit) for k, v in value.items()}
    return value


def main(argv):
    pick, limit, pretty = None, 20, sys.stdout.isatty()
    args = list(argv)
    while args:
        a = args.pop(0)
        if a == "--pick":
            pick = args.pop(0)
        elif a == "--limit":
            limit = int(args.pop(0))
        elif a == "--compact":
            pretty = False
        elif a == "--pretty":
            pretty = True
        else:
            sys.exit(f"jsonout.py: unknown option {a}")

    try:
        body = json.load(sys.stdin)
    except ValueError:
        return 1

    if pick:
        body = apply(body, parse(pick))
    body = truncate(body, limit)

    scalar = lambda v: v is None or isinstance(v, (str, int, float, bool))
    if scalar(body):
        print("null" if body is None else body if isinstance(body, str) else json.dumps(body))
    elif isinstance(body, list) and all(scalar(v) for v in body):
        for v in body:
            print("null" if v is None else v if isinstance(v, str) else json.dumps(v))
    elif pretty:
        print(json.dumps(body, indent=2, ensure_ascii=False))
    else:
        print(json.dumps(body, separators=(",", ":"), ensure_ascii=False))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
