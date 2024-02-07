<template>
  <report-edit
    :report.sync="editReport"
    :modules="modules"
    :supported-metrics="1"
  >
    <template #dimension-options="{ index, dimension, field }">
      <b-form-group
        v-if="showPicker(field)"
        :label="$t('edit.dimension.options.label')"
        label-class="text-primary"
      >
        <c-item-picker
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
      </b-form-group>
    </template>

    <template #metric-options="{ metric, presetFormattedOptions }">
      <b-row>
        <b-col
          cols="12"
          md="6"
          lg="6"
        >
          <b-form-group
            :label="$t('numberFormat')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.metricFormatter.numberFormat"
              placeholder="0.00"
            />
          </b-form-group>
        </b-col>

        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('edit.metric.options.label')"
            label-class="text-primary"
          >
            <b-form-checkbox
              v-model="metric.cumulative"
            >
              {{ $t('edit.metric.cumulative') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="metric.relativeValue"
            >
              {{ $t('edit.metric.relative') }}
            </b-form-checkbox>

            <b-form-checkbox
              v-model="metric.fixTooltips"
            >
              {{ $t('edit.metric.fixTooltips') }}
            </b-form-checkbox>
          </b-form-group>
        </b-col>
      </b-row>

      <b-row>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('prefix')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.metricFormatter.prefix"
              placeholder="USD/mo"
            />
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          md="6"
        >
          <b-form-group
            :label="$t('suffix')"
            label-class="text-primary"
          >
            <b-input
              v-model="report.metricFormatter.suffix"
              placeholder="$"
            />
          </b-form-group>
        </b-col>
        <b-col
          cols="12"
          lg="6"
        >
          <b-form-group
            :label="$t('edit.additionalConfig.tooltip.formatting.presetFormats.label')"
            label-class="text-primary"
            style="white-space: pre-line;"
            :description="presetFormattedOptions.formattedOptionsDescription"
          >
            <b-form-select
              v-model="report.metricFormatter.presetFormat"
              :options="presetFormattedOptions.formatOptions"
            />
          </b-form-group>
        </b-col>
      </b-row>
    </template>

    <template #additional-config="{ report, presetFormattedOptions }">
      <hr>
      <div class="px-3">
        <h5 class="mb-3">
          {{ $t('edit.additionalConfig.tooltip.label') }}
        </h5>

        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('numberFormat')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltipFormatter.numberFormat"
                placeholder="0.00"
              />
            </b-form-group>
          </b-col>

          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('prefix')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltipFormatter.prefix"
                placeholder="USD/mo"
              />
            </b-form-group>
          </b-col>
        </b-row>

        <b-row>
          <b-col
            cols="12"
            md="6"
          >
            <b-form-group
              :label="$t('suffix')"
              label-class="text-primary"
            >
              <b-input
                v-model="report.tooltipFormatter.suffix"
                placeholder="$"
              />
            </b-form-group>
          </b-col>
          <b-col
            cols="12"
            lg="6"
          >
            <b-form-group
              :label="$t('edit.additionalConfig.tooltip.formatting.presetFormats.label')"
              label-class="text-primary"
              style="white-space: pre-line;"
              :description="presetFormattedOptions.formattedOptionsDescription"
            >
              <b-form-select
                v-model="report.tooltipFormatter.presetFormat"
                :options="presetFormattedOptions.formatOptions"
              />
            </b-form-group>
          </b-col>
        </b-row>
      </div>
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
