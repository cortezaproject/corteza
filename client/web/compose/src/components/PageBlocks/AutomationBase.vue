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

      abortableRequests: [],
    }
  },

  computed: {
    hasUIHooks () {
      return this.$UIHooks.set && !!this.$UIHooks.set.length
    },
  },

  beforeDestroy () {
    this.abortRequests()
    this.setDefaultValues()
  },

  created () {
    if (this.$UIHooks.set && !!this.$UIHooks.set.length) {
      return
    }

    this.fetchAutomationLists()
  },

  methods: {
    fetchAutomationLists () {
      this.processing = true

      const { response, cancel } = this.$ComposeAPI
        .automationListCancellable({ eventTypes: ['onManual'], excludeInvalid: true })

      this.abortableRequests.push(cancel)

      return response()
        .then(({ set = [] }) => {
          this.automationScripts = set
        })
        .finally(() => {
          this.processing = false
        })
    },

    setDefaultValues () {
      this.processing = false
      this.automationScripts = []
      this.abortableRequests = []
    },

    abortRequests () {
      this.abortableRequests.forEach((cancel) => {
        cancel()
      })
    },
  },
}
</script>
