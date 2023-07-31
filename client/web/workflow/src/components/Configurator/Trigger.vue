<template>
  <div>
    <b-card
      class="flex-grow-1 border-bottom border-light rounded-0"
    >
      <b-card-header
        header-tag="header"
        class="bg-white p-0 mb-3"
      >
        <h5
          class="mb-0"
        >
          {{ $t('configurator:configuration') }}
        </h5>
      </b-card-header>
      <b-card-body
        class="p-0"
      >
        <b-form-group
          :label="$t('steps:trigger.configurator.resource*')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="item.triggers.resourceType"
            :options="resourceTypeOptions"
            :get-option-key="getOptionTypeKey"
            label="text"
            :reduce="r => r.value"
            :filter="resTypeFilter"
            :placeholder="$t('steps:trigger.configurator.select-resource-type')"
            @input="resourceChanged"
          />
        </b-form-group>

        <b-form-group
          v-if="item.triggers.resourceType"
          :label="$t('steps:trigger.configurator.event*')"
          label-class="text-primary"
        >
          <c-input-select
            v-model="item.triggers.eventType"
            :options="eventTypeOptions"
            :get-option-key="getOptionEventTypeKey"
            label="eventType"
            :reduce="e => e.eventType"
            :filter="evtTypeFilter"
            :placeholder="$t('steps:trigger.configurator.select-event-type')"
            @input="eventChanged"
          />
        </b-form-group>

        <b-form-group
          class="mb-0"
        >
          <b-form-checkbox
            v-model="item.triggers.enabled"
            :disabled="isSubworkflow && !item.triggers.enabled"
            class="text-primary"
            @change="enabledChanged()"
          >
            {{ $t('general:enabled') }}
          </b-form-checkbox>
        </b-form-group>
      </b-card-body>
    </b-card>

    <b-card
      v-if="showConstraints"
      class="flex-grow-1 border-bottom border-light rounded-0"
      body-class="p-0"
    >
      <b-card-header
        header-tag="header"
        class="d-flex align-items-center bg-white p-4"
      >
        <h5
          class="mb-0"
        >
          {{ $t('steps:trigger.configurator.constraints') }}
        </h5>
        <b-button
          v-if="constraintNameTypes.length"
          variant="primary"
          class="align-top border-0 ml-3"
          @click="addConstraint()"
        >
          {{ $t('steps:trigger.configurator.add-constraints') }}
        </b-button>
      </b-card-header>
      <b-card-body
        class="p-0"
      >
        <b-table
          v-if="constraintNameTypes.length"
          id="constraints"
          fixed
          borderless
          hover
          head-row-variant="secondary"
          details-td-class="bg-white"
          :items="item.triggers.constraints"
          :fields="constraintFields"
          :tbody-tr-class="rowClass"
          @row-clicked="item=>$set(item, '_showDetails', !item._showDetails)"
        >
          <template #cell(name)="{ item: c }">
            <samp v-if="c.name">
              {{ c.name.split('.').map(s => {
                return s[0].toUpperCase() + s.slice(1).toLowerCase()
              }).join(' ')
              }}
            </samp>
          </template>

          <template #cell(values)="{ item: c, index }">
            <div
              class="text-truncate"
              :class="{ 'w-75': c._showDetails}"
            >
              <samp>{{ c.values.join(' or ') }}</samp>
            </div>

            <c-input-confirm
              v-if="c._showDetails"
              show-icon
              class="position-absolute trash"
              size="md"
              @confirmed="removeConstraint(index)"
            />
          </template>

          <template #row-details="{ item: c }">
            <div class="arrow-up" />
            <b-card
              class="bg-light"
            >
              <b-form-group
                :label="$t('steps:trigger.configurator.resource')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="c.name"
                  :options="constraintNameTypes"
                  :get-option-key="getOptionTypeKey"
                  label="text"
                  :reduce="c => c.value"
                  :filter="constrFilter"
                  :placeholder="$t('steps:trigger.configurator.select-constraint-type')"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-group
                :label="$t('steps:trigger.configurator.operator')"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="c.op"
                  :options="constraintOperatorTypes"
                  :get-option-key="getOptionTypeKey"
                  label="text"
                  :reduce="c => c.value"
                  :placeholder="$t('steps:trigger.configurator.select-operator')"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-group>
                <template #label>
                  <div
                    class="d-flex text-primary"
                  >
                    Values
                    <b-button
                      variant="link"
                      class="align-top border-0 p-0 ml-2"
                      @click="c.values.push('')"
                    >
                      {{ $t('steps:trigger.configurator.add') }}
                    </b-button>
                  </div>
                </template>

                <b-input-group
                  v-for="(value, index) in c.values"
                  :key="index"
                  class="mb-2"
                >
                  <b-form-input
                    v-model="c.values[index]"
                    @input="$root.$emit('change-detected')"
                  />

                  <c-input-confirm
                    show-icon
                    @confirmed="c.values.splice(index, 1)"
                  />
                </b-input-group>
              </b-form-group>
            </b-card>
          </template>
        </b-table>

        <b-form-group
          v-else-if="item.triggers.constraints[0]"
          label-class="d-flex align-items-center text-primary"
          class="mt-0 mb-4 mx-4"
        >
          <template #label>
            {{ item.triggers.eventType.replace('on', '') }}
            <a
              :href="intervalDocumentationURL"
              target="_blank"
              class="d-flex align-items-center h6 mb-0 ml-1 pointer"
            >
              <font-awesome-icon
                :icon="['far', 'question-circle']"
              />
            </a>
          </template>
          <c-input-date-time
            v-if="item.triggers.eventType === 'onTimestamp'"
            v-model="item.triggers.constraints[0].values[0]"
            :labels="{
              clear: $t('general:clear'),
              none: $t('general:none'),
              now: $t('general:now'),
              today: $t('general:today'),
            }"
            @input="$root.$emit('change-detected')"
          />

          <b-form-input
            v-else
            v-model="item.triggers.constraints[0].values[0]"
            @input="$root.$emit('change-detected')"
          />
        </b-form-group>
      </b-card-body>
    </b-card>

    <b-card
      v-if="(eventType.properties || []).length"
      class="flex-grow-1 rounded-0"
      body-class="p-0"
    >
      <b-card-header
        header-tag="header"
        class="bg-white p-4"
      >
        <h5
          class="mb-0"
        >
          {{ $t('steps:trigger.configurator.initial-scope') }}
        </h5>
      </b-card-header>
      <b-card-body
        class="p-0"
      >
        <b-table
          id="variable"
          fixed
          borderless
          head-row-variant="secondary"
          class="mb-4"
          :items="eventType.properties || []"
          :fields="scopeFields"
        >
          <template #cell(type)="{ item: v }">
            <var>{{ v.type }}</var>
          </template>
        </b-table>
      </b-card-body>
    </b-card>
  </div>
</template>

<script>
import base from './base'
import { components } from '@cortezaproject/corteza-vue'
import { objectSearchMaker } from '../../lib/filter'
const { CInputDateTime } = components

export default {
  components: {
    CInputDateTime,
  },

  extends: base,

  data () {
    return {
      modules: [],

      eventTypes: [],
      resourceTypes: [],
    }
  },

  computed: {
    resourceTypeOptions () {
      return this.resourceTypes
    },

    eventTypeOptions () {
      return this.eventTypes.filter(({ resourceType }) => resourceType === this.item.triggers.resourceType)
    },

    eventType () {
      return this.eventTypes.find(({ resourceType, eventType }) => resourceType === this.item.triggers.resourceType && eventType === this.item.triggers.eventType) || {}
    },

    showConstraints () {
      if (this.item.triggers.resourceType && this.item.triggers.eventType) {
        return this.constraintNameTypes.length ? true : this.item.triggers.eventType !== 'onManual'
      }
      return false
    },

    constraintFields () {
      return [
        {
          key: 'name',
          thClass: 'pl-3 py-2 w-auto',
          tdClass: 'pr-0 text-truncate pointer',
        },
        {
          key: 'op',
          label: this.$t('steps:trigger.configurator.operator'),
          thClass: 'py-2 operator text-center',
          tdClass: 'pl-0 text-truncate text-center pointer',
        },
        {
          key: 'values',
          thClass: 'pr-3 py-2 w-auto text-center',
          tdClass: 'position-relative pointer text-center',
        },
      ]
    },

    scopeFields () {
      return [
        {
          key: 'name',
          thClass: 'pl-3 py-2',
          tdClass: 'text-truncate',
        },
        {
          key: 'type',
          thClass: 'pr-3 py-2',
          tdClass: 'text-truncate',
        },
      ]
    },

    constraintNameTypes () {
      const constraints = this.eventType.constraints || []

      return constraints.reduce((cons, { name }) => {
        if (!name.includes('*')) {
          cons.push({
            value: name,
            text: name.split('.').map(s => {
              return s[0].toUpperCase() + s.slice(1).toLowerCase()
            }).join(' '),
          })
        }

        return cons
      }, [])
    },

    constraintOperatorTypes () {
      return [
        { value: '=', text: this.$t('steps:trigger.configurator.equal') },
        { value: '!=', text: this.$t('steps:trigger.configurator.not-equal') },
        { value: 'like', text: this.$t('steps:trigger.configurator.like') },
        { value: 'not like', text: this.$t('steps:trigger.configurator.not-like') },

      ]
    },

    intervalDocumentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/automation/workflows/index.html#deferred-interval`
    },
  },

  async created () {
    if (!this.item.triggers) {
      this.$set(this.item, 'triggers', {
        resourceType: null,
        eventType: null,
        constraints: [],
        enabled: true,
      })
    }

    await this.getEventTypes()
  },

  methods: {
    resTypeFilter: objectSearchMaker('text'),
    evtTypeFilter: objectSearchMaker('eventType'),
    constrFilter: objectSearchMaker('text'),

    async getEventTypes () {
      return this.$AutomationAPI.eventTypesList()
        .then(({ set }) => {
          this.eventTypes = set
          const resourceTypes = new Set(set.map(({ resourceType }) => resourceType))
          this.resourceTypes = [...resourceTypes].map(resourceType => {
            return {
              value: resourceType,
              text: resourceType.split(':').map(s => {
                return s[0].toUpperCase() + s.slice(1).toLowerCase()
              }).join(' '),
            }
          })
        })
        .catch(this.toastErrorHandler(this.$t('steps:trigger.configurator.failed-fetch-event-types')))
    },

    addConstraint () {
      this.item.triggers.constraints.push({
        name: null,
        op: '=',
        values: [''],
        _showDetails: true,
      })

      this.$root.$emit('change-detected')
    },

    removeConstraint (index) {
      this.item.triggers.constraints.splice(index, 1)
      this.$root.$emit('change-detected')
    },

    resourceChanged () {
      this.item.triggers.eventType = null
      this.item.triggers.constraints = []
      this.$root.$emit('change-detected')
      this.updateDefaultName()
    },

    eventChanged () {
      this.item.triggers.constraints = []

      if (['onTimestamp', 'onInterval'].includes(this.item.triggers.eventType)) {
        this.addConstraint()
      }

      this.$root.$emit('change-detected')
      this.updateDefaultName()
    },

    enabledChanged () {
      this.$root.$emit('trigger-updated', this.item.node)
      this.$root.$emit('change-detected')
    },

    rowClass (item, type) {
      if (type === 'row') {
        return item._showDetails ? 'border-thick' : 'border-thick-transparent'
      } else if (type === 'row-details') {
        return ''
      }
    },

    updateDefaultName () {
      let { resourceType, eventType } = this.item.triggers

      if (resourceType) {
        eventType = eventType || ''
        let value = [resourceType.split(':').join(' '), eventType].filter(v => v).join(' - ')
        value = value.charAt(0).toUpperCase() + value.slice(1)
        this.$emit('update-default-value', { value, force: !this.item.node.value })
      }
    },

    getOptionTypeKey ({ value }) {
      return value
    },

    getOptionEventTypeKey ({ eventType }) {
      return eventType
    },
  },
}
</script>

<style lang="scss" scoped>
.operator {
  width: 100px;
}
</style>
