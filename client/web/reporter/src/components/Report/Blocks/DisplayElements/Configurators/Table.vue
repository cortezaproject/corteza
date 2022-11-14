<template>
  <div
    v-if="options"
  >
    <div
      class="mb-3"
    >
      <h5 class="text-primary mb-2">
        {{ $t('display-element:table.configurator.general') }}
      </h5>

      <b-row no-gutters>
        <b-col class="pr-3">
          <b-form-group
            :label="$t('display-element:table.configurator.table.variant')"
            label-class="text-primary"
          >
            <b-form-select
              v-model="options.tableVariant"
              :options="tableVariants"
            />
          </b-form-group>
        </b-col>
        <b-col>
          <b-form-group
            :label="$t('display-element:table.configurator.head-variant')"
            label-class="text-primary"
          >
            <b-form-radio-group
              v-model="options.headVariant"
              class="mt-lg-2"
            >
              <b-form-radio
                :value="null"
                inline
              >
                {{ $t('display-element:table.configurator.none') }}
              </b-form-radio>
              <b-form-radio
                value="light"
                inline
              >
                {{ $t('display-element:table.configurator.light') }}
              </b-form-radio>
              <b-form-radio
                value="dark"
                inline
              >
                {{ $t('display-element:table.configurator.dark') }}
              </b-form-radio>
            </b-form-radio-group>
          </b-form-group>
        </b-col>
      </b-row>

      <b-form-group
        :label="$t('display-element:table.configurator.table.options.label')"
        label-class="text-primary"
      >
        <b-form-checkbox
          v-model="options.striped"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.striped') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.bordered"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.bordered') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.borderless"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.borderless') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.small"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.small') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.hover"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.hover') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.dark"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.dark') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.responsive"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.responsive') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.fixed"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.fixed') }}
        </b-form-checkbox>
        <b-form-checkbox
          v-model="options.noCollapse"
          inline
        >
          {{ $t('display-element:table.configurator.table.options.no-collapse') }}
        </b-form-checkbox>
      </b-form-group>
    </div>

    <hr>

    <div
      class="mb-3"
    >
      <h5 class="text-primary mb-2">
        {{ $t('display-element:table.configurator.data') }}
      </h5>

      <b-form-group
        v-if="options.datasources.length > 1"
        :label="$t('display-element:table.configurator.joined-datasource-handling')"
        label-class="text-primary"
      >
        <b-form-select
          v-model="currentConfigurableDatasourceName"
          :options="options.datasources"
          text-field="name"
          value-field="name"
        />
      </b-form-group>

      <b-form-group
        v-if="currentConfigurableDatasourceName && currentColumns.length"
        :label="$t('display-element:table.configurator.columns')"
        label-class="text-primary"
      >
        <column-picker
          :all-items="currentColumns"
          :selected-items.sync="currentSelectedColumns"
          class="d-flex flex-column"
        />
      </b-form-group>
    </div>
  </div>
</template>

<script>
import base from './base'
import ColumnPicker from 'corteza-webapp-reporter/src/components/Common/ColumnPicker'

export default {
  components: {
    ColumnPicker,
  },

  extends: base,

  data () {
    return {
      currentConfigurableDatasourceName: undefined,
    }
  },

  computed: {
    tableVariants () {
      return [
        { value: '', text: this.$t('display-element:table.configurator.none') },
        { value: 'primary', text: this.$t('display-element:table.configurator.table.variants.primary') },
        { value: 'secondary', text: this.$t('display-element:table.configurator.table.variants.secondary') },
        { value: 'info', text: this.$t('display-element:table.configurator.table.variants.info') },
        { value: 'danger', text: this.$t('display-element:table.configurator.table.variants.danger') },
        { value: 'warning', text: this.$t('display-element:table.configurator.table.variants.warning') },
        { value: 'success', text: this.$t('display-element:table.configurator.table.variants.success') },
        { value: 'light', text: this.$t('display-element:table.configurator.table.variants.light') },
        { value: 'dark', text: this.$t('display-element:table.configurator.table.variants.dark') },
      ]
    },

    currentColumns () {
      if (this.currentConfigurableDatasourceName && this.columns) {
        const datasourceIndex = this.options.datasources.findIndex(ds => ds.name === this.currentConfigurableDatasourceName)
        if (datasourceIndex >= 0) {
          return this.columns[datasourceIndex] || []
        }
      }

      return []
    },

    currentSelectedColumns: {
      get () {
        return this.currentConfigurableDatasourceName ? this.options.columns[this.currentConfigurableDatasourceName] : []
      },

      set (columns) {
        if (this.currentConfigurableDatasourceName) {
          this.$set(this.options.columns, this.currentConfigurableDatasourceName, columns || [])
        }
      },
    },
  },

  watch: {
    'options.datasources': {
      immediate: true,
      handler (datasources) {
        datasources.forEach(({ name }) => {
          this.$set(this.options.columns, name, this.options.columns[name] || [])
        })

        this.currentConfigurableDatasourceName = (datasources[0] || {}).name
      },
    },
  },
}
</script>
