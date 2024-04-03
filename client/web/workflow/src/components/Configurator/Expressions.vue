<template>
  <div>
    <b-card
      class="flex-grow-1 rounded-0"
      body-class="p-0"
    >
      <b-card-header
        header-tag="header"
        class="d-flex align-items-center py-4"
      >
        <h5
          class="d-flex align-items-center mb-0"
        >
          {{ $t('steps:expressions.label') }}
          <a
            :href="documentationURL"
            target="_blank"
            class="d-flex align-items-center h6 mb-0 ml-1"
          >
            <font-awesome-icon
              :icon="['far', 'question-circle']"
            />
          </a>
        </h5>

        <portal to="sidebar-footer">
          <b-button
            variant="primary"
            class="align-top border-0 ml-auto"
            @click="addArgument()"
          >
            {{ $t('steps:expressions.configurator.add-expression') }}
          </b-button>
        </portal>
      </b-card-header>

      <b-card-body
        v-if="hasArguments"
        class="p-0"
      >
        <expression-table
          value-field="expr"
          :items="item.config.arguments"
          :fields="argumentFields"
          :types="fieldTypes"
          @remove="removeArgument"
          @open-editor="openInEditor"
        />
      </b-card-body>
    </b-card>

    <b-modal
      id="expression-editor"
      :visible="!!expressionEditor.currentExpression"
      :title="$t('editor:editor')"
      size="lg"
      :ok-title="$t('general:save')"
      :cancel-title="$t('general:cancel')"
      cancel-variant="light"
      body-class="p-0"
      no-fade
      @ok="saveExpression"
      @hidden="resetExpression"
    >
      <c-ace-editor
        v-model="currentExpressionValue"
        height="500"
        lang="javascript"
        font-size="18px"
        show-line-numbers
        auto-complete
        :border="false"
        :show-popout="false"
        :auto-complete-suggestions="expressionAutoCompleteValues"
      />
    </b-modal>
  </div>
</template>

<script>
import base from './base'
import ExpressionTable from '../ExpressionTable.vue'
import { components } from '@cortezaproject/corteza-vue'
import { EXPRESSION_EDITOR_AUTO_COMPLETE_VALUES } from '../../lib/editor-auto-complete.js'

const { CAceEditor } = components

export default {
  components: {
    CAceEditor,
    ExpressionTable,
  },

  extends: base,

  data () {
    return {
      fieldTypes: [],

      expressionEditor: {
        currentIndex: undefined,
        currentExpression: undefined,
      },

      expressionAutoCompleteValues: EXPRESSION_EDITOR_AUTO_COMPLETE_VALUES,
    }
  },

  computed: {
    currentExpressionValue: {
      get () {
        return this.expressionEditor.currentExpression ? this.expressionEditor.currentExpression.expr : ''
      },

      set (value) {
        if (this.expressionEditor.currentExpression) {
          this.expressionEditor.currentExpression.expr = value
        }
      },
    },

    argumentFields () {
      return [
        {
          key: 'target',
          label: this.$t('steps:expressions.configurator.target'),
          thClass: 'pl-4 ml-1',
          formatter: (item) => {
            return `${item.target}(${item.type})`
          },
        },
        {
          key: 'expr',
          label: this.$t('steps:expressions.configurator.expression'),
          thClass: 'pl-1 mr-3',
        },
      ]
    },

    hasArguments () {
      const { config } = this.item || {}
      return (config && (config.arguments || []).length) || []
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
      handler () {
        this.$set(this.item.config, 'arguments', this.item.config.arguments || [])
      },
    },
  },

  created () {
    this.getTypes()
  },

  methods: {
    addArgument () {
      this.item.config.arguments.push({
        target: '',
        expr: '',
        type: 'Any',
        _showDetails: true,
      })
      this.$root.$emit('change-detected')
    },

    removeArgument (index) {
      this.item.config.arguments.splice(index, 1)
      this.$root.$emit('change-detected')
    },

    openInEditor (index = -1) {
      this.expressionEditor = {
        currentIndex: index >= -1 ? index : undefined,
        currentExpression: index >= 0 ? { ...this.item.config.arguments[index] } : undefined,
      }
    },

    saveExpression () {
      if (this.expressionEditor.currentIndex >= 0) {
        const args = [...this.item.config.arguments]
        args[this.expressionEditor.currentIndex] = this.expressionEditor.currentExpression
        this.$set(this.item.config, 'arguments', args)
        this.$root.$emit('change-detected')
      }

      this.resetExpression()
    },

    resetExpression () {
      this.expressionEditor = {
        currentIndex: undefined,
        currentExpression: undefined,
      }
    },

    async getTypes () {
      return this.$AutomationAPI.typeList()
        .then(({ set }) => {
          this.fieldTypes = set
        })
        .catch(this.toastErrorHandler(this.$t('notification:fetch-types-failed')))
    },

    getTypeDescription (type) {
      // This will be moved to backend field type information
      const typeDescriptions = {
        ID: 'Make sure to provide the ID in double quotes if you\'re using a literal value. Example "123"',
      }

      return typeDescriptions[type]
    },
  },
}
</script>
