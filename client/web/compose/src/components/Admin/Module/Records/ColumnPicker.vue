<template>
  <div class="d-flex">
    <b-button
      :size="size"
      :variant="variant"
      class="flex-fill"
      :disabled="disabled"
      @click="showModal = true"
    >
      <slot>
        {{ $t('allRecords.columns.title') }}
      </slot>
    </b-button>

    <b-modal
      id="columns"
      v-model="showModal"
      size="lg"
      scrollable
      :ok-title="$t('general.label.saveAndClose')"
      cancel-variant="link"
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

      <field-picker
        :module="module"
        :fields.sync="filteredFields"
        :field-subset="fieldSubset"
        style="height: 71vh;"
      />
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

    disabled: {
      type: Boolean,
      default: false,
    },

    size: {
      type: String,
      default: 'lg',
    },

    variant: {
      type: String,
      default: 'light',
    },

    fieldSubset: {
      type: Array,
      required: false,
      default: () => null,
    },
  },

  data () {
    return {
      showModal: false,

      filteredFields: [],
    }
  },

  watch: {
    fields: {
      immediate: true,
      handler (fields) {
        if (fields) {
          this.filteredFields = this.module.filterFields(fields)
        }
      },
    },
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
