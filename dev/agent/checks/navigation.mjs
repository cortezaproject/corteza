// Navigation checks, as a worked example of drive.mjs.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/navigation.mjs
//      node dev/agent/drive.mjs dev/agent/checks/navigation.mjs --only "back"
//
// Each of these caught a real defect: a disabled namespace trapped the user in
// its editor, every editor's Back dead-ended when deep-linked, and the module
// editor dirtied itself on load.
import { drive, check, expectPath, ids } from '../drive.mjs'

const NS = 'catalogue'
const id = ids(NS)

drive('back from a deep-linked namespace editor reaches the list', async page => {
  await page.open(`/compose/namespaces/edit/${NS}`)
  await page.back()
  expectPath(page, '/compose/namespaces')
})

drive('back from a deep-linked module editor reaches its list', async page => {
  await page.open(`/compose/namespace/${NS}/admin/modules/${id.module.catalogue_field}/edit`)
  await page.back()
  expectPath(page, `/compose/namespace/${NS}/admin/modules`)
})

drive('the namespace editor sidebar lists namespaces and searches them', async page => {
  await page.expandSidebar('compose')
  await page.open(`/compose/namespaces/edit/${NS}`)

  const all = await page.sidebarRows()
  check('namespaces are listed', all.length > 2, `${all.length} rows`)
  check('the group header is the list link', all[0] === 'Namespaces', all[0])

  await page.fill('[data-testid="app-sidebar"] input', 'catalog')
  const found = await page.sidebarRows()
  check('search narrows the list', found.length < all.length, `${all.length} -> ${found.length}`)
})

drive('an untouched module editor does not report unsaved changes', async page => {
  await page.open(`/compose/namespace/${NS}/admin/modules/${id.module.catalogue_field}/edit`)
  page.clearDialog()
  await page.open('/compose/namespaces')
  check('no unsaved-changes prompt', page.dialog() === null, page.dialog() || '')
})
