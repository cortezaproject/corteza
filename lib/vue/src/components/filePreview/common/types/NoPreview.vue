<template>
  <div class="d-flex align-items-center justify-content-center h-100 mb-2">
    <font-awesome-icon
      v-if="inline"
      :title="name"
      :icon="['far', `file-${icon}`]"
      :style="previewStyle"
      class="inline-icon d-block text-secondary"
      @click.stop="$emit('openPreview')"
    />

    <p
      v-else
      class="text-secondary"
    >
      {{ labels.previewUnavailable || 'Preview unavailable' }}
    </p>
  </div>
</template>

<script lang="js">
import { getExtensionIconType } from '../index.js'
import base from '../base.vue'

export default {
  extends: base,

  props: {
    previewStyle: {
      type: Object,
      default: () => ({}),
    },

    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  computed: {
    icon () {
      const { original = {} } = this.meta || {}
      const { ext } = original || {}
      return getExtensionIconType(ext)
    },
  },
}
</script>

<style lang="scss" scoped>
.inline-icon {
  cursor: zoom-in;
}
</style>
