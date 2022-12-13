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
        :body-class="`bg-white ${isExpanded ? 'py-2 px-3' : ''}`"
        :footer-class="`bg-white rounded-right ${isExpanded ? 'p-2' : ''}`"
        :no-header="!isExpanded"
        :backdrop="isMobile"
        :shadow="isExpanded"
        no-slide
        :right="right"
        no-close-on-route-change
        no-close-on-esc
      >
        <template #header>
          <div
            class="d-flex align-items-center justify-content-between px-2"
            style="height: 50px;"
          >
            <img
              data-test-id="img-main-logo"
              class="logo w-auto border-0"
              :src="logo"
            >

            <b-button
              v-if="isMobile"
              variant="outline-light border-0"
              class="d-flex align-items-center justify-content-center p-2"
              @click="closeSidebar()"
            >
              <font-awesome-icon
                :icon="['fas', 'times']"
                class="h6 mb-0 text-dark"
              />
            </b-button>

            <b-button
              v-else
              data-test-id="button-pin-icon"
              variant="outline-light border-0"
              class="d-flex align-items-center justify-content-center p-2"
              @click="togglePin()"
            >
              <font-awesome-icon
                data-test-id="pin-icon"
                :icon="['fas', 'thumbtack']"
                :class="`h6 mb-0 ${isPinned ? 'text-primary' : 'text-secondary'}`"
              />
            </b-button>
          </div>

          <div
            v-if="!isExpanded"
            class="d-flex align-items-center justify-content-center my-3"
          >
            <b-button
              variant="link"
              @click="togglePin()"
            >
              <font-awesome-icon
                :icon="['fas', 'chevron-right']"
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
        v-if="expandOnHover && !disabledRoutes.includes($route.name)"
        data-test-id="button-sidebar-open"
        variant="outline-light"
        size="lg"
        class="d-flex align-items-center border-0"
        @mouseover="onHover(true)"
      >
        <font-awesome-icon
          :icon="['fas', 'bars']"
          class="h4 mb-0 text-primary"
        />
      </b-button>

      <b-button
        v-else-if="!disabledRoutes.includes($route.name)"
        data-test-id="button-home"
        variant="outline-light"
        size="lg"
        class="d-flex align-items-center p-2 border-0"
        :to="{ name: 'root' }"
      >
        <font-awesome-icon
          :icon="['fas', 'home']"
          class="h4 mb-0 text-primary"
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

    expandOnHover: {
      type: Boolean,
      default: false,
    },

    disabledRoutes: {
      type: Array,
      default: () => [],
    },

    icon: {
      type: String,
      default: () => ''
    },

    logo: {
      type: String,
      default: () => ''
    },

    right: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      sidebar_settings : {}
    }
  },

  computed: {
    isExpanded: {
      get () {
        return this.expanded
      },

      set (expanded) {
        this.$emit('update:expanded', expanded)
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

    isMobile () {
      return window.innerWidth < 576
    },
  },

  created () {
    this.$root.$on('close-sidebar', () => {
      this.isExpanded = false
      this.isPinned = false
    })
  },

  beforeDestroy () {
    this.$root.$off('close-sidebar')
  },

  watch: {
    '$route.name': {
      immediate: true,
      handler (name) {
        // If sidebar should be disabled on route, close and unpin when navigating to route
        if (this.disabledRoutes.includes(name)) {
          this.isPinned = false
          this.isExpanded = false
        } else if(this.expandOnHover){
          this.defaultSidebarAppearance()
        }
      },
    },
  },

  methods: {
    onHover: throttle(function (expand) {
      if (!this.pinned && this.expandOnHover) {
        setTimeout(() => {
          this.isExpanded = expand
        }, expand ? 0 : 100)
      }
    }, 300),

    togglePin () {
      this.saveSettings(!this.isPinned)
      this.isPinned = !this.isPinned
    },

    defaultSidebarAppearance () {
      const localstorage_settings = JSON.parse(window.localStorage.getItem('sidebar_settings'))
      if (localstorage_settings) {
        this.sidebar_settings = localstorage_settings
      }
      const app_sidebar = (localstorage_settings || {})[this.$root.$options.name]
      if (!this.isMobile) {
        if (app_sidebar) {
          this.isExpanded = app_sidebar.pinned
          this.isPinned = app_sidebar.pinned
        } else {
          this.openSidebar()
        }
      } else {
        this.closeSidebar()
      }
    },

    saveSettings (pinned) {
      if (this.sidebar_settings[this.$root.$options.name]) {
        this.sidebar_settings[this.$root.$options.name].pinned = pinned
      } else {
        this.sidebar_settings[this.$root.$options.name] = { pinned: pinned }
      }
      window.localStorage.setItem('sidebar_settings', JSON.stringify(this.sidebar_settings))
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
$sidebar-bg: #F4F7FA;

.sidebar {
  display: flex !important;
  left: calc(-#{$nav-width}) !important;
  -webkit-transition: left 0.15s ease-in-out;
  -moz-transition: left 0.15s ease-in-out;
  -o-transition: left 0.15s ease-in-out;
  transition: left 0.15s ease-in-out;

  header {
    background-color: white;

    &.expanded {
      background-color: $sidebar-bg !important;
    }
  }

  &.expanded {
    left: 0 !important;
    -webkit-transition: left 0.2s ease-in-out;
    -moz-transition: left 0.2s ease-in-out;
    -o-transition: left 0.2s ease-in-out;
    transition: left 0.2s ease-in-out;
  }
}

[dir="rtl"] {
  .sidebar {
    right: calc(-#{$nav-width}) !important;
    -webkit-transition: right 0.15s ease-in-out;
    -moz-transition: right 0.15s ease-in-out;
    -o-transition: right 0.15s ease-in-out;
    transition: right 0.15s ease-in-out;

    header {
      background-color: white;

      &.expanded {
        background-color: $sidebar-bg !important;
      }
    }

    &.expanded {
      right: 0 !important;
      -webkit-transition: right 0.2s ease-in-out;
      -moz-transition: right 0.2s ease-in-out;
      -o-transition: right 0.2s ease-in-out;
      transition: right 0.2s ease-in-out;
    }
  }
}
</style>
