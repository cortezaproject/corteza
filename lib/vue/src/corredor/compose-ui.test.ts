import { describe, it, expect, vi } from 'vitest'
import { compose } from '@planetcrust/human-js'
import { ComposeUIHelper } from './compose-ui'

const namespace = new compose.Namespace({ namespaceID: '1', slug: 'agent-sandbox' })
const module = new compose.Module(
  { moduleID: '2', namespaceID: '1', fields: [{ name: 'name', kind: 'String' }] },
  namespace,
)
const recordPage = new compose.Page({ pageID: '3', moduleID: '2', namespaceID: '1' })
const otherPage = new compose.Page({ pageID: '4', moduleID: '9', namespaceID: '1' })
const record = new compose.Record(module, { recordID: '5' })

function helper(overrides = {}) {
  const routePusher = vi.fn()
  const toast = { success: vi.fn(), warning: vi.fn() }

  const ui = new ComposeUIHelper({
    $namespace: namespace,
    $module: module,
    $record: record,
    pages: () => [otherPage, recordPage],
    toast,
    routePusher,
    ...overrides,
  })

  return { ui, routePusher, toast }
}

describe('ComposeUIHelper navigation', () => {
  it('opens the record page of the record module', () => {
    const { ui, routePusher } = helper()
    ui.gotoRecordViewer()

    expect(routePusher).toHaveBeenCalledWith({
      name: 'page.record',
      params: { slug: 'agent-sandbox', pageID: '3', recordID: '5' },
    })
  })

  it('opens the record page in edit mode', () => {
    const { ui, routePusher } = helper()
    ui.gotoRecordEditor()

    expect(routePusher).toHaveBeenCalledWith({
      name: 'page.record',
      params: { slug: 'agent-sandbox', pageID: '3', recordID: '5' },
      query: { edit: '1' },
    })
  })

  it('takes the record it is handed over the one in context', () => {
    const { ui, routePusher } = helper()
    ui.gotoRecordViewer(new compose.Record(module, { recordID: '7' }))

    expect(routePusher.mock.calls[0][0].params.recordID).toBe('7')
  })

  it('refuses a module with no record page', () => {
    const { ui, routePusher } = helper({ pages: () => [otherPage] })

    expect(() => ui.gotoRecordViewer()).toThrow('record page does not exist')
    expect(routePusher).not.toHaveBeenCalled()
  })

  it('refuses a record that has not been saved', () => {
    const { ui } = helper({ $record: undefined })

    expect(() => ui.gotoRecordViewer()).toThrow()
  })
})

describe('ComposeUIHelper messages', () => {
  it('shows a success and a warning message', () => {
    const { ui, toast } = helper()

    ui.success('saved')
    ui.warning('careful')

    expect(toast.success).toHaveBeenCalledWith('saved')
    expect(toast.warning).toHaveBeenCalledWith('careful')
  })

  it('is safe without a toast', () => {
    const { ui } = helper({ toast: undefined })

    expect(() => ui.success('saved')).not.toThrow()
  })
})
