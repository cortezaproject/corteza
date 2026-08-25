#!/usr/bin/env python3
"""Build compose charts + pages on the local dev server from a JSON spec.

Usage: pagebuild.py <namespace-slug> <spec.json>

Why this exists: envoy YAML import cannot produce working pages — block
option refs (module/chart handles) are never resolved to IDs (the encoder's
ref labels don't match the SetValue paths), so presentation must be built
through REST where real IDs are known. This tool is the standard way to
create/refresh pages and charts, both for fixtures (dev/agent/seed.sh) and
for system-design tasks (/sys-design).

Spec format:
{
  "charts": [
    {"handle": "...", "name": "...", "config": {"reports": [{"module": "<module-handle>", ...}]}}
  ],
  "pages": [
    {"handle": "...", "title": "...", "visible": true, "weight": 0,
     "module": "<module-handle>",               // makes it a record page
     "blocks": [
       {"kind": "RecordList", "title": "...", "xywh": [x,y,w,h],
        "options": {"module": "<module-handle>"}}
     ],
     "children": [ ...same shape, created with selfID=parent... ]}
  ]
}

Geometry is written to the page AND to its primary layout: the layout is what
the webapp renders from, and it is not refreshed by a page write, so without
that second step an edited xywh silently applies to nothing. The same goes for
a block added to a page that already exists — the page write stores it and the
layout would never learn of it, so it is appended to the layout here too.

Anywhere in chart configs / block options, {"module": "<handle>"} and
{"chart": "<handle>"} are replaced with resolved {"moduleID"}/{"chartID"}.
Unknown handles are an error. Every block must have xywh (48-col grid,
cell height 10px; blocks CLIP when too short — Metric needs h>=20,
RecordList/Chart h>=30). A per-row estimate only holds for single-line
cells: a multi-value field renders one line per value, so those rows run
several times taller and a list sized by estimate hides most of them behind
an inner scrollbar. Size lists from what verify-ui.mjs measures, not from
the estimate. Blocks sit side by side by stepping x and keeping
y — a row of four tiles is [0,0,12,20], [12,0,12,20], [24,0,12,20],
[36,0,12,20] — so lay a page out as rows and give the full 48 only to
lists and forms. Pages get no nav icon: the webapp draws one as an image,
and only a human uploading one in the page editor produces a URL for it.
"""

import json
import os
import subprocess
import sys
import urllib.parse
import urllib.request

import stack

AGENT_DIR = os.path.dirname(os.path.abspath(__file__))
# Resolved from this checkout, so a worktree talks to its own server rather
# than the primary's; an exported HUMAN_API still wins. See stack.py.
API = stack.api()

if not urllib.parse.urlparse(API).hostname in ("localhost", "127.0.0.1", "::1"):
    sys.exit(f"pagebuild is local-only; refusing to touch {API}")

TOKEN = subprocess.run(
    [os.path.join(AGENT_DIR, "token.sh")], capture_output=True, text=True, check=True
).stdout.strip()


def api(method, path, body=None):
    req = urllib.request.Request(
        API + path,
        method=method,
        data=None if body is None else json.dumps(body).encode(),
        headers={
            "Authorization": f"Bearer {TOKEN}",
            "Content-Type": "application/json",
        },
    )
    with urllib.request.urlopen(req, timeout=30) as r:
        payload = json.load(r)
    if "error" in payload:
        raise RuntimeError(f"{method} {path}: {payload['error'].get('message')}")
    return payload.get("response")


def fetch_handle_map(path, id_key):
    out = {}
    for item in api("GET", path + "?limit=500")["set"]:
        if item.get("handle"):
            out[item["handle"]] = item[id_key]
    return out


def resolve_refs(obj, modules, charts):
    """Recursively swap {"module": handle} -> {"moduleID": id} and
    {"chart": handle} -> {"chartID": id}."""
    if isinstance(obj, list):
        return [resolve_refs(v, modules, charts) for v in obj]
    if not isinstance(obj, dict):
        return obj
    out = {}
    for k, v in obj.items():
        if k == "module" and isinstance(v, str):
            if v not in modules:
                raise RuntimeError(f"unknown module handle: {v}")
            out["moduleID"] = modules[v]
        elif k == "chart" and isinstance(v, str):
            if v not in charts:
                raise RuntimeError(f"unknown chart handle: {v}")
            out["chartID"] = charts[v]
        else:
            out[k] = resolve_refs(v, modules, charts)
    return out


def upsert_chart(nsid, spec, modules):
    existing = {
        c["handle"]: c
        for c in api("GET", f"/compose/namespace/{nsid}/chart/?limit=500")["set"]
        if c.get("handle")
    }
    payload = {
        "handle": spec["handle"],
        "name": spec.get("name", spec["handle"]),
        "config": resolve_refs(spec["config"], modules, {}),
    }
    if spec["handle"] in existing:
        cid = existing[spec["handle"]]["chartID"]
        payload["updatedAt"] = existing[spec["handle"]].get("updatedAt")
        api("POST", f"/compose/namespace/{nsid}/chart/{cid}", payload)
        print(f"chart {spec['handle']} updated (ID {cid})")
    else:
        cid = api("POST", f"/compose/namespace/{nsid}/chart/", payload)["chartID"]
        print(f"chart {spec['handle']} created (ID {cid})")
    return cid


def build_blocks(spec_blocks, modules, charts):
    blocks = []
    for i, b in enumerate(spec_blocks):
        if "xywh" not in b or len(b["xywh"]) != 4:
            raise RuntimeError(f"block '{b.get('title', b.get('kind'))}' needs xywh [x,y,w,h]")
        blocks.append(
            {
                "blockID": str(i + 1),
                "kind": b["kind"],
                "title": b.get("title", ""),
                "xywh": b["xywh"],
                "style": b.get("style", {}),
                "options": resolve_refs(b.get("options", {}), modules, charts),
            }
        )
    return blocks


def sync_layout(nsid, pid, blocks):
    """Push the spec's geometry onto the page's primary layout.

    The layout — not page.blocks — is what the webapp renders from. On the
    first create the server derives one from the blocks, so a fresh page looks
    right and the two agree; every later write leaves the layout untouched, so
    an edited xywh applies to nothing and pagebuild still reports success.
    Blocks are matched by blockID, which build_blocks assigns by position.

    A block the spec adds to a page that already exists is the same trap one
    step further along: the page write stores it, the layout knows nothing of
    it, and the webapp renders every block but that one. So blocks missing from
    the layout are appended here rather than only re-placed.
    """
    layouts = [
        l
        for l in api("GET", f"/compose/namespace/{nsid}/page-layout?limit=500")["set"]
        if l["pageID"] == pid
    ]
    if not layouts:
        return None, None
    # The server derives a layout handled "primary" but leaves the primary flag
    # false, so neither signal alone finds it. Prefer the flag, then the handle,
    # then the lowest weight — taking the first hit picks an arbitrary layout
    # once a page has more than one.
    layout = sorted(
        layouts,
        key=lambda l: (
            not l.get("primary"),
            l.get("handle") != "primary",
            l.get("weight", 0),
        ),
    )[0]
    want = {b["blockID"]: b["xywh"] for b in blocks}
    lblocks = layout.get("blocks") or []
    changed = 0
    for lb in lblocks:
        xywh = want.get(lb.get("blockID"))
        if xywh and lb.get("xywh") != xywh:
            lb["xywh"] = xywh
            changed += 1

    have = {lb.get("blockID") for lb in lblocks}
    added = [
        {"blockID": b["blockID"], "xywh": b["xywh"]}
        for b in blocks
        if b["blockID"] not in have
    ]
    lblocks.extend(added)
    layout["blocks"] = lblocks
    changed += len(added)

    if not changed:
        return 0, 0
    body = {
        k: layout[k]
        for k in (
            "pageLayoutID",
            "pageID",
            "namespaceID",
            "handle",
            "blocks",
            "config",
            "meta",
            "weight",
            "primary",
            "ownedBy",
        )
        if k in layout
    }
    body["updatedAt"] = layout.get("updatedAt")
    api(
        "POST",
        f"/compose/namespace/{nsid}/page/{pid}/layout/{layout['pageLayoutID']}",
        body,
    )
    return changed - len(added), len(added)


def upsert_page(nsid, spec, modules, charts, existing_pages, self_id="0"):
    payload = {
        "selfID": self_id,
        "title": spec.get("title", spec["handle"]),
        "handle": spec["handle"],
        "visible": spec.get("visible", True),
        "weight": spec.get("weight", 0),
        "moduleID": modules[spec["module"]] if "module" in spec else "0",
        "blocks": build_blocks(spec.get("blocks", []), modules, charts),
        "config": spec.get("config", {}),
        "meta": spec.get("meta", {}),
    }
    if "icon" in spec:
        raise RuntimeError(
            f"page '{spec['handle']}': drop the \"icon\" key. It used to write a "
            '{"type":"library","src":"font-awesome://..."} icon, which the webapp '
            "renders as a broken image — it draws nav icons as <img>, and a library "
            "icon has no URL. Icons are set by a person in the page editor, which "
            "uploads the image first."
        )

    if spec["handle"] in existing_pages:
        pid = existing_pages[spec["handle"]]["pageID"]
        payload["updatedAt"] = existing_pages[spec["handle"]].get("updatedAt")
        api("POST", f"/compose/namespace/{nsid}/page/{pid}", payload)
        moved, added = sync_layout(nsid, pid, payload["blocks"])
        bits = []
        if moved:
            bits.append(f"{moved} block(s) re-placed")
        if added:
            bits.append(f"{added} block(s) added to layout")
        note = f", {', '.join(bits)}" if bits else ""
        print(f"page {spec['handle']} updated (ID {pid}){note}")
    else:
        pid = api("POST", f"/compose/namespace/{nsid}/page/", payload)["pageID"]
        sync_layout(nsid, pid, payload["blocks"])
        print(f"page {spec['handle']} created (ID {pid})")

    for child in spec.get("children", []):
        upsert_page(nsid, child, modules, charts, existing_pages, self_id=pid)
    return pid


def main():
    if len(sys.argv) != 3:
        sys.exit(__doc__)
    slug, spec_path = sys.argv[1], sys.argv[2]
    spec = json.load(open(spec_path))

    ns = api("GET", f"/compose/namespace/?slug={slug}&limit=1")["set"]
    if not ns:
        sys.exit(f"namespace not found: {slug}")
    nsid = ns[0]["namespaceID"]

    modules = fetch_handle_map(f"/compose/namespace/{nsid}/module/", "moduleID")

    for c in spec.get("charts", []):
        upsert_chart(nsid, c, modules)

    charts = fetch_handle_map(f"/compose/namespace/{nsid}/chart/", "chartID")
    existing_pages = {
        p["handle"]: p
        for p in api("GET", f"/compose/namespace/{nsid}/page/?limit=500")["set"]
        if p.get("handle")
    }

    for p in spec.get("pages", []):
        upsert_page(nsid, p, modules, charts, existing_pages)


if __name__ == "__main__":
    main()
