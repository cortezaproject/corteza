import { onBeforeUnmount, onMounted, ref, toValue, type MaybeRefOrGetter } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { useI18n } from 'vue-i18n'

interface UnsavedGuardOptions {
  isDirty: MaybeRefOrGetter<boolean>
  messageKey: string
  tabClose?: boolean
}

export function useUnsavedGuard(options: UnsavedGuardOptions) {
  const { isDirty, messageKey, tabClose = true } = options
  const { t } = useI18n()
  const navigatingAfterSave = ref(false)

  function markSaved() {
    navigatingAfterSave.value = true
  }

  const beforeUnloadHandler = (e: BeforeUnloadEvent) => {
    if (!toValue(isDirty)) return
    e.preventDefault()
    e.returnValue = ''
  }

  onMounted(() => {
    if (tabClose) window.addEventListener('beforeunload', beforeUnloadHandler)
  })

  onBeforeUnmount(() => {
    if (tabClose) window.removeEventListener('beforeunload', beforeUnloadHandler)
  })

  onBeforeRouteLeave((to, from, next) => {
    if (navigatingAfterSave.value) {
      navigatingAfterSave.value = false
      next()
      return
    }
    if (!toValue(isDirty)) {
      next()
      return
    }
    next(window.confirm(t(messageKey)))
  })

  return { markSaved }
}
