<template>
  <div>
    <div class="configurator-section">
      <div class="configurator-section__title">{{ $t('configurator.configuration') }}</div>
      <div class="flex flex-col gap-1">
        <label class="font-medium text-primary">
          {{ $t('configurator.delay.duration.label') }}
        </label>
        <expression-editor
          v-model="item.config.arguments[0].expr"
          font-size="18px"
          show-line-numbers
          :show-popout="false"
          @input="valueChanged"
        />
        <small class="text-muted-color">
          {{ $t('configurator.delay.duration.description') }}
        </small>
      </div>
    </div>
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

  watch: {
    'item.config.stepID': {
      immediate: true,
      handler () {
        let args = [{
          target: 'offset',
          type: 'Duration',
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
    },
  },

  methods: {
    valueChanged (value) {
      this.$emit('update-default-value', {
        value: `Delay workflow execution for ${value}`,
        force: !this.item.node.value,
      })

      eventBus.emit('change-detected')
    },
  },
}
</script>
