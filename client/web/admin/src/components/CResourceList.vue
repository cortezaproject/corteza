<template>
  <b-card
    no-body
    class="shadow-sm"
    header-bg-variant="white"
  >
    <b-card-body
      class="p-0"
    >
      <b-table
        id="resource-list"
        hover
        responsive
        head-variant="light"
        class="mb-0"
        :primary-key="primaryKey"
        :sort-by.sync="sorting.sortBy"
        :sort-desc.sync="sorting.sortDesc"
        :items="items"
        :fields="fields"
        :empty-text="$t('notFound')"
        show-empty
        no-sort-reset
        no-local-sorting
      >
        <template #table-busy>
          <div class="text-center m-5">
            <div>
              <b-spinner
                class="align-middle m-2"
              />
            </div>
            <div>{{ loadingText }}</div>
          </div>
        </template>
        <template #cell(tags)="row">
          <div
            class="d-flex flex-wrap h3"
          >
            <b-badge
              v-for="(t, index) in row.item.tags"
              :key="index"
              variant="warning"
              class="rounded mr-1 py-1 px-2"
            >
              {{ t }}
            </b-badge>
          </div>
        </template>
        <template #cell(actions)="row">
          <b-button
            v-if="row.item.nodeID === row.item.sharedNodeID && (row.item.status || '').toLowerCase() === 'pair_requested'"
            size="sm"
            variant="link"
            class="p-0 pr-1"
            @click="$emit('confirm-pending', row.item)"
          >
            <font-awesome-icon
              :icon="['fas', 'exclamation-triangle']"
              class="text-danger"
            />
          </b-button>
          <b-button
            v-if="row.item.roleID !== '1'"
            size="sm"
            variant="link"
            :to="{ name: editRoute, params: { [primaryKey]: row.item[primaryKey] } }"
          >
            <font-awesome-icon
              :icon="['fas', 'pen']"
            />
          </b-button>
        </template>
      </b-table>
    </b-card-body>

    <!--
      Card header
    -->
    <template #header>
      <b-container
        class="p-0"
        fluid
      >
        <b-row
          align-v="end"
        >
          <b-col
            lg="9"
          >
            <slot
              name="filter"
            />
          </b-col>
        </b-row>
      </b-container>
    </template>

    <!--
      Card footer
    -->
    <template #footer>
      <b-button-group class="float-right">
        <b-button
          :disabled="hasPrevPage"
          variant="link"
          class="text-dark"
          @click="goToPage()"
        >
          <font-awesome-icon :icon="['fas', 'angle-double-left']" />
        </b-button>
        <b-button
          :disabled="hasPrevPage"
          data-test-id="button-previous-page"
          variant="link"
          class="text-dark"
          @click="goToPage('prevPage')"
        >
          <font-awesome-icon :icon="['fas', 'angle-left']" />
          {{ $t('pagination.prev') }}
        </b-button>
        <b-button
          :disabled="hasNextPage"
          data-test-id="button-next-page"
          variant="link"
          class="text-dark"
          @click="goToPage('nextPage')"
        >
          {{ $t('pagination.next') }}
          <font-awesome-icon :icon="['fas', 'angle-right']" />
        </b-button>
      </b-button-group>
    </template>
  </b-card>
</template>
<script>

export default {
  name: 'CResourceList',

  props: {
    loadingText: {
      type: String,
      default: 'Loading',
    },

    editRoute: {
      type: String,
      required: true,
    },

    primaryKey: {
      type: String,
      required: true,
    },

    sorting: {
      type: Object,
      required: true,
    },

    paging: {
      type: Object,
      required: true,
    },

    fields: {
      type: Array,
      required: true,
    },

    items: {
      type: Function,
      required: true,
    },
  },

  i18nOptions: {
    namespaces: 'admin',
    keyPrefix: 'general',
  },

  computed: {
    hasPrevPage () {
      return !this.paging.prevPage
    },

    hasNextPage () {
      return !this.paging.nextPage
    },
  },

  methods: {
    goToPage (page) {
      let pageCursor = this.paging[page] || ''
      this.$router.push({ path: this.$route.path, query: { ...this.$route.query, pageCursor } })
    },
  },
}
</script>
