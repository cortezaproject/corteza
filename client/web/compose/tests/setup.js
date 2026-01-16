/* global jest */

jest.mock('@cortezaproject/corteza-js', () => ({
  compose: {
    Namespace: class {},
    Module: class {},
    ModuleField: class {},
    Record: class {},
  },
}), { virtual: true })
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
    CExtendSession: {
      name: 'c-extend-session',
      render: () => {},
    },
    CTranslationModal: {
      name: 'c-translation-modal',
      render: () => {},
    },
    CNamespaceSidebar: {
      name: 'c-namespace-sidebar',
      render: () => {},
    },
    CNotificationSidebar: {
      name: 'c-notification-sidebar',
      render: () => {},
    },
    CDraftSidebar: {
      name: 'c-draft-sidebar',
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
