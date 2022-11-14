<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
    :dimension-field-kind="['Select']"
    :supported-metrics="1"
    un-skippable
  >
    <template #dimension-options="{ index, dimension, field }">
      <c-item-picker
        v-if="showPicker(field)"
        :value="getOptions(dimension)"
        :options="field.options.options"
        :labels="{
          searchPlaceholder:$t('edit.dimension.optionsPicker.searchPlaceholder'),
          availableItems: $t('edit.dimension.optionsPicker.availableItems'),
          selectAllItems: $t('edit.dimension.optionsPicker.selectAllItems'),
          selectedItems: $t('edit.dimension.optionsPicker.selectedItems'),
          unselectAllItems: $t('edit.dimension.optionsPicker.unselectAllItems'),
          noItemsFound: $t('edit.dimension.optionsPicker.noItemsFound'),
        }"
        class="d-flex flex-column"
        style="max-height: 45vh;"
        @update:value="setOptions(index, field, $event)"
      />
    </template>

    <template #metric-options="{ metric }">
      <b-form-group
        horizontal
        :label-cols="2"
        breakpoint="md"
      >
        <b-form-checkbox
          v-model="metric.cumulative"
        >
          {{ $t('edit.metric.cumulative') }}
        </b-form-checkbox>
      </b-form-group>
    </template>
  </report-edit>
</template>

<script>
import base from './base'
import ReportEdit from './ReportEdit'
import { components } from '@cortezaproject/corteza-vue'
const { CItemPicker } = components

export default {
  components: {
    ReportEdit,
    CItemPicker,
  },

  extends: base,

  methods: {
    showPicker (field) {
      return field && field.kind === 'Select' && field.options.options
    },

    getOptions ({ meta = {} }) {
      const { fields = [] } = meta
      return fields.map(({ value }) => value)
    },

    setOptions (index, field, fields) {
      this.editReport.dimensions[index].meta.fields = fields.map(f => {
        const { options = [] } = field.options || {}
        return options.find(({ value }) => value === f)
      })
    },
  },
}
</script>

<style scoped lang="scss">
  .cursor-pointer {
    cursor: pointer;
  }
</style>
