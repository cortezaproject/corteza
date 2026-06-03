<template>
  <div class="flex flex-col gap-4">
    <namespace-module-selector
      ref="selector"
      :namespace-labels="localNamespaceLabels"
      :module-labels="localModuleLabels"
      @change="handleChange"
    />

    <div class="flex gap-2 justify-between">
      <Button
        :label="$t('general.reset')"
        severity="secondary"
        @click="handleReset"
      />

      <div class="flex gap-2">
        <Button
          :label="$t('general.cancel')"
          severity="secondary"
          @click="$emit('close')"
        />

        <Button
          :label="$t('general.filter.apply')"
          @click="handleApply"
        />
      </div>
    </div>
  </div>
</template>

<script>
import NamespaceModuleSelector from './NamespaceModuleSelector.vue'

export default {
  name: 'WorkflowFilterModal',

  components: {
    NamespaceModuleSelector,
  },

  props: {
    namespaceLabels: {
      type: Array,
      default: () => [],
    },

    moduleLabels: {
      type: Array,
      default: () => [],
    },
  },

  data () {
    return {
      localNamespaceLabels: [],
      localModuleLabels: [],
    }
  },

  watch: {
    namespaceLabels: {
      immediate: true,
      handler (val) {
        this.localNamespaceLabels = [...(val || [])]
      },
    },
    moduleLabels: {
      immediate: true,
      handler (val) {
        this.localModuleLabels = [...(val || [])]
      },
    },
  },

  methods: {
    handleChange ({ namespaceLabels, moduleLabels }) {
      this.localNamespaceLabels = namespaceLabels
      this.localModuleLabels = moduleLabels
    },

    handleApply () {
      this.$emit('apply', {
        namespaceLabels: this.localNamespaceLabels,
        moduleLabels: this.localModuleLabels,
      })
      this.$emit('close')
    },

    handleReset () {
      this.localNamespaceLabels = []
      this.localModuleLabels = []
      if (this.$refs.selector) {
        this.$refs.selector.reset()
      }
    },
  },
}
</script>
