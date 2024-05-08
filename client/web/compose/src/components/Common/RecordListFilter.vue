<template>
  <div>
    <b-button
      :id="popoverTarget"
      v-b-tooltip.noninteractive.hover="{ title: $t('recordList.filter.title'), container: '#body' }"
      :variant="variant"
      class="d-flex align-items-center d-print-none border-0 px-1 h-100"
      :class="buttonClass"
      :style="buttonStyle"
      @click.stop
    >
      <font-awesome-icon
        :icon="['fas', 'filter']"
        :class="[inFilter ? 'text-primary' : inactiveIconClass]"
      />
    </b-button>

    <b-popover
      ref="popover"
      custom-class="record-list-filter shadow-sm"
      triggers="click blur"
      placement="bottom"
      delay="0"
      boundary="window"
      boundary-padding="2"
      :target="popoverTarget"
      @show="onOpen"
      @hide="onHide"
    >
      <b-card
        no-body
        class="position-static w-100"
      >
        <b-card-body
          class="px-1 pb-0 overflow-auto"
        >
          <b-table-simple
            v-if="componentFilter.length"
            borderless
            class="mb-0"
          >
            <template
              v-for="(filterGroup, groupIndex) in componentFilter"
            >
              <template v-if="filterGroup.filter.length">
                <b-tr
                  v-for="(filter, index) in filterGroup.filter"
                  :key="`${groupIndex}-${index}`"
                  class="pb-2"
                >
                  <b-td
                    class="px-2"
                  >
                    <c-input-select
                      v-model="filter.name"
                      :options="fieldOptions"
                      :get-option-key="getOptionKey"
                      :clearable="false"
                      :placeholder="$t('recordList.filter.fieldPlaceholder')"
                      :reduce="f => f.name"
                      :class="{ 'filter-field-picker': !!filter.name }"
                      @input="onChange($event, groupIndex, index)"
                    />
                  </b-td>

                  <b-td
                    v-if="getField(filter.name)"
                    :class="{ 'pr-2': getField(filter.name) }"
                  >
                    <b-form-select
                      v-if="getField(filter.name)"
                      v-model="filter.operator"
                      :options="getOperators(filter.kind, getField(filter.name))"
                      class="d-flex field-operator w-100"
                      @change="updateFilterProperties(filter)"
                    />
                  </b-td>
                  <b-td
                    v-if="getField(filter.name)"
                  >
                    <template v-if="isBetweenOperator(filter.operator)">
                      <template
                        v-if="getField(`${filter.name}-start`)"
                      >
                        <field-editor
                          v-bind="mock"
                          class="field-editor mb-0"
                          value-only
                          :field="getField(`${filter.name}-start`)"
                          :record="filter.record"
                          @change="onValueChange"
                        />
                        <span class="text-center my-1 w-100">
                          {{ $t('general.label.and') }}
                        </span>
                        <field-editor
                          v-bind="mock"
                          class="field-editor mb-0"
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
                        class="field-editor mb-0"
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
                      variant="link"
                      class="d-flex align-items-center"
                      @click="deleteFilter(groupIndex, index)"
                    >
                      <font-awesome-icon
                        :icon="['far', 'trash-alt']"
                        size="sm"
                      />
                    </b-button>
                  </b-td>
                </b-tr>

                <b-tr :key="`addFilter-${groupIndex}`">
                  <b-td class="pb-0">
                    <b-button
                      variant="link text-decoration-none d-block mr-auto"
                      style="min-height: 38px; min-width: 84px;"
                      @click="addFilter(groupIndex)"
                    >
                      <font-awesome-icon
                        :icon="['fas', 'plus']"
                        size="sm"
                        class="mr-1"
                      />
                      {{ $t('general.label.add') }}
                    </b-button>
                  </b-td>
                </b-tr>

                <b-tr
                  :key="`groupCondtion-${groupIndex}`"
                >
                  <b-td
                    colspan="100%"
                    class="p-0 justify-content-center"
                    :class="{ 'pb-3': filterGroup.groupCondition }"
                  >
                    <div
                      class="group-separator"
                    >
                      <div style="height: 20px; width: 100%;" />

                      <b-button
                        v-if="groupIndex === (componentFilter.length - 1)"
                        variant="outline-primary"
                        class="btn-add-group bg-white py-2 px-3"
                        @click="addGroup()"
                      >
                        <font-awesome-icon
                          :icon="['fas', 'plus']"
                          class="h6 mb-0 "
                        />
                      </b-button>
                    </div>
                  </b-td>
                </b-tr>
              </template>
            </template>
          </b-table-simple>
        </b-card-body>

        <b-card-footer
          class="d-flex justify-content-between shadow-sm rounded"
        >
          <b-button
            variant="light"
            @click="resetFilter"
          >
            {{ $t('general:label.reset') }}
          </b-button>

          <div class="d-flex">
            <b-button
              v-if="allowFilterPresetSave"
              variant="outline-primary"
              class="mr-2"
              @click="onSave(true, 'filter-preset')"
            >
              {{ $t('recordList.filter.addFilterToPreset') }}
            </b-button>
            <b-button
              ref="btnSave"
              variant="primary"
              @click="onSave"
            >
              {{ $t('general.label.save') }}
            </b-button>
          </div>
        </b-card-footer>
      </b-card>

      <a
        ref="focusMe"
        href=""
        disabled
      />
    </b-popover>
  </div>
</template>
<script>
import FieldEditor from '../ModuleFields/Editor'
import { compose, validator } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldEditor,
  },

  props: {
    target: {
      type: String,
      default: '',
    },

    selectedField: {
      type: Object,
      required: true,
    },

    namespace: {
      type: Object,
      required: true,
    },

    module: {
      type: Object,
      required: true,
    },

    recordListFilter: {
      type: Array,
      required: true,
    },

    variant: {
      type: String,
      default: 'outline-light',
    },

    inactiveIconClass: {
      type: String,
      default: 'text-secondary',
    },

    buttonClass: {
      type: String,
      default: '',
    },

    buttonStyle: {
      type: String,
      default: '',
    },

    allowFilterPresetSave: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      componentFilter: [],

      conditions: [
        { value: 'AND', text: this.$t('recordList.filter.conditions.and') },
        { value: 'OR', text: this.$t('recordList.filter.conditions.or') },
      ],

      mock: {},

      // Used to prevent unwanted closure of popover
      preventPopoverClose: false,
    }
  },

  computed: {
    fields () {
      return [
        ...[...this.module.fields].sort((a, b) =>
          (a.label || a.name).localeCompare(b.label || b.name),
        ),
        ...this.module.systemFields().map(sf => {
          sf.label = this.$t(`field:system.${sf.name}`)
          return sf
        }),
      ].filter(({ isFilterable }) => isFilterable)
    },

    fieldOptions () {
      return this.fields.map(({ name, label }) => ({ name, label: label || name }))
    },

    inFilter () {
      return this.recordListFilter.some(({ filter }) => {
        return filter.some(({ name }) => name === this.selectedField.name)
      })
    },

    popoverTarget () {
      return `${this.target || '0'}-${this.selectedField.name}`
    },
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  created () {
    // Change all module fields to single value to keep multi value fields and single value
    const module = JSON.parse(JSON.stringify(this.module || {}))

    module.fields = [
      ...[...module.fields].map(f => {
        f.multi = f.isMulti
        f.isMulti = false

        // Disable edge case options
        if (f.kind === 'DateTime') {
          f.options.onlyFutureValues = false
          f.options.onlyPastValues = false
        }

        return f
      }),
      ...this.module.systemFields().map(sf => {
        return { ...sf, label: this.$t(`field:system.${sf.name}`) }
      }),
    ]

    this.mock = {
      namespace: this.namespace,
      module,
      errors: new validator.Validated(),
    }
  },

  methods: {
    onHide (e) {
      if (this.preventPopoverClose) {
        e.preventDefault()
        // Focuses invisible element to refocus popover (problems with closing it otherwise)
        this.$nextTick(() => {
          this.$refs.focusMe.focus()
        })
      }
    },

    onValueChange () {
      this.preventPopoverClose = true

      setTimeout(() => {
        this.preventPopoverClose = false
      }, 100)
    },

    onChange (fieldName, groupIndex, index) {
      const field = this.getField(fieldName)
      const filterExists = !!(this.componentFilter[groupIndex] || { filter: [] }).filter[index]
      if (field && filterExists) {
        const tempFilter = [...this.componentFilter]
        tempFilter[groupIndex].filter[index].kind = field.kind
        tempFilter[groupIndex].filter[index].name = field.name
        tempFilter[groupIndex].filter[index].value = undefined
        tempFilter[groupIndex].filter[index].operator = field.multi ? 'IN' : '='
        this.componentFilter = tempFilter
      }
      this.onValueChange()
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

    getField (name = '') {
      const field = name ? this.mock.module.fields.find(f => f.name === name) : undefined

      return field ? { ...field } : undefined
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

    createDefaultFilterGroup (groupCondition = undefined, field) {
      return {
        groupCondition,
        filter: [
          this.createDefaultFilter('Where', field),
        ],
      }
    },

    addFilter (groupIndex) {
      if ((this.componentFilter[groupIndex] || {}).filter) {
        this.componentFilter[groupIndex].filter.push(this.createDefaultFilter('AND', this.selectedField))
      }
    },

    addGroup () {
      this.$refs.btnSave.focus()

      this.componentFilter[this.componentFilter.length - 1].groupCondition = 'AND'
      this.componentFilter.push(this.createDefaultFilterGroup(undefined, this.selectedField))
    },

    resetFilter () {
      this.componentFilter = [
        this.createDefaultFilterGroup(),
      ]
      this.$emit('reset')
    },

    deleteFilter (groupIndex, index) {
      const filterExists = !!(this.componentFilter[groupIndex] || { filter: [] }).filter[index]

      if (filterExists) {
        // Set focus to previous element
        this.onSetFocus(groupIndex, index)
        // Delete filter from filterGroup
        this.componentFilter[groupIndex].filter.splice(index, 1)

        // If filter was last in filterGroup, delete filterGroup
        if (!this.componentFilter[groupIndex].filter.length) {
          this.componentFilter.splice(groupIndex, 1)

          // If no more filterGroups, add default back
          if (!this.componentFilter.length) {
            this.resetFilter()
          } else if (groupIndex === this.componentFilter.length) {
            // Reset first filterGroup groupCondition if last filterGroup was deleted
            this.componentFilter[groupIndex - 1].groupCondition = undefined
          }
        }
      }
    },

    onSetFocus () {
      this.$refs.focusMe.focus()
    },

    onOpen () {
      // Create record and fill its values property if value exists
      this.componentFilter = this.recordListFilter
        .filter(({ filter = [] }) => filter.some(f => f.name))
        .map(({ groupCondition, filter = [], name }) => {
          filter = filter.map(({ value, ...f }) => {
            f.record = new compose.Record(this.mock.module, {})

            if (this.isBetweenOperator(f.operator)) {
              if (this.getField(f.name).isSystem) {
                f.record[`${f.name}-start`] = value.start
                f.record[`${f.name}-end`] = value.end
              } else {
                f.record.values[`${f.name}-start`] = value.start
                f.record.values[`${f.name}-end`] = value.end
              }

              const field = this.mock.module.fields.find(field => field.name === f.name)

              this.mock.module.fields.push({ ...field, name: `${f.name}-end` })
              this.mock.module.fields.push({ ...field, name: `${f.name}-start` })
            } else if (Object.keys(f.record).includes(f.name)) {
              // If its a system field add value to root of record
              f.record[f.name] = value
            } else {
              f.record.values[f.name] = value
            }

            return f
          })

          return { groupCondition, filter, name }
        })

      // If no filterGroups, add default
      if (!this.componentFilter.length) {
        this.componentFilter.push(this.createDefaultFilterGroup(undefined, this.selectedField))
      } else if (!this.inFilter) {
        this.addFilter(0)
      }
    },

    processFilter () {
      return this.componentFilter.map(({ groupCondition, filter = [], name }) => {
        filter = filter.map(({ record, ...f }) => {
          if (!f.name) {
            return
          }

          if (record) {
            f.value = record[f.name] || record.values[f.name]
          }

          if (this.isBetweenOperator(f.operator)) {
            f.value = {
              start: this.getField(f.name).isSystem ? record[`${f.name}-start`] : record.values[`${f.name}-start`],
              end: this.getField(f.name).isSystem ? record[`${f.name}-end`] : record.values[`${f.name}-end`],
            }
          }

          return f
        }).filter(f => f)

        return { groupCondition, filter, name }
      }).filter(({ filter }) => filter.length)
    },

    onSave (close = true, type = 'filter') {
      if (close) {
        this.$refs.popover.$emit('close')
      }

      // Emit only value and not whole record with every filter
      this.$emit(type, this.processFilter())
    },

    updateFilterProperties (filter) {
      if (this.isBetweenOperator(filter.operator)) {
        filter.record.values[`${filter.name}-start`] = filter.record.values[`${filter.name}-start`]
        filter.record.values[`${filter.name}-end`] = filter.record.values[`${filter.name}-end`]

        const field = this.mock.module.fields.find(f => f.name === filter.name)

        this.mock.module.fields.push({ ...field, name: `${filter.name}-end` })
        this.mock.module.fields.push({ ...field, name: `${filter.name}-start` })
      }
    },

    isBetweenOperator (op) {
      return ['BETWEEN', 'NOT BETWEEN'].includes(op)
    },

    getOptionKey ({ name }) {
      return name
    },

    setDefaultValues () {
      this.componentFilter = []
      this.conditions = []
      this.mock = {}
      this.preventPopoverClose = false
    },
  },
}
</script>

<style lang="scss">
.record-list-filter {
  z-index: 1040;
  max-width: 800px !important;
  opacity: 1 !important;
  border-color: transparent;

  .popover-body {
    display: flex;
    width: 800px;
    min-width: min(99vw, 350px);
    max-width: 99vw;
    max-height: 60vh;
    padding: 0;
    color: var(--black);
    text-align: center;
    background: var(--white);
    border-radius: 0.25rem;
    opacity: 1 !important;
    box-shadow: 0 3px 48px #00000026;
    font-size: 0.9rem;
  }

  .v-select, .field-operator, .field-editor {
    min-width: 120px;
  }

  .arrow {
    &::before {
      border-bottom-color: var(--white);
      border-top-color: var(--white);
    }

    &::after {
      border-top-color: var(--white);
    }
  }
}
</style>

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
  &:hover, &:active {
    background-color: var(--primary) !important;
    color: var(--white) !important;
  }
}
</style>
