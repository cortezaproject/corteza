<template>
  <div class="d-flex">
    <b-button
      data-test-id="button-import"
      variant="light"
      size="lg"
      class="flex-fill"
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
          v-bind="$props"
          @imported="onImported"
          @reset="onReset"
          @close="onClose"
        />
      </keep-alive>
    </b-modal>
  </div>
</template>

<script>
import FileUpload from './FileUpload'

export default {
  i18nOptions: {
    namespaces: 'system.users',
  },

  name: 'CUserImport',

  data () {
    return {
      step: 0,
      showModal: false,
      components: [FileUpload],
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
      this.showModal = false
    },

    onImported (e) {
      this.$emit('imported')
      this.onReset()
      this.onClose()
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
