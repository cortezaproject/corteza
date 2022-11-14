/* eslint-disable no-unused-expressions */
/* eslint-disable no-unused-vars */
import { expect } from 'chai'
import sinon from 'sinon'
import Layout from 'corteza-webapp-one/src/views/Layout'
import { shallowMount } from 'corteza-webapp-one/tests/lib/helpers'
import fp from 'flush-promises'

describe('/views/Layout.vue', () => {
  afterEach(() => {
    sinon.restore()
  })

  let $auth, $Settings

  beforeEach(() => {
    $auth = {
      user: {},
    }

    $Settings = {
      get: () => ({}),
      attachment: () => '',
    }
  })

  const mountLayout = (opt) => shallowMount(Layout, {
    mocks: { $auth, $Settings },
    ...opt,
  })

  describe('Init', () => {
    it('It renders', async () => {
      const wrapper = mountLayout()

      await fp()
      expect(wrapper.find('div')).to.exist
    })
  })
})
