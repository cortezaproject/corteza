<template>
  <div>
    <span
      v-for="(user, index) in resolvedUsers"
      :key="user.userID || index"
      :class="{ block: isNewlineDelimiter, 'mt-1': isNewlineDelimiter && index !== 0 }"
    >
      {{ formatUser(user)
      }}{{ index !== resolvedUsers.length - 1 && !isNewlineDelimiter ? delimiter : '' }}
    </span>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useUserResolver } from '../../../composables/useUserResolver'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  record: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  valueOnly: {
    type: Boolean,
    default: false,
  },
  extraOptions: {
    type: Object,
    default: () => ({}),
  },
  disableClick: {
    type: Boolean,
    default: false,
  },
})

const { formatUser, findCached, resolveUsers } = useUserResolver()

const userIDs = computed(() => {
  const v = props.field.isSystem
    ? props.record[props.field.name]
    : props.record?.values?.[props.field.name]

  if (!v) return []
  if (props.field.isMulti && Array.isArray(v)) return v.filter(Boolean)
  return [v].filter(Boolean)
})

const delimiter = computed(() => props.field.options?.multiDelimiter || ', ')
const isNewlineDelimiter = computed(() => delimiter.value === '\n')

const resolvedUsers = computed(() => {
  return userIDs.value.map(id => findCached(id) || { userID: id })
})

// Trigger resolution of any users not yet in the store
watch(
  userIDs,
  ids => {
    if (ids.length) {
      resolveUsers(ids)
    }
  },
  { immediate: true },
)
</script>
