<template>
  <Dialog
    v-model:visible="visible"
    modal
    header="Add a connection"
    :style="{ width: '60rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-end' } }"
  >
    <p class="text-sm text-muted-color mb-3">
      Pick a connection your system needs. Credentials are configured later in Admin.
    </p>

    <IconField class="mb-3">
      <InputIcon class="pi pi-search" />
      <InputText v-model="query" placeholder="Search connections…" fluid autofocus />
    </IconField>

    <div class="max-h-[26rem] overflow-auto p-1">
      <div class="grid grid-cols-2 lg:grid-cols-3 gap-3">
        <button
          v-for="c in filtered"
          :key="c.id"
          type="button"
          class="flex items-start gap-3 p-3 rounded-xl border border-surface hover:border-primary hover:bg-emphasis transition-colors text-left"
          @click="$emit('pick', c)"
        >
          <span class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 ring-surface bg-emphasis shrink-0">
            <i :class="[c.icon, 'text-lg text-primary']" />
          </span>
          <div class="min-w-0">
            <div class="font-medium text-sm leading-tight">{{ c.label }}</div>
            <div class="text-xs text-muted-color line-clamp-2">{{ c.description }}</div>
          </div>
        </button>
      </div>
      <div v-if="!filtered.length" class="text-sm text-muted-color italic text-center py-6">
        No connections match “{{ query }}”.
      </div>
    </div>

    <template #footer>
      <Button label="Cancel" severity="secondary" outlined size="small" @click="visible = false" />
    </template>
  </Dialog>
</template>

<script setup>
import { CONNECTORS } from '@/sections/project-poc/config/connectors'
import { computed, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue', 'pick'])

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const query = ref('')
watch(visible, v => {
  if (v) query.value = ''
})

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return CONNECTORS
  return CONNECTORS.filter(
    c =>
      c.label.toLowerCase().includes(q) ||
      c.description.toLowerCase().includes(q) ||
      (c.tags || []).some(t => t.includes(q)),
  )
})
</script>
