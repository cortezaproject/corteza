// Synthetic Record blocks for the admin record screens.
//
// The screens have no page behind them, so they hand Grid a block list built
// here: the module's own fields in one card, the record's system fields in a
// second card below it.

// Grid geometry, from Grid.vue: 48 columns, 10px cells, and a 6px margin that
// insets an item's content on every side.
const GRID_COLUMNS = 48
const CELL_HEIGHT = 10
const ITEM_MARGIN = 6

// What a Record block draws per field: a label above the value, a gap to the
// next field, and the container's padding at each end. An editor's value box is
// taller than a viewer's, and a block is sized for the taller of the two —
// gridstack keeps the height it was given, so a screen that switches between
// viewing and editing gets one height for both.
const VIEWER_ROW = 53
const EDITOR_ROW = 63
const ROW_GAP = 20
const CONTENT_PADDING = 32

// The Card header a titled block draws above its fields.
const CARD_HEADER = 45

// The one system field a record's owner can change; the rest are read-only
// wherever they are drawn.
const EDITABLE_SYSTEM_FIELD = 'ownedBy'

// The server assigns these when it stores the record. On a record that does not
// exist yet they hold the unset sentinel, which draws as a literal 0.
const ASSIGNED_ON_SAVE = ['recordID', 'revision']

// Cell height that fits the given row heights without the block scrolling.
function cellsForRows(rows, titled) {
  const content =
    CONTENT_PADDING +
    rows.reduce((sum, h) => sum + h, 0) +
    Math.max(rows.length - 1, 0) * ROW_GAP +
    (titled ? CARD_HEADER : 0)

  return Math.ceil((content + 2 * ITEM_MARGIN) / CELL_HEIGHT)
}

/**
 * The two stacked blocks an admin record screen renders.
 *
 * A second block is what puts Grid into gridstack layout, where a block is as
 * tall as its xywh says rather than as tall as the viewport — hence the
 * computed heights. A field that draws taller than a plain input (a file list,
 * a map, several values) scrolls inside its block.
 *
 * @param {object} recordModule the module whose record is on screen
 * @param {string} opts.idPrefix distinguishes the screen's blocks in the grid
 * @param {string} opts.systemTitle heading for the system-fields card
 * @param {boolean} opts.isNew the record is not stored yet
 */
export function adminRecordBlocks(recordModule, { idPrefix, systemTitle, isNew = false }) {
  const ownFields = recordModule?.fields || []
  const systemFields = (recordModule?.systemFields?.() || []).filter(
    f => !(isNew && ASSIGNED_ON_SAVE.includes(f.name)),
  )

  const ownHeight = cellsForRows(
    ownFields.map(() => EDITOR_ROW),
    false,
  )
  const systemHeight = cellsForRows(
    systemFields.map(f => (f.name === EDITABLE_SYSTEM_FIELD ? EDITOR_ROW : VIEWER_ROW)),
    true,
  )

  return [
    block({
      id: `${idPrefix}_fields`,
      title: '',
      fields: [],
      xywh: [0, 0, GRID_COLUMNS, ownHeight],
    }),
    block({
      id: `${idPrefix}_system`,
      title: systemTitle,
      fields: systemFields.map(f => f.name),
      xywh: [0, ownHeight, GRID_COLUMNS, systemHeight],
    }),
  ]
}

function block({ id, title, fields, xywh }) {
  return {
    blockID: id,
    kind: 'Record',
    title,
    description: '',
    style: { wrap: { kind: 'card' } },
    options: { fields },
    xywh,
    meta: { tempID: id },
  }
}
