<template>
  <div>
    <div class="flex flex-col gap-1">
      <label class="font-medium text-primary">
        {{ $t('general.label.label') }}
      </label>
      <InputText
        v-model="label"
        @input="$emit('update-value', $event)"
      />
    </div>
  </div>
</template>

<script>
import base from './base.vue'

export default {
  extends: base,

  computed: {
    // Ignores exclusiveGateway indexes (#n)
    label: {
      get () {
        if (this.getSourceType) {
          if (this.getSourceType === 'gatewayExclusive') {
            /* eslint-disable no-unused-vars */
            const [edgeID, ...rest] = this.item.node.value.split(' - ')
            return rest.join(' - ')
          }
        }

        return this.item.node.value
      },

      set (label) {
        if (this.getSourceType) {
          if (this.getSourceType === 'gatewayExclusive') {
            /* eslint-disable no-unused-vars */
            const [edgeID, ...rest] = this.item.node.value.split(' - ')
            const newLabel = [edgeID]
            if (label) {
              newLabel.push(label)
            }
            label = newLabel.join(' - ')
          }
        }

        this.item.node.value = label
      },
    },

    getSourceType () {
      const { source } = this.item.node
      if (source && source.style) {
        return source.style
      }
      return undefined
    },
  },
}
</script>
