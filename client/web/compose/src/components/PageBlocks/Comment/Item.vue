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
      body-class="comment-card-body d-flex rounded"
      class="comment-card rounded-lg position-relative"
    >
      <div
        v-if="!isEditing"
        class="comment-toolbox d-flex align-items-center justify-content-end bg-light rounded-lg gap"
      >
        <b-button
          v-b-tooltip.hover="{ title: $t('comment.tooltip.reply'), delay: { show: 300, hide: 0 } }"
          variant="extra-light"
          size="sm"
          class="py-1"
          @click.stop="$emit('reply')"
        >
          <font-awesome-icon :icon="['fas', 'reply']" />
        </b-button>
        <b-button
          v-if="canEdit"
          v-b-tooltip.hover="{ title: $t('comment.tooltip.edit'), delay: { show: 300, hide: 0 } }"
          variant="extra-light"
          size="sm"
          class="py-1"
          @click.stop="onEdit"
        >
          <font-awesome-icon :icon="['fas', 'pen']" />
        </b-button>
      </div>

      <div
        :title="commentFullDateTime"
        :class="['comment-time', 'text-nowrap', 'text-muted', 'ml-1', 'overflow-hidden', { 'always-visible': showTimeAlways }]"
      >
        <small>{{ commentTime }}</small>
      </div>

      <div class="d-flex flex-column w-100 overflow-hidden gap">
        <comment-reply
          v-if="comment.reply"
          :reply="comment.reply"
          :title-field="titleField"
          :content-field="contentField"
          :namespace="namespace"
          @click.native="$emit('reply-click', comment.reply.recordID)"
        />

        <div
          v-if="isEditing"
          class="d-flex flex-column"
        >
          <b-form-input
            v-if="titleField"
            v-model="editValue.title"
            class="mb-1"
            :placeholder="$t('comment.title.placeholder')"
          />

          <c-rich-text-input
            v-if="contentField"
            v-model="editValue.content"
            hide-toolbar
            :placeholder="$t('comment.content.placeholder')"
            :labels="{
              urlPlaceholder: $t('content.urlPlaceholder'),
              ok: $t('content.ok'),
            }"
            min-body-height="4rem"
            max-body-height="10rem"
            body-class="overflow-auto"
            style="border: none !important;"
          />

          <div class="d-flex justify-content-end gap-1 my-1">
            <b-button
              variant="light"
              size="sm"
              @click="onCancel"
            >
              {{ $t('general.label.cancel') }}
            </b-button>

            <b-button
              variant="primary"
              size="sm"
              :disabled="isProcessing || !isValid"
              @click="onSave"
            >
              {{ $t('general.label.save') }}
            </b-button>
          </div>
        </div>

        <template v-else-if="!isEditing">
          <field-viewer
            v-if="shouldShowTitle"
            :field="titleField"
            :record="comment"
            :namespace="namespace"
            value-only
            class="font-weight-bold text-muted h5"
          />

          <small
            v-else-if="showTitle && titleField && !titleField.canReadRecordValue"
            class="text-secondary"
          >
            {{ $t('field.noPermission') }}
          </small>

          <field-viewer
            v-if="shouldShowContent"
            :field="contentField"
            :record="comment"
            :namespace="namespace"
            value-only
            class="multiline"
          />
          <small
            v-else-if="showContent && contentField && !contentField.canReadRecordValue"
            class="text-secondary"
          >
            {{ $t('field.noPermission') }}
          </small>

          <field-viewer
            v-if="showAttachments"
            :field="attachmentField"
            :record="comment"
            :namespace="namespace"
            value-only
          />
        </template>
      </div>
    </b-card>
  </div>
</template>

<script>
import { fmt } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import CommentReply from './Reply.vue'

const { CRichTextInput } = components

export default {
  name: 'CommentItem',

  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
    CommentReply,
    CRichTextInput,
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

    attachmentField: {
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

    isProcessing: {
      type: Boolean,
      default: false,
    },
  },

  data () {
    return {
      isEditing: false,
      editValue: {
        title: '',
        content: '',
      },
    }
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

    canEdit () {
      return this.authorIsCurrentUser && !this.comment.deletedAt
    },

    shouldShowTitle () {
      if (!this.showTitle || !this.titleField || !this.titleField.canReadRecordValue) {
        return false
      }

      const v = this.comment.values[this.titleField.name]
      return !!v && v.toString().trim().length > 0
    },

    shouldShowContent () {
      if (!this.showContent || !this.contentField || !this.contentField.canReadRecordValue) {
        return false
      }

      const v = this.comment.values[this.contentField.name]
      return !!v && v.toString().trim().length > 0
    },

    showAttachments () {
      if (!this.attachmentField || !this.attachmentField.canReadRecordValue) {
        return false
      }

      const v = this.comment.values[this.attachmentField.name]

      if (this.attachmentField.isMulti) {
        return Array.isArray(v) && v.length > 0
      }

      return !!v
    },

    isValid () {
      return !!this.editValue.title || !!this.editValue.content
    },
  },

  methods: {
    onEdit () {
      this.isEditing = true
      this.editValue.title = this.titleField ? this.comment.values[this.titleField.name] : ''
      this.editValue.content = this.contentField ? this.comment.values[this.contentField.name] : ''
    },

    onCancel () {
      this.isEditing = false
    },

    onSave () {
      this.$emit('edit', {
        title: this.editValue.title,
        content: this.editValue.content,
      })
      this.isEditing = false
    },
  },
}
</script>

<style lang="scss" scoped>
.avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  user-select: none;
}

.comment-item {
  .comment-time {
    display: block;
    min-width: 3.25rem;
    opacity: 0;
    transition: opacity 0.2s ease;
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
    .comment-card-body {
      padding: 0.2rem 0.25rem;
    }

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
