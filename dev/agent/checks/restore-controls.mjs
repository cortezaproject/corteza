// Does the restore control a deleted resource offers actually restore it.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/restore-controls.mjs
//
// The sweep that standardised these buttons is guarded by a source-level test,
// which can only say the label and the glyph agree everywhere. It cannot say
// the button is reachable on a deleted resource, that the server takes the
// call, or that the screen stops calling the resource deleted afterwards.
// This does, against the real API, on two of the screens that gained one.
import { drive, check } from '../drive.mjs'

const EMAIL = 'agent-restore-check@local.dev'
// Earlier runs of this check left duplicates behind — a soft delete does not
// free the email, and users have no hard delete. The row is addressed by a name
// only the probe carries, so those leftovers cannot be mistaken for it.
const NAME = 'Restore Check Probe'

async function makeDeletedUser(page) {
  return page.evaluate(
    async ({ email, name }) => {
      const { $SystemAPI } = window.__human.globals()

      // Reuse the probe rather than making another.
      const existing = await $SystemAPI.userList({ email, deleted: 1, limit: 50 })
      let user = (existing.set || []).find(u => u.email === email)

      if (!user) {
        // The JS client refuses an absent userGroupID; the REST API does not.
        user = await $SystemAPI.userCreate({ email, name, handle: '', userGroupID: '0' })
      }

      // Naming it takes a live resource, so restore it, name it, delete it.
      const read = await $SystemAPI.userRead({ userID: user.userID })
      if (read.deletedAt) await $SystemAPI.userUndelete({ userID: user.userID })
      if (read.name !== name) {
        await $SystemAPI.userUpdate({ ...read, deletedAt: undefined, name })
      }
      await $SystemAPI.userDelete({ userID: user.userID })

      const after = await $SystemAPI.userRead({ userID: user.userID })
      return { userID: user.userID, deletedAt: after.deletedAt || null, name: after.name }
    },
    { email: EMAIL, name: NAME },
  )
}

async function readUser(page, userID) {
  return page.evaluate(async id => {
    const { $SystemAPI } = window.__human.globals()
    const u = await $SystemAPI.userRead({ userID: id })
    return { deletedAt: u.deletedAt || null }
  }, userID)
}

drive('a deleted user is restored from its editor', async page => {
  await page.open('/admin/system/users')
  const made = await makeDeletedUser(page)
  check('the probe user starts out deleted', !!made.deletedAt, `deletedAt=${made.deletedAt}`)

  await page.open(`/admin/system/users/${made.userID}`)

  // PrimeVue draws a button icon as a span, so `i.pi-replay` matches nothing.
  const restore = 'button:has(span.pi-replay)'
  const label = await page.text(restore)
  check('the editor offers Restore', label.trim() === 'Restore', `read "${label.trim()}"`)

  const deleteGone = await page.raw.locator('button:has(span.pi-trash)').count()
  check('delete is not offered beside it', deleteGone === 0, `${deleteGone} delete buttons`)

  await page.click(restore, { settle: 1500 })

  const after = await readUser(page, made.userID)
  check('the user is no longer deleted', after.deletedAt === null, `deletedAt=${after.deletedAt}`)

  const stillThere = await page.raw.locator(restore).count()
  check('restore stops being offered once restored', stillThere === 0, `${stillThere} still shown`)

  await page.evaluate(async id => {
    const { $SystemAPI } = window.__human.globals()
    await $SystemAPI.userDelete({ userID: id })
  }, made.userID)
})

drive('a deleted user offers Restore in the list row menu', async page => {
  await page.open('/admin/system/users')
  const made = await makeDeletedUser(page)

  // The list hides deleted rows by default; the filter lives in a popover.
  await page.open('/admin/system/users')
  await page.click('button:has(span.pi-filter)', { settle: 600 })
  await page.click('label[for="del2"]', { settle: 1500 })

  // Address the row by the user on it: the list still holds other deleted
  // rows, and the actions popup is shared by all of them.
  const row = page.raw.locator('tr', { hasText: NAME })
  const rows = await row.count()
  check('the deleted user is listed', rows === 1, `${rows} rows matched ${NAME}`)

  await row.locator('.row-action-btn').click()
  await page.raw.locator('.p-menu').waitFor({ state: 'visible' })

  const menu = await page.text('.p-menu')
  check('the row menu offers Restore', menu.includes('Restore'), `menu read "${menu.trim()}"`)
  check('and does not say Undelete', !menu.includes('Undelete'), menu.trim())

  const glyph = await page.raw.locator('.p-menu span.pi-replay').count()
  check('with the replay glyph', glyph === 1, `${glyph} replay glyphs`)

  await page.evaluate(async id => {
    const { $SystemAPI } = window.__human.globals()
    await $SystemAPI.userDelete({ userID: id }).catch(() => {})
  }, made.userID)
})
