<template>
  <b-tab class="p-0">
    <template #title>
      {{ $t('automation.label') }}
      <b-badge
        v-if="buttons.length > 0"
        pill
        variant="light"
      >
        {{ buttons.length }}
      </b-badge>
    </template>

    <b-container class="pt-3">
      <b-row>
        <b-col lg="6">
          <b-card
            :title="$t('automation.configuredButtons')"
            footer-class="text-right"
            class="border"
          >
            <draggable
              :list.sync="buttons"
              group="buttons"
              filter=".disabled"
            >
              <b-button
                v-for="(b, i) in buttons"
                :key="i"
                :variant="b.variant || 'primary'"
                class="cursor-move m-1"
                @click="currentButton=b"
              >
                {{ b.label || '-' }}
              </b-button>
            </draggable>
            <template #footer>
              <b-button
                variant="link"
                @click="appendButton({ label: $t('automation.dummyButtonLabel'), variant: 'danger' })"
              >
                {{ $t('automation.addPlaceholderLabel') }}
              </b-button>
              <c-input-confirm
                v-if="buttons.length"
                :text="$t('automation.removeAll')"
                variant="link"
                size="md"
                @confirmed="removeAllButtons"
              />
            </template>
          </b-card>
        </b-col>
        <b-col lg="6">
          <button-editor
            v-if="currentButton"
            :page="page"
            :block.sync="block"
            :button="currentButton"
            :script="currentScript"
            :trigger="currentTrigger"
            :record="record"
            :module="module"
            @delete="deleteButton(currentButton)"
          />
        </b-col>
      </b-row>

      <b-row class="mt-4">
        <b-col>
          <b-card
            v-if="available.length > 0"
            :title="$t('automation.availableScriptsAndWorkflow', { count: available.length })"
            class="border"
          >
            <c-input-search
              v-model="queryAvailable"
              class="mb-1"
              :placeholder="$t('automation.searchPlaceholder')"
            />

            <b-list-group
              v-for="(b, index) in filtered"
              :key="index"
              class="mb-2 cursor-pointer"
              no-gutters
              @click.prevent="appendButton(b)"
            >
              <b-list-group-item>
                <div class="d-flex align-items-center w-100 justify-content-between">
                  <h6>
                    {{ b.label || b.script }}
                    <b-badge
                      v-if="b.workflowID"
                      variant="info"
                    >
                      {{ $t('automation.badge.workflow') }}
                    </b-badge>
                    <b-badge
                      v-else-if="b.script"
                      variant="light"
                    >
                      {{ $t('automation.badge.script') }}
                    </b-badge>
                  </h6>
                  <code v-if="b.label && b.script">{{ b.script }}</code>
                </div>
                <p
                  class="my-2"
                >
                  <span
                    v-if="b.description"
                  >
                    {{ b.description }}
                  </span>

                  <i v-else>
                    {{ $t('automation.noDescription') }}
                  </i>
                </p>

                <var
                  v-if="b.stepID"
                >
                  {{ $t('automation.stepID', { stepID: b.stepID }) }}
                </var>
              </b-list-group-item>
            </b-list-group>
          </b-card>
          <p
            v-else-if="buttons.length === 0"
          >
            {{ $t('automation.noScripts') }}
          </p>
        </b-col>
      </b-row>
    </b-container>
  </b-tab>
</template>
<script>
import { compose } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import draggable from 'vuedraggable'
import base from '../base'
import { words } from 'lodash'
import ButtonEditor from './AutomationTabButtonEditor'
const { CInputSearch } = components

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    ButtonEditor,
    draggable,
    CInputSearch,
  },

  extends: base,

  props: {
    buttons: {
      type: Array,
      required: true,
    },

    page: {
      type: compose.Page,
      required: true,
    },

    block: {
      type: compose.PageBlock,
      required: true,
    },
  },

  data () {
    return {
      currentButton: null,
      queryAvailable: '',

      // Filled on create, see fetchTriggers fn
      triggerButtons: [],
    }
  },

  computed: {
    currentScript () {
      const c = this.currentButton
      if (!c.script) {
        return undefined
      }

      return this.scriptButtons
        .filter(({ script }) => script)
        .find(({ script }) => script === c.script)
    },

    currentTrigger () {
      const c = this.currentButton
      if (!c.workflowID) {
        return undefined
      }

      return this.triggerButtons
        .filter(({ workflowID, stepID }) => workflowID && stepID)
        .find(t => t.workflowID === c.workflowID && t.stepID === c.stepID)
    },

    resourceTypes () {
      const resourceTypes = [
        // Three base types we always include when loading list of
        // available automation scripts
        'compose',
        'compose:namespace',
        'compose:page',
      ]

      if (this.module) {
        resourceTypes.push('compose:module')
      }

      if (this.record) {
        resourceTypes.push('compose:record')
      }

      return resourceTypes
    },

    // Returns all compatible buttons from automation scripts
    //
    // This will be deprecated in the future and the only way to add buttons to the UI will be via workflows
    scriptButtons () {
      // @todo this is not a complete implementation
      //       we need to do a proper filtering via constraint matching
      //       for now, all (available) buttons can be configured
      return this.$UIHooks.Find(this.resourceTypes)
    },

    // Available buttons (compatible w/o ones already added)
    available () {
      const existingScripts = this.buttons.map(b => b.script || `${b.workflowID}-${b.stepID}`)

      return [
        ...this.scriptButtons,
        ...this.triggerButtons,
      ].filter(b => !existingScripts.includes(b.script || `${b.workflowID}-${b.stepID}`))
    },

    filtered () {
      if (!this.queryAvailable) {
        return this.available
      }

      const q = words(this.queryAvailable.toLowerCase())
      return this.available
        .filter(({ script = '', label, description }) => q.every(q => `${script} ${label} ${description}`.toLowerCase().search(q) > -1))
    },
  },

  created () {
    this.fetchTriggers()
  },

  methods: {
    appendButton (newButton) {
      this.currentButton = { ...newButton, variant: newButton.variant || 'primary' }
      this.buttons.push(this.currentButton)
    },

    deleteButton (button = {}) {
      const i = this.buttons.indexOf(button)
      this.buttons.splice(i, 1)
      this.currentButton = undefined
    },

    async fetchTriggers () {
      let aux = []

      // Fetch triggers & workflows a
      return this.$AutomationAPI.triggerList({ eventType: 'onManual' })
        .then(({ set } = {}) => {
          aux = set.map(({ triggerID, workflowID, resourceType, stepID }) => ({ triggerID, workflowID, resourceType, stepID }))

          // Pass on simple array of workflow IDs that we can use
          // in the next query
          return set.map(({ workflowID }) => workflowID)
        })
        .then((workflowID) => {
          // Fetch all related workflows
          return this.$AutomationAPI.workflowList({ workflowID })
        })
        .then(({ set = [] } = {}) => {
          // Map triggers, join them with workflows and extract information
          // pieces needed to construct automation buttons
          this.triggerButtons = aux.map(trigger => {
            const { triggerID, workflowID, stepID, resourceType } = trigger
            const workflow = set.find(wf => wf.workflowID === workflowID)
            if (!workflow) {
              // Can not link to workflow (might be disabled or missing)
              console.warn(
                'trigger referencing an non existing workflow',
                { triggerID, workflowID: trigger.workflowID },
              )
              return null
            }

            const { handle, meta: { name, description } = {} } = workflow
            // Try to get label from workflow name stored in meta obj or from the handle
            let label = name || handle

            const step = workflow.steps.find(s => s.stepID === stepID)
            if (!step) {
              // Can not link to step
              console.warn(
                'trigger referencing an non existing step',
                { triggerID, workflowID, stepID },
              )
              return null
            } else if (step.meta && step.meta.label) {
              // There might be more than
              label = `${label} (${step.meta.label})`
            }

            return {
              label,

              // Trigger info
              workflowID,
              stepID,
              resourceType,

              // Description from workflow; for filtering
              description,

              workflow,
            }
          }).filter(t => !!t)
        })
        .catch(err => {
          console.error(err)
        })
    },

    removeAllButtons () {
      this.buttons.splice(0)
      this.currentButton = undefined
    },
  },
}
</script>
<style lang="scss" scoped>
.cursor-move {
  cursor: move !important;
}

.cursor-pointer {
  cursor: pointer !important;
}
</style>
