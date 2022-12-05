<template>
  <div>
    <b-button
      data-test-id="button-import"
      size="lg"
      variant="light"
      @click="showModal=true"
    >
      {{ $t('import.buttonLabel') }}
    </b-button>

    <b-modal
      :visible="showModal"
      size="lg"
      :title="$t('import.title')"
      hide-footer
      body-class="p-0"
      @hide="onModalHide"
    >
      <keep-alive>
        <component
          :is="stepComponent"
          v-if="!importing"
          v-bind="$props"
          :session="session"
          @fileUploaded="onFileUploaded"
          @configured="onConfigured"
          @back="onBack"
          @reset="onReset"
          @close="onClose"
          v-on="$listeners"
        />
        <div
          v-else
          class="p-5 h-100 d-flex align-items-center justify-content-center"
        >
          <b-spinner />
        </div>
      </keep-alive>
    </b-modal>
  </div>
</template>

<script>
import FileUpload from './FileUpload'
import ImportConfiguration from './ImportConfiguration.vue'

export default {
  i18nOptions: {
    namespaces: 'namespace',
  },

  name: 'Importer',

  data () {
    return {
      step: 0,
      showModal: false,
      session: {},
      components: [FileUpload, ImportConfiguration],

      importing: false,
    }
  },

  computed: {
    stepComponent () {
      return this.components[this.step]
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

    async onConfigured (e) {
      this.importing = true

      try {
        const out = await this.$ComposeAPI.namespaceImportRun({
          sessionID: this.session.sessionID,
          name: e.name,
          slug: e.slug,
        })

        this.$emit('imported', out)
      } catch (err) {
        this.$emit('failed', err)
      }

      this.onReset()
      this.onClose()
      this.importing = false
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
