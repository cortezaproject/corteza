<template>
  <div class="flex flex-col gap-3">
    <div class="configurator-section">
      <div class="configurator-section__title">{{ $t('configurator.configuration') }}</div>
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1">
          <label class="font-medium text-primary">
            {{ $t('steps.trigger.configurator.resource*') }}
          </label>
          <Select
            v-model="item.triggers.resourceType"
            :options="resourceTypeOptions"
            optionLabel="text"
            optionValue="value"
            filter
            :placeholder="$t('steps.trigger.configurator.select-resource-type')"
            class="w-full"
            @change="resourceChanged"
          />
        </div>

        <div
          v-if="item.triggers.resourceType"
          class="flex flex-col gap-1"
        >
          <label class="font-medium text-primary">
            {{ $t('steps.trigger.configurator.event*') }}
          </label>
          <Select
            v-model="item.triggers.eventType"
            :options="eventTypeOptions"
            :optionLabel="getEventTypeLabel"
            optionValue="eventType"
            filter
            :placeholder="$t('steps.trigger.configurator.select-event-type')"
            class="w-full"
            @change="eventChanged"
          />
        </div>

        <div class="flex items-center gap-2">
          <Checkbox
            v-model="item.triggers.enabled"
            :binary="true"
            inputId="trigger-enabled"
            :disabled="isSubworkflow && !item.triggers.enabled"
            @change="enabledChanged()"
          />
          <label for="trigger-enabled" class="text-primary">
            {{ $t('general.enabled') }}
          </label>
        </div>
      </div>
    </div>

    <div
      v-if="showConstraints"
      class="configurator-section"
    >
      <div class="configurator-section__title">
        {{ $t('steps.trigger.configurator.constraints') }}
        <Button
          v-if="constraintNameTypes.length"
          :label="$t('steps.trigger.configurator.add-constraints')"
          severity="secondary"
          size="small"
          text
          class="ml-auto"
          @click="addConstraint()"
        />
      </div>

      <div v-if="constraintNameTypes.length">
        <div
          v-if="!item.triggers.constraints.length"
          class="text-center text-muted-color py-4"
        >
          {{ $t('steps.trigger.configurator.no-constraints') }}
        </div>

        <div
          v-for="(c, index) in item.triggers.constraints"
          :key="index"
          class="border border-surface rounded-border mb-3"
        >
          <div
            class="flex items-center justify-between px-3 py-2 cursor-pointer hover:bg-emphasis"
            @click="c._showDetails = !c._showDetails"
          >
            <span class="truncate">{{ getConstraintNameLabel(c.name) }}</span>
            <span class="text-muted-color text-sm">{{ getConstraintOperatorLabel(c.op) }}</span>
            <span class="truncate flex-1 text-right">{{ c.values.join(' or ') }}</span>
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              rounded
              size="small"
              class="ml-2"
              @click.stop="removeConstraint(index)"
            />
          </div>

          <div
            v-if="c._showDetails"
            class="px-3 py-3 border-t border-surface bg-emphasis flex flex-col gap-3"
          >
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary text-sm">
                {{ $t('steps.trigger.configurator.resource') }}
              </label>
              <Select
                v-model="c.name"
                :options="constraintNameTypes"
                optionLabel="text"
                optionValue="value"
                filter
                :placeholder="$t('steps.trigger.configurator.select-constraint-type')"
                class="w-full"
                @change="emitChange"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary text-sm">
                {{ $t('steps.trigger.configurator.operator') }}
              </label>
              <Select
                v-model="c.op"
                :options="constraintOperatorTypes"
                optionLabel="text"
                optionValue="value"
                :placeholder="$t('steps.trigger.configurator.select-operator')"
                class="w-full"
                @change="emitChange"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary text-sm">Values</label>
              <div
                v-for="(value, vIndex) in c.values"
                :key="vIndex"
                class="mb-2"
              >
                <p
                  v-if="vIndex > 0"
                  class="text-center uppercase text-muted-color text-sm mb-2"
                >
                  {{ $t('general.label.or') }}
                </p>

                <div class="flex items-center gap-2">
                  <InputText
                    v-model="c.values[vIndex]"
                    class="flex-1"
                    @input="emitChange"
                  />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    rounded
                    size="small"
                    @click="c.values.splice(vIndex, 1)"
                  />
                </div>
              </div>

              <Button
                :label="$t('steps.trigger.configurator.add')"
                size="small"
                @click="c.values.push('')"
              />
            </div>
          </div>
        </div>
      </div>

      <div
        v-else-if="item.triggers.constraints[0]"
        class="flex flex-col gap-1"
      >
        <label class="flex items-center gap-2 font-medium text-primary">
          {{ item.triggers.eventType.replace('on', '') }}
          <a
            :href="intervalDocumentationURL"
            target="_blank"
            class="text-muted-color hover:text-color"
          >
            <i class="pi pi-question-circle" />
          </a>
        </label>

        <CInputDateTime
          v-if="item.triggers.eventType === 'onTimestamp'"
          v-model="item.triggers.constraints[0].values[0]"
          :labels="{
            clear: $t('general.clear'),
            none: $t('general.none'),
            now: $t('general.now'),
            today: $t('general.today'),
          }"
          @update:modelValue="emitChange"
        />

        <InputText
          v-else
          v-model="item.triggers.constraints[0].values[0]"
          @input="emitChange"
        />
      </div>
    </div>

    <div
      v-if="(eventType.properties || []).length"
      class="configurator-section"
    >
      <div class="configurator-section__title">{{ $t('steps.trigger.configurator.initial-scope') }}</div>
      <DataTable
        :value="eventType.properties || []"
        class="border border-surface rounded-border"
      >
        <Column field="name" :header="$t('general.label.name')" />
        <Column field="type" :header="$t('general.label.type')">
          <template #body="{ data }">
            <var>{{ data.type || $t('general.label.any') }}</var>
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<script>
import base from './base.vue'
import { components } from '@planetcrust/human-vue'
import { getConstraintNameLabel } from '../../lib/constraint'
import { getDocumentationURL } from '../../lib/version'
import { camelToTitle } from '../../lib/string'
import eventBus from '../../lib/eventBus'
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
        { value: '=', text: this.$t('steps.trigger.configurator.equal') },
        { value: '!=', text: this.$t('steps.trigger.configurator.not-equal') },
        { value: 'like', text: this.$t('steps.trigger.configurator.like') },
        { value: 'not like', text: this.$t('steps.trigger.configurator.not-like') },
      ]
    },

    intervalDocumentationURL () {
      return getDocumentationURL('integrator-guide/automation/workflows/index.html#deferred-interval')
    },
  },

  async created () {
    if (!this.item.triggers) {
      this.item['triggers'] = {
        resourceType: null,
        eventType: null,
        constraints: [],
        enabled: true,
      }
    }

    await this.getEventTypes()
  },

  methods: {
    getConstraintNameLabel,

    emitChange () {
      eventBus.emit('change-detected')
    },

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
        .catch(e => this.toast.add({ severity: 'error', summary: this.$t('steps.trigger.configurator.failed-fetch-event-types'), detail: e?.message, life: 5000 }))
    },

    addConstraint () {
      this.item.triggers.constraints.push({
        name: '',
        op: '=',
        values: [''],
        _showDetails: true,
      })

      eventBus.emit('change-detected')
    },

    removeConstraint (index) {
      this.item.triggers.constraints.splice(index, 1)
      eventBus.emit('change-detected')
    },

    resourceChanged () {
      this.item.triggers.eventType = null
      this.item.triggers.constraints = []
      eventBus.emit('change-detected')
      this.updateDefaultName()
    },

    eventChanged () {
      if (['onTimestamp', 'onInterval'].includes(this.item.triggers.eventType)) {
        this.item.triggers.constraints = []
        this.addConstraint()
      }

      eventBus.emit('change-detected')
      this.updateDefaultName()
    },

    enabledChanged () {
      eventBus.emit('trigger-updated', this.item.node)
      eventBus.emit('change-detected')
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
