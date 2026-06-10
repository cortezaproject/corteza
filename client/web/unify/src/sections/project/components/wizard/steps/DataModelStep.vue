<template>
  <div class="h-full overflow-auto p-6">
    <div>
      <CFormGroup label="Modules">
        <template #actions>
          <Button
            v-if="!disabled"
            icon="pi pi-plus"
            label="Add module"
            severity="secondary"
            size="small"
            @click="createResource?.('module')"
          />
        </template>

        <CFormItemList
          :items="modules"
          item-key="id"
          empty-message="No modules yet."
          :hide-remove="disabled"
          remove-label="Remove module"
          @select="m => configureResource?.(m.id)"
          @remove="removeModule"
        >
        <template #default="{ item }">
          <div class="flex items-center gap-3 min-w-0">
            <span
              class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
              :class="[cfg.bg, cfg.ring]"
            >
              <i :class="[cfg.icon, cfg.text]" />
            </span>
            <div class="min-w-0">
              <div class="font-medium truncate">{{ item.name }}</div>
              <div class="text-xs text-muted-color">{{ fieldSummary(item) }}</div>
            </div>
          </div>
        </template>
        </CFormItemList>
      </CFormGroup>
    </div>
  </div>
</template>

<script setup>
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useToast } from 'primevue/usetoast'
import { computed, inject } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const toast = useToast()
const configureResource = inject('configureResource', null)
const createResource = inject('createResource', null)

const cfg = kindConfig('module')
const modules = computed(() => (props.project.resources || []).filter(r => r.kind === 'module'))

const fieldSummary = m => {
  const n = (m.fields || []).length
  return n === 1 ? '1 field' : `${n} fields`
}

async function removeModule(m) {
  try {
    await store.removeResource(props.project.id, m.id)
    toast.add({ severity: 'success', summary: 'Module removed', detail: m.name, life: 2500 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Could not remove module', detail: err.message, life: 4000 })
  }
}
</script>
