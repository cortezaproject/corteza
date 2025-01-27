/* eslint-disable no-unused-expressions */
/* global jest */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import Layout from '../../../src/views/Layout'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import Vuex from 'vuex'
import sinon from 'sinon'

describe('Layout.vue', () => {
  let localVue
  let store
  let rbacModule

  beforeEach(() => {
    localVue = createLocalVue()
    localVue.use(BootstrapVue)
    localVue.use(PortalVue)
    localVue.use(Vuex)

    // Mock the Vuex store
    rbacModule = {
      namespaced: true,
      getters: {
        can: () => () => true,
      },
    }

    store = new Vuex.Store({
      modules: {
        rbac: rbacModule,
      },
    })
  })

  afterEach(() => {
    sinon.restore()
  })

  it('renders', () => {
    const wrapper = shallowMount(Layout, {
      localVue,
      store,
      mocks: {
        $auth: {
          user: {},
        },
        $Settings: {
          get: () => ({}),
          attachment: () => '',
        },
        $t: (key) => key,
        textDirectionality: () => 'ltr',
        $root: {
          $on: jest.fn(),
        },
      },
      stubs: [
        'router-view',
        'portal-target',
        'c-topbar',
        'c-sidebar',
        'c-prompts',
        'c-permissions-modal',
      ],
    })

    expect(wrapper.find('div').exists()).to.be.true
    expect(wrapper.find('div').classes()).to.include.members([
      'd-flex',
      'flex-column',
      'w-100',
      'vh-100',
      'overflow-hidden',
    ])
  })
})
