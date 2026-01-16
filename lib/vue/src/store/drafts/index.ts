import { system } from '@cortezaproject/corteza-js'
import {
  loadAllDraftsFromStorage,
  saveDraftToStorage,
  removeDraftFromStorage,
  clearAllDraftsFromStorage,
} from './storage'

const { Revision } = system

export interface DraftEntry {
  revision: system.Revision
  source: 'local' | 'backend'
}

interface DraftsState {
  drafts: { [key: string]: DraftEntry }
  loading: boolean
  visible: boolean
}

const types = {
  SET_DRAFT: 'SET_DRAFT',
  REMOVE_DRAFT: 'REMOVE_DRAFT',
  SET_LOADING: 'SET_LOADING',
  setVisible: 'setVisible',
  CLEAR_DRAFTS: 'CLEAR_DRAFTS',
}

export default function () {
  return {
    namespaced: true,

    state: (): DraftsState => ({
      drafts: {},
      loading: false,
      visible: false,
    }),

    getters: {
      getDraft: (state: DraftsState) => (changeID: string): DraftEntry | undefined => {
        return state.drafts[changeID]
      },

      hasDraft: (state: DraftsState) => (changeID: string): boolean => {
        return !!state.drafts[changeID]
      },

      getAllDrafts: (state: DraftsState): DraftEntry[] => {
        return Object.values(state.drafts)
      },

      getAllDraftsMap: (state: DraftsState): { [key: string]: DraftEntry } => {
        return state.drafts
      },

      getDraftsByResourceType: (state: DraftsState) => (resourceType: string): DraftEntry[] => {
        return Object.values(state.drafts).filter(
          entry => entry.revision.resource.startsWith(resourceType),
        )
      },

      getDraftsBySource: (state: DraftsState) => (source: 'local' | 'backend'): DraftEntry[] => {
        return Object.values(state.drafts).filter(entry => entry.source === source)
      },

      getDraftsForRecord: (state: DraftsState) => (recordID: string): DraftEntry[] => {
        return Object.values(state.drafts).filter(
          entry => entry.revision.resource.endsWith(`/${recordID}`),
        )
      },

      isLoading: (state: DraftsState): boolean => {
        return state.loading
      },

      visible: (state: DraftsState): boolean => {
        return state.visible
      },
    },

    actions: {
      async init (
        { dispatch }: { dispatch: (action: string, payload?: any) => Promise<void> },
        { resourceType }: { resourceType?: string },
      ): Promise<void> {
        await dispatch('loadAllDrafts', { resourceType })
      },

      async loadAllDrafts (
        { commit, dispatch }: { commit: (mutation: string, payload?: any) => void; dispatch: (action: string, payload?: any) => Promise<void> },
        { resourceType }: { resourceType?: string } = {},
      ): Promise<void> {
        commit(types.SET_LOADING, true)

        try {
          await dispatch('loadLocalDrafts')
        } finally {
          commit(types.SET_LOADING, false)
        }
      },

      loadLocalDrafts ({ commit }: { commit: (mutation: string, payload?: any) => void }): void {
        const localDrafts = loadAllDraftsFromStorage()

        localDrafts.forEach((revisionData, changeID) => {
          const revision = new Revision(revisionData)
          const entry: DraftEntry = {
            revision,
            source: 'local',
          }
          commit(types.SET_DRAFT, { changeID: String(revision.changeID), entry })
        })
      },

      saveDraft (
        { commit }: { commit: (mutation: string, payload?: any) => void },
        { revision }: { revision: system.Revision },
      ): void {
        const changeID = String(revision.changeID)
        saveDraftToStorage(changeID, revision)

        const entry: DraftEntry = {
          revision,
          source: 'local',
        }
        commit(types.SET_DRAFT, { changeID, entry })
      },

      async removeDraft (
        { commit }: { commit: (mutation: string, payload?: any) => void },
        { changeID }: { changeID: string },
      ): Promise<void> {
        removeDraftFromStorage(changeID)
        commit(types.REMOVE_DRAFT, changeID)
      },

      async clearDrafts (
        { commit }: { commit: (mutation: string, payload?: any) => void },
      ): Promise<void> {
        clearAllDraftsFromStorage()
        commit(types.CLEAR_DRAFTS)
      },

      toggleVisibility ({ commit, state }: { commit: (mutation: string, payload?: any) => void; state: DraftsState }): void {
        commit(types.setVisible, !state.visible)
      },
    },

    mutations: {
      [types.SET_DRAFT] (state: DraftsState, { changeID, entry }: { changeID: string; entry: DraftEntry }): void {
        state.drafts = { ...state.drafts, [changeID]: entry }
      },

      [types.REMOVE_DRAFT] (state: DraftsState, changeID: string): void {
        const drafts = { ...state.drafts }
        delete drafts[changeID]
        state.drafts = drafts
      },

      [types.CLEAR_DRAFTS] (state: DraftsState): void {
        state.drafts = {}
      },

      [types.SET_LOADING] (state: DraftsState, loading: boolean): void {
        state.loading = loading
      },

      [types.setVisible] (state: DraftsState, visible: boolean): void {
        state.visible = visible
      },
    },
  }
}

export { getDraftFromStorage } from './storage'
