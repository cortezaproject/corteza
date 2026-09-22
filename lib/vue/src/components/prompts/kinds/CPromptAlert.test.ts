import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import CPromptAlert from './CPromptAlert.vue'

const ButtonStub = { props: { label: String }, template: '<button>{{ label }}</button>' }

// Every prompt kind renders its message as HTML, and a workflow author wrote
// it; the base component is where that is made safe, so the check belongs on
// one of the kinds that show it.
describe('CPromptAlert', () => {
  const render = (message: string) =>
    mount(CPromptAlert, {
      // a prompt payload is automation.Vars: every value is wrapped
      props: { payload: { message: { '@value': message, '@type': 'String' } } },
      global: { stubs: { Button: ButtonStub }, mocks: { $t: (k: string) => k } },
    })

  it('shows the markup the author wrote', () => {
    expect(render('<b>careful</b>').html()).toContain('<b>careful</b>')
  })

  it('lets nothing run', () => {
    const html = render('<img src=x onerror="window.pwned = 1">').html()
    expect(html).not.toContain('onerror')
  })
})
