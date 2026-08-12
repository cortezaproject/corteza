<template>
  <Button
    data-test-id="button-export-workflow"
    :label="$t('general.export')"
    icon="pi pi-download"
    :severity="severity"
    :size="size"
    @click="jsonExport(workflows)"
  >
    <slot />
  </Button>
</template>

<script>
import { saveAs } from 'file-saver'
import { useToast } from 'primevue/usetoast'

export default {
  props: {
    workflows: {
      type: Array,
      default: () => [],
    },

    fileName: {
      type: String,
      default: 'workflows-export',
    },

    size: {
      type: String,
      default: undefined,
    },

    severity: {
      type: String,
      default: 'secondary',
    },
  },

  setup() {
    const toast = useToast()
    return { toast }
  },

  methods: {
    async jsonExport(workflowID = []) {
      const triggers = {}
      let workflows = []

      // Get workflow triggers
      await this.$AutomationAPI
        .triggerList({ workflowID, disabled: 1 })
        .then(({ set = [] }) => {
          set.forEach(
            ({ workflowID, resourceType, eventType, constraints, enabled, stepID, meta }) => {
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
            },
          )
        })
        .catch(e =>
          this.toast.add({
            severity: 'error',
            summary: this.$t('notification.failed-fetch-triggers'),
            detail: e?.message,
            life: 5000,
          }),
        )

      // Get workflows, add related triggers
      await this.$AutomationAPI
        .workflowList({ workflowID, disabled: 1, subWorkflow: 1 })
        .then(({ set = [] }) => {
          workflows = set.map(
            ({ workflowID, handle, enabled, keepSessions, steps, paths, meta }) => {
              return {
                handle,
                enabled,
                meta,
                keepSessions,
                steps,
                paths,
                triggers: triggers[workflowID],
              }
            },
          )
        })
        .catch(e =>
          this.toast.add({
            severity: 'error',
            summary: this.$t('notification.failed-fetch-workflows'),
            detail: e?.message,
            life: 5000,
          }),
        )

      // Save file
      const blob = new Blob([JSON.stringify({ workflows }, null, 2)], { type: 'application/json' })
      const filename = this.fileName.replace(/[/\\?%*:|"<>]/g, '')
      saveAs(blob, `${filename}.json`)
    },
  },
}
</script>
