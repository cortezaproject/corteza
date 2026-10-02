#!/usr/bin/env python3
"""Print the shape of a REST endpoint before calling it.

Usage: shape.py                         services and their endpoint groups
       shape.py compose                 endpoint groups of one service
       shape.py compose/record          every call of one group: name, method, path
       shape.py compose/record update   one call: method, path, parameters, traps
       shape.py --grep undelete         calls whose name or path contains a word

The source is the codegen input the handlers and request structs are generated
from, server/<service>/rest.yaml, so it is current by construction. The things
a definition cannot say — which verb an update really takes, where an error
hides, what a client requires that REST does not — live in traps.yaml beside
this script, keyed by service, service/group or service/group/call, with `*`
for everything.
"""
import fnmatch
import os
import sys

import yaml

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.environ.get("HUMAN_REPO_ROOT") or os.path.abspath(os.path.join(HERE, "..", ".."))
SERVICES = ("system", "compose", "automation", "federation")


def load():
    out = {}
    for svc in SERVICES:
        f = os.path.join(ROOT, "server", svc, "rest.yaml")
        if not os.path.exists(f):
            continue
        d = yaml.safe_load(open(f)) or {}
        groups = {}
        for ep in d.get("endpoints", []):
            name = ep.get("entrypoint") or ep.get("title", "").lower()
            scope = (ep.get("parameters") or {}).get("path") or []
            calls = {}
            for api in ep.get("apis", []):
                p = (api.get("parameters") or {})
                calls[api["name"]] = {
                    "method": api.get("method", "GET").upper(),
                    "path": "/" + svc + ep.get("path", "") + api.get("path", ""),
                    "title": api.get("title", ""),
                    "params": {
                        "path": scope + (p.get("path") or []),
                        "query": p.get("get") or [],
                        "body": p.get("post") or [],
                    },
                }
            groups[name] = {"title": ep.get("title", ""), "path": "/" + svc + ep.get("path", ""), "calls": calls}
        out[svc] = groups
    return out


def traps_for(key):
    f = os.path.join(HERE, "traps.yaml")
    if not os.path.exists(f):
        return []
    t = yaml.safe_load(open(f)) or {}
    hits = []
    for pattern, notes in t.items():
        pat = str(pattern)
        if fnmatch.fnmatch(key, pat) or fnmatch.fnmatch(key, pat + "/*") or pat == "*":
            hits.extend(notes if isinstance(notes, list) else [notes])
    return hits


def find_group(data, svc, group):
    groups = data.get(svc)
    if groups is None:
        sys.exit(f"no service {svc!r}; have {', '.join(data)}")
    low = {k.lower(): k for k in groups}
    want = group.lower().replace("-", "").replace("_", "")
    for k in low:
        if k.replace("-", "") == want:
            return low[k], groups[low[k]]
    near = [k for k in groups if want in k.lower()]
    sys.exit(f"no group {group!r} in {svc}; " + (f"did you mean {', '.join(near)}?" if near else f"have {', '.join(groups)}"))


def fmt_param(p):
    req = "*" if p.get("required") else " "
    return f"  {req} {p['name']:<22} {str(p.get('type', '')):<28} {p.get('title', '')}"


def show_call(svc, gname, cname, call):
    print(f"{call['method']} {call['path']}" + (f"    # {call['title']}" if call["title"] else ""))
    for loc, label in (("path", "path"), ("query", "query"), ("body", "body (JSON)")):
        ps = call["params"][loc]
        if ps:
            print(f"{label}:")
            for p in ps:
                print(fmt_param(p))
    traps = traps_for(f"{svc}/{gname}/{cname}")
    if traps:
        print("traps:")
        for t in traps:
            print(f"  - {t}")


def main(argv):
    data = load()
    if argv and argv[0] == "--grep":
        word = argv[1].lower()
        for svc, groups in data.items():
            for g, grp in groups.items():
                for c, call in grp["calls"].items():
                    if word in c.lower() or word in call["path"].lower():
                        print(f"{svc}/{g} {c:<16} {call['method']:<6} {call['path']}")
        return 0
    if not argv:
        for svc, groups in data.items():
            print(f"{svc}: {' '.join(groups)}")
        return 0
    target = argv[0].strip("/").split("/")
    svc = target[0]
    if svc not in data:
        sys.exit(f"no service {svc!r}; have {', '.join(data)}")
    if len(target) == 1 and len(argv) == 1:
        for g, grp in data[svc].items():
            print(f"{g:<24} {grp['path']}    {grp['title']}")
        traps = traps_for(svc)
        if traps:
            print("traps:")
            for t in traps:
                print(f"  - {t}")
        return 0
    gname, grp = find_group(data, svc, target[1])
    cname = target[2] if len(target) > 2 else (argv[1] if len(argv) > 1 else None)
    if cname is None:
        for c, call in grp["calls"].items():
            print(f"{c:<16} {call['method']:<6} {call['path']}")
        traps = traps_for(f"{svc}/{gname}")
        if traps:
            print("traps:")
            for t in traps:
                print(f"  - {t}")
        return 0
    low = {k.lower(): k for k in grp["calls"]}
    if cname.lower() not in low:
        sys.exit(f"no call {cname!r} in {svc}/{gname}; have {', '.join(grp['calls'])}")
    show_call(svc, gname, low[cname.lower()], grp["calls"][low[cname.lower()]])
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
