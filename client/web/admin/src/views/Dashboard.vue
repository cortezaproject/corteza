<template>
  <b-container
    class="py-3"
  >
    <c-content-header
      :title="$t('title')"
    >
      <c-corredor-manual-buttons
        ui-page="dashboard"
        ui-slot="toolbar"
        resource-type="system'"
        class="mr-1"
        @click="dispatchCortezaSystemEvent($event)"
      />
    </c-content-header>
    <b-row>
      <b-col
        v-show="users.total"
        cols="12"
      >
        <b-card
          class="shadow-sm mb-3"
        >
          <b-card-title title-tag="h3">
            <router-link
              class="display-3"
              :to="{ name: 'system.user.list' }"
              :area-label="users.valid + ' ' + $t('users.title')"
            >
              {{ users.valid }}
            </router-link>
          </b-card-title>
          <b-card-sub-title sub-title-tag="h4">
            {{ $t('users.title') }}
          </b-card-sub-title>

          <canvas
            ref="userChart"
            style="height: 100px; width: 100%"
          />

          <b-container class="mt-3">
            <b-row>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.user.list', query: { deleted: 1, suspended: 1 } }"
                  :aria-label="users.total + ' ' + $t('users.users') + ' ' + $t('users.total')"
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
                >
                  {{ users.deleted }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('users.deleted') }}
                </span>
              </b-col>
            </b-row>
          </b-container>
        </b-card>
      </b-col>
      <b-col
        v-show="roles.total"
        cols="12"
        md="6"
      >
        <b-card
          class="shadow-sm mb-3"
        >
          <b-card-title title-tag="h3">
            <router-link
              class="display-4"
              :to="{ name: 'system.role.list' }"
              :aria-label="roles.valid + ' ' + $t('roles.title')"
            >
              {{ roles.valid }}
            </router-link>
          </b-card-title>
          <b-card-sub-title sub-title-tag="h4">
            {{ $t('roles.title') }}
          </b-card-sub-title>

          <b-container class="mt-3">
            <b-row>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.role.list', query: { deleted: 1, archived: 1 } }"
                  :aria-label="roles.total + ' ' + $t('roles.roles') + ' ' + $t('roles.total')"
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
                >
                  {{ roles.deleted }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('roles.deleted') }}
                </span>
              </b-col>
            </b-row>
          </b-container>
        </b-card>
      </b-col>
      <b-col
        v-show="applications.total"
        cols="12"
        md="6"
      >
        <b-card
          class="shadow-sm mb-3"
        >
          <b-card-title title-tag="h3">
            <router-link
              class="display-4"
              :to="{ name: 'system.application.list' }"
              :aria-label="applications.valid + ' ' + $t('applications.title')"
            >
              {{ applications.valid }}
            </router-link>
          </b-card-title>
          <b-card-sub-title sub-title-tag="h4">
            {{ $t('applications.title') }}
          </b-card-sub-title>

          <b-container class="mt-3">
            <b-row>
              <b-col
                cols="12"
                sm="4"
                class="mb-2 mb-sm-0"
              >
                <router-link
                  :to="{ name: 'system.application.list', query: { deleted: 1 } }"
                  :aria-label="applications.total + ' ' + $t('applications.applications') + ' ' + $t('applications.total')"
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
                >
                  {{ applications.deleted }}
                </router-link>
                <span class="text-secondary d-sm-block">
                  {{ $t('applications.deleted') }}
                </span>
              </b-col>
            </b-row>
          </b-container>
        </b-card>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import ChartJS from 'chart.js'
import * as moment from 'moment'

export default {
  i18nOptions: {
    namespaces: 'dashboard',
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

  computed: {
    userTimeline () {
      const d = this.users.dailyCreated
      const unit = this.getComfortableTimeUnit(d)

      let aux = {}
      for (let i = 0; i < d.length; i += 2) {
        const ts = moment.unix(d[i]).startOf(unit).toISOString()
        aux[ts] = (aux[ts] || 0) + d[i + 1]
      }

      let fmt = []
      for (let x in aux) {
        fmt.push({ x, y: aux[x] })
      }

      return fmt
    },
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

    // this.$nextTick(() => {
    //   this.$EventBus.Dispatch({
    //     resourceType: 'ui:admin:dashboard',
    //     eventType: 'afterMount',
    //     args: { data: this.$data, $el: this.$el },
    //   })
    // })
  },

  // updated () {
  //   this.$nextTick(() => {
  //     this.$EventBus.Dispatch({
  //       resourceType: 'ui:admin:dashboard',
  //       eventType: 'afterUpdate',
  //       args: { data: this.users, $el: this.$el },
  //     })
  //   })
  // },

  methods: {
    initUserChart () {
      if (this.users.total === 0) {
        return
      }

      const d = this.users.dailyCreated
      const unit = this.getComfortableTimeUnit(d)
      const ctx = this.$refs.userChart.getContext('2d')

      this.userChart = new ChartJS(ctx, {
        type: 'line',
        data: {
          datasets: [
            {
              label: 'Created users',
              data: this.userTimeline,
              fill: false,
              borderColor: 'gray',
            },
          ],
        },
        options: {
          layout: {
            padding: {
              left: 10,
              right: 10,
              top: 10,
            },
          },
          legend: {
            display: false,
          },
          scales: {
            xAxes: [
              {
                type: 'time',
                time: {
                  unit: unit,
                  minUnit: unit,
                },
              },
            ],
            yAxes: [
              {
                display: false,
              },
            ],
          },
        },
      })
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

      if (diffInDays > 360) {
        return 'year'
      }

      if (diffInDays > 120) {
        return 'month'
      }

      if (diffInDays > 14) {
        return 'week'
      }

      return 'day'
    },
  },
}
</script>
