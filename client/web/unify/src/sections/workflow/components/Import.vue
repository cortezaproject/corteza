<template>
  <div>
    <Button
      :label="$t('general.import.label')"
      icon="pi pi-upload"
      size="small"
      severity="secondary"
      :disabled="disabled"
      @click="showDialog = true"
    />

    <Dialog
      v-model:visible="showDialog"
      :header="$t('general.import.label')"
      modal
      :closable="!processing"
      :style="{ width: '32rem' }"
      @hide="onDialogHide"
    >
      <div class="flex flex-col gap-4">
        <CFormGroup :label="$t('general.import.json')">
          <CFileDropZone
            accept=".json"
            :uploading="processing"
            :error="parseError"
            :icon="fileName ? 'pi pi-check-circle' : 'pi pi-cloud-upload'"
            :drop-label="fileName || $t('general.import.upload-files')"
            @select="onFilesSelected"
          />

          <small class="text-muted-color">
            {{ $t('general.import.reassign-run-as') }}
          </small>
        </CFormGroup>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            outlined
            size="small"
            :disabled="processing"
            @click="showDialog = false"
          />
          <Button
            :label="$t('general.import.label')"
            size="small"
            :loading="processing"
            :disabled="!workflows.length || processing"
            @click="onImport"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script>
import { useToast } from 'primevue/usetoast'
import { components } from '@planetcrust/human-vue'

const { CFileDropZone } = components

export default {
  components: {
    CFileDropZone,
  },

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
      showDialog: false,
      workflows: [],
      fileName: '',
      parseError: '',
      processing: false,
    }
  },

  methods: {
    onDialogHide () {
      if (!this.processing) {
        this.clearFile()
      }
    },

    clearFile () {
      this.workflows = []
      this.fileName = ''
      this.parseError = ''
    },

    onFilesSelected (files = []) {
      const file = files[0]
      if (!file) return

      this.processing = true
      this.parseError = ''
      const reader = new FileReader()

      reader.readAsText(file)

      reader.onload = (evt) => {
        try {
          const { workflows = [] } = JSON.parse(evt.target.result)
          if (!workflows.length) {
            this.parseError = this.$t('general.import.no-workflows-in-file')
            this.workflows = []
            this.fileName = ''
            return
          }
          this.workflows = workflows
          this.fileName = file.name
        } catch (err) {
          this.parseError = err?.message || ''
          this.workflows = []
          this.fileName = ''
          this.toast.add({ severity: 'error', summary: this.$t('notification.general.warning'), detail: err?.message, life: 5000 })
        } finally {
          this.processing = false
        }
      }

      reader.onerror = () => {
        this.parseError = this.$t('notification.failed-load-file')
        this.toast.add({ severity: 'error', summary: this.$t('notification.failed-load-file'), life: 5000 })
        this.processing = false
      }
    },

    onImport () {
      this.$emit('import', this.workflows)
      this.showDialog = false
    },
  },
}
</script>
