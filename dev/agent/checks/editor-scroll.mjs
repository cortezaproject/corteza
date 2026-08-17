// Every compose editor scrolls across its full width.
//
// The scroller used to be the centred `container mx-auto`, so the strips either
// side of it belonged to no scrollable element: a wheel with the pointer in a
// gutter scrolled nothing, and the pointer had to be over the column itself.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/editor-scroll.mjs
//      node dev/agent/drive.mjs dev/agent/checks/editor-scroll.mjs --only "chart"
import { drive, check, ids } from '../drive.mjs'

const NS = 'catalogue'
const id = ids(NS)

// The editor's own scroll container: a direct child of the form that scrolls.
// Deliberately not "any overflowing descendant" — with a short editor nothing
// overflows and that search latches onto an inner textarea or overlay instead,
// then reports the layout is wrong when it is merely small.
const FIND_SCROLLER = `
  const el = [...document.querySelectorAll('form > div, main > div, form > form > div')].find(e => {
    const o = getComputedStyle(e).overflowY
    return o === 'auto' || o === 'scroll'
  })
`

const geometry = page =>
  page.evaluate(
    new Function(`
    ${FIND_SCROLLER}
    if (!el) return null
    const inner = el.querySelector('.container')
    const r = el.getBoundingClientRect()
    const c = inner?.getBoundingClientRect()
    return {
      top: el.scrollTop,
      overflows: el.scrollHeight > el.clientHeight + 20,
      scroller: { left: r.left, right: r.right, width: r.width },
      column: c ? { left: c.left, right: c.right, width: c.width } : null,
    }
  `),
  )

const resetScroll = page =>
  page.evaluate(
    new Function(`
    ${FIND_SCROLLER}
    if (el) el.scrollTop = 0
  `),
  )

/** Assert the open editor scrolls with the pointer in either gutter. */
async function expectGutterScroll(page, label) {
  const start = await geometry(page)

  check(`${label}: the editor owns a scroll container`, !!start, start ? 'found' : 'none found')
  if (!start) return

  // The gutters only exist if the scroller is wider than what it centres. A
  // scroller with no .container inside it IS the centred column — the old shape.
  const hasGutters = !!start.column && start.scroller.width > start.column.width + 20
  check(
    `${label}: the scroller is wider than the centred column`,
    hasGutters,
    start.column
      ? `scroller ${Math.round(start.scroller.width)} vs column ${Math.round(start.column.width)}`
      : 'the scroller is the centred column itself',
  )
  if (!hasGutters) return

  check(`${label}: there is enough content to scroll`, start.overflows, `${start.overflows}`)
  if (!start.overflows) return

  for (const [side, x] of [
    ['left', Math.round(start.scroller.left + (start.column.left - start.scroller.left) / 2)],
    ['right', Math.round(start.column.right + (start.scroller.right - start.column.right) / 2)],
  ]) {
    await resetScroll(page)
    await page.raw.mouse.move(x, 300)
    await page.raw.mouse.wheel(0, 600)
    await page.raw.waitForTimeout(600)

    const after = await geometry(page)
    check(`${label}: a wheel in the ${side} gutter scrolls`, after?.top > 0, `top=${after?.top}`)
  }
}

// Short enough that every editor here has more content than fits, so "nothing
// scrolled" always means the wheel was ignored rather than the page being small.
const shortViewport = page => page.raw.setViewportSize({ width: 1400, height: 560 })

drive('the namespace editor scrolls from its gutters', async page => {
  await shortViewport(page)
  await page.open(`/compose/namespaces/edit/${NS}`, { settle: 2000 })
  await expectGutterScroll(page, 'namespace')
})

drive('the page editor scrolls from its gutters', async page => {
  await shortViewport(page)
  await page.open(`/compose/namespace/${NS}/admin/pages/${id.page.catalogue_content}/edit`, {
    settle: 2500,
  })
  await expectGutterScroll(page, 'page')
})

drive('the chart builder scrolls from its gutters', async page => {
  await shortViewport(page)
  await page.open(`/compose/namespace/${NS}/admin/charts/${id.chart.cat_chart_filtered}/edit`, {
    settle: 2500,
  })
  await expectGutterScroll(page, 'chart')
})
