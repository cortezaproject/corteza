<template>
  <c-form-table-wrapper hide-add-button>
    <b-table-simple
      v-if="render && filter.ref"
      responsive
      borderless
      small
    >
      <template
        v-for="(group, groupIndex) in filter.args"
      >
        <tr
          v-for="(arg, argIndex) in group.args[0].args"
          :key="`${groupIndex}-${argIndex}`"
        >
          <td
            class="fit text-center align-middle pl-0"
          >
            <b-form-select
              v-if="argIndex === 1"
              v-model="group.args[0].ref"
              :options="conditions"
              class="w-auto"
              @change="reRender()"
            />

            <span
              v-else
              class="px-3"
              style="min-width: 60px;"
            >
              {{ argIndex === 0 ? 'Where' : `${group.args[0].ref[0].toUpperCase() + group.args[0].ref.slice(1).toLowerCase()}` }}
            </span>
          </td>

          <template v-if="Object.keys(arg).includes('raw')">
            <td>
              <b-input-group>
                <b-input-group-prepend>
                  <b-button
                    variant="outline-primary"
                    class="d-flex justify-content-center align-items-center"
                    @click="toggleMode(groupIndex, argIndex)"
                  >
                    <font-awesome-icon
                      :icon="['fas', 'filter']"
                      size="sm"
                    />
                  </b-button>
                </b-input-group-prepend>

                <b-form-input
                  v-model="arg.raw"
                  :placeholder="$t('builder:filter-expression')"
                />
              </b-input-group>
            </td>
          </template>

          <template v-else-if="group">
            <td>
              <b-input-group>
                <b-input-group-prepend>
                  <b-button
                    variant="primary"
                    class="d-flex justify-content-center align-items-center"
                    @click="toggleMode(groupIndex, argIndex)"
                  >
                    <font-awesome-icon
                      :icon="['fas', 'filter']"
                      size="sm"
                    />
                  </b-button>
                </b-input-group-prepend>

                <template v-if="group.args[0].args[argIndex].args[0].args[0].value && getColumnData(group.args[0].args[argIndex].args[0].args[1]).multivalue">
                  <column-selector
                    :columns="columns"
                    :value="group.args[0].args[argIndex].args[0].args[1].symbol"
                    @input="setType(groupIndex, argIndex, $event, group.args[0].args[argIndex].args[0].args[0].value['@value'])"
                  />

                  <b-form-select
                    v-model="group.args[0].args[argIndex].args[0].ref"
                    :options="getOperators(getColumnData(group.args[0].args[argIndex].args[0].args[1]))"
                    style="max-width: 120px;"
                  />

                  <b-form-input
                    v-model="group.args[0].args[argIndex].args[0].args[0].value['@value']"
                    :placeholder="$t('builder:value')"
                  />
                </template>

                <template v-else>
                  <column-selector
                    :value="group.args[0].args[argIndex].args[0].args[0].symbol"
                    :columns="columns"
                    @input="setType(groupIndex, argIndex, $event, group.args[0].args[argIndex].args[0].args[1].value['@value'])"
                  />

                  <b-form-select
                    v-model="group.args[0].args[argIndex].args[0].ref"
                    :options="getOperators(getColumnData(group.args[0].args[argIndex].args[0].args[0]))"
                    style="max-width: 120px;"
                  />

                  <b-form-input
                    v-model="group.args[0].args[argIndex].args[0].args[1].value['@value']"
                    :placeholder="$t('builder:value')"
                  />
                </template>
              </b-input-group>
            </td>
          </template>

          <td
            class="fit text-center align-middle pl-2 pr-0"
          >
            <c-input-confirm
              show-icon
              @confirmed="deleteFilter(groupIndex, argIndex)"
            />
          </td>
        </tr>

        <tr
          :key="`${groupIndex}-add`"
        >
          <td
            class="fit align-middle pl-0"
            :class="{ 'text-center': group.args[0].args && group.args[0].args.length }"
          >
            <b-button
              variant="primary"
              size="sm"
              class="mt-1"
              @click="addFilter(groupIndex)"
            >
              <font-awesome-icon
                :icon="['fas', 'plus']"
                size="sm"
                class="mr-1"
              />
              {{ $t('general:label.add') }}
            </b-button>
          </td>
        </tr>

        <tr
          v-if="group.args[0].args && group.args[0].args.length"
          :key="`${groupIndex}-addGroup`"
        >
          <td
            colspan="100%"
            class="p-0 filter-border text-center"
            :class="{ 'pb-1': groupIndex < filter.args.length - 1 }"
          >
            <b-form-select
              v-if="groupIndex < filter.args.length - 1"
              v-model="filter.ref"
              :options="conditions"
              class="w-auto"
              @change="reRender()"
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
          </td>
        </tr>
      </template>
    </b-table-simple>

    <tr
      v-else
    >
      <b-button
        variant="primary"
        size="sm"
        @click="initFilter()"
      >
        <font-awesome-icon
          :icon="['fas', 'plus']"
          size="sm"
          class="mr-1"
        />
        {{ $t('general:label.add') }}
      </b-button>
    </tr>
  </c-form-table-wrapper>
</template>

<script>
import ColumnSelector from 'corteza-webapp-reporter/src/components/Common/ColumnSelector.vue'

export default {
  components: {
    ColumnSelector,
  },

  props: {
    filter: {
      type: Object,
      required: true,
    },

    columns: {
      type: Array,
      default: () => [],
    },
  },

  data () {
    return {
      render: true,

      defaultFilter: {
        ref: 'group',
        // raw: ''
        args: [{
          ref: 'eq',
          args: [
            { symbol: '' },
            { value: { '@type': 'String', '@value': '' } },
          ],
        }],
      },

      conditions: [
        { value: 'and', text: this.$t('general:label.and') },
        { value: 'or', text: this.$t('general:label.or') },
      ],

      operators: [
        {
          value: 'eq',
          text: this.$t('builder:filter.operators.equal'),
          isMulti: false,
        },
        {
          value: 'ne',
          text: this.$t('builder:filter.operators.notEqual'),
          isMulti: false,
        },
        {
          value: 'lt',
          text: this.$t('builder:filter.operators.lessThan'),
          isMulti: false,
        },
        {
          value: 'le',
          text: this.$t('builder:filter.operators.lessThanEqualTo'),
          isMulti: false,
        },
        {
          value: 'gt',
          text: this.$t('builder:filter.operators.greaterThan'),
          isMulti: false,
        },
        {
          value: 'ge',
          text: this.$t('builder:filter.operators.greaterThanEqualTo'),
          isMulti: false,
        },
        {
          value: 'in',
          text: this.$t('builder:filter.operators.contains'),
          isMulti: true,
        },
        {
          value: 'nin',
          text: this.$t('builder:filter.operators.notContains'),
          isMulti: true,
        },
      ],
    }
  },

  methods: {
    initFilter () {
      this.filter.ref = 'and'
      this.filter.args = []
      this.addGroup()
    },

    addGroup () {
      if (this.filter.args) {
        this.filter.args.push({
          ref: 'group',
          args: [
            {
              ref: 'or',
              args: [{
                ref: 'group',
                // raw: ''
                args: [{
                  ref: 'eq',
                  args: [
                    { symbol: '' },
                    { value: { '@type': '', '@value': '' } },
                  ],
                }],
              }],
            },
          ],
        })
      }
      this.reRender()
    },

    addFilter (groupIndex) {
      if (!this.filter.args[groupIndex].args[0].args) {
        this.filter.args[groupIndex].args[0].args = []
      }

      this.filter.args[groupIndex].args[0].args.push({
        ref: 'group',
        // raw: ''
        args: [{
          ref: 'eq',
          args: [
            { symbol: '' },
            { value: { '@type': '', '@value': '' } },
          ],
        }],
      })

      this.reRender()
    },

    deleteFilter (groupIndex, argIndex) {
      const { args } = this.filter.args[groupIndex].args[0]

      if (args) {
        // If last group and last filter, set filter to default
        if (this.filter.args.length === 1 && args.length === 1) {
          delete this.filter.ref
          delete this.filter.args
        } else if (args.length === 1) {
          // If only one left in group, remove group
          this.filter.args.splice(groupIndex, 1)
        } else {
          // Remove filter from group
          args.splice(argIndex, 1)
        }
      }

      this.reRender()
    },

    toggleMode (groupIndex, argIndex) {
      const { args } = this.filter.args[groupIndex].args[0]

      if (args[argIndex]) {
        if (Object.keys(args[argIndex]).includes('raw')) {
          args[argIndex].args = [{
            ref: 'eq',
            args: [
              { symbol: '' },
              { value: { '@type': '', '@value': '' } },
            ],
          }]

          delete args[argIndex].raw
        } else {
          args[argIndex].raw = ''

          delete args[argIndex].args
        }

        this.reRender()
      }
    },

    setType (groupIndex, argIndex, symbol, value) {
      if (!this.filter.args[groupIndex].args[0].args[argIndex]) {
        return
      }

      // Get type
      const { kind = '', multivalue } = this.columns.find(({ name }) => name === symbol) || {}

      // Set type
      if (multivalue) {
        this.filter.args[groupIndex].args[0].args[argIndex].args[0].args[0] = { value: { '@type': kind, '@value': value } }
        this.filter.args[groupIndex].args[0].args[argIndex].args[0].args[1] = { symbol }
        this.filter.args[groupIndex].args[0].args[argIndex].args[0].ref = 'in'
      } else {
        this.filter.args[groupIndex].args[0].args[argIndex].args[0].args[0] = { symbol }
        this.filter.args[groupIndex].args[0].args[argIndex].args[0].args[1] = { value: { '@type': kind, '@value': value } }
        this.filter.args[groupIndex].args[0].args[argIndex].args[0].ref = 'eq'
      }

      this.reRender()
    },

    reRender () {
      this.render = false
      this.$nextTick().then(() => {
        this.render = true
      })
    },

    getColumnData (group) {
      return this.columns.find(({ name }) => name === group.symbol)
    },

    getOperators (column) {
      return column ? this.operators.filter(value => value.isMulti === column.multivalue) : this.operators
    },
  },
}
</script>

<style lang="scss" scoped>
.table td.fit,
.table th.fit {
  white-space: nowrap;
  width: 1%;
}

.btn-add-group {
  &:hover, &:active {
    background-color: var(--primary) !important;
    color: var(--white) !important;
  }
}

.filter-border {
  background-image: linear-gradient(to left, lightgray, lightgray);
  background-repeat: no-repeat;
  background-size: 100% 1px;
  background-position: center;
}
</style>

<style lang="scss">
.prefilter .column-selector {
  .vs__dropdown-toggle {
    border-right: 0;
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
  }
}
</style>
