/* global jest */

jest.mock('@cortezaproject/corteza-js', () => ({}), { virtual: true })
jest.mock('@cortezaproject/corteza-vue', () => ({
  components: {
    CPrompts: {
      name: 'c-prompts',
      render: () => {},
    },
    CTopbar: {
      name: 'c-topbar',
      render: () => {},
    },
    CLoaderLogo: {
      name: 'c-loader-logo',
      render: () => {},
    },
  },
}), { virtual: true })
