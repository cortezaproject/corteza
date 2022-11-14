<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    >
      <b-button
        v-if="canCreate"
        variant="primary"
        class="mr-2"
        :to="{ name: 'federation.nodes.new' }"
      >
        {{ $t('new') }}
      </b-button>

      <b-button
        variant="light"
        @click="openPairModal()"
      >
        {{ $t('pair.label') }}
      </b-button>
    </c-content-header>

    <c-resource-list
      primary-key="nodeID"
      edit-route="federation.nodes.edit"
      :loading-text="$t('loading')"
      :paging="paging"
      :sorting="sorting"
      :items="items"
      :fields="fields"
      @confirm-pending="openConfirmPending($event)"
    >
      <template #filter>
        <b-form-group
          class="p-0 m-0"
        >
          <b-input-group>
            <b-form-input
              v-model.trim="filter.query"
              :placeholder="$t('filterForm.query.placeholder')"
              @keyup="filterList"
            />
          </b-input-group>
        </b-form-group>
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
            <c-submit-button
              button-class="px-4"
              variant="outline-primary"
              icon-variant="text-primary"
              :disabled="!pair.url"
              :processing="pair.processing"
              :success="pair.success"
              @submit="pairNode()"
            >
              {{ $t('pair.confirm') }}
            </c-submit-button>
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
        <c-submit-button
          button-class="px-5 mt-4"
          variant="outline-primary"
          icon-variant="text-primary"
          :processing="pair.processing"
          :success="pair.success"
          @submit="confirmPending()"
        >
          {{ $t('pair.confirm') }}
        </c-submit-button>
      </div>
    </b-modal>
  </b-container>
</template>

<script>
import moment from 'moment'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import CSubmitButton from 'corteza-webapp-admin/src/components/CSubmitButton'
import { mapGetters } from 'vuex'

export default {
  name: 'FederationList',

  i18nOptions: {
    namespaces: 'federation.nodes',
    keyPrefix: 'list',
  },

  components: {
    CSubmitButton,
  },

  mixins: [
    listHelpers,
  ],

  data () {
    return {
      id: 'federation',

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
          label: '',
          tdClass: 'text-right',
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
