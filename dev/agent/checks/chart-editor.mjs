// Chart builder form: grouping, required fields, and what save does with them.
//
// The editor used to be one long card of hand-rolled labels, with `name`
// painted red before the user had typed anything and no rule at all on the
// source module — so a chart with no module saved happily and then rendered
// nothing wherever it was placed.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/chart-editor.mjs
import { drive, check, ids } from '../drive.mjs'

const NS = 'catalogue'
const id = ids(NS)
const MODULE_ID = id.module.catalogue_field

const REPORT = {
  dimensions: [{ field: 'createdAt', modifier: 'MONTH' }],
  metrics: [{ field: 'count', aggregate: '', type: 'bar', label: 'count' }],
}

const countCharts = page =>
  page.evaluate(
    a =>
      window.__human
        .globals()
        .$ComposeAPI.chartList(a)
        .then(r => (r.set || []).length),
    { namespaceID: id.namespaceID, limit: 100 },
  )

/** Create a chart, run body against its editor, remove it. */
async function withChart(page, { name }, body) {
  const moduleID = MODULE_ID
  await page.open('/compose/namespaces')

  const chartID = await page.evaluate(
    async a => {
      const c = await window.__human.globals().$ComposeAPI.chartCreate({
        namespaceID: a.namespaceID,
        name: a.name,
        config: { reports: [{ ...a.report, moduleID: a.moduleID }] },
      })
      return c.chartID
    },
    { namespaceID: id.namespaceID, moduleID, name, report: REPORT },
  )

  try {
    await page.open(`/compose/namespace/${NS}/admin/charts/${chartID}/edit`, { settle: 2500 })
    await body(chartID)
  } finally {
    await page
      .evaluate(a => window.__human.globals().$ComposeAPI.chartDelete(a), {
        namespaceID: id.namespaceID,
        chartID,
      })
      .catch(() => {})
  }
}

// Header text of every panel, plus whether it is open. A collapsed PrimeVue
// Panel keeps its content element in the DOM and hides it, so presence of
// .p-panel-content answers nothing — the toggler's aria-expanded does.
const panels = page =>
  page.evaluate(() =>
    [...document.querySelectorAll('.p-panel')].map(p => ({
      header: p.querySelector('.p-panel-header')?.textContent?.trim().replace(/\s+/g, ' ') || '',
      open: p.querySelector('[aria-expanded]')?.getAttribute('aria-expanded') === 'true',
    })),
  )

// A required CFormGroup renders a red asterisk inside its label.
const requiredLabels = page =>
  page.evaluate(() =>
    [...document.querySelectorAll('label')]
      .filter(l => l.querySelector('.text-red-500'))
      .map(l => l.textContent.replace('*', '').trim()),
  )

drive('settings are grouped into panels, advanced ones collapsed', async page => {
  await withChart(page, { name: 'panel layout' }, async () => {
    const found = await panels(page)
    const headers = found.map(p => p.header)

    check(
      'the core groups are panels',
      ['General settings', 'Data Source', 'Dimensions', 'Metrics'].every(h =>
        headers.some(x => x.includes(h)),
      ),
      headers.join(' | '),
    )

    const openOf = h => found.find(p => p.header.includes(h))?.open

    check('General settings starts open', openOf('General settings') === true, '')
    check('Data Source starts open', openOf('Data Source') === true, '')
    check('Metrics starts open', openOf('Metrics') === true, '')
    check('Y-axis starts collapsed', openOf('Y-axis') === false, '')
    check('Legend starts collapsed', openOf('Legend') === false, '')
    check('Tools starts collapsed', openOf('Tools') === false, '')

    const body = await page.evaluate(() => document.body.textContent || '')
    check('the preview no longer says "Load data"', !body.includes('Load data'), '')
  })
})

// The scroller used to be the centred `container mx-auto`, so the gutters
// beside it belonged to no scrollable element and a wheel over them did
// nothing — the pointer had to be over the column itself.
drive('the page scrolls with the pointer in either gutter', async page => {
  await withChart(page, { name: 'scroll surface' }, async () => {
    // Geometry of the element that actually scrolls, and of the centred column
    // inside it. The gutters are the strips of the former either side of the
    // latter — measured off the scroller, not the window, since the section
    // sits right of the app sidebar.
    const geometry = () =>
      page.evaluate(() => {
        const el = [...document.querySelectorAll('form *')].find(
          e => e.scrollHeight > e.clientHeight + 20 && getComputedStyle(e).overflowY !== 'visible',
        )
        if (!el) return null
        const inner = el.querySelector('.container')
        const r = el.getBoundingClientRect()
        const c = inner?.getBoundingClientRect()
        return {
          top: el.scrollTop,
          scroller: { left: r.left, right: r.right, width: r.width },
          column: c ? { left: c.left, right: c.right, width: c.width } : null,
        }
      })

    const resetScroll = () =>
      page.evaluate(() => {
        const el = [...document.querySelectorAll('form *')].find(
          e => e.scrollHeight > e.clientHeight + 20 && getComputedStyle(e).overflowY !== 'visible',
        )
        if (el) el.scrollTop = 0
      })

    const start = await geometry()
    check('the page has something to scroll', !!start, start ? `${start.scroller.width}px` : 'none')

    // No inner column means the scroller IS the centred container — the shape
    // that made the gutters dead in the first place.
    const hasGutters = !!start?.column && start.scroller.width > start.column.width + 20
    check(
      'the scroller is wider than the centred column, so gutters belong to it',
      hasGutters,
      start?.column
        ? `scroller ${Math.round(start.scroller.width)} vs column ${Math.round(start.column.width)}`
        : 'the scroller is the centred column itself',
    )
    if (!hasGutters) return

    for (const [side, x] of [
      ['left', Math.round(start.scroller.left + (start.column.left - start.scroller.left) / 2)],
      ['right', Math.round(start.column.right + (start.scroller.right - start.column.right) / 2)],
    ]) {
      await resetScroll()
      await page.raw.mouse.move(x, 450)
      await page.raw.mouse.wheel(0, 600)
      await page.raw.waitForTimeout(600)

      const after = await geometry()
      check(
        `a wheel in the ${side} gutter scrolls`,
        after && after.top > 0,
        `x=${x} top=${after?.top}`,
      )
    }
  })
})

drive('the preview is just the chart', async page => {
  await withChart(page, { name: 'no refresh button' }, async () => {
    const refreshButtons = await page.evaluate(
      () => document.querySelectorAll('.pi-refresh').length,
    )
    check('the refresh button is gone', refreshButtons === 0, `${refreshButtons} found`)

    const canvases = await page.evaluate(() => document.querySelectorAll('canvas').length)
    check('the chart still renders', canvases > 0, `${canvases} canvas`)
  })
})

drive('name and module are marked required', async page => {
  await withChart(page, { name: 'required markers' }, async () => {
    const labels = await requiredLabels(page)

    check('Chart name is marked required', labels.includes('Chart name'), labels.join(' | '))
    check('Module is marked required', labels.includes('Module'), labels.join(' | '))
    check('only the genuinely required fields are marked', labels.length === 2, `${labels.length}`)
  })
})

drive('an untouched editor does not report unsaved changes', async page => {
  await withChart(page, { name: 'clean on load' }, async () => {
    page.clearDialog()
    await page.back()
    const d = page.dialog()
    check('no confirm on a clean leave', !d, d || '')
  })
})

drive('a valid edit saves', async page => {
  await withChart(page, { name: 'before rename' }, async chartID => {
    await page.fill('#name', 'after rename')
    await page.click('button[type="submit"]', { settle: 3000 })

    const stored = await page.evaluate(
      a =>
        window.__human
          .globals()
          .$ComposeAPI.chartRead(a)
          .then(c => c.name),
      { namespaceID: id.namespaceID, chartID },
    )
    check('the new name reached the server', stored === 'after rename', stored)
  })
})

drive('an empty name blocks save and says why', async page => {
  await withChart(page, { name: 'will be emptied' }, async chartID => {
    await page.fill('#name', '')
    await page.click('button[type="submit"]', { settle: 2000 })

    const body = await page.evaluate(() => document.body.textContent || '')
    check(
      'the form reports an error',
      /required/i.test(body),
      body.match(/[^.\n]*[Rr]equired[^.\n]*/)?.[0]?.trim().slice(0, 120) || '(no message)',
    )

    const stored = await page.evaluate(
      a =>
        window.__human
          .globals()
          .$ComposeAPI.chartRead(a)
          .then(c => c.name),
      { namespaceID: id.namespaceID, chartID },
    )
    check('the empty name never reached the server', stored === 'will be emptied', stored)
  })
})

// The realistic case, and the one that used to save happily: a brand new chart
// where the author names it and never picks a module. It is driven through the
// create route rather than seeded, because the API rejects a module-less report
// outright — only the builder could ever produce one.
drive('a new chart with no source module cannot be saved', async page => {
  await page.open('/compose/namespaces')
  const before = await countCharts(page)

  await page.open(`/compose/namespace/${NS}/admin/charts/create`, { settle: 2500 })
  await page.fill('#name', 'never picked a module')
  await page.click('button[type="submit"]', { settle: 2500 })

  check(
    'the form reports an error',
    /required/i.test(await page.evaluate(() => document.body.textContent || '')),
    '',
  )
  check(
    'it stayed on the create route',
    page.raw.url().includes('/charts/create'),
    page.raw.url(),
  )
  check('no chart was created', (await countCharts(page)) === before, '')
})
