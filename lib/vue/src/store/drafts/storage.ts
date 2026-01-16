import { system } from '@cortezaproject/corteza-js'

const { Revision } = system

export const DRAFT_STORAGE_PREFIX = 'corteza:revision:draft:'

export function buildStorageKey (changeID: string): string {
  return `${DRAFT_STORAGE_PREFIX}${changeID}`
}

export function isDraftStorageKey (key: string): boolean {
  return key.startsWith(DRAFT_STORAGE_PREFIX)
}

export function extractChangeID (key: string): string | null {
  if (!isDraftStorageKey(key)) return null
  return key.slice(DRAFT_STORAGE_PREFIX.length)
}

export function getAllDraftKeys (): string[] {
  const keys: string[] = []
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i)
    if (key && isDraftStorageKey(key)) {
      keys.push(key)
    }
  }
  return keys
}

export function getDraftFromStorage (changeID: string): system.Revision | null {
  try {
    const key = buildStorageKey(changeID)
    const item = localStorage.getItem(key)
    if (!item) return null
    return new Revision(JSON.parse(item))
  } catch (e) {
    console.warn('Failed to parse draft:', changeID, e)
    return null
  }
}

export function saveDraftToStorage (changeID: string, revision: system.Revision): void {
  try {
    const key = buildStorageKey(changeID)
    localStorage.setItem(key, JSON.stringify(revision))
  } catch (e) {
    console.error('Failed to save draft:', changeID, e)
  }
}

export function removeDraftFromStorage (changeID: string): void {
  localStorage.removeItem(buildStorageKey(changeID))
}

export function clearAllDraftsFromStorage (): void {
  for (const storageKey of getAllDraftKeys()) {
    localStorage.removeItem(storageKey)
  }
}

export function loadAllDraftsFromStorage (): Map<string, system.Revision> {
  const drafts = new Map<string, system.Revision>()

  for (const storageKey of getAllDraftKeys()) {
    const changeID = extractChangeID(storageKey)
    if (!changeID) continue

    try {
      const item = localStorage.getItem(storageKey)
      if (item) {
        drafts.set(changeID, new Revision(JSON.parse(item)))
      }
    } catch (e) {
      console.warn('Failed to parse draft:', storageKey, e)
    }
  }

  return drafts
}
