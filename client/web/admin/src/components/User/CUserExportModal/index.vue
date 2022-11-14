<template>
  <div>
    <b-button
      data-test-id="button-export"
      variant="light"
      @click="showModal=true"
    >
      {{ $t('export.buttonLabel') }}
    </b-button>

    <b-modal
      :visible="showModal"
      size="lg"
      :title="$t('export.title')"
      hide-footer
      body-class="p-0"
      @hide="onModalHide"
    >
      <keep-alive>
        <component
          :is="stepComponent"
          v-if="!processing"
          v-bind="$props"
          :session="session"
          @configured="onConfigured"
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
import ExportConfiguration from './ExportConfiguration.vue'

export default {
  i18nOptions: {
    namespaces: 'system.users',
  },

  name: 'CUserExport',

  data () {
    return {
      step: 0,
      showModal: false,
      session: {},
      components: [ExportConfiguration],

      processing: false,
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

    async onConfigured (e) {
      this.processing = true
      this.$emit('export', e)
      this.onReset()
      this.onClose()
      this.processing = false
    },

    onReset () {
      this.step = 0
      this.$emit('reset')
    },

    onClose () {
      this.showModal = false
    },
  },
}
</script>
