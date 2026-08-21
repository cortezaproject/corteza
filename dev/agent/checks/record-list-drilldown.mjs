// Does a metric tile narrow the record list, say so, and let go again.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/record-list-drilldown.mjs
//
// The tile and the table are two blocks that only meet over the event bus, so
// nothing either of them can be mounted with proves the round trip. This does,
// against the real API: that a tile filters the table, that the chip it leaves
// behind clears it, and that a metric counting deleted records takes the table
// into deleted mode and brings it back.
import { drive, check } from '../drive.mjs'

// Reused across runs rather than remade — a namespace has no hard delete, so a
// fresh one per run would pile up.
const NS_SLUG = 'agent_drilldown_probe'
const MODULE = 'Drilldown Probe'
const LIVE = ['DP-3', 'DP-4', 'DP-5']
const DELETED = ['DP-1', 'DP-2']

async function provision(page) {
  return page.evaluate(
    async ({ slug, moduleName, live, deleted }) => {
      const { $ComposeAPI } = window.__human.globals()

      const found = await $ComposeAPI.namespaceList({ slug, limit: 50 })
      let ns = (found.set || []).find(n => n.slug === slug)
      if (!ns) {
        ns = await $ComposeAPI.namespaceCreate({
          name: 'Drilldown Probe',
          slug,
          enabled: true,
          meta: {},
        })
      }

      const mods = await $ComposeAPI.moduleList({ namespaceID: ns.namespaceID, limit: 50 })
      let mod = (mods.set || []).find(m => m.name === moduleName)
      if (!mod) {
        mod = await $ComposeAPI.moduleCreate({
          namespaceID: ns.namespaceID,
          name: moduleName,
          handle: 'drilldown_probe',
          meta: {},
          fields: [{ name: 'code', label: 'Code', kind: 'String', options: {} }],
        })
      }

      const ids = { namespaceID: ns.namespaceID, moduleID: mod.moduleID }
      const codeOf = r => (r.values || []).find(v => v.name === 'code')?.value

      // Every probe record, whatever state it is in
      const all = await $ComposeAPI.recordList({ ...ids, deleted: 1, limit: 100 })
      const byCode = new Map((all.set || []).map(r => [codeOf(r), r]))

      for (const code of [...live, ...deleted]) {
        if (!byCode.has(code)) {
          const made = await $ComposeAPI.recordCreate({
            ...ids,
            values: [{ name: 'code', value: code }],
          })
          byCode.set(code, made)
        }
      }

      // Put each one back in the state this check expects to find it in
      for (const code of live) {
        const r = byCode.get(code)
        if (r.deletedAt) await $ComposeAPI.recordUndelete({ ...ids, recordID: r.recordID })
      }
      for (const code of deleted) {
        const r = byCode.get(code)
        if (!r.deletedAt) await $ComposeAPI.recordDelete({ ...ids, recordID: r.recordID })
      }

      return ids
    },
    { slug: NS_SLUG, moduleName: MODULE, live: LIVE, deleted: DELETED },
  )
}

const state = () => ({
  codes: [...document.querySelectorAll('tbody tr')]
    .map(r => (r.innerText.match(/DP-\d+/) || [''])[0])
    .filter(Boolean)
    .sort(),
  chips: [...document.querySelectorAll('.p-chip')].map(c => c.innerText.trim()),
  onDeleted: !!document.querySelector('button.p-button-warn'),
  buttons: [...document.querySelectorAll('button')].map(b => b.innerText.trim()).filter(Boolean),
  // What each tile is reporting, to check the list against
  tiles: Object.fromEntries(
    [...document.querySelectorAll('.metric-item')].map(m => [
      m.querySelector('span')?.innerText.trim(),
      m.querySelector('.metric-value')?.innerText.trim(),
    ]),
  ),
})

const tile = name => `.metric-item:has(span:text-is("${name}"))`
const REMOVE_CHIP = '.p-chip [data-pc-section="removeicon"], .p-chip .p-chip-remove-icon'

drive('a tile narrows the list and names itself on it', async page => {
  await page.open('/')
  const { namespaceID, moduleID } = await provision(page)
  await page.open(`/compose/namespace/${namespaceID}/admin/modules/${moduleID}/records`, {
    settle: 3000,
  })

  const before = await page.evaluate(state)
  check('the live records are listed', before.codes.join() === LIVE.join())
  check('nothing is narrowing it yet', before.chips.length === 0)

  await page.click(tile('Owned by me'), { settle: 2000 })

  const after = await page.evaluate(state)
  // The probe's records belong to whoever the check runs as, so this tile does
  // not always narrow. What it must always do is agree with its own count.
  check(
    'the list shows what the tile counted',
    after.codes.length === Number(after.tiles['Owned by me']),
  )
  check(
    'a chip says which tile did it',
    after.chips.some(c => /Owned by me/.test(c)),
  )
})

drive('removing the chip puts the list back', async page => {
  await page.open('/')
  const { namespaceID, moduleID } = await provision(page)
  await page.open(`/compose/namespace/${namespaceID}/admin/modules/${moduleID}/records`, {
    settle: 3000,
  })

  await page.click(tile('Owned by me'), { settle: 2000 })
  await page.click(REMOVE_CHIP, { settle: 2000 })

  const out = await page.evaluate(state)
  check('the records are back', out.codes.join() === LIVE.join())
  check('and the chip is gone', out.chips.length === 0)
})

drive('a deleted-records tile takes the list into deleted mode and back', async page => {
  await page.open('/')
  const { namespaceID, moduleID } = await provision(page)
  await page.open(`/compose/namespace/${namespaceID}/admin/modules/${moduleID}/records`, {
    settle: 3000,
  })

  await page.click(tile('Deleted records'), { settle: 2500 })

  const shown = await page.evaluate(state)
  check('the deleted records are listed', shown.codes.join() === DELETED.join())
  check(
    'a chip says so',
    shown.chips.some(c => /Deleted records/.test(c)),
  )
  check('the list is on deleted records', shown.onDeleted)
  // What it made would land among the existing ones, out of this view's reach
  check('and offers no new record', !shown.buttons.includes('New Record'))

  await page.click(REMOVE_CHIP, { settle: 2500 })

  const back = await page.evaluate(state)
  check('clearing returns to the existing records', back.codes.join() === LIVE.join())
  check('and to existing mode', !back.onDeleted)
  check('and offers a new record again', back.buttons.includes('New Record'))
})

drive('the whole-module tile reads as a reset', async page => {
  await page.open('/')
  const { namespaceID, moduleID } = await provision(page)
  await page.open(`/compose/namespace/${namespaceID}/admin/modules/${moduleID}/records`, {
    settle: 3000,
  })

  await page.click(tile('Deleted records'), { settle: 2500 })
  await page.click(tile('Total records'), { settle: 2500 })

  const out = await page.evaluate(state)
  check('everything is listed again', out.codes.join() === LIVE.join())
  check('with nothing narrowing it', out.chips.length === 0)
  check('and back on existing records', !out.onDeleted)
})
