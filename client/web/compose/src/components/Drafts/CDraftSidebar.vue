<template>
  <b-sidebar
    v-model="isVisible"
    header-class="d-flex align-items-center justify-content-between bg-white pr-2 pl-3 py-3 border-bottom"
    body-class="d-flex flex-column overflow-hidden bg-white"
    sidebar-class="draft-sidebar"
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
        <b>{{ $t('title') }}</b>
      </h5>

      <b-button
        v-b-tooltip.hover="{ title: $t('general:label.close'), delay: { show: 500, hide: 0 } }"
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

    <drafts />
  </b-sidebar>
</template>

<script>
import Drafts from 'corteza-webapp-compose/src/components/Drafts/Drafts.vue'
import { mapGetters, mapMutations } from 'vuex'
import { throttle } from 'lodash'

export default {
  i18nOptions: {
    namespaces: 'drafts',
  },

  components: {
    Drafts,
  },

  data () {
    return {
      isMobile: false,
    }
  },

  computed: {
    ...mapGetters({
      visible: 'drafts/visible',
    }),

    isVisible: {
      get () {
        return this.visible
      },

      set (visible) {
        this.setVisible(visible)
      },
    },
  },

  watch: {
    isVisible (visible) {
      if (visible) {
        this.$root.$emit('right-sidebar:opened', 'drafts')
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
    ...mapMutations({
      setVisible: 'drafts/setVisible',
    }),

    checkIfMobile: throttle(function () {
      this.isMobile = window.innerWidth < 1024
    }, 500),

    handleSidebarOpened (name) {
      if (name !== 'drafts') {
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
  .b-sidebar.draft-sidebar {
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
