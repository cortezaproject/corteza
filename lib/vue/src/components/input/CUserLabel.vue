<template>
  <span :class="stateClass">{{ label }}</span>
</template>

<script>
import { NoID } from '@cortezaproject/corteza-js'

export default {
  name: 'CUserLabel',

  props: {
    user: {
      type: Object,
      default: () => ({}),
    },
  },

  computed: {
    label () {
      const { userID, name, username, email, handle } = this.user || {}
      return name || email || username || handle || (userID && userID !== NoID ? `<@${userID}>` : '')
    },

    stateClass () {
      const { deletedAt, suspendedAt } = this.user || {}
      if (deletedAt) return 'text-danger'
      if (suspendedAt) return 'text-muted'
      return ''
    },
  },
}
</script>
