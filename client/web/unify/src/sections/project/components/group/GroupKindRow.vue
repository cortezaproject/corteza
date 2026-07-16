<template>
  <div class="px-4 py-2">
    <div class="text-xs uppercase tracking-wider text-muted-color mb-1.5 flex items-center gap-1.5">
      <i :class="[cfg.icon, cfg.text]" />
      <span>{{ pluralLabel }}</span>
    </div>
    <div class="flex flex-wrap items-center gap-2">
      <ResourceBadge
        v-for="item in linkedItems"
        :key="item.id"
        :resource="item.resource"
        hide-icon
        clickable
        :removable="!disabled"
        @click="configureResource?.(item.id)"
        @remove="$emit('unlink', item.id)"
      />

      <button
        v-if="!disabled"
        type="button"
        class="shrink-0 flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-dashed border-surface text-muted-color hover:border-primary hover:text-primary transition-colors text-sm"
        @click="openPicker"
      >
        <i class="pi pi-plus text-xs" />
        {{ t('project.groups.kindRow.add', { kind: singularLabel }) }}
      </button>

      <span v-if="disabled && !linkedItems.length" class="text-sm text-muted-color italic">
        {{ t('project.groups.kindRow.empty') }}
      </span>

      <Popover ref="pickerRef" :pt="{ content: { class: 'p-0' } }">
        <div class="w-64 flex flex-col">
          <div class="flex-1 max-h-72 overflow-auto py-2">
            <button
              v-for="opt in addable"
              :key="opt.id"
              type="button"
              class="w-full px-3 py-1.5 flex items-center gap-2 hover:bg-emphasis text-left"
              @click="link(opt.id)"
            >
              <i :class="[cfg.icon, cfg.text]" />
              <span class="text-sm">{{ opt.name }}</span>
            </button>
            <div
              v-if="!addable.length"
              class="px-3 py-3 text-sm text-muted-color italic text-center"
            >
              {{ t('project.groups.kindRow.pickerEmpty', { kind: pluralLabel.toLowerCase() }) }}
            </div>
          </div>
          <div v-if="allowCreate" class="border-t border-surface p-2 bg-emphasis">
            <Button
              icon="pi pi-plus"
              :label="t('project.groups.kindRow.createNew', { kind: singularLabel })"
              severity="secondary"
              size="small"
              class="w-full"
              @click="openCreate"
            />
          </div>
        </div>
      </Popover>
    </div>

    <Dialog
      v-model:visible="createOpen"
      modal
      :header="t('project.groups.kindRow.newTitle', { kind: singularLabel })"
      :style="{ width: '24rem' }"
      :pt="{ footer: { class: 'flex justify-end gap-2' } }"
    >
      <InputText
        v-model="createName"
        :placeholder="t('project.groups.kindRow.namePlaceholder', { kind: singularLabel })"
        fluid
        autofocus
        @keydown.enter="create"
      />
      <template #footer>
        <Button :label="t('general.label.cancel')" severity="secondary" text size="small" @click="createOpen = false" />
        <Button :label="t('project.groups.kindRow.createLink')" size="small" :disabled="!createName.trim()" @click="create" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import ResourceBadge from '@/sections/project/components/ResourceBadge.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Provided by ArchitectureStep: opens the link/permission dialog for a resource.
const configureResource = inject('configureResource', null)

const props = defineProps({
  project: { type: Object, required: true },
  kind: { type: String, required: true },
  linkedIds: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
  // A resource id to keep out of the picker (e.g. don't let a resource link to
  // itself). Null = no exclusion.
  excludeId: { type: String, default: null },
  // Show the "Create new …" shortcut in the picker. Links use existing-only.
  allowCreate: { type: Boolean, default: true },
})

const emit = defineEmits(['link', 'unlink'])

const store = useProjectsStore()

const pickerRef = ref(null)
const createOpen = ref(false)
const createName = ref('')

const cfg = computed(() => kindConfig(props.kind))
// kindConfig only carries i18n keys (labelKey plural, singularKey singular);
// resolve them here for display.
const singularLabel = computed(() => t(cfg.value.singularKey).toLowerCase())
const pluralLabel = computed(() => t(cfg.value.labelKey))

const allOfKind = computed(() =>
  (props.project?.resources || []).filter(r => r.kind === props.kind),
)

const linkedItems = computed(() =>
  props.linkedIds
    .map(id => allOfKind.value.find(r => r.id === id))
    .filter(Boolean)
    .map(r => ({ id: r.id, resource: r })),
)

const addable = computed(() => {
  const linked = new Set(props.linkedIds)
  return allOfKind.value.filter(r => !linked.has(r.id) && r.id !== props.excludeId)
})

const openPicker = ev => pickerRef.value?.toggle(ev)

const link = id => {
  emit('link', id)
  pickerRef.value?.hide()
}

const openCreate = () => {
  pickerRef.value?.hide()
  createName.value = ''
  createOpen.value = true
}

const create = () => {
  if (!createName.value.trim()) return
  const id = store.addResource(props.project.projectID, { kind: props.kind, name: createName.value.trim() })
  if (id) emit('link', id)
  createOpen.value = false
}
</script>
