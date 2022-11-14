<template>
  <div>
    <div v-if="mode === 'list'">
      <a :href="attachment.download">
        <font-awesome-icon :icon="['fas', 'download']" />
        {{ attachment.name }}
      </a>
      <i18next
        path="label.attachmentFileInfo"
        tag="label"
      >
        <span>{{ attachment.size }}</span>
        <span>{{ attachment.changedAt }}</span>
      </i18next>
    </div>

    <div v-if="mode === 'grid'">
      <a :href="aattachment.download">
        <font-awesome-icon
          :icon="['far', 'file-'+ext(a)]"
          :title="$t('label.openBookmarks')"
        />
        {{ aattachment.name }}
      </a>
      <i18next
        path="label.attachmentFileInfo"
        tag="label"
      >
        <span>{{ attachment.size }}</span>
        <span>{{ attachment.changedAt }}</span>
      </i18next>
    </div>

    <div
      v-if="mode === 'single' || 'gallery'"
      class="single"
    >
      <div v-if="isImage(a)">
        <img
          :src="attachment.previewUrl"
          @click="openLightbox(index)"
        >
      </div>
      <div v-else>
        <font-awesome-icon
          :icon="['far', 'file-'+ext(a)]"
          :title="$t('label.openBookmarks')"
        />
        <a :href="attachment.download">
          {{ $t('label.download') }}
        </a>
      </div>
      {{ a.name }}
    </div>
  </div>
</template>
<script>

export default {
  i18nOptions: {
    namespaces: 'general',
  },

  props: {
    kind: {
      type: String,
      required: true,
    },

    mode: {
      type: String,
      required: true,
    },

    value: {
      type: [Object, String],
      required: true,
    },
  },

  data () {
    return {
      attachment: {},
    }
  },

  watch: {
    value: {
      immediate: true,
      handler (value) {
        // On input change resolve/load all attachments
        if (typeof value === 'string') {
          this.$ComposeAPI.attachmentRead({ kind: this.kind, attachmentID: value }).then(a => {
            this.attachment = a
          })
        } else {
          this.attachment = value
        }
      },
    },
  },
}
</script>
