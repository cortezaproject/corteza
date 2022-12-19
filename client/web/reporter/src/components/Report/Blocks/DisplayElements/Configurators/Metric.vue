<template>
  <div>
    <b-form-group
      v-if="options.valueColumn !== undefined"
      :label="$t('display-element:metric.configurator.label-column')"
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
          :label="$t('display-element:metric.configurator.format')"
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
          :label="$t('display-element:metric.configurator.prefix')"
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
          :label="$t('display-element:metric.configurator.suffix')"
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
          :label="$t('display-element:metric.configurator.color.text')"
        >
          <b-form-input
            v-model="options.color"
            type="color"
            debounce="300"
          />
        </b-form-group>
      </b-col>
      <b-col>
        <b-form-group
          :label="$t('display-element:metric.configurator.color.background')"
        >
          <b-form-input
            v-model="options.backgroundColor"
            type="color"
            debounce="300"
          />
        </b-form-group>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import base from './base'
import ColumnSelector from 'corteza-webapp-reporter/src/components/Common/ColumnSelector.vue'

export default {
  components: {
    ColumnSelector,
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
