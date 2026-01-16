<template>
  <b-modal
    id="workflow-filter"
    :title="$t('general:filter.title')"
    size="lg"
    no-fade
    body-class="p-3"
    @show="onShow"
  >
    <namespace-module-selector
      ref="selector"
      :namespace-labels="localNamespaceLabels"
      :module-labels="localModuleLabels"
      @change="handleChange"
    />

    <template #modal-footer>
      <div class="d-flex gap-1 w-100">
        <b-button
          variant="light"
          @click="handleReset"
        >
          {{ $t('general:reset') }}
        </b-button>

        <div class="d-flex ml-auto gap-1">
          <b-button
            variant="light"
            @click="handleCancel"
          >
            {{ $t('general:cancel') }}
          </b-button>

          <b-button
            variant="primary"
            @click="handleApply"
          >
            {{ $t('general:filter.apply') }}
          </b-button>
        </div>
      </div>
    </template>
  </b-modal>
</template>

<script>
import NamespaceModuleSelector from './NamespaceModuleSelector'

export default {
  name: 'WorkflowFilterModal',

  components: {
    NamespaceModuleSelector,
  },

  props: {
    /**
     * Current namespace labels filter
     */
    namespaceLabels: {
      type: Array,
      default: () => [],
    },

    /**
     * Current module labels filter
     */
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

  methods: {
    onShow () {
      // Reset local state to current props when modal opens
      this.localNamespaceLabels = [...this.namespaceLabels]
      this.localModuleLabels = [...this.moduleLabels]
    },

    handleChange ({ namespaceLabels, moduleLabels }) {
      this.localNamespaceLabels = namespaceLabels
      this.localModuleLabels = moduleLabels
    },

    handleApply () {
      this.$emit('apply', {
        namespaceLabels: this.localNamespaceLabels,
        moduleLabels: this.localModuleLabels,
      })
      this.$bvModal.hide('workflow-filter')
    },

    handleReset () {
      this.localNamespaceLabels = []
      this.localModuleLabels = []
      if (this.$refs.selector) {
        this.$refs.selector.reset()
      }
    },

    handleCancel () {
      this.$bvModal.hide('workflow-filter')
    },
  },
}
</script>
