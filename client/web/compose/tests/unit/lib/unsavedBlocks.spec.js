/* global jest, expect */
import Vue from 'vue'
import Vuex from 'vuex'
import { shallowMount, createLocalVue } from '@vue/test-utils'
import ui from 'corteza-webapp-compose/src/store/ui'
import unsavedBlocks from 'corteza-webapp-compose/src/mixins/unsavedBlocks'

Vue.use(Vuex)

const localVue = createLocalVue()
localVue.use(Vuex)

const Host = {
  mixins: [unsavedBlocks],
  render: h => h('div'),
}

describe('mixins/unsavedBlocks', () => {
  let store, confirmSpy, addSpy, removeSpy

  beforeEach(() => {
    store = new Vuex.Store({ modules: { ui: ui({}) } })
    confirmSpy = jest.spyOn(window, 'confirm')
    addSpy = jest.spyOn(window, 'addEventListener')
    removeSpy = jest.spyOn(window, 'removeEventListener')
  })

  afterEach(() => {
    jest.restoreAllMocks()
  })

  const mount = () => shallowMount(Host, { localVue, store, mocks: { $t: (k) => k } })

  it('lets the user leave without asking when nothing is unsaved', () => {
    const wrapper = mount()
    expect(wrapper.vm.confirmUnsavedBlocks()).toBe(true)
    expect(confirmSpy).not.toHaveBeenCalled()
  })

  it('asks before leaving with unsaved blocks and stays when the user declines', async () => {
    const wrapper = mount()
    await store.dispatch('ui/setBlockUnsaved', { id: 'list', unsaved: true })

    confirmSpy.mockReturnValue(false)
    expect(wrapper.vm.confirmUnsavedBlocks()).toBe(false)
    expect(confirmSpy).toHaveBeenCalledWith('general:record.unsavedChanges')
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(true)
  })

  it('leaves and forgets the unsaved blocks when the user agrees', async () => {
    const wrapper = mount()
    await store.dispatch('ui/setBlockUnsaved', { id: 'list', unsaved: true })

    confirmSpy.mockReturnValue(true)
    expect(wrapper.vm.confirmUnsavedBlocks()).toBe(true)
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(false)
  })

  it('warns the browser on reload only while something is unsaved', async () => {
    const wrapper = mount()
    await store.dispatch('ui/setBlockUnsaved', { id: 'list', unsaved: true })
    await wrapper.vm.$nextTick()
    expect(addSpy).toHaveBeenCalledWith('beforeunload', wrapper.vm.warnBeforeUnload)

    await store.dispatch('ui/setBlockUnsaved', { id: 'list', unsaved: false })
    await wrapper.vm.$nextTick()
    expect(removeSpy).toHaveBeenCalledWith('beforeunload', wrapper.vm.warnBeforeUnload)

    const e = { preventDefault: jest.fn(), returnValue: undefined }
    wrapper.vm.warnBeforeUnload(e)
    expect(e.preventDefault).toHaveBeenCalled()
    expect(e.returnValue).toBe('')
  })
})
