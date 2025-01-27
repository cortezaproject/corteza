/* eslint-disable no-unused-expressions */
/* global jest */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import Layout from 'corteza-webapp-admin/src/views/Layout'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import Vuex from 'vuex'

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

  it('renders', () => {
    // Mock the CTheMainNav component
    const CTheMainNav = {
      name: 'c-the-main-nav',
      render: () => {},
    }

    // Register the component locally
    localVue.component('c-the-main-nav', CTheMainNav)

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
        'c-the-main-nav',
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
