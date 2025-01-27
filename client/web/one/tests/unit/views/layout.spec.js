/* eslint-disable no-unused-expressions */
import { expect } from 'chai'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import Layout from 'corteza-webapp-one/src/views/Layout'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import Vuex from 'vuex'
import sinon from 'sinon'

describe('Layout.vue', () => {
  let localVue
  let store
  let applicationsModule

  beforeEach(() => {
    localVue = createLocalVue()
    localVue.use(BootstrapVue)
    localVue.use(PortalVue)
    localVue.use(Vuex)

    // Mock the Vuex store
    applicationsModule = {
      namespaced: true,
      actions: {
        load: sinon.stub().resolves(),
      },
    }

    store = new Vuex.Store({
      modules: {
        applications: applicationsModule,
      },
    })
  })

  afterEach(() => {
    sinon.restore()
  })

  it('renders', () => {
    // Mock the CAppSelector component
    const CAppSelector = {
      name: 'c-app-selector',
      render: () => {},
    }

    // Register the component locally
    localVue.component('c-app-selector', CAppSelector)

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
      },
      stubs: [
        'router-view',
        'portal-target',
        'c-topbar',
        'c-prompts',
        'c-loader-logo',
        'c-app-selector',
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
