// What a calendar block's event sources are, and whether each one survives
// being chosen and then renders.
//
// Run: node dev/agent/drive.mjs dev/agent/checks/calendar-feeds.mjs
//
// A feed names its source with the resource it reads (`compose:record`,
// `system:reminder`) and the block runtime dispatches on that. The configurator
// once wrote a short `resourceType` of its own instead, which the model does not
// declare — so the block was re-made without it on every open and every save,
// and picking Reminders silently reverted to Records (issue #44). The reminder
// loader existed in lib/js and was wired to nothing, so the option could not
// have rendered anything either.
//
// The check builds its own namespace so it depends on no fixture, and deletes it
// again; a reminder is per-user, so one is created for whoever the browser is
// logged in as and removed with it.
import { drive, check } from '../drive.mjs'

const SLUG = 'agent_calendar_feeds'
const RECORD_TITLES = ['Kickoff meeting', 'Design review']

// A reminder feed is a date range, not a due list: the block asks for whatever
// falls inside the month on screen, on both sides of now. One reminder each side
// is what says so — a check with only a past one passes just as well against a
// feed that drops everything still to come.
const REMINDER_TITLES = { past: 'Reminder already passed', future: 'Reminder still to come' }

// Everything this check needs, created through the app's own API clients.
async function build(page) {
  return page.evaluate(
    async ({ slug, reminderTitles, recordTitles }) => {
      const { $ComposeAPI, $SystemAPI, $Auth } = window.__human.globals()

      // A namespace left behind by a crashed run would take the slug.
      const stale = await $ComposeAPI.namespaceList({ slug, limit: 1 })
      for (const ns of stale.set || []) {
        await $ComposeAPI.namespaceDelete({ namespaceID: ns.namespaceID })
      }

      const ns = await $ComposeAPI.namespaceCreate({
        name: 'Agent calendar feeds',
        slug,
        enabled: true,
        meta: {},
      })
      const namespaceID = ns.namespaceID

      const mod = await $ComposeAPI.moduleCreate({
        namespaceID,
        name: 'Calendar events',
        handle: 'calendar_events',
        meta: {},
        fields: [
          { name: 'name', label: 'Name', kind: 'String', options: {} },
          { name: 'starts_at', label: 'Starts at', kind: 'DateTime', options: {} },
        ],
      })

      // Midday today, so the event lands inside the month the calendar opens on.
      const at = new Date()
      at.setHours(12, 0, 0, 0)

      for (const name of recordTitles) {
        await $ComposeAPI.recordCreate({
          namespaceID,
          moduleID: mod.moduleID,
          values: [
            { name: 'name', value: name },
            { name: 'starts_at', value: at.toISOString() },
          ],
        })
      }

      // Both ends of today, so each sits inside the month the calendar opens on
      // whatever the date. Only a run in the first or last five minutes of a day
      // blurs which side of now they are on, and both still render there.
      const dayEdge = h => {
        const t = new Date()
        t.setHours(h, h === 0 ? 5 : 55, 0, 0)
        return t.toISOString()
      }

      const reminders = {}
      for (const [when, title] of Object.entries(reminderTitles)) {
        reminders[when] = (
          await $SystemAPI.reminderCreate({
            resource: 'agent:calendar-feeds',
            assignedTo: $Auth.user.userID,
            payload: { title },
            remindAt: dayEdge(when === 'past' ? 0 : 23),
          })
        ).reminderID
      }

      const feed = resource => ({
        resource,
        startField: 'starts_at',
        endField: '',
        titleField: 'name',
        allDay: false,
        options: { moduleID: mod.moduleID, color: '#09344E', prefilter: '' },
      })

      // A page per feed kind: one calendar cannot be two things at once, and
      // the builder half of the check saves the page it is driving.
      const mk = async (handle, title, resource) =>
        (
          await $ComposeAPI.pageCreate({
            namespaceID,
            title,
            handle,
            visible: true,
            blocks: [
              {
                kind: 'Calendar',
                title,
                xywh: [0, 0, 48, 30],
                options: { defaultView: 'dayGridMonth', header: {}, locale: 'en-gb', feeds: [feed(resource)] },
              },
            ],
          })
        ).pageID

      return {
        namespaceID,
        reminderIDs: Object.values(reminders),
        editPage: await mk('feed_edit', 'Feed edit', 'compose:record'),
        recordPage: await mk('feed_record', 'Record feed', 'compose:record'),
        reminderPage: await mk('feed_reminder', 'Reminder feed', 'system:reminder'),
      }
    },
    { slug: SLUG, reminderTitles: REMINDER_TITLES, recordTitles: RECORD_TITLES },
  )
}

const tearDown = (page, built) =>
  page
    .evaluate(async b => {
      const { $ComposeAPI, $SystemAPI } = window.__human.globals()
      for (const reminderID of b.reminderIDs) {
        await $SystemAPI.reminderDelete({ reminderID }).catch(() => {})
      }
      await $ComposeAPI.namespaceDelete({ namespaceID: b.namespaceID }).catch(() => {})
    }, built)
    .catch(() => {})

/** Build the fixture, run body against it, remove it. */
async function withFixture(page, body) {
  await page.open('/compose/namespaces', { settle: 1200 })
  const built = await build(page)
  try {
    await body(built)
  } finally {
    await tearDown(page, built)
  }
}

const eventTitles = page =>
  page.evaluate(() =>
    [...document.querySelectorAll('.fc-event')].map(e => e.textContent.trim().replace(/\s+/g, ' ')),
  )

const firstEventClass = page =>
  page.evaluate(() => document.querySelector('.fc-event')?.className || '(no event)')

// Mark the event-source Select so a click can name it: the same dialog holds an
// on-click Select whose value ("Open record in the same tab") also says record.
const markSource = page =>
  page.evaluate(() => {
    document.querySelectorAll('[data-check-source]').forEach(e => e.removeAttribute('data-check-source'))
    const label = [...document.querySelectorAll('.p-dialog label')].find(
      l => l.textContent.trim().toLowerCase().replace('*', '') === 'event source',
    )
    if (!label) return '(no Event source field)'
    const sel = label.parentElement?.parentElement?.querySelector('.p-select')
    if (!sel) return '(no select)'
    sel.setAttribute('data-check-source', '')
    return sel.querySelector('.p-select-label')?.textContent?.trim() ?? '(no value)'
  })

// Which record-feed mapping fields the dialog is asking for.
const mappingLabels = page =>
  page.evaluate(() =>
    [...document.querySelectorAll('.p-dialog label')]
      .map(l => l.textContent.trim())
      .filter(t => /module|event start|event end|event title|prefilter/i.test(t)),
  )

const openBlock = page => page.click('.block-toolbox .pi-pencil', { settle: 900 })

drive('a chosen event source survives the block being saved and reopened', async page => {
  await withFixture(page, async built => {
    await page.open(
      `/compose/namespace/${SLUG}/admin/pages/${built.editPage}/builder`,
      { settle: 2500 },
    )
    await openBlock(page)

    const before = await markSource(page)
    check('the block editor offers an event source', !before.startsWith('('), before)
    check('a record feed shows as Records', before === 'Records', before)
    check('a record feed asks for its field mapping', (await mappingLabels(page)).length > 0, '')

    await page.click('[data-check-source]', { settle: 400 })
    await page.click('.p-select-option', { hasText: 'Reminders', settle: 400 })
    check('picking Reminders takes', (await markSource(page)) === 'Reminders', '')

    const mapping = await mappingLabels(page)
    check(
      'a reminder feed asks for no field mapping, having nothing to map',
      mapping.length === 0,
      JSON.stringify(mapping),
    )

    await page.click('.p-dialog-footer button', { hasText: 'Save', settle: 1200 })
    await page.click('[data-testid="editor-actions"] .pi-save', { settle: 2000 })

    await openBlock(page)
    check('reopening the block still shows Reminders', (await markSource(page)) === 'Reminders', '')

    const stored = await page.evaluate(async b => {
      const { $ComposeAPI } = window.__human.globals()
      const p = await $ComposeAPI.pageRead({ namespaceID: b.namespaceID, pageID: b.editPage })
      return (p.blocks || []).find(x => x.kind === 'Calendar').options.feeds[0].resource
    }, built)

    check('the page stores the resource the block runtime dispatches on', stored === 'system:reminder', stored)
  })
})

drive('each event source renders what it reads', async page => {
  await withFixture(page, async built => {
    await page.open(`/compose/namespace/${SLUG}/pages/${built.recordPage}`, { settle: 3500 })

    const records = await eventTitles(page)
    check('a record feed renders its records', records.length === RECORD_TITLES.length, JSON.stringify(records))
    check(
      'each carries its title field',
      RECORD_TITLES.every(t => records.some(r => r.includes(t))),
      JSON.stringify(records),
    )
    check('they are styled as record events', (await firstEventClass(page)).includes('event-record'), '')

    await page.open(`/compose/namespace/${SLUG}/pages/${built.reminderPage}`, { settle: 3500 })

    const reminders = await eventTitles(page)
    check(
      'a reminder feed renders a reminder already passed',
      reminders.some(t => t.includes(REMINDER_TITLES.past)),
      JSON.stringify(reminders),
    )
    check(
      'and one still to come, the range being what decides',
      reminders.some(t => t.includes(REMINDER_TITLES.future)),
      JSON.stringify(reminders),
    )
    check(
      'they are styled as reminder events',
      (await firstEventClass(page)).includes('event-reminder'),
      '',
    )
  })
})
