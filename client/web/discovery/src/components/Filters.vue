<template>
  <div>
    <div class="my-2">
      <h6 class="text-primary mb-2">
        {{ this.$t('types.title') }}
      </h6>
      <b-form-checkbox-group
        v-model="types"
        name="types"
        :disabled="$store.state.processing"
        stacked
        class="mt-1"
        @change="updateTypes()"
      >
        <b-form-checkbox
          v-for="option in options"
          :key="option.value"
          :value="option.value"
          class="ml-2"
        >
          <div
            class="d-flex align-items-center mb-0"
          >
            <span class="d-inline-block text-truncate mr-3">
              {{ option.text }}
            </span>
          </div>
        </b-form-checkbox>
      </b-form-checkbox-group>
    </div>

    <div
      v-for="agg in aggregationOptions"
      :key="agg.resource"
      class="mt-4"
    >
      <div
        v-if="agg.items.length"
        class="d-flex justify-content-between align-items-center"
        style="min-height: 25px;"
      >
        <h6
          class="text-primary d-flex mb-0"
        >
          {{ agg.name }}
          <b-badge
            v-if="groups[agg.name].length"
            variant="dark"
            pill
            class="ml-1 align-self-center"
          >
            {{ groups[agg.name].length }}
          </b-badge>
        </h6>
        <b-button
          v-if="groups[agg.name].length"
          variant="link"
          class="text-muted p-0 m-0"
          size="sm"
          @click="clearGroup(agg.name)"
        >
          {{ $t('reset') }}
        </b-button>
      </div>

      <b-form-checkbox-group
        v-model="groups[agg.name]"
        stacked
        class="mt-1 ml-2"
        :disabled="$store.state.processing"
        @change="updateGroup(agg.name)"
      >
        <b-form-checkbox
          v-for="(resource, i) in agg.items"
          :key="i"
          :value="resource.name"
          class="mb-1"
        >
          <div
            class="d-flex align-items-center"
          >
            <span class="d-inline-block text-truncate">
              {{ getResourceDisplayName(agg.resource, resource) }}
            </span>
            <span
              class="pl-3 ml-auto text-muted"
            >
              {{ resource.hits }}
            </span>
          </div>
        </b-form-checkbox>
      </b-form-checkbox-group>
    </div>
  </div>
</template>

<script>
export default {
  i18nOptions: {
    namespaces: 'filters',
  },

  data () {
    return {
      types: [],
      groups: {
        Module: [],
        Namespace: [],
      },
    }
  },

  computed: {
    options () {
      return [
        { text: this.$t('types.namespace'), value: 'compose:namespace' },
        { text: this.$t('types.module'), value: 'compose:module' },
        { text: this.$t('types.record'), value: 'compose:record' },
        { text: this.$t('types.user'), value: 'system:user' },
      ]
    },

    aggregationOptions () {
      let namespaceOptions = this.$store.state.aggregations.find(({ resource }) => resource === 'compose:namespace') || {}
      let moduleOptions = this.$store.state.aggregations.find(({ resource }) => resource === 'compose:module') || {}

      namespaceOptions = {
        resource: 'compose:namespace',
        name: 'Namespace',
        hits: namespaceOptions.hits || 0,
        items: namespaceOptions.resource_name || [],
      }

      // Get all modules that are missing from store aggregations but are in filter
      const missingModuleOptions = this.groups.Module.filter(name => !(moduleOptions.resource_name || []).some(o => o.name === name))
        .map(name => ({ name }))

      moduleOptions = {
        resource: 'compose:module',
        name: 'Module',
        hits: moduleOptions.hits || 0,
        items: [
          ...missingModuleOptions,
          ...(moduleOptions.resource_name || []),
        ],
      }

      return [namespaceOptions, moduleOptions]
    },
  },

  watch: {
    '$store.state.namespaces': {
      immediate: true,
      handler (namespace) {
        this.groups.Namespace = namespace
      },
    },

    '$store.state.modules': {
      immediate: true,
      handler (module) {
        this.groups.Module = module
      },
    },
  },
  methods: {
    getResourceDisplayName (type, { name, handle, slug }) {
      if (type === 'compose:namespace') {
        return name || slug || 'Unnamed namespace'
      } else if (type === 'compose:module') {
        return handle || name || 'Unnamed module'
      }
    },

    updateTypes () {
      this.$store.commit('updateTypes', this.types)
    },

    updateGroup (name) {
      this.$store.commit(`update${name}s`, this.groups[name])
    },

    clearGroup (name) {
      this.groups[name] = []
      this.$store.commit(`update${name}s`, [])
    },
  },
}
</script>
