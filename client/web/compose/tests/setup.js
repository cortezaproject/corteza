/* global jest */

jest.mock('@cortezaproject/corteza-js', () => ({}), { virtual: true })
jest.mock('@cortezaproject/corteza-vue', () => ({
  components: {
    CToaster: jest.fn(),
    CPrompts: {
      name: 'c-prompts',
      render: () => {},
    },
    CPermissionsModal: {
      name: 'c-permissions-modal',
      render: () => {},
    },
    CTopbar: {
      name: 'c-topbar',
      render: () => {},
    },
    CSidebar: {
      name: 'c-sidebar',
      render: () => {},
    },
  },
  mixins: {
    corredor: {
      methods: {
        triggerScript: jest.fn(),
      },
    },
  },
}), { virtual: true })
