<template>
  <div class="flex flex-col gap-3">
    <!-- Preview slot (shown when there's an existing value) -->
    <div v-if="hasPreview">
      <slot name="preview">
        <div
          v-if="previewUrl"
          class="relative flex items-center justify-center border border-surface rounded-border"
          :style="previewContainerStyle"
        >
          <img :src="previewUrl" :alt="label" class="w-auto h-full object-contain p-4" />

          <Button
            v-if="!disabled && clearable"
            icon="pi pi-trash"
            severity="danger"
            text
            size="small"
            class="!absolute top-1 right-1"
            @click="$emit('clear')"
          />
        </div>
      </slot>
    </div>

    <!-- Drop zone -->
    <div
      class="flex flex-col items-center justify-center gap-3 border-2 border-dashed rounded-border cursor-pointer transition-colors duration-200"
      :class="[
        disabled ? 'opacity-50 cursor-not-allowed' : '',
        dragOver ? 'border-primary bg-primary/5' : 'border-surface-300 hover:border-primary',
        compact ? 'p-4' : 'p-6',
      ]"
      @click="!disabled && openFileDialog()"
      @dragover.prevent="!disabled && (dragOver = true)"
      @dragleave.prevent="dragOver = false"
      @drop.prevent="onDrop"
    >
      <ProgressSpinner v-if="uploading" style="width: 2rem; height: 2rem" />
      <i v-else :class="[icon, 'text-4xl text-muted-color']" />
      <span class="text-muted-color text-sm text-center">
        {{ uploading ? uploadingLabel : dropLabel }}
      </span>
    </div>

    <input
      ref="fileInputRef"
      type="file"
      :accept="accept || undefined"
      :multiple="multiple"
      :disabled="disabled"
      class="hidden"
      @change="onFileInputChange"
    />

    <!-- Error message -->
    <Message v-if="error" severity="error" size="small" :closable="false">
      {{ error }}
    </Message>
  </div>
</template>

<script setup>
import { ref, computed, useSlots } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const $slots = useSlots()

const props = defineProps({
  /** File type filter (e.g., '.zip', 'image/*') */
  accept: {
    type: String,
    default: '',
  },
  /** Allow multiple file selection */
  multiple: {
    type: Boolean,
    default: false,
  },
  /** External uploading state */
  uploading: {
    type: Boolean,
    default: false,
  },
  /** Disable interaction */
  disabled: {
    type: Boolean,
    default: false,
  },
  /** Show the clear/delete button */
  clearable: {
    type: Boolean,
    default: false,
  },
  /** Current preview URL (for image uploads) */
  previewUrl: {
    type: String,
    default: '',
  },
  /** Drop zone label text */
  dropLabel: {
    type: String,
    default: '',
  },
  /** Uploading state label text */
  uploadingLabel: {
    type: String,
    default: '',
  },
  /** Upload error message */
  error: {
    type: String,
    default: '',
  },
  /** PrimeVue icon class */
  icon: {
    type: String,
    default: 'pi pi-cloud-upload',
  },
  /** Use compact padding */
  compact: {
    type: Boolean,
    default: false,
  },
  /** Preview container max dimensions */
  previewMaxWidth: {
    type: String,
    default: '200px',
  },
  previewMaxHeight: {
    type: String,
    default: '100px',
  },
  /** Label for alt text */
  label: {
    type: String,
    default: 'Preview',
  },
})

const emit = defineEmits(['select', 'clear'])

const fileInputRef = ref(null)
const dragOver = ref(false)

const hasPreview = computed(() => {
  return !!props.previewUrl || !!$slots.preview
})

const previewContainerStyle = computed(() => ({
  width: props.previewMaxWidth,
  height: props.previewMaxHeight,
}))

const dropLabel = computed(() => {
  return props.dropLabel || t('general.label.dropFiles')
})

const uploadingLabel = computed(() => {
  return props.uploadingLabel || t('general.label.uploading')
})

function openFileDialog() {
  fileInputRef.value?.click()
}

function onFileInputChange(event) {
  const files = Array.from(event.target?.files || [])
  if (files.length) {
    emit('select', files)
  }
  // Reset input so the same file can be re-selected
  if (fileInputRef.value) fileInputRef.value.value = ''
}

function onDrop(event) {
  if (props.disabled) return
  dragOver.value = false
  const files = Array.from(event.dataTransfer?.files || [])
  if (files.length) {
    emit('select', files)
  }
}

defineExpose({ openFileDialog })
</script>
