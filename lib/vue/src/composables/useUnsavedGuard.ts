import { onBeforeUnmount, onMounted, toValue, type MaybeRefOrGetter } from 'vue'
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
    if (!toValue(isDirty)) {
      next()
      return
    }
    next(window.confirm(t(messageKey)))
  })
}
