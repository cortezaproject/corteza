<template>
  <div>
    <div class="configurator-section">
      <CFormGroup :label="$t('general.error-expression')">
        <expression-editor
          v-model="item.config.arguments[0].expr"
          font-size="18px"
          show-line-numbers
          @open="openInEditor"
          @input="valueChanged"
        />
      </CFormGroup>
    </div>

    <Dialog
      :visible="!!expressionEditor.currentExpression"
      :header="$t('editor.editor')"
      modal
      class="w-full max-w-5xl"
      @update:visible="resetExpression"
    >
      <expression-editor
        v-model="expressionEditor.currentExpression"
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
import ExpressionEditor from '../ExpressionEditor.vue'
import eventBus from '../../lib/eventBus'

export default {
  components: {
    ExpressionEditor,
  },

  extends: base,

  data () {
    return {
      expressionEditor: {
        currentExpression: undefined,
      },
    }
  },

  created () {
    let args = [{
      target: 'message',
      type: 'String',
      expr: '',
    }]

    if (this.item.config.arguments && this.item.config.arguments.length) {
      args = this.item.config.arguments.map(({ target, type, value, expr }) => {
        return {
          target,
          type,
          expr: expr || (value ? `"${value}"` : ''),
        }
      })
    }

    this.item.config['arguments'] = args
  },

  methods: {
    valueChanged (value) {
      this.$emit('update-default-value', {
        value: `Stop workflow with error: ${value}`,
        force: !this.item.node.value,
      })
      eventBus.emit('change-detected')
    },

    openInEditor () {
      this.expressionEditor.currentExpression = this.item.config.arguments[0].expr
    },

    saveExpression () {
      const { currentExpression } = this.expressionEditor
      this.item.config.arguments[0]['expr'] = currentExpression
      eventBus.emit('change-detected')

      this.resetExpression()
    },

    resetExpression () {
      this.expressionEditor.currentExpression = undefined
    },
  },
}
</script>
