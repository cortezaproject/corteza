<template>
  <b-table-simple
    borderless
    class="mb-0"
  >
    <template v-for="(filterGroup, groupIndex) in value">
      <template v-if="filterGroup.filter.length">
        <b-tr
          v-for="(filter, index) in filterGroup.filter"
          :key="`${groupIndex}-${index}`"
          class="pb-2"
        >
          <b-td style="width: 250px;">
            <c-input-select
              v-model="filter.name"
              :options="fieldOptions"
              :get-option-key="getOptionKey"
              :clearable="false"
              :placeholder="$t('recordList.filter.fieldPlaceholder')"
              :reduce="(f) => f.name"
              :class="{ 'filter-field-picker': !!filter.name }"
              @input="onChange($event, groupIndex, index)"
            />
          </b-td>

          <b-td
            v-if="getField(filter.name)"
            style="width: 250px;"
            :class="{ 'px-2': getField(filter.name) }"
          >
            <b-form-select
              v-if="getField(filter.name)"
              v-model="filter.operator"
              :options="getOperators(filter.kind, getField(filter.name))"
              class="d-flex field-operator w-100"
              @change="updateFilterProperties(filter)"
            />
          </b-td>

          <b-td v-if="getField(filter.name)">
            <template v-if="isBetweenOperator(filter.operator)">
              <template v-if="getField(`${filter.name}-start`)">
                <field-editor
                  v-bind="mock"
                  class="mb-0 field-editor"
                  value-only
                  :field="getField(`${filter.name}-start`)"
                  :record="filter.record"
                  @change="onValueChange"
                />
                <span class="my-1 text-center w-100">
                  {{ $t("general.label.and") }}
                </span>
                <field-editor
                  v-bind="mock"
                  class="mb-0 field-editor"
                  value-only
                  :field="getField(`${filter.name}-end`)"
                  :record="filter.record"
                  @change="onValueChange"
                />
              </template>
            </template>

            <template v-else>
              <field-editor
                v-bind="mock"
                class="mb-0 field-editor"
                value-only
                :field="getField(filter.name)"
                :record="filter.record"
                @change="onValueChange"
              />
            </template>
          </b-td>

          <b-td
            v-if="getField(filter.name)"
            class="align-middle"
            style="width: 1%;"
          >
            <b-button
              :id="`${groupIndex}-${index}`"
              ref="delete"
              variant="outline-extra-light"
              class="d-block text-dark border-0 h-full ml-2 px-2"
              @click="deleteFilter(groupIndex, index)"
            >
              <font-awesome-icon
                :icon="['far', 'trash-alt']"
                size="sm"
              />
            </b-button>
          </b-td>
        </b-tr>

        <b-tr
          v-if="showAddCondition"
          :key="`addFilter-${groupIndex}`"
        >
          <b-td class="pb-0">
            <b-button
              variant="primary"
              size="sm"
              class="d-block mr-auto"
              @click="addFilter(groupIndex)"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                class="mr-1"
              />
              {{ $t("general.label.add") }}
            </b-button>
          </b-td>
        </b-tr>

        <b-tr :key="`groupCondtion-${groupIndex}`">
          <b-td
            colspan="100%"
            class="p-0 justify-content-center"
            :class="{ 'pb-3': filterGroup.groupCondition }"
          >
            <div class="group-separator">
              <div style="height: 20px; width: 100%;" />

              <b-button
                v-if="groupIndex === (value.length - 1)"
                variant="outline-primary"
                class="btn-add-group d-block py-2 px-3 m-auto bg-white  "
                @click="addGroup()"
              >
                <font-awesome-icon
                  :icon="['fas', 'plus']"
                  class="mb-0 h6"
                />
              </b-button>
            </div>
          </b-td>
        </b-tr>
      </template>
    </template>
  </b-table-simple>
</template>

<script>
import { compose } from '@cortezaproject/corteza-js'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'FilterToolbox',

  components: {
    FieldEditor,
  },

  props: {
    value: {
      type: Array,
      required: true,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    selectedField: {
      type: Object,
      default: undefined,
    },

    mock: {
      type: Object,
      required: true,
    },

    resetFilterOnCreated: {
      type: Boolean,
      default: false,
    },

    startEmpty: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      conditions: [
        { value: 'AND', text: this.$t('recordList.filter.conditions.and') },
        { value: 'OR', text: this.$t('recordList.filter.conditions.or') },
      ],
    }
  },

  computed: {
    fieldOptions () {
      return this.fields.map(({ name, label }) => ({
        name,
        label: label || name,
      }))
    },

    fields () {
      return [
        ...[...this.module.fields].sort((a, b) =>
          (a.label || a.name).localeCompare(b.label || b.name)
        ),
        ...this.module.systemFields().map((sf) => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(({ isFilterable }) => isFilterable)
    },

    resolvedSelectedField () {
      if (this.selectedField) {
        return this.selectedField
      } else if (this.fields.length) {
        return this.fields[0]
      }

      return {}
    },

    showAddCondition () {
      return this.value.length >= 1 && this.value[0].filter[0].name
    },
  },

  created () {
    if (this.resetFilterOnCreated || !this.startEmpty) {
      this.$emit('input', [this.createDefaultFilterGroup(undefined, this.startEmpty ? undefined : this.resolvedSelectedField)])
    }
  },

  methods: {
    getField (name = '') {
      const field = name
        ? this.mock.module.fields.find((f) => f.name === name)
        : undefined

      return field ? { ...field } : undefined
    },

    onChange (fieldName, groupIndex, index) {
      const field = this.getField(fieldName)
      const filterExists = !!(
        this.value[groupIndex] || { filter: [] }
      ).filter[index]
      if (field && filterExists) {
        const value = this.value
        const tempFilter = [...value]
        tempFilter[groupIndex].filter[index].kind = field.kind
        tempFilter[groupIndex].filter[index].name = field.name
        tempFilter[groupIndex].filter[index].value = undefined
        tempFilter[groupIndex].filter[index].operator = field.multi
          ? 'IN'
          : '='

        this.$emit('input', value)
      }
      this.$emit('prevent-close')
    },

    onValueChange () {
      this.$emit('prevent-close')
    },

    getOperators (kind, field) {
      const operators = [
        {
          value: '=',
          text: this.$t('recordList.filter.operators.equal'),
        },
        {
          value: '!=',
          text: this.$t('recordList.filter.operators.notEqual'),
        },
      ]

      const inOperators = [
        {
          value: 'IN',
          text: this.$t('recordList.filter.operators.contains'),
        },
        {
          value: 'NOT IN',
          text: this.$t('recordList.filter.operators.notContains'),
        },
      ]

      const lgOperators = [
        {
          value: '>',
          text: this.$t('recordList.filter.operators.greaterThan'),
        },
        {
          value: '<',
          text: this.$t('recordList.filter.operators.lessThan'),
        },
      ]
      const matchOperators = [
        {
          value: 'LIKE',
          text: this.$t('recordList.filter.operators.like'),
        },
        {
          value: 'NOT LIKE',
          text: this.$t('recordList.filter.operators.notLike'),
        },
      ]

      const betweenOperators = [
        {
          value: 'BETWEEN',
          text: this.$t('recordList.filter.operators.between'),
        },
        {
          value: 'NOT BETWEEN',
          text: this.$t('recordList.filter.operators.notBetween'),
        },
      ]

      if (field.multi) {
        return inOperators
      }

      switch (kind) {
        case 'Number':
          return [...operators, ...lgOperators, ...betweenOperators]

        case 'DateTime':
          return [...operators, ...lgOperators, ...betweenOperators]

        case 'String':
        case 'Url':
        case 'Email':
          return [...operators, ...matchOperators]

        default:
          return operators
      }
    },

    updateFilterProperties (filter) {
      if (this.isBetweenOperator(filter.operator)) {
        filter.record.values[`${filter.name}-start`] =
          filter.record.values[`${filter.name}-start`]
        filter.record.values[`${filter.name}-end`] =
          filter.record.values[`${filter.name}-end`]

        const field = this.mock.module.fields.find(
          (f) => f.name === filter.name
        )

        this.mock.module.fields.push({ ...field, name: `${filter.name}-end` })
        this.mock.module.fields.push({
          ...field,
          name: `${filter.name}-start`,
        })
      }
    },

    isBetweenOperator (op) {
      return ['BETWEEN', 'NOT BETWEEN'].includes(op)
    },

    deleteFilter (groupIndex, index) {
      let value = this.value
      const filterExists = !!(
        value[groupIndex] || { filter: [] }
      ).filter[index]

      if (filterExists) {
        // Delete filter from filterGroup
        value[groupIndex].filter.splice(index, 1)

        // If filter was last in filterGroup, delete filterGroup
        if (!value[groupIndex].filter.length) {
          value.splice(groupIndex, 1)

          // If no more filterGroups, add default back
          if (!value.length) {
            value = [this.createDefaultFilterGroup()]
          } else if (groupIndex === value.length) {
            // Reset first filterGroup groupCondition if last filterGroup was deleted
            value[groupIndex - 1].groupCondition = undefined
          }
        }

        this.$emit('input', value)
      }

      this.$emit('prevent-close')
    },

    createDefaultFilter (condition, field = {}) {
      return {
        condition,
        name: field.name,
        operator: field.isMulti ? 'IN' : '=',
        value: undefined,
        kind: field.kind,
        record: new compose.Record(this.mock.module, {}),
      }
    },

    getOptionKey ({ name }) {
      return name
    },

    addFilter (groupIndex) {
      const value = this.value

      if ((value[groupIndex] || {}).filter) {
        value[groupIndex].filter.push(
          this.createDefaultFilter('AND', this.resolvedSelectedField)
        )
      }

      this.$emit('input', value)
      this.$emit('prevent-close')
    },

    createDefaultFilterGroup (groupCondition = undefined, field) {
      return {
        groupCondition,
        filter: [this.createDefaultFilter('Where', field)],
      }
    },

    addGroup () {
      const value = this.value
      value[value.length - 1].groupCondition =
        'AND'
      value.push(
        this.createDefaultFilterGroup(undefined, this.resolvedSelectedField)
      )

      this.$emit('input', value)
      this.$emit('prevent-close')
    },

    processFilter (filterGroup = this.value) {
      return filterGroup.map(({ groupCondition, filter = [], name }) => {
        filter = filter.map(({ record, ...f }) => {
          if (record) {
            f.value = record[f.name] || record.values[f.name]
          }

          if (this.isBetweenOperator(f.operator)) {
            f.value = {
              start: this.getField(f.name).isSystem
                ? record[`${f.name}-start`]
                : record.values[`${f.name}-start`],
              end: this.getField(f.name).isSystem
                ? record[`${f.name}-end`]
                : record.values[`${f.name}-end`],
            }
          }

          return f
        })

        return { groupCondition, filter, name }
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.group-separator {
  background-image: linear-gradient(to left, lightgray, lightgray);
  background-repeat: no-repeat;
  background-size: 100% 1px;
  background-position: center;
}

td {
  padding: 0;
  padding-bottom: 0.5rem;
  vertical-align: middle;
}

.btn-add-group {
  &:hover,
  &:active {
    background-color: var(--primary) !important;
    color: var(--white) !important;
  }
}
</style>
