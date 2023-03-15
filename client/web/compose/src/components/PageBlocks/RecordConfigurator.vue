<template>
  <b-tab
    :title="$t('record.label')"
  >
    <fieldset
      class="form-group"
    >
      <label>
        {{ $t('general.module') }}
      </label>

      <input
        v-if="module"
        v-model="module.name"
        class="form-control"
        type="text"
        readonly
      >
    </fieldset>

    <b-form-group
      :label="$t('module:general.fields')"
    >
      <field-picker
        v-if="module"
        :module="module"
        :fields.sync="options.fields"
        style="max-height: 52vh;"
      />
    </b-form-group>

    <hr>

    <div>
      <b-form-group
        label-size="lg"
      >
        <template #label>
          <div
            class="d-flex"
          >
            {{ $t('record.fieldConditions.label') }}

            <b-button
              variant="link"
              class="p-0 ml-1 text-decoration-none border-0"
              :disabled="addRuleDisabled"
              @click="addRule"
            >
              {{ $t('record.fieldConditions.action') }}
            </b-button>

            <b-button
              variant="link"
              :href="`${documentationURL}#value-sanitizers`"
              target="_blank"
              class="p-0 ml-auto"
            >
              {{ $t('record.fieldConditions.help') }}
            </b-button>
          </div>
        </template>
      </b-form-group>

      <b-table-simple
        v-if="block.options.fieldConditions.length"
        borderless
        small
        responsive="lg"
      >
        <b-thead>
          <b-tr>
            <b-th
              class="text-primary"
            >
              {{ $t('record.fieldConditions.field') }}
            </b-th>

            <b-th
              class="text-primary"
            >
              {{ $t('record.fieldConditions.condition') }}
            </b-th>

            <b-th />
          </b-tr>
        </b-thead>

        <b-tbody>
          <b-tr
            v-for="(condition, i) in block.options.fieldConditions"
            :key="i"
          >
            <b-td
              class="align-middle"
              style="width: 33%; min-width: 250px;"
            >
              <vue-select
                v-model="condition.field"
                :options="block.options.fields"
                append-to-body
                :placeholder="$t('record.fieldConditions.selectPlaceholder')"
                :selectable="option => isSelectable(option)"
                :get-option-label="getOptionLabel"
                :reduce="option => option.isSystem ? option.name : option.fieldID"
                :calculate-position="calculatePosition"
                class="field-selector bg-white"
              />
            </b-td>

            <b-td
              class="align-middle"
              style="min-width: 300px;"
            >
              <b-input-group>
                <b-input-group-prepend>
                  <b-button variant="dark">
                    ƒ
                  </b-button>
                </b-input-group-prepend>
                <b-form-input
                  v-model="condition.condition"
                  :placeholder="$t('record.fieldConditions.placeholder')"
                />
              </b-input-group>
            </b-td>

            <b-td
              class="text-center align-middle pr-2"
              style="width: 100px;"
            >
              <c-input-confirm
                @confirmed="deleteRule(i)"
              />
            </b-td>
          </b-tr>
        </b-tbody>
      </b-table-simple>
    </div>
  </b-tab>
</template>
<script>
import base from './base'
import FieldPicker from 'corteza-webapp-compose/src/components/Common/FieldPicker'
import { VueSelect } from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'Record',

  components: {
    FieldPicker,
    VueSelect,
  },

  extends: base,

  computed: {
    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/compose-configuration/index.html`
    },

    addRuleDisabled () {
      return this.block.options.fields.filter(f => !f.isRequired).length === this.block.options.fieldConditions.length
    },
  },

  methods: {
    addRule () {
      this.options.fieldConditions.push({
        field: undefined,
        condition: '',
      })
    },

    deleteRule (i) {
      this.options.fieldConditions.splice(i, 1)
    },

    isSelectable (option) {
      return !this.block.options.fieldConditions.find(({ field }) => field === option.fieldID || field === option.name) && !option.isRequired
    },

    getOptionLabel (option) {
      return option.label || option.name
    },
  },
}
</script>

<style lang="scss">
.field-selector {
  .vs__selected-options {
    flex-wrap: nowrap;
  }

  .vs__selected {
    max-width: 200px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.vs__dropdown-menu .vs__dropdown-option {
  text-overflow: ellipsis;
  overflow: hidden !important;
}
</style>
