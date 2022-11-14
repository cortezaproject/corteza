import Vue from 'vue'
import { system } from '@cortezaproject/corteza-js'

// Fixes event with the used resourceType and extend with arguments
function fix (ev, resourceType, args) {
  return {
    ...ev,
    // Override with the resource type from the trigger
    resourceType,
    args,
  }
}

Vue.mixin({
  methods: {
    dispatchCortezaSystemEvent ({ script, resourceType }, args = {}) {
      this.$EventBus.Dispatch(fix(system.UserEvent(), resourceType, args), script)
    },

    dispatchCortezaSystemUserEvent ({ script, resourceType }, args) {
      this.$EventBus.Dispatch(fix(system.UserEvent(args.user), resourceType, args), script)
    },

    dispatchCortezaSystemRoleEvent ({ script, resourceType }, args) {
      this.$EventBus.Dispatch(fix(system.RoleEvent(args.role), resourceType, args), script)
    },
  },
})
