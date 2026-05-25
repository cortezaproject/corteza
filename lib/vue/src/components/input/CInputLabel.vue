<template>
  <div class="flex flex-col gap-3">
    <!-- MultiSelect + Create button -->
    <div class="flex gap-2">
      <MultiSelect
        :model-value="selectedLabelObjects"
        :options="allLabels"
        option-label="name"
        :placeholder="placeholder"
        :disabled="disabled"
        :loading="loading"
        class="flex-1"
        display="chip"
        filter
        fluid
        @update:model-value="onSelectionChange"
        @show="onShow"
      >
        <template #option="{ option }">
          <span>{{ option.name }}</span>
        </template>
      </MultiSelect>
      <Button
        icon="pi pi-plus"
        :aria-label="createLabel"
        v-tooltip.top="createLabel"
        @click="openCreateDialog"
        :disabled="disabled"
      />
    </div>

    <!-- Create Label Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="createDialogLabel"
      modal
      :style="{ width: '28rem' }"
    >
      <Form
        ref="dialogFormRef"
        v-slot="$form"
        :resolver="dialogResolver"
        :initialValues="dialogInitialValues"
        @submit="onDialogSubmit"
      >
        <div class="flex flex-col gap-4">
          <FormField name="name" class="flex flex-col gap-1">
            <label for="label-name" class="font-medium text-primary text-sm">
              {{ nameLabel }}
            </label>
            <InputText
              id="label-name"
              name="name"
              v-model="dialogForm.name"
            />
            <Message
              v-if="$form.name?.invalid"
              severity="error"
              size="small"
              variant="simple"
            >
              {{ $form.name.error?.message }}
            </Message>
          </FormField>
        </div>
      </Form>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="cancelBtnLabel"
            severity="secondary"
            outlined
            size="small"
            @click="dialogVisible = false"
          />
          <Button
            :label="saveBtnLabel"
            severity="primary"
            size="small"
            @click="dialogFormRef?.submit?.()"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, onMounted, ref, watch } from 'vue'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  placeholder: {
    type: String,
    default: 'Select labels...',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  kind: {
    type: String,
    default: '',
  },
  // Labels (i18n injected from parent)
  createLabel: { type: String, default: 'Create new' },
  createDialogLabel: { type: String, default: 'Create Label' },
  nameLabel: { type: String, default: 'Name' },
  saveBtnLabel: { type: String, default: 'Save' },
  cancelBtnLabel: { type: String, default: 'Cancel' },
  requiredMessage: { type: String, default: 'Name is required' },
  invalidHandleMessage: {
    type: String,
    default: 'Should be at least 2 characters long. Can contain only letters, numbers, underscores and dots. Must end with letter or number',
  },
})

const HANDLE_PATTERN = /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI')

const allLabels = ref([])
const loading = ref(false)

// Dialog state
const dialogVisible = ref(false)
const dialogForm = ref({ name: '' })
const dialogFormRef = ref(null)

let cancelCurrentRequest = null

// Convert modelValue map keys into label objects for MultiSelect binding
const selectedLabelObjects = computed(() => {
  if (!props.modelValue || typeof props.modelValue !== 'object') return []
  return Object.keys(props.modelValue)
    .map(name => allLabels.value.find(l => l.name === name))
    .filter(Boolean)
})

// --- Fetching Labels ---

async function fetchLabels() {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  try {
    const params = { limit: 200 }
    if (props.kind) {
      params.kind = props.kind
    }

    const { response, cancel } = $SystemAPI.labelListCancellable(params)
    cancelCurrentRequest = cancel

    const result = await response()
    const set = Array.isArray(result) ? result : result.set || []

    allLabels.value = set.map(l => ({
      name: l.name || '',
    }))

    // Ensure any labels in modelValue that aren't in the API response are still in the options
    ensureModelValueInOptions()
  } catch (e) {
    if (e?.message !== 'canceled') {
      allLabels.value = []
    }
  } finally {
    loading.value = false
    cancelCurrentRequest = null
  }
}

function ensureModelValueInOptions() {
  if (!props.modelValue) return
  for (const name of Object.keys(props.modelValue)) {
    if (!allLabels.value.find(l => l.name === name)) {
      allLabels.value.push({ name })
    }
  }
}

function onShow() {
  if (allLabels.value.length === 0) {
    fetchLabels()
  }
}

// --- Selection handling ---

function onSelectionChange(selected) {
  const labels = {}
  for (const item of selected) {
    labels[item.name] = props.modelValue?.[item.name] || ''
  }
  emit('update:modelValue', labels)
}

// --- Dialog ---

const dialogInitialValues = computed(() => ({ name: '' }))

function dialogResolver({ values }) {
  const errors = {}
  const name = (values?.name || '').trim()
  if (!name) {
    errors.name = [{ message: props.requiredMessage }]
  } else if (!HANDLE_PATTERN.test(name)) {
    errors.name = [{ message: props.invalidHandleMessage }]
  }
  return { errors }
}

function openCreateDialog() {
  dialogForm.value = { name: '' }
  dialogVisible.value = true
}

function onDialogSubmit({ valid, values }) {
  if (!valid) return
  const name = (values?.name || dialogForm.value.name || '').trim()
  if (!name) return

  // Add to options if not already there
  if (!allLabels.value.find(l => l.name === name)) {
    allLabels.value = [...allLabels.value, { name }]
  }

  // Add to selection
  const updated = { ...(props.modelValue || {}) }
  updated[name] = ''
  emit('update:modelValue', updated)

  dialogVisible.value = false
}

// --- Watchers ---

watch(
  () => props.modelValue,
  () => ensureModelValueInOptions(),
  { deep: true },
)

onMounted(() => {
  fetchLabels()
})

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
