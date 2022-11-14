<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <b-card
      class="shadow-sm"
      body-class="p-0"
      header-bg-variant="white"
      footer-bg-variant="white"
      footer-class="d-flex align-items-center justify-content-center"
    >
      <template #header>
        <b-form
          @submit.prevent="search"
        >
          <b-form-group
            label-cols-lg="2"
            :label="$t('filter.from')"
          >
            <c-input-date-time
              v-model="filter.from"
              data-test-id="filter-starting-from"
              :labels="{
                clear: $t('general:label.clear'),
                none: $t('general:label.none'),
                now: $t('general:label.now'),
                today: $t('general:label.today'),
              }"
            />
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
            :label="$t('filter.to')"
          >
            <c-input-date-time
              v-model="filter.to"
              data-test-id="filter-ending-at"
              only-past
              :labels="{
                clear: $t('general:label.clear'),
                none: $t('general:label.none'),
                now: $t('general:label.now'),
                today: $t('general:label.today'),
              }"
            />
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
            :label="$t('filter.resource')"
          >
            <b-form-input
              v-model="filter.resource"
              data-test-id="input-resource"
              size="sm"
            />
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
            :label="$t('filter.action')"
          >
            <b-form-input
              v-model="filter.action"
              data-test-id="input-action"
              size="sm"
            />
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
            :label="$t('filter.actor')"
          >
            <b-form-input
              v-model="filter.actorID"
              data-test-id="input-user-id"
              size="sm"
            />
          </b-form-group>

          <b-form-group
            label-cols-lg="2"
          >
            <b-button
              data-test-id="button-submit"
              type="submit"
            >
              {{ $t('filter.search') }}
            </b-button>
          </b-form-group>
        </b-form>
      </template>

      <b-table
        id="resource-list"
        responsive
        hover
        class="mb-0 small"
        head-variant="light"
        :items="items"
        :fields="fields"
        tbody-tr-class="pointer"
        :empty-text="$t('admin:general.notFound')"
        show-empty
        @row-clicked="item=>$set(item, '_showDetails', !item._showDetails)"
      >
        <template #table-busy>
          <div class="text-center m-5">
            <div>
              <b-spinner
                small
                class="align-middle m-2"
              />
            </div>
            <div>{{ $t('loading') }}</div>
          </div>
        </template>
        <template #cell(timestamp)="{ item: a }">
          {{ a.timestamp | locFullDateTime }}
        </template>
        <template #cell(actor)="{ item: a }">
          <router-link
            v-if="a.actorID && a.actorID !== '0'"
            data-test-id="item-user-id"
            :to="drillDownLink({ actorID: a.actorID })"
          >
            {{ a.actor || a.actorID }}
          </router-link>
        </template>
        <template #cell(resource)="{ item: a }">
          <router-link
            data-test-id="item-resource"
            :to="drillDownLink({ resource: a.resource })"
          >
            {{ a.resource }}
          </router-link>
        </template>
        <template #cell(action)="{ item: a }">
          <router-link
            data-test-id="item-action"
            :to="drillDownLink({ action: a.action })"
          >
            {{ a.action }}
          </router-link>
        </template>
        <template #row-details="{ item: a }">
          <b-card-group class="m-3 mb-5">
            <b-card :header="$t('details.header')">
              <b-row>
                <b-col cols="4">
                  {{ $t('details.id') }}
                </b-col>
                <b-col cols="8">
                  {{ a.actionID }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.timestamp') }}
                </b-col>
                <b-col
                  data-test-id="details-timestamp"
                  cols="8"
                >
                  {{ a.timestamp | locFullDateTime }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.requestOrigin') }}
                </b-col>
                <b-col
                  data-test-id="details-request-origin"
                  cols="8"
                >
                  {{ a.requestOrigin }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.requestID') }}
                </b-col>
                <b-col cols="8">
                  {{ a.requestID }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.actorIPAddr') }}
                </b-col>
                <b-col cols="8">
                  {{ a.actorIPAddr }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.actor') }}
                </b-col>
                <b-col
                  data-test-id="details-user"
                  cols="8"
                >
                  {{ a.actor }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.actorID') }}
                </b-col>
                <b-col
                  data-test-id="details-user-id"
                  cols="8"
                >
                  {{ a.actorID }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.severity') }}
                </b-col>
                <b-col
                  data-test-id="details-severity"
                  cols="8"
                >
                  {{ getSeverityLabel(a.severity) }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.resource') }}
                </b-col>
                <b-col
                  data-test-id="details-resource"
                  cols="8"
                >
                  {{ a.resource }}
                </b-col>
              </b-row>
              <b-row>
                <b-col cols="4">
                  {{ $t('details.action') }}
                </b-col>
                <b-col
                  data-test-id="details-action"
                  cols="8"
                >
                  {{ a.action }}
                </b-col>
              </b-row>
            </b-card>
            <b-card :header="$t('details.headerAdditional')">
              <b-row>
                <b-col cols="4">
                  {{ $t('details.description') }}
                </b-col>
                <b-col cols="8">
                  {{ a.description }}
                </b-col>
              </b-row>
              <b-row v-if="a.error">
                <b-col cols="4">
                  {{ $t('details.error') }}
                </b-col>
                <b-col cols="8">
                  {{ a.error }}
                </b-col>
              </b-row>
              <hr>
              <b-row
                v-for="(val, key) in a.meta"
                :key="key"
              >
                <b-col cols="4">
                  <code>{{ key }}</code>
                </b-col>
                <b-col cols="8">
                  <code>{{ val }}</code>
                </b-col>
              </b-row>
            </b-card>
          </b-card-group>
        </template>
      </b-table>

      <template #footer>
        <b-button
          v-if="items.length"
          data-test-id="button-load-older-actions"
          variant="light"
          @click.stop="load()"
        >
          {{ $t('loadOlder') }}
        </b-button>
      </template>
    </b-card>
  </b-container>
</template>

<script>
import { debounce } from 'lodash'
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import { components } from '@cortezaproject/corteza-vue'
const { CInputDateTime } = components

export default {
  components: {
    CInputDateTime,
  },

  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'system.actionlog',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'actionlog',

      filter: {
        from: undefined,
        to: undefined,
        beforeActionID: undefined,
        actorID: undefined,
        resource: undefined,
        action: undefined,
      },

      pagination: {
        limit: 10,
      },

      fields: [
        {
          key: 'timestamp',
          tdClass: 'text-nowrap',
          // formatter: (v) => moment(v).fromNow(),
        },
        {
          key: 'actor',
        },
        {
          key: 'requestOrigin',
        },
        {
          key: 'resource',
        },
        {
          key: 'action',
        },
        {
          key: 'description',
        },
        {
          key: 'severity',
          label: '',
          tdClass: (v, k, item) => ['text-right', this.severity[item.severity].class],
          formatter: (v) => this.severity[v].label,
        },
      ].map(c => ({
        // Generate column label translation key
        label: this.$t(`columns.${c.key}`),
        ...c,
      })),

      items: [],

      severity: [
        { label: this.$t('severity.emergency'),
          class: 'text-danger' },
        { label: this.$t('severity.alert'),
          class: 'text-danger' },
        { label: this.$t('severity.critical'),
          class: 'text-danger' },
        { label: this.$t('severity.error'),
          class: 'text-danger' },
        { label: this.$t('severity.warning'),
          class: 'text-warning' },
        { label: this.$t('severity.notice'),
          class: 'text-success' },
        { label: this.$t('severity.info'),
          class: 'text-success' },
        { label: this.$t('severity.debug'),
          class: '' },
      ],
    }
  },

  mounted () {
    this.load()
  },

  methods: {
    search () {
      // Do a complete search, not just beforeActionID
      this.load(true)
    },

    // Overwrites mixin method
    encodeRouteParams () {
      return {
        query: {
          ...this.pagination,
          ...this.filter,
        },
      }
    },

    load: debounce(function (reset = false) {
      if (reset) {
        this.items.length = 0
        this.pagination.beforeActionID = undefined
      } else {
        const len = this.items.length
        if (len > 0) {
          this.pagination.beforeActionID = (this.items[len - 1] || {}).actionID
        }
      }

      if (!this.filter.actorID) {
        this.$delete(this.filter, 'actorID')
      }

      if (!this.filter.action) {
        this.$delete(this.filter, 'action')
      }

      if (!this.filter.resource) {
        this.$delete(this.filter, 'resource')
      }

      this.procListResults(this.$SystemAPI.actionlogList({ ...this.filter, ...this.pagination }))
        .then(rr => {
          this.items.push(...rr)
        })
    }, 300),

    // Resets pagination & sorting and adds filtering params for drill-down
    drillDownLink (query = {}) {
      return {
        name: 'system.actionlog',
        query: {
          ...this.$route.query,
          ...query,
          sort: undefined,
        },
      }
    },

    getSeverityLabel (index = -1) {
      if (index >= 0) {
        return this.severity[index] ? this.severity[index].label.toLowerCase() : index
      }
    },
  },
}
</script>
