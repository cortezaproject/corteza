<template>
  <div>
    <b-form-group
      v-if="options.valueColumn !== undefined"
      :label="$t('label-column')"
      label-class="text-primary"
    >
      <column-selector
        v-model="options.valueColumn"
        :columns="valueColumns"
        style="min-width: 100% !important;"
      />
    </b-form-group>

    <b-row>
      <b-col>
        <b-form-group
          :label="$t('format')"
          label-class="text-primary"
        >
          <b-form-input
            v-model="options.format"
            placeholder="0.00"
          />
        </b-form-group>
      </b-col>
      <b-col>
        <b-form-group
          :label="$t('prefix')"
          label-class="text-primary"
        >
          <b-form-input
            v-model="options.prefix"
            placeholder="$"
          />
        </b-form-group>
      </b-col>
      <b-col>
        <b-form-group
          :label="$t('suffix')"
          label-class="text-primary"
        >
          <b-form-input
            v-model="options.suffix"
            placeholder="USD/mo"
          />
        </b-form-group>
      </b-col>
    </b-row>

    <b-row>
      <b-col>
        <b-form-group
          :label="$t('color.text')"
          label-class="text-primary"
        >
          <c-input-color-picker
            v-model="options.color"
            :translations="{
              modalTitle: $t('color.picker'),
              cancelBtnLabel: $t('general:label.cancel'),
              saveBtnLabel: $t('general:label.saveAndClose')
            }"
          />
        </b-form-group>
      </b-col>
      <b-col>
        <b-form-group
          :label="$t('color.background')"
          label-class="text-primary"
        >
          <c-input-color-picker
            v-model="options.backgroundColor"
            :translations="{
              modalTitle: $t('color.picker'),
              cancelBtnLabel: $t('general:label.cancel'),
              saveBtnLabel: $t('general:label.saveAndClose')
            }"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import base from './base'
import ColumnSelector from 'corteza-webapp-reporter/src/components/Common/ColumnSelector.vue'
import { components } from '@cortezaproject/corteza-vue'
const { CInputColorPicker } = components

export default {
  i18nOptions: {
    namespaces: 'display-element',
    keyPrefix: 'metric.configurator',
  },

  components: {
    ColumnSelector,
    CInputColorPicker,
  },

  extends: base,

  computed: {
    valueColumns () {
      const columns = this.columns.length ? this.columns[0] : []
      return [
        ...columns.filter(({ kind }) => ['Number'].includes(kind)),
      ].sort((a, b) => a.label.localeCompare(b.label))
    },
  },
}
</script>
