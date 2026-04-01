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
          v-if="reactionsField"
          :id="reactionPickerId"
          v-b-tooltip.hover="{ title: 'React', delay: { show: 300, hide: 0 } }"
          variant="extra-light"
          size="sm"
          class="py-1"
        >
          <font-awesome-icon :icon="['far', 'face-smile']" />
        </b-button>
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

      <!-- Reaction Picker Popover -->
      <b-popover
        v-if="reactionsField"
        :target="reactionPickerId"
        triggers="click blur"
        placement="bottom"
        container="body"
        custom-class="reaction-emoji-popover border-light"
        @shown="onReactionPickerShown"
      >
        <c-emoji-picker
          ref="reactionPicker"
          :emojis="emojiData"
          :viewport-height="200"
          :viewport-width="240"
          :labels="{
            search: $t('comment.emojiPicker.search'),
            searchResults: $t('comment.emojiPicker.searchResults'),
            frequentlyUsed: $t('comment.emojiPicker.frequentlyUsed'),
            noResults: $t('comment.emojiPicker.noResults'),
            quickReactions: $t('comment.emojiPicker.quickReactions'),
          }"
          @select="onReactionSelect"
        />
      </b-popover>

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
              emojiPicker: {
                search: $t('content.emojiPicker.search'),
                searchResults: $t('content.emojiPicker.searchResults'),
                frequentlyUsed: $t('content.emojiPicker.frequentlyUsed'),
                noResults: $t('content.emojiPicker.noResults'),
                quickReactions: $t('content.emojiPicker.quickReactions'),
              },
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

          <!-- Reaction badges -->
          <div
            v-if="hasReactions"
            class="comment-reactions d-flex flex-wrap align-items-center gap-1 mt-1"
          >
            <button
              v-for="(userIDs, emoji) in reactions"
              :key="emoji"
              v-b-tooltip.hover="{ title: reactionTooltip(emoji, userIDs), delay: { show: 200, hide: 0 } }"
              type="button"
              class="reaction-badge"
              :class="{ 'reaction-mine': userIDs.includes(currentUserID) }"
              @click.stop="$emit('react', emoji)"
            >
              <span class="reaction-emoji">{{ emoji }}</span>
              <span class="reaction-count">{{ userIDs.length }}</span>
            </button>
          </div>
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

const { CRichTextInput, CEmojiPicker } = components

let reactionPickerCounter = 0

export default {
  name: 'CommentItem',

  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
    CommentReply,
    CRichTextInput,
    CEmojiPicker,
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

    reactionsField: {
      type: Object,
      default: undefined,
    },

    emojiData: {
      type: Array,
      default: () => [],
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

    currentUserID: {
      type: String,
      default: '',
    },

    findUserByID: {
      type: Function,
      default: () => () => undefined,
    },
  },

  data () {
    return {
      isEditing: false,
      editValue: {
        title: '',
        content: '',
      },
      reactionPickerId: `reaction-picker-${++reactionPickerCounter}`,
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

    reactions () {
      if (!this.reactionsField) return {}

      try {
        const val = this.comment.values[this.reactionsField.name]
        return JSON.parse(val || '{}') || {}
      } catch {
        return {}
      }
    },

    hasReactions () {
      return Object.keys(this.reactions).length > 0
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

    reactionTooltip (emoji, userIDs) {
      const names = userIDs.map(id => {
        if (id === this.currentUserID) return 'You'
        const user = this.findUserByID(id)
        return (user || {}).name || (user || {}).handle || (user || {}).email || 'Unknown'
      })

      return names.join(', ')
    },

    onReactionPickerShown () {
      this.$nextTick(() => {
        if (this.$refs.reactionPicker) {
          this.$refs.reactionPicker.reset()
        }
      })
    },

    onReactionSelect (emoji) {
      if (emoji && emoji.emoji) {
        this.$emit('react', emoji.emoji)
      }
      this.$root.$emit('bv::hide::popover', this.reactionPickerId)
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
    position: sticky;
    top: 0;
    align-self: flex-start;
    margin-left: auto;
    opacity: 0;
    transition: opacity 0.2s ease;
    z-index: 1;
    flex-shrink: 0;
    order: 3;
  }

  .comment-card {
    overflow: visible;

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

// Reaction badges
.comment-reactions {
  .reaction-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.2rem;
    padding: 0.1rem 0.4rem;
    border: 1px solid var(--extra-light, #e0e0e0);
    border-radius: 1rem;
    background: var(--white, #fff);
    cursor: pointer;
    font-size: 0.8rem;
    line-height: 1.4;
    transition: background-color 0.15s, border-color 0.15s;

    &:hover {
      background-color: var(--light, #f5f5f5);
      border-color: var(--secondary, #ccc);
    }

    &.reaction-mine {
      background-color: rgba(var(--primary-rgb, 64, 128, 255), 0.08);
      border-color: var(--primary, #4080ff);
    }

    .reaction-emoji {
      font-size: 0.9rem;
    }

    .reaction-count {
      font-size: 0.75rem;
      font-weight: 600;
      color: var(--secondary, #888);
    }

    &.reaction-mine .reaction-count {
      color: var(--primary, #4080ff);
    }
  }
}
</style>

<style lang="scss">
// Unscoped for body-appended popover
.reaction-emoji-popover {
  background: var(--white, #fff);

  .popover-body {
    padding: 0;
  }

  .arrow::after {
    border-bottom-color: var(--white, #fff);
  }
}
</style>
