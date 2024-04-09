<template>
  <div
    v-if="!processing"
  >
    <b-card
      v-if="showFunctionList"
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
        v-if="functionTypes.length"
        class="p-0"
      >
        <b-form-group
          :label="$t('steps:function.configurator.type*')"
          label-class="text-primary"
          class="mb-0"
        >
          <c-input-select
            v-model="functionRef"
            :options="functionTypes"
            :get-option-key="getOptionTypeKey"
            label="text"
            :selectable="f => !f.disabled"
            :reduce="f => f.value"
            :filter="functionFilter"
            :placeholder="$t('steps:function.configurator.select-function')"
            @input="functionChanged"
          />
        </b-form-group>

        <p
          v-if="functionDescription"
          class="mt-3 mb-0"
        >
          {{ functionDescription }}
        </p>
      </b-card-body>
    </b-card>

    <b-card
      v-if="args.length"
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
          {{ $t('steps:function.configurator.arguments') }}
        </h5>
      </b-card-header>

      <b-card-body
        class="p-0"
      >
        <b-table
          id="arguments"
          fixed
          borderless
          hover
          head-variant="light"
          details-td-class="bg-white"
          :items="args"
          :fields="argumentFields"
          :tbody-tr-class="rowClass"
          class="mb-4"
          @row-clicked="item=>$set(item, '_showDetails', !item._showDetails)"
        >
          <template #cell(target)="{ item: a }">
            <var>{{ `${a.target}${a.required ? '*' : ''}` }}</var>
            <samp v-if="!isWhileIterator"> ({{ a.type }})</samp>
          </template>

          <template #cell(type)="{ item: a }">
            <var>{{ a.type }}</var>
          </template>

          <template #cell(value)="{ item: a }">
            <samp>{{ a[a.valueType] }}</samp>
          </template>

          <template #row-details="{ item: a, index }">
            <div class="arrow-up" />

            <b-card
              class="bg-light"
              body-class="px-4 pb-3"
            >
              <b-form-group
                v-if="(paramTypes[functionRef][a.target] || []).length > 1"
                label-class="text-primary"
              >
                <c-input-select
                  v-model="a.type"
                  :options="(paramTypes[functionRef][a.target] || [])"
                  :get-option-key="getOptionParamKey"
                  :filter="argTypeFilter"
                  :clearable="false"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-group
                label-class="d-flex align-items-center text-primary"
                class="mb-0"
              >
                <div
                  v-if="a.valueType === 'value'"
                >
                  <c-input-select
                    v-if="a.target === 'workflow'"
                    v-model="a.value"
                    :options="workflowOptions"
                    :get-option-label="getWorkflowLabel"
                    :get-option-key="getWorkflowKey"
                    :reduce="wf => a.type === 'ID' ? wf.workflowID : wf.handle"
                    :placeholder="$t('steps:function.configurator.search-workflow')"
                    :filterable="false"
                    @input="$root.$emit('change-detected')"
                    @search="searchWorkflows"
                  />
                  <!-- Clearable -->
                  <c-input-select
                    v-else-if="a.input.type === 'select'"
                    v-model="a.value"
                    :options="a.input.properties.options"
                    :get-option-key="getOptionTypeKey"
                    label="text"
                    :filter="varFilter"
                    :reduce="a => a.value"
                    :placeholder="$t('steps:function.configurator.option-select')"
                    :clearable="false"
                    @input="$root.$emit('change-detected')"
                  />

                  <b-form-checkbox
                    v-else-if="a.type === 'Boolean'"
                    v-model="a.value"
                    value="true"
                    size="md"
                    unchecked-value="false"
                    @input="$root.$emit('change-detected')"
                  >
                    {{ a.target }}
                  </b-form-checkbox>

                  <c-ace-editor
                    v-else
                    v-model="a.value"
                    @open="openInEditor(index)"
                    @input="$root.$emit('change-detected')"
                  />
                </div>

                <c-ace-editor
                  v-else-if="a.valueType === 'expr'"
                  v-model="a.expr"
                  lang="javascript"
                  show-line-numbers
                  @open="openInEditor(index)"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>

              <b-form-checkbox
                v-if="!isWhileIterator"
                v-model="a.valueType"
                value="expr"
                unchecked-value="value"
                switch
                size="sm"
                class="float-right mr-2 mt-2"
                @change="valueTypeChanged($event, index)"
              >
                <div
                  class="d-flex"
                >
                  {{ $t('steps:function.configurator.expression') }}
                  <a
                    :href="documentationURL"
                    target="_blank"
                    class="d-flex align-items-center h6 mb-0 ml-1"
                  >
                    <font-awesome-icon
                      :icon="['far', 'question-circle']"
                      class="ml-1"
                    />
                  </a>
                </div>
              </b-form-checkbox>
            </b-card>
          </template>
        </b-table>
      </b-card-body>
    </b-card>

    <b-card
      v-if="expressionResults || results.length"
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
          {{ $t('steps:function.configurator.results') }}
        </h5>
      </b-card-header>

      <b-card-body
        v-if="results.length"
        class="p-0"
      >
        <expression-table
          v-if="expressionResults"
          value-field="expr"
          :items="results"
          :fields="resultFields"
          :types="fieldTypes"
          @remove="removeResult"
          @open-editor="openInEditor"
        />

        <b-table
          v-else
          id="results"
          fixed
          borderless
          hover
          head-variant="light"
          details-td-class="bg-white"
          class="mb-4"
          :items="results"
          :fields="resultFields"
          :tbody-tr-class="rowClass"
          @row-clicked="item=>$set(item, '_showDetails', !item._showDetails)"
        >
          <template #cell(type)="{ item: a }">
            <var>{{ a.type }}</var>
          </template>

          <template #cell(value)="{ item: a }">
            <samp>{{ a.expr }}</samp>
          </template>

          <template #row-details="{ item: a }">
            <div class="arrow-up" />

            <b-card
              class="bg-light"
              body-class="px-4 pb-3"
            >
              <b-form-group
                class="mb-0"
              >
                <b-form-input
                  v-model="a.target"
                  :placeholder="$t('configurator:target')"
                  @input="$root.$emit('change-detected')"
                />
              </b-form-group>
            </b-card>
          </template>
        </b-table>
      </b-card-body>
    </b-card>

    <portal to="sidebar-footer">
      <b-button
        v-if="expressionResults"
        variant="primary"
        class="align-top border-0 ml-auto"
        @click="addResult()"
      >
        {{ $t('steps:function.configurator.add-result') }}
      </b-button>
    </portal>

    <b-modal
      id="expression-editor"
      :visible="!!expressionEditor.currentExpression"
      :title="$t('editor:editor')"
      size="lg"
      :ok-title="$t('general:save')"
      :cancel-title="$t('general:cancel')"
      body-class="p-0"
      no-fade
      @ok="saveExpression"
      @hidden="resetExpression"
    >
      <c-ace-editor
        v-model="currentExpressionValue"
        :lang="expressionEditor.lang"
        height="500"
        font-size="18px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </b-modal>
  </div>
</template>

<script>
import base from './base'
import ExpressionTable from '../ExpressionTable.vue'
import { objectSearchMaker, stringSearchMaker } from '../../lib/filter'
import { components } from '@cortezaproject/corteza-vue'

const { CAceEditor } = components

export default {
  components: {
    CAceEditor,
    ExpressionTable,
  },

  extends: base,

  data () {
    return {
      processing: true,

      showFunctionList: true,
      expressionResults: false,
      functionRef: undefined,

      functions: [],
      args: [],
      results: [],

      fieldTypes: [],

      paramTypes: {},
      resultTypes: {},

      expressionEditor: {
        currentIndex: undefined,
        currentExpression: undefined,
        lang: 'javascript',
      },
    }
  },

  computed: {
    // Used for expression editor modal
    currentExpressionValue: {
      get () {
        const { currentExpression } = this.expressionEditor
        return currentExpression ? currentExpression[currentExpression.valueType] : ''
      },

      set (value) {
        const { currentExpression } = this.expressionEditor

        if (currentExpression) {
          currentExpression[currentExpression.valueType] = value
        }
      },
    },

    functionTypes () {
      return this.functions.map(({ ref, meta, disabled = false }) => ({ value: ref, text: meta.short, disabled }))
    },

    argumentFields () {
      return [
        {
          key: 'target',
          label: this.$t('steps:function.configurator.name'),
          thClass: 'pl-3 py-2',
          tdClass: 'text-truncate pointer',
        },
        {
          key: 'value',
          thClass: 'pr-3 py-2',
          tdClass: 'text-truncate pointer',
        },
      ]
    },

    resultFields () {
      return [
        {
          key: 'target',
          label: this.$t('steps:function.configurator.target'),
          thClass: 'pl-3',
          tdClass: 'text-truncate pointer',
        },
        {
          key: 'type',
          label: this.$t('steps:function.configurator.type'),
          tdClass: 'text-truncate pointer',
        },
        {
          key: 'expr',
          label: this.$t('steps:function.configurator.result'),
          thClass: 'mr-3',
          tdClass: 'position-relative pointer',
        },
      ]
    },

    valueTypes () {
      return [
        { text: this.$t('steps:function.configurator.expression'), value: 'expr' },
        { text: this.$t('steps:function.configurator.constant'), value: 'value' },
      ]
    },

    defaultOptions () {
      return [{ value: null, text: this.$t('steps:function.configurator.option-select'), disabled: true }]
    },

    functionDescription () {
      return (this.functions.find(({ ref }) => ref === this.functionRef) || { meta: {} }).meta.description
    },

    isWhileIterator () {
      if (this.item.config) {
        return this.item.config.kind === 'iterator' && this.functionRef === 'loopDo'
      }
      return false
    },

    documentationURL () {
      // eslint-disable-next-line no-undef
      const [year, month] = VERSION.split('.')
      return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/integrator-guide/expr/index.html`
    },
  },

  watch: {
    'item.config.stepID': {
      immediate: true,
      async handler () {
        this.processing = true

        this.$set(this.item.config, 'arguments', this.item.config.arguments || [])
        this.$set(this.item.config, 'results', this.item.config.results || [])

        await this.getFunctionTypes()
        await this.getTypes()

        this.functionRef = this.item.config.ref || this.functionRef

        this.setParams(this.functionRef, true)

        this.processing = false
      },
    },

    args: {
      deep: true,
      handler (args) {
        this.item.config.arguments = args.filter(({ value, source, expr }) => value || source || expr)
          .map(arg => {
            const argMapped = {
              target: arg.target,
              type: arg.type,
            }

            argMapped[arg.valueType] = arg[arg.valueType]

            return argMapped
          })
      },
    },

    results: {
      deep: true,
      handler (res) {
        this.item.config.results = res.filter(({ target }) => target).map(({ target, expr, type }) => ({ target, type, expr }))
      },
    },
  },

  methods: {
    functionFilter: objectSearchMaker('text'),
    argTypeFilter: stringSearchMaker(),
    varFilter: objectSearchMaker('text'),

    setParams (fName, immediate = false) {
      this.args = []
      this.results = []

      if (!immediate) {
        this.$root.$emit('change-detected')
      }

      if (fName) {
        const func = this.functions.find(({ ref }) => ref === fName)

        // Set parameters
        if (!this.paramTypes[func.ref] && func.parameters) {
          this.paramTypes[func.ref] = {}
          func.parameters.forEach(({ name, types }) => {
            this.paramTypes[func.ref][name] = types || []
          })
        }

        this.args = func.parameters?.map(param => {
          const arg = this.item.config.arguments.find(({ target }) => target === param.name) || {}
          const { input = {} } = (param.meta || {}).visual || {}
          return {
            name: param.name,
            target: param.name,
            type: arg.type || this.paramTypes[func.ref][param.name][0],
            valueType: this.getValueType(arg, arg.type || this.paramTypes[func.ref][param.name][0], input),
            value: arg.value || input.default || null,
            expr: arg.expr || arg.source || null,
            required: param.required || false,
            input,
          }
        }) || []

        // Set results
        if (!this.expressionResults) {
          if (!this.resultTypes[func.ref] && func.results) {
            this.resultTypes[func.ref] = {}
            func.results.forEach(({ name, types }) => {
              this.resultTypes[func.ref][name] = types || []
            })
          }

          this.results = func.results?.map(result => {
            const res = this.item.config.results.find(({ expr }) => expr === result.name) || {}
            return {
              name: result.name,
              valueType: 'expr',
              target: res.target || undefined,
              type: this.resultTypes[func.ref][result.name][0],
              expr: res.expr || result.name,
            }
          }) || []
        } else {
          this.results = this.item.config.results.map(({ target, type, expr }) => {
            return {
              valueType: 'expr',
              target,
              type,
              expr,
            }
          }) || []
        }
      }
    },

    openInEditor (index = -1) {
      this.expressionEditor = {
        currentIndex: index >= -1 ? index : undefined,
        currentExpression: index >= 0 ? { ...this.args[index] } : undefined,
      }

      this.expressionEditor.lang = this.expressionEditor.currentExpression.valueType === 'expr' ? 'javascript' : 'text'
    },

    saveExpression () {
      const { currentIndex = -1, currentExpression } = this.expressionEditor
      if (currentIndex >= 0) {
        this.args[currentIndex] = currentExpression
        this.$set(this.args, currentIndex, currentExpression)
        this.$root.$emit('change-detected')
      }

      this.resetExpression()
    },

    resetExpression () {
      this.expressionEditor = {
        currentIndex: undefined,
        currentExpression: undefined,
        lang: 'javascript',
      }
    },

    async getFunctionTypes () {
      return this.$AutomationAPI.functionList()
        .then(({ set }) => {
          this.functions = set.filter(({ kind = '' }) => kind !== 'iterator').sort((a, b) => a.meta.short.localeCompare(b.meta.short))
        })
        .catch(this.toastErrorHandler(this.$t('notification:failed-fetch-functions')))
    },

    async getTypes () {
      return this.$AutomationAPI.typeList()
        .then(({ set }) => {
          this.fieldTypes = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:fetch-types-failed')))
    },

    functionChanged (functionRef) {
      this.item.config.ref = functionRef

      this.setParams(functionRef)

      this.$emit('update-default-value', {
        value: (this.functionTypes.find(({ value }) => value === functionRef) || { meta: {} }).text,
        force: !this.item.node.value,
      })
    },

    valueTypeChanged (valueType, index) {
      const oldType = valueType === 'value' ? 'expr' : 'value'
      this.args[index][valueType] = this.args[index][oldType]

      if (!this.args[index].value && this.args[index].type === 'Boolean' && valueType === 'value') {
        this.args[index].value = 'false'
      }

      this.$root.$emit('change-detected')
    },

    getValueType (item, type, input = {}) {
      if (['Boolean'].includes(type) || ['select'].includes(input.type) || this.isWhileIterator) {
        return item.expr ? 'expr' : 'value'
      } else {
        return item.value ? 'value' : 'expr'
      }
    },

    rowClass (item, type) {
      return item._showDetails && type === 'row' ? 'border-thick' : 'border-thick-transparent'
    },

    addResult () {
      this.results.push({
        target: '',
        expr: '',
        type: 'Any',
        _showDetails: true,
      })
      this.$root.$emit('change-detected')
    },

    removeResult (index) {
      this.results.splice(index, 1)
      this.$root.$emit('change-detected')
    },

    getTypeDescription (type) {
      // This will be moved to backend field type information
      const typeDescriptions = {
        ID: 'Make sure to provide the ID in double quotes if you\'re using a literal value. Example "123"',
      }

      return typeDescriptions[type]
    },

    getOptionTypeKey ({ value }) {
      return value
    },

    getOptionEWorkflowLabelKey ({ workflowID }) {
      return workflowID
    },

    getOptionParamKey (type) {
      return type
    },
  },
}
</script>
