<template>
  <AutoComplete
    :model-value="selectedUserGroup"
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
  clearOnSelect: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'select'])

const $SystemAPI = inject('$SystemAPI')

const suggestions = ref([])
const selectedUserGroup = ref(null)
const loading = ref(false)

function getOptionLabel(userGroup) {
  if (!userGroup) return ''
  return userGroup.name || userGroup.meta?.short || userGroup.handle || userGroup.userGroupID
}

const search = debounce(async event => {
  loading.value = true
  try {
    const response = await $SystemAPI.userGroupList({
      query: event.query,
      limit: 20,
    })
    suggestions.value = response.set || []
  } catch {
    suggestions.value = []
  } finally {
    loading.value = false
  }
}, 300)

function onUpdateModelValue(value) {
  selectedUserGroup.value = value
  emit('update:modelValue', value?.userGroupID || null)
}

function onItemSelect(event) {
  emit('select', event.value)
  if (props.clearOnSelect) {
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
  } catch {
    // UserGroup not found or API error - leave selectedUserGroup as null
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
  if (
    props.modelValue &&
    (!selectedUserGroup.value || selectedUserGroup.value.userGroupID !== props.modelValue)
  ) {
    loadUserGroupById(props.modelValue)
  }
})
</script>
