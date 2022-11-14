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
      <b-dropdown-item
        data-test-id="dropdown-item-reminders"
        @click="remindersVisible = true"
      >
        {{ $t('reminder.listLabel') }}
      </b-dropdown-item>
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
import NamespaceSidebar from 'corteza-webapp-compose/src/components/Namespaces/NamespaceSidebar'
import Reminders from 'corteza-webapp-compose/src/components/Namespaces/Reminders'
import { components } from '@cortezaproject/corteza-vue'
const { CReminderSidebar } = components

export default {
  i18nOptions: {
    namespaces: 'general',
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
      namespaces: [],

      remindersVisible: false,
    }
  },

  created () {
    this.$store.dispatch('namespace/load', { force: true }).then(namespaces => {
      this.namespaces = namespaces
      this.loaded = true
    }).catch(this.toastErrorHandler(this.$t('notification:general.composeAccessNotAllowed')))

    this.$root.$on('reminders.show', () => {
      this.remindersVisible = true
    })
  },
}
</script>
