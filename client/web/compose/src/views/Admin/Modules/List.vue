"<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column py-3"
  >
    <portal to="topbar-title">
      {{ $t('navigation.module') }}
    </portal>

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
        resourceSingle: $t('general:label.module.single'),
        resourcePlural: $t('general:label.module.plural'),
      }"
      clickable
      sticky-header
      class="h-100 flex-fill"
      @search="filterList"
      @row-clicked="handleRowClicked"
    >
      <template #header>
        <b-btn
          v-if="namespace.canCreateModule"
          data-test-id="button-create"
          variant="primary"
          size="lg"
          :to="{ name: 'admin.modules.create' }"
        >
          {{ $t('createLabel') }}
        </b-btn>

        <import
          v-if="namespace.canCreateModule"
          :namespace="namespace"
          type="module"
          @importSuccessful="onImportSuccessful"
        />

        <export
          v-if="namespace.canExportModules"
          :list="modules"
          type="module"
        />

        <b-dropdown
          v-if="namespace.canGrant"
          size="lg"
          variant="light"
          class="permissions-dropdown"
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
              :button-label="$t('general:label.module.single')"
              :show-button-icon="false"
              button-variant="outline-light"
              class="border-0 text-dark text-left w-100"
            />
          </b-dropdown-item>

          <b-dropdown-item>
            <c-permissions-button
              :resource="`corteza::compose:module-field/${namespace.namespaceID}/*/*`"
              :button-label="$t('general:label.field')"
              :show-button-icon="false"
              button-variant="outline-light"
              class="border-0 text-dark text-left w-100"
            />
          </b-dropdown-item>

          <b-dropdown-item>
            <c-permissions-button
              :resource="`corteza::compose:record/${namespace.namespaceID}/*/*`"
              :button-label="$t('general:label.record')"
              :show-button-icon="false"
              button-variant="outline-light"
              class="border-0 text-dark text-left w-100"
            />
          </b-dropdown-item>
        </b-dropdown>
      </template>

      <template #actions="{ item: m }">
        <related-pages
          :namespace="namespace"
          :module="m"
          size="sm"
          boundary="scrollParent"
          class="d-inline-block"
        />

        <b-dropdown
          v-if="m.canGrant"
          data-test-id="dropdown-permissions"
          size="sm"
          variant="extra-light"
          :title="$t('permissions:resources.compose.module.tooltip')"
          class="permissions-dropdown ml-2"
        >
          <template #button-content>
            <font-awesome-icon :icon="['fas', 'lock']" />
          </template>

          <b-dropdown-item>
            <c-permissions-button
              :title="m.name || m.handle || m.moduleID"
              :target="m.name || m.handle || m.moduleID"
              :resource="`corteza::compose:module/${namespace.namespaceID}/${m.moduleID}`"
              :button-label="$t('general:label.module.single')"
              :show-button-icon="false"
              button-variant="outline-light"
              class="border-0 text-dark text-left w-100"
            />
          </b-dropdown-item>

          <b-dropdown-item>
            <c-permissions-button
              :title="m.name || m.handle || m.moduleID"
              :target="m.name || m.handle || m.moduleID"
              :resource="`corteza::compose:module-field/${namespace.namespaceID}/${m.moduleID}/*`"
              :button-label="$t('general:label.field')"
              :show-button-icon="false"
              all-specific
              button-variant="outline-light"
              class="border-0 text-dark text-left w-100"
            />
          </b-dropdown-item>

          <b-dropdown-item>
            <c-permissions-button
              :title="m.name || m.handle || m.moduleID"
              :target="m.name || m.handle || m.moduleID"
              :resource="`corteza::compose:record/${namespace.namespaceID}/${m.moduleID}/*`"
              :button-label="$t('general:label.record')"
              :show-button-icon="false"
              all-specific
              button-variant="outline-light"
              class="border-0 text-dark text-left w-100"
            />
          </b-dropdown-item>
        </b-dropdown>

        <b-dropdown
          variant="outline-extra-light"
          toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2 ml-2"
          no-caret
          dropleft
          lazy
          menu-class="m-0"
        >
          <template #button-content>
            <font-awesome-icon
              :icon="['fas', 'ellipsis-v']"
            />
          </template>

          <b-dropdown-item
            data-test-id="button-all-records"
            :to="{name: 'admin.modules.record.list', params: { moduleID: m.moduleID }}"
          >
            <font-awesome-icon
              :icon="['fas', 'columns']"
              class="text-primary"
            />
            {{ $t('allRecords.label') }}
          </b-dropdown-item>

          <c-input-confirm
            v-if="m.canDeleteModule"
            :text="$t('list.delete')"
            show-icon
            borderless
            variant="link"
            size="md"
            button-class="dropdown-item text-decoration-none text-dark regular-font rounded-0"
            icon-class="text-danger"
            class="w-100"
            @confirmed="handleDelete(m)"
          />
        </b-dropdown>
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
  </b-container>
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
          tdClass: 'text-right text-nowrap actions',
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
      deletePage: 'page/delete',
      deleteModule: 'module/delete',
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

    handleDelete (module) {
      this.deleteModule(module).then(() => {
        const moduleRecordPage = this.pages.find(p => p.moduleID === module.moduleID)
        if (moduleRecordPage) {
          return this.deletePage({ ...moduleRecordPage, strategy: 'rebase' })
        }
      }).catch(this.toastErrorHandler(this.$t('notification:module.deleteFailed')))
        .finally(() => {
          this.toastSuccess(this.$t('notification:module.deleted'))
          this.filterList()
        })
    },
  },
}
</script>
