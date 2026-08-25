"""Where this checkout's dev stack listens, for the python side of the toolkit.

One resolver, shared with the bash scripts and stack.mjs: stack.sh reads
server/.env and .env.e2e, so a tool run inside a worktree talks to that
worktree's server rather than the primary's. An exported HUMAN_API still wins —
stack.sh inherits this process's environment.

Defaulting to port 1043 instead is not a harmless fallback. Every other script
here resolves, so `token.sh` mints against the worktree while a hardcoded
default sends the call to the primary: the token's row is in the worktree's
database and the primary's lookup cannot find it, which surfaces as an
unauthorised MCP call that looks like broken auth. Where the token does
validate — anything reached with the primary's own token — the call succeeds
against the WRONG SERVER, and a worktree session reads and writes the primary's
data while believing it is isolated.
"""

import os
import subprocess

HERE = os.path.dirname(os.path.abspath(__file__))


def stack():
    """Every URL stack.sh resolves, as a dict."""
    out = subprocess.run(
        ["bash", os.path.join(HERE, "stack.sh")],
        capture_output=True,
        text=True,
        check=True,
    ).stdout
    env = {}
    for line in out.splitlines():
        if "=" in line:
            k, v = line.split("=", 1)
            env[k] = v
    return env


def api():
    """This checkout's API base, e.g. http://localhost:1143/api in slot 1."""
    return stack()["HUMAN_API"]
