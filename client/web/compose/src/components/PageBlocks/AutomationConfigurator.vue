<template>
  <automation-tab
    v-bind="$props"
    :buttons.sync="options.buttons"
  />
</template>
<script>
import base from './base'
import AutomationTab from './Shared/AutomationTab'

export default {
  name: 'Automation',

  components: {
    AutomationTab,
  },

  extends: base,
}
</script>
