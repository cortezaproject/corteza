<template>
  <div
    class="comment-item"
    :class="{ 'comment-highlighted': highlighted, 'no-hover': disableHover }"
  >
    <div
      v-if="showHeader"
      class="comment-header d-flex align-items-center gap-2 px-2"
    >
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

    <b-card
      body-class="d-flex rounded p-0 p-1"
      class="comment-card rounded-lg position-relative"
    >
      <div class="comment-toolbox bg-light rounded-lg">
        <b-button
          variant="extra-light"
          size="sm"
          class="py-1"
          @click.stop="$emit('reply')"
        >
          <font-awesome-icon :icon="['fas', 'reply']" />
        </b-button>
      </div>

      <div
        :title="commentFullDateTime"
        :class="['comment-time', 'text-nowrap', 'text-muted', 'ml-1', 'overflow-hidden', { 'always-visible': showTimeAlways }]"
      >
        <small>{{ commentTime }}</small>
      </div>

      <div class="d-flex flex-column w-100 overflow-hidden gap-1">
        <comment-reply
          v-if="comment.reply"
          :reply="comment.reply"
          :title-field="titleField"
          :content-field="contentField"
          :namespace="namespace"
          @click.native="$emit('reply-click', comment.reply.recordID)"
        />
        <field-viewer
          v-if="showTitle"
          :field="titleField"
          :record="comment"
          :namespace="namespace"
          value-only
          class="font-weight-bold text-muted h5"
        />
        <small
          v-else-if="titleField && !titleField.canReadRecordValue"
          class="text-secondary"
        >
          {{ $t('field.noPermission') }}
        </small>

        <field-viewer
          v-if="showContent"
          :field="contentField"
          :record="comment"
          :namespace="namespace"
          value-only
          class="multiline"
        />
        <small
          v-else-if="contentField && !contentField.canReadRecordValue"
          class="text-secondary"
        >
          {{ $t('field.noPermission') }}
        </small>
      </div>
    </b-card>
  </div>
</template>

<script>
import { fmt } from '@cortezaproject/corteza-js'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import CommentReply from './Reply.vue'

export default {
  name: 'CommentItem',

  components: {
    FieldViewer,
    CommentReply,
  },

  props: {
    comment: {
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

    showHeader: {
      type: Boolean,
      default: true,
    },

    showTimeAlways: {
      type: Boolean,
      default: false,
    },

    showTitle: {
      type: Boolean,
      default: true,
    },

    showContent: {
      type: Boolean,
      default: true,
    },

    highlighted: {
      type: Boolean,
      default: false,
    },

    disableHover: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    commentTime () {
      return fmt.time((this.comment || {}).updatedAt || (this.comment || {}).createdAt)
    },

    commentFullDateTime () {
      return fmt.fullDateTime((this.comment || {}).updatedAt || (this.comment || {}).createdAt)
    },

    authorName () {
      return ((this.comment || {}).author || {}).name || ''
    },

    authorInitials () {
      return ((this.comment || {}).author || {}).initials || ''
    },

    authorIsCurrentUser () {
      return Boolean(((this.comment || {}).author || {}).isCurrentUser)
    },
  },
}
</script>

<style lang="scss" scoped>
.avatar {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  user-select: none;
}

.comment-item {
  .comment-time {
    display: block;
    width: 3.8rem;
    opacity: 0;
    transition: opacity 0.2s ease;
  }

  .comment-time.always-visible {
    opacity: 1;
  }

  .comment-toolbox {
    position: absolute;
    top: 0;
    right: 0;
    opacity: 0;
    transition: opacity 0.2s ease;
    z-index: 1;
  }

  .comment-card {
    &.comment-highlighted {
      background-color: var(--light);
    }

    transition: background-color 0.2s ease;
  }

  &.comment-highlighted {
    .comment-card {
      background-color: var(--light);
    }
  }

  &:hover {
    .comment-toolbox {
      opacity: 1;
    }
  }

  &:not(.no-hover):hover {
    .comment-card {
      background-color: var(--light);
    }

    .comment-time {
      opacity: 1;
    }
  }
}
</style>
