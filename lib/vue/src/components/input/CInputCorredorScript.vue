<template>
  <Select
    :model-value="modelValue"
    @update:model-value="onSelect"
    :options="options"
    option-label="label"
    option-value="name"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    editable
    filter
    :filter-fields="['label', 'name']"
    fluid
    showClear
    @show="onShow"
    :pt="{ overlay: { class: 'max-w-lg' } }"
  >
    <template #option="{ option }">
      <div class="flex flex-col w-full min-w-0 whitespace-normal break-words">
        <span>{{ option.label }}</span>
        <small v-if="option.description" class="text-muted-color">
          {{ option.description }}
        </small>
      </div>
    </template>
  </Select>
</template>

<script setup>
import { inject, onBeforeUnmount, onMounted, ref } from 'vue'

defineOptions({ inheritAttrs: false })

// The model is a Corredor script name — its path inside the extension, e.g.
// `/server-scripts/SystemPing.js:default`. The Select is editable, so a name
// the list does not carry (deployed later, or hidden from the caller) can be
// typed and is kept as typed.
defineProps({
  modelValue: {
    type: String,
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
})

const emit = defineEmits(['update:modelValue'])

const $SystemAPI = inject('$SystemAPI', null)

const options = ref([])
const loading = ref(false)

let cancelCurrentRequest = null

// Server-side scripts with a manual trigger on the system resource; the server
// leaves out the ones the caller may not run.
const SCRIPT_FILTER = {
  eventTypes: ['onManual'],
  resourceTypes: ['system'],
  excludeInvalid: true,
  excludeClientScripts: true,
}

function getOptionLabel(script) {
  if (!script) return ''
  return script.label || script.name
}

async function fetchScripts() {
  if (!$SystemAPI) return

  if (cancelCurrentRequest) {
    cancelCurrentRequest()
    cancelCurrentRequest = null
  }

  loading.value = true
  const { response, cancel } = $SystemAPI.automationListCancellable(SCRIPT_FILTER)
  cancelCurrentRequest = cancel

  try {
    const result = await response()
    const scripts = Array.isArray(result) ? result : result?.set || []

    options.value = scripts
      .filter(script => script?.name)
      .map(script => ({ ...script, label: getOptionLabel(script) }))
      .sort((a, b) => (a.label || '').localeCompare(b.label || ''))
  } catch {
    // Corredor unreachable or the request cancelled; a typed name still works
  } finally {
    if (cancelCurrentRequest === cancel) {
      loading.value = false
      cancelCurrentRequest = null
    }
  }
}

function onShow() {
  if (options.value.length === 0) {
    fetchScripts()
  }
}

function onSelect(name) {
  emit('update:modelValue', name || null)
}

onMounted(fetchScripts)

onBeforeUnmount(() => {
  if (cancelCurrentRequest) {
    cancelCurrentRequest()
  }
})
</script>
