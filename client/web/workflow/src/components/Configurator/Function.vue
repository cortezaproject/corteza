<template>
  <div
    v-if="!processing"
  >
    <div
      v-if="showFunctionList"
    >
      <div
        v-if="functionTypes.length"
        class="flex flex-col gap-3"
      >
        <CFormGroup :label="$t('steps.function.configurator.type')">
          <Select
            v-model="functionRef"
            :options="functionTypes"
            optionLabel="text"
            optionValue="value"
            filter
            :optionDisabled="f => f.disabled"
            :placeholder="$t('steps.function.configurator.select-function')"
            class="w-full"
          @change="functionChanged(functionRef)"
          />
        </CFormGroup>

        <p
          v-if="functionDescription"
          class="text-muted-color text-sm"
        >
          {{ functionDescription }}
        </p>
      </div>
    </div>

    <Divider v-if="showFunctionList && args.length" />

    <div
      v-if="args.length"
    >
      <label class="text-xs font-semibold uppercase tracking-wider text-muted-color mb-2 block">
        {{ $te('steps.function.configurator.parameters') ? $t('steps.function.configurator.parameters') : 'Parameters' }}
      </label>
      <div
        v-for="(a, index) in args"
        :key="index"
        class="border border-surface rounded-border mb-3"
      >
        <div
          class="flex items-start justify-between gap-2 px-3 py-2 cursor-pointer hover:bg-emphasis"
          @click="a._showDetails = !a._showDetails"
        >
          <div class="flex-1 min-w-0 flex flex-col gap-0.5">
            <div class="flex items-center gap-2 min-w-0">
              <div class="font-medium text-primary truncate">
                <span>{{ a.target }}</span><span v-if="a.required" class="text-red-500">*</span>
              </div>
              <div class="ml-auto flex items-center gap-1.5 shrink-0">
                <span
                  v-if="a.valueType === 'expr'"
                  v-tooltip="$t('steps.function.configurator.expression')"
                  class="inline-flex items-center justify-center w-5 h-5 rounded-full bg-primary text-primary-contrast text-xs font-bold"
                >
                  e
                </span>
                <div v-if="!isWhileIterator" class="text-xs text-muted-color truncate">
                  ({{ a.type }})
                </div>
              </div>
            </div>
            <div class="text-sm truncate">
              <samp v-if="a.valueType === 'expr'" class="text-color">{{ a.expr }}</samp>
              <span v-else class="text-color">{{ a.value }}</span>
            </div>
          </div>
        </div>

        <div
          v-if="a._showDetails"
          class="px-4 py-3 border-t border-surface bg-emphasis flex flex-col gap-3"
        >
          <div
            v-if="(paramTypes[functionRef][a.target] || []).length > 1"
            class="flex flex-col gap-1"
          >
            <Select
              v-model="a.type"
              :options="(paramTypes[functionRef][a.target] || [])"
              filter
              class="w-full"
              @change="emitChange"
            />
          </div>

          <div class="flex flex-col gap-1">
            <div
              v-if="a.valueType === 'value'"
            >
              <Select
                v-if="a.target === 'workflow'"
                v-model="a.value"
                :options="workflowOptions"
                :optionLabel="getWorkflowLabel"
                :optionValue="wf => a.type === 'ID' ? wf.workflowID : wf.handle"
                filter
                :placeholder="$t('steps.function.configurator.search-workflow')"
                class="w-full"
                @filter="searchWorkflows"
                @change="emitChange"
              />

              <Select
                v-else-if="a.input.type === 'select'"
                v-model="a.value"
                :options="a.input.properties.options"
                optionLabel="text"
                optionValue="value"
                filter
                :placeholder="$t('steps.function.configurator.option-select')"
                class="w-full"
                @change="emitChange"
              />

              <div
                v-else-if="a.type === 'Boolean'"
                class="flex items-center gap-2"
              >
                <Checkbox
                  v-model="a.value"
                  :binary="true"
                  :true-value="'true'"
                  :false-value="'false'"
                  :inputId="'arg-bool-' + index"
                  @change="emitChange"
                />
                <label :for="'arg-bool-' + index">{{ a.target }}</label>
              </div>

              <expression-editor
                v-else
                v-model="a.value"
                :auto-complete="false"
                @open="openInEditor(index)"
                @input="emitChange"
              />
            </div>

            <expression-editor
              v-else-if="a.valueType === 'expr'"
              v-model="a.expr"
              show-line-numbers
              @open="openInEditor(index)"
              @input="emitChange"
            />
          </div>

          <div
            v-if="!isWhileIterator"
            class="flex items-center gap-2 justify-end"
          >
            <ToggleSwitch
              v-model="a.valueType"
              :true-value="'expr'"
              :false-value="'value'"
              class="scale-75 origin-right"
              @change="valueTypeChanged(a.valueType, index)"
            />
            <span class="text-sm">{{ $t('steps.function.configurator.expression') }}</span>
          </div>
        </div>
      </div>
    </div>

    <Divider v-if="args.length && (expressionResults || results.length)" />

    <div
      v-if="expressionResults || results.length"
    >
      <div v-if="results.length || expressionResults" class="flex items-center justify-between mb-2">
        <label class="text-xs font-semibold uppercase tracking-wider text-muted-color">
          {{ $te('steps.function.configurator.results') ? $t('steps.function.configurator.results') : 'Results' }}
        </label>
        <Button
          v-if="expressionResults"
          :label="$t('steps.function.configurator.add-result')"
          severity="secondary"
          size="small"
          @click="addResult()"
        />
      </div>
      <div v-if="results.length">
        <expression-table
          v-if="expressionResults"
          value-field="expr"
          :items="results"
          :fields="resultFields"
          :types="fieldTypes"
          @update:items="results = $event"
          @remove="removeResult"
          @open-editor="openInEditor"
        />

        <div v-else>
          <div
            v-for="(a, index) in results"
            :key="index"
            class="border border-surface rounded-border mb-3"
          >
            <div
              class="flex items-start justify-between gap-2 px-3 py-2 cursor-pointer hover:bg-emphasis"
              @click="a._showDetails = !a._showDetails"
            >
              <div class="flex-1 min-w-0 flex flex-col gap-0.5">
                <div class="flex items-center gap-2 min-w-0">
                  <div class="font-medium text-primary truncate">{{ a.target }}</div>
                  <div class="ml-auto text-xs text-muted-color truncate shrink-0">({{ a.type }})</div>
                </div>
                <div class="text-sm truncate">
                  <samp class="text-color">{{ a.expr }}</samp>
                </div>
              </div>
            </div>

            <div
              v-if="a._showDetails"
              class="px-4 py-3 border-t border-surface bg-emphasis"
            >
              <InputText
                v-model="a.target"
                :placeholder="$t('configurator.target')"
                class="w-full"
                @input="emitChange"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <Dialog
      :visible="!!expressionEditor.currentExpression"
      :header="$t('editor.editor')"
      modal
      class="w-full max-w-5xl"
      @update:visible="resetExpression"
    >
      <expression-editor
        v-model="currentExpressionValue"
        :lang="expressionEditor.lang"
        min-height="80vh"
        font-size="18px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />

      <template #footer>
        <div class="flex items-center justify-end w-full gap-2">
          <Button
            :label="$t('general.cancel')"
            severity="secondary"
            text
            @click="resetExpression"
          />
          <Button
            :label="$t('general.save')"
            @click="saveExpression"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script>
import base from './base.vue'
import ExpressionTable from '../ExpressionTable.vue'
import ExpressionEditor from '../ExpressionEditor.vue'
import eventBus from '../../lib/eventBus'

export default {
  components: {
    ExpressionEditor,
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

    resultFields () {
      return [
        {
          key: 'target',
          label: this.$t('steps.function.configurator.target'),
        },
        {
          key: 'type',
          label: this.$t('steps.function.configurator.type'),
        },
        {
          key: 'expr',
          label: this.$t('steps.function.configurator.result'),
        },
      ]
    },

    valueTypes () {
      return [
        { text: this.$t('steps.function.configurator.expression'), value: 'expr' },
        { text: this.$t('steps.function.configurator.constant'), value: 'value' },
      ]
    },

    defaultOptions () {
      return [{ value: null, text: this.$t('steps.function.configurator.option-select'), disabled: true }]
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
  },

  watch: {
    'item.config.stepID': {
      immediate: true,
      async handler () {
        this.processing = true

        this.item.config['arguments'] = this.item.config.arguments || []
        this.item.config['results'] = this.item.config.results || []

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

    emitChange () {
      eventBus.emit('change-detected')
    },

    setParams (fName, immediate = false) {
      this.args = []
      this.results = []

      if (!immediate) {
        eventBus.emit('change-detected')
      }

      if (fName) {
        const func = this.functions.find(({ ref }) => ref === fName)

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
            valueType: arg.expr !== undefined ? 'expr' : 'value',
            value: arg.value || input.default || null,
            expr: arg.expr || arg.source || null,
            required: param.required || false,
            input,
          }
        }) || []

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
        this.args[currentIndex] = currentExpression
        eventBus.emit('change-detected')
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
        .catch(e => this.toast.add({ severity: 'error', summary: this.$t('notification.failed-fetch-functions'), detail: e?.message, life: 5000 }))
    },

    async getTypes () {
      return this.$AutomationAPI.typeList()
        .then(({ set }) => {
          this.fieldTypes = set
        })
        .catch(e => this.toast.add({ severity: 'error', summary: this.$t('notification.fetch-types-failed'), detail: e?.message, life: 5000 }))
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

      eventBus.emit('change-detected')
    },

    addResult () {
      this.results.push({
        target: '',
        expr: '',
        type: 'Any',
        _showDetails: true,
      })
      eventBus.emit('change-detected')
    },

    removeResult (index) {
      this.results.splice(index, 1)
      eventBus.emit('change-detected')
    },

    getTypeDescription (type) {
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
