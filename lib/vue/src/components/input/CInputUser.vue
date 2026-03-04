<template>
  <AutoComplete
    :model-value="selectedUser"
    @update:model-value="onUpdateModelValue"
    @item-select="onItemSelect"
    :suggestions="suggestions"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    @complete="search"
    class="w-full"
    dropdown
    :complete-on-focus="true"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { debounce } from 'lodash-es'
import { inject, onMounted, ref, watch } from 'vue'
import { useUserResolver } from '../../composables/useUserResolver'

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
  clearOnSelect: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')
const { formatUser, resolveUser, cacheUsers } = useUserResolver()

const suggestions = ref([])
const selectedUser = ref(null)
const loading = ref(false)

function getOptionLabel(user) {
  return formatUser(user)
}

const search = debounce(async event => {
  // If empty query (focus click) and we already have suggestions, don't refetch
  if (!event.query && suggestions.value.length > 0) return

  loading.value = true
  try {
    const response = await $SystemAPI.userList({
      query: event.query,
      limit: 20,
    })
    suggestions.value = response.set || []
    cacheUsers(suggestions.value)
  } catch {
    suggestions.value = []
  } finally {
    loading.value = false
  }
}, 300)

function onUpdateModelValue(value) {
  // AutoComplete emits string when typing, and object when selected
  selectedUser.value = value
  emit('update:modelValue', value?.userID || null)
}

function onItemSelect(event) {
  emit('select', event.value)
  if (props.clearOnSelect) {
    // defer clearing to allow event to propagate and component to finish updates
    setTimeout(() => {
      selectedUser.value = null
      emit('update:modelValue', null)
    }, 0)
  }
}

async function loadUserById(userID) {
  if (!userID) return

  loading.value = true
  try {
    const user = await resolveUser(userID)
    if (user) selectedUser.value = user
  } catch {
    // User not found or API error - leave selectedUser as null
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedUser.value || selectedUser.value.userID !== newVal)) {
      loadUserById(newVal)
    } else if (!newVal) {
      selectedUser.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  if (props.modelValue) {
    loadUserById(props.modelValue)
  }
})
</script>
