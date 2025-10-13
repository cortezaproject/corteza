<template>
  <div>
    <b-card
      class="flex-grow-1 border-bottom border-light rounded-0"
    >
      <b-card-header
        header-tag="header"
        class="p-0 mb-3"
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
            :clearable="false"
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
            :get-option-label="getEventTypeLabel"
            :reduce="e => e.eventType"
            :filter="evtTypeFilter"
            :placeholder="$t('steps:trigger.configurator.select-event-type')"
            :clearable="false"
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
        class="d-flex align-items-center"
      >
        <h5
          class="mb-0"
        >
          {{ $t('steps:trigger.configurator.constraints') }}
        </h5>

        <b-button
          v-if="constraintNameTypes.length"
          variant="primary"
          class="ml-3"
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
          head-variant="light"
          details-td-class="bg-white"
          :items="item.triggers.constraints"
          :fields="constraintFields"
          :tbody-tr-class="rowClass"
          show-empty
          class="mb-0"
          @row-clicked="item=>$set(item, '_showDetails', !item._showDetails)"
        >
          <template #empty>
            <p class="text-center text-muted m-4">
              {{ $t('steps:trigger.configurator.no-constraints') }}
            </p>
          </template>

          <template #cell(name)="{ item: c }">
            {{ getConstraintNameLabel(c.name) }}
          </template>

          <template #cell(op)="{ item: c }">
            <p style="margin-bottom: 0.4rem;">
              {{ getConstraintOperatorLabel(c.op) }}
            </p>
          </template>

          <template #cell(values)="{ item: c, index }">
            {{ c.values.join(' or ') }}

            <c-input-confirm
              show-icon
              class="delete-btn ml-auto"
              @confirmed="removeConstraint(index)"
            />
          </template>

          <template #row-details="{ item: c }">
            <div class="arrow-up" />

            <b-card class="bg-light">
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
                  :clearable="false"
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
                  :clearable="false"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-group
                label="Values"
                label-class="text-primary"
              >
                <div
                  v-for="(value, index) in c.values"
                  :key="index"
                  class="mb-2"
                >
                  <p
                    v-if="index > 0"
                    class="text-center text-uppercase text-muted mb-2"
                  >
                    {{ $t('general:label.or') }}
                  </p>

                  <div class="d-flex align-items-center gap-1">
                    <b-form-input
                      v-model="c.values[index]"
                      @input="$root.$emit('change-detected')"
                    />

                    <c-input-confirm
                      show-icon
                      @confirmed="c.values.splice(index, 1)"
                    />
                  </div>
                </div>

                <b-button
                  variant="primary"
                  size="sm"
                  class="mr-auto"
                  @click="c.values.push('')"
                >
                  {{ $t('steps:trigger.configurator.add') }}
                </b-button>
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
          head-variant="light"
          class="mb-4"
          :items="eventType.properties || []"
          :fields="scopeFields"
        >
          <template #cell(type)="{ item: v }">
            <var>{{ v.type || $t('general:label.any') }}</var>
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
import { getConstraintNameLabel } from '../../lib/constraint'
import { getDocumentationURL } from '../../lib/version'
import { camelToTitle } from '../../lib/string'
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
          tdClass: 'pointer',
        },
        {
          key: 'op',
          label: this.$t('steps:trigger.configurator.operator'),
          thClass: 'text-center',
          tdClass: 'text-truncate text-center pointer',
        },
        {
          key: 'values',
          tdClass: 'd-flex align-items-start gap-2 pointer pr-2',
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
            text: this.getConstraintNameLabel(name),
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
      return getDocumentationURL('integrator-guide/automation/workflows/index.html#deferred-interval')
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
    getConstraintNameLabel,

    async getEventTypes () {
      return this.$AutomationAPI.eventTypesList()
        .then(({ set }) => {
          this.eventTypes = set
          const resourceTypes = new Set(set.map(({ resourceType }) => resourceType))
          this.resourceTypes = [...resourceTypes].map(resourceType => {
            return {
              value: resourceType,
              text: this.getResourceTypeLabel(resourceType),
            }
          })
        })
        .catch(this.toastErrorHandler(this.$t('steps:trigger.configurator.failed-fetch-event-types')))
    },

    addConstraint () {
      this.item.triggers.constraints.push({
        name: '',
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
      if (['onTimestamp', 'onInterval'].includes(this.item.triggers.eventType)) {
        this.item.triggers.constraints = []
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
      const { resourceType, eventType } = this.item.triggers

      if (resourceType) {
        let value = [this.getResourceTypeLabel(resourceType), this.getEventTypeLabel({ eventType })].filter(v => v).join(' - ')
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

    getResourceTypeLabel (resourceType) {
      if (!resourceType) return ''

      return resourceType
        .split(':')
        .map(part => part
          .split('-')
          .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
          .join(' '),
        )
        .join(' - ')
    },

    getEventTypeLabel ({ eventType = '' } = {}) {
      if (!eventType) return ''

      return camelToTitle(eventType.replace('on', ''))
    },

    getConstraintOperatorLabel (op) {
      const operator = this.constraintOperatorTypes.find(type => type.value === op)
      return operator ? operator.text : op
    },
  },
}
</script>

<style lang="scss" scoped>
#constraints {
  tbody tr .delete-btn {
    display: none !important;
  }

  tbody tr:hover .delete-btn {
    display: inline-flex !important;
  }
}
</style>
