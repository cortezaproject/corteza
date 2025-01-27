/* eslint-disable no-unused-expressions */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import Layout from 'corteza-webapp-compose/src/views/Layout'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import Vuex from 'vuex'

describe('Layout.vue', () => {
  it('renders', () => {
    const localVue = createLocalVue()
    localVue.use(BootstrapVue)
    localVue.use(PortalVue)
    localVue.use(Vuex)

    // Create a mock store
    const store = new Vuex.Store({
      modules: {
        ui: {
          namespaced: true,
          state: {
            namespaceSlug: 'test-namespace',
            pageHandle: 'test-page',
            layoutHandle: 'test-layout',
          },
          getters: {
            namespaceSlug: state => state.namespaceSlug,
            pageHandle: state => state.pageHandle,
            layoutHandle: state => state.layoutHandle,
          },
          actions: {
            setNamespaceSlug: () => {},
          },
        },
      },
    })

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
        $route: {
          params: {
            slug: 'test-namespace',
          },
        },
        $root: {
          $on: () => {},
        },
        $t: (key) => key,
        textDirectionality: () => 'ltr',
      },
      stubs: ['router-view', 'portal-target'],
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
