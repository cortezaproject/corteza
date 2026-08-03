#!/usr/bin/env python3
"""Verify the local Human MCP surface end to end.

Usage:
  mcp-verify.py             # all sections
  mcp-verify.py auth        # auth matrix only
  mcp-verify.py schema      # declared-contract checks only
  mcp-verify.py scope       # group filtering and risk ceiling
  mcp-verify.py exercise    # create/read/page/delete against real data
  mcp-verify.py cost        # tool-list size in tokens

Why this exists: unit tests assert what the code declares, not what the server
answers. The /api/mcp auth fix passed its unit tests and did nothing, because
HttpTokenValidator lets a request with no token through — that was only visible
by asking the running server. Run this after every batch.

Exercise mode writes to the dev server using the `agent-` prefix required by
CLAUDE.md, and deletes what it creates. Local-only.
"""

import json
import os
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

AGENT_DIR = os.path.dirname(os.path.abspath(__file__))
API = os.environ.get("HUMAN_API", "http://localhost:1043/api")
MCP = API + "/mcp"

if urllib.parse.urlparse(API).hostname not in ("localhost", "127.0.0.1", "::1"):
    sys.exit(f"mcp-verify.py is local-only; refusing to touch {API}")

TOKEN = subprocess.run(
    [os.path.join(AGENT_DIR, "token.sh")], capture_output=True, text=True, check=True
).stdout.strip()

PREFIX = f"agent-verify-{int(time.time())}"

failures = []
checks = 0


def check(ok, label, detail=""):
    global checks
    checks += 1
    if ok:
        print(f"  \033[32m✓\033[0m {label}")
    else:
        print(f"  \033[31m✗\033[0m {label}{(' — ' + detail) if detail else ''}")
        failures.append(label)


def post(body, sid=None, token=TOKEN, raw_status=False, url=None):
    h = {"Content-Type": "application/json", "Accept": "application/json, text/event-stream"}
    if token is not None:
        h["Authorization"] = f"Bearer {token}"
    if sid:
        h["Mcp-Session-Id"] = sid
    req = urllib.request.Request(url or MCP, data=json.dumps(body).encode(), headers=h, method="POST")
    try:
        with urllib.request.urlopen(req, timeout=60) as r:
            body_text, status, new_sid = r.read().decode(), r.status, r.headers.get("Mcp-Session-Id")
    except urllib.error.HTTPError as e:
        if raw_status:
            return None, None, e.code
        raise
    if body_text.startswith(("event:", "data:")) or "\ndata:" in body_text:
        body_text = next((l[5:].strip() for l in body_text.splitlines() if l.startswith("data:")), "")
    parsed = json.loads(body_text) if body_text.strip() else None
    return (parsed, new_sid, status) if raw_status else (parsed, new_sid)


def rpc(method, params, sid=None, rid=1, url=None):
    resp, new_sid = post({"jsonrpc": "2.0", "id": rid, "method": method, "params": params}, sid, url=url)
    if resp and "error" in resp:
        raise RuntimeError(json.dumps(resp["error"]))
    return (resp or {}).get("result"), new_sid


def session(url=None):
    _, sid = rpc(
        "initialize",
        {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "mcp-verify", "version": "1"}},
        url=url,
    )
    post({"jsonrpc": "2.0", "method": "notifications/initialized"}, sid, url=url)
    return sid


def list_tools_at(url):
    sid = session(url)
    tools, _ = rpc("tools/list", {}, sid, rid=2, url=url)
    return sid, sorted(tools["tools"], key=lambda t: t["name"])


def call(sid, name, args, rid=90, url=None):
    result, _ = rpc("tools/call", {"name": name, "arguments": args}, sid, rid=rid, url=url)
    text = "".join(c.get("text", "") for c in result.get("content", []))
    if result.get("isError"):
        raise RuntimeError(f"{name}: {text}")
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        return text


# ---------------------------------------------------------------- auth matrix

INIT = {
    "jsonrpc": "2.0", "id": 1, "method": "initialize",
    "params": {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "v", "version": "1"}},
}


def status_with(token):
    _, _, status = post(INIT, token=token, raw_status=True)
    return status


def verify_auth():
    print("\nauth matrix")
    check(status_with(None) == 401, "no token is rejected",
          "tool discovery lists every tool and description, so this must not be public")
    check(status_with(TOKEN) == 200, "valid token is accepted")

    # A malformed token currently yields 500 rather than 401. That is global to
    # the API (/system/settings/current behaves the same) and lives in the
    # shared verifier, so it is reported, not asserted.
    bogus = status_with("not-a-real-token")
    if bogus != 401:
        print(f"  \033[33m·\033[0m malformed token returns {bogus}, not 401 — known, global, not MCP-specific")


# ------------------------------------------------------------ declared schema


def verify_schema(sid):
    print("\ndeclared contracts")
    tools, _ = rpc("tools/list", {}, sid, rid=2)
    tools = sorted(tools["tools"], key=lambda t: t["name"])
    check(len(tools) > 0, f"tools/list returns {len(tools)} tools")

    untagged, bad_annotation, bad_ids, bad_lookup = [], [], [], []

    for t in tools:
        name, meta = t["name"], t.get("_meta") or {}
        risk = meta.get("human.dev/risk")
        groups = meta.get("human.dev/groups")
        ann = t.get("annotations") or {}
        props = (t.get("inputSchema") or {}).get("properties") or {}
        required = (t.get("inputSchema") or {}).get("required") or []

        if not risk or not groups:
            untagged.append(name)

        if ann.get("openWorldHint") is not False:
            bad_annotation.append(f"{name}: openWorldHint")
        if risk == "read" and (ann.get("readOnlyHint") is not True or ann.get("destructiveHint") is not False):
            bad_annotation.append(f"{name}: read tool annotated {ann}")
        if risk == "destructive" and ann.get("destructiveHint") is not True:
            bad_annotation.append(f"{name}: destructive tool not annotated destructive")

        for p, spec in props.items():
            if (p.endswith("ID") or p.endswith("IDs")) and spec.get("type") != "string":
                bad_ids.append(f"{name}.{p} is {spec.get('type')}")

        if name.endswith("_lookup"):
            resource = name[name.index("_") + 1: name.rindex("_")]
            for ref in (resource, resource + "ID"):
                if ref in required:
                    bad_lookup.append(f"{name} requires its own ref {ref}")
            if "limit" not in props:
                bad_lookup.append(f"{name} declares no limit")

    check(not untagged, "every tool carries group and risk in _meta", ", ".join(untagged))
    check(not bad_annotation, "annotations agree with declared risk", "; ".join(bad_annotation))
    check(not bad_ids, "ID params are declared as strings", "; ".join(bad_ids))
    check(not bad_lookup, "lookups are bounded and their ref is optional", "; ".join(bad_lookup))
    return tools


# ------------------------------------------------------------------- exercise


def verify_exercise(sid):
    """Create real data, read it back, page it, then remove it.

    Paging is the point: a cursor that cannot be fed back is worse than no
    cursor, and that defect shipped once already.
    """
    print("\nexercise (writes to the dev server, cleans up after)")
    ns = mod = None
    try:
        ns = call(sid, "compose_namespace_create", {"name": PREFIX, "slug": PREFIX})
        ns_id = ns.get("namespaceID") or ns.get("ID")
        check(bool(ns_id), "namespace created", json.dumps(ns)[:120])

        mod = call(sid, "compose_module_create", {
            "namespace": PREFIX, "name": f"{PREFIX}-mod", "handle": "agent_verify_mod",
            "fields": json.dumps([{"name": "title", "kind": "String"}]),
        })
        check(bool(mod), "module created")

        for i in range(3):
            call(sid, "compose_record_create", {
                "namespace": PREFIX, "module": "agent_verify_mod",
                "values": json.dumps({"title": f"row {i}"}),
            })

        first = call(sid, "compose_record_lookup", {
            "namespace": PREFIX, "module": "agent_verify_mod", "limit": "2",
        })
        rows = first.get("records") or []
        cursor = first.get("nextPageCursor")
        check(len(rows) == 2, f"limit is honoured (got {len(rows)})")
        check(bool(cursor), "a next-page cursor is returned")

        if cursor:
            second = call(sid, "compose_record_lookup", {
                "namespace": PREFIX, "module": "agent_verify_mod",
                "limit": "2", "pageCursor": cursor,
            })
            check(isinstance(second.get("records"), list),
                  "the cursor round-trips", "a cursor that cannot be fed back is the defect this catches")

        listing = call(sid, "compose_namespace_lookup", {})
        check("namespaces" in listing,
              "namespace list mode is reachable with no ref",
              "it used to be unreachable because the ref was marked required")
    finally:
        if ns is not None:
            try:
                call(sid, "compose_namespace_delete", {"namespace": PREFIX})
                print(f"  \033[90m·\033[0m cleaned up {PREFIX}")
            except Exception as e:  # noqa: BLE001 - cleanup must not mask a real failure
                print(f"  \033[33m·\033[0m cleanup failed, remove {PREFIX} by hand: {e}")


# ----------------------------------------------------------------- token cost


def verify_scope():
    """Group narrows the listing; risk narrows it AND refuses on dispatch.

    The distinction is the point: a filter only affects tools/list, so a client
    that already knew a tool's name could still call it. If the ceiling did not
    refuse at dispatch it would mean nothing.
    """
    print("\nscope (group filtering and risk ceiling)")

    _, all_tools = list_tools_at(MCP)
    names = {t["name"] for t in all_tools}

    _, configuring = list_tools_at(MCP + "/configuring")
    _, usage = list_tools_at(MCP + "/usage")

    cfg, use = {t["name"] for t in configuring}, {t["name"] for t in usage}
    check(cfg and use, f"group endpoints list subsets (configuring {len(cfg)}, usage {len(use)})")
    check(cfg < names and use < names, "each group is a strict subset of the full list")
    check(not (cfg & use), "configuring and usage do not overlap",
          ", ".join(sorted(cfg & use)))
    check(cfg | use == names, "the groups together account for every tool",
          "missing: " + ", ".join(sorted(names - (cfg | use))))

    _, read_only = list_tools_at(MCP + "?maxRisk=read")
    risks = {(t.get("_meta") or {}).get("human.dev/risk") for t in read_only}
    check(risks <= {"read"}, f"a read ceiling lists only read tools (saw {sorted(r for r in risks if r)})")

    # The part a filter cannot do: refuse a call by name.
    sid = session(MCP + "?maxRisk=read")
    writer = next((t["name"] for t in all_tools
                   if (t.get("_meta") or {}).get("human.dev/risk") == "destructive"), None)
    if writer:
        try:
            call(sid, writer, {}, url=MCP + "?maxRisk=read")
            check(False, f"a capped session refuses to dispatch {writer}", "the call went through")
        except RuntimeError as e:
            check("capped at" in str(e), f"a capped session refuses to dispatch {writer}", str(e)[:120])


def verify_cost(tools):
    print("\ntool-list cost")
    payload = json.dumps({"tools": tools})
    # ~4 chars per token is the usual rule of thumb; exact enough to argue about
    # whether group filtering is worth building.
    approx = len(payload) // 4
    print(f"  {len(tools)} tools, {len(payload):,} chars, ~{approx:,} tokens per request")
    if tools:
        worst = max(tools, key=lambda t: len(json.dumps(t)))
        print(f"  largest: {worst['name']} (~{len(json.dumps(worst)) // 4:,} tokens)")
    print(f"  projected at 200 tools: ~{approx * 200 // max(len(tools), 1):,} tokens per request")


def main():
    want = sys.argv[1] if len(sys.argv) > 1 else "all"
    print(f"verifying {MCP}")

    if want in ("all", "auth"):
        verify_auth()

    sid = session()
    tools = []
    if want in ("all", "schema", "cost"):
        tools = verify_schema(sid)
    if want in ("all", "scope"):
        verify_scope()
    if want in ("all", "exercise"):
        verify_exercise(sid)
    if want in ("all", "cost"):
        verify_cost(tools)

    print(f"\n{checks - len(failures)}/{checks} checks passed")
    if failures:
        print("failed: " + "; ".join(failures))
        sys.exit(1)


if __name__ == "__main__":
    main()
