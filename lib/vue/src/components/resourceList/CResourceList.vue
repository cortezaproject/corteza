<template>
  <b-card
    no-body
    header-bg-variant="white"
    footer-bg-variant="white"
    footer-class="p-0 border-top"
    :header-class="cardHeaderClass"
    class="shadow-sm"
  >
    <template #header>
      <b-container
        fluid
        class="d-flex flex-column p-0 gap-2 d-print-none"
      >
        <b-row
          no-gutters
          class="d-flex align-items-center justify-content-between gap-1"
        >
          <div :class="`d-flex align-items-center flex-grow-1 flex-wrap flex-fill-child gap-1 ${headerClass}`">
            <slot
              name="header"
              :selected="selected"
            />
          </div>

          <div
            v-if="!hideSearch"
            class="flex-fill"
          >
            <c-input-search
              v-model.trim="filter[queryField]"
              :placeholder="translations.searchPlaceholder"
              :debounce="300"
              @input="$emit('search')"
            />
          </div>
        </b-row>

        <b-row
          v-if="$slots.toolbar"
          class="gap-1"
        >
          <slot
            name="toolbar"
          />
        </b-row>
      </b-container>
    </template>

    <b-card-body
      class="p-0"
    >
      <b-table
        id="resource-list"
        ref="resourceList"
        head-variant="light"
        :fields="_fields"
        :items="_items"
        :sort-by.sync="sorting.sortBy"
        :sort-desc.sync="sorting.sortDesc"
        :sticky-header="stickyHeader"
        :tbody-tr-class="tableRowClasses"
        hover
        responsive
        show-empty
        no-sort-reset
        no-local-sorting
        :primary-key="primaryKey"
        class="mh-100 h-100 mb-0"
        @sort-changed="pagination.page = 1"
        @row-clicked="$emit('row-clicked', $event)"
      >
        <template #empty>
          <p
            data-test-id="no-matches"
            class="text-center text-dark"
            style="margin-top: 1vh;"
          >
            {{ translations.noItems }}
          </p>
        </template>

        <template #table-busy>
          <div class="text-center m-5">
            <div>
              <b-spinner
                class="align-middle m-2"
              />
            </div>
            {{ translations.loading }}
          </div>
        </template>

        <template
          v-if="selectable"
          #head(select)
        >
          <b-checkbox
            :disabled="disableSelectAll"
            :checked="allRowsSelected && !disableSelectAll"
            class="ml-1"
            @change="selectAllRows"
          />
        </template>

        <template #cell(select)="{ item }">
          <b-form-checkbox
            v-if="isItemSelectable(item)"
            class="ml-1"
            :checked="selected.includes(item[primaryKey])"
            @change="onSelectRow($event, item[primaryKey])"
          />
        </template>

        <!-- Magic; Make slots if parent provides them -->
        <template
          v-for="field in customFieldSlots"
          #[`cell(${field})`]="{ item }"
        >
          <slot
            :name="field"
            :item="item"
          />
        </template>
      </b-table>
    </b-card-body>

    <template
      v-if="showFooter"
      #footer
    >
      <div
        class="d-flex align-items-center flex-wrap justify-content-between p-2 w-100"
      >
        <div class="d-flex gap-col-3 align-items-center flex-wrap">
          <div
            v-if="!hideTotal"
            class="text-nowrap ml-2"
          >
            {{ getPagination }}
          </div>
        
          <div class="d-flex align-items-center ml-2 my-1 gap-1 text-nowrap">
            <span>
                {{ $t('general:resourceList.pagination.recordsPerPage') }}
            </span>

            <b-form-select
              :value="pagination.limit"
              :options="perPageOptions"
              @change="handlePerPageChange"
            />
          </div>
        </div>

        <div
          v-if="!hidePagination"
          class="d-flex align-items-center justify-content-end"
        >
          <b-button-group>
            <b-button
              :disabled="!hasPrevPage"
              variant="outline-light"
              class="d-flex align-items-center text-primary border-0"
              @click="goToPage()"
            >
              <font-awesome-icon :icon="['fas', 'angle-double-left']" />
            </b-button>

            <b-button
              :disabled="!hasPrevPage"
              variant="outline-light"
              class="d-flex align-items-center text-primary border-0"
              @click="goToPage('prevPage')"
            >
              <font-awesome-icon
                :icon="['fas', 'angle-left']"
                class="mr-1"
              />

              {{ translations.prevPagination }}
            </b-button>

            <b-button
              :disabled="!hasNextPage"
              variant="outline-light"
              class="d-flex align-items-center justify-content-center text-primary border-0"
              @click="goToPage('nextPage')"
            >
              {{ translations.nextPagination }}

              <font-awesome-icon
                :icon="['fas', 'angle-right']"
                class="ml-1"
              />
            </b-button>
          </b-button-group>
        </div>
      </div>
    </template>
  </b-card>
</template>
<script>
import CInputSearch from '../input/CInputSearch.vue'

export default {
  name: 'ResourceList',

  components: {
    CInputSearch,
  },

  props: {
    primaryKey: {
      type: String,
      required: true,
    },

    filter: {
      type: Object,
      required: true,
    },

    sorting: {
      type: Object,
      required: true,
    },

    pagination: {
      type: Object,
      required: true,
    },

    fields: {
      type: Array,
      required: true,
    },

    // Promise that resolves with an array
    items: {
      type: Function,
      required: true,
    },

    hideSearch: {
      type: Boolean,
    },

    hideTotal: {
      type: Boolean,
    },

    hidePagination: {
      type: Boolean,
    },

    stickyHeader: {
      type: Boolean,
    },

    // Are rows clickable
    clickable: {
      type: Boolean,
      default: false,
    },

    selectable: {
      type: Boolean,
      default: false,
    },

    isItemSelectable: {
      type: Function,
      default: () => true,
    },

    cardHeaderClass: {
      type: String,
      default: '',
    },

    headerClass: {
      type: String,
      default: '',
    },

    rowClass: {
      type: Function,
      default: () => {},
    },

    translations: {
      type: Object,
      required: true,
    },

    queryField: {
      type: String,
      default: 'query',
    },

    hidePerPageOption: {
      type: Boolean,
      default: false
    }
  },

  data () {
    return {
      selected: [],
      selectableItemIDs: [],
    }
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  computed: {
    _fields () {
      const select = this.selectable ? [
        {
          key: 'select',
          label: '',
          thStyle: 'width: 0; white-space: nowrap;',
        }
      ] : []

      return [
        ...select,
        ...this.fields,
      ].map(f => {
        return { ...f, thClass: `${f.thClass || 'border-0'}` }
      })
    },

    customFieldSlots () {
      return [
        ...Object.keys(this.$slots),
        ...Object.keys(this.$scopedSlots),
      ].filter(s => s !== 'header')
    },

    perPageOptions () {
      const defaultText = this.pagination.limit === 0 ? this.$t('general:label.all') : this.pagination.limit.toString()

      return [
        { text: defaultText, value: this.pagination.limit },
        { text: '25', value: 25 },
        { text: '50', value: 50 },
        { text: '100', value: 100 },
      ].filter((v, i) => i === 0 || v.value !== this.pagination.limit).sort((a, b) => {
        if (a.value === 0) return 1
        if (b.value === 0) return -1
        return a.value - b.value
      })
    },

    disableSelectAll () {
      return !this.selectableItemIDs.length
    },

    allRowsSelected () {
      return this.selected.length === this.selectableItemIDs.length
    },

    getPagination () {
      let { total = 0, limit = 10, page = 1 } = this.pagination
      total = isNaN(total) ? 0 : total

      const pagination = {
        from: ((page - 1) * limit) + 1,
        to: limit > 0 ? Math.min((page * limit), total) : total,
        count: total,
        data: total == 1 ? this.translations.resourceSingle : this.translations.resourcePlural
      }

      return this.$t(this.translations[total > limit ? 'showingPagination' : 'singlePluralPagination'], pagination)
    },

    hasPrevPage () {
      return !!this.pagination.prevPage
    },

    hasNextPage () {
      return !!this.pagination.nextPage
    },

    showFooter () {
      return !(this.hideTotal && this.hidePagination && this.hidePerPageOption)
    }
  },

  methods: {
    tableRowClasses (item = {}) {
      return {
        'pointer': this.clickable,
        ...this.rowClass(item),
      }
    },

    _items () {
      this.selected = []
      this.selectableItemIDs = []

      return this.items().then((items = []) => {
        this.selectableItemIDs = items.filter(this.isItemSelectable).map(i => i[this.primaryKey])
        return items
      })
    },

    onSelectRow (selected, itemID) {
      if (selected) {
        if (this.selected.includes(itemID)) {
          return
        }

        this.selected.push(itemID)
      } else {
        const i = this.selected.indexOf(itemID)
        if (i < 0) {
          return
        }

        this.selected.splice(i, 1)
      }
    },

    selectAllRows (allSelected = false) {
      if (allSelected) {
        this.selected = this.selectableItemIDs
      } else {
        this.selected = []
      }
    },

    goToPage (destination) {
      const pageCursor = this.pagination[destination] || ''

      let { page = 1 } = this.pagination

      if (destination === 'nextPage') {
        page++
      } else if (destination === 'prevPage') {
        page--
      } else {
        page = 1
      }

      this.$router.replace({ query: { ...this.$route.query, page, pageCursor } })
    },

    handlePerPageChange (limit) {
      this.$router.replace({ query: { ...this.$route.query, page: 1, limit } })
      this.$refs.resourceList.refresh()
    },

    setDefaultValues () {
      this.selected = [],
      this.selectableItemIDs = []
    },
  },
}
</script>

<style lang="scss">
#resource-list {
  th {
    background-color: var(--gray-200) !important;
  }

  td.actions {
    padding-top: 8px;
    right: 0;
    opacity: 0;
    position: sticky;
    transition: opacity 0.25s;
    width: 1%;

    .regular-font {
      font-family: 'Poppins-Regular' !important;
    }
  }

  tr:hover td.actions {
    opacity: 1;
    background-color: var(--gray-200);
    z-index: 1;
  }
}

@media (max-width: 576px) {
  .flex-fill-child > * {
    flex: 1 1 auto !important;
  }
}
</style>
