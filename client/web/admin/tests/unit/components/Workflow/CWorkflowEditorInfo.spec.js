/* eslint-disable no-unused-expressions */
/* global jest */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import BootstrapVue from 'bootstrap-vue'
import CWorkflowEditorInfo from 'corteza-webapp-admin/src/components/Workflow/CWorkflowEditorInfo'
import { url } from '@cortezaproject/corteza-vue'

// tests/setup.js stubs corteza-vue with components only; this component needs
// the handle helpers and the app URL builder (jest hoists these above imports)
jest.mock('@cortezaproject/corteza-vue', () => ({
  handle: { IsValid: () => true, handleState: () => null },
  url: {
    MakeAppURL: jest.fn(({ app, path }) => `https://host.tld/process/${app}/${path}`),
  },
}), { virtual: true })

jest.mock('@cortezaproject/corteza-js', () => ({ NoID: '0' }), { virtual: true })

describe('components/Workflow/CWorkflowEditorInfo.vue', () => {
  let localVue

  beforeEach(() => {
    localVue = createLocalVue()
    localVue.use(BootstrapVue)
    url.MakeAppURL.mockClear()
  })

  it('opens the builder through the webapp base so a path prefix is kept', () => {
    const open = jest.spyOn(window, 'open').mockImplementation(() => null)

    const wrapper = shallowMount(CWorkflowEditorInfo, {
      localVue,
      propsData: { workflow: { workflowID: '42', handle: 'wf', meta: {} } },
      mocks: { $t: (key) => key },
      stubs: ['c-input-confirm', 'c-button-submit', 'c-input-checkbox', 'font-awesome-icon'],
    })

    wrapper.vm.openWorkflowBuilder()

    expect(url.MakeAppURL.mock.calls[0][0]).to.deep.equal({ app: 'workflow', path: '42/edit' })
    expect(open.mock.calls[0][0]).to.equal('https://host.tld/process/workflow/42/edit')
    expect(open.mock.calls[0][1]).to.equal('_blank')

    open.mockRestore()
  })
})
