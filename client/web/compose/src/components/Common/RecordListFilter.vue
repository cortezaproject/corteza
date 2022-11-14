<template>
  <div>
    <b-button
      :id="popoverTarget"
      variant="link p-0 ml-1"
      :class="[inFilter ? 'text-primary' : 'text-secondary']"
      @click.stop
    >
      <font-awesome-icon :icon="['fas', 'filter']" />
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
                    style="width: 1%;"
                  >
                    <h6
                      v-if="index === 0"
                      class="mb-0"
                    >
                      {{ $t('recordList.filter.where') }}
                    </h6>
                    <b-form-select
                      v-else
                      v-model="filter.condition"
                      :options="conditions"
                    />
                  </b-td>
                  <b-td
                    class="px-2"
                  >
                    <vue-select
                      v-model="filter.name"
                      :options="fieldOptions"
                      :clearable="false"
                      :placeholder="$t('recordList.filter.fieldPlaceholder')"
                      option-value="name"
                      option-text="label"
                      :reduce="f => f.name"
                      append-to-body
                      :calculate-position="calculatePosition"
                      :class="{ 'filter-field-picker': !!filter.name }"
                      class="field-selector bg-white"
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
                      :options="getOperators(filter.kind)"
                      class="d-flex field-operator w-100"
                    />
                  </b-td>
                  <b-td
                    v-if="getField(filter.name)"
                  >
                    <field-editor
                      v-bind="mock"
                      class="field-editor mb-0"
                      value-only
                      :field="getField(filter.name)"
                      :record="filter.record"
                      @change="onValueChange"
                    />
                  </b-td>
                  <b-td
                    v-if="getField(filter.name)"
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
                      variant="link text-decoration-none"
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
                      <b-form-select
                        v-if="filterGroup.groupCondition"
                        v-model="filterGroup.groupCondition"
                        class="w-auto"
                        :options="conditions"
                      />

                      <b-button
                        v-else
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
          class="d-flex justify-content-between bg-white shadow-sm rounded"
        >
          <b-button
            variant="light"
            @click="resetFilter"
          >
            {{ $t('general:label.reset') }}
          </b-button>

          <b-button
            ref="btnSave"
            variant="primary"
            @click="onSave"
          >
            {{ $t('general.label.save') }}
          </b-button>
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
import { VueSelect } from 'vue-select'
import calculatePosition from 'corteza-webapp-compose/src/mixins/vue-select-position'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldEditor,
    VueSelect,
  },

  mixins: [
    calculatePosition,
  ],

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

  created () {
    // Change all module fields to single value to keep multi value fields and single value
    const module = JSON.parse(JSON.stringify(this.module || {}))

    module.fields = [
      ...[...module.fields].map(f => {
        f.isMulti = false
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
        tempFilter[groupIndex].filter[index].operator = '='
        this.componentFilter = tempFilter
      }
      this.onValueChange()
    },

    getOperators (kind) {
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

      if (['Number', 'DateTime'].includes(kind)) {
        return [
          ...operators,
          {
            value: '>',
            text: this.$t('recordList.filter.operators.greaterThan'),
          },
          {
            value: '<',
            text: this.$t('recordList.filter.operators.lessThan'),
          },
        ]
      } else if (['String', 'Url', 'Email'].includes(kind)) {
        return [
          ...operators,
          {
            value: 'LIKE',
            text: this.$t('recordList.filter.operators.contains'),
          },
        ]
      }

      return operators
    },

    getField (name = '') {
      const field = name ? this.mock.module.fields.find(f => f.name === name) : undefined

      return field ? { ...field } : undefined
    },

    createDefaultFilter (condition, field = {}) {
      return {
        condition,
        name: field.name,
        operator: '=',
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

      this.onSave(false)
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

    onSetFocus (groupIndex, index) {
      let focusIndex = this.$refs.delete.findIndex(r => r.id === `${groupIndex}-${index}`)
      if (focusIndex > 0) {
        focusIndex--
      }
      this.$refs.delete[focusIndex].focus()
    },

    onOpen () {
      if (this.recordListFilter.length) {
        // Create record and fill its values property if value exists
        this.componentFilter = this.recordListFilter
          .filter(({ filter = [] }) => filter.some(f => f.name))
          .map(({ groupCondition, filter = [] }) => {
            filter = filter.map(({ value, ...f }) => {
              f.record = new compose.Record(this.mock.module, {})
              f.record.values[f.name] = value
              return f
            })

            return { groupCondition, filter }
          })
      }

      // If no filterGroups, add default
      if (!this.componentFilter.length) {
        this.componentFilter.push(this.createDefaultFilterGroup(undefined, this.selectedField))
      }
    },

    onSave (close = true) {
      if (close) {
        this.$refs.popover.$emit('close')
      }

      // Emit only value and not whole record with every filter
      this.$emit('filter', this.componentFilter.map(({ groupCondition, filter = [] }) => {
        filter = filter.map(({ record, ...f }) => {
          if (record) {
            f.value = record.values[f.name]
          }
          return f
        })

        return { groupCondition, filter }
      }))
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
    color: #2d2d2d;
    text-align: center;
    background: white;
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
      border-bottom-color: white;
      border-top-color: white;
    }

    &::after {
      border-top-color: white;
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
    background-color: $primary !important;
    color: white !important;
  }
}
</style>
