<template>
  <div
    class="comment-item"
    :class="{ 'comment-highlighted': highlighted, 'no-hover': disableHover }"
    @mouseleave="$emit('mouseleave')"
  >
    <!-- Author header (shown only for first message in a group) -->
    <div v-if="showHeader" class="flex items-center gap-4 px-2">
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
          :class="[
            'comment-time',
            'text-muted-color',
            'whitespace-nowrap',
            'overflow-hidden',
            { 'always-visible': showTimeAlways },
          ]"
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
            <div v-if="showTitle && titleField" class="font-semibold text-muted-color text-lg">
              <CFieldViewer :field="titleField" :record="comment" :namespace="namespace" />
            </div>

            <div v-if="showContent && contentField" class="comment-content">
              <CFieldViewer :field="contentField" :record="comment" :namespace="namespace" />
            </div>
          </template>
        </div>

        <!-- Hover toolbox -->
        <div v-if="!isEditing" class="comment-toolbox flex items-center justify-end gap-1">
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
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import CommentReply from './CommentReply.vue'
import { components } from '@cortezaproject/corteza-vue-next'

const { CRichTextInput, CFieldViewer } = components

const props = defineProps({
  comment: { type: Object, required: true },
  titleField: { type: Object, default: undefined },
  contentField: { type: Object, default: undefined },
  namespace: { type: Object, required: true },
  showHeader: { type: Boolean, default: true },
  showTimeAlways: { type: Boolean, default: false },
  showTitle: { type: Boolean, default: true },
  showContent: { type: Boolean, default: true },
  highlighted: { type: Boolean, default: false },
  disableHover: { type: Boolean, default: false },
})

const emit = defineEmits(['reply', 'edit', 'reply-click', 'mouseleave'])

const isEditing = ref(false)
const editTitle = ref('')
const editContent = ref('')

const authorName = computed(() => props.comment?.author?.name || '')
const authorInitials = computed(() => props.comment?.author?.initials || '?')
const authorIsCurrentUser = computed(() => Boolean(props.comment?.author?.isCurrentUser))

const canEdit = computed(() => authorIsCurrentUser.value && !props.comment?.deletedAt)

const titleValue = computed(() => {
  if (!props.titleField) return ''
  return props.comment?.values?.[props.titleField.name] || ''
})

const contentValue = computed(() => {
  if (!props.contentField) return ''
  return props.comment?.values?.[props.contentField.name] || ''
})

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

function onEdit() {
  isEditing.value = true
  editTitle.value = titleValue.value
  editContent.value = contentValue.value
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
  position: sticky;
  top: 0;
  align-self: flex-start;
  margin-left: auto;
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 1;
  flex-shrink: 0;
  order: 3;
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
</style>
