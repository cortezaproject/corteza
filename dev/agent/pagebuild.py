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
     "icon": "font-awesome://rocket",           // optional nav icon
     "module": "<module-handle>",               // makes it a record page
     "blocks": [
       {"kind": "RecordList", "title": "...", "xywh": [x,y,w,h],
        "options": {"module": "<module-handle>"}}
     ],
     "children": [ ...same shape, created with selfID=parent... ]}
  ]
}

Anywhere in chart configs / block options, {"module": "<handle>"} and
{"chart": "<handle>"} are replaced with resolved {"moduleID"}/{"chartID"}.
Unknown handles are an error. Every block must have xywh (48-col grid,
cell height 10px; blocks CLIP when too short — Metric needs h>=20,
RecordList/Chart h>=30).
"""

import json
import os
import subprocess
import sys
import urllib.parse
import urllib.request

AGENT_DIR = os.path.dirname(os.path.abspath(__file__))
API = os.environ.get("HUMAN_API", "http://localhost:1043/api")

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
        payload["config"].setdefault("navItem", {})["icon"] = {
            "type": "library",
            "src": spec["icon"],
        }

    if spec["handle"] in existing_pages:
        pid = existing_pages[spec["handle"]]["pageID"]
        payload["updatedAt"] = existing_pages[spec["handle"]].get("updatedAt")
        api("POST", f"/compose/namespace/{nsid}/page/{pid}", payload)
        print(f"page {spec['handle']} updated (ID {pid})")
    else:
        pid = api("POST", f"/compose/namespace/{nsid}/page/", payload)["pageID"]
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
