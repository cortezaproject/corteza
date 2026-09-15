<template>
  <!-- Multi-select mode -->
  <MultiSelect
    v-if="multiple"
    :model-value="selectedUsers"
    @update:model-value="onMultiSelect"
    :options="filteredOptions"
    option-label="label"
    option-value="userID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :size="size"
    :filter-fields="FILTER_FIELDS"
    fluid
    display="chip"
    @filter="onFilter"
    @show="onShow"
  >
    <template v-if="$slots.option" #option="slotProps">
      <slot name="option" v-bind="slotProps" />
    </template>
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
  </MultiSelect>

  <!-- Single-select mode -->
  <Select
    v-else
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
    :filter-fields="FILTER_FIELDS"
    fluid
    showClear
    @filter="onFilter"
    @show="onShow"
  >
    <template v-if="$slots.option" #option="slotProps">
      <slot name="option" v-bind="slotProps" />
    </template>
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
    type: [String, Number, Array],
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
  multiple: {
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
  excludeUsers: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')
const { formatUser, resolveUser, cacheUsers } = useUserResolver()

// Both branches filter on the same fields; naming it once keeps them in step.
const FILTER_FIELDS = ['label', 'name', 'handle', 'email', 'username']

const options = ref([])
const selectedUser = ref(null)
const selectedUsers = ref([])

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

    if (props.roleID.length) {
      params.roleID = props.roleID
    }

    const { response, cancel } = $SystemAPI.userListCancellable(params)
    cancelCurrentRequest = cancel

    const result = await response()
    const users = result.set || []

    // Add an explicit string label property for PrimeVue to easily render/filter
    options.value = users.map(u => ({ ...u, label: formatUser(u) }))

    // A page of results replaces the options wholesale. Anything already picked
    // has to go back in, or the value renders as a bare ID — or, single-select,
    // as nothing at all.
    if (props.multiple) {
      loadUsersByIds(selectedUsers.value)
    } else if (props.modelValue && !options.value.some(u => u.userID === props.modelValue)) {
      loadUserById(props.modelValue)
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
  fetchUsers()
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

function onMultiSelect(userIDs) {
  selectedUsers.value = userIDs || []
  emit('update:modelValue', selectedUsers.value)
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

async function loadUsersByIds(userIDs) {
  const missing = (userIDs || []).filter(id => id && !options.value.some(u => u.userID === id))
  if (!missing.length) return

  loading.value = true
  try {
    const resolved = await Promise.all(missing.map(id => resolveUser(id)))
    // Re-checked after the await: pinning runs from the watcher and from every
    // page fetch, so by now another pass may have added the same user.
    const have = new Set(options.value.map(u => u.userID))
    const found = resolved
      .filter(u => u && !have.has(u.userID))
      .map(u => ({ ...u, label: formatUser(u) }))
    if (found.length) {
      options.value = [...options.value, ...found]
    }
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (props.multiple) {
      const ids = Array.isArray(newVal) ? newVal : newVal ? [newVal] : []
      if (JSON.stringify(ids) !== JSON.stringify(selectedUsers.value)) {
        selectedUsers.value = ids
        loadUsersByIds(ids)
      }
      return
    }

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
  if (!props.multiple && props.modelValue && props.modelValue !== selectedUser.value) {
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
