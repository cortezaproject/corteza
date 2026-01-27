<template>
  <AutoComplete
    :model-value="selectedUser"
    @update:model-value="onSelect"
    :suggestions="suggestions"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    @complete="search"
    class="w-full"
    dropdown
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ getOptionLabel(option) }}</span>
        <span v-if="option.email" class="text-muted-color text-xs">{{ option.email }}</span>
      </div>
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { debounce } from 'lodash-es'
import { inject, onMounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: null,
  },
  placeholder: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')

const suggestions = ref([])
const selectedUser = ref(null)
const loading = ref(false)

function getOptionLabel(user) {
  if (!user) return ''
  return user.name || user.handle || user.email || user.userID
}

const search = debounce(async (event) => {
  loading.value = true
  try {
    const response = await $SystemAPI.userList({
      query: event.query,
      limit: 20,
    })
    suggestions.value = response.set || []
  } catch (_e) {
    suggestions.value = []
  } finally {
    loading.value = false
  }
}, 300)

function onSelect(value) {
  selectedUser.value = value
  emit('update:modelValue', value?.userID || null)
}

async function loadUserById(userID) {
  if (!userID || !$SystemAPI) return
  loading.value = true
  try {
    const user = await $SystemAPI.userRead({ userID })
    selectedUser.value = user
  } catch (_e) {
    // User not found or API error - leave selectedUser as null
  } finally {
    loading.value = false
  }
}

watch(() => props.modelValue, (newVal) => {
  if (newVal && (!selectedUser.value || selectedUser.value.userID !== newVal)) {
    loadUserById(newVal)
  } else if (!newVal) {
    selectedUser.value = null
  }
}, { immediate: true })

onMounted(() => {
  if (props.modelValue) {
    loadUserById(props.modelValue)
  }
})
</script>
