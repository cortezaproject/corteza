<template>
  <div
    class="app-selector d-flex flex-column h-100 mt-3"
  >
    <div class="d-flex justify-content-center align-items-center">
      <b-img
        :src="logo"
        class="logo px-3"
      />
    </div>

    <div class="search w-100 mx-auto my-3 px-5">
      <div class="flex-grow-1 mt-1">
        <c-input-search
          v-model.trim="query"
          data-v-onboarding="app-list"
          :aria-label="$t('search')"
          :placeholder="$t('search')"
          :debounce="200"
        />
      </div>
    </div>

    <div
      class="flex-fill overflow-auto"
    >
      <b-container
        class="h-100"
      >
        <draggable
          v-if="filteredApps.length"
          v-model="appList"
          :disabled="!canCreateApplication || query || isMobileResolution"
          group="apps"
          class="h-100 w-100"
          @end="onDrop"
        >
          <transition-group
            name="apps"
            tag="b-row"
            class="d-flex flex-wrap align-items-stretch justify-content-center mx-2"
          >
            <b-col
              v-for="app in filteredApps"
              :key="app.applicationID"
              cols="12"
              sm="6"
              md="6"
              lg="4"
              xl="3"
              class="p-0 mb-3 mt-1"
              :data-v-onboarding="getStepName(app.unify.url)"
            >
              <b-card
                no-body
                overlay
                class="app h-100"
                @mouseover="hovered = app.applicationID"
                @mouseleave="hovered = undefined"
              >
                <div
                  class="align-content-center d-flex flex-grow-1 flex-wrap"
                >
                  <b-card-img
                    class="rounded-bottom thumbnail"
                    :src="logoUrl(app)"
                    :alt="app.unify.name || app.name"
                  />
                </div>

                <b-card-text
                  class="text-center my-4 h6"
                >
                  {{ app.unify.name || app.name }}
                </b-card-text>

                <b-link
                  :data-test-id="app.name"
                  :disabled="!app.enabled"
                  :href="app.unify.url"
                  :target="openAppInNewTab(app.unify.url)"
                  :style="[{ cursor: `${app.enabled ? 'pointer': canCreateApplication ? 'grab' : 'default'}` }]"
                  class="stretched-link"
                />
              </b-card>
            </b-col>
          </transition-group>
        </draggable>

        <div
          v-else
          class="d-flex justify-content-center w-100"
        >
          <h4
            data-test-id="heading-no-apps"
            class="mt-5"
          >
            {{ $t('no-applications') }}
          </h4>
        </div>
      </b-container>
    </div>

    <portal
      to="topbar-help-dropdown"
    >
      <b-dropdown-item
        data-test-id="dropdown-helper-tour"
        @click="$refs.tour.onStartClick()"
      >
        {{ $t('start-tour') }}
      </b-dropdown-item>
    </portal>

    <tour-start
      @start="startTour"
    />

    <tour
      ref="tour"
      name="app-list"
      :steps="filteredSteps"
    />
  </div>
</template>
<script>
import { mapGetters, mapActions } from 'vuex'
import Draggable from 'vuedraggable'
import { components, url } from '@cortezaproject/corteza-vue'
const { Tour, TourStart, CInputSearch } = components

export default {
  i18nOptions: {
    namespaces: 'layout',
  },

  components: {
    Draggable,
    Tour,
    TourStart,
    CInputSearch,
  },

  props: {
    logo: {
      type: String,
      default: () => '',
    },
  },

  data () {
    return {
      query: '',

      appList: [],

      canCreateApplication: false,
      canPin: false,

      hovered: undefined,

      isMobileResolution: false,

      steps: [
        { name: 'app-list', dynamic: false },
        { name: 'low-code', dynamic: true },
        { name: 'crm', dynamic: true },
        { name: 'reporter', dynamic: true },
        { name: 'workflow', dynamic: true },
        { name: 'profile', dynamic: false },
      ],

    }
  },

  computed: {
    ...mapGetters({
      apps: 'applications/unifyOnly',
    }),

    filteredApps () {
      const query = (this.query || '').toUpperCase()
      return this.query
        ? this.appList.filter(({ name }) => (name.toUpperCase()).includes(query))
        : this.appList
    },

    filteredSteps () {
      return this.steps.filter(step => {
        if (step.dynamic) {
          return this.filteredApps.some(app => {
            return this.getStepName(app.unify.url) === step.name
          })
        }
        return true
      }).map(s => { return s.name })
    },
  },

  watch: {
    'apps': {
      immediate: true,
      handler (apps) {
        this.appList = apps
      },
    },
  },

  created () {
    this.fetchEffective()
    if (window.innerWidth < 576) {
      this.isMobileResolution = true
    }
  },

  methods: {
    ...mapActions({
      reorderApp: 'applications/reorder',
      pinApp: 'applications/pin',
      unpinApp: 'applications/unpin',
    }),

    getStepName (url) {
      switch (url) {
        case 'compose/':
          return 'low-code'
        case 'compose/ns/crm/pages':
          return 'crm'
        case 'reporter/':
          return 'reporter'
        case 'workflow/':
          return 'workflow'
      }
    },

    fetchEffective () {
      this.$SystemAPI.permissionsEffective({ resource: 'application' })
        .then(p => {
          this.canCreateApplication = p.find(per => per.operation === 'application.create').allow || false
          // this.canPin = p.find(({ resource, operation, allow }) => resource === 'system' && operation === 'application.flag.self').allow
        })
    },

    handlePin (pin = true, applicationID) {
      if (pin) {
        this.unpinApp({ applicationID, ownedBy: this.$auth.user.userID })
      } else {
        this.pinApp({ applicationID, ownedBy: this.$auth.user.userID })
      }
    },

    async onDrop () {
      const applicationIDs = this.appList.map(({ applicationID }) => applicationID)
      await this.reorderApp(applicationIDs)
    },

    startTour () {
      this.$refs.tour.onStart()
    },

    logoUrl (app) {
      if (!app.unify.logo) {
        return 'applications/default-app.png'
      }

      const apiSystem = '/api/system'
      const apiBaseUrl = (new URL(url.Make({ url: this.$SystemAPI.baseURL }))).toString()

      // Properly handle uploaded logos
      // but cut away only /api/system (without any potential base-url prefix)
      if (app.unify.logo.startsWith(apiSystem)) {
        // remove path from the URL
        return apiBaseUrl.substring(0, apiBaseUrl.length - apiSystem.length) + app.unify.logo
      }

      // Provisioned app logos
      return app.unify.logo
    },

    openAppInNewTab (route) {
      return !route.includes('jitsi') ? '' : '_blank'
    },
  },
}
</script>
<style lang="scss" scoped>
.app-selector {
  .logo {
    max-height: 20vh;
    max-width: 500px;
    width: auto;
  }

  @media only screen and (max-width: 576px) {
    .logo {
      max-width: 100%;
    }
  }

  .search {
    max-width: 600px;
  }

  .app {
    min-height: 13rem;
    transition: all 0.2s ease;
    box-shadow: 0;
    top: 0;
    margin: 0 0.625rem;

    .thumbnail {
      max-width: 100%;
      max-height: 150px;
      object-fit: contain;
    }

    &:hover {
      transition: all 0.2s ease;
      box-shadow: 0px 4px 8px rgba(38, 38, 38, 0.2);
      top: -2px;
    }
  }

  .star {
    position: absolute;
    top: .2rem;
    left: .2rem;
    padding: 0;
    margin: 0;
    background-color: transparent;
    border: none;
    .star-icon {
      fill: var(--warning);
      width: 1.2rem;
      height: 1.2rem;
    }
  }

  .apps-leave-active {
    position: absolute;
    transition: opacity 0.25s ease;
  }
  .apps-enter, .apps-leave-to {
    opacity: 0;
  }

  .apps-move {
    transition: transform 0.25s ease;
  }
}
</style>
