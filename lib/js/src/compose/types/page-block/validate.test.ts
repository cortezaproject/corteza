import { expect } from 'chai'

import { PageBlockMaker, PageBlockRegistry, BlockIssue } from './index'
import { PageBlockCalendar } from './calendar'

/**
 * A block's own class is the one place its required options are written down —
 * the builder reads `validate()` to mark the configurator, badge the block on
 * the grid and refuse the page save. These pin what each kind demands, so a
 * kind cannot quietly stop asking for the option it cannot render without.
 */

const optionsOf = (kind: string, options: Record<string, unknown> = {}): string[] =>
  PageBlockMaker({ kind, options } as { kind: string })
    .validate()
    .map((i: BlockIssue) => i.option)

const MODULE = '100001'

describe('page block required options', () => {
  it('names what a freshly added block of each kind is missing', () => {
    const missing = Object.fromEntries(
      [...PageBlockRegistry.keys()].map(kind => [kind, optionsOf(kind)]),
    )

    expect(missing).to.deep.eq({
      Automation: ['buttons'],
      AgentChat: ['allowedAgentIDs'],
      Calendar: ['feeds'],
      Chart: ['chartID'],
      ChatbotInbox: ['chatbotIDs'],
      Comment: ['moduleID', 'contentField'],
      Content: ['body'],
      Geometry: ['feeds'],
      IFrame: ['src'],
      Metric: ['metrics'],
      Navigation: ['navigationItems'],
      Progress: [],
      Record: [],
      RecordList: ['moduleID'],
      RecordOrganizer: ['moduleID'],
      RecordRevisions: [],
      Tabs: ['tabs'],
      File: [],
    })
  })

  // HumanID casts a blank to NoID, so an unset module is '0', not ''. A plain
  // falsy check passes it and the guard never fires.
  it('reads NoID as unset, not as a module', () => {
    expect(optionsOf('RecordList', { moduleID: '0' })).to.deep.eq(['moduleID'])
    expect(optionsOf('RecordList', { moduleID: MODULE })).to.deep.eq([])
  })

  it('accepts either an iframe source or a source field', () => {
    expect(optionsOf('IFrame', { src: 'https://example.test' })).to.deep.eq([])
    expect(optionsOf('IFrame', { srcField: 'link' })).to.deep.eq([])
  })

  it('asks each metric for its own module', () => {
    expect(optionsOf('Metric', { metrics: [{ moduleID: MODULE }, {}] })).to.deep.eq([
      'metrics.1.moduleID',
    ])
  })

  it('asks a calendar record feed for a module and a start, and a reminder feed for neither', () => {
    const record = PageBlockCalendar.feedResources.record
    const reminder = PageBlockCalendar.feedResources.reminder

    expect(optionsOf('Calendar', { feeds: [{ resource: record }] })).to.deep.eq([
      'feeds.0.moduleID',
      'feeds.0.startField',
    ])
    expect(optionsOf('Calendar', { feeds: [{ resource: reminder }] })).to.deep.eq([])
  })

  it('asks each geometry feed for a module and a geometry field', () => {
    expect(optionsOf('Geometry', { feeds: [{}] })).to.deep.eq([
      'feeds.0.moduleID',
      'feeds.0.geometryField',
    ])
  })

  it('names an automation button that reaches nothing', () => {
    expect(optionsOf('Automation', { buttons: [{ workflowID: '1' }, {}] })).to.deep.eq([
      'buttons.1',
    ])
  })

  it('carries an i18n key for every issue it reports', () => {
    for (const kind of PageBlockRegistry.keys()) {
      for (const issue of PageBlockMaker({ kind } as { kind: string }).validate()) {
        expect(issue.labelKey, `${kind}.${issue.option}`).to.match(/^block\.issue\./)
      }
    }
  })
})
