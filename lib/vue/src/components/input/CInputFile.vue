<template>
  <div class="flex flex-col gap-3">
    <div
      class="flex flex-col items-center justify-center gap-3 p-6 border-2 border-dashed rounded-border cursor-pointer transition-colors duration-200"
      :class="[
        disabled ? 'opacity-50 cursor-not-allowed' : '',
        dragOver ? 'border-primary bg-primary/5' : 'border-surface-300 hover:border-primary',
      ]"
      @click="!disabled && $refs.fileInput.click()"
      @dragover.prevent="!disabled && (dragOver = true)"
      @dragleave.prevent="dragOver = false"
      @drop.prevent="onDrop"
    >
      <i class="pi pi-cloud-upload text-4xl text-muted-color" />
      <span class="text-muted-color text-sm text-center">
        {{ $t('general.label.dropFiles') }}
      </span>
    </div>

    <input
      ref="fileInput"
      type="file"
      :accept="accept || undefined"
      :multiple="multiple"
      :disabled="disabled"
      class="hidden"
      @change="onFileInputChange"
    />

    <!-- Validation errors -->
    <div v-if="validationErrors.length" class="flex flex-col gap-1">
      <span v-for="(err, i) in validationErrors" :key="i" class="text-red-500 text-xs">
        {{ err }}
      </span>
    </div>

    <!-- File list (staged + uploaded) -->
    <div v-if="stagedFiles.length || (modelValue && modelValue.length)" class="flex flex-col gap-1">
      <!-- Staged (pending) files -->
      <div
        v-for="staged in stagedFiles"
        :key="'staged-' + staged.id"
        class="flex items-center gap-2 px-3 py-2 border border-surface rounded"
      >
        <i class="pi pi-file text-muted-color" />
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium truncate">{{ staged.name }}</div>
          <div class="text-xs text-muted-color">{{ formatSize(staged.size) }}</div>
        </div>
        <span class="text-xs text-muted-color italic">pending</span>
        <button
          v-if="!disabled"
          type="button"
          class="text-muted-color hover:text-red-500 transition-colors"
          @click="removeStagedFile(staged.id)"
        >
          <i class="pi pi-trash text-sm" />
        </button>
      </div>

      <!-- Already-uploaded attachments -->
      <div
        v-for="(attID, index) in modelValue"
        :key="attID"
        class="flex items-center gap-2 px-3 py-2 border border-surface rounded"
      >
        <i class="pi pi-file text-primary" />
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium truncate">{{ attachmentInfo[attID]?.name || attID }}</div>
          <div v-if="attachmentInfo[attID]?.size" class="text-xs text-muted-color">
            {{ formatSize(attachmentInfo[attID].size) }}
          </div>
        </div>
        <button
          v-if="attachmentInfo[attID]?.downloadUrl"
          type="button"
          class="text-muted-color hover:text-primary transition-colors"
          @click="downloadFile(attachmentInfo[attID].downloadUrl)"
        >
          <i class="pi pi-download text-sm" />
        </button>
        <button
          v-if="!disabled"
          type="button"
          class="text-muted-color hover:text-red-500 transition-colors"
          @click="removeUploaded(index)"
        >
          <i class="pi pi-trash text-sm" />
        </button>
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
  attachmentInfo: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue', 'stage-files'])

const fileInput = ref(null)
const dragOver = ref(false)
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

function processFiles(files) {
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
  emit(
    'stage-files',
    stagedFiles.value.map(s => s.file),
  )
}

function onFileInputChange(event) {
  const files = Array.from(event.target?.files || [])
  if (files.length) processFiles(files)
  // Reset input so the same file can be re-selected
  if (fileInput.value) fileInput.value.value = ''
}

function onDrop(event) {
  if (props.disabled) return
  dragOver.value = false
  const files = Array.from(event.dataTransfer?.files || [])
  if (files.length) processFiles(files)
}

function removeStagedFile(id) {
  const index = stagedFiles.value.findIndex(s => s.id === id)
  if (index !== -1) stagedFiles.value.splice(index, 1)
  emit(
    'stage-files',
    stagedFiles.value.map(s => s.file),
  )
}

function downloadFile(url) {
  if (url) window.open(url, '_blank')
}

function removeUploaded(index) {
  const updated = [...(props.modelValue || [])]
  updated.splice(index, 1)
  emit('update:modelValue', updated)
}
</script>
