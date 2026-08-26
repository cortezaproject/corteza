/* eslint-disable no-unused-expressions */
/* global jest */
import { expect } from 'chai'
import { createLocalVue, shallowMount } from '@vue/test-utils'
import BootstrapVue from 'bootstrap-vue'
import PortalVue from 'portal-vue'
import Vuex from 'vuex'
import UserList from 'corteza-webapp-admin/src/views/System/User/List'

jest.mock('axios', () => ({ isCancel: () => false }))

const ResourceListStub = {
  name: 'CResourceList',
  render (h) {
    const children = this.$scopedSlots.header ? this.$scopedSlots.header() : []
    if (this.$scopedSlots.actions) {
      const can = this.$store.getters['rbac/can']
      children.push(this.$scopedSlots.actions({
        item: { userID: 'user-fixture-001', canDeleteUser: can('system/', 'user.delete') },
      }))
    }
    return h('div', children)
  },
}

const ImportStub = {
  name: 'CUserImportModal',
  template: '<button data-test-id="button-import">Import</button>',
}

const ExportStub = {
  name: 'CUserExportModal',
  template: '<button data-test-id="button-export">Export</button>',
}

const DropdownStub = {
  name: 'BDropdown',
  template: '<div data-test-id="user-actions"><slot name="button-content"/><slot /></div>',
}

const ConfirmStub = {
  name: 'CInputConfirm',
  template: '<button data-test-id="button-delete-user">Delete</button>',
}

function mountUserList (can) {
  const localVue = createLocalVue()
  localVue.use(BootstrapVue)
  localVue.use(PortalVue)
  localVue.use(Vuex)

  const store = new Vuex.Store({
    modules: {
      rbac: {
        namespaced: true,
        getters: {
          can: () => can,
        },
      },
    },
  })

  return shallowMount(UserList, {
    localVue,
    store,
    mocks: {
      $route: { query: {} },
      $router: { replace: jest.fn(), push: jest.fn() },
      $root: { $emit: jest.fn() },
      $t: key => key,
      $auth: { accessToken: '' },
      $SystemAPI: {},
    },
    stubs: {
      CResourceList: ResourceListStub,
      CUserImportModal: ImportStub,
      CUserExportModal: ExportStub,
      'c-content-header': true,
      'c-permissions-button': true,
      'c-corredor-manual-buttons': true,
      'c-resource-list-status-filter': true,
      'b-dropdown': DropdownStub,
      'c-input-confirm': ConfirmStub,
    },
  })
}

describe('System user list permissions', () => {
  it('shows management controls to an administrator', () => {
    const wrapper = mountUserList(() => true)

    expect(wrapper.find('[data-test-id="button-new-user"]').exists()).to.equal(true)
    expect(wrapper.find('[data-test-id="button-import"]').exists()).to.equal(true)
    expect(wrapper.find('[data-test-id="button-export"]').exists()).to.equal(true)
    expect(wrapper.find('[data-test-id="user-actions"]').exists()).to.equal(true)
    expect(wrapper.find('[data-test-id="button-delete-user"]').exists()).to.equal(true)
  })

  it('keeps the list readable but hides management controls for a restricted user', () => {
    const wrapper = mountUserList((_resource, operation) => operation === 'users.search')

    expect(wrapper.find('[data-test-id="button-new-user"]').exists()).to.equal(false)
    expect(wrapper.find('[data-test-id="button-import"]').exists()).to.equal(false)
    expect(wrapper.find('[data-test-id="button-export"]').exists()).to.equal(false)
    expect(wrapper.find('[data-test-id="user-actions"]').exists()).to.equal(false)
    expect(wrapper.find('[data-test-id="button-delete-user"]').exists()).to.equal(false)
  })
})
