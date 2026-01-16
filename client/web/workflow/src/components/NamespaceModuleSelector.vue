<template>
  <div>
    <b-form-group
      :label="$t('filter.namespace.label')"
      label-class="text-primary"
    >
      <c-input-select
        class="namespace-selector"
        :options="namespace.options"
        :value="namespace.values"
        :get-option-label="getNamespaceOptionLabel"
        :get-option-key="n => `corteza::compose:namespace/${n.namespaceID}`"
        :reduce="n => `corteza::compose:namespace/${n.namespaceID}`"
        :placeholder="$t('filter.namespace.placeholder')"
        :loading="namespace.processing"
        :filterable="false"
        multiple
        @search="searchNamespaces"
        @input="updateNamespaces"
      />
    </b-form-group>

    <b-form-group
      v-for="ns in namespace.values"
      :key="ns"
      :label="getModuleLabel(ns)"
      label-class="text-primary"
    >
      <c-input-select
        class="module-selector"
        :options="getModulesForNamespace(ns.split('/')[1])"
        :value="getModuleValuesForNamespace(ns.split('/')[1])"
        :get-option-key="m => `corteza::compose:module/${m.namespaceID}/${m.moduleID}`"
        :get-option-label="getModuleOptionLabel"
        :reduce="m => `corteza::compose:module/${m.namespaceID}/${m.moduleID}`"
        :placeholder="$t('filter.module.placeholder')"
        :loading="module.processing"
        :filterable="false"
        multiple
        @search="query => searchModulesForNamespace(query, ns.split('/')[1])"
        @input="modules => updateModulesForNamespace(modules, ns.split('/')[1])"
      />
    </b-form-group>
  </div>
</template>

<script>
import { debounce } from 'lodash'
import { components } from '@cortezaproject/corteza-vue'

const { CInputSelect } = components

export default {
  name: 'NamespaceModuleSelector',

  components: {
    CInputSelect,
  },

  i18nOptions: {
    namespaces: 'general',
  },

  props: {
    /**
     * Initial namespace labels in format: ["corteza::compose:namespace/id1", ...]
     */
    namespaceLabels: {
      type: Array,
      default: () => [],
    },

    /**
     * Initial module labels in format: ["corteza::compose:module/nsID/modID", ...]
     */
    moduleLabels: {
      type: Array,
      default: () => [],
    },
  },

  data () {
    return {
      namespace: {
        processing: false,
        values: [],
        options: [],
        filter: {
          query: null,
          limit: 20,
          sort: 'name DESC',
        },
      },

      module: {
        processing: false,
        values: [],
        options: [],
        filter: {
          query: null,
          limit: 20,
          sort: 'name DESC',
        },
      },
    }
  },

  computed: {
    modulesByNamespace () {
      const grouped = {}

      this.module.options.forEach(module => {
        if (!grouped[module.namespaceID]) {
          grouped[module.namespaceID] = []
        }
        grouped[module.namespaceID].push(module)
      })

      return grouped
    },
  },

  watch: {
    namespaceLabels: {
      handler (newVal) {
        if (newVal && newVal.length > 0 && this.namespace.options.length > 0) {
          this.initializeFromProps()
        }
      },
      immediate: false,
    },
  },

  created () {
    this.fetchNamespaces().then(() => {
      this.initializeFromProps()
    })
  },

  methods: {
    initializeFromProps () {
      // Set namespace values directly from labels
      this.namespace.values = this.namespaceLabels || []

      // Set module values directly from labels
      if (this.namespace.values.length > 0) {
        this.fetchModules().then(() => {
          this.module.values = this.moduleLabels || []
        })
      }
    },

    fetchNamespaces () {
      this.namespace.processing = true

      return this.$ComposeAPI.namespaceList(this.namespace.filter).then(({ set = [] } = {}) => {
        const namespacePromises = []

        // Fetch any missing namespaces from initial labels
        if (this.namespaceLabels && this.namespaceLabels.length > 0 && !this.namespace.filter.query) {
          const namespaceIDs = this.namespaceLabels.map(label => label.split('/')[1]).filter(Boolean)

          namespaceIDs.forEach(namespaceID => {
            if (!set.some(n => n.namespaceID === namespaceID)) {
              namespacePromises.push(
                this.$ComposeAPI.namespaceRead({ namespaceID })
                  .then(n => [n])
                  .catch(() => []),
              )
            }
          })
        }

        return Promise.all(namespacePromises).then(results => {
          this.namespace.options = [...set, ...results.flat()].sort((a, b) =>
            (a.name || '').localeCompare(b.name || ''),
          )
        }).catch(() => {
          this.namespace.options = []
        })
      }).finally(() => {
        this.namespace.processing = false
      })
    },

    fetchModules () {
      if (!this.namespace.values || this.namespace.values.length === 0) {
        this.module.options = []
        return Promise.resolve()
      }

      this.module.processing = true

      // Extract namespace IDs from label strings
      const namespaceIDs = this.namespace.values.map(label => label.split('/')[1])

      const promises = namespaceIDs.map(namespaceID =>
        this.$ComposeAPI.moduleList({
          namespaceID,
          ...this.module.filter,
        }).then(({ set }) => set),
      )

      return Promise.all(promises).then(results => {
        this.module.options = results.flat()
      }).catch(() => {
        this.module.options = []
      }).finally(() => {
        this.module.processing = false
      })
    },

    searchNamespaces: debounce(function (query) {
      if (query !== this.namespace.filter.query) {
        this.namespace.filter.query = query
      }
      this.fetchNamespaces()
    }, 300),

    searchModulesForNamespace: debounce(function (query, namespaceID) {
      if (query !== this.module.filter.query) {
        this.module.filter.query = query
      }
      this.fetchModules()
    }, 300),

    updateNamespaces (namespaceLabels) {
      this.namespace.values = namespaceLabels || []

      if (this.namespace.values.length > 0) {
        // Filter out modules from namespaces that are no longer selected
        const selectedNsIDs = new Set(this.namespace.values.map(label => label.split('/')[1]))
        this.module.values = this.module.values.filter(moduleLabel => {
          const nsID = moduleLabel.split('/')[1]
          return selectedNsIDs.has(nsID)
        })

        this.fetchModules()
      } else {
        this.module.options = []
        this.module.values = []
      }

      this.emitChange()
    },

    updateModulesForNamespace (moduleLabels, namespaceID) {
      // Remove old modules for this namespace
      this.module.values = this.module.values.filter(label => {
        const nsID = label.split('/')[1]
        return nsID !== namespaceID
      })

      // Add new module labels for this namespace
      if (moduleLabels && moduleLabels.length > 0) {
        this.module.values = [...this.module.values, ...moduleLabels]
      }

      this.emitChange()
    },

    emitChange () {
      // Values are already in label format - emit directly
      this.$emit('change', {
        namespaceLabels: this.namespace.values,
        moduleLabels: this.module.values,
      })
    },

    getNamespaceOptionLabel ({ name, handle } = {}) {
      return name || handle || 'Unnamed Namespace'
    },

    getModuleOptionLabel (module) {
      return module.name || module.handle || 'Unnamed Module'
    },

    getModuleLabel (namespaceLabel) {
      const namespaceID = namespaceLabel.split('/')[1]
      const namespace = this.namespace.options.find(n => n.namespaceID === namespaceID)
      const nsLabel = namespace ? this.getNamespaceOptionLabel(namespace) : namespaceID

      return this.$t('filter.module.template', { namespace: nsLabel })
    },

    getModulesForNamespace (namespaceID) {
      return this.modulesByNamespace[namespaceID] || []
    },

    getModuleValuesForNamespace (namespaceID) {
      return this.module.values.filter(label => {
        const nsID = label.split('/')[1]
        return nsID === namespaceID
      })
    },

    /**
     * Reset the selector to empty state
     */
    reset () {
      this.namespace.values = []
      this.module.values = []
      this.module.options = []
      this.emitChange()
    },
  },
}
</script>

<style scoped>
/* Style namespace select tags (primary blue) */
::v-deep .namespace-selector .vs__selected {
  background-color: var(--primary) !important;
  color: white !important;
}

/* Style module select tags (extra-light with dark text) */
::v-deep .module-selector .vs__selected {
  background-color: var(--extra-light) !important;
  color: var(--dark) !important;
}
</style>
