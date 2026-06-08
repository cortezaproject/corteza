<template>
  <div>
    <div class="configurator-section">
      <div class="flex items-center justify-between mb-2">
        <label class="text-xs font-semibold uppercase tracking-wider text-muted-color">
          {{ $t('steps.expressions.label') }}
        </label>
        <Button
          :label="$t('steps.expressions.configurator.add-expression')"
          severity="secondary"
          size="small"
          @click="addArgument()"
        />
      </div>

      <div v-if="hasArguments">
        <expression-table
          value-field="expr"
          :items="item.config.arguments"
          :fields="argumentFields"
          :types="fieldTypes"
          @update:items="item.config.arguments = $event"
          @remove="removeArgument"
          @open-editor="openInEditor"
        />
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
            outlined
            size="small"
            @click="resetExpression"
          />
          <Button
            :label="$t('general.save')"
            size="small"
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
      fieldTypes: [],

      expressionEditor: {
        currentIndex: undefined,
        currentExpression: undefined,
      },
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
          label: this.$t('steps.expressions.configurator.target'),
          thClass: 'pl-4 ml-1',
          formatter: (item) => {
            return `${item.target}(${item.type})`
          },
        },
        {
          key: 'expr',
          label: this.$t('steps.expressions.configurator.expression'),
          thClass: 'pl-1 mr-3',
        },
      ]
    },

    hasArguments () {
      const { config } = this.item || {}
      return (config && (config.arguments || []).length) || []
    },
  },

  watch: {
    'item.config.stepID': {
      immediate: true,
      handler () {
        this.item.config['arguments'] = this.item.config.arguments || []
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
      eventBus.emit('change-detected')
    },

    removeArgument (index) {
      this.item.config.arguments.splice(index, 1)
      eventBus.emit('change-detected')
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
        this.item.config['arguments'] = args
        eventBus.emit('change-detected')
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
        .then(({ set = [] }) => {
          this.fieldTypes = set
        })
        .catch(e => this.toast.add({ severity: 'error', summary: this.$t('notification.fetch-types-failed'), detail: e?.message, life: 5000 }))
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
