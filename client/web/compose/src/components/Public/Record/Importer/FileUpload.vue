<template>
  <b-card
    header-bg-variant="white"
    footer-bg-variant="white"
  >
    <b-form-group
      :label="$t('recordList.import.uploadFile')"
      label-class="text-primary"
    >
      <uploader
        class="uploader"
        :label="dzLabel"
        :endpoint="endpoint"
        :accepted-files="['application/json', 'text/csv']"
        :max-filesize="$s('compose.Record.Attachments.MaxSize', 100)"
        @uploaded="onUploaded"
      />
    </b-form-group>

    <b-form-group
      :label="$t('recordList.import.onError')"
      label-class="text-primary"
    >
      <b-form-select
        v-model="onError"
        class="w-auto"
      >
        <option value="FAIL">
          {{ $t('recordList.import.onErrorFail') }}
        </option>

        <option value="SKIP">
          {{ $t('recordList.import.onErrorSkip') }}
        </option>
      </b-form-select>
    </b-form-group>

    <div
      slot="footer"
      class="text-right"
    >
      <b-button
        :disabled="!canContinue"
        @click="fileUploaded"
      >
        {{ $t('general.label.next') }}
      </b-button>
    </div>
  </b-card>
</template>

<script>
import Uploader from 'corteza-webapp-compose/src/components/Public/Page/Attachment/Uploader'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    Uploader,
  },

  props: {
    namespace: {
      type: Object,
      required: true,
      default: () => ({}),
    },
    module: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  data () {
    return {
      session: null,
      onError: 'FAIL',
      sessionFile: null,
    }
  },

  computed: {
    endpoint () {
      return this.$ComposeAPI.recordImportInitEndpoint({
        namespaceID: this.namespace.namespaceID,
        moduleID: this.module.moduleID,
      })
    },

    canContinue () {
      return !!this.session
    },

    dzLabel () {
      if (this.sessionFile) {
        return this.$t('recordList.import.dropzoneFileAdded', { name: this.sessionFile.name, count: this.session.progress.entryCount })
      }

      return this.$t('recordList.import.dropzoneLabel')
    },
  },

  methods: {
    onUploaded (e, f) {
      this.session = e
      this.sessionFile = f
    },

    fileUploaded () {
      this.$emit('fileUploaded', {
        ...this.session || {},
        onError: this.onError,
      })
    },
  },
}
</script>
