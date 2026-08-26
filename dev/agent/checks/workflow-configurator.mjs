// The workflow configuration dialog keeps its actions in view.
//
// Delete, Cancel and Save used to be the last row of the dialog's own scroll
// body, so a workflow with a handful of declared inputs pushed them hundreds of
// pixels below the fold: saving meant scrolling to the bottom of a form that had
// just been read from the top.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/workflow-configurator.mjs
import { drive, check } from '../drive.mjs'

const HANDLE = 'agent_wf_configurator'

// Short enough that the form overflows whatever the theme's spacing is, so
// "the actions are on screen" always means they were pinned rather than that
// the dialog happened to fit.
const VIEWPORT = { width: 1400, height: 620 }

async function build(page) {
  return page.evaluate(async handle => {
    const { $AutomationAPI } = window.__human.globals()
    if (!$AutomationAPI) throw new Error('no $AutomationAPI global')

    // One left behind by a crashed run would take the handle.
    const stale = await $AutomationAPI.workflowList({ query: handle, limit: 10 })
    for (const w of stale.set || []) {
      if (w.handle === handle) await $AutomationAPI.workflowDelete({ workflowID: w.workflowID })
    }

    const wf = await $AutomationAPI.workflowCreate({
      handle,
      enabled: true,
      // The client refuses an absent runAs/ownedBy; '0' is their unset
      // sentinel and the server fills the owner in from the caller.
      runAs: '0',
      ownedBy: '0',
      scope: {},
      meta: {
        name: 'Agent configurator footer',
        description: 'Built by dev/agent/checks/workflow-configurator.mjs',
        // Declared inputs are the cheapest way to make the form taller than
        // the dialog; each renders a row of four controls.
        input: ['alpha', 'beta', 'gamma'].map(name => ({
          name,
          label: name,
          types: ['String'],
          required: false,
        })),
      },
      steps: [],
      paths: [],
    })
    return wf.workflowID
  }, HANDLE)
}

const teardown = (page, workflowID) =>
  page
    .evaluate(
      id => window.__human.globals().$AutomationAPI.workflowDelete({ workflowID: id }),
      workflowID,
    )
    .catch(() => {})

const measure = page =>
  page.evaluate(() => {
    const dlg = document.querySelector('.p-dialog')
    const content = dlg.querySelector('.p-dialog-content')
    const footer = dlg.querySelector('.p-dialog-footer')
    const save = dlg.querySelector('[data-test-id="button-save-workflow"]')
    const label = t =>
      [...(footer?.querySelectorAll('button') || [])].some(b => t.test(b.innerText))
    const r = save?.getBoundingClientRect()
    return {
      bodyOverflows: content.scrollHeight > content.clientHeight + 1,
      bodyScrollTop: content.scrollTop,
      saveInFooter: !!(footer && save && footer.contains(save)),
      saveInBody: !!(save && content.contains(save)),
      hasCancel: label(/cancel/i),
      hasDelete: label(/delete/i),
      saveTop: r ? Math.round(r.top) : null,
      saveBottom: r ? Math.round(r.bottom) : null,
      viewportH: window.innerHeight,
    }
  })

drive('the workflow configurator keeps Save in the dialog footer', async page => {
  await page.raw.setViewportSize(VIEWPORT)
  await page.open('/workflow/list', { settle: 1000 })
  const workflowID = await build(page)

  try {
    await page.open(`/workflow/${workflowID}/edit`, {
      until: '[data-test-id="button-configure-workflow"]',
    })
    // Waited for wherever it renders, so a regression is reported by the
    // assertions below rather than as a timeout on a selector.
    await page.click('[data-test-id="button-configure-workflow"]', {
      until: '.p-dialog [data-test-id="button-save-workflow"]',
    })

    const m = await measure(page)

    check('the form is taller than the dialog', m.bodyOverflows)
    check('save is in the footer, not the scroll body', m.saveInFooter && !m.saveInBody)
    check('cancel is in the footer', m.hasCancel)
    check('delete is in the footer', m.hasDelete)
    check(
      'save is on screen with the form unscrolled',
      m.bodyScrollTop === 0 && m.saveTop >= 0 && m.saveBottom <= m.viewportH,
      `save ${m.saveTop}-${m.saveBottom} in ${m.viewportH}px, scrollTop ${m.bodyScrollTop}`,
    )

    // The footer is a Teleport target in the editor, so the buttons render one
    // component up from the state they act on; a broken wiring still draws them.
    await page.click('.p-dialog-footer button', { hasText: 'Cancel', settle: 600 })
    check(
      'cancel from the footer closes the dialog',
      await page.evaluate(() => !document.querySelector('.p-dialog')),
    )
  } finally {
    await teardown(page, workflowID)
  }
})
