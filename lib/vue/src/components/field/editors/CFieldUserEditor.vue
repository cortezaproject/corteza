<template>
  <!-- Multi-value: MultiSelect mode (selectType === 'multiple') -->
  <MultiSelect
    v-if="isMultipleType"
    :model-value="multiValue"
    :options="userOptions"
    :option-label="getOptionLabel"
    option-value="userID"
    :placeholder="$t('field.kind.user.suggestionPlaceholder')"
    :disabled="disabled"
    :loading="loading"
    filter
    class="w-full"
    @filter="onFilter"
    @update:model-value="onMultiSelectChange"
  >
    <template #option="{ option }">
      {{ getOptionLabel(option) }}
    </template>
  </MultiSelect>

  <!-- Single-value or default/each multi: CInputUser per slot -->
  <CInputUser
    v-else
    :model-value="modelValue"
    :placeholder="$t('field.kind.user.suggestionPlaceholder')"
    :disabled="disabled"
    :role-id="roles"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed, inject, onMounted, onBeforeUnmount, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import CInputUser from '../../input/CInputUser.vue'
import { useUserResolver } from '../../../composables/useUserResolver'

const { t: $t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: [String, Array],
    default: '',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth', null)
const { formatUser, resolveUser, cacheUsers } = useUserResolver()

const userOptions = ref([])
const loading = ref(false)
let cancelCurrentRequest = null
let searchTimeout = null

const roles = computed(() => props.field.options?.roles || [])
const isMultipleType = computed(
  () => props.field.isMulti && props.field.options?.selectType === 'multiple',
)

function getOptionLabel(user) {
  return formatUser(user)
}

// Normalize array model value for MultiSelect
const multiValue = computed(() => {
  if (Array.isArray(props.modelValue)) return props.modelValue.filter(Boolean)
  return props.modelValue ? [props.modelValue] : []
})

function onMultiSelectChange(values) {
  if (props.field.options?.isUniqueMultiValue) {
    values = [...new Set(values)]
  }
  emit('update:modelValue', values)
}

async function fetchUsers(query = '') {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const params = { query, limit: 50 }
    if (roles.value.length > 0) {
      params.roleID = roles.value
    }

    const { response, cancel } = $SystemAPI.userListCancellable(params)
    cancelCurrentRequest = cancel
    const result = await response()
    userOptions.value = result.set || []
    cacheUsers(userOptions.value)
  } catch (e) {
    if (e?.message !== 'canceled') {
      userOptions.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function onFilter(event) {
  const query = event.value || ''
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchUsers(query), 300)
}

// Resolve existing user IDs so MultiSelect can display labels
async function resolveExistingUsers() {
  const ids = multiValue.value
  if (!ids.length) return

  for (const id of ids) {
    if (id && !userOptions.value.find((u) => u.userID === id)) {
      try {
        const user = await resolveUser(id)
        if (user) {
          userOptions.value = [...userOptions.value, user]
        }
      } catch {
        // ignore
      }
    }
  }
}

// Preset with authenticated user on new records
onMounted(async () => {
  if (isMultipleType.value) {
    await fetchUsers()
    await resolveExistingUsers()
  }

  // presetWithAuthenticated: auto-fill with current user if value is empty
  if (
    props.field.options?.presetWithAuthenticated &&
    $Auth?.user?.userID &&
    !props.modelValue
  ) {
    emit('update:modelValue', $Auth.user.userID)
  }
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) cancelCurrentRequest()
  if (searchTimeout) clearTimeout(searchTimeout)
})
</script>
