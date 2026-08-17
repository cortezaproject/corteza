// Every guarded editor, audited the same way: open an existing record from its
// list, leave without touching anything, and assert it stays quiet — then type
// into it and assert it warns.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/editor-guards.mjs
//      node dev/agent/drive.mjs dev/agent/checks/editor-guards.mjs --only "role"
//
// Both halves are the point. An editor that never warns loses work; one that
// warns on load — because it dirties itself resolving defaults, or captures its
// baseline before the data lands — trains people to click through the warning,
// which loses work just as effectively.
//
// A red run here means a defect. The one thing reported as SKIPPED instead is
// what this server does not have to offer: a list with no rows leaves nothing to
// open, and a check that cannot go green is one people learn to skip, after which
// it protects nothing. Skips print on green runs and each is a hole in the
// coverage — but a row that opens nothing, an editor without a Back button, or an
// unexpected console error is a defect and fails.
import { drive, check } from '../drive.mjs'

// Editors reached from a list. Only the list path is named: the check clicks
// the first row, so no resource ID is pinned here to go stale.
//
// `noise` allows console errors an editor emits for reasons of its own, so the
// check still fails on any error that is NOT on the list. Each entry names the
// bug it is waiting on and should go when that is fixed.
const EDITORS = [
  { name: 'user', list: '/admin/system/users' },
  { name: 'user group', list: '/admin/system/user-groups' },
  { name: 'role', list: '/admin/system/roles' },
  { name: 'application', list: '/admin/system/applications' },
  { name: 'auth client', list: '/admin/system/auth-clients' },
  { name: 'llm provider', list: '/admin/system/llm-providers' },
  { name: 'template', list: '/admin/system/templates' },
  { name: 'queue', list: '/admin/system/queues' },
  { name: 'api gateway route', list: '/admin/system/api-gateway' },
  { name: 'connection', list: '/admin/system/connections' },
  {
    name: 'data source',
    list: '/admin/system/data-sources',
    // Its labels carry `{{module}}` and a JSON example, both of which vue-i18n
    // reads as placeholders and fails to compile
    // (locale/en/human-webapp/system.yaml:590,598).
    noise: [/Message compilation error/],
  },
  { name: 'workflow', list: '/admin/automation/workflows' },
  { name: 'taq', list: '/admin/automation/taq' },
  { name: 'federation node', list: '/admin/federation/nodes' },
]

/** Record a hole in the coverage: printed and counted, but not a failure. */
function skip(why) {
  return check(`SKIPPED — ${why}`, true)
}

/** Console and network errors the editor is not entitled to. */
function assertQuiet(page, editor) {
  const unexpected = page.problems().filter(p => !(editor.noise || []).some(re => re.test(p)))
  check('no unexpected console or network errors', !unexpected.length, unexpected.join(' | '))
}

// Opens the first row of a list, or returns null when there is nothing to open.
// An empty table still renders a row — its "no records" message — and clicking
// that navigates nowhere, which reads as a broken editor rather than a gap.
async function openFirstRow(page, list) {
  await page.open(list)
  const rows = page.locator('.p-datatable-tbody > tr:not(.p-datatable-empty-message)')
  await rows
    .first()
    .waitFor({ state: 'visible', timeout: 20000 })
    .catch(() => {})
  if ((await rows.count()) === 0) return null

  await rows.first().click()
  await page.raw.waitForTimeout(2500)
  return page.path()
}

async function leaveViaBack(page) {
  page.clearDialog()
  const back = page.locator('[data-testid="editor-back"]').first()
  if (!(await back.count())) return { left: false, dialog: null }
  await back.click()
  await page.raw.waitForTimeout(2000)
  return { left: true, dialog: page.dialog() }
}

// drive's blanket console-error check is off because it cannot tell an editor's
// known noise from a new fault; assertQuiet does that per editor instead.
const OPTS = { allowProblems: true }

for (const editor of EDITORS) {
  drive(
    `${editor.name}: untouched does not warn`,
    async page => {
      const at = await openFirstRow(page, editor.list)
      if (!at) return skip(`no ${editor.name} on this server to open`)
      if (!check('the row opened an editor', at !== editor.list, `still at ${at}`)) return

      const { left, dialog } = await leaveViaBack(page)
      if (!check('the editor has a Back button', left, 'no [data-testid="editor-back"]')) return

      check('no confirm on a clean leave', !dialog, dialog || '')
      assertQuiet(page, editor)
    },
    OPTS,
  )

  drive(
    `${editor.name}: warns once edited`,
    async page => {
      const at = await openFirstRow(page, editor.list)
      if (!at) return skip(`no ${editor.name} on this server to open`)
      if (!check('the row opened an editor', at !== editor.list, `still at ${at}`)) return

      const input = page.locator('input[type="text"]:not([readonly])').first()
      if (!(await input.count())) return skip(`the ${editor.name} editor has no text field to type in`)

      // Typed, not filled: a fill() sets the value with no pointer or key event,
      // so it exercises a path no user can take.
      await input.click()
      await page.raw.keyboard.type('drive check edit')
      await page.raw.keyboard.press('Tab')
      await page.raw.waitForTimeout(1200)

      const { left, dialog } = await leaveViaBack(page)
      if (!check('the editor has a Back button', left, 'no [data-testid="editor-back"]')) return

      check('leaving prompts a confirm', !!dialog, dialog || 'no dialog')
      assertQuiet(page, editor)
    },
    OPTS,
  )
}
