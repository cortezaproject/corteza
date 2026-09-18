import { inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

// Opens the record a `record` notification points at, in the notification's
// openMode; any other kind has nothing to open.
export function useOpenNotification() {
  const route = useRoute()
  const router = useRouter()
  const $ComposeAPI = inject<any>('$ComposeAPI', null)
  const $toast = inject<any>('$toast', null)
  const { t } = useI18n()

  return async function openNotification(notification: { kind?: string; config?: any }) {
    const { config = {}, kind } = notification
    if (kind !== 'record' || !$ComposeAPI) {
      return
    }

    const { namespaceID, moduleID, recordID, openMode, edit } = config

    try {
      const namespace = await $ComposeAPI.namespaceRead({ namespaceID })
      if (!namespace) {
        $toast?.toastDanger?.(t('notifications.namespaceNotFound'))
        return
      }

      const { set: recordPages = [] } = await $ComposeAPI.pageList({ namespaceID, moduleID })
      if (!recordPages.length) {
        $toast?.toastDanger?.(t('notifications.pageNotFound'))
        return
      }

      if (recordID && recordID !== '0') {
        try {
          const record = await $ComposeAPI.recordRead({ namespaceID, moduleID, recordID })
          if (!record) {
            $toast?.toastDanger?.(t('notifications.recordNotFound'))
            return
          }
        } catch {
          $toast?.toastDanger?.(t('notifications.recordNotFound'))
          return
        }
      }

      const slug = namespace.slug || namespace.namespaceID
      const pageID = recordPages[0].pageID
      const recID = !recordID || recordID === '0' ? '0' : recordID

      const hasComposeRoute = router.hasRoute('page.record')
      const u = new URL(window.location.href)
      let externalUrl = `${u.origin}/compose/namespace/${slug}/pages/${pageID}/records/${recID}`
      if (edit) {
        externalUrl += '?edit=1'
      }

      const routeLocation = {
        name: 'page.record',
        params: { slug, pageID, recordID: recID },
        query: edit ? { edit: '1' } : {},
      }

      // Modal
      if (openMode === 'modal') {
        if (hasComposeRoute) {
          const modalQuery = {
            recordPageID: pageID,
            recordID: recID,
            ...(edit ? { edit: '1' } : {}),
          }
          // If targeting a different namespace, navigate there first so RecordModal
          // opens in the correct namespace context with its page store populated.
          if (route.params.slug !== slug) {
            return router.push({ name: 'pages', params: { slug }, query: modalQuery })
          }
          return router.push({
            query: { ...route.query, ...modalQuery },
          })
        }

        // Outside Compose — redirect to namespace root with modal query params
        const modalUrl = new URL(`${u.origin}/compose/namespace/${slug}`)
        modalUrl.searchParams.set('recordPageID', pageID)
        modalUrl.searchParams.set('recordID', recID)
        if (edit) modalUrl.searchParams.set('edit', '1')
        window.location.href = modalUrl.toString()
        return
      }

      // New tab
      if (openMode === 'newTab') {
        return hasComposeRoute
          ? window.open(router.resolve(routeLocation).href, '_blank', 'noopener')
          : window.open(externalUrl, '_blank', 'noopener')
      }

      // Same tab (default) or modal fallback outside Compose
      if (hasComposeRoute) {
        return router.push(routeLocation)
      }
      window.location.href = externalUrl
    } catch {
      $toast?.toastDanger?.(t('notifications.recordRedirectError'))
    }
  }
}
