<template>
  <div class="flex flex-col gap-3">
    <FileUpload
      ref="fileUploadRef"
      mode="advanced"
      :custom-upload="true"
      :auto="false"
      :multiple="multiple"
      :accept="accept || undefined"
      :show-upload-button="false"
      :show-cancel-button="false"
      :disabled="disabled"
      @select="onSelect"
    >
      <template #content>
        <div
          v-if="!stagedFiles.length && (!modelValue || !modelValue.length)"
          class="flex flex-col items-center py-6 text-muted-color"
        >
          <i class="pi pi-cloud-upload text-4xl mb-2" />
          <span class="text-sm">Drag and drop files here</span>
        </div>
        <div v-else class="p-2 text-sm text-muted-color">
          {{ stagedFiles.length + (modelValue?.length || 0) }} file(s) selected
        </div>
      </template>
    </FileUpload>

    <!-- Validation errors -->
    <div v-if="validationErrors.length" class="flex flex-col gap-1">
      <span
        v-for="(err, i) in validationErrors"
        :key="i"
        class="text-red-500 text-xs"
      >{{ err }}</span>
    </div>

    <!-- Staged (pending) files -->
    <div v-if="stagedFiles.length" class="flex flex-col gap-1">
      <span class="text-xs font-semibold text-muted-color uppercase tracking-wide">Pending upload</span>
      <div class="flex flex-wrap gap-2">
        <div
          v-for="staged in stagedFiles"
          :key="staged.id"
          class="flex items-center gap-1 bg-highlight border-surface rounded px-2 py-1 text-sm"
        >
          <i class="pi pi-file text-xs text-muted-color" />
          <span class="max-w-[160px] truncate">{{ staged.name }}</span>
          <span class="text-muted-color text-xs">({{ formatSize(staged.size) }})</span>
          <button
            type="button"
            class="ml-1 text-muted-color hover:text-red-500 transition-colors"
            :disabled="disabled"
            @click="removeStagedFile(staged.id)"
          >
            <i class="pi pi-times text-xs" />
          </button>
        </div>
      </div>
    </div>

    <!-- Already-uploaded attachment IDs -->
    <div v-if="modelValue && modelValue.length" class="flex flex-col gap-1">
      <span class="text-xs font-semibold text-muted-color uppercase tracking-wide">Uploaded</span>
      <div class="flex flex-wrap gap-2">
        <div
          v-for="(attID, index) in modelValue"
          :key="attID"
          class="flex items-center gap-1 bg-blue-50 border border-blue-200 rounded px-2 py-1 text-sm"
        >
          <i class="pi pi-paperclip text-xs text-blue-400" />
          <span class="max-w-[160px] truncate font-mono text-xs">{{ attID }}</span>
          <button
            type="button"
            class="ml-1 text-muted-color hover:text-red-500 transition-colors"
            :disabled="disabled"
            @click="removeUploaded(index)"
          >
            <i class="pi pi-times text-xs" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => [],
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  multiple: {
    type: Boolean,
    default: false,
  },
  accept: {
    type: String,
    default: '',
  },
  maxFileSize: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['update:modelValue', 'stage-files'])

const fileUploadRef = ref(null)
let nextId = 0
const stagedFiles = ref([])
const validationErrors = ref([])

function formatSize(bytes) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function validateFile(file) {
  if (props.accept) {
    const accepted = props.accept.split(',').some(type => {
      const t = type.trim()
      if (t.endsWith('/*')) return file.type.startsWith(t.slice(0, -2))
      return file.type === t || file.name.endsWith(t.replace('*', ''))
    })
    if (!accepted) return `"${file.name}" has an unsupported file type.`
  }
  if (props.maxFileSize > 0 && file.size > props.maxFileSize) {
    return `"${file.name}" exceeds the maximum size of ${formatSize(props.maxFileSize)}.`
  }
  return null
}

function onSelect(event) {
  const files = event.files || []
  const errors = []

  for (const file of files) {
    const error = validateFile(file)
    if (error) {
      errors.push(error)
      continue
    }
    const isDuplicate = stagedFiles.value.some(s => s.name === file.name && s.size === file.size)
    if (isDuplicate) continue
    stagedFiles.value.push({ id: nextId++, file, name: file.name, size: file.size })
  }

  validationErrors.value = errors

  // Clear PrimeVue's internal queue so it doesn't display its own list
  fileUploadRef.value?.clear()

  emit('stage-files', stagedFiles.value.map(s => s.file))
}

function removeStagedFile(id) {
  const index = stagedFiles.value.findIndex(s => s.id === id)
  if (index !== -1) stagedFiles.value.splice(index, 1)
  emit('stage-files', stagedFiles.value.map(s => s.file))
}

function removeUploaded(index) {
  const updated = [...(props.modelValue || [])]
  updated.splice(index, 1)
  emit('update:modelValue', updated)
}
</script>
