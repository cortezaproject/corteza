# Docs screenshots

The product screenshots in `docs/public/screenshots/` are generated, never taken
by hand. Re-run the pipeline and every image is refreshed from the same data.

## Prerequisites

A running local dev server from the agent toolkit (`dev/agent/`): the API and
the webapp of this checkout, e.g. a worktree started with
`dev/agent/worktree.sh up <name>`. Both scripts refuse anything that is not
`localhost`.

## Run

```sh
docs/tools/screenshots/setup.py        # seed data, users, agent, TAQ, chatbot (idempotent)
node docs/tools/screenshots/shoot.mjs  # capture every scene, light and dark
```

`setup.py` imports the `dev/fixtures/docs-showcase` fixture (a small Sales
namespace), builds its pages and charts with `dev/agent/pagebuild.py`, and
creates a screenshot user (`docs@example.com`, password in
`dev/agent/.state/docs-password`) together with an LLM provider, an agent, a
TAQ and a chatbot. A second run changes nothing it does not have to.

`shoot.mjs` logs in as that user and writes
`docs/public/screenshots/<scene>-<light|dark>.webp` at 1440×900, 2× scale. It
switches the user's theme between passes and leaves it on light.

- `--only home,taq-builder` captures just those scenes.
- `--png` keeps the lossless capture next to each WebP (don't commit those).

Look at every image before committing: a spinner, an empty list or a toast in
a screenshot means a scene's `ready` selector needs tightening.

## Add a scene

Add an entry to `SCENES` in `shoot.mjs`: a `name`, the webapp `path`, a `ready`
selector that only matches once the page's content has rendered, and
optionally `prepare(page)` to open a panel or scroll first. Then use it in a
page:

```md
<Screenshot name="my-scene" alt="What the image shows" />
```

## Clean up

`setup.py` records everything it creates in the toolkit's ledger under the
session `docs-screenshots`. `dev/agent/cleanup.sh --session docs-screenshots`
deletes the Sales namespace from it, but `cleanup.sh` only deletes namespaces:
the screenshot user, the LLM provider, the agent, the TAQ and the chatbot stay.
Run the pipeline in a worktree when you want nothing left behind —
`worktree.sh rm` drops its whole database.
