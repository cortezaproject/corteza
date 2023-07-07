<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <list-loader
      kind="page"
      :set="options.attachments"
      :namespace="namespace"
      :mode="options.mode"
      :hide-file-name="options.hideFileName"
      :preview-options="options"
      class="h-100 p-2"
    />
  </wrap>
</template>
<script>
import base from './base'
import ListLoader from '../Public/Page/Attachment/ListLoader'

export default {
  components: {
    ListLoader,
  },

  extends: base,
}
</script>
