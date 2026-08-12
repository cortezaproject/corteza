<template>
  <div
    class="comment-reply flex flex-col gap-1 overflow-hidden bg-surface border rounded-lg p-2 cursor-pointer"
    @click="$emit('click')"
  >
    <div class="flex items-center gap-1">
      <div :title="authorName" class="avatar flex items-center justify-center font-semibold">
        {{ authorInitials }}
      </div>

      <span
        :class="authorIsCurrentUser ? 'text-primary' : 'text-muted-color'"
        class="whitespace-nowrap font-semibold truncate text-sm"
      >
        {{ authorName }}
      </span>
    </div>

    <div class="text-sm text-muted-color truncate">
      {{ contentPreview }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  reply: { type: Object, required: true },
  titleField: { type: Object, default: undefined },
  contentField: { type: Object, default: undefined },
})

defineEmits(['click'])

const authorName = computed(() => props.reply?.author?.name || '')
const authorInitials = computed(() => props.reply?.author?.initials || '?')
const authorIsCurrentUser = computed(() => Boolean(props.reply?.author?.isCurrentUser))

const contentPreview = computed(() => {
  if (!props.contentField || !props.reply) return ''
  const val = props.reply.values?.[props.contentField.name] || ''
  // Strip HTML tags for preview
  return val.replace(/<[^>]*>/g, '').substring(0, 200)
})
</script>

<style scoped>
.avatar {
  width: 1.75rem;
  height: 1.75rem;
  font-size: 0.7rem;
  border-radius: 50%;
  user-select: none;
  background: var(--p-content-hover-background);
}

.comment-reply {
  max-height: 6.5rem;
  border-left-color: var(--p-primary-color) !important;
  border-left-width: 3px !important;
}

.comment-reply:hover {
  background: var(--p-content-hover-background);
}
</style>
