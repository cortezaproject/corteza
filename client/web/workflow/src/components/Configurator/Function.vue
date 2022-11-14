<template>
  <div
    v-if="!processing"
  >
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
          :label="$t('steps:function.configurator.type*')"
          label-class="text-primary"
          class="mb-0"
        >
          <vue-select
            v-model="item.config.ref"
            :options="functionTypes"
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
        class="d-flex align-items-center bg-white p-4"
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
          head-row-variant="secondary"
          details-td-class="bg-white"
          class="mb-4"
          :items="args"
          :fields="argumentFields"
          :tbody-tr-class="rowClass"
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
                v-if="(paramTypes[item.config.ref][a.target] || []).length > 1"
                label-class="text-primary"
              >
                <vue-select
                  v-model="a.type"
                  :options="(paramTypes[item.config.ref][a.target] || [])"
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
                  <vue-select
                    v-if="a.input.type === 'select'"
                    v-model="a.value"
                    :options="a.input.properties.options"
                    label="text"
                    :filter="varFilter"
                    :reduce="a => a.value"
                    :placeholder="$t('steps:function.configurator.option-select')"
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

                  <expression-editor
                    v-else
                    :value.sync="a.value"
                    @open="openInEditor(index)"
                    @input="$root.$emit('change-detected')"
                  />
                </div>

                <expression-editor
                  v-else-if="a.valueType === 'expr'"
                  :value.sync="a.expr"
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
      v-if="results.length"
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
          {{ $t('steps:function.configurator.results') }}
        </h5>
      </b-card-header>

      <b-card-body
        class="p-0"
      >
        <b-table
          id="results"
          fixed
          borderless
          hover
          head-row-variant="secondary"
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

    <b-modal
      id="expression-editor"
      :visible="!!expressionEditor.currentExpression"
      :title="$t('editor:editor')"
      size="lg"
      :ok-title="$t('general:save')"
      :cancel-title="$t('general:cancel')"
      body-class="p-0"
      @ok="saveExpression"
      @hidden="resetExpression"
    >
      <expression-editor
        :value.sync="currentExpressionValue"
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
import { VueSelect } from 'vue-select'
import ExpressionEditor from '../ExpressionEditor.vue'
import { objectSearchMaker, stringSearchMaker } from '../../lib/filter'

export default {
  components: {
    ExpressionEditor,
    VueSelect,
  },

  extends: base,

  data () {
    return {
      processing: true,

      functions: [],
      args: [],
      results: [],

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
        // {
        //   key: 'type',
        //   thClass: "py-2",
        //   tdClass: 'text-truncate pointer'
        // },
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
          thClass: 'pl-3 py-2',
          tdClass: 'text-truncate pointer',
        },
        {
          key: 'type',
          thClass: 'py-2',
          tdClass: 'text-truncate pointer',
        },
        {
          key: 'expr',
          label: this.$t('steps:function.configurator.result'),
          thClass: 'pr-3 py-2',
          tdClass: 'text-truncate pointer',
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
      return (this.functions.find(({ ref }) => ref === this.item.config.ref) || { meta: {} }).meta.description
    },

    isWhileIterator () {
      if (this.item.config) {
        return this.item.config.kind === 'iterator' && this.item.config.ref === 'loopDo'
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

        this.setParams(this.item.config.ref, true)

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
        this.item.config.results = res.filter(({ target }) => target).map(({ target, expr }) => ({ target, expr }))
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
            value: arg.value || null,
            expr: arg.expr || arg.source || null,
            required: param.required || false,
            input,
          }
        }) || []

        // Set results
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
        .catch(this.defaultErrorHandler(this.$t('notification:failed-fetch-functions')))
    },

    async getTypes () {
      return this.$AutomationAPI.typeList()
        .then(({ set }) => {
          this.types = set
        })
        .catch(this.defaultErrorHandler(this.$t('notification:fetch-types-failed')))
    },

    functionChanged (functionRef) {
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
  },
}
</script>

<style lang="scss">
tr.b-table-details > td {
  padding-top: 0;
}

.border-thick {
  border-left: 4px solid #A7D0E3;
}

.border-thick-transparent {
  border-left: none;
}

.border-thick-transparent td:first-child {
  padding-left: calc(0.75rem + 2px);
}
</style>

<style lang="scss" scoped>
.trash {
  right: 0;
  left: 1;
  top: 0;
  bottom: 0;
}

.arrow-up {
  width: 0;
  height: 0;
  margin: 0 auto;
  border-left: 10px solid transparent;
  border-right: 10px solid transparent;
  border-bottom: 10px solid $light;
}
</style>
