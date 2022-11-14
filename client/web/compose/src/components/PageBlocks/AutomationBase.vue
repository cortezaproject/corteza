<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <div
      v-if="processing"
      class="d-flex align-items-center justify-content-center h-100"
    >
      <b-spinner />
    </div>

    <automation-buttons
      v-else
      class="d-flex flex-wrap my-2"
      button-class="my-1 mx-3 flex-fill"
      :buttons="options.buttons"
      :automation-scripts="automationScripts"
      v-bind="$props"
    />
  </wrap>
</template>
<script>
import base from './base'
import AutomationButtons from './Shared/AutomationButtons'

export default {
  components: { AutomationButtons },
  extends: base,

  props: {
    extraEventArgs: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      processing: false,

      automationScripts: [],
    }
  },

  computed: {
    hasUIHooks () {
      return this.$UIHooks.set && !!this.$UIHooks.set.length
    },
  },

  created () {
    if (this.$UIHooks.set && !!this.$UIHooks.set.length) {
      return
    }

    this.processing = true
    return this.$ComposeAPI.automationList({ eventTypes: ['onManual'], excludeInvalid: true })
      .then(({ set = [] }) => {
        this.automationScripts = set
      })
      .finally(() => {
        this.processing = false
      })
  },
}
</script>
