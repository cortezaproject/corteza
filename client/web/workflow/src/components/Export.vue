<template>
  <b-button
    v-if="workflows.length > 0"
    variant="light"
    size="lg"
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
      required: true,
    },

    fileName: {
      type: String,
      default: 'workflows-export',
    },
  },

  methods: {
    async jsonExport (workflowIDs) {
      const triggers = {}
      let workflows = []

      // Get workflow triggers
      await this.$AutomationAPI.triggerList({ workflowID: workflowIDs, disabled: 1 })
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
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-triggers')))

      // Get workflows, add related triggers
      await this.$AutomationAPI.workflowList({ workflowID: workflowIDs, disabled: 1 })
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
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-workflows')))

      // Save file
      const blob = new Blob([JSON.stringify({ workflows }, null, 2)], { type: 'application/json' })
      saveAs(blob, `${this.fileName}.json`)
    },
  },
}
</script>
