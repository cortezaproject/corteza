---
name: page_layout
description: A page holds what its blocks are; a layout holds where they go — and how the grid places a block you do not position.
triggers:
  - compose_page_create
  - compose_page_update
  - compose_page_remove_blocks
  - compose_page_layout_lookup
  - compose_page_layout_create
  - compose_page_layout_update
---

# Pages and layouts

They are two objects and the difference decides whether anything renders.

- A **page** holds what its blocks ARE: `kind`, `title`, `options`.
- A **layout** holds where they go, as `{blockID, xywh}`.

The viewer draws the layout. A block the layout does not name is never
rendered, however complete it looks on the page.

Every page is created with a `primary` layout seeded from the blocks it was
created with, so a fresh page needs nothing else.

## Adding a block to an existing page

`compose_page_update` merges blocks by `blockID` and places the new ones for
you. Which layout they land on:

- **One layout on the page** — it is used, and `layout` can be omitted.
- **Several layouts** — pass `layout` (its handle or ID). Omitting it is
  refused rather than guessed: each layout is a deliberately different subset of
  the page, so a guess would put the block somewhere you did not ask for and
  hide it everywhere else.

Blocks already on the page are never re-placed, so changing one block's
`options` does not move it.

`compose_page_remove_blocks` drops the block from every layout as well as from
the page. That one needs no `layout` argument — a block the page no longer has
cannot be placed anywhere.

## The grid

48 columns wide; a cell is 10px tall. Full width is `w=48`, half `24`, a
quarter `12`. Blocks sit side by side by stepping `x` and keeping `y` equal.

**`xywh` is honoured.** A block that carries a width is left exactly where you
put it — this is how a dashboard is built. Only the coordinate the grid cannot
render is normalised: a negative `x`/`y` is clamped to 0, a width above 48 to
48, and a missing height filled in per kind.

Omit `xywh` and the block is auto-placed instead: it flows left to right at its
kind's default width, wraps when the row fills, and starts below the lowest
block you positioned yourself.

| Kind                                                                   | Default width | Default height |
| ---------------------------------------------------------------------- | ------------- | -------------- |
| Metric, Progress                                                       | 12 (quarter)  | 20             |
| Content, Automation                                                    | 24 (half)     | 20             |
| Chart, Calendar, Comment, SocialFeed                                   | 24 (half)     | 30             |
| RecordList, Record, RecordOrganizer, ChatbotInbox and every other kind | 48 (full)     | 30             |

Blocks CLIP when too short: give Metric `h>=20` and RecordList or Chart
`h>=30`.

## Reading before writing

`compose_page_lookup` reports the blocks and their blockIDs;
`compose_page_layout_lookup` reports which of them each layout places.
blockIDs are assigned per page, so they mean nothing on another page.
