<template>
  <slot :trigger="trigger">
    <Button
      :label="label"
      :icon="icon"
      :severity="severity"
      :text="text"
      :outlined="outlined"
      :size="size"
      :disabled="disabled"
      @click="trigger"
    />
  </slot>
</template>

<script setup>
import { useConfirmDelete } from '../../composables/useConfirmDelete'

const { confirmDelete } = useConfirmDelete()

const props = defineProps({
  // Button props (used by default slot)
  label: {
    type: String,
    default: '',
  },
  icon: {
    type: String,
    default: 'pi pi-trash',
  },
  severity: {
    type: String,
    default: 'danger',
  },
  text: {
    type: Boolean,
    default: false,
  },
  outlined: {
    type: Boolean,
    default: false,
  },
  size: {
    type: String,
    default: undefined,
  },
  disabled: {
    type: Boolean,
    default: false,
  },

  // Confirm dialog props
  message: {
    type: String,
    default: '',
  },
  header: {
    type: String,
    default: '',
  },
  confirmIcon: {
    type: String,
    default: 'pi pi-exclamation-triangle',
  },
  acceptProps: {
    type: Object,
    default: undefined,
  },
  rejectProps: {
    type: Object,
    default: undefined,
  },
})

const emit = defineEmits(['confirm'])

function trigger() {
  if (props.disabled) return

  confirmDelete({
    message: props.message,
    header: props.header,
    icon: props.confirmIcon,
    ...(props.acceptProps && { acceptProps: props.acceptProps }),
    ...(props.rejectProps && { rejectProps: props.rejectProps }),
    onConfirm: () => emit('confirm'),
  })
}
</script>
