/* global expect */
import Vue from 'vue'
import Vuex from 'vuex'
import ui from 'corteza-webapp-compose/src/store/ui'

Vue.use(Vuex)

describe('store/ui unsaved blocks', () => {
  let store

  beforeEach(() => {
    store = new Vuex.Store({ modules: { ui: ui({}) } })
  })

  it('reports no unsaved blocks by default', () => {
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(false)
  })

  it('tracks blocks that report unsaved changes and forgets them once saved', async () => {
    await store.dispatch('ui/setBlockUnsaved', { id: 'list-a', unsaved: true })
    await store.dispatch('ui/setBlockUnsaved', { id: 'list-b', unsaved: true })
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(true)

    await store.dispatch('ui/setBlockUnsaved', { id: 'list-a', unsaved: false })
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(true)

    await store.dispatch('ui/setBlockUnsaved', { id: 'list-b', unsaved: false })
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(false)
  })

  it('does not duplicate a block that reports twice', async () => {
    await store.dispatch('ui/setBlockUnsaved', { id: 'list-a', unsaved: true })
    await store.dispatch('ui/setBlockUnsaved', { id: 'list-a', unsaved: true })
    expect(store.state.ui.unsavedBlocks).toEqual(['list-a'])
  })

  it('ignores blocks without an id', async () => {
    await store.dispatch('ui/setBlockUnsaved', { id: undefined, unsaved: true })
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(false)
  })

  it('clears everything when the page is left', async () => {
    await store.dispatch('ui/setBlockUnsaved', { id: 'list-a', unsaved: true })
    await store.dispatch('ui/clearUnsavedBlocks')
    expect(store.getters['ui/hasUnsavedBlocks']).toBe(false)
  })
})
