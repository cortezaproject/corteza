"<template>
  <div class="py-3">
    <portal to="topbar-title">
      {{ $t('navigation.module') }}
    </portal>

    <b-container fluid="xl">
      <b-row>
        <b-col>
          <c-resource-list
            data-test-id="table-modules-list"
            :primary-key="primaryKey"
            :filter="filter"
            :sorting="sorting"
            :pagination="pagination"
            :fields="fields"
            :items="items"
            :translations="{
              searchPlaceholder: $t('searchPlaceholder'),
              notFound: $t('general:resourceList.notFound'),
              noItems: $t('general:resourceList.noItems'),
              loading: $t('general:label.loading'),
              showingPagination: 'general:resourceList.pagination.showing',
              singlePluralPagination: 'general:resourceList.pagination.single',
              prevPagination: $t('general:resourceList.pagination.prev'),
              nextPagination: $t('general:resourceList.pagination.next'),
            }"
            clickable
            @search="filterList"
            @row-clicked="handleRowClicked"
          >
            <template #header>
              <div class="flex-grow-1">
                <div
                  class="wrap-with-vertical-gutters"
                >
                  <b-btn
                    v-if="namespace.canCreateModule"
                    data-test-id="button-create"
                    variant="primary"
                    size="lg"
                    class="mr-1 float-left"
                    :to="{ name: 'admin.modules.create' }"
                  >
                    {{ $t('createLabel') }}
                  </b-btn>

                  <import
                    v-if="namespace.canCreateModule"
                    :namespace="namespace"
                    type="module"
                    class="mr-1 float-left"
                    @importSuccessful="onImportSuccessful"
                  />

                  <export
                    :list="modules"
                    type="module"
                    class="mr-1 float-left"
                  />

                  <b-dropdown
                    v-if="namespace.canGrant"
                    size="lg"
                    variant="light"
                    class="permissions-dropdown mr-1"
                  >
                    <template #button-content>
                      <font-awesome-icon :icon="['fas', 'lock']" />
                      <span>
                        {{ $t('general:label.permissions') }}
                      </span>
                    </template>

                    <b-dropdown-item>
                      <c-permissions-button
                        :resource="`corteza::compose:module/${namespace.namespaceID}/*`"
                        :button-label="$t('general:label.module')"
                        :show-button-icon="false"
                        button-variant="white text-left w-100"
                      />
                    </b-dropdown-item>

                    <b-dropdown-item>
                      <c-permissions-button
                        :resource="`corteza::compose:module-field/${namespace.namespaceID}/*/*`"
                        :button-label="$t('general:label.field')"
                        :show-button-icon="false"
                        button-variant="white text-left w-100"
                      />
                    </b-dropdown-item>

                    <b-dropdown-item>
                      <c-permissions-button
                        :resource="`corteza::compose:record/${namespace.namespaceID}/*/*`"
                        :button-label="$t('general:label.record')"
                        :show-button-icon="false"
                        button-variant="white text-left w-100"
                      />
                    </b-dropdown-item>
                  </b-dropdown>
                </div>
              </div>
            </template>

            <template #actions="{ item: m }">
              <related-pages
                :namespace="namespace"
                :module="m"
              />
              <b-button
                data-test-id="button-all-records"
                variant="link"
                :to="{name: 'admin.modules.record.list', params: { moduleID: m.moduleID }}"
                class="text-dark text-decoration-none"
              >
                {{ $t('allRecords.label') }}
              </b-button>
              <c-permissions-button
                v-if="m.canGrant"
                :title="m.name || m.handle || m.moduleID"
                :target="m.name || m.handle || m.moduleID"
                :resource="`corteza::compose:module/${m.namespaceID}/${m.moduleID}`"
                :tooltip="$t('permissions:resources.compose.module.tooltip')"
                class="btn px-2"
                link
              />
            </template>

            <template #name="{ item: m }">
              <div
                class="d-flex align-items-center"
              >
                {{ m.name }}
                <h5
                  class="ml-2 mb-0"
                >
                  <b-badge
                    v-if="Object.keys(m.labels || {}).includes('federation')"
                    pill
                    variant="primary"
                  >
                    {{ $t('federated') }}
                  </b-badge>
                </h5>
              </div>
            </template>

            <template #changedAt="{ item }">
              {{ (item.deletedAt || item.updatedAt || item.createdAt) | locFullDateTime }}
            </template>
          </c-resource-list>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import { compose } from '@cortezaproject/corteza-js'
import listHelpers from 'corteza-webapp-compose/src/mixins/listHelpers'
import RelatedPages from 'corteza-webapp-compose/src/components/Admin/Module/RelatedPages'
import Import from 'corteza-webapp-compose/src/components/Admin/Import'
import Export from 'corteza-webapp-compose/src/components/Admin/Export'

export default {
  i18nOptions: {
    namespaces: 'module',
  },

  name: 'ModuleList',

  components: {
    Import,
    Export,
    RelatedPages,
  },

  mixins: [
    listHelpers,
  ],

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },
  },

  data () {
    return {
      primaryKey: 'moduleID',

      filter: {
        query: '',
        namespaceID: this.namespace.namespaceID,
      },

      sorting: {
        sortBy: 'name',
        sortDesc: false,
      },

      creatingRecordPage: false,
    }
  },

  computed: {
    ...mapGetters({
      modules: 'module/set',
      pages: 'page/set',
    }),

    fields () {
      return [
        {
          key: 'name',
          label: this.$t('list.columns.name'),
          sortable: true,
          tdClass: 'text-nowrap',
        },
        {
          key: 'handle',
          label: this.$t('list.columns.handle'),
          sortable: true,
        },
        {
          key: 'changedAt',
          label: this.$t('list.columns.changedAt'),
          sortable: true,
          class: 'text-right text-nowrap',
        },
        {
          key: 'actions',
          label: '',
          tdClass: 'text-right text-nowrap',
        },
      ]
    },

    recordPage () {
      return (moduleID) => this.pages.find(p => p.moduleID === moduleID)
    },
  },

  methods: {
    ...mapActions({
      createPage: 'page/create',
    }),

    handleRowClicked ({ moduleID, canUpdateModule, canDeleteModule }) {
      if (!(canUpdateModule || canDeleteModule)) {
        return
      }
      this.$router.push({
        name: 'admin.modules.edit',
        params: { moduleID },
        query: null,
      })
    },

    encodeRouteParams () {
      const { query } = this.filter
      const { limit, pageCursor, page } = this.pagination

      return {
        query: {
          limit,
          ...this.sorting,
          query,
          page,
          pageCursor,
        },
      }
    },

    items () {
      return this.procListResults(this.$ComposeAPI.moduleList(this.encodeListParams()))
    },

    onImportSuccessful () {
      this.filterList()
      this.toastSuccess(this.$t('notification:general.import.successful'))
    },
  },
}
</script>

<style lang="scss">
.permissions-dropdown {
  .dropdown-item {
    padding: 0;
  }
}
</style>
