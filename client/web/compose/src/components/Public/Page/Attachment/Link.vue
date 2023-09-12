<template>
  <a
    v-if="canPreview && attachment.clickToView"
    :href="attachment.url"
    @click.exact.prevent="openLightbox({ ...attachment, ...$event })"
  >
    <slot>
      {{ attachment.name }}
    </slot>
  </a>

  <a
    v-else
  >
    <slot>
      {{ attachment.name }}
    </slot>
  </a>
</template>
<script>
import { components } from '@cortezaproject/corteza-vue'
const { canPreview } = components

export default {
  props: {
    attachment: {
      type: Object,
      required: true,
    },
  },

  computed: {
    canPreview () {
      const meta = this.attachment.meta || {}
      const type = (meta.preview || meta.original || {}).mimetype
      const src = (this.attachment.meta.original && this.attachment.meta.original.ext === 'pdf' ? this.attachment.download : this.attachment.url)
      return canPreview({ type, src, name: this.attachment.name })
    },
  },

  methods: {
    openLightbox (e) {
      this.$root.$emit('showAttachmentsModal', e)
    },
  },
}
</script>
