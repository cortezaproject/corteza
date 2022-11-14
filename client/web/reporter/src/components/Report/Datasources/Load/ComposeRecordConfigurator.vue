<template>
  <div>
    <b-row>
      <b-col>
        <b-form-group
          :label="$t('datasources:namespace')"
          class="text-primary"
        >
          <b-form-select
            v-model="namespace"
            :options="namespaces"
            text-field="name"
            value-field="namespaceID"
            @change="selectNamespace"
          >
            <template #first>
              <b-form-select-option
                :value="undefined"
              >
                {{ $t('general:label.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-form-group>
      </b-col>
      <b-col>
        <b-form-group
          v-if="namespace"
          :label="$t('datasources:module')"
          class="text-primary"
        >
          <b-form-select
            v-model="module"
            :options="modules"
            text-field="name"
            value-field="moduleID"
          >
            <template #first>
              <b-form-select-option
                :value="undefined"
              >
                {{ $t('general:label.none') }}
              </b-form-select-option>
            </template>
          </b-form-select>
        </b-form-group>
      </b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  props: {
    definition: {
      type: Object,
      required: true,
      default: () => ({}),
    },
  },

  data () {
    return {
      processing: false,

      namespaces: [],
      modules: [],
    }
  },

  computed: {
    namespace: {
      get () {
        return this.definition.namespaceID
      },

      set (namespaceID) {
        this.$emit('update:definition', { ...this.definition, namespaceID })
      },
    },

    module: {
      get () {
        return this.definition.moduleID
      },

      set (moduleID) {
        this.$emit('update:definition', { ...this.definition, moduleID })
      },
    },

    filter: {
      get () {
        return this.definition.filter
      },

      set (filter) {
        this.$emit('update:definition', { ...this.definition, filter })
      },
    },

    sort: {
      get () {
        return this.definition.sort
      },

      set (sort) {
        this.$emit('update:definition', { ...this.definition, sort })
      },
    },
  },

  watch: {
    namespace: {
      immediate: true,
      handler (namespace) {
        if (namespace) {
          this.processing = true

          this.fetchModules(namespace)
            .finally(() => {
              this.processing = false
            })
        }
      },
    },
  },

  created () {
    this.processing = true

    this.fetchNamespaces()
      .then(() => {
        if (this.namespace) {
          return this.fetchModules(this.namespace)
        }
      }).finally(() => {
        this.processing = false
      })
  },

  methods: {
    selectNamespace () {
      this.module = undefined
    },

    fetchNamespaces () {
      return this.$ComposeAPI.namespaceList({ sort: 'name' }).then(({ set = [] }) => {
        this.namespaces = set
      }).catch((e) => {
        this.toastErrorHandler(this.$t('notification:namespace.fetch-failed'))(e)
      })
    },

    fetchModules (namespaceID) {
      return this.$ComposeAPI.moduleList({ namespaceID, sort: 'name' }).then(({ set = [] }) => {
        this.modules = set
      }).catch((e) => {
        this.toastErrorHandler(this.$t('notification:module.fetch-failed'))(e)
      })
    },
  },
}
</script>
