<template>
  <div class="d-flex">
    <b-button
      variant="light"
      size="lg"
      class="flex-fill"
      @click="show = true"
    >
      {{ $t('general:import.label') }}
    </b-button>

    <b-modal
      id="import"
      v-model="show"
      size="lg"
      :title="$t('general:import.json')"
      ok-only
      no-fade
      class="d-none"
      @ok="$emit('import', workflows)"
    >
      <b-form-group
        :description="$t('general:import.reassign-run-as')"
        class="mb-0"
      >
        <b-form-file
          :placeholder="$t('general:import.upload-files')"
          @change="fileUpload"
        />
      </b-form-group>

      <template #modal-footer>
        <b-button
          variant="primary"
          size="lg"
          :disabled="!workflows.length || processing"
          class="d-flex justify-content-center align-items-center"
          @click="$emit('import', workflows)"
        >
          <b-spinner
            v-if="processing"
            small
            type="grow"
          />

          <span
            v-else
          >
            {{ $t('general:import.label') }}
          </span>
        </b-button>
      </template>
    </b-modal>
  </div>
</template>

<script>
export default {
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      show: false,
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
            err.message = this.$t('notification:failed-load-file')
            this.toastErrorHandler(this.$t('notification:general.warning'))(err)
          } finally {
            this.processing = false
          }
        }

        reader.onerror = () => {
          this.toastErrorHandler(this.$t('notification:failed-load-file'))
          this.processing = false
        }
      }
    },
  },
}
</script>

<style>

</style>
