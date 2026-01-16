<template>
  <div
    class="d-flex w-100"
  >
    <namespace-sidebar
      v-if="namespaces.length"
      :namespaces="namespaces"
    />

    <portal
      to="topbar-avatar-dropdown"
    >
      <b-dropdown-item-button
        data-test-id="dropdown-item-reminders"
        @click="remindersVisible = true"
      >
        {{ $t('reminder.listLabel') }}
      </b-dropdown-item-button>
    </portal>

    <c-reminder-sidebar
      :title="$t('reminder.listLabel')"
      :visible.sync="remindersVisible"
    >
      <reminders />
    </c-reminder-sidebar>

    <router-view
      v-if="loaded"
    />
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import NamespaceSidebar from 'corteza-webapp-compose/src/components/Namespaces/NamespaceSidebar'
import Reminders from 'corteza-webapp-compose/src/components/Namespaces/Reminders'
import { components } from '@cortezaproject/corteza-vue'
const { CReminderSidebar } = components

export default {
  i18nOptions: {
    namespaces: ['general', 'drafts'],
  },

  components: {
    NamespaceSidebar,
    CReminderSidebar,
    Reminders,
  },

  data () {
    return {
      loaded: false,

      query: '',

      remindersVisible: false,
    }
  },

  computed: {
    ...mapGetters({
      namespaces: 'namespace/set',
    }),

    showDrafts () {
      return this.$Settings.get('ui.topbar.showDrafts', false)
    },
  },

  created () {
    // Preload first 500 users
    this.loadUsers({ limit: 500 })

    this.loadNamespaces({ force: true }).finally(() => {
      this.loaded = true
    }).catch(this.toastErrorHandler(this.$t('notification:general.composeAccessNotAllowed')))

    this.$root.$on('reminders.show', () => {
      this.remindersVisible = true
    })
  },

  beforeDestroy () {
    this.setDefaultValues()
  },

  methods: {
    ...mapActions({
      loadUsers: 'user/load',
      loadNamespaces: 'namespace/load',
      toggleDrafts: 'drafts/toggleVisibility',
    }),

    setDefaultValues () {
      this.loaded = false
      this.query = ''
      this.remindersVisible = false
    },
  },
}
</script>

<style scoped>
</style>
