<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    />

    <b-card
      no-body
      class="shadow-sm"
      header-bg-variant="white"
    >
      <template #header>
        <b-form
          @submit.prevent="search"
        >
          <b-form-group
            label-cols-lg="2"
            :label="$t('filter.searchQuery')"
          >
            <b-form-input
              v-model="filter.query"
              size="sm"
            />
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
          >
            <b-form-checkbox
              v-model="filter.incScriptsWithErrors"
              size="sm"
            >
              {{ $t('filter.incScriptsWithErrors', { count: totalScriptsWithErrors}) }}
            </b-form-checkbox>
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
          >
            <b-form-checkbox
              v-model="filter.incScriptsWithTriggers"
              size="sm"
            >
              {{ $t('filter.incScriptsWithTriggers', { count: totalScriptsWithTriggers}) }}
            </b-form-checkbox>
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
          >
            <b-form-checkbox
              v-model="filter.incScriptsWithIterator"
              size="sm"
            >
              {{ $t('filter.incScriptsWithIterator', { count: totalScriptsWithIterator}) }}
            </b-form-checkbox>
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
          >
            <b-form-checkbox
              v-model="filter.incScriptsWithSecurity"
              size="sm"
            >
              {{ $t('filter.incScriptsWithSecurity', { count: totalScriptsWithSecurity}) }}
            </b-form-checkbox>
          </b-form-group>
          <b-form-group
            label-cols-lg="2"
          >
            <b-form-checkbox
              v-model="filter.absoluteTime"
              size="sm"
            >
              {{ $t('filter.absoluteTime') }}
            </b-form-checkbox>
          </b-form-group>
        </b-form>
      </template>

      <b-card-body
        class="p-0"
      >
        <b-table
          id="resource-list"
          hover
          responsive
          class="mb-0"
          head-variant="light"
          :items="filtered"
          :fields="fields"
          :empty-text="$t('admin:general.notFound')"
          show-empty
          no-sort-reset
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
          <template #cell(label)="{}" />
          <template #cell(name)="{ item: { label, name, errors, description, ...r }, toggleDetails }">
            <div>
              <span
                v-if="label"
              >
                {{ label }}
              </span>
              <span
                v-else
                class="text-secondary"
              >{{ $t('labelMissing') }}
              </span>
              <b-badge
                v-if="r.security"
                class="rounded m-1 py-1 px-2 pointer"
                variant="primary"
                @click="toggleDetails"
              >
                {{ $t('flags.security') }}
              </b-badge>
              <b-badge
                v-if="r.triggers"
                class="rounded m-1 py-1 px-2 pointer"
                variant="primary"
                @click="toggleDetails"
              >
                {{ $t('flags.triggers') }}
              </b-badge>
              <b-badge
                v-if="r.iterator"
                class="rounded m-1 py-1 px-2 pointer"
                variant="primary"
                @click="toggleDetails"
              >
                {{ $t('flags.iterator') }}
              </b-badge>
            </div>
            <p
              v-if="description"
              class="text-secondary"
            >
              {{ description }}
            </p>
            <div><small><code>{{ name }}</code></small></div>
            <b-alert
              v-for="(error, i) in errors"
              :key="i"
              show
              variant="warning"
            >
              {{ error }}
            </b-alert>
          </template>
          <template v-slot:row-details="{ item: r }">
            <b-card>
              <pre>{{ r.triggers }}</pre>
              <pre>{{ r.iterator }}</pre>
              <pre>{{ r.security }}</pre>
            </b-card>
          </template>
          <template #cell(updatedAt)="{ value }">
            <time
              :datetime="value.toISOString()"
              :title="value"
            >
              {{ filter.absoluteTime ? value : value.fromNow() }}
            </time>
          </template>
        </b-table>
      </b-card-body>
    </b-card>
  </b-container>
</template>

<script>
import listHelpers from 'corteza-webapp-admin/src/mixins/listHelpers'
import moment from 'moment'

export default {
  mixins: [
    listHelpers,
  ],

  i18nOptions: {
    namespaces: 'automation.scripts',
    keyPrefix: 'list',
  },

  data () {
    return {
      id: 'automation',

      items: [],

      filter: {
        query: '',
        incScriptsWithErrors: false,
        incScriptsWithTriggers: false,
        incScriptsWithIterator: false,
        incScriptsWithSecurity: false,
        absoluteTime: false,
      },

      fields: [
        {
          key: 'name',
          label: '',
          sortable: true,
        },
        {
          key: 'updatedAt',
          sortable: true,
          tdClass: 'text-right text-nowrap',
          formatter: (v) => moment(v),
        },
      ].map(c => ({
        // Generate column label translation key
        label: this.$t(`columns.${c.key}`),
        ...c,
      })),
    }
  },

  computed: {
    filtered () {
      const {
        query,
        incScriptsWithErrors,
        incScriptsWithTriggers,
        incScriptsWithIterator,
        incScriptsWithSecurity,
      } = this.filter
      const lcQuery = query.toLocaleLowerCase()

      return this.items
        .filter(({ name, label }) => (lcQuery.length === 0 || (name + ' ' + label).toLocaleLowerCase().indexOf(lcQuery) > -1))
        .filter(({ errors }) => (incScriptsWithErrors === false || (errors && errors.length > 0)))
        .filter(({ triggers }) => (incScriptsWithTriggers === false || !!triggers))
        .filter(({ iterator }) => (incScriptsWithIterator === false || !!iterator))
        .filter(({ security }) => (incScriptsWithSecurity === false || !!security))
    },

    totalScriptsWithErrors () {
      return this.items.filter(({ errors }) => (errors && errors.length > 0)).length
    },

    totalScriptsWithSecurity () {
      return this.items.filter(({ security }) => (security)).length
    },

    totalScriptsWithTriggers () {
      return this.items.filter(({ triggers }) => (triggers)).length
    },

    totalScriptsWithIterator () {
      return this.items.filter(({ iterator }) => (iterator)).length
    },
  },

  created () {
    this.procListResults(this.$SystemAPI.automationList(this.encodeListParams()))
      .then(set => { this.items = set || [] })
  },

  methods: {
    events (tt) {
      const ee = []

      if (!Array.isArray(tt) || tt.length === 0) {
        return ee
      }

      tt.forEach(({ events }) => ee.push(...(events || [])))
      return ee.filter((v, i) => ee.indexOf(v) === i)
    },
  },
}
</script>
<style lang="scss">
.pointer {
  cursor: pointer;
}
</style>
