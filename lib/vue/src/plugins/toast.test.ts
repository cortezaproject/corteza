import { describe, expect, it } from 'vitest'
import { ToastPlugin } from './toast'

function install() {
  const added: any[] = []
  const app: any = {
    config: { globalProperties: { $toast: { add: (o: any) => added.push(o) } } },
    provide: () => {},
  }

  ToastPlugin.install(app)

  return { added, $toast: app.config.globalProperties.$toast }
}

describe('toastErrorHandler', () => {
  it('shows the prefix, the message and the given title', () => {
    const { added, $toast } = install()
    $toast.toastErrorHandler('Could not save', 'Record')({ message: 'nope' })
    expect(added).toEqual([
      { severity: 'error', summary: 'Record', detail: 'Could not save: nope', life: 7000 },
    ])
  })

  it('prefers the title carried by the error (workflow error step)', () => {
    const { added, $toast } = install()
    $toast.toastErrorHandler(
      'Could not save',
      'Record',
    )({ message: 'order is not ready', meta: { title: 'Order check' } })
    expect(added[0].summary).toBe('Order check')
    expect(added[0].detail).toBe('Could not save: order is not ready')
  })

  it('falls back to the given title when the error carries none', () => {
    const { added, $toast } = install()
    $toast.toastErrorHandler('Could not save', 'Record')({ message: 'nope', meta: {} })
    expect(added[0].summary).toBe('Record')
  })
})
