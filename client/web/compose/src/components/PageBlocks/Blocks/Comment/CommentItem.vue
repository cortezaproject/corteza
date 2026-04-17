<template>
  <div
    class="comment-item"
    :class="{ 'comment-highlighted': highlighted, 'no-hover': disableHover }"
    @mouseleave="$emit('mouseleave')"
  >
    <!-- Author header (shown only for first message in a group) -->
    <div v-if="showHeader" class="flex items-center gap-5 px-2">
      <div :title="authorName" class="avatar flex items-center justify-center font-semibold">
        {{ authorInitials }}
      </div>

      <span
        :class="authorIsCurrentUser ? 'text-primary' : 'text-muted-color'"
        class="whitespace-nowrap font-semibold truncate"
      >
        {{ authorName }}
      </span>
    </div>

    <!-- Comment card -->
    <div class="comment-card relative rounded-lg">
      <div class="comment-card-body flex rounded">
        <!-- Timestamp -->
        <div
          :title="commentFullDateTime"
          :class="{ 'always-visible': showTimeAlways }"
          class="comment-time text-muted-color whitespace-nowrap overflow-hidden mr-1"
        >
          <small>{{ commentTime }}</small>
        </div>

        <div class="flex flex-col w-full overflow-hidden gap-1">
          <CommentReply
            v-if="comment.reply"
            :reply="comment.reply"
            :title-field="titleField"
            :content-field="contentField"
            @click="$emit('reply-click', comment.reply.recordID)"
          />

          <div v-if="isEditing" class="flex flex-col gap-1">
            <InputText
              v-if="titleField"
              v-model="editTitle"
              :placeholder="$t('block.comment.title.placeholder')"
              class="w-full"
            />

            <CRichTextInput
              v-model="editContent"
              :placeholder="$t('block.comment.content.placeholder')"
              hide-toolbar
              min-body-height="4rem"
              max-body-height="10rem"
              body-class="overflow-auto"
            />

            <div class="flex justify-end gap-1">
              <Button
                :label="$t('general.label.cancel')"
                size="small"
                severity="secondary"
                @click="onCancel"
              />
              <Button
                :label="$t('general.label.save')"
                size="small"
                :disabled="!isValid"
                @click="onSave"
              />
            </div>
          </div>

          <template v-else>
            <!-- Title with RBAC check -->
            <div v-if="shouldShowTitle" class="font-semibold text-muted-color text-lg">
              <CFieldViewer :field="titleField" :record="comment" :namespace="namespace" />
            </div>
            <small
              v-else-if="showTitle && titleField && !titleField.canReadRecordValue"
              class="text-muted-color"
            >
              {{ $t('block.field.noPermission') }}
            </small>

            <!-- Content with RBAC check -->
            <div v-if="shouldShowContent" class="comment-content">
              <CFieldViewer :field="contentField" :record="comment" :namespace="namespace" />
            </div>
            <small
              v-else-if="showContent && contentField && !contentField.canReadRecordValue"
              class="text-muted-color"
            >
              {{ $t('block.field.noPermission') }}
            </small>

            <!-- Attachment display -->
            <CFieldViewer
              v-if="showAttachments"
              :field="attachmentField"
              :record="comment"
              :namespace="namespace"
            />

            <!-- Reaction badges -->
            <div
              v-if="hasReactions"
              class="comment-reactions flex flex-wrap items-center gap-1 mt-1"
            >
              <button
                v-for="(userIDs, emoji) in reactions"
                :key="emoji"
                v-tooltip.top="{ value: reactionTooltip(emoji, userIDs), showDelay: 200 }"
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

        <!-- Hover toolbox -->
        <div v-if="!isEditing" class="comment-toolbox flex items-center justify-end gap-1">
          <Button
            v-if="reactionsField"
            v-tooltip.top="{ value: $t('block.comment.tooltip.react'), showDelay: 300 }"
            icon="pi pi-face-smile"
            text
            size="small"
            severity="secondary"
            @click.stop="toggleEmojiPicker"
          />
          <Button
            v-tooltip.top="{ value: $t('block.comment.tooltip.reply'), showDelay: 300 }"
            icon="pi pi-reply"
            text
            size="small"
            severity="secondary"
            @click.stop="$emit('reply')"
          />
          <Button
            v-if="canEdit"
            v-tooltip.top="{ value: $t('block.comment.tooltip.edit'), showDelay: 300 }"
            icon="pi pi-pencil"
            text
            size="small"
            severity="secondary"
            @click.stop="onEdit"
          />
        </div>
      </div>

      <!-- Emoji picker popover -->
      <Popover ref="emojiPopoverRef" @show="onEmojiPopoverShow">
        <CEmojiPicker
          ref="emojiPickerComp"
          :emojis="allEmojis"
          :show-quick-reactions="true"
          @select="selectEmoji"
        />
      </Popover>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import CommentReply from './CommentReply.vue'
import { components, CEmojiPicker, emojiData } from '@planetcrust/human-vue'

const { CRichTextInput, CFieldViewer } = components

const props = defineProps({
  comment: { type: Object, required: true },
  titleField: { type: Object, default: undefined },
  contentField: { type: Object, default: undefined },
  attachmentField: { type: Object, default: undefined },
  reactionsField: { type: Object, default: undefined },
  namespace: { type: Object, required: true },
  showHeader: { type: Boolean, default: true },
  showTimeAlways: { type: Boolean, default: false },
  showTitle: { type: Boolean, default: true },
  showContent: { type: Boolean, default: true },
  highlighted: { type: Boolean, default: false },
  disableHover: { type: Boolean, default: false },
  currentUserID: { type: String, default: '' },
  findUserByID: { type: Function, default: () => () => undefined },
})

const emit = defineEmits(['reply', 'edit', 'react', 'reply-click', 'mouseleave'])

const isEditing = ref(false)
const editTitle = ref('')
const editContent = ref('')
const emojiPopoverRef = ref(null)
const emojiPickerComp = ref(null)

const allEmojis = computed(() => emojiData || [])

const authorName = computed(() => props.comment?.author?.name || '')
const authorInitials = computed(() => props.comment?.author?.initials || '?')
const authorIsCurrentUser = computed(() => Boolean(props.comment?.author?.isCurrentUser))

const canEdit = computed(() => authorIsCurrentUser.value && !props.comment?.deletedAt)

const commentTime = computed(() => {
  const dt = props.comment?.updatedAt || props.comment?.createdAt
  if (!dt) return ''
  try {
    return new Date(dt).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })
  } catch {
    return ''
  }
})

const commentFullDateTime = computed(() => {
  const dt = props.comment?.updatedAt || props.comment?.createdAt
  if (!dt) return ''
  try {
    return new Date(dt).toLocaleString()
  } catch {
    return ''
  }
})

const isValid = computed(() => !!editTitle.value || !!editContent.value)

// RBAC-aware display checks
const shouldShowTitle = computed(() => {
  if (!props.showTitle || !props.titleField || !props.titleField.canReadRecordValue) return false
  const v = props.comment?.values?.[props.titleField.name]
  return !!v && v.toString().trim().length > 0
})

const shouldShowContent = computed(() => {
  if (!props.showContent || !props.contentField || !props.contentField.canReadRecordValue)
    return false
  const v = props.comment?.values?.[props.contentField.name]
  return !!v && v.toString().trim().length > 0
})

const showAttachments = computed(() => {
  if (!props.attachmentField || !props.attachmentField.canReadRecordValue) return false
  const v = props.comment?.values?.[props.attachmentField.name]
  if (props.attachmentField.isMulti) {
    return Array.isArray(v) && v.length > 0
  }
  return !!v
})

// Reactions
const reactions = computed(() => {
  if (!props.reactionsField) return {}
  try {
    const val = props.comment?.values?.[props.reactionsField.name]
    return JSON.parse(val || '{}') || {}
  } catch {
    return {}
  }
})

const hasReactions = computed(() => Object.keys(reactions.value).length > 0)

function reactionTooltip(emoji, userIDs) {
  const names = userIDs.map(id => {
    if (id === props.currentUserID) return 'You'
    const user = props.findUserByID(id)
    return user?.name || user?.handle || user?.email || 'Unknown'
  })
  return names.join(', ')
}

function toggleEmojiPicker(event) {
  emojiPopoverRef.value?.toggle(event)
}

function onEmojiPopoverShow() {
  nextTick(() => {
    if (emojiPickerComp.value) {
      emojiPickerComp.value.reset()
    }
  })
}

function selectEmoji(emojiObj) {
  if (emojiObj && emojiObj.emoji) {
    emit('react', emojiObj.emoji)
  }
  emojiPopoverRef.value?.hide()
}

function onEdit() {
  isEditing.value = true
  editTitle.value = props.titleField ? props.comment?.values?.[props.titleField.name] || '' : ''
  editContent.value = props.contentField
    ? props.comment?.values?.[props.contentField.name] || ''
    : ''
}

function onCancel() {
  isEditing.value = false
}

function onSave() {
  emit('edit', {
    title: editTitle.value,
    content: editContent.value,
  })
  isEditing.value = false
}
</script>

<style scoped>
.avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  user-select: none;
  background: var(--p-content-hover-background);
  font-size: 0.8rem;
  flex-shrink: 0;
}

.comment-card {
  overflow: visible;
  transition: background-color 0.2s ease;
}

.comment-card-body {
  padding: 0.2rem 0.25rem;
}

.comment-item .comment-time {
  display: block;
  min-width: 3.5rem;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.comment-item .comment-time.always-visible {
  opacity: 1;
}

.comment-item .comment-toolbox {
  position: absolute;
  top: 0;
  right: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 1;
  background: var(--p-content-hover-background);
  border-radius: 0.5rem;
}

.comment-item.comment-highlighted .comment-card {
  background: var(--p-content-hover-background);
}

.comment-item:hover .comment-toolbox {
  opacity: 1;
}

.comment-item:hover .comment-card {
  background: var(--p-content-hover-background);
}

.comment-item:hover .comment-time {
  opacity: 1;
}

/* Disable hover effects */
.comment-item.no-hover:hover .comment-toolbox {
  opacity: 0;
}

.comment-item.no-hover:hover .comment-card {
  background: transparent;
}

.comment-item.no-hover:hover .comment-time {
  opacity: 0;
}

.comment-content :deep(p) {
  margin: 0;
}

.comment-content {
  overflow-wrap: break-word;
  word-wrap: break-word;
  word-break: break-word;
  white-space: pre-wrap;
}

/* Reaction badges */
.comment-reactions .reaction-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.2rem;
  padding: 0.1rem 0.4rem;
  border: 1px solid var(--p-content-border-color);
  border-radius: 1rem;
  background: var(--p-content-background);
  cursor: pointer;
  font-size: 0.8rem;
  line-height: 1.4;
  transition:
    background-color 0.15s,
    border-color 0.15s;
}

.comment-reactions .reaction-badge:hover {
  background-color: var(--p-content-hover-background);
  border-color: var(--p-text-muted-color);
}

.comment-reactions .reaction-badge.reaction-mine {
  background-color: color-mix(in srgb, var(--p-primary-color) 8%, transparent);
  border-color: var(--p-primary-color);
}

.comment-reactions .reaction-emoji {
  font-size: 0.9rem;
}

.comment-reactions .reaction-count {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--p-text-muted-color);
}

.comment-reactions .reaction-badge.reaction-mine .reaction-count {
  color: var(--p-primary-color);
}
</style>
