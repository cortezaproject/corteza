<!-- Template: folder-named intent doc. Delete comments. Cap: 60 prose lines. -->
---
kind: folder
covers: recursive # or "." if children carry their own docs
owner: fe # fe | be | shared
depends-on: [] # repo-relative paths this area relies on
touched-by: [] # known consumers of this area
tests: [] # specs that must pass when this area changes
---

# <Area name>

## Intention

<!-- Why this folder exists. What user/system need it serves. 2–5 lines. -->

## Map

<!-- One line per child (folder or key file): name — role. -->

## Data touched

<!-- Stores, APIs, resources this area reads/writes. Omit section if none. -->

## When changing this

<!-- Invariants to preserve, traps, what must be re-tested. -->
