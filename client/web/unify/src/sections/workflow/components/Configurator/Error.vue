<template>
  <div>
    <div class="configurator-section">
      <CFormGroup :label="$t('general.error-expression')">
        <expression-editor
          v-model="item.config.arguments[0].expr"
          font-size="18px"
          show-line-numbers
          @open="openInEditor(0)"
          @input="valueChanged"
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('general.error-title-expression')"
        :description="$t('general.error-title-description')"
      >
        <expression-editor
          v-model="item.config.arguments[1].expr"
          font-size="18px"
          show-line-numbers
          @open="openInEditor(1)"
          @input="eventBus.emit('change-detected')"
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
            size="small"
            @click="resetExpression"
          />
          <Button :label="$t('general.save')" size="small" @click="saveExpression" />
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

  data() {
    return {
      eventBus,
      expressionEditor: {
        currentExpression: undefined,
        index: 0,
      },
    }
  },

  created() {
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

    this.item.config['arguments'] = defaults.map(
      d => configured.find(a => a.target === d.target) || d,
    )
  },

  methods: {
    valueChanged(value) {
      this.$emit('update-default-value', {
        value: `Stop workflow with error: ${value}`,
        force: !this.item.node.value,
      })
      eventBus.emit('change-detected')
    },

    openInEditor(index = 0) {
      this.expressionEditor.index = index
      this.expressionEditor.currentExpression = this.item.config.arguments[index].expr
    },

    saveExpression() {
      const { currentExpression, index } = this.expressionEditor
      this.item.config.arguments[index]['expr'] = currentExpression
      eventBus.emit('change-detected')

      this.resetExpression()
    },

    resetExpression() {
      this.expressionEditor.currentExpression = undefined
    },
  },
}
</script>
