<template>
  <Select
    :model-value="selectedUserGroup"
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
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')

const options = ref([])
const selectedUserGroup = ref(null)
const loading = ref(false)

let cancelCurrentRequest = null

function getOptionLabel(userGroup) {
  if (!userGroup) return ''
  return userGroup.name || userGroup.meta?.short || userGroup.handle || userGroup.userGroupID
}

async function fetchUserGroups() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const { response, cancel } = $SystemAPI.userGroupListCancellable({
      query: '',
      limit: 100,
    })
    cancelCurrentRequest = cancel

    const result = await response()
    options.value = result.set || []
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
    fetchUserGroups()
  }
}

function onSelect(value) {
  selectedUserGroup.value = value
  emit('update:modelValue', value?.userGroupID || null)

  if (value) {
    emit('select', value)
  }

  if (props.clearOnSelect && value) {
    setTimeout(() => {
      selectedUserGroup.value = null
      emit('update:modelValue', null)
    }, 0)
  }
}

async function loadUserGroupById(userGroupID) {
  if (!userGroupID || !$SystemAPI) return
  loading.value = true
  try {
    const userGroup = await $SystemAPI.userGroupRead({ userGroupID })
    selectedUserGroup.value = userGroup
    if (!options.value.find(g => g.userGroupID === userGroupID)) {
      options.value = [...options.value, userGroup]
    }
  } catch {
    // UserGroup not found or API error
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedUserGroup.value || selectedUserGroup.value.userGroupID !== newVal)) {
      loadUserGroupById(newVal)
    } else if (!newVal) {
      selectedUserGroup.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  fetchUserGroups()
  if (
    props.modelValue &&
    (!selectedUserGroup.value || selectedUserGroup.value.userGroupID !== props.modelValue)
  ) {
    loadUserGroupById(props.modelValue)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
