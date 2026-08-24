<template>
  <CInputUser
    :model-value="modelValue"
    :multiple="isMultipleType"
    :placeholder="$t('field.kind.user.suggestionPlaceholder')"
    :disabled="disabled"
    :roleID="roles"
    @update:model-value="onChange"
  />
</template>

<script setup>
import { computed, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import CInputUser from '../../input/CInputUser.vue'

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

const $Auth = inject('$Auth', null)

const roles = computed(() => props.field.options?.roles || [])

// A multi-value field still edits one slot at a time unless it asks for the
// combined picker; `each` renders one editor per slot instead.
const isMultipleType = computed(
  () => props.field.isMulti && props.field.options?.selectType === 'multiple',
)

function onChange(value) {
  if (isMultipleType.value && props.field.options?.isUniqueMultiValue && Array.isArray(value)) {
    value = [...new Set(value)]
  }
  emit('update:modelValue', value)
}

// Preset with authenticated user on new records
onMounted(() => {
  if (props.field.options?.presetWithAuthenticated && $Auth?.user?.userID && !props.modelValue) {
    emit('update:modelValue', $Auth.user.userID)
  }
})
</script>
