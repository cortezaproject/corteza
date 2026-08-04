#!/usr/bin/env python3
"""Verify the local Human MCP surface end to end.

Usage:
  mcp-verify.py             # all sections
  mcp-verify.py auth        # auth matrix only
  mcp-verify.py schema      # declared-contract checks only
  mcp-verify.py scope       # group filtering and risk ceiling
  mcp-verify.py exercise    # create/read/page/delete against real data
  mcp-verify.py identity    # create/wire/read/delete users, groups, roles
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
# Progressive disclosure means tools/list shows ~5 by default. Auditing the
# surface needs the documented tooling opt-out; agents must never use it.
MCP_ALL = MCP + "?tools=all"

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
    # Match the reply by id. Since progressive disclosure fires
    # notifications/tools/list_changed, a response frame is no longer
    # necessarily the first data: line in the stream.
    msgs = []
    if body_text.startswith(("event:", "data:")) or "\ndata:" in body_text:
        for line in body_text.splitlines():
            if line.startswith("data:"):
                try:
                    msgs.append(json.loads(line[5:].strip()))
                except json.JSONDecodeError:
                    pass
    elif body_text.strip():
        msgs = [json.loads(body_text)]

    want = body.get("id")
    parsed = next((m for m in msgs if m.get("id") == want), msgs[0] if msgs else None)
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

    # A malformed token used to yield 500: jwtauth's parse error carries no
    # Kind, so it rendered as a server fault rather than a rejected credential.
    check(status_with("not-a-real-token") == 401, "a malformed token is rejected as 401, not 500")


# ------------------------------------------------------------ declared schema


def verify_schema(_sid):
    print("\ndeclared contracts")
    _, tools = list_tools_at(MCP_ALL)
    check(len(tools) > 0, f"the full surface is {len(tools)} tools")

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
    load_all(sid)
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

    _, all_tools = list_tools_at(MCP_ALL)
    names = {t["name"] for t in all_tools}

    _, configuring = list_tools_at(MCP + "/configuring?tools=all")
    _, usage = list_tools_at(MCP + "/usage?tools=all")

    meta = {"human_tool_search", "human_tool_load"}
    cfg = {t["name"] for t in configuring} - meta
    use = {t["name"] for t in usage} - meta
    names = names - meta
    check(cfg and use, f"group endpoints list subsets (configuring {len(cfg)}, usage {len(use)})")
    check(cfg < names and use < names, "each group is a strict subset of the full list")
    check(not (cfg & use), "configuring and usage do not overlap",
          ", ".join(sorted(cfg & use)))
    check(cfg | use == names, "the groups together account for every tool",
          "missing: " + ", ".join(sorted(names - (cfg | use))))

    _, read_only = list_tools_at(MCP + "?maxRisk=read&tools=all")
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


def verify_identity(sid):
    """Exercise the identity tools against real data.

    These 38 tools shipped build-verified only. Identity is the batch where a
    contract mistake is least visible from a schema and most consequential in
    effect, so this creates a user, a group and a role, wires them together,
    reads the wiring back, and removes all three.
    """
    print("\nidentity (writes to the dev server, cleans up after)")
    load_all(sid)

    email = f"{PREFIX}@local.dev"
    created = {}

    try:
        user = call(sid, "system_user_create", {
            "email": email, "handle": PREFIX.replace("-", "_"), "name": "Verify User",
        })
        created["user"] = user.get("userID")
        check(bool(created["user"]), "user created", json.dumps(user)[:120])

        role = call(sid, "system_role_create", {
            "name": f"{PREFIX} role", "handle": PREFIX.replace("-", "_") + "_role",
        })
        created["role"] = role.get("roleID")
        check(bool(created["role"]), "role created")

        # Every non-root group needs a parent — checkPaths rejects a group with
        # no paths — so find the root and hang the test group off it.
        groups = call(sid, "system_user_group_lookup", {"limit": "100"})
        root = next((g for g in (groups.get("userGroups") or groups.get("groups") or [])
                     if g.get("isRoot")), None)
        check(root is not None, "the root user group is discoverable via isRoot",
              json.dumps(groups)[:160])

        if root:
            group = call(sid, "system_user_group_create", {
                "handle": PREFIX.replace("-", "_") + "_grp",
                "short": "verify group",
                "parents": json.dumps([root.get("handle") or root.get("userGroupID")]),
            })
            created["group"] = group.get("userGroupID")
            check(bool(created["group"]), "user group created under the root group")

        # The wiring is the part a schema cannot prove: a member added through
        # one tool must be visible through another.
        call(sid, "system_role_member_add", {"role": created["role"], "user": created["user"]})
        members = call(sid, "system_role_member_list", {"role": created["role"]})
        blob = json.dumps(members)
        check(str(created["user"]) in blob, "an added member shows up in member_list", blob[:160])

        # And the reverse direction, which goes through a different filter path.
        holders = call(sid, "system_user_lookup", {"role": created["role"]})
        check(str(created["user"]) in json.dumps(holders),
              "the member is findable by filtering users on that role")

        found = call(sid, "system_user_lookup", {"user": email})
        check(str(created["user"]) in json.dumps(found), "a user resolves by email")

        # Suspend is not delete, and suspended users are hidden by default.
        call(sid, "system_user_suspend", {"user": created["user"]})
        visible = json.dumps(call(sid, "system_user_lookup", {"query": PREFIX}))
        hidden = str(created["user"]) not in visible
        check(hidden, "a suspended user drops out of the default listing")

        withsusp = call(sid, "system_user_lookup", {"query": PREFIX, "includeSuspended": True})
        check(str(created["user"]) in json.dumps(withsusp), "includeSuspended brings it back")
        call(sid, "system_user_unsuspend", {"user": created["user"]})

    finally:
        for tool, key in (("system_role_delete", "role"),
                          ("system_user_group_delete", "group"),
                          ("system_user_delete", "user")):
            if created.get(key):
                try:
                    call(sid, tool, {_refname(tool): created[key]})
                except Exception as e:  # noqa: BLE001 - cleanup must not mask a failure
                    print(f"  \033[33m·\033[0m cleanup of {key} failed, remove by hand: {e}")
        print(f"  \033[90m·\033[0m cleaned up {PREFIX} user/group/role")


def _refname(tool):
    """The ref param for a delete tool: system_user_group_delete -> userGroup."""
    resource = tool[len("system_"): -len("_delete")]
    head, *rest = resource.split("_")
    return head + "".join(p.title() for p in rest)


def load_all(sid):
    """Pull every tool into this session so the exercises can call them."""
    _, every = list_tools_at(MCP_ALL)
    names = [t["name"] for t in every if not t["name"].startswith("human_tool_")]
    for i in range(0, len(names), 40):
        call(sid, "human_tool_load", {"names": json.dumps(names[i:i + 40])}, rid=80 + i)


def verify_cost(tools):
    print("\ntool-list cost")
    payload = json.dumps({"tools": tools})
    # ~4 chars per token is the usual rule of thumb; exact enough to argue about
    # whether group filtering is worth building.
    approx = len(payload) // 4
    print(f"  full surface: {len(tools)} tools, ~{approx:,} tokens")

    _, initial = list_tools_at(MCP)
    start = len(json.dumps(initial)) // 4
    print(f"  what a session actually starts with: {len(initial)} tools, ~{start:,} tokens"
          f"  ({approx / max(start, 1):.0f}x smaller)")
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
    if want in ("all", "identity"):
        verify_identity(sid)
    if want in ("all", "cost"):
        verify_cost(tools)

    print(f"\n{checks - len(failures)}/{checks} checks passed")
    if failures:
        print("failed: " + "; ".join(failures))
        sys.exit(1)


if __name__ == "__main__":
    main()
