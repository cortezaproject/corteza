<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column flex-fill pt-2 pb-3"
  >
    <c-content-header :title="$t('title')" />

    <c-resource-list
      :primary-key="primaryKey"
      :fields="fields"
      :filter="filter"
      :pagination="pagination"
      :sorting="sorting"
      :items="items"
      :translations="{
        searchPlaceholder: $t('filterForm.query.placeholder'),
        notFound: $t('admin:general.notFound'),
        noItems: $t('admin:general.resource-list.no-items'),
        loading: $t('loading'),
        showingPagination: 'admin:general.pagination.showing',
        singlePluralPagination: 'admin:general.pagination.single',
        prevPagination: $t('admin:general.pagination.prev'),
        nextPagination: $t('admin:general.pagination.next'),
      }"
      clickable
      sticky-header
      hide-total
      class="custom-resource-list-height flex-fill"
      @row-clicked="handleRowClicked"
      @search="filterList"
    >
      <template #header>
        <b-button
          v-if="canCreate"
          variant="primary"
          size="lg"
          :to="{ name: 'federation.nodes.new' }"
        >
          {{ $t('new') }}
        </b-button>

        <b-button
          v-if="canCreate"
          variant="light"
          size="lg"
          @click="openPairModal()"
        >
          {{ $t('pair.label') }}
        </b-button>
      </template>

      <template #actions="{ item: n }">
        <b-dropdown
          v-if="n.nodeID === n.sharedNodeID && (n.status || '').toLowerCase() === 'pair_requested'"
          variant="outline-light"
          toggle-class="d-flex align-items-center justify-content-center text-primary border-0 py-2"
          no-caret
          dropleft
          menu-class="m-0"
        >
          <template #button-content>
            <font-awesome-icon
              :icon="['fas', 'ellipsis-v']"
            />
          </template>

          <b-dropdown-item
            link-class="p-0"
          >
            <b-button
              size="sm"
              variant="link"
              class="text-decoration-none"
              @click="openConfirmPending(n)"
            >
              <font-awesome-icon
                :icon="['fas', 'exclamation-triangle']"
                class="text-danger"
              />
              {{ $t('pair.confirm') }}
            </b-button>
          </b-dropdown-item>
        </b-dropdown>
      </template>
    </c-resource-list>

    <b-modal
      v-model="pair.modal"
      hide-header
      hide-footer
      centered
      size="lg"
      body-class="px-5"
    >
      <div
        v-if="!pair.status"
      >
        <div
          class="text-center px-5"
        >
          <font-awesome-icon
            size="7x"
            :icon="['fas', 'share-alt']"
            class="text-light mb-2"
          />
          <h2>
            {{ $t('pair.status.none.description') }}
          </h2>
        </div>

        <b-input-group
          size="xl"
          class="mt-5"
        >
          <b-form-input
            v-model="pair.url"
            type="url"
            placeholder=""
          />
          <b-input-group-append>
            <c-button-submit
              :disabled="!pair.url"
              :text="$t('pair.confirm')"
              :processing="pair.processing"
              :success="pair.success"
              @submit="pairNode()"
            />
          </b-input-group-append>
        </b-input-group>

        <p
          class="mt-4"
        >
          <strong>{{ $t('pair.note') }}</strong> {{ $t('pair.networkEstablished') }}
        </p>
      </div>

      <div
        v-else-if="pair.status === 'pair-successful'"
      >
        <div
          class="text-center px-5"
        >
          <font-awesome-icon
            size="7x"
            :icon="['far', 'check-circle']"
            class="text-light mb-4"
          />
          <h2>
            {{ $t('pair.status.pending.description') }}
          </h2>
        </div>
      </div>

      <div
        v-else-if="pair.status === 'confirm-pending'"
        class="text-center"
      >
        <div
          class="px-5"
        >
          <font-awesome-icon
            size="7x"
            :icon="['fas', 'share-alt']"
            class="text-light mb-4"
          />
          <h2>
            {{ $t(pair.node.email ? 'pair.status.confirmPending.description' : 'pair.status.confirmPending.descriptionNoMail', pair.node) }}
          </h2>
        </div>

        <c-button-submit
          :processing="pair.processing"
          :success="pair.success"
          :text="$t('pair.confirm')"
          @submit="confirmPending()"
        />
      </div>
    </b-modal>
  </b-container>
</template>

<script>
import moment from 'moment'
import { mapGetters } from 'vuex'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  name: 'FederationList',

  i18nOptions: {
    namespaces: 'federation.nodes',
    keyPrefix: 'list',
  },

  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  data () {
    return {
      id: 'federation',

      primaryKey: 'nodeID',
      editRoute: 'federation.nodes.edit',

      pair: {
        modal: false,
        processing: false,
        success: false,

        url: '',
        status: undefined,
        node: undefined,
      },

      filter: {
        query: '',
        suspended: 0,
        deleted: 0,
      },

      sorting: {
        sortBy: 'createdAt',
        sortDesc: true,
      },

      fields: [
        {
          key: 'name',
          sortable: true,
        },
        // {
        //   key: 'enabled',
        //   sortable: true,
        // },
        {
          key: 'status',
          sortable: true,
        },
        {
          key: 'createdAt',
          label: 'Created',
          sortable: true,
          formatter: (v) => moment(v).fromNow(),
        },
        // {
        //   key: 'tags',
        //   label: '',
        //   sortable: false,
        //   tdClass: 'w-25',
        // },
        {
          key: 'actions',
          class: 'actions',
        },
      ].map(c => ({
        ...c,
        // Generate column label translation key
        label: this.$t(`columns.${c.key}`),
      })),
    }
  },

  computed: {
    ...mapGetters({
      can: 'rbac/can',
    }),

    canCreate () {
      return this.can('federation/', 'node.create')
    },
  },

  methods: {
    items () {
      return this.procListResults(this.$FederationAPI.nodeSearch(this.encodeListParams()))
    },

    openPairModal () {
      this.pair.status = undefined
      this.pair.modal = true
    },

    async pairNode () {
      this.pair.processing = true

      await this.$FederationAPI.nodeCreate({ pairingURI: this.pair.url })
        .then(async node => {
          await this.$FederationAPI.nodePair(node)
          this.pair.url = ''
          this.pair.status = 'pair-successful'
          this.pair.node = node

          // Refetch list
          this.$root.$emit('bv::refresh::table', 'resource-list')
          this.toastSuccess(this.$t('notification:federation.pair.success'))
        })
        .catch(this.toastErrorHandler(this.$t('notification:federation.pair.error')))
        .finally(() => {
          this.pair.processing = false
        })
    },

    openConfirmPending (node) {
      this.pair.node = node
      this.pair.status = 'confirm-pending'
      this.pair.modal = true
    },

    async confirmPending () {
      this.pair.processing = true

      await this.$FederationAPI.nodeHandshakeConfirm({ nodeID: this.pair.node.nodeID })
        .then(() => {
          this.pair.success = true

          // Refetch list
          this.$root.$emit('bv::refresh::table', 'resource-list')
          this.toastSuccess(this.$t('notification:federation.handshake.success'))

          setTimeout(() => {
            this.pair.success = false
          }, 2000)

          setTimeout(() => {
            this.pair.node = undefined
            this.pair.status = undefined
            this.pair.modal = false
          }, 1000)
        })
        .catch(this.toastErrorHandler(this.$t('notification:federation.handshake.error')))
        .finally(() => {
          this.pair.processing = false
        })
    },
  },
}
</script>

<style lang="scss">
.pointer {
  cursor: pointer;
}
</style>
