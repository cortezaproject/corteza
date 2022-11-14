/* eslint-disable @typescript-eslint/ban-ts-ignore */

import alert from './CPromptAlert.vue'
import choice from './CPromptChoice.vue'
import composeRecordPicker from './CPromptComposeRecordPicker.vue'
import input from './CPromptInput.vue'
import notification from './CPromptNotification.vue'
import options from './CPromptOptions.vue'
import { Component } from 'vue'
import { pType, pVal } from '../utils'
import { automation } from '@cortezaproject/corteza-js'
import { KV } from '@cortezaproject/corteza-js/dist/compose/types/chart/util'

interface Handler {
  (this: Component, input: automation.Vars): void|Promise<void>;
}

interface PromptDefinition {
  component?: Component;

  /**
   *
   */
  handler?: Handler;

  /**
   * Passive prompt, will not be listed
   *
   * Also, when displaying toasts, we'll display all
   * passive toasts first and then, at the and one single
   * non-passive toast
   */
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
    passive: true,
    component: notification,
  },

  options: {
    component: options,
  },

  redirect: {
    handler: function (v): void {
      const url = pVal(v, 'url')
      const delay = (pVal(v, 'delay') || 0) as number
      if (url !== undefined) {
        console.debug('redirect to %s via prompt in %d sec', url, delay)
        setTimeout(() => {
          // @ts-ignore
          window.location = url
        }, delay * 1000)
      }
    },
  },

  reroute: {
    handler: function (v): void {
      const name = pVal(v, 'name')
      const params = pVal(v, 'params')
      const query = pVal(v, 'query')
      const delay = (pVal(v, 'delay') || 0) as number
      if (name !== undefined) {
        console.debug('reroute to %s via prompt in %d sec', name, delay, { params, query })
        setTimeout(() => {
          // @ts-ignore
          this.$router.push({ name, params, query })
        }, delay * 1000)
      }
    },
  },

  recordPage: {
    handler: async function (v): Promise<void> {
      const module = pVal(v, 'module')
      const namespace = pVal(v, 'namespace')
      const record = pVal(v, 'record')
      const edit = !!pVal(v, 'edit')
      const delay = (pVal(v, 'delay') || 0) as number

      let namespaceID = ''
      let slug = ''
      let moduleID = ''
      let recordID = ''
      let pageID = ''

      // We can extract almost anything for a record
      if (pType(v, 'record') === 'ComposeRecord') {
        namespaceID = (record as KV).namespaceID
        moduleID = (record as KV).moduleID
        recordID = (record as KV).recordID
      } else {
        // Resolve recordID
        if (pType(v, 'record') === 'ID') {
          recordID = record as string
        }

        // Resolve namespaceID
        if (pType(v, 'namespace') === 'ID') {
          namespaceID = namespace as string
        } else if (pType(v, 'namespace') === 'ComposeNamespace') {
          namespaceID = (namespace as KV).namespaceID as string
          slug = (namespace as KV).slug as string
        } else {
          // @ts-ignore
          const { set: nn } = await this.$ComposeAPI.namespaceList({ slug: namespace as string })
          if (!nn || nn.length !== 1) {
            throw new Error('namespace not resolved')
          }

          namespaceID = nn[0].namespaceID
          slug = nn[0].slug
        }

        // Resolve moduleID
        if (pType(v, 'module') === 'ID') {
          moduleID = module as string
        } else if (pType(v, 'module') === 'ComposeModule') {
          moduleID = (module as KV).moduleID as string
        } else {
          // @ts-ignore
          const { set: mm } = await this.$ComposeAPI.moduleList({ handle: module as string, namespaceID })
          if (!mm || mm.length !== 1) {
            throw new Error('module not resolved')
          }

          moduleID = mm[0].moduleID
        }
      }

      if (!slug) {
        // @ts-ignore
        const ns = await this.$ComposeAPI.namespaceRead({ namespaceID })
        slug = ns.slug
      }

      // Resolve record page
      // @ts-ignore
      const { set: pp } = await this.$ComposeAPI.pageList({ moduleID, namespaceID })
      if (!pp || pp.length !== 1) {
        throw new Error('record page not resolved')
      }
      pageID = pp[0].pageID

      // @ts-ignore
      if (this.$root.$options.name === 'compose') {
        if (!edit && !recordID) {
          throw new Error('invalid record page prompt configuration')
        }

        let name = 'page.record'
        if (edit) {
          name += recordID ? '.edit' : '.create'
        }

        // If name and params match, make sure to refresh page instead of push
        // @ts-ignore
        const reloadPage = name === this.$route.name && slug === this.$route.params.slug && pageID === this.$route.params.pageID && recordID === this.$route.params.recordID

        setTimeout(() => {
          console.debug('reroute to %s via prompt in %d sec', name, delay, { namespaceID, slug, moduleID, recordID })
          if (reloadPage) {
            window.location.reload()
          } else {
            // @ts-ignore
            this.$router.push({ name, params: { recordID, pageID, slug } })
          }
        }, delay * 1000)
      } else {
        // @ts-ignore
        const u = new URL(window.location)
        // Generate direct link
        setTimeout(() => {
          const url = `${u.origin}/compose/ns/crm/pages/${pageID}/record/${recordID}/${edit ? 'edit' : ''}`
          console.debug('redirect to %s via prompt in %d sec', url, delay)
          // @ts-ignore
          window.location = url
        }, delay * 1000)
      }
    },
  },

  refetchRecords: {
    handler: function (): void {
      // @ts-ignore
      this.$root.$emit('refetch-records')
    },
  },
}

export default definitions
