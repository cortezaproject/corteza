<template>
  <div>
    <b-button
      v-b-modal.columns
      size="lg"
      variant="light"
    >
      {{ $t('allRecords.columns.title') }}
    </b-button>
    <b-modal
      id="columns"
      size="lg"
      scrollable
      :title="$t('allRecords.columns.title')"
      :ok-title="$t('general.label.saveAndClose')"
      cancel-variant="link"
      body-class="p-0"
      @ok="onSave"
    >
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
      filteredFields: [],
    }
  },

  created () {
    this.filteredFields = this.fields.map(f => {
      return { ...f.moduleField }
    })
  },

  methods: {
    onSave () {
      this.$emit('updateFields', this.filteredFields)
    },
  },
}
</script>

<style lang="scss" scoped>
.fit-modal {
  max-height: calc(100% - 3.5rem);
}
</style>
