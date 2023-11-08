<template>
  <b-container
    class="d-flex flex-column pt-2 pb-3 flex-fill"
  >
    <c-content-header
      :title="$t('title')"
    >
      <c-corredor-manual-buttons
        ui-page="dashboard"
        ui-slot="toolbar"
        resource-type="system'"
        @click="dispatchCortezaSystemEvent($event)"
      />
    </c-content-header>

    <b-row
      class="flex-fill"
    >
      <b-col
        cols="12"
      >
        <b-card
          body-class="position-relative p-0"
          header-class="bg-white"
          footer-class="bg-white"
          class="shadow-sm h-100"
        >
          <template #header>
            <b-card-title title-tag="h3">
              <router-link
                :to="{ name: 'system.user.list' }"
                :area-label="`${users.valid} ${$t('users.title')}`"
                class="display-3 text-decoration-none"
              >
                {{ users.valid }}
              </router-link>
            </b-card-title>
            <b-card-sub-title sub-title-tag="h4">
              {{ $t('users.title') }}
            </b-card-sub-title>
          </template>

          <c-chart
            v-if="userChart"
            :chart="userChart"
          />

          <template #footer>
            <b-row>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.user.list', query: { deleted: 1, suspended: 1 } }"
                  :aria-label="users.total + ' ' + $t('users.users') + ' ' + $t('users.total')"
                  class="text-decoration-none"
                >
                  {{ users.total }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('users.total') }}
                </span>
              </b-col>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.user.list', query: { deleted: 1, suspended: 2 } }"
                  :aria-label="users.suspended + ' ' + $t('users.users') + ' ' + $t('users.suspended')"
                  class="text-decoration-none"
                >
                  {{ users.suspended }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('users.suspended') }}
                </span>
              </b-col>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.user.list', query: { deleted: 2, suspended: 1 } }"
                  :aria-label="users.deleted + ' ' + $t('users.users') + ' ' + $t('users.deleted')"
                  class="text-decoration-none"
                >
                  {{ users.deleted }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('users.deleted') }}
                </span>
              </b-col>
            </b-row>
          </template>
        </b-card>
      </b-col>
    </b-row>

    <b-row align-v="stretch">
      <b-col
        v-show="roles.total"
        cols="12"
        md="6"
        class="mt-3"
      >
        <b-card
          no-body
          header-class="bg-white"
          footer-class="bg-white"
          class="shadow-sm"
        >
          <template #header>
            <b-card-title title-tag="h3">
              <router-link
                :to="{ name: 'system.role.list' }"
                :aria-label="roles.valid + ' ' + $t('roles.title')"
                class="display-4 text-decoration-none"
              >
                {{ roles.valid }}
              </router-link>
            </b-card-title>
            <b-card-sub-title sub-title-tag="h4">
              {{ $t('roles.title') }}
            </b-card-sub-title>
          </template>

          <template #footer>
            <b-row>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.role.list', query: { deleted: 1, archived: 1 } }"
                  :aria-label="roles.total + ' ' + $t('roles.roles') + ' ' + $t('roles.total')"
                  class="text-decoration-none"
                >
                  {{ roles.total }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('roles.total') }}
                </span>
              </b-col>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.role.list', query: { deleted: 1, archived: 2 } }"
                  :aria-label="roles.archived + ' ' + $t('roles.roles') + ' ' + $t('roles.archived')"
                  class="text-decoration-none"
                >
                  {{ roles.archived }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('roles.archived') }}
                </span>
              </b-col>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.role.list', query: { deleted: 2, archived: 1 } }"
                  :aria-label="roles.deleted + ' ' + $t('roles.roles') + ' ' + $t('roles.deleted')"
                  class="text-decoration-none"
                >
                  {{ roles.deleted }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('roles.deleted') }}
                </span>
              </b-col>
            </b-row>
          </template>
        </b-card>
      </b-col>

      <b-col
        v-show="applications.total"
        cols="12"
        md="6"
        class="mt-3"
      >
        <b-card
          no-body
          header-class="bg-white"
          footer-class="bg-white"
          class="shadow-sm"
        >
          <template #header>
            <b-card-title title-tag="h3">
              <router-link
                :to="{ name: 'system.application.list' }"
                :aria-label="applications.valid + ' ' + $t('applications.title')"
                class="display-4 text-decoration-none"
              >
                {{ applications.valid }}
              </router-link>
            </b-card-title>
            <b-card-sub-title sub-title-tag="h4">
              {{ $t('applications.title') }}
            </b-card-sub-title>
          </template>

          <template #footer>
            <b-row>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.application.list', query: { deleted: 1 } }"
                  :aria-label="applications.total + ' ' + $t('applications.applications') + ' ' + $t('applications.total')"
                  class="text-decoration-none"
                >
                  {{ applications.total }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('applications.total') }}
                </span>
              </b-col>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.application.list', query: { deleted: 2 } }"
                  :aria-label="applications.deleted + ' ' + $t('applications.applications') + ' ' + $t('applications.deleted')"
                  class="text-decoration-none"
                >
                  {{ applications.deleted }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('applications.deleted') }}
                </span>
              </b-col>
            </b-row>
          </template>
        </b-card>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import * as moment from 'moment'
import { components } from '@cortezaproject/corteza-vue'
const { CChart } = components

export default {
  i18nOptions: {
    namespaces: 'dashboard',
  },

  components: {
    CChart,
  },

  data () {
    return {
      userChart: null,
      users: {
        total: 0,
        valid: 0,
        deleted: 0,
        suspended: 0,
        dailyCreated: [],
        dailyUpdated: [],
        dailySuspended: [],
        dailyDeleted: [],
      },
      roles: {
        total: 0,
        valid: 0,
        archived: 0,
        deleted: 0,
      },
      applications: {
        total: 0,
        valid: 0,
        deleted: 0,
      },
    }
  },

  mounted () {
    this.$SystemAPI.statsList().then(({ users, roles, applications }) => {
      if (users) {
        this.users = users
      }

      if (roles) {
        this.roles = roles
      }

      if (applications) {
        this.applications = applications
      }

      this.initUserChart()
    })
  },

  methods: {
    initUserChart () {
      if (this.users.total === 0) {
        return
      }

      const { dates, values } = this.getUserTimeline()

      this.userChart = {
        tooltip: {
          trigger: 'axis',
        },
        textStyle: {
          fontFamily: 'Poppins-Regular',
        },
        xAxis: {
          type: 'category',
          data: dates,
          boundaryGap: false,
        },
        yAxis: {
          type: 'value',
        },
        grid: {
          top: 20,
          right: 50,
          bottom: 20,
          left: 40,
          containLabel: true,
        },
        series: [
          {
            name: this.$t('users.created'),
            type: 'line',
            data: values,
            smooth: 0.5,
            areaStyle: {
              opacity: 0.5,
            },
          },
        ],
      }
    },

    getUserTimeline () {
      const data = this.users.dailyCreated
      const unit = this.getComfortableTimeUnit(data)

      let aux = {}
      for (let i = 0; i < data.length; i += 2) {
        const ts = moment.unix(data[i]).startOf(unit).format(unit === 'month' ? 'MMM YYYY' : 'D MMM YYYY')
        aux[ts] = (aux[ts] || 0) + data[i + 1]
      }

      const dates = []
      const values = []

      for (let date in aux) {
        dates.push(date)
        values.push(aux[date])
      }

      return { dates, values }
    },

    /**
     * Takes standard array with date+value pairs and returns comfortable time unit
     * depending on time span between min & max unit
     *
     * @param {number[]} range
     */
    getComfortableTimeUnit (range) {
      if (range.length === 0) {
        return undefined
      }

      if (range.length === 2) {
        return 'day'
      }

      const ts = range.filter((v, i) => i % 2 === 0).sort()
      const min = ts[0]
      const max = ts[ts.length - 1]
      const diffInDays = (max - min) / (60 * 60 * 24)

      if (diffInDays > 120) {
        return 'month'
      }

      return 'day'
    },
  },
}
</script>
