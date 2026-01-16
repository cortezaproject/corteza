<template>
  <b-sidebar
    v-model="isVisible"
    :title="title"
    header-class="d-flex align-items-center justify-content-between bg-white pr-2 pl-3 py-3 border-bottom"
    body-class="d-flex flex-column overflow-hidden bg-white"
    sidebar-class="reminder-sidebar"
    bg-variant="white"
    :backdrop="isMobile"
    backdrop-variant="white"
    no-footer
    right
    shadow
    no-close-on-route-change
    no-close-on-esc
    width="400px"
  >
    <template #header>
      <h5
        class="text-primary mb-0"
      >
        <b>{{ title }}</b>
      </h5>

      <b-button
        variant="outline-light"
        class="d-flex align-items-center justify-content-center p-2 border-0 text-secondary"
        @click="isVisible = false"
      >
        <font-awesome-icon
          :icon="['fas', 'times']"
          class="h6 mb-0"
        />
      </b-button>
    </template>

    <slot />
  </b-sidebar>
</template>

<script lang="js">
import { throttle } from 'lodash'

export default {
  props: {
    title: {
      type: String,
      default: '',
    },

    visible: {
      type: Boolean,
      required: false,
      default: false,
    },
  },

  data () {
    return {
      isMobile: false,
    }
  },

  computed: {
    isVisible: {
      get () {
        return this.visible
      },

      set (visible) {
        this.$emit('update:visible', visible)
      },
    },
  },

  watch: {
    isVisible (visible) {
      if (visible) {
        this.$root.$emit('right-sidebar:opened', 'reminders')
      }
    },
  },

  created () {
    this.$root.$on('right-sidebar:opened', this.handleSidebarOpened)
  },

  mounted () {
    this.checkIfMobile()
    window.addEventListener('resize', this.checkIfMobile)
  },

  beforeDestroy () {
    this.$root.$off('right-sidebar:opened', this.handleSidebarOpened)
    window.removeEventListener('resize', this.checkIfMobile)
  },

  methods: {
    checkIfMobile: throttle(function () {
      this.isMobile = window.innerWidth < 1024
    }, 500),

    handleSidebarOpened (name) {
      if (name !== 'reminders') {
        this.isVisible = false
      }
    },
  },
}
</script>

<style lang="scss">
.b-sidebar-backdrop {
  opacity: 0.75 !important;
}

@media (min-width: 1024px) {
  .b-sidebar.reminder-sidebar {
    top: calc(var(--topbar-height) + 0.5rem) !important;
    right: 0.5rem !important;
    height: calc(100% - var(--topbar-height) - 1rem) !important;
    border-radius: 1rem !important;
    border: none !important;
    z-index: 1048 !important;

    .b-sidebar-header {
      border-top-left-radius: 1rem !important;
      border-top-right-radius: 1rem !important;
    }

    .b-sidebar-body {
      border-bottom-left-radius: 1rem !important;
      border-bottom-right-radius: 1rem !important;
    }
  }
}
</style>
