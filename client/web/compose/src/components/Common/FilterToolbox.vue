<template>
  <b-table-simple
    v-if="mockModule"
    borderless
    class="mb-0"
  >
    <template v-for="(filterGroup, groupIndex) in internalFilter">
      <template v-if="filterGroup.filter.length">
        <b-tr
          v-for="(filter, index) in filterGroup.filter"
          :key="`${groupIndex}-${index}`"
          class="pb-2"
        >
          <b-td style="width: 250px;">
            <c-input-select
              v-model="filter.name"
              :options="fields"
              :get-option-key="getOptionKey"
              :clearable="false"
              :placeholder="$t('recordList.filter.fieldPlaceholder')"
              :reduce="(f) => f.name"
              :class="{ 'filter-field-picker': !!filter.name }"
              @input="onChange($event, groupIndex, index)"
            />
          </b-td>

          <b-td
            v-if="getPreparedField(filter.name)"
            style="width: 250px;"
            :class="{ 'px-2': getPreparedField(filter.name) }"
          >
            <b-form-select
              v-if="getPreparedField(filter.name)"
              v-model="filter.operator"
              :options="getOperators(filter.kind, getPreparedField(filter.name))"
              class="d-flex field-operator w-100"
            />
          </b-td>

          <b-td
            v-if="getPreparedField(filter.name)"
            :key="`${getPreparedField(filter.name)?.fieldID}-${filter.name}`"
          >
            <template v-if="isBetweenOperator(filter.operator)">
              <template v-if="getPreparedField(`${filter.name}-start`)">
                <field-editor
                  :field="getPreparedField(`${filter.name}-start`)"
                  :record="filter.record"
                  :module="mockModule"
                  :namespace="namespace"
                  :errors="errors"
                  value-only
                  class="mb-0 field-editor"
                  @change="onValueChange"
                />
                <div class="my-1 text-center w-100">
                  {{ $t("general.label.and") }}
                </div>
                <field-editor
                  :field="getPreparedField(`${filter.name}-end`)"
                  :record="filter.record"
                  :module="mockModule"
                  :namespace="namespace"
                  :errors="errors"
                  value-only
                  class="mb-0 field-editor"
                  @change="onValueChange"
                />
              </template>
            </template>

            <template v-else>
              <field-editor
                :field="getPreparedField(filter.name)"
                :errors="errors"
                :record="filter.record"
                :module="mockModule"
                :namespace="namespace"
                value-only
                class="mb-0 field-editor"
                @change="onValueChange"
              />
            </template>
          </b-td>

          <b-td
            v-if="getPreparedField(filter.name)"
            style="width: 1%;"
          >
            <b-button
              :id="`${groupIndex}-${index}`"
              ref="delete"
              variant="outline-extra-light"
              class="d-block text-dark border-0 h-full ml-2 px-2 mt-1"
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
            :class="{ 'pb-2': groupIndex !== (internalFilter.length - 1) }"
          >
            <div class="group-separator">
              <b-button
                v-if="groupIndex === (internalFilter.length - 1)"
                variant="outline-primary"
                class="btn-add-group d-block py-2 px-3 m-auto bg-white"
                @click="addGroup()"
              >
                <font-awesome-icon
                  :icon="['fas', 'plus']"
                  class="mb-0 h6"
                />
              </b-button>

              <b-form-select
                v-else
                v-model="internalFilter[groupIndex + 1].groupCondition"
                :options="groupConditionOptions"
                size="sm"
                class="group-condition-select"
                @change="onGroupConditionChange"
              />
            </div>
          </b-td>
        </b-tr>
      </template>
    </template>
  </b-table-simple>
</template>

<script>
import { compose, validator } from '@cortezaproject/corteza-js'
import FieldEditor from 'corteza-webapp-compose/src/components/ModuleFields/Editor'
import { isBetweenOperator } from 'corteza-webapp-compose/src/lib/record-filter.js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  name: 'FilterToolbox',

  components: {
    FieldEditor,
  },

  props: {
    // Raw filter data (without Record objects) - same format as parent components use
    value: {
      type: Array,
      default: undefined,
    },

    module: {
      type: compose.Module,
      required: true,
    },

    namespace: {
      type: compose.Namespace,
      required: true,
    },

    selectedField: {
      type: Object,
      default: undefined,
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

      errors: new validator.Validated(),

      mockModule: undefined,
      preparedFields: [],

      // Internal filter format (with Record objects) for editing
      internalFilter: [],

      // Flag to prevent circular updates when loading external data
      isLoadingExternalData: false,
    }
  },

  computed: {
    fields () {
      return [
        ...[...this.mockModule.fields].sort((a, b) => (a.label || a.name).localeCompare(b.label || b.name)),
        ...this.mockModule.systemFields().map((sf) => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(({ isFilterable, canReadRecordValue }) => {
        // Field must be filterable AND user must have permission to read its values
        // Using !== false allows system fields (undefined) and explicitly allowed fields
        return isFilterable && canReadRecordValue !== false
      })
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
      return this.internalFilter.length >= 1 && this.internalFilter[0].filter[0].name
    },

    groupConditionOptions () {
      return [
        { value: 'OR', text: this.$t('recordList.filter.conditions.or') },
        { value: 'AND', text: this.$t('recordList.filter.conditions.and') },
      ]
    },
  },

  watch: {
    module: {
      immediate: true,
      handler (newModule) {
        this.mockModule = new compose.Module(newModule)

        if (this.mockModule) {
          this.prepareFields()
        }
      },
    },

    value: {
      immediate: true,
      deep: true,
      handler (rawFilter, oldValue) {
        let internal = []

        if (!rawFilter) {
          internal = [this.createDefaultFilterGroup()]
        } else {
          internal = this.rawFilterToInternal(rawFilter)

          if (!internal.length) {
            internal = [this.createDefaultFilterGroup(this.startEmpty ? undefined : this.resolvedSelectedField)]
          } else {
            this.isLoadingExternalData = true
          }
        }

        this.internalFilter = internal

        this.$nextTick(() => {
          this.isLoadingExternalData = false
        })
      },
    },

    internalFilter: {
      deep: true,
      handler (internalFilter) {
        if (this.isLoadingExternalData) {
          return
        }

        this.$emit('input', this.processInternalFilter(internalFilter))
      },
    },
  },

  methods: {
    isBetweenOperator,

    /**
     * Converts raw filter data (without Record objects) to internal format (with Record objects)
     * Used when loading filter data from parent component
     */
    rawFilterToInternal (recordListFilter = []) {
      if (!recordListFilter.length || !this.mockModule) {
        return []
      }

      return recordListFilter.map(({ filter = [], name, groupCondition = 'OR' }) => {
        filter = filter.map(({ value, ...f } = {}) => {
          f.record = new compose.Record(this.mockModule, {})

          if (this.isBetweenOperator(f.operator)) {
            const field = this.getPreparedField(f.name)
            if (field && field.isSystem) {
              f.record[`${f.name}-start`] = value?.start
              f.record[`${f.name}-end`] = value?.end
            } else {
              f.record.values[`${f.name}-start`] = value?.start
              f.record.values[`${f.name}-end`] = value?.end
            }
          } else if (Object.keys(f.record.values).includes(f.name)) {
            f.record.values[f.name] = value
          } else if (Object.keys(f.record).includes(f.name)) {
            // If its a system field add value to root of record
            f.record[f.name] = value
          }

          return f
        })

        return { filter, name, groupCondition }
      })
    },

    /**
     * Processes internal filter format (with Record objects) to output format (without Record objects)
     * Used when emitting filter data to parent component
     */
    processInternalFilter (filter = []) {
      if (!filter.length || !this.mockModule) {
        return []
      }

      return filter.map(({ filter = [], name, groupCondition }) => {
        filter = filter.map(({ record, ...f }) => {
          if (!f.name || !record) {
            return undefined
          }

          if (this.isBetweenOperator(f.operator)) {
            const field = this.getPreparedField(f.name)

            if (field) {
              f.value = {
                start: field.isSystem ? record[`${f.name}-start`] : record.values[`${f.name}-start`],
                end: field.isSystem ? record[`${f.name}-end`] : record.values[`${f.name}-end`],
              }
            }
          } else if (Object.keys(record.values).includes(f.name)) {
            f.value = record.values[f.name]
          } else if (Object.keys(record).includes(f.name)) {
            f.value = record[f.name]
          }

          return f
        })

        return { filter, name, groupCondition }
      })
    },

    prepareFields () {
      const fields = []

      this.fields.forEach(f => {
        if (f.isMulti) {
          f.isMulti = false
          f.multi = true
        }

        if (f.kind === 'Record') {
          f.options.prefilter = ''
        }

        if (f.kind === 'DateTime') {
          f.options.onlyFutureValues = false
          f.options.onlyPastValues = false
        }

        if (f.kind === 'Number') {
          f.options.min = undefined
          f.options.max = undefined
        }

        fields.push(f)

        if (f.kind === 'DateTime' || f.kind === 'Number') {
          fields.push({ ...f, name: `${f.name}-start` })
          fields.push({ ...f, name: `${f.name}-end` })
        }
      })

      this.preparedFields = fields
    },

    onChange (fieldName, groupIndex, index) {
      const field = this.getPreparedField(fieldName)

      const filterExists = !!(
        this.internalFilter[groupIndex] || { filter: [] }
      ).filter[index]
      if (field && filterExists) {
        this.internalFilter[groupIndex].filter[index].kind = field.kind
        this.internalFilter[groupIndex].filter[index].name = field.name
        this.internalFilter[groupIndex].filter[index].value = undefined
        this.internalFilter[groupIndex].filter[index].operator = field.multi
          ? 'IN'
          : '='
      }
    },

    onValueChange () {
      this.$emit('value-change')
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

    deleteFilter (groupIndex, index) {
      const filterExists = !!(
        this.internalFilter[groupIndex] || { filter: [] }
      ).filter[index]

      if (filterExists) {
        // Delete filter from filterGroup
        this.internalFilter[groupIndex].filter.splice(index, 1)

        // If filter was last in filterGroup, delete filterGroup
        if (!this.internalFilter[groupIndex].filter.length) {
          this.internalFilter.splice(groupIndex, 1)

          // If no more filterGroups, add default back
          if (!this.internalFilter.length) {
            this.internalFilter = [this.createDefaultFilterGroup()]
          }
        }
      }

      this.$emit('value-change')
    },

    getOptionKey ({ name }) {
      return name
    },

    getPreparedField (name = '') {
      if (!this.preparedFields.length) {
        return undefined
      }

      return this.preparedFields.find(f => f.name === name)
    },

    addFilter (groupIndex) {
      if ((this.internalFilter[groupIndex] || {}).filter) {
        this.internalFilter[groupIndex].filter.push(
          this.createDefaultFilter(this.resolvedSelectedField),
        )
      }
    },

    createDefaultFilter (field = {}) {
      return {
        name: field.name,
        operator: field.isMulti ? 'IN' : '=',
        value: undefined,
        kind: field.kind,
        record: new compose.Record(this.mockModule, {}),
      }
    },

    createDefaultFilterGroup (field, groupCondition = 'OR') {
      return {
        filter: [this.createDefaultFilter(field)],
        groupCondition,
      }
    },

    onGroupConditionChange () {
      this.$emit('value-change')
    },

    addGroup () {
      this.internalFilter.push(this.createDefaultFilterGroup(this.resolvedSelectedField))
      this.$emit('value-change')
    },

    setDefaultValues () {
      this.mockModule = undefined
    },
  },
}
</script>

<style lang="scss" scoped>
.group-separator {
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: linear-gradient(to left, lightgray, lightgray);
  background-repeat: no-repeat;
  background-size: 100% 1px;
  background-position: center;
}

td {
  padding: 0;
  padding-bottom: 0.5rem;
}

.btn-add-group {
  &:hover,
  &:active {
    background-color: var(--primary) !important;
    color: var(--white) !important;
  }
}

.group-condition-select {
  width: auto;
  min-width: 80px;
  border: 1px solid var(--light);
}
</style>
