<template>
  <b-container
    v-if="isDC !== null"
    fluid
    class="d-flex flex-column p-3"
  >
    <portal to="topbar-title">
      {{ $t('title') }}
    </portal>

    <c-resource-list
      primary-key="requestID"
      :items="items"
      :fields="fields"
      :filter="filter"
      :sorting="sorting"
      :pagination="pagination"
      :translations="{
        notFound: $t('general:resourceList.notFound'),
        noItems: $t('general:resourceList.noItems'),
        loading: $t('general:resourceList.loading'),
        searchPlaceholder: $t('general:resourceList.search.placeholder'),
        showingPagination: 'general:resourceList.pagination.showing',
        singlePluralPagination: 'general:resourceList.pagination.single',
        prevPagination: $t('general:resourceList.pagination.prev'),
        nextPagination: $t('general:resourceList.pagination.next'),
        resourceSingle: $t('general:label.privacy_request.single'),
        resourcePlural: $t('general:label.privacy_request.plural'),
      }"
      :is-item-selectable="isItemSelectable"
      selectable
      clickable
      hide-total
      class="flex-grow-1"
      @search="filterList"
      @row-clicked="rowClicked"
    >
      <template #header="{ selected = [] }">
        <template v-if="isDC">
          <c-input-confirm
            v-if="isDC"
            :disabled="processing || !selected.length"
            :processing="processingApprove"
            :text="$t('request.approve')"
            variant="primary"
            size="lg"
            size-confirm="lg"
            @confirmed="handleSelectedRequests(selected, 'approved')"
          />

          <c-input-confirm
            v-if="isDC"
            :disabled="processing || !selected.length"
            :processing="processingReject"
            :text="$t('request.reject')"
            variant="danger"
            size="lg"
            size-confirm="lg"
            @confirmed="handleSelectedRequests(selected, 'rejected')"
          />
        </template>

        <template v-else>
          <c-input-confirm
            :disabled="processing || !selected.length"
            :processing="processingCancel"
            :text="$t('request.cancel')"
            variant="light"
            size="lg"
            size-confirm="lg"
            @confirmed="handleSelectedRequests(selected, 'canceled')"
          />

          <!-- <b-button
            :disabled="processing"
            variant="light"
            size="lg"
            class="ml-1"
            @click="exportRequests()"
          >
            {{ $t('export') }}
          </b-button> -->
        </template>
      </template>

      <template #status="{ item }">
        <div
          class="d-flex align-items-baseline"
        >
          <span
            class="d-inline-block rounded-circle mr-1"
            :class="`bg-${statusVariants[item.status]}`"
            style="width: 0.6rem; height: 0.6rem;"
          />
          {{ $t(`request:status.${item.status}`) }}
        </div>
      </template>
    </c-resource-list>
  </b-container>
</template>

<script>
import moment from 'moment'
import listHelpers from 'corteza-webapp-privacy/src/mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CResourceList } = components

export default {
  name: 'RequestList',

  components: {
    CResourceList,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'request',
    keyPrefix: 'list',
  },

  data () {
    return {
      processing: false,
      processingApprove: false,
      processingReject: false,
      processingCancel: false,

      isDC: null,

      users: {},

      filter: {
        requestedBy: [],
      },

      sorting: {
        sortBy: 'requestedAt',
        sortDesc: true,
      },

      statusVariants: {
        canceled: 'secondary',
        pending: 'warning',
        rejected: 'danger',
        approved: 'success',
      },
    }
  },

  computed: {
    fields () {
      return [
        {
          key: 'kind',
          sortable: true,
          formatter: kind => this.$t(`request:kind.${kind}`),
        },
        {
          key: 'requestedAt',
          sortable: true,
          formatter: requestedAt => moment(requestedAt).fromNow(),
        },
        {
          hide: !this.isDC,
          key: 'requestedBy',
          sortable: false,
          formatter: requestedBy => this.formatUser(requestedBy),
        },
        {
          key: 'status',
          sortable: true,
        },
      ].filter(({ hide }) => !hide)
        .map(c => ({
          ...c,
          // Generate column label translation key
          label: c.label || this.$t(`columns.${c.key}`),
        }))
    },
  },

  created () {
    this.checkIsDC()
  },

  methods: {
    checkIsDC () {
      this.$SystemAPI.roleList({ query: 'data-privacy-officer', memberID: this.$auth.user.userID })
        .then(({ set = [] }) => {
          if (!set.length) {
            this.filter.requestedBy = [this.$auth.user.userID]
          }

          this.isDC = !!set.length
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
      return this.procListResults(this.$SystemAPI.dataPrivacyRequestList(this.encodeListParams())
        .then(async ({ filter, set }) => {
          if (this.isDC) {
            await this.fetchUsers(set.map(({ requestedBy }) => requestedBy))
          }
          return { filter, set }
        }))
    },

    handleSelectedRequests (selected, status) {
      this.processing = true

      if (status === 'approved') {
        this.processingApprove = true
      } else if (status === 'rejected') {
        this.processingReject = true
      } else {
        this.processingCancel = true
      }

      Promise.all(selected.map(requestID => {
        return this.$SystemAPI.dataPrivacyRequestUpdateStatus({ requestID, status })
      }))
        .then(() => {
          this.$root.$emit('bv::refresh::table', 'resource-list')
        })
        .finally(() => {
          this.processing = false

          if (status === 'approved') {
            this.processingApprove = false
          } else if (status === 'rejected') {
            this.processingReject = false
          } else {
            this.processingCancel = false
          }
        })
    },

    fetchUsers (userID = []) {
      userID = [...new Set(userID)]
      return this.$SystemAPI.userList({ userID })
        .then(({ set }) => {
          set.forEach(user => {
            this.users[user.userID] = user
          })
        })
    },

    isItemSelectable (item) {
      return item.status === 'pending'
    },

    formatUser (userID) {
      const { name, username, email, handle } = this.users[userID]
      return name || username || email || handle || userID || ''
    },

    rowClicked ({ requestID, kind }) {
      this.$router.push({ name: 'request.view', params: { requestID, kind } })
    },
  },
}
</script>
