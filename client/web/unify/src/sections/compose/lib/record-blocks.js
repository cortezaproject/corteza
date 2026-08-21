// Synthetic page blocks for the admin record screens.
//
// The screens have no page behind them, so they hand Grid a block list built
// here. The blocks are real PageBlock instances rather than plain objects:
// Metric and RecordRevisions read their data through methods on the class.

import { compose } from '@planetcrust/human-js'

// Grid geometry, from Grid.vue: 48 columns, 10px cells, and a 6px margin that
// insets an item's content on every side.
const GRID_COLUMNS = 48
const CELL_HEIGHT = 10
const ITEM_MARGIN = 6

// What a Record block draws per field: a label above the value. An editor's
// value box is taller than a viewer's, and a block is sized for the taller of
// the two — gridstack keeps the height it was given, so a screen that switches
// between viewing and editing gets one height for both.
const EDITOR_ROW = 63
const CONTENT_PADDING = 32
const CARD_HEADER = 45

// The `wrap` field layout draws fields in responsive columns — four at the
// widths these screens are used at (RecordBlock's `applyColumnClasses`), each
// row carrying its own bottom margin rather than a flex gap. A narrower window
// gets fewer columns and so more rows than are budgeted here, and the card
// scrolls the difference.
const WRAP_COLUMNS = 4
const WRAP_ROW_MARGIN = 16

// A field that needs the width of the page to read at all: a map, an
// attachment list, a paragraph, or anything holding several values.
const WIDE_FIELD_HEIGHT = { Geometry: 320, File: 180 }
const WIDE_FIELD_DEFAULT_HEIGHT = 130

// A metric tile is a label over one number, and a revision list is a table.
// MetricItem sizes its number against the tile, so a short tile draws a
// smaller number rather than clipping it.
const TILE_HEIGHT = 12
const REVISIONS_HEIGHT = 30

// The one system field a record's owner can change; the rest are read-only
// wherever they are drawn.
const EDITABLE_SYSTEM_FIELD = 'ownedBy'

// The server assigns these when it stores the record. On a record that does not
// exist yet they hold the unset sentinel, which draws as a literal 0.
const ASSIGNED_ON_SAVE = ['recordID', 'revision']

// Metric tiles, in the order they sit across the top of the record list. An
// empty filter counts the whole module.
//
// `deleted` is the record report's own state constraint: 0 leaves deleted
// records out, 2 counts nothing else. It is not expressible as a filter — the
// report drops deleted records before a filter is ever applied — which is why
// the last tile carries it rather than a `deletedAt IS NOT NULL`.
const DELETED_ONLY = 2

const LIST_METRICS = [
  { key: 'total', filter: '' },
  { key: 'createdRecently', filter: 'createdAt >= DATE_SUB(NOW(), INTERVAL 7 DAY)' },
  // Interpolated by MetricBlock against the signed-in user
  { key: 'ownedByMe', filter: 'ownedBy = ${userID}' },
  { key: 'deleted', filter: '', deleted: DELETED_ONLY },
]

/** A field the compact multi-column layout cannot hold. */
export function isWideField(field) {
  if (field.isMulti) return true
  if (field.kind === 'File' || field.kind === 'Geometry') return true

  return !!(field.options?.multiLine || field.options?.useRichTextEditor)
}

function cells(contentHeight) {
  return Math.ceil((contentHeight + 2 * ITEM_MARGIN) / CELL_HEIGHT)
}

/** Height of a card whose fields are drawn in wrapped columns. */
function wrappedHeight(fieldCount, titled) {
  const rows = Math.max(Math.ceil(fieldCount / WRAP_COLUMNS), 1)

  return CONTENT_PADDING + rows * (EDITOR_ROW + WRAP_ROW_MARGIN) + (titled ? CARD_HEADER : 0)
}

function makeBlock({ id, kind, title = '', options, xywh }) {
  return compose.PageBlockMaker({
    blockID: id,
    kind,
    title,
    description: '',
    style: { wrap: { kind: 'card' } },
    options,
    xywh,
    meta: { tempID: id },
  })
}

/**
 * The stacked blocks an admin record screen renders: the module's fields in a
 * multi-column card, a card of its own for every field too wide to sit in a
 * column, the record's system fields, and its revisions where the module keeps
 * them.
 *
 * Grid lays a lone block out with CSS and anything more through gridstack,
 * where a block is as tall as its `xywh` says rather than as tall as the
 * viewport — hence the computed heights.
 *
 * @param {object} recordModule the module whose record is on screen
 * @param {string} opts.idPrefix distinguishes the screen's blocks in the grid
 * @param {string} opts.systemTitle heading for the system-fields card
 * @param {string} opts.revisionsTitle heading for the revisions card
 * @param {boolean} opts.isNew the record is not stored yet
 * @param {boolean} opts.withRevisions offer revisions where the module has them
 */
export function adminRecordBlocks(
  recordModule,
  { idPrefix, systemTitle, revisionsTitle, isNew = false, withRevisions = false },
) {
  const fields = recordModule?.fields || []
  const compact = fields.filter(f => !isWideField(f))
  const wide = fields.filter(isWideField)
  const systemFields = (recordModule?.systemFields?.() || []).filter(
    f => !(isNew && ASSIGNED_ON_SAVE.includes(f.name)),
  )

  const blocks = []
  let y = 0

  const stack = ({ id, kind = 'Record', title, options, height }) => {
    blocks.push(makeBlock({ id, kind, title, options, xywh: [0, y, GRID_COLUMNS, height] }))
    y += height
  }

  // An empty field list is what a Record block reads as "every module field",
  // so a card with nothing to put in it is left out rather than drawn empty.
  if (compact.length) {
    stack({
      id: `${idPrefix}_fields`,
      options: { fields: compact.map(f => f.name), recordFieldLayoutOption: 'wrap' },
      height: cells(wrappedHeight(compact.length, false)),
    })
  }

  for (const field of wide) {
    stack({
      id: `${idPrefix}_field_${field.name}`,
      options: { fields: [field.name] },
      height: cells(CONTENT_PADDING + (WIDE_FIELD_HEIGHT[field.kind] || WIDE_FIELD_DEFAULT_HEIGHT)),
    })
  }

  stack({
    id: `${idPrefix}_system`,
    title: systemTitle,
    options: { fields: systemFields.map(f => f.name), recordFieldLayoutOption: 'wrap' },
    height: cells(wrappedHeight(systemFields.length, true)),
  })

  if (withRevisions && recordModule?.config?.recordRevisions?.enabled) {
    stack({
      id: `${idPrefix}_revisions`,
      kind: 'RecordRevisions',
      title: revisionsTitle,
      options: { preload: true, sortDirection: 'desc' },
      height: REVISIONS_HEIGHT,
    })
  }

  return blocks
}

/**
 * The blocks the "all records" screen renders: a row of metric tiles, and the
 * table under them. A tile opens what it counted in a dialog of its own, which
 * leaves the table below it as the user left it.
 *
 * They are returned apart because the screen lays them out apart: the tiles are
 * a fixed-height row and the table takes whatever is left, which is a thing
 * Grid does for a lone block and gridstack does not do at all.
 *
 * @param {string} opts.moduleID the module being listed
 * @param {string} opts.idPrefix distinguishes the screen's blocks in the grid
 * @param {object} opts.metricLabels label per metric key
 */
export function adminRecordListBlocks({ moduleID, idPrefix, metricLabels = {} }) {
  const listID = `${idPrefix}_list`

  // The tiles share the row, so their width is the row divided among them
  const tileWidth = Math.floor(GRID_COLUMNS / LIST_METRICS.length)

  const tiles = LIST_METRICS.map(({ key, filter, deleted = 0 }, i) =>
    makeBlock({
      id: `${idPrefix}_metric_${key}`,
      kind: 'Metric',
      options: {
        metrics: [
          {
            label: metricLabels[key] || key,
            moduleID,
            metricField: 'count',
            operation: 'count',
            filter,
            deleted,
            // No blockID: the records open in the tile's own dialog
            drillDown: { enabled: true, recordListOptions: { fields: [] } },
          },
        ],
      },
      xywh: [i * tileWidth, 0, tileWidth, TILE_HEIGHT],
    }),
  )

  const list = makeBlock({
    id: listID,
    kind: 'RecordList',
    options: {
      moduleID,
      fields: [],
      perPage: 20,
      selectable: true,
      allowExport: true,
      // Choosing columns is the point of a module's own record table, and the
      // block hides that button unless told otherwise
      hideConfigureFieldsButton: false,
      // The admin list is the module's whole record set, deleted rows included
      showDeletedRecordsOption: true,
    },
    xywh: [0, 0, GRID_COLUMNS, TILE_HEIGHT],
  })

  return { tiles, list }
}
