import { expect } from 'chai'
import toast from './toast'

// the mixin's methods are plain functions bound to a fake component
function handler (title: string, calls: Array<[string, string]>) {
  const vm = {
    ...(toast as any).methods,
    $t: (k: string) => k,
    toastDanger (msg: string, t: string) { calls.push([msg, t]) },
  }

  return vm.toastErrorHandler(title)
}

describe('toastErrorHandler', () => {
  it('uses the given title and the error message', () => {
    const calls: Array<[string, string]> = []
    handler('Could not save record', calls)({ message: 'something broke' })
    expect(calls).to.deep.eq([['Something broke', 'Could not save record']])
  })

  it('prefers the title carried by the error (workflow error step)', () => {
    const calls: Array<[string, string]> = []
    handler('Could not save record', calls)({ message: 'order is not ready', meta: { title: 'Order check' } })
    expect(calls).to.deep.eq([['Order is not ready', 'Order check']])
  })

  it('falls back to the title when the error carries none', () => {
    const calls: Array<[string, string]> = []
    handler('Could not save record', calls)({ message: 'nope', meta: {} })
    expect(calls).to.deep.eq([['Nope', 'Could not save record']])
  })
})
