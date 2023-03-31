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
        style="height: 100% !important;"
        @update:value="setOptions(index, field, $event)"
      />
    </template>

    <template #metric-options="{ metric }">
      <b-row
        class="text-primary"
      >
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.cumulative')"
          >
            <c-input-checkbox
              v-model="metric.cumulative"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.relative')"
          >
            <c-input-checkbox
              v-model="metric.relativeValue"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('edit.metric.fixTooltips')"
          >
            <c-input-checkbox
              v-model="metric.fixTooltips"
              switch
              :labels="checkboxLabel"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>
  </report-edit>
</template>

<script>
import base from './base'
import ReportEdit from './ReportEdit'
import { components } from '@cortezaproject/corteza-vue'
const { CItemPicker } = components

export default {
  name: 'FunnelChart',

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
