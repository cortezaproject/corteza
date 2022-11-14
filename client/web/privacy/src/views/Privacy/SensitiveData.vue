<template>
  <b-container
    fluid
    class="d-flex flex-column p-3"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <c-resource-list
      primary-key="moduleID"
      :items="items"
      :fields="fields"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :selectable="false"
      :translations="{
        noItems: $t('general:resourceList.noItems'),
        loading: $t('general:resourceList.loading'),
        searchPlaceholder: $t('general:resourceList.search.placeholder'),
        showingPagination: 'general:resourceList.pagination.showing',
        singlePluralPagination: 'general:resourceList.pagination.single',
        prevPagination: $t('general:resourceList.pagination.prev'),
        nextPagination: $t('general:resourceList.pagination.next'),
      }"
      clickable
      hide-total
      class="flex-grow-1"
      @search="filterList"
      @row-clicked="rowClicked"
    />
  </b-container>
</template>

<script>
import listHelpers from 'corteza-webapp-privacy/src/mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  name: 'SensitiveData',

  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'sensitive-data',
    keyPrefix: 'list',
  },

  data () {
    return {
      fields: [
        {
          key: 'module',
          formatter: module => {
            return module ? module.name : ''
          },
        },
        {
          key: 'namespace',
          formatter: namespace => {
            return namespace ? namespace.name : ''
          },
        },
        {
          key: 'connection',
          formatter: connection => {
            const { name } = connection.meta || {}
            return name
          },
        },
        {
          key: 'location',
          formatter: (value, key, item) => {
            const { meta = {} } = item.connection || {}
            const { properties = {} } = meta.location || {}
            return properties.name
          },
        },
        {
          key: 'ownership',
          formatter: (value, key, item) => {
            const { meta = {} } = item.connection || {}
            const { ownership } = meta || {}
            return ownership
          },
        },
      ].map(c => ({
        ...c,
        // Generate column label translation key
        label: c.label || this.$t(`columns.${c.key}`),
      })),
    }
  },

  methods: {
    items () {
      return this.procListResults(this.$ComposeAPI.dataPrivacyModuleList(this.encodeListParams()))
    },

    rowClicked ({ namespace, module }) {
      window.open(`/compose/ns/${namespace.slug}/admin/modules/${module.moduleID}/edit`, '_blank')
    },
  },
}
</script>
