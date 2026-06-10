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
          @remove="m => store.removeResource(project.id, m.id)"
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
import { kindConfig } from '@/sections/project-poc/config/kinds'
import { useProjectsStore } from '@/sections/project-poc/stores/projects'
import { computed, inject } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const configureResource = inject('configureResource', null)
const createResource = inject('createResource', null)

const cfg = kindConfig('module')
const modules = computed(() => (props.project.resources || []).filter(r => r.kind === 'module'))

const fieldSummary = m => {
  const n = (m.fields || []).length
  return n === 1 ? '1 field' : `${n} fields`
}
</script>
