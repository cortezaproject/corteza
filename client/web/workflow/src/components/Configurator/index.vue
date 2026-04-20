<template>
  <div class="flex flex-col">
    <component
      :is="stepComponent"
      v-if="stepComponent"
      v-model:item="item"
      v-model:edges="edges"
      :out-edges="outEdges"
      :is-subworkflow="isSubworkflow"
      @update-value="$emit('update-value', $event)"
      @update-default-value="updateDefaultName"
    />
  </div>
</template>
<script>
import base from './base.vue'
import * as Configurators from './loader'

export default {
  components: {
    ...Configurators,
  },

  extends: base,

  data () {
    return {
      collapse: {
        basic: true,
        configurator: true,
      },
    }
  },

  computed: {
    stepComponent () {
      return Configurators[this.kind]
    },

    kind () {
      const { kind, ref } = this.item.config

      if (kind === 'exec-workflow') {
        return 'ExecWorkflow'
      }

      if (kind === 'error-handler') {
        return 'ErrorHandler'
      }

      if (kind === 'visual' && ref === 'content') {
        return 'Content'
      }

      if (kind) {
        return kind.charAt(0).toUpperCase() + kind.slice(1)
      }

      return undefined
    },
  },

  methods: {
    updateDefaultName ({ value, force = false }) {
      if (force || this.item.config.defaultName || this.item.config.defaultName === undefined) {
        this.$emit('update-default-value', value)
      }
    },
  },
}
</script>
