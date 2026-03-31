<template>
  <div>
    <div class="flex flex-col gap-3">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('general.import.json') }}
        </label>
        <input
          ref="fileInput"
          type="file"
          accept=".json"
          class="block w-full text-sm text-color border border-surface rounded-border p-2 cursor-pointer"
          @change="fileUpload"
        />
        <small class="text-muted-color">
          {{ $t('general.import.reassign-run-as') }}
        </small>
      </div>

      <div class="flex justify-end">
        <Button
          :label="$t('general.import.label')"
          :loading="processing"
          :disabled="!workflows.length || processing"
          @click="$emit('import', workflows)"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { useToast } from 'primevue/usetoast'

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  setup () {
    const toast = useToast()
    return { toast }
  },

  data () {
    return {
      workflows: [],
      processing: false,
    }
  },

  methods: {
    fileUpload (e = {}) {
      const { files = [] } = (e.type === 'drop' ? e.dataTransfer : e.target) || {}

      if (files[0]) {
        this.processing = true
        const reader = new FileReader()

        reader.readAsText(files[0])

        reader.onload = (evt) => {
          try {
            const { workflows = [] } = JSON.parse(evt.target.result)
            this.workflows = workflows
          } catch (err) {
            this.toast.add({ severity: 'error', summary: this.$t('notification.general.warning'), detail: err?.message, life: 5000 })
          } finally {
            this.processing = false
          }
        }

        reader.onerror = () => {
          this.toast.add({ severity: 'error', summary: this.$t('notification.failed-load-file'), life: 5000 })
          this.processing = false
        }
      }
    },
  },
}
</script>
