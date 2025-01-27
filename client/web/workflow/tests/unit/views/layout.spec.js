/* eslint-disable no-unused-expressions */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import Layout from '../../../src/views/Layout'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import sinon from 'sinon'

describe('Layout.vue', () => {
  let localVue

  beforeEach(() => {
    localVue = createLocalVue()
    localVue.use(BootstrapVue)
    localVue.use(PortalVue)
  })

  afterEach(() => {
    sinon.restore()
  })

  it('renders', () => {
    const wrapper = shallowMount(Layout, {
      localVue,
      mocks: {
        $auth: {
          user: {},
        },
        $Settings: {
          get: () => ({}),
          attachment: () => '',
        },
        $t: (key) => key,
      },
      stubs: [
        'router-view',
        'portal-target',
        'c-topbar',
        'c-sidebar',
        'c-permissions-modal',
      ],
    })

    expect(wrapper.find('div').classes()).to.include.members([
      'd-flex',
      'flex-column',
      'w-100',
      'vh-100',
      'overflow-hidden',
    ])
  })
})
