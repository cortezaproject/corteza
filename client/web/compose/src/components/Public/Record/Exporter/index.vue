<template>
  <div>
    <b-button
      size="lg"
      variant="light"
      @click="toggleModal"
    >
      {{ $t('label.export') }}
    </b-button>
    <b-modal
      :visible="showExportModal"
      size="lg"
      :title="$t('label.export')"
      hide-footer
      body-class="p-0"
      no-fade
      @hide="toggleModal"
    >
      <configurator
        :processing="processing"
        v-bind="$attrs"
        v-on="$listeners"
      />
    </b-modal>
  </div>
</template>

<script>
import Configurator from './Configurator'

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  components: {
    Configurator,
  },

  inheritAttrs: true,

  props: {
    processing: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      showExportModal: false,
    }
  },

  methods: {
    toggleModal () {
      this.showExportModal = !this.showExportModal
    },
  },
}
</script>
