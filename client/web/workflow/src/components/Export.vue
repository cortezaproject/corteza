<template>
  <b-button
    :variant="variant"
    :size="size"
    @click="jsonExport(workflows)"
  >
    {{ $t('general:export') }}
  </b-button>
</template>

<script>
import { saveAs } from 'file-saver'

export default {
  props: {
    workflows: {
      type: Array,
      default: () => ([]),
    },

    fileName: {
      type: String,
      default: 'workflows-export',
    },

    size: {
      type: String,
      default: 'md',
    },

    variant: {
      type: String,
      default: 'light',
    },
  },

  methods: {
    async jsonExport (workflowID = []) {
      const triggers = {}
      let workflows = []

      // Get workflow triggers
      await this.$AutomationAPI.triggerList({ workflowID, disabled: 1 })
        .then(({ set = [] }) => {
          set.forEach(({ workflowID, resourceType, eventType, constraints, enabled, stepID, meta }) => {
            if (!triggers[workflowID]) {
              triggers[workflowID] = []
            }

            triggers[workflowID].push({
              resourceType,
              eventType,
              constraints,
              enabled,
              stepID,
              meta,
            })
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:failed-fetch-triggers')))

      // Get workflows, add related triggers
      await this.$AutomationAPI.workflowList({ workflowID, disabled: 1, subWorkflow: 1 })
        .then(({ set = [] }) => {
          workflows = set.map(({ workflowID, handle, enabled, keepSessions, steps, paths, meta }) => {
            return {
              handle,
              enabled,
              meta,
              keepSessions,
              steps,
              paths,
              triggers: triggers[workflowID],
            }
          })
        })
        .catch(this.toastErrorHandler(this.$t('notification:failed-fetch-workflows')))

      // Save file
      const blob = new Blob([JSON.stringify({ workflows }, null, 2)], { type: 'application/json' })
      const filename = this.fileName.replace(/[/\\?%*:|"<>]/g, '')
      saveAs(blob, `${filename}.json`)
    },
  },
}
</script>
