import { markRaw, type Component } from 'vue'
import { resolveInternalPath } from '../../../utils/internalNav'
import { pType, pVal } from '../utils'
import alertCmp from './CPromptAlert.vue'
import choiceCmp from './CPromptChoice.vue'
import composeRecordPickerCmp from './CPromptComposeRecordPicker.vue'
import inputCmp from './CPromptInput.vue'
import notificationCmp from './CPromptNotification.vue'
import optionsCmp from './CPromptOptions.vue'

// markRaw keeps Vue's reactivity system from deep-walking these component
// definitions when they're stored inside a reactive ref (e.g. passivePrompts
// in CPromptToast.vue). Without it, Vue logs "Component that was made a
// reactive object" and pays the cost of proxying every option on the SFC.
const alert = markRaw(alertCmp)
const choice = markRaw(choiceCmp)
const composeRecordPicker = markRaw(composeRecordPickerCmp)
const input = markRaw(inputCmp)
const notification = markRaw(notificationCmp)
const options = markRaw(optionsCmp)

interface PromptDefinition {
  component?: Component;
  handler?: (_this: any, _input: any) => void | Promise<void>;
  passive?: boolean;
}

const definitions: Record<string, PromptDefinition> = {
  alert: {
    component: alert,
  },
  choice: {
    component: choice,
  },
  composeRecordPicker: {
    component: composeRecordPicker,
  },
  input: {
    component: input,
  },
  notification: {
    component: notification,
    passive: true,
  },
  options: {
    component: options,
  },
  redirect: {
    handler(v) {
      const url = pVal(v, 'url')
      const delay = Number(pVal(v, 'delay', 0) || 0)
      const openMode = pVal(v, 'openMode', 'sameTab')

      if (!url) {
        return
      }

      window.setTimeout(() => {
        if (openMode === 'newTab') {
          window.open(String(url), '_blank', 'noopener')
          return
        }
        // Navigate client-side when the URL is an internal route of this app;
        // otherwise do a full-page redirect (other app / external URL).
        const internal = resolveInternalPath(this?.$router, String(url))
        if (internal && this?.$router) {
          this.$router.push(internal)
        } else {
          window.location.assign(String(url))
        }
      }, delay * 1000)
    },
  },
  reroute: {
    handler(v) {
      const name = pVal(v, 'name')
      const params = pVal(v, 'params')
      const query = pVal(v, 'query')
      const delay = Number(pVal(v, 'delay', 0) || 0)
      const openMode = pVal(v, 'openMode', 'sameTab')

      if (!name) {
        return
      }

      window.setTimeout(() => {
        const target = { name, params, query }
        if (openMode === 'newTab') {
          const resolved = this.$router.resolve(target)
          window.open(resolved.href, '_blank', 'noopener')
        } else {
          this.$router.push(target)
        }
      }, delay * 1000)
    },
  },
  recordPage: {
    async handler(v) {
      const namespace = pVal(v, 'namespace')
      const record = pVal(v, 'record')
      const edit = !!pVal(v, 'edit')
      const delay = Number(pVal(v, 'delay', 0) || 0)
      const openMode = pVal(v, 'openMode', 'sameTab')

      let namespaceID = ''
      let slug = ''
      let recordID = ''
      let moduleID = ''

      if (pType(v, 'record') === 'ComposeRecord') {
        namespaceID = record.namespaceID
        moduleID = record.moduleID
        recordID = record.recordID
      } else {
        recordID = String(record || '')
        if (pType(v, 'namespace') === 'ComposeNamespace') {
          namespaceID = namespace.namespaceID
          slug = namespace.slug || namespace.namespaceID
        } else {
          namespaceID = String(namespace || '')
        }
      }

      if (!slug && namespaceID) {
        const ns = await this.$ComposeAPI.namespaceRead({ namespaceID })
        slug = ns.slug || ns.namespaceID
      }

      const { set: pages = [] } = await this.$ComposeAPI.pageList({ namespaceID, moduleID })
      if (!pages.length) {
        this.$toast?.toastDanger?.('Record page not resolved', 'Prompt error')
        return
      }

      const target = {
        name: 'page.record',
        params: { slug, pageID: pages[0].pageID, recordID },
        query: edit ? { edit: '1' } : {},
      }

      window.setTimeout(() => {
        // Human does not yet have record modal parity; modal falls back to same-tab navigation.
        if (openMode === 'newTab') {
          const resolved = this.$router.resolve(target)
          window.open(resolved.href, '_blank', 'noopener')
        } else {
          this.$router.push(target)
        }
      }, delay * 1000)
    },
  },
  refetchRecords: {
    handler() {
      this.$eventBus?.emit?.('refetch-records')
    },
  },
}

export default definitions
