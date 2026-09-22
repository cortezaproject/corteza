#!/usr/bin/env python3
"""Score whether a skill still makes a fresh session write deployable apps.

Usage:
  skill_eval.py [evals/custom_app.json] [--only ID] [--model sonnet]

Each prompt runs in a fresh `claude -p` session that sees nothing but this
checkout's MCP server — no project CLAUDE.md, no pasted skill — which is how a
sales rep in Claude Web meets it. What the session writes is scored on what a
custom app needs to deploy and run: the skill was read before the page was
written, the page is plain HTML with the bridge snippet unchanged, it carries
sample data and says which mode it is in, and the deploy guard accepts it.

Costs money: every prompt is a real session (about $0.30 on Sonnet). Never runs
on Fable; pass --model to use another model.

Transcripts and pages land in .state/evals/<timestamp>/.
"""

import argparse
import json
import os
import re
import subprocess
import sys
import tempfile
import time

import mcp

AGENT_DIR = os.path.dirname(os.path.abspath(__file__))

# The line of the bridge snippet that only an unchanged copy carries.
SNIPPET = "window.human = window.human ||"
HELLO_V2 = re.compile(r"human:hello['\"]?\s*,\s*v\s*:\s*2\b")


def call_tool(sid, name, args):
    """A tool call whose refusal comes back as text rather than an exit."""
    resp, _ = mcp.post(
        {"jsonrpc": "2.0", "id": 9, "method": "tools/call", "params": {"name": name, "arguments": args}}, sid
    )
    if resp and "error" in resp:
        return None, resp["error"].get("message", json.dumps(resp["error"]))
    result = (resp or {}).get("result") or {}
    text = "".join(c.get("text", "") for c in result.get("content", []))
    if result.get("isError"):
        return None, text
    return json.loads(text) if text.startswith("{") else text, None


def namespaces_present(sid, slugs):
    found, _ = call_tool(sid, "compose_namespace_lookup", {"limit": "200"})
    have = {n.get("slug") for n in (found or {}).get("namespaces", [])}
    return [s for s in slugs if s not in have]


def run_session(prompt, model, config, workdir, transcript):
    cmd = [
        "claude", "-p", prompt,
        "--model", model,
        "--strict-mcp-config", "--mcp-config", config,
        "--allowedTools", "mcp__human__*",
        "--setting-sources", "",
        "--no-session-persistence",
        "--output-format", "stream-json", "--verbose",
        "--max-turns", "30",
    ]
    with open(transcript, "w") as out:
        subprocess.run(cmd, cwd=workdir, stdin=subprocess.DEVNULL, stdout=out, stderr=subprocess.DEVNULL, timeout=900)
    return read_transcript(transcript)


def read_transcript(transcript):
    """The tool calls, the final answer and the cost of one stream-json run."""
    calls, final, cost = [], "", None
    for line in open(transcript):
        try:
            event = json.loads(line)
        except ValueError:
            continue
        if event.get("type") == "assistant":
            for part in event["message"].get("content", []):
                if part.get("type") == "tool_use":
                    calls.append((part["name"].replace("mcp__human__", ""), part.get("input", {})))
        if event.get("type") == "result":
            final = event.get("result") or ""
            cost = event.get("total_cost_usd")
    return calls, final, cost


def page_of(text):
    fenced = re.findall(r"```html\s*\n(.*?)```", text, re.S | re.I)
    if fenced:
        return max(fenced, key=len)
    doc = re.search(r"<!doctype html.*?</html>", text, re.S | re.I)
    return doc.group(0) if doc else ""


def read_skill_first(calls, skill):
    """The skill was read, and before anything was written."""
    for name, args in calls:
        if name == "system_skill_lookup" and args.get("skill") == skill:
            return True
        if name in ("Write", "Edit"):
            return False
    return False


def guard_accepts(sid, page):
    stamp = f"skill-eval-{int(time.time() * 1000)}"
    app, err = call_tool(sid, "system_application_create", {
        "name": stamp, "enabled": False,
        "unify": json.dumps({"name": stamp, "listed": False, "kind": "custom"}),
    })
    if err:
        return False, f"could not create a scratch app: {err}"
    try:
        _, err = call_tool(sid, "system_application_source_set", {"application": app["applicationID"], "source": page})
        return err is None, err or ""
    finally:
        call_tool(sid, "system_application_delete", {"application": app["applicationID"]})


def written_page(calls):
    """The last page a session tried to save; the harness allows no writes, so
    a session that saves its page instead of printing it leaves it here."""
    for name, args in reversed(calls):
        if name == "Write" and str(args.get("file_path", "")).lower().endswith((".html", ".htm")):
            return args.get("content", "")
    return ""


def score(sid, skill, calls, final):
    page = page_of(final) or written_page(calls)
    checks = {
        "read skill first": read_skill_first(calls, skill),
        "plain html": bool(page) and "<script type=\"module\"" not in page and "import React" not in page,
        "snippet unchanged": SNIPPET in page and bool(HELLO_V2.search(page)),
        "sample data": "window.SAMPLE" in page,
        "mode badge": bool(re.search(r"sample data", page, re.I)) and "human.ready" in page,
    }
    accepted, why = guard_accepts(sid, page) if page else (False, "no page in the answer")
    checks["guard accepts"] = accepted
    return page, checks, why


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("suite", nargs="?", default=os.path.join(AGENT_DIR, "evals", "custom_app.json"))
    ap.add_argument("--only")
    ap.add_argument("--model", default="sonnet")
    opts = ap.parse_args()

    if "fable" in opts.model.lower():
        sys.exit("skill_eval.py does not run on Fable; pick another --model")

    suite = json.load(open(opts.suite))
    sid = mcp.session()

    missing = namespaces_present(sid, suite.get("needs", []))
    if missing:
        sys.exit(f"missing fixtures: {', '.join(missing)} — seed them with dev/agent/seed.sh {' '.join(missing)}")

    outdir = os.path.join(AGENT_DIR, ".state", "evals", time.strftime("%Y%m%d-%H%M%S"))
    os.makedirs(outdir, exist_ok=True)
    workdir = tempfile.mkdtemp(prefix="skill-eval-")
    config = os.path.join(outdir, "mcp.json")
    with open(config, "w") as fh:
        json.dump({"mcpServers": {"human": {
            "type": "http", "url": mcp.API + "/mcp", "headers": {"Authorization": "Bearer " + mcp.TOKEN},
        }}}, fh)
    os.chmod(config, 0o600)

    failed, spent = 0, 0.0
    for p in suite["prompts"]:
        if opts.only and p["id"] != opts.only:
            continue
        calls, final, cost = run_session(p["prompt"], opts.model, config, workdir, os.path.join(outdir, p["id"] + ".jsonl"))
        page, checks, why = score(sid, suite["skill"], calls, final)
        with open(os.path.join(outdir, p["id"] + ".html"), "w") as fh:
            fh.write(page)

        spent += cost or 0
        ok = all(checks.values())
        failed += not ok
        marks = "  ".join(("✓ " if v else "✗ ") + k for k, v in checks.items())
        print(f"{'PASS' if ok else 'FAIL'}  {p['id']:<26} ${cost or 0:.2f}  {marks}")
        if not checks["guard accepts"] and why:
            print(f"      guard: {why[:160]}")

    os.remove(config)
    print(f"\n{failed} failed · ${spent:.2f} · transcripts in {outdir}")
    sys.exit(1 if failed else 0)


if __name__ == "__main__":
    main()
