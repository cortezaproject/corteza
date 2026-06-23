---
name: page_layout
description: How blocks are laid out on a page — the server places them automatically based on order.
importance: high
triggers:
  - compose_page_create
  - compose_page_update
---

# Page Layout

The server automatically positions blocks based on their order in the array. You do not need to specify `xywh`. Just send blocks in the order you want them to appear.

## How the server places blocks

Blocks are placed in a 2-column grid (each column is half the page width). The first block goes left, the second goes right, the third starts a new row on the left, the fourth goes right on that row, and so on.

```
[block 0]  [block 1]
[block 2]  [block 3]
[block 4]  [block 5]
```

Height is set automatically per block kind. Row height is determined by the taller of the two blocks in that row.

## Order matters

The order of blocks in your array determines their visual position. Put the most important content first. If you want something at the top left, put it first in the array.

## Adding blocks to an existing page

When updating a page with new blocks (no `blockID`), they are appended after existing blocks. The layout of existing blocks is not changed. Look up the page first with `compose_page_lookup` if you need to see what is already there.

## Do not specify xywh

Do not set `xywh` — the server ignores it and computes positions from block order. Just provide `kind`, `title`, and `options`.
