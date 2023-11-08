<template>
  <div class="d-flex">
    <b-button
      size="lg"
      variant="light"
      class="flex-fill"
      @click="showModal = true"
    >
      {{ $t('allRecords.columns.title') }}
    </b-button>

    <b-modal
      id="columns"
      v-model="showModal"
      size="lg"
      scrollable
      :ok-title="$t('general.label.saveAndClose')"
      cancel-variant="link"
      body-class="p-0"
      title-class="d-flex align-items-center p-0"
      @ok="onSave"
    >
      <template #modal-title>
        {{ $t('allRecords.columns.title') }}
        <c-hint
          :tooltip="$t('allRecords.tooltip.configureColumns')"
          icon-class="text-warning"
        />
      </template>

      <b-card-body
        class="d-flex flex-column mh-100"
      >
        <field-picker
          :module="module"
          :fields.sync="filteredFields"
          style="height: 71vh;"
        />
      </b-card-body>
    </b-modal>
  </div>
</template>

<script>
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  components: {
    FieldPicker,
  },

  props: {
    module: {
      type: Object,
      required: true,
      default: () => ({}),
    },

    fields: {
      type: Array,
      required: true,
      default: () => [],
    },
  },

  data () {
    return {
      showModal: false,

      filteredFields: [],
    }
  },

  created () {
    this.filteredFields = this.fields.map(f => {
      return { ...f.moduleField }
    })
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    onSave () {
      this.$emit('updateFields', this.filteredFields)
    },

    setDefaultValues () {
      this.filteredFields = []
    },
  },
}
</script>

<style lang="scss" scoped>
.fit-modal {
  max-height: calc(100% - 3.5rem);
}
</style>
