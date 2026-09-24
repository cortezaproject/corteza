#!/usr/bin/env python3
"""Prepare the local dev server for the docs screenshots.

Usage: docs/tools/screenshots/setup.py

Idempotent: every resource is looked up by its handle, email or name first and
created only when missing, so a second run creates nothing. Everything goes
through the agent toolkit (dev/agent/api.sh, mcp.py, pagebuild.py, the server
CLI import) and is recorded in dev/agent/.state/created.jsonl under the session
`docs-screenshots` (override with DOCS_SHOTS_SESSION); `dev/agent/cleanup.sh --session docs-screenshots`
removes its namespace (see README.md for the rest).

Creates:
  - the Sales namespace from dev/fixtures/docs-showcase (records, charts, pages)
  - users Alex Morgan (docs@example.com, the screenshot login), Priya Shah and
    Daniel Okafor, and the role "Sales team" holding all three; Alex is also a
    super-admin member and gets a password (dev/agent/.state/docs-password)
  - an Anthropic LLM provider with a placeholder key, the "Sales assistant"
    agent, the "Notify on large deals" TAQ and the "Website chat" chatbot
  - two deals created after the TAQ, so its notifications reach Alex's inbox
"""

import csv
import json
import os
import secrets
import shutil
import subprocess
import sys
import time

HERE = os.path.dirname(os.path.abspath(__file__))
REPO = os.path.abspath(os.path.join(HERE, "..", "..", ".."))
AGENT = os.path.join(REPO, "dev", "agent")
FIXTURE = os.path.join(REPO, "dev", "fixtures", "docs-showcase")
STATE = os.path.join(AGENT, ".state")
LEDGER = os.path.join(STATE, "created.jsonl")
PASSWORD_FILE = os.path.join(STATE, "docs-password")

SESSION = os.environ.get("DOCS_SHOTS_SESSION", "docs-screenshots")
ENV = dict(os.environ, CLAUDE_CODE_SESSION_ID=SESSION)

SLUG = "sales"
SHOT_EMAIL = "docs@example.com"
USERS = [
    ("Alex Morgan", SHOT_EMAIL, "alex_morgan"),
    ("Priya Shah", "priya.shah@example.com", "priya_shah"),
    ("Daniel Okafor", "daniel.okafor@example.com", "daniel_okafor"),
]
ROLE_HANDLE = "sales_team"
PROVIDER_HANDLE = "docs_anthropic"
AGENT_HANDLE = "sales_assistant"
TAQ_HANDLE = "notify_large_deals"
CHATBOT_HANDLE = "website_chat"
MODEL = "claude-sonnet-4-5"

# Deals created once the TAQ listens, so its notifications are real ones.
LIVE_DEALS = [
    {"name": "Warehouse robotics pilot", "company": "Tidewater Robotics",
     "contact": "Felix Brandt", "stage": "qualification", "value": "18500",
     "close_date": "2027-01-15"},
    {"name": "Store loyalty app", "company": "Copperleaf Bakeries",
     "contact": "Camille Rousseau", "stage": "proposal", "value": "42000",
     "close_date": "2026-12-09"},
]

created = 0


def say(msg):
    print(f"  {msg}", flush=True)


def ledger(kind, rid, slug=""):
    global created
    created += 1
    with open(LEDGER, "a") as fh:
        fh.write(json.dumps({
            "kind": kind, "id": str(rid), "slug": slug, "session": SESSION,
            "ts": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
        }) + "\n")
    say(f"created {kind} {slug or rid} (ID {rid})")


def api(method, path, body=None, retry=6):
    """dev/agent/api.sh --json; returns response, raises on an API error."""
    cmd = [os.path.join(AGENT, "api.sh"), "--json", method, path]
    if body is not None:
        cmd += ["-d", json.dumps(body)]
    for _ in range(retry):
        p = subprocess.run(cmd, capture_output=True, text=True, env=ENV)
        if not p.stdout.strip():
            time.sleep(5)
            continue
        out = json.loads(p.stdout)
        err = out.get("error") if isinstance(out, dict) else None
        if err and err.get("kind") == "restarting":
            time.sleep(5)
            continue
        if err:
            raise RuntimeError(f"{method} {path}: {err.get('message')}")
        return out.get("response", out)
    raise RuntimeError(f"{method} {path}: server kept restarting")


def mcp(tool, args):
    """dev/agent/mcp.py call; returns the parsed payload."""
    for _ in range(6):
        p = subprocess.run(
            ["python3", os.path.join(AGENT, "mcp.py"), "call", tool, json.dumps(args)],
            capture_output=True, text=True, env=ENV,
        )
        try:
            out = json.loads(p.stdout)
        except json.JSONDecodeError:
            if p.returncode == 0:
                return {"text": p.stdout.strip()}
            raise RuntimeError(f"{tool}: {p.stdout.strip()} {p.stderr.strip()}")
        if isinstance(out, dict) and "error" in out and out["error"].get("kind") in ("http", "unreachable"):
            time.sleep(5)
            continue
        if p.returncode != 0:
            raise RuntimeError(f"{tool}: {json.dumps(out)[:600]}")
        return out
    raise RuntimeError(f"{tool}: server unreachable")


def cli_import(stage):
    subprocess.run(
        ["bash", "-c", f'source "{AGENT}/common.sh"; server_cli import --skip-existing "$1"', "_", stage],
        check=True, env=ENV,
    )


# ── Users and the Sales team role ─────────────────────────────────────────


def ensure_users():
    ids = {}
    for name, email, handle in USERS:
        found = api("GET", f"/system/users/?email={email}")["set"]
        if found:
            ids[email] = found[0]["userID"]
            continue
        u = mcp("system_user_create", {"email": email, "name": name, "handle": handle})
        uid = u.get("userID") or u.get("user", {}).get("userID")
        ids[email] = uid
        ledger("user", uid, email)
        mcp("system_user_set_email_confirmed", {"user": uid, "confirmed": True})
    return ids


def ensure_role(user_ids):
    roles = api("GET", "/system/roles/?limit=500")["set"]
    role = next((r for r in roles if r["handle"] == ROLE_HANDLE), None)
    if role:
        rid = role["roleID"]
    else:
        r = mcp("system_role_create", {
            "handle": ROLE_HANDLE, "name": "Sales team",
            "description": "Account executives working the Sales namespace; sees the Sales assistant.",
        })
        rid = r.get("roleID") or r.get("role", {}).get("roleID")
        ledger("role", rid, ROLE_HANDLE)
    members = set(api("GET", f"/system/roles/{rid}/members"))
    for uid in user_ids.values():
        if uid not in members:
            mcp("system_role_member_add", {"role": rid, "user": uid})
    # A few system rules, so the role's column in the permission grid is not blank.
    want = [("users.search", "allow"), ("roles.search", "allow"), ("action-log.read", "allow"),
            ("role.create", "deny"), ("user.create", "deny")]
    have = {(r["operation"], r["access"]) for r in
            mcp("system_permission_lookup", {"role": rid}).get("rules") or []
            if r.get("resource") == "corteza::system/"}
    if not set(want) <= have:
        mcp("system_permission_grant", {"role": rid, "rules": json.dumps([
            {"resource": "corteza::system/", "operation": op, "access": acc} for op, acc in want])})
        say("granted the Sales team its system rules")
    return rid


def ensure_login(uid):
    admins = api("GET", "/system/roles/?limit=500")["set"]
    sa = next(r["roleID"] for r in admins if r["handle"] == "super-admin")
    if uid not in set(api("GET", f"/system/roles/{sa}/members")):
        mcp("system_role_member_add", {"role": sa, "user": uid})
    if not os.path.exists(PASSWORD_FILE):
        with open(PASSWORD_FILE, "w") as fh:
            fh.write(f"Docs-{secrets.token_hex(8)}-Shots9!")
    pw = open(PASSWORD_FILE).read().strip()
    api("POST", f"/system/users/{uid}/password", {"password": pw})


# ── The Sales namespace ───────────────────────────────────────────────────


def namespace():
    found = api("GET", f"/compose/namespace/?slug={SLUG}&limit=1")["set"]
    return found[0]["namespaceID"] if found else None


def modules(nsid):
    return {m["handle"]: m for m in api("GET", f"/compose/namespace/{nsid}/module/?limit=100")["set"]}


def ensure_namespace():
    nsid = namespace()
    if nsid:
        mods = modules(nsid)
        deal = mods.get("deal")
        if not deal or not any(f["name"] == "close_date" for f in deal["fields"]):
            sys.exit(f"a namespace '{SLUG}' exists that is not the docs showcase (ID {nsid}); "
                     "remove or rename it first")
    else:
        stage = os.path.join(STATE, "sessions", SESSION, "import-docs-showcase")
        shutil.rmtree(stage, ignore_errors=True)
        os.makedirs(stage)
        for f in os.listdir(FIXTURE):
            if f != "ui.json":
                shutil.copy(os.path.join(FIXTURE, f), stage)
        cli_import(stage)
        nsid = namespace()
        if not nsid:
            sys.exit("import ran but the sales namespace is not there")
        ledger("namespace", nsid, SLUG)
        # The CLI import bypasses the running server's DAL registry; a no-op
        # module update registers each model with it.
        for m in modules(nsid).values():
            api("POST", f"/compose/namespace/{nsid}/module/{m['moduleID']}", m)

    ensure_fields(nsid)
    subprocess.run(["python3", os.path.join(AGENT, "pagebuild.py"), SLUG,
                    os.path.join(FIXTURE, "ui.json")], check=True, env=ENV)
    subprocess.run([os.path.join(AGENT, "ids.sh"), SLUG], check=True, env=ENV,
                   stdout=subprocess.DEVNULL)
    return nsid


def ensure_fields(nsid):
    """Field labels and order as def.yaml lists them.

    The envoy import keeps neither: a module field has no label to decode into,
    and the YAML mapping reaches the store in no particular order.
    """
    import yaml

    spec = yaml.safe_load(open(os.path.join(FIXTURE, "def.yaml")))["namespace"][SLUG]["modules"]
    for handle, m in modules(nsid).items():
        want = [(name, f.get("label", name)) for name, f in spec[handle]["fields"].items()]
        have = {f["name"]: f for f in m["fields"]}
        fields = [dict(have[name], label=label) for name, label in want]
        if [(f["name"], f.get("label")) for f in m["fields"]] == want:
            continue
        api("POST", f"/compose/namespace/{nsid}/module/{m['moduleID']}", dict(m, fields=fields))
        say(f"set field labels and order on module {handle}")


def records(nsid, mid):
    return api("GET", f"/compose/namespace/{nsid}/module/{mid}/record/?limit=200")["set"]


def value(rec, name):
    return next((v["value"] for v in rec.get("values", []) if v["name"] == name), None)


def ensure_owners(nsid, user_ids):
    deal = modules(nsid)["deal"]
    owners = {}
    with open(os.path.join(FIXTURE, "deals.csv")) as fh:
        for row in csv.DictReader(fh):
            owners[row["name"]] = user_ids[row["owner_email"]]
    for rec in records(nsid, deal["moduleID"]):
        want = owners.get(value(rec, "name"))
        if not want or (value(rec, "owner") == want and rec.get("ownedBy") == want):
            continue
        rec["values"] = [v for v in rec["values"] if v["name"] != "owner"] + [{"name": "owner", "value": want}]
        rec["ownedBy"] = want
        api("POST", f"/compose/namespace/{nsid}/module/{deal['moduleID']}/record/{rec['recordID']}", rec)


# ── Assistant, TAQ, chatbot ───────────────────────────────────────────────


def ensure_provider():
    for p in api("GET", "/system/llm-providers/")["set"]:
        if p.get("handle") == PROVIDER_HANDLE:
            return p["llmProviderID"]
    p = api("POST", "/system/llm-providers/", {
        "handle": PROVIDER_HANDLE, "provider": "anthropic", "status": "active",
        "apiKey": "sk-ant-docs-placeholder",
        "meta": {"short": "Anthropic", "description": "Claude models for the Sales assistant"},
        "config": {"model": MODEL, "timeout": "60s"},
    })
    ledger("llm_provider", p["llmProviderID"], PROVIDER_HANDLE)
    return p["llmProviderID"]


def ensure_agent(nsid, provider_id, role_id):
    for a in api("GET", "/system/agents/?limit=0")["set"]:
        if a.get("handle") == AGENT_HANDLE:
            return a["agentID"]
    scope = [{"namespaceID": nsid, "moduleIDs": []}]
    tool = lambda name, perm, desc: {"name": name, "permission": perm, "description": desc, "allow": scope}
    a = mcp("system_agent_create", {
        "handle": AGENT_HANDLE,
        "meta": json.dumps({
            "short": "Sales assistant",
            "description": "Answers questions about companies, contacts and deals, and keeps deal records up to date.",
            "sidebarRoles": [role_id],
        }),
        "behavior": json.dumps({"systemPrompt": (
            "You help the sales team work their pipeline in the Sales namespace. "
            "Look records up before answering, quote deal values in euros, and name the "
            "company and deal owner when you mention a deal. Ask before you change a record, "
            "and say what you changed afterwards. Keep answers short."
        )}),
        "execution": json.dumps({
            "model": {"llmProviderID": provider_id, "model": MODEL},
            "limits": {"maxIterations": 10, "timeout": "2m"},
        }),
        "access": json.dumps({"allow": scope, "tools": [
            tool("compose_record_lookup", "always", "Find companies, contacts and deals"),
            tool("compose_record_report", "always", "Sum and count the pipeline"),
            tool("compose_record_create", "ask", "Log a new deal or contact"),
            tool("compose_record_update", "ask", "Move a deal to another stage"),
        ]}),
        "invocation": json.dumps({"user": {"enabled": True}, "system": {"enabled": False}}),
    })
    aid = a.get("agentID") or a.get("agent", {}).get("agentID")
    ledger("agent", aid, AGENT_HANDLE)
    return aid


def ensure_taq(nsid, deal_mid):
    found = mcp("automation_taq_lookup", {"query": TAQ_HANDLE}).get("taqs") or []
    hit = next((t for t in found if t.get("handle") == TAQ_HANDLE), None)
    if hit:
        return hit["automationID"]
    trigger = "onDealCreated"
    t = mcp("automation_taq_create", {
        "handle": TAQ_HANDLE,
        "name": "Notify on large deals",
        "description": "Tells the deal owner when a deal worth more than €10,000 is created.",
        "enabled": True,
        "triggers": json.dumps([{
            "triggerID": "1001", "handle": trigger, "enabled": True,
            "resourceType": "compose:record", "eventType": "afterCreate",
            "meta": {"short": "After Record Create"},
            "constraints": [
                {"name": "namespace", "op": "=", "values": [{"@type": "ID", "@value": nsid}]},
                {"name": "module", "op": "=", "values": [{"@type": "ID", "@value": deal_mid}]},
            ],
        }]),
        "steps": json.dumps([
            {"stepID": "1", "handle": "isLarge", "kind": "gatewayExclusive",
             "meta": {"short": "First Match"}},
            {"stepID": "2", "handle": "notifyOwner", "kind": "function", "ref": "notificationSendRecord",
             "meta": {"short": "Send Record Notification"},
             "arguments": [
                 {"argumentName": "recipient", "expr": "record.ownedBy", "type": "ID"},
                 {"argumentName": "title", "value": "Large deal created", "type": "String"},
                 {"argumentName": "description", "expr": "record.values.name + \" is worth more than €10,000\"",
                  "type": "String"},
                 {"argumentName": "namespace", "value": nsid, "type": "ID"},
                 {"argumentName": "module", "value": deal_mid, "type": "ID"},
                 {"argumentName": "record", "expr": "record.recordID", "type": "ID"},
             ]},
            {"stepID": "3", "handle": "smallDeal", "kind": "termination",
             "meta": {"short": "Small deal, nothing to send"}},
        ]),
        "paths": json.dumps([
            {"parentID": "1001", "childID": "1"},
            {"parentID": "1", "childID": "2", "condition": {
                "ref": "gt",
                "args": [{"symbol": "record.values.value", "meta": {"scope": trigger}},
                         {"value": {"@type": "Integer", "@value": 10000}}]}},
            {"parentID": "1", "childID": "3"},
        ]),
    })
    if t.get("issues"):
        raise RuntimeError(f"TAQ stored with issues: {json.dumps(t['issues'])}")
    tid = t.get("automationID")
    ledger("taq", tid, TAQ_HANDLE)
    return tid


def ensure_chatbot(agent_id):
    found = api("GET", f"/system/chatbots/?handle={CHATBOT_HANDLE}")["set"] or []
    if found:
        return found[0]["chatbotID"]
    c = mcp("system_chatbot_create", {
        "handle": CHATBOT_HANDLE,
        "name": "Website chat",
        "enabled": True,
        "allowedOrigins": json.dumps(["https://www.example.com"]),
        "sessionTTL": "30m",
        "scenarios": json.dumps([
            {"id": "welcome", "name": "Welcome", "type": "static_message",
             "config": {"message": "Hi! I can answer questions about our plans and pricing.", "autoAdvanceMs": 800}},
            {"id": "details", "name": "Your details", "type": "form",
             "config": {"submitLabel": "Start chat", "fields": [
                 {"name": "name", "label": "Your name", "type": "text", "required": True},
                 {"name": "email", "label": "Work email", "type": "email", "required": True},
                 {"name": "company", "label": "Company", "type": "text", "required": False},
             ]}},
            {"id": "chat", "name": "Chat with sales", "type": "conversation", "agentID": agent_id,
             "config": {"placeholder": "Ask about plans, pricing or a demo…",
                        "initialPrompt": "The visitor filled in the form above. Greet them by name.",
                        "typingIndicator": True}},
        ]),
    })
    cid = c.get("chatbotID") or c.get("chatbot", {}).get("chatbotID")
    ledger("chatbot", cid, CHATBOT_HANDLE)
    return cid


def ensure_live_deals(nsid, alex):
    mods = modules(nsid)
    deal, company, contact = mods["deal"], mods["company"], mods["contact"]
    have = {value(r, "name") for r in records(nsid, deal["moduleID"])}
    by = lambda recs, f: {value(r, f): r["recordID"] for r in recs}
    companies = by(records(nsid, company["moduleID"]), "name")
    contacts = by(records(nsid, contact["moduleID"]), "full_name")
    for d in LIVE_DEALS:
        if d["name"] in have:
            continue
        vals = dict(d, company=companies[d["company"]], contact=contacts[d["contact"]], owner=alex)
        api("POST", f"/compose/namespace/{nsid}/module/{deal['moduleID']}/record/", {
            "ownedBy": alex, "values": [{"name": k, "value": v} for k, v in vals.items()],
        })
        say(f"created deal {d['name']}")


def main():
    print(f"▸ docs screenshots setup (session {SESSION})", flush=True)
    user_ids = ensure_users()
    alex = user_ids[SHOT_EMAIL]
    role_id = ensure_role(user_ids)
    ensure_login(alex)
    nsid = ensure_namespace()
    ensure_owners(nsid, user_ids)
    provider_id = ensure_provider()
    agent_id = ensure_agent(nsid, provider_id, role_id)
    ensure_taq(nsid, modules(nsid)["deal"]["moduleID"])
    ensure_chatbot(agent_id)
    ensure_live_deals(nsid, alex)
    print(f"✓ ready — {created} resource(s) created; log in as {SHOT_EMAIL} "
          f"(password in {os.path.relpath(PASSWORD_FILE, REPO)})")


if __name__ == "__main__":
    main()
