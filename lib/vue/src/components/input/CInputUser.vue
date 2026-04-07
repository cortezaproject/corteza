<template>
  <Select
    :model-value="selectedUser"
    @update:model-value="onSelect"
    :options="filteredOptions"
    option-label="label"
    option-value="userID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :size="size"
    :filter-fields="['label', 'name', 'handle', 'email', 'username']"
    fluid
    showClear
    @filter="onFilter"
    @show="onShow"
  >
    <template #footer>
      <div
        v-if="hasNextPage || hasPrevPage"
        class="flex justify-between items-center px-3 py-2 border-t border-surface"
      >
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
import { inject, onBeforeUnmount, onMounted, ref, watch, nextTick, computed } from 'vue'
import { useUserResolver } from '../../composables/useUserResolver'

defineOptions({ inheritAttrs: false })

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
  roleID: {
    type: Array,
    default: () => [],
  },
  size: {
    type: String,
    default: '',
  },
  roleId: {
    // for backwards compatibility
    type: Array,
    default: () => [],
  },
  excludeUsers: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')
const { formatUser, resolveUser, cacheUsers } = useUserResolver()

const options = ref([])
const selectedUser = ref(null)

const filteredOptions = computed(() => {
  if (!props.excludeUsers || props.excludeUsers.length === 0) return options.value
  return options.value.filter(
    u =>
      !props.excludeUsers.includes(u.userID) &&
      !props.excludeUsers.includes(u.handle) &&
      !props.excludeUsers.includes(u.email),
  )
})

const loading = ref(false)

let cancelCurrentRequest = null
let searchTimeout = null

// Pagination state
const nextPage = ref('')
const prevPage = ref('')
const hasNextPage = ref(false)
const hasPrevPage = ref(false)
const currentQuery = ref('')

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
      sort: 'name ASC',
    }

    if (pageCursor) {
      params.pageCursor = pageCursor
    }

    const effectiveRoleID = props.roleID?.length ? props.roleID : props.roleId
    if (effectiveRoleID && effectiveRoleID.length > 0) {
      params.roleID = effectiveRoleID
    }

    const { response, cancel } = $SystemAPI.userListCancellable(params)
    cancelCurrentRequest = cancel

    const result = await response()
    const users = result.set || []

    // Add an explicit string label property for PrimeVue to easily render/filter
    options.value = users.map(u => ({ ...u, label: formatUser(u) }))

    // If we have a bound ID that isn't in this new page, append it to prevent deselection
    if (props.modelValue) {
      if (!options.value.some(u => u.userID === props.modelValue)) {
        loadUserById(props.modelValue)
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

function onSelect(userID) {
  const selectedNode = options.value.find(u => u.userID === userID)

  if (props.clearOnSelect) {
    selectedUser.value = userID
    if (selectedNode) emit('select', selectedNode)

    nextTick(() => {
      selectedUser.value = null
      emit('update:modelValue', null)
    })
  } else {
    selectedUser.value = userID
    emit('update:modelValue', userID || null)
    if (selectedNode) emit('select', selectedNode)
  }
}

async function loadUserById(userID) {
  if (!userID) return

  loading.value = true
  try {
    const user = await resolveUser(userID)
    if (user) {
      user.label = formatUser(user)
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
    if (newVal) {
      if (newVal !== selectedUser.value) {
        selectedUser.value = newVal
        loadUserById(newVal)
      }
    } else {
      selectedUser.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchUsers()
  if (props.modelValue && props.modelValue !== selectedUser.value) {
    selectedUser.value = props.modelValue
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
