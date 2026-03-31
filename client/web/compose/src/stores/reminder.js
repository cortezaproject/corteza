import { system, NoID } from '@cortezaproject/corteza-js-next'
import { defineStore } from 'pinia'
import { computed, inject, reactive, toRef } from 'vue'

function normalizeReminder (raw = {}) {
  return raw instanceof system.Reminder ? raw : new system.Reminder(raw)
}

function sortReminders (items = []) {
  return [...items].sort((a, b) => {
    if (!!a.dismissedAt !== !!b.dismissedAt) {
      return a.dismissedAt ? 1 : -1
    }

    if (!a.dismissedAt && !b.dismissedAt) {
      const aTime = a.remindAt ? new Date(a.remindAt).getTime() : Number.MAX_SAFE_INTEGER
      const bTime = b.remindAt ? new Date(b.remindAt).getTime() : Number.MAX_SAFE_INTEGER
      return aTime - bTime
    }

    const aDismissed = a.dismissedAt ? new Date(a.dismissedAt).getTime() : 0
    const bDismissed = b.dismissedAt ? new Date(b.dismissedAt).getTime() : 0
    return bDismissed - aDismissed
  })
}

function reminderVersion (reminder) {
  return [
    reminder.remindAt ? new Date(reminder.remindAt).toISOString() : '',
    reminder.dismissedAt ? new Date(reminder.dismissedAt).toISOString() : '',
    reminder.snoozeCount || 0,
    reminder.payload?.title || '',
    reminder.payload?.notes || '',
    reminder.payload?.link?.label || '',
  ].join('|')
}

export const useReminderStore = defineStore('compose-reminder', () => {
  const $SystemAPI = inject('$SystemAPI')
  const $Auth = inject('$Auth', inject('$auth', {}))

  const state = reactive({
    reminders: [],
    toasts: [],
    visible: false,
    editing: null,
    processing: false,
    shownVersions: new Map(),
    timer: null,
  })

  const currentUserID = computed(() => $Auth?.user?.userID || NoID)
  const activeCount = computed(() => state.reminders.filter(({ dismissedAt }) => !dismissedAt).length)

  function clearTimer () {
    if (state.timer) {
      window.clearTimeout(state.timer)
      state.timer = null
    }
  }

  function findReminder (reminderID) {
    return state.reminders.find(({ reminderID: id }) => id === reminderID)
  }

  function setReminders (reminders = []) {
    state.reminders = sortReminders(reminders.map(normalizeReminder))

    state.toasts = state.toasts
      .map(({ reminderID }) => findReminder(reminderID))
      .filter(reminder => reminder && !reminder.dismissedAt)

    scheduleDueProcessing()
  }

  function upsertReminder (raw) {
    const reminder = normalizeReminder(raw)
    const next = [...state.reminders]
    const index = next.findIndex(({ reminderID }) => reminderID === reminder.reminderID)

    if (index > -1) {
      next.splice(index, 1, reminder)
    } else {
      next.push(reminder)
    }

    setReminders(next)
    return reminder
  }

  function removeReminder (reminderID) {
    state.reminders = state.reminders.filter(({ reminderID: id }) => id !== reminderID)
    state.toasts = state.toasts.filter(({ reminderID: id }) => id !== reminderID)
    state.shownVersions.delete(reminderID)

    if (state.editing?.reminderID === reminderID) {
      state.editing = null
    }

    scheduleDueProcessing()
  }

  function showToast (raw) {
    const reminder = normalizeReminder(raw)
    if (reminder.dismissedAt) {
      hideToast(reminder.reminderID)
      state.shownVersions.delete(reminder.reminderID)
      return
    }

    const version = reminderVersion(reminder)
    const index = state.toasts.findIndex(({ reminderID }) => reminderID === reminder.reminderID)

    state.shownVersions.set(reminder.reminderID, version)

    if (index > -1) {
      state.toasts.splice(index, 1, reminder)
    } else {
      state.toasts.unshift(reminder)
    }
  }

  function hideToast (reminderID) {
    state.toasts = state.toasts.filter(({ reminderID: id }) => id !== reminderID)
  }

  function processDueReminders (now = new Date()) {
    state.reminders.forEach(reminder => {
      if (reminder.dismissedAt || !reminder.remindAt) {
        return
      }

      if (new Date(reminder.remindAt) <= now) {
        const version = reminderVersion(reminder)
        if (state.shownVersions.get(reminder.reminderID) !== version) {
          showToast(reminder)
        }
      }
    })

    scheduleDueProcessing()
  }

  function scheduleDueProcessing () {
    clearTimer()

    const now = new Date()
    let nextDueAt = null

    state.reminders.forEach(reminder => {
      if (reminder.dismissedAt || !reminder.remindAt) {
        return
      }

      const remindAt = new Date(reminder.remindAt)
      if (remindAt <= now) {
        nextDueAt = now
        return
      }

      if (!nextDueAt || remindAt < nextDueAt) {
        nextDueAt = remindAt
      }
    })

    if (!nextDueAt) {
      return
    }

    const delay = Math.max(0, nextDueAt.getTime() - now.getTime())
    state.timer = window.setTimeout(() => processDueReminders(new Date()), delay)
  }

  async function fetchReminders () {
    if (!$SystemAPI || !currentUserID.value || currentUserID.value === NoID) {
      setReminders([])
      return []
    }

    const { set = [] } = await $SystemAPI.reminderList({
      assignedTo: currentUserID.value,
      limit: 0,
    })

    const reminders = set.map(reminder => normalizeReminder(reminder))
    setReminders(reminders)
    processDueReminders()

    return reminders
  }

  function setVisible (visible) {
    state.visible = visible
  }

  function toggleVisibility () {
    state.visible = !state.visible
  }

  function clearEdit () {
    state.editing = null
  }

  function startCreate ({ resource, assignedTo, payload = {}, remindAt } = {}) {
    state.editing = new system.Reminder({
      resource,
      assignedTo: assignedTo || currentUserID.value || NoID,
      payload,
      remindAt,
    })
    state.visible = true
  }

  function startEdit (reminder) {
    state.editing = normalizeReminder(reminder)
    state.visible = true
  }

  async function saveReminder (reminder) {
    if (!$SystemAPI) return

    state.processing = true

    try {
      const endpoint = reminder.reminderID && reminder.reminderID !== NoID ? 'reminderUpdate' : 'reminderCreate'
      await $SystemAPI[endpoint]({
        reminderID: reminder.reminderID,
        resource: reminder.resource,
        assignedTo: reminder.assignedTo,
        payload: reminder.payload,
        remindAt: reminder.remindAt ? new Date(reminder.remindAt).toISOString() : undefined,
      })

      await fetchReminders()
      state.editing = null
      state.visible = true
    } finally {
      state.processing = false
    }
  }

  async function setDismissed (reminder, value) {
    if (!$SystemAPI) return

    const endpoint = value ? 'reminderDismiss' : 'reminderUndismiss'
    await $SystemAPI[endpoint]({ reminderID: reminder.reminderID })

    hideToast(reminder.reminderID)
    if (!value) {
      state.shownVersions.delete(reminder.reminderID)
    }

    await fetchReminders()
  }

  async function deleteReminder (reminder) {
    if (!$SystemAPI) return

    await $SystemAPI.reminderDelete({ reminderID: reminder.reminderID })
    removeReminder(reminder.reminderID)
  }

  async function snoozeReminder (reminder, duration) {
    if (!$SystemAPI) return

    const remindAt = new Date(Date.now() + duration).toISOString()

    await $SystemAPI.reminderSnooze({
      reminderID: reminder.reminderID,
      remindAt,
    })

    hideToast(reminder.reminderID)
    await fetchReminders()
  }

  function handleRealtimeReminder (raw) {
    const reminder = upsertReminder(raw)

    if (reminder.dismissedAt) {
      hideToast(reminder.reminderID)
      state.shownVersions.delete(reminder.reminderID)
      return
    }

    if (reminder.remindAt && new Date(reminder.remindAt) <= new Date()) {
      showToast(reminder)
    }
  }

  function dispose () {
    clearTimer()
  }

  return {
    reminders: toRef(state, 'reminders'),
    toasts: toRef(state, 'toasts'),
    visible: toRef(state, 'visible'),
    editing: toRef(state, 'editing'),
    processing: toRef(state, 'processing'),
    currentUserID,
    activeCount,
    fetchReminders,
    setVisible,
    toggleVisibility,
    clearEdit,
    startCreate,
    startEdit,
    saveReminder,
    setDismissed,
    deleteReminder,
    snoozeReminder,
    handleRealtimeReminder,
    hideToast,
    dispose,
  }
})
