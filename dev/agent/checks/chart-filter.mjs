// Chart builder preview: filters that carry template expressions.
//
// A chart is namespace-level, so its report filter may read `${recordID}` /
// `${record.values.x}` / `${ownerID}` — variables that only resolve wherever
// the chart is later placed. The builder has no record, and sending such a
// filter unevaluated came back from the server as "found an illegal token",
// which the preview then rendered as [object Object].
//
// The charts are made and removed here rather than seeded, so the check needs
// nothing of the namespace beyond a module to report on.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/chart-filter.mjs
import { drive, check, ids } from '../drive.mjs'

const NS = 'catalogue'
const id = ids(NS)
const MODULE_ID = id.module.catalogue_field

const bodyText = page => page.evaluate(() => document.body.textContent || '')

// Every XHR the document made, read from the resource timeline. Vite serves
// hundreds of modules on a cold load and fills the default 250-entry buffer
// before the first XHR is made, so the buffer is emptied and the preview
// re-fired — that is what makes "no request" an observation, not an overflow.
async function reportCallsOnRefresh(page) {
  await page.evaluate(() => {
    performance.setResourceTimingBufferSize(1000)
    performance.clearResourceTimings()
  })
  await page.click('button:has(.pi-refresh)', { settle: 2500 })
  return page.evaluate(() =>
    performance
      .getEntriesByType('resource')
      .map(e => e.name)
      .filter(u => u.includes('/record/report')),
  )
}

/** Build a chart with the given filter, run body against its editor, remove it. */
async function withChart(page, { filter, name }, body) {
  await page.open('/compose/namespaces')

  const chartID = await page.evaluate(
    async a => {
      const chart = await window.__human.globals().$ComposeAPI.chartCreate({
        namespaceID: a.namespaceID,
        name: a.name,
        config: {
          reports: [
            {
              moduleID: a.moduleID,
              filter: a.filter,
              dimensions: [{ field: 'createdAt', modifier: 'MONTH' }],
              metrics: [{ field: 'count', aggregate: '', type: 'bar', label: 'count' }],
            },
          ],
        },
      })
      return chart.chartID
    },
    { namespaceID: id.namespaceID, moduleID: MODULE_ID, filter, name },
  )

  try {
    await page.open(`/compose/namespace/${NS}/admin/charts/${chartID}/edit`, { settle: 2500 })
    await body()
  } finally {
    await page
      .evaluate(a => window.__human.globals().$ComposeAPI.chartDelete(a), {
        namespaceID: id.namespaceID,
        chartID,
      })
      .catch(() => {})
  }
}

const noObjectObject = body =>
  check(
    'nothing renders as [object Object]',
    !body.includes('[object Object]'),
    body.includes('[object Object]') ? 'preview shows [object Object]' : 'no [object Object]',
  )

drive('a record-variable filter is refused rather than sent unevaluated', async page => {
  const spec = { filter: 'recordID = ${recordID}', name: 'record-var filter' }

  await withChart(page, spec, async () => {
    const body = await bodyText(page)
    const calls = await reportCallsOnRefresh(page)

    check(
      'the preview says record variables resolve at placement',
      /record variables/i.test(body) && /where the chart is placed/i.test(body),
      body.match(/This filter uses[^.]*\.\s*[^.]*\./)?.[0] || '(notice not found)',
    )
    check('no report request was fired', calls.length === 0, `${calls.length} call(s)`)
    noObjectObject(body)
    check(
      'no illegal-token error from the server',
      !/illegal token/i.test(body),
      body.match(/[^.\n]*illegal token[^.\n]*/i)?.[0]?.trim() || 'no server error',
    )
  })
})

drive('a genuine server error reaches the preview as its message', async page => {
  await withChart(page, { filter: 'nosuchfield = 5', name: 'bad field filter' }, async () => {
    const body = await bodyText(page)

    noObjectObject(body)
    check(
      'the server message is shown',
      /nosuchfield|unknown field|does not exist/i.test(body),
      body.match(/[^.\n]*nosuchfield[^.\n]*/i)?.[0]?.trim().slice(0, 160) || '(no message)',
    )
  })
})

drive('a filter with no record variables still reports', async page => {
  await withChart(page, { filter: "sel_badge = 'live'", name: 'plain filter' }, async () => {
    const body = await bodyText(page)
    const calls = await reportCallsOnRefresh(page)

    noObjectObject(body)
    check('the report was actually requested', calls.length > 0, `${calls.length} call(s)`)
  })
})
