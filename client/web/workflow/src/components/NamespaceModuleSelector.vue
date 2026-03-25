<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1">
      <label class="font-medium text-primary">
        {{ $t('general.filter.namespace.label') }}
      </label>
      <MultiSelect
        class="namespace-selector w-full"
        :options="namespace.options"
        :modelValue="namespace.values"
        :optionLabel="getNamespaceOptionLabel"
        :optionValue="n => `corteza::compose:namespace/${n.namespaceID}`"
        :placeholder="$t('general.filter.namespace.placeholder')"
        :loading="namespace.processing"
        filter
        @filter="searchNamespaces($event.value)"
        @update:modelValue="updateNamespaces"
      />
    </div>

    <div
      v-for="ns in namespace.values"
      :key="ns"
      class="flex flex-col gap-1"
    >
      <label class="font-medium text-primary">
        {{ getModuleLabel(ns) }}
      </label>
      <MultiSelect
        class="module-selector w-full"
        :options="getModulesForNamespace(ns.split('/')[1])"
        :modelValue="getModuleValuesForNamespace(ns.split('/')[1])"
        :optionLabel="getModuleOptionLabel"
        :optionValue="m => `corteza::compose:module/${m.namespaceID}/${m.moduleID}`"
        :placeholder="$t('general.filter.module.placeholder')"
        :loading="module.processing"
        filter
        @filter="e => searchModulesForNamespace(e.value, ns.split('/')[1])"
        @update:modelValue="modules => updateModulesForNamespace(modules, ns.split('/')[1])"
      />
    </div>
  </div>
</template>

<script>
import { debounce } from 'lodash-es'
export default {
  name: 'NamespaceModuleSelector',

  components: {},

  i18nOptions: {
    namespaces: 'general',
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
      this.namespace.values = this.namespaceLabels || []

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
      this.module.values = this.module.values.filter(label => {
        const nsID = label.split('/')[1]
        return nsID !== namespaceID
      })

      if (moduleLabels && moduleLabels.length > 0) {
        this.module.values = [...this.module.values, ...moduleLabels]
      }

      this.emitChange()
    },

    emitChange () {
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

      return this.$t('general.filter.module.template', { namespace: nsLabel })
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

    reset () {
      this.namespace.values = []
      this.module.values = []
      this.module.options = []
      this.emitChange()
    },
  },
}
</script>
