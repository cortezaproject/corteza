<template>
  <div>
    <div
      v-for="(rule, index) in rules"
      :key="index"
    >
      <hr v-if="index">
      <div>
        <div class="d-flex justify-content-between align-items-center">
          <h5 v-html="$t('uniqueValueConstraint', { index: index + 1, interpolation: { escapeValue: false } })" />

          <div class="px-4">
            <c-input-confirm
              @confirmed="rules.splice(index, 1)"
            />
          </div>
        </div>
        <b-form-checkbox
          v-model="rule.strict"
          switch
          class="mt-3"
        >
          {{ $t("preventRecordsSave") }}
        </b-form-checkbox>
      </div>
      <div class="mt-3">
        <b-table-simple
          v-if="rule.constraints && rule.constraints.length > 0"
          borderless
        >
          <thead>
            <tr>
              <th>
                {{ $t("field") }}
              </th>
              <th>
                {{ $t("type") }}
              </th>
              <th style="width: 250px;">
                {{ $t("valueModifiers") }}
              </th>
              <th style="width: 250px;">
                {{ $t("multiValues") }}
              </th>
            </tr>
          </thead>
          <tbody v-if="rule.constraints">
            <tr
              v-for="(constraint, consIndex) in rule.constraints"
              :key="`constraint-${consIndex}`"
            >
              <td>{{ getOptionLabel(getField(constraint.attribute)) }}</td>
              <td>{{ getField(constraint.attribute).kind }}</td>
              <td>
                <b-form-select
                  v-model="constraint.modifier"
                  :options="modifierOptions"
                  size="sm"
                />
              </td>
              <td>
                <b-form-select
                  v-model="constraint.multiValue"
                  :options="multiValueOptions"
                  :disabled="!getField(constraint.attribute).isMulti"
                  size="sm"
                />
              </td>
              <td class="text-right p-0 px-4 align-middle">
                <c-input-confirm
                  button-class="text-right"
                  @confirmed="rule.constraints.splice(consIndex, 1)"
                />
              </td>
            </tr>
          </tbody>
        </b-table-simple>

        <b-form-group>
          <b-input-group>
            <vue-select
              v-model="rule.currentField"
              :placeholder="$t('searchFields')"
              :get-option-label="getOptionLabel"
              :get-option-key="getOptionKey"
              :options="filterFieldOptions(rule)"
              :calculate-position="calculateDropdownPosition"
              class="bg-white"
            />

            <b-input-group-append>
              <b-button
                variant="primary"
                class="px-4"
                @click="updateRuleConstraint(rule)"
              >
                {{ $t("add") }}
              </b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
      </div>
    </div>

    <hr>

    <div class="d-flex justify-content-end">
      <b-button
        size="lg"
        variant="link"
        class="d-flex align-items-center text-decoration-none p-0 mt-3"
        @click="addNewConstraint"
      >
        <font-awesome-icon
          :icon="['fas', 'plus']"
          class="mr-2"
        />
        {{ $t("addNewConstraint") }}
      </b-button>
    </div>
  </div>
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'
import { VueSelect } from 'vue-select'

export default {
  i18nOptions: {
    namespaces: 'module',
    keyPrefix: 'edit.config.uniqueValues',
  },

  components: {
    VueSelect,
  },

  props: {
    module: {
      type: compose.Module,
      required: true,
    },
  },

  computed: {
    rules: {
      get () {
        return this.module.config.recordDeDup.rules
      },
      set (value) {
        this.module.config.recordDeDup.rules = value
      },
    },

    modifierOptions () {
      return [
        { value: 'ignore-case', text: this.$t('ignoreCase') },
        { value: 'fuzzy-match', text: this.$t('fuzzyMatch') },
        { value: 'sounds-like', text: this.$t('soundsLike') },
        { value: 'case-sensitive', text: this.$t('caseSensitive') },
      ]
    },

    multiValueOptions () {
      return [
        { value: 'one-of', text: this.$t('oneOf') },
        { value: 'equal', text: this.$t('equal') },
      ]
    },
  },

  methods: {
    addNewConstraint () {
      this.rules.push({
        name: '',
        strict: true,
        constraints: [],
      })
    },

    updateRuleConstraint (rule) {
      if (!rule.constraints) {
        rule.constraints = []
      }

      rule.constraints.push({
        attribute: rule.currentField.name,
        modifier: 'case-sensitive',
        multiValue: 'equal',
        type: rule.currentField.kind,
        isMulti: rule.currentField.isMulti,
      })

      rule.currentField = undefined
    },

    filterFieldOptions (rule) {
      const selectedFields = rule.constraints ? rule.constraints.map(({ attribute }) => attribute) : []
      return this.module.fields.filter(({ name }) => !selectedFields.includes(name))
    },

    getField (attribute) {
      const field = this.module.fields.find(
        ({ name }) => name === attribute,
      )

      return field || {}
    },

    getOptionLabel ({ kind, label, name }) {
      return label || name || kind
    },

    getOptionKey ({ fieldID }) {
      return fieldID
    },
  },
}
</script>
