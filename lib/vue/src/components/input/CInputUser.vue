<template>
  <Select
    :model-value="selectedUser"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    fluid
    showClear
    @filter="onFilter"
    @show="onShow"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
    <template v-if="hasNextPage || hasPrevPage" #footer>
      <div class="flex justify-between items-center px-3 py-2 border-t border-surface">
        <Button
          icon="pi pi-angle-left"
          text
          size="small"
          :disabled="!hasPrevPage"
          @click="goToPage(false)"
        />
        <Button
          icon="pi pi-angle-right"
          text
          size="small"
          :disabled="!hasNextPage"
          @click="goToPage(true)"
        />
      </div>
    </template>
  </Select>
</template>

<script setup>
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'
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
  roleId: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')
const { formatUser, resolveUser, cacheUsers } = useUserResolver()

const options = ref([])
const selectedUser = ref(null)
const loading = ref(false)

let cancelCurrentRequest = null
let searchTimeout = null

// Pagination state
const nextPage = ref('')
const prevPage = ref('')
const hasNextPage = ref(false)
const hasPrevPage = ref(false)
const currentQuery = ref('')

function getOptionLabel(user) {
  return formatUser(user)
}

async function fetchUsers(query = '', pageCursor = '') {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const params = {
      query,
      limit: 15,
    }

    if (pageCursor) {
      params.pageCursor = pageCursor
    }

    if (props.roleId && props.roleId.length > 0) {
      params.roleID = props.roleId
    }

    const { response, cancel } = $SystemAPI.userListCancellable(params)
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = result.set || []

    // Re-sync selectedUser reference with the matching option from the new set
    // so PrimeVue Select can match it by reference
    if (selectedUser.value) {
      const match = options.value.find(u => u.userID === selectedUser.value.userID)
      if (match) {
        selectedUser.value = match
      } else {
        // Selected user not in current page — keep them in options
        options.value = [...options.value, selectedUser.value]
      }
    }

    nextPage.value = result.filter?.nextPage || ''
    prevPage.value = result.filter?.prevPage || ''
    hasNextPage.value = !!nextPage.value
    hasPrevPage.value = !!prevPage.value
    cacheUsers(options.value)
  } catch (e) {
    if (e?.message !== 'canceled') {
      options.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onFilter(event) {
  currentQuery.value = event.value || ''
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchUsers(currentQuery.value)
  }, 300)
}

function onShow() {
  if (options.value.length === 0) {
    fetchUsers()
  }
}

function goToPage(next) {
  const cursor = next ? nextPage.value : prevPage.value
  if (cursor) {
    fetchUsers(currentQuery.value, cursor)
  }
}

function onSelect(value) {
  selectedUser.value = value
  emit('update:modelValue', value?.userID || null)

  if (value) {
    emit('select', value)
  }

  if (props.clearOnSelect && value) {
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
    if (user) {
      selectedUser.value = user
      if (!options.value.find(u => u.userID === userID)) {
        options.value = [...options.value, user]
      }
    }
  } catch {
    // User not found or API error
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
  fetchUsers()
  if (props.modelValue) {
    loadUserById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
})
</script>

