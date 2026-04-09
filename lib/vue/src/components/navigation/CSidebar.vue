<template>
  <div>
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
      z-index="1037"
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

    <div
      class="d-flex align-items-center justify-content-center tab position-absolute p-2"
    >
      <div
        v-if="disabledRoutes.includes($route.name)"
        class="d-flex align-items-center border-0 p-2"
      >
        <img
          class="icon w-auto border-0"
          :src="icon"
        >
      </div>

      <b-button
        v-else-if="expandOnClick"
        data-test-id="button-sidebar-open"
        variant="outline-extra-light"
        size="lg"
        class="d-flex align-items-center border-0 text-primary"
        @click="openSidebar()"
      >
        <font-awesome-icon
          :icon="['fas', 'bars']"
          class="h4 mb-0"
        />
      </b-button>

      <b-button
        v-else
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

    storageKey: {
      type: String,
      default: 'sidebar-expanded',
    },
  },

  data () {
    return {
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

  mounted () {
    this.checkSidebar(true)
    this.checkIfMobile()

    this.$root.$on('close-sidebar', this.closeSidebar)
    window.addEventListener('resize', this.checkIfMobile)
  },

  beforeUnmount () {
    this.$root.$off('close-sidebar', this.closeSidebar)
    window.removeEventListener('resize', this.checkIfMobile)
  },

  methods: {
    checkIfMobile: throttle(function () {
      this.isMobile = window.innerWidth < 1024
    }, 500),

    checkSidebar (initial = false) {
      // If sidebar should be disabled on route, close and unpin when navigating to route
      if (this.disabledRoutes.includes(this.$route.name)) {
        this.isExpanded = false
      } else if (!this.isMobile && initial) {
        const stored = localStorage.getItem(this.storageKey)
        // Default to closed if no stored value
        this.isExpanded = stored ? stored === 'true' : true
      }
    },

    openSidebar () {
      this.isExpanded = true
      localStorage.setItem(this.storageKey, 'true')
    },

    closeSidebar () {
      this.isExpanded = false
      localStorage.setItem(this.storageKey, 'false')
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
$nav-width-mobile: 400px;

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

// Mobile sidebar should be 400px wide to match notifications and drafts
@media (max-width: 1023px) {
  .sidebar {
    width: $nav-width-mobile !important;
    left: calc(-#{$nav-width-mobile}) !important;
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

  @media (max-width: 1023px) {
    .sidebar {
      width: $nav-width-mobile !important;
      right: calc(-#{$nav-width-mobile}) !important;
    }
  }
}
</style>
