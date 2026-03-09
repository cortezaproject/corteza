<template>
  <AutoComplete
    :model-value="selectedRole"
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
    showClear
  >
    <template #option="{ option }">
      <div class="flex items-center gap-2">
        <span>{{ getOptionLabel(option) }}</span>
      </div>
    </template>
  </AutoComplete>
</template>

<script setup>
import AutoComplete from 'primevue/autocomplete'
import { debounce } from 'lodash-es'
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

const suggestions = ref([])
const selectedRole = ref(null)
const loading = ref(false)

// Store cancel function for current request
let cancelCurrentRequest = null

function getOptionLabel(role) {
  if (!role) return ''
  return role.name || role.handle || role.roleID
}

const search = debounce(async event => {
  // Cancel previous request if pending
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.roleListCancellable({
      query: event.query,
      limit: 20,
    })
    cancelCurrentRequest = cancel

    const result = await response()

    if (props.filterContextRoles) {
      suggestions.value = (result.set || []).filter(r => {
        // Filter out context-based roles (Authenticated, Anonymous, etc.)
        const isContext = r.meta?.context?.expr || r.meta?.context?.resourceTypes?.length > 0
        const isSpecialHandle = ['authenticated', 'anonymous', 'everyone'].includes(
          (r.handle || '').toLowerCase(),
        )
        return !(isContext || isSpecialHandle)
      })
    } else {
      suggestions.value = result.set || []
    }
  } catch (e) {
    // Ignore cancelled requests
    if (e?.message !== 'canceled') {
      suggestions.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}, 300)

function onUpdateModelValue(value) {
  selectedRole.value = value
  emit('update:modelValue', value?.roleID || null)
}

function onItemSelect(event) {
  emit('select', event.value)
  if (props.clearOnSelect) {
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
  } catch {
    // Role not found or API error - leave selectedRole as null
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
  if (props.modelValue && (!selectedRole.value || selectedRole.value.roleID !== props.modelValue)) {
    loadRoleById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  // Cancel any pending request
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
  search.cancel()
})
</script>
