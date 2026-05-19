<template>
  <div>
    <div class="configurator-section">
      <CFormGroup :label="$t('configurator.delay.duration.label')">
        <InputText
          v-model="item.config.arguments[0].expr"
          :placeholder="$t('configurator.delay.duration.placeholder')"
          @update:model-value="valueChanged"
        />
      </CFormGroup>
    </div>
  </div>
</template>

<script>
import InputText from 'primevue/inputtext'
import base from './base.vue'
import eventBus from '../../lib/eventBus'

export default {
  components: {
    InputText,
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
