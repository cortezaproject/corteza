<template>
  <div>
    <div
      @mouseleave="onHover(false)"
    >
      <b-sidebar
        v-model="isExpanded"
        data-test-id="sidebar"
        :sidebar-class="`sidebar ${isExpanded ? 'expanded' : ''}`"
        :header-class="`d-block sidebar-header ${isExpanded ? 'expanded border-bottom p-2' : ''}`"
        :body-class="`${isExpanded ? 'px-3' : ''}`"
        :footer-class="`rounded-right ${isExpanded ? 'px-2' : ''}`"
        :no-header="!isExpanded"
        :backdrop="isMobile"
        backdrop-variant="white"
        :shadow="isExpanded && 'sm'"
        no-slide
        :right="right"
        no-close-on-route-change
        no-close-on-esc
      >
        <template #header>
          <div
            class="d-flex align-items-center justify-content-between pl-2"
            style="height: 47px;"
          >
            <img
              data-test-id="img-main-logo"
              class="logo w-auto border-0"
              :src="logo"
            >

            <b-button
              variant="outline-light"
              class="d-flex align-items-center justify-content-center p-2 border-0 text-secondary"
              @click="closeSidebar()"
            >
              <font-awesome-icon
                :icon="['fas', 'times']"
                class="h6 mb-0"
              />
            </b-button>
          </div>

          <div
            v-if="isExpanded"
            class="px-2"
          >
            <slot
              name="header-expanded"
            />
          </div>

          <hr
            v-if="!isExpanded"
            class="my-2"
          >
        </template>

        <slot
          v-if="isExpanded"
          name="body-expanded"
        />

        <template #footer>
          <slot
            v-if="isExpanded"
            name="footer-expanded"
          />
        </template>
      </b-sidebar>
    </div>

    <div
      class="d-flex align-items-center justify-content-center tab position-absolute p-2"
    >
      <b-button
        v-if="expandOnClick && !disabledRoutes.includes($route.name)"
        data-test-id="button-sidebar-open"
        variant="outline-extra-light"
        size="lg"
        class="d-flex align-items-center border-0 text-primary"
        @click="togglePin()"
      >
        <font-awesome-icon
          :icon="['fas', 'bars']"
          class="h4 mb-0"
        />
      </b-button>

      <b-button
        v-else-if="!disabledRoutes.includes($route.name)"
        data-test-id="button-home"
        variant="outline-extra-light"
        size="lg"
        class="d-flex align-items-center p-2 border-0 text-primary"
        :to="{ name: 'root' }"
      >
        <font-awesome-icon
          :icon="['fas', 'home']"
          class="h4 mb-0"
        />
      </b-button>

      <div
        v-else
        class="d-flex align-items-center border-0 p-2"
      >
        <img
          class="icon w-auto border-0"
          :src="icon"
        >
      </div>
    </div>
  </div>
</template>

<script>
import { throttle } from 'lodash'

export default {
  props: {
    expanded: {
      type: Boolean,
      default: false,
    },

    pinned: {
      type: Boolean,
      default: false,
    },

    expandOnClick: {
      type: Boolean,
      default: false,
    },

    disabledRoutes: {
      type: Array,
      default: () => [],
    },

    icon: {
      type: String,
      default: () => '',
    },

    logo: {
      type: String,
      default: () => '',
    },

    right: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      sidebarSettings: {},
      isMobile: false,
    }
  },

  computed: {
    isExpanded: {
      get () {
        return this.expanded
      },

      set (expanded) {
        this.$emit('update:expanded', expanded)

        if (!expanded) {
          this.isPinned = false
        }
      },
    },

    isPinned: {
      get () {
        return this.pinned
      },

      set (pinned) {
        this.$emit('update:pinned', pinned)
      },
    },
  },

  watch: {
    '$route.name': {
      handler () {
        this.checkSidebar()
      },
    },

    disabledRoutes: {
      handler () {
        this.checkSidebar()
      },
    },
  },

  created () {
    this.checkSidebar()

    this.$root.$on('close-sidebar', this.closeSidebar)
    window.addEventListener('resize', this.checkIfMobile)
  },

  beforeUnmount () {
    this.$root.$off('close-sidebar', this.closeSidebar)
    window.removeEventListener('resize', this.checkIfMobile)
  },

  methods: {
    checkIfMobile: throttle(function () {
      this.isMobile = window.innerWidth < 1024 || /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)

      if (this.isMobile) {
        this.closeSidebar()
      }
    }, 500),

    checkSidebar () {
      // If sidebar should be disabled on route, close and unpin when navigating to route
      if (this.disabledRoutes.includes(this.$route.name)) {
        this.isPinned = false
        this.isExpanded = false
      } else if (this.expandOnClick && !this.isExpanded) {
        this.defaultSidebarAppearance()
      }

      this.checkIfMobile()
    },

    onHover: throttle(function (expand) {
      if (!expand && !this.pinned && this.expandOnClick) {
        setTimeout(() => {
          this.isExpanded = expand
        }, expand ? 0 : 100)
      }
    }, 300),

    togglePin () {
      if (this.expandOnClick && !this.isExpanded) {
        this.isExpanded = true
      }

      if (!this.isMobile) {
        this.isPinned = !this.isPinned
        this.saveSettings(this.isPinned)
      }
    },

    defaultSidebarAppearance () {
      const localStorageSettings = JSON.parse(window.localStorage.getItem('sidebarSettings'))

      if (localStorageSettings) {
        this.sidebarSettings = localStorageSettings
      }

      const appSidebar = (localStorageSettings || {})[this.$root.$options.name]

      if (!this.isMobile) {
        if (appSidebar) {
          this.isExpanded = appSidebar.pinned
          this.isPinned = appSidebar.pinned
        } else {
          this.openSidebar()
        }
      } else {
        this.closeSidebar()
      }
    },

    saveSettings (pinned) {
      if (this.sidebarSettings[this.$root.$options.name]) {
        this.sidebarSettings[this.$root.$options.name].pinned = pinned
      } else {
        this.sidebarSettings[this.$root.$options.name] = { pinned }
      }
      window.localStorage.setItem('sidebarSettings', JSON.stringify(this.sidebarSettings))
    },

    openSidebar () {
      this.isPinned = true
      this.isExpanded = true
    },

    closeSidebar () {
      this.isPinned = false
      this.isExpanded = false
    },
  },
}
</script>

<style lang="scss" scoped>
$header-height: 64px;

.tab {
  z-index: 1021;
  top: 0;
  height: $header-height;
  width: 66px;
}

.icon {
  max-height: 40px;
  max-width: 40px;
}

.logo {
  max-height: 40px;
}

.sidebar-header {
  height: $header-height;
}
</style>

<style lang="scss">
$nav-width: 320px;

.b-sidebar {
  background-color: var(--white) !important;
}

.b-sidebar-backdrop {
  opacity: 0.75 !important;
}

.sidebar {
  display: flex !important;
  left: calc(-#{$nav-width}) !important;
  transition: left 0.2s cubic-bezier(0.4, 0, 0.2, 1);

  &.expanded {
    left: 0 !important;
    transition: left 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }
}

[dir="rtl"] {
  .sidebar {
    right: calc(-#{$nav-width}) !important;
    left: auto !important;
    transition: right 0.2s cubic-bezier(0.4, 0, 0.2, 1);

    &.expanded {
      right: 0 !important;
      transition: right 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    }
  }
}
</style>
