import { NoID } from '@planetcrust/human-js'

// Where a carried page lands, from the rows on screen and the pointer.
//
// A row is { key, parentKey, depth, top, bottom } in document order;
// `top`/`bottom` and the pointer share one coordinate space. Rows of the
// carried page's own subtree are left out by the caller: a page cannot land
// inside itself.
//
// On a row's middle band the page goes into that row, first among its
// sub-pages. Between rows it goes after the last row whose middle the pointer
// has passed (P), or before the first row, at the depth the pointer's x asks
// for, clamped to what the tree allows there: no deeper than one level under
// P, and no shallower than the row below, which would otherwise be cut off
// from its parent. The line is drawn `gap / 2` below the row it follows.
export function projectDrop(rows, x, y, { indent, baseX, gap = 12 }) {
  if (!rows.length) return { parentKey: NoID, afterKey: null, depth: 0, lineY: 0, intoKey: null }

  const onto = rows.find(r => {
    const band = (r.bottom - r.top) / 4
    return y >= r.top + band && y <= r.bottom - band
  })
  if (onto) {
    return {
      parentKey: onto.key,
      afterKey: null,
      depth: onto.depth + 1,
      lineY: onto.bottom + gap / 2,
      intoKey: onto.key,
    }
  }

  let p = -1
  for (let i = 0; i < rows.length; i++) {
    if (y >= (rows[i].top + rows[i].bottom) / 2) p = i
    else break
  }
  const P = rows[p]
  const next = rows[p + 1]

  if (!P) {
    return {
      parentKey: NoID,
      afterKey: null,
      depth: 0,
      lineY: rows[0].top - gap / 2,
      intoKey: null,
    }
  }

  const minDepth = next ? next.depth : 0
  const maxDepth = P.depth + 1
  const wanted = Math.round((x - baseX) / indent)
  const depth = Math.max(minDepth, Math.min(maxDepth, wanted))
  const lineY = P.bottom + gap / 2

  if (depth === P.depth + 1)
    return { parentKey: P.key, afterKey: null, depth, lineY, intoKey: null }

  // The sibling the page follows: P's ancestor-or-self at that depth
  let i = p
  while (i >= 0 && rows[i].depth > depth) i--
  const S = rows[i]
  return { parentKey: S.parentKey, afterKey: S.key, depth, lineY, intoKey: null }
}

// The saved shape of a landing: the parent and the order of its sub-pages
// with the carried page in place. `children(parentKey)` returns that parent's
// current sub-page keys.
export function dropPlan(spot, draggedKey, children) {
  const siblings = children(spot.parentKey).filter(k => k !== draggedKey)
  const index = spot.afterKey ? siblings.indexOf(spot.afterKey) + 1 : 0
  siblings.splice(index, 0, draggedKey)
  return { parentKey: spot.parentKey, pageIDs: siblings, index }
}
