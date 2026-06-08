<template>
  <Select
    :model-value="selectedNamespace"
    @update:model-value="onSelect"
    :options="options"
    :option-label="getOptionLabel"
    data-key="namespaceID"
    :placeholder="placeholder"
    :disabled="disabled"
    :loading="loading"
    class="w-full"
    filter
    :filter-fields="['name', 'slug', 'namespaceID']"
    fluid
    showClear
    @show="onShow"
  >
    <template #option="{ option }">
      {{ option.name }}
    </template>
  </Select>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useNamespaceStore } from '../../stores/useNamespaceStore'

defineOptions({ inheritAttrs: false })

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
})

const emit = defineEmits(['update:modelValue'])

const store = useNamespaceStore()
const selectedNamespace = ref(null)
const loading = ref(false)

const options = computed(() => store.set)

function getOptionLabel(namespace) {
  if (!namespace) return ''
  return namespace.name || namespace.slug || namespace.namespaceID
}

async function ensureLoaded() {
  if (store.set.length > 0) return
  loading.value = true
  try {
    await store.load()
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function onShow() {
  ensureLoaded()
}

function onSelect(value) {
  selectedNamespace.value = value
  emit('update:modelValue', value?.namespaceID || null)
}

async function loadNamespaceById(namespaceID) {
  if (!namespaceID) return

  const existing = store.set.find(ns => ns.namespaceID === namespaceID)
  if (existing) {
    selectedNamespace.value = existing
    return
  }

  loading.value = true
  try {
    const namespace = await store.findByID({ namespaceID })
    if (namespace) {
      selectedNamespace.value = store.set.find(ns => ns.namespaceID === namespaceID) || namespace
    }
  } catch {
    // namespace not found
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  newVal => {
    if (newVal && (!selectedNamespace.value || selectedNamespace.value.namespaceID !== newVal)) {
      loadNamespaceById(newVal)
    } else if (!newVal) {
      selectedNamespace.value = null
    }
  },
  { immediate: true },
)

onMounted(() => {
  ensureLoaded()

  if (props.modelValue) {
    loadNamespaceById(props.modelValue)
  }
})
</script>
