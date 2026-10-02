/* eslint-disable no-unused-expressions */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import BootstrapVue from 'bootstrap-vue'
import CSystemEmailServer from 'corteza-webapp-admin/src/components/Settings/Email/CSystemEmailServer'

describe('components/Settings/Email/CSystemEmailServer.vue', () => {
  let localVue

  const mount = (value) => shallowMount(CSystemEmailServer, {
    localVue,
    propsData: { value, disabled: false },
    mocks: { $t: (key) => key },
    stubs: ['c-hint', 'c-input-checkbox', 'c-input-confirm', 'c-button-submit', 'font-awesome-icon'],
  })

  beforeEach(() => {
    localVue = createLocalVue()
    localVue.use(BootstrapVue)
  })

  const padded = {
    host: '  smtp.example.com ',
    port: 587,
    user: ' mailer ',
    pass: ' keep spaces ',
    from: ' noreply@example.com ',
    tlsInsecure: false,
    tlsServerName: ' smtp.example.com ',
  }

  it('trims host and addresses when submitting', () => {
    const wrapper = mount(padded)

    wrapper.vm.submit()

    const [server] = wrapper.emitted('submit')[0]
    expect(server.host).to.equal('smtp.example.com')
    expect(server.user).to.equal('mailer')
    expect(server.from).to.equal('noreply@example.com')
    expect(server.tlsServerName).to.equal('smtp.example.com')
    // the password is sent as typed
    expect(server.pass).to.equal(' keep spaces ')
    expect(server.port).to.equal(587)
  })

  it('trims host and addresses when testing the connection', () => {
    const wrapper = mount(padded)

    wrapper.vm.smtpConnectionCheck()

    const [server] = wrapper.emitted('smtpConnectionCheck')[0]
    expect(server.host).to.equal('smtp.example.com')
    expect(server.from).to.equal('noreply@example.com')
  })
})
