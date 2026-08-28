// The field configurator's Validation tab has to say which way round `test` is.
//
// A validator's `test` names the condition under which the value is REJECTED,
// so the obvious "value >= 0 && value <= 5" admits 7 and refuses 3. The panel
// used to describe the box as testing validity, which confirms the wrong
// reading.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/field-validator-wording.mjs
import { drive, check } from '../drive.mjs'

const NS = 'fieldexpr_repro'
const MODULE = 'd_val'

drive('the validation panel names the rejection condition', async page => {
  // __human is only on the page once the app has booted
  await page.open(`/compose/namespace/${NS}/admin/modules`)

  const moduleID = await page.evaluate(
    async ([ns, handle]) => {
      const { $ComposeAPI } = window.__human.globals()
      const nn = await $ComposeAPI.namespaceList({ query: ns, limit: 50 })
      const namespace = (nn.set || []).find(n => n.slug === ns)
      if (!namespace) throw new Error(`namespace ${ns} not found`)

      const mm = await $ComposeAPI.moduleList({ namespaceID: namespace.namespaceID, limit: 100 })
      const mod = (mm.set || []).find(m => m.handle === handle)
      if (!mod) throw new Error(`module ${handle} not found`)

      return mod.moduleID
    },
    [NS, MODULE],
  )

  await page.open(`/compose/namespace/${NS}/admin/modules/${moduleID}/edit`)
  await page.locator('button:has(span.pi-cog)').first().click()
  await page.locator('[role="tab"]', { hasText: 'Value validation' }).click()

  const body = await page.locator('body').innerText()

  check(
    'says the expression names the rejection condition',
    /condition under which the value is rejected/i.test(body),
    body.slice(0, 400),
  )
  check(
    'no longer describes the box as testing validity',
    !/test validity of the input value/i.test(body),
  )
  // vue-i18n reads a bare | as the plural separator, so the operator is written
  // as a literal in the message. If that escape is dropped the string does not
  // compile at all and the panel renders empty.
  check(
    'the || operator survives i18n compilation',
    body.includes('value < 0 || value > 5'),
    body.slice(0, 400),
  )
})
