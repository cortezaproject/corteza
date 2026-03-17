<template>
  <Select
    :model-value="selectedRole"
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
import { inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

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
  filterContextRoles: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')

const options = ref([])
const selectedRole = ref(null)
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
    })
    cancelCurrentRequest = cancel

    const result = await response()

    if (props.filterContextRoles) {
      options.value = (result.set || []).filter(r => {
        const isContext = r.meta?.context?.expr || r.meta?.context?.resourceTypes?.length > 0
        const isSpecialHandle = ['authenticated', 'anonymous', 'everyone'].includes(
          (r.handle || '').toLowerCase(),
        )
        return !(isContext || isSpecialHandle)
      })
    } else {
      options.value = result.set || []
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

function onSelect(value) {
  selectedRole.value = value
  emit('update:modelValue', value?.roleID || null)

  if (value) {
    emit('select', value)
  }

  if (props.clearOnSelect && value) {
    setTimeout(() => {
      selectedRole.value = null
      emit('update:modelValue', null)
    }, 0)
  }
}

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

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedRole.value || selectedRole.value.roleID !== newVal)) {
      loadRoleById(newVal)
    } else if (!newVal) {
      selectedRole.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchRoles()
  if (props.modelValue && (!selectedRole.value || selectedRole.value.roleID !== props.modelValue)) {
    loadRoleById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
