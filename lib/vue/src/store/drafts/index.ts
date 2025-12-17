import { system } from '@cortezaproject/corteza-js'
import {
  loadAllDraftsFromStorage,
  saveDraftToStorage,
  removeDraftFromStorage,
} from './storage'

const { Revision } = system

interface RevisionChange {
  key: string
  old: unknown[]
  new: unknown[]
}

export interface DraftContext {
  namespaceID: string
  moduleID: string
  recordID?: string
  isNew?: boolean
}

export interface DraftEntry {
  revision: system.Revision
  source: 'local' | 'backend'
  context?: DraftContext
}

interface DraftsState {
  drafts: Map<string, DraftEntry>
  syncing: Set<string>
  loading: boolean
  syncIntervalId: ReturnType<typeof setInterval> | null
  composeAPI: any | null
  systemAPI: any | null
}

const types = {
  SET_DRAFT: 'SET_DRAFT',
  REMOVE_DRAFT: 'REMOVE_DRAFT',
  SET_SYNCING: 'SET_SYNCING',
  SET_LOADING: 'SET_LOADING',
  SET_SYNC_INTERVAL: 'SET_SYNC_INTERVAL',
  SET_APIS: 'SET_APIS',
}

const SYNC_INTERVAL_MS = 60000

export default function () {
  return {
    namespaced: true,

    state: (): DraftsState => ({
      drafts: new Map(),
      syncing: new Set(),
      loading: false,
      syncIntervalId: null,
      composeAPI: null,
      systemAPI: null,
    }),

    getters: {
      getDraft: (state: DraftsState) => (changeID: string): DraftEntry | undefined => {
        return state.drafts.get(changeID)
      },

      hasDraft: (state: DraftsState) => (changeID: string): boolean => {
        return state.drafts.has(changeID)
      },

      getAllDrafts: (state: DraftsState): DraftEntry[] => {
        return Array.from(state.drafts.values())
      },

      getAllDraftsMap: (state: DraftsState): Map<string, DraftEntry> => {
        return state.drafts
      },

      getDraftsByResourceType: (state: DraftsState) => (resourceType: string): DraftEntry[] => {
        return Array.from(state.drafts.values()).filter(
          entry => entry.revision.resourceType === resourceType,
        )
      },

      getDraftsBySource: (state: DraftsState) => (source: 'local' | 'backend'): DraftEntry[] => {
        return Array.from(state.drafts.values()).filter(entry => entry.source === source)
      },

      getDraftsForRecord: (state: DraftsState) => (recordID: string): DraftEntry[] => {
        return Array.from(state.drafts.values()).filter(
          entry => entry.revision.resourceID === recordID,
        )
      },

      isLoading: (state: DraftsState): boolean => {
        return state.loading
      },

      isSyncing: (state: DraftsState) => (changeID: string): boolean => {
        return state.syncing.has(changeID)
      },

      isSyncRunning: (state: DraftsState): boolean => {
        return state.syncIntervalId !== null
      },
    },

    actions: {
      async init (
        { commit, dispatch }: { commit: Function; dispatch: Function },
        { composeAPI, systemAPI, resourceType }: { composeAPI: any; systemAPI: any; resourceType?: string },
      ): Promise<void> {
        commit(types.SET_APIS, { composeAPI, systemAPI })
        await dispatch('loadAllDrafts', { api: systemAPI, resourceType })
        dispatch('startBackgroundSync')
      },

      async loadAllDrafts (
        { commit, dispatch }: { commit: Function; dispatch: Function },
        { api, resourceType, resourceID }: { api?: any; resourceType?: string; resourceID?: string } = {},
      ): Promise<void> {
        commit(types.SET_LOADING, true)

        try {
          dispatch('loadLocalDrafts')
          if (api) {
            await dispatch('fetchBackendDrafts', { api, resourceType, resourceID })
          }
        } finally {
          commit(types.SET_LOADING, false)
        }
      },

      loadLocalDrafts ({ commit }: { commit: Function }): void {
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

      async fetchBackendDrafts (
        { commit }: { commit: Function },
        { api, resourceType, resourceID }: { api: any; resourceType?: string; resourceID?: string },
      ): Promise<void> {
        try {
          const response = await api.revisionList({
            status: 'draft',
            resourceType,
            resourceID,
          })

          const { set = [] } = response || {}

          set.forEach((revisionData: any) => {
            const revision = new Revision(revisionData)
            const entry: DraftEntry = {
              revision,
              source: 'backend',
            }
            commit(types.SET_DRAFT, { changeID: String(revision.changeID), entry })
          })
        } catch (e) {
          console.error('Failed to fetch backend drafts:', e)
        }
      },

      saveDraft (
        { commit }: { commit: Function },
        { revision, context }: { revision: system.Revision; context?: DraftContext },
      ): void {
        const changeID = String(revision.changeID)
        saveDraftToStorage(changeID, revision)

        const entry: DraftEntry = {
          revision,
          source: 'local',
          context,
        }
        commit(types.SET_DRAFT, { changeID, entry })
      },

      async removeDraft (
        { commit, state }: { commit: Function; state: DraftsState },
        { changeID }: { changeID: string },
      ): Promise<void> {
        const entry = state.drafts.get(changeID)
        removeDraftFromStorage(changeID)

        if (entry?.source === 'backend' && state.systemAPI) {
          try {
            await state.systemAPI.revisionDelete({ revisionID: changeID })
          } catch (e) {
            console.error('Failed to delete backend draft:', changeID, e)
          }
        }

        commit(types.REMOVE_DRAFT, changeID)
      },

      startBackgroundSync (
        { commit, dispatch, state }: { commit: Function; dispatch: Function; state: DraftsState },
      ): void {
        if (state.syncIntervalId) return

        const intervalId = setInterval(() => {
          dispatch('syncAllLocalDrafts')
        }, SYNC_INTERVAL_MS)

        commit(types.SET_SYNC_INTERVAL, intervalId)
      },

      stopBackgroundSync (
        { commit, state }: { commit: Function; state: DraftsState },
      ): void {
        if (state.syncIntervalId) {
          clearInterval(state.syncIntervalId)
          commit(types.SET_SYNC_INTERVAL, null)
        }
      },

      async syncAllLocalDrafts (
        { dispatch, state }: { dispatch: Function; state: DraftsState },
      ): Promise<void> {
        if (!state.composeAPI) return

        const localDrafts = Array.from(state.drafts.entries()).filter(([, entry]) => {
          return (
            entry.source === 'local' &&
            entry.context &&
            entry.context.recordID &&
            !entry.context.isNew
          )
        })

        for (const [changeID] of localDrafts) {
          try {
            await dispatch('syncDraft', { changeID })
          } catch (e) {
            console.warn('Failed to sync draft:', changeID, e)
          }
        }
      },

      async syncDraft (
        { commit, state }: { commit: Function; state: DraftsState },
        { changeID }: { changeID: string },
      ): Promise<void> {
        const entry = state.drafts.get(changeID)
        if (!entry || entry.source !== 'local') return

        const ctx = entry.context
        if (!ctx?.recordID || ctx.isNew) return
        if (!state.systemAPI) return
        if (state.syncing.has(changeID)) return

        commit(types.SET_SYNCING, { changeID, syncing: true })

        try {
          const revision = entry.revision

          await state.systemAPI.revisionCreate({
            resourceType: revision.resourceType,
            resourceID: ctx.recordID,
            status: 'draft',
            changes: revision.changes,
            comment: revision.comment || '',
          })

          removeDraftFromStorage(changeID)

          const updatedEntry: DraftEntry = { ...entry, source: 'backend' }
          commit(types.SET_DRAFT, { changeID, entry: updatedEntry })
        } catch (e) {
          console.error('Failed to sync draft revision:', changeID, e)
          throw e
        } finally {
          commit(types.SET_SYNCING, { changeID, syncing: false })
        }
      },
    },

    mutations: {
      [types.SET_DRAFT] (state: DraftsState, { changeID, entry }: { changeID: string; entry: DraftEntry }): void {
        state.drafts.set(changeID, entry)
      },

      [types.REMOVE_DRAFT] (state: DraftsState, changeID: string): void {
        state.drafts.delete(changeID)
      },

      [types.SET_SYNCING] (state: DraftsState, { changeID, syncing }: { changeID: string; syncing: boolean }): void {
        if (syncing) {
          state.syncing.add(changeID)
        } else {
          state.syncing.delete(changeID)
        }
      },

      [types.SET_LOADING] (state: DraftsState, loading: boolean): void {
        state.loading = loading
      },

      [types.SET_SYNC_INTERVAL] (state: DraftsState, intervalId: ReturnType<typeof setInterval> | null): void {
        state.syncIntervalId = intervalId
      },

      [types.SET_APIS] (state: DraftsState, { composeAPI, systemAPI }: { composeAPI: any; systemAPI: any }): void {
        state.composeAPI = composeAPI
        state.systemAPI = systemAPI
      },
    },
  }
}

export { generateLocalChangeID, getDraftFromStorage } from './storage'
