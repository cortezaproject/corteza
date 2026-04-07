<template>
  <!-- Multi-select mode -->
  <MultiSelect
    v-if="multiple"
    :model-value="selectedRoles"
    @update:model-value="onMultiSelect"
    :options="filteredOptions"
    :option-label="getOptionLabel"
    option-value="roleID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    :size="size"
    filter
    fluid
    display="chip"
    @show="onShow"
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ getOptionLabel(option) }}</span>
      </div>
    </template>
  </MultiSelect>

  <!-- Single-select mode -->
  <Select
    v-else
    :model-value="selectedRole"
    @update:model-value="onSelect"
    :options="filteredOptions"
    :option-label="getOptionLabel"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    fluid
    :showClear="showClear"
    @show="onShow"
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ getOptionLabel(option) }}</span>
      </div>
    </template>
  </Select>
</template>

<script setup>
import { inject, onBeforeUnmount, onMounted, ref, watch, nextTick, computed } from 'vue'

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
  filterContextRoles: {
    type: Boolean,
    default: false,
  },
  excludeRoles: {
    type: Array,
    default: () => [],
  },
  size: {
    type: String,
    default: '',
  },
  showClear: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')

const selectedRole = ref(null)
const selectedRoles = ref([])
const options = ref([])

const filteredOptions = computed(() => {
  let list = options.value
  if (props.excludeRoles && props.excludeRoles.length > 0) {
    list = list.filter(
      r =>
        !props.excludeRoles.includes(r.handle) &&
        !props.excludeRoles.includes(r.roleID) &&
        !props.excludeRoles.includes(r.name),
    )
  }
  return list
})

const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(role) {
  if (!role) return ''
  return role.name || role.handle || role.roleID
}

async function fetchRoles() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.roleListCancellable({
      query: '',
      limit: 100,
      sort: 'name ASC',
    })
    cancelCurrentRequest = cancel

    const result = await response()
    let fetchedOptions = result.set || []

    if (props.filterContextRoles) {
      options.value = fetchedOptions.filter(r => {
        const isContext = r.meta?.context?.expr || r.meta?.context?.resourceTypes?.length > 0
        const isSpecialHandle = ['authenticated', 'anonymous', 'everyone'].includes(
          (r.handle || '').toLowerCase(),
        )
        return !(isContext || isSpecialHandle)
      })
    } else {
      options.value = fetchedOptions
    }
  } catch (e) {
    if (e?.message !== 'canceled') {
      options.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onShow() {
  if (options.value.length === 0) {
    fetchRoles()
  }
}

// --- Single-select handlers ---
function onSelect(value) {
  if (props.clearOnSelect) {
    selectedRole.value = value
    if (value) emit('select', value)

    nextTick(() => {
      selectedRole.value = null
      emit('update:modelValue', null)
    })
  } else {
    selectedRole.value = value
    emit('update:modelValue', value?.roleID || null)
    if (value) emit('select', value)
  }
}

// --- Multi-select handlers ---
function onMultiSelect(roleIDs) {
  selectedRoles.value = roleIDs || []
  emit('update:modelValue', selectedRoles.value)
}

// --- Load role(s) by ID ---
async function loadRoleById(roleID) {
  if (!roleID || !$SystemAPI) return
  loading.value = true
  try {
    const role = await $SystemAPI.roleRead({ roleID })
    selectedRole.value = role
    if (!options.value.find(r => r.roleID === roleID)) {
      options.value = [...options.value, role]
    }
  } catch {
    // Role not found or API error
  } finally {
    loading.value = false
  }
}

async function loadRolesByIds(roleIDs) {
  if (!roleIDs?.length || !$SystemAPI) return
  loading.value = true
  try {
    // Ensure all roleIDs are in options
    const missing = roleIDs.filter(id => !options.value.find(r => r.roleID === id))
    if (missing.length) {
      const results = await Promise.all(
        missing.map(roleID => $SystemAPI.roleRead({ roleID }).catch(() => null)),
      )
      const loaded = results.filter(Boolean)
      if (loaded.length) {
        options.value = [...options.value, ...loaded]
      }
    }
    selectedRoles.value = roleIDs
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

// --- Watchers ---
watch(
  () => props.modelValue,
  newVal => {
    if (props.multiple) {
      const arr = Array.isArray(newVal) ? newVal : newVal ? [newVal] : []
      if (JSON.stringify(arr) !== JSON.stringify(selectedRoles.value)) {
        loadRolesByIds(arr)
      }
    } else {
      if (newVal && (!selectedRole.value || selectedRole.value.roleID !== newVal)) {
        loadRoleById(newVal)
      } else if (!newVal) {
        selectedRole.value = null
      }
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchRoles()
  if (props.multiple) {
    const arr = Array.isArray(props.modelValue) ? props.modelValue : []
    if (arr.length) {
      loadRolesByIds(arr)
    }
  } else if (
    props.modelValue &&
    (!selectedRole.value || selectedRole.value.roleID !== props.modelValue)
  ) {
    loadRoleById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
