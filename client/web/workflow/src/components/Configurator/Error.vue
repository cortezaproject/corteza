<template>
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
      <!-- Title -->
      <b-form-group
        :label="$t('general:error-step.title.label')"
        :description="$t('general:error-step.title.description')"
        label-class="text-primary"
      >
        <expression-editor
          v-model="titleArg.expr"
          font-size="16px"
          show-line-numbers
          @open="openInEditor('title')"
          @input="onFieldInput"
        />
      </b-form-group>

      <!-- Severity -->
      <b-form-group
        :label="$t('general:error-step.severity.label')"
        :description="$t('general:error-step.severity.description')"
        label-class="text-primary"
      >
        <b-form-select
          v-if="severityIsLiteral"
          v-model="severityValue"
          :options="severityOptions"
          @change="onSeverityChange"
        />
        <div
          v-else
          class="d-flex align-items-center"
        >
          <b-form-input
            v-model="severityArg.expr"
            class="flex-grow-1"
            :placeholder="$t('general:error-step.severity.label')"
            @input="onFieldInput"
          />
          <b-button
            v-b-tooltip.hover
            variant="link"
            size="sm"
            :title="$t('general:error-step.severity.reset-to-literal')"
            class="ml-2 p-0"
            @click="resetSeverityToLiteral"
          >
            <font-awesome-icon :icon="['fas', 'undo']" />
          </b-button>
        </div>
      </b-form-group>

      <!-- Message -->
      <b-form-group
        :label="$t('general:error-step.message.label')"
        :description="$t('general:error-step.message.description')"
        label-class="text-primary"
        class="mb-0"
      >
        <expression-editor
          v-model="messageArg.expr"
          font-size="16px"
          show-line-numbers
          @open="openInEditor('message')"
          @input="onFieldInput"
        />
      </b-form-group>
    </b-card-body>

    <b-modal
      id="expression-editor"
      :visible="!!expressionEditor.currentExpression"
      :title="$t('editor:editor')"
      size="xl"
      scrollable
      :ok-title="$t('general:save')"
      :cancel-title="$t('general:cancel')"
      cancel-variant="light"
      body-class="p-0"
      no-fade
      @ok="saveExpression"
      @hidden="resetExpression"
    >
      <expression-editor
        v-model="expressionEditor.currentExpression"
        min-height="80vh"
        font-size="18px"
        show-line-numbers
        :border="false"
        :show-popout="false"
      />
    </b-modal>
  </b-card>
</template>

<script>
import base from './base'
import ExpressionEditor from '../ExpressionEditor.vue'

const ALLOWED_TARGETS = ['message', 'title', 'severity']
const SEVERITY_VALUES = ['error', 'warning', 'info']

function makeArg (target, expr = '') {
  return { target, type: 'String', expr }
}

export default {
  components: {
    ExpressionEditor,
  },

  extends: base,

  data () {
    return {
      expressionEditor: {
        currentExpression: undefined,
        currentField: undefined,
      },
    }
  },

  computed: {
    messageArg () {
      return this.findOrCreateArg('message')
    },
    titleArg () {
      return this.findOrCreateArg('title')
    },
    severityArg () {
      return this.findOrCreateArg('severity')
    },

    // severityIsLiteral: true when the severity expression is either
    // empty or a plain quoted literal we can round-trip through the
    // dropdown. When it's a real expression (e.g. `vars.level`), the
    // configurator falls back to a raw text input so we never clobber
    // the author's work.
    severityIsLiteral () {
      const raw = (this.severityArg.expr || '').trim()
      if (!raw) return true
      return /^["'](error|warning|info)["']$/.test(raw)
    },

    severityValue: {
      get () {
        const raw = (this.severityArg.expr || '').trim()
        if (!raw) return 'error'
        const m = raw.match(/^["'](error|warning|info)["']$/)
        return m ? m[1] : 'error'
      },
      set (v) {
        if (!SEVERITY_VALUES.includes(v)) v = 'error'
        const arg = this.severityArg
        this.$set(arg, 'expr', `"${v}"`)
      },
    },

    severityOptions () {
      return [
        { value: 'error', text: this.$t('general:error-step.severity.options.error') },
        { value: 'warning', text: this.$t('general:error-step.severity.options.warning') },
        { value: 'info', text: this.$t('general:error-step.severity.options.info') },
      ]
    },
  },

  created () {
    // Normalise config.arguments: keep only allowed targets, ensure all of
    // them exist (even if empty) so v-model bindings are stable.
    const existing = Array.isArray(this.item.config.arguments) ? this.item.config.arguments : []
    const byTarget = {}
    existing.forEach(({ target, type, value, expr }) => {
      if (!ALLOWED_TARGETS.includes(target)) return
      byTarget[target] = {
        target,
        type: type || 'String',
        expr: expr || (value ? `"${value}"` : ''),
      }
    })

    const args = ALLOWED_TARGETS.map(t => byTarget[t] || makeArg(t))
    this.$set(this.item.config, 'arguments', args)
  },

  methods: {
    findOrCreateArg (target) {
      const args = this.item.config.arguments || []
      let arg = args.find(a => a.target === target)
      if (!arg) {
        arg = makeArg(target)
        args.push(arg)
        this.$set(this.item.config, 'arguments', args)
      }
      return arg
    },

    // stripLiteralQuotes returns the unquoted content of a simple
    // quoted string expression (e.g. "foo" -> foo) so the canvas node
    // label preview reads cleanly. Non-literal expressions are passed
    // through unchanged.
    stripLiteralQuotes (s) {
      const t = (s || '').trim()
      const m = t.match(/^(["'])((?:[^\\]|\\.)*)\1$/)
      return m ? m[2] : t
    },

    onFieldInput () {
      const pick = this.titleArg.expr || this.messageArg.expr || ''
      const preview = this.stripLiteralQuotes(pick)
      this.$emit('update-default-value', {
        value: preview ? `Stop workflow with error: ${preview}` : 'Stop workflow with error',
        force: !this.item.node.value,
      })
      this.$root.$emit('change-detected')
    },

    onSeverityChange () {
      this.$root.$emit('change-detected')
    },

    resetSeverityToLiteral () {
      this.$set(this.severityArg, 'expr', '"error"')
      this.$root.$emit('change-detected')
    },

    openInEditor (field) {
      const arg = this.findOrCreateArg(field)
      this.expressionEditor.currentField = field
      this.expressionEditor.currentExpression = arg.expr || ''
    },

    saveExpression () {
      const { currentExpression, currentField } = this.expressionEditor
      if (currentField) {
        const arg = this.findOrCreateArg(currentField)
        this.$set(arg, 'expr', currentExpression)
        this.$root.$emit('change-detected')
      }
      this.resetExpression()
    },

    resetExpression () {
      this.expressionEditor.currentExpression = undefined
      this.expressionEditor.currentField = undefined
    },
  },
}
</script>
