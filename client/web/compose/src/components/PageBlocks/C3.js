import { compose } from '@cortezaproject/corteza-js'
import PageBlockFactory from './index.js'

const options = {
  Automation: {
    buttons: [
      // { enabled: true, label: 'Enabled script', script: 'someScript' },
      { enabled: true, label: 'Enabled workflow', workflowID: '42' },
      { enabled: false, label: 'Disabled workflow', workflowID: '42' },
    ],
  },

  Chart: {},

  Content: {
    body: '<h1>Lorem ipsum</h1>...',
  },

  File: {
    set: [
      {
        attachmentID: 'string',
      },
    ],
  },
}

export const pageBlockConfigurators = Object.fromEntries([
  // all page block kinds
  ...compose.PageBlockRegistry.keys(),
].map(name => {
  const BlockClass = compose.PageBlockRegistry.get(name)
  const block = new BlockClass({ options: (options[name] || {}) })
  const module = new compose.Module({
    fields: [
      {
        name: 'field1',
        canUpdateRecordValue: true,
        canReadRecordValue: true,
      },
    ],
    canUpdateModule: true,
    canDeleteModule: true,
    canCreateRecord: true,
    canReadRecord: true,
    canUpdateRecord: true,
    canDeleteRecord: true,
    canGrant: true,
  })

  /**
   * @todo This approach still raises errors with store, plugins and mixins
   *       Find a way how to inject them OR rewrite components to handle
   *       c3 environment without crashing
   */
  return [name + 'Configurator', {
    name,
    wip: true,
    group: ['Page blocks', 'Configurator'],
    component: PageBlockFactory,
    controls: [],
    props: {
      block,
      mode: 'configurator',
      namespace: new compose.Namespace(),
      module,
      page: new compose.Page(),
      record: new compose.Record(module),
    },
  }]
}))

export const pageBlockBase = Object.fromEntries([
  // all page block kinds
  ...compose.PageBlockRegistry.keys(),
].map(name => {
  const BlockClass = compose.PageBlockRegistry.get(name)
  const block = new BlockClass({ options: (options[name] || {}) })
  const module = new compose.Module({
    fields: [
      {
        name: 'field1',
        canUpdateRecordValue: true,
        canReadRecordValue: true,
      },
    ],
    canUpdateModule: true,
    canDeleteModule: true,
    canCreateRecord: true,
    canReadRecord: true,
    canUpdateRecord: true,
    canDeleteRecord: true,
    canGrant: true,
  })

  /**
   * @todo This approach still raises errors with store, plugins and mixins
   *       Find a way how to inject them OR rewrite components to handle
   *       c3 environment without crashing
   */
  return [name, {
    name,
    wip: true,
    group: ['Page blocks', 'View-mode'],
    component: PageBlockFactory,
    controls: [],
    props: {
      block,
      namespace: new compose.Namespace(),
      module,
      page: new compose.Page(),
      record: new compose.Record(module),
    },
  }]
}))
