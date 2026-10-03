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
      <b-form-group
        :label="$t('general:error-expression')"
        label-class="text-primary"
      >
        <expression-editor
          v-model="item.config.arguments[0].expr"
          font-size="18px"
          show-line-numbers
          @open="openInEditor(0)"
          @input="valueChanged"
        />
      </b-form-group>

      <b-form-group
        :label="$t('general:error-title-expression')"
        :description="$t('general:error-title-description')"
        label-class="text-primary"
        class="mb-0"
      >
        <expression-editor
          v-model="item.config.arguments[1].expr"
          font-size="18px"
          show-line-numbers
          @open="openInEditor(1)"
          @input="$root.$emit('change-detected')"
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

export default {
  components: {
    ExpressionEditor,
  },

  extends: base,

  data () {
    return {
      expressionEditor: {
        currentExpression: undefined,
        index: 0,
      },
    }
  },

  created () {
    // message is what the error says; the optional title is shown as the
    // heading of the notification instead of the generic one
    const defaults = [
      { target: 'message', type: 'String', expr: '' },
      { target: 'title', type: 'String', expr: '' },
    ]

    const configured = (this.item.config.arguments || []).map(({ target, type, value, expr }) => {
      return {
        target,
        type,
        expr: expr || (value ? `"${value}"` : ''),
      }
    })

    const args = defaults.map(d => configured.find(a => a.target === d.target) || d)

    this.$set(this.item.config, 'arguments', args)
  },

  methods: {
    valueChanged (value) {
      this.$emit('update-default-value', {
        value: `Stop workflow with error: ${value}`,
        force: !this.item.node.value,
      })
      this.$root.$emit('change-detected')
    },

    openInEditor (index = 0) {
      this.expressionEditor.index = index
      this.expressionEditor.currentExpression = this.item.config.arguments[index].expr
    },

    saveExpression () {
      const { currentExpression, index } = this.expressionEditor
      this.$set(this.item.config.arguments[index], 'expr', currentExpression)
      this.$root.$emit('change-detected')

      this.resetExpression()
    },

    resetExpression () {
      this.expressionEditor.currentExpression = undefined
    },
  },
}
</script>
