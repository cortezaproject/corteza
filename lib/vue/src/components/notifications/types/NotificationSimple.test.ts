import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import NotificationSimple from './NotificationSimple.vue'

// A notification's description is HTML somebody else wrote, and the server
// stores it as it came, so nothing but this component stands between its
// author and whoever it is shown to.
describe('NotificationSimple', () => {
  const render = (description: string) =>
    mount(NotificationSimple, {
      props: { notification: { config: { title: 'Hello', description } } },
    })

  it('shows the markup its author wrote', () => {
    expect(render('<b>done</b>').html()).toContain('<b>done</b>')
  })

  it('lets nothing run', () => {
    const html = render(
      '<img src=x onerror="window.pwned = 1"><script>window.pwned = 1</script>',
    ).html()
    expect(html).not.toContain('onerror')
    expect(html).not.toContain('<script')
  })

  it('renders nothing in particular when there is no description', () => {
    expect(render('').html()).toContain('Hello')
  })
})
