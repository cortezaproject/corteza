<template>
  <div class="comment-reply d-flex flex-column gap-1 overflow-hidden border rounded-lg p-2 bg-white pointer">
    <div class="d-flex align-items-center gap-1">
      <div
        :title="authorName"
        class="avatar d-flex align-items-center justify-content-center bg-light font-weight-bold"
      >
        {{ authorInitials }}
      </div>

      <span
        :class="authorIsCurrentUser ? 'text-primary' : 'text-muted'"
        class="text-nowrap font-weight-bold text-truncate"
      >
        {{ authorName }}
      </span>
    </div>

    <div class="d-flex flex-column gap-2 overflow-hidden">
      <field-viewer
        v-if="titleField"
        :field="titleField"
        :record="reply"
        :namespace="namespace"
        value-only
        class="font-weight-bold text-muted h5"
      />

      <field-viewer
        v-if="contentField"
        :field="contentField"
        :record="reply"
        :namespace="namespace"
        value-only
        class="reply-content text-muted text-truncate w-100"
      />
    </div>
  </div>
</template>

<script>
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'

export default {
  name: 'CommentReply',

  components: {
    FieldViewer,
  },

  props: {
    reply: {
      type: Object,
      required: true,
    },

    titleField: {
      type: Object,
      default: undefined,
    },

    contentField: {
      type: Object,
      default: undefined,
    },

    namespace: {
      type: Object,
      required: true,
    },
  },

  computed: {
    authorName () {
      return (((this.reply || {}).author || {}).name) || ''
    },

    authorInitials () {
      return (((this.reply || {}).author || {}).initials) || ''
    },

    authorIsCurrentUser () {
      return Boolean(((this.reply || {}).author || {}).isCurrentUser)
    },
  },
}
</script>

<style lang="scss" scoped>
.avatar {
  width: 2rem;
  height: 2rem;
  font-size: 0.8rem;
  border-radius: 50%;
  user-select: none;
}

.comment-reply {
  border-left-color: var(--primary) !important;
  border-left-width: 3px !important;

  .reply-content {
    transition: color 0.2s ease;
  }

  &:hover {
    .reply-content {
      color: var(--dark) !important;
    }
  }
}
</style>
