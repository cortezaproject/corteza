<template>
  <div>
    <b-button
      size="lg"
      variant="light"
      @click="showModal=true"
    >
      {{ $t('general.label.import') }}
    </b-button>

    <b-modal
      :visible="showModal"
      size="lg"
      :title="$t('recordList.import.to', { modulename: module.name })"
      hide-footer
      body-class="p-0"
      @hide="onModalHide"
    >
      <keep-alive>
        <component
          :is="stepComponent"
          v-bind="$props"
          :session="session"
          @fileUploaded="onFileUploaded"
          @fieldsMatched="onFieldsMatched"
          @importFailed="onImportFailed"
          @back="onBack"
          @reset="onReset"
          @close="onClose"
          v-on="$listeners"
        >
          <label
            v-if="progress.failed"
            slot="uploadLabel"
            class="text-danger"
          >

            {{ $t('recordList.import.failed', progress) }}
          </label>
        </component>
      </keep-alive>
    </b-modal>
  </div>
</template>

<script>
import FileUpload from './FileUpload'
import FieldMatch from './FieldMatch'
import Progress from './Progress'
import ErrorReport from './ErrorReport'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Importer',
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
      step: 0,
      showModal: false,
      session: {},
      components: [FileUpload, FieldMatch, Progress, ErrorReport],
    }
  },

  computed: {
    stepComponent () {
      return this.components[this.step]
    },

    progress () {
      return this.session.progress || {}
    },
  },

  methods: {
    onModalHide () {
      this.step = 0
      this.session = {}
      this.showModal = false
    },

    onBack () {
      this.step = Math.max(0, this.step - 1)
    },

    onFileUploaded (e) {
      this.session = e
      this.step = 1
    },

    onFieldsMatched (e) {
      this.session.fields = e
      this.step = 2

      this.$ComposeAPI.recordImportRun(this.session)
    },

    onImportFailed (e) {
      this.session.progress = e
      this.step = 3
    },

    onReset () {
      this.step = 0
      this.$set(this, 'session', {})
      this.$emit('reset')
    },

    onClose () {
      this.showModal = false
    },
  },
}
</script>
