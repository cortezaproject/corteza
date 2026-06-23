<template>
  <div class="h-full overflow-auto p-4">
    <div>
      <CFormGroup :label="$t('project.dataModel.label')">
        <template #actions>
          <Button
            v-if="!disabled"
            icon="pi pi-plus"
            :label="$t('project.dataModel.addModule')"
            severity="secondary"
            size="small"
            @click="createResource?.('module')"
          />
        </template>

        <div
          v-if="!modules.length"
          class="text-muted-color p-4 border rounded-lg bg-emphasis text-center"
        >
          {{ $t('project.dataModel.empty') }}
        </div>

        <!-- One card per module — click the header to configure it; each card
             lists the module's fields with their type. -->
        <div v-else class="flex flex-col gap-4 mt-1">
          <div
            v-for="m in modules"
            :key="m.id"
            class="border border-surface rounded-border shadow-sm overflow-hidden"
          >
            <div
              class="group flex items-center gap-3 p-3 cursor-pointer hover:bg-emphasis transition-colors"
              @click="toggle(m.id)"
            >
              <i
                class="pi pi-chevron-down text-xs text-muted-color shrink-0 transition-transform duration-200"
                :class="{ '-rotate-90': isCollapsed(m.id) }"
              />
              <span
                class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
                :class="[cfg.bg, cfg.ring]"
              >
                <i :class="[cfg.icon, cfg.text]" />
              </span>
              <div class="min-w-0">
                <div class="font-medium truncate">{{ m.name }}</div>
                <div class="text-xs text-muted-color">{{ fieldSummary(m) }}</div>
              </div>
              <Button
                v-if="!disabled"
                icon="pi pi-pencil"
                severity="secondary"
                text
                size="small"
                class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                :aria-label="$t('general.label.edit')"
                :title="$t('general.label.edit')"
                @click.stop="configureResource?.(m.id)"
              />
              <div class="flex-1" />
              <Button
                v-if="!disabled"
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                :aria-label="$t('general.label.remove')"
                :title="$t('general.label.remove')"
                @click.stop="removeModule(m)"
              />
            </div>
            <!-- Animated collapse via grid-template-rows 0fr <-> 1fr -->
            <div
              class="grid transition-[grid-template-rows] duration-200 ease-in-out"
              :class="isCollapsed(m.id) ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'"
            >
              <div class="overflow-hidden min-h-0">
                <div class="divide-y divide-surface border-t border-surface">
                  <div
                    v-for="f in m.fields || []"
                    :key="f.id"
                    class="group flex items-center gap-3 px-3 min-h-[2.75rem] hover:bg-emphasis transition-colors cursor-pointer"
                    @click="editField?.(m.id, f.id)"
                  >
                    <div class="min-w-0 flex-1 flex items-center gap-2">
                      <span class="text-sm truncate">
                        {{ f.name || $t('project.dataModel.untitledField') }}
                      </span>
                      <span
                        class="shrink-0 text-[11px] font-medium px-1.5 py-0.5 rounded bg-emphasis text-muted-color border border-surface"
                      >
                        {{ fieldTypeText(f.type) }}
                      </span>
                      <template v-if="f.required">
                        <span class="shrink-0 text-muted-color/50 text-[11px]">·</span>
                        <span
                          class="shrink-0 text-[11px] font-medium text-amber-600 dark:text-amber-400"
                        >
                          {{ $t('project.dataModel.required') }}
                        </span>
                      </template>
                      <template v-if="f.multi">
                        <span class="shrink-0 text-muted-color/50 text-[11px]">·</span>
                        <span
                          class="shrink-0 text-[11px] font-medium text-sky-600 dark:text-sky-400"
                        >
                          {{ $t('project.dataModel.multiple') }}
                        </span>
                      </template>
                    </div>
                    <Button
                      v-if="!disabled"
                      icon="pi pi-trash"
                      severity="danger"
                      text
                      size="small"
                      class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                      :aria-label="$t('general.label.remove')"
                      :title="$t('general.label.remove')"
                      @click.stop="removeField(m, f)"
                    />
                  </div>
                  <!-- No fields: a message with the add-field action under it
                       (the action is hidden on read-only steps). -->
                  <div
                    v-if="!(m.fields || []).length"
                    class="px-3 py-5 flex flex-col items-center gap-3"
                  >
                    <span class="text-sm text-muted-color italic">
                      {{ $t('project.dataModel.emptyFields') }}
                    </span>
                    <Button
                      v-if="!disabled"
                      icon="pi pi-plus"
                      :label="$t('project.dataModel.addField')"
                      severity="secondary"
                      size="small"
                      @click="createField?.(m.id)"
                    />
                  </div>
                  <!-- With fields: the add-field action sits below the list. -->
                  <div v-else-if="!disabled" class="p-3">
                    <Button
                      icon="pi pi-plus"
                      :label="$t('project.dataModel.addField')"
                      severity="secondary"
                      size="small"
                      @click="createField?.(m.id)"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CFormGroup>
    </div>
  </div>
</template>

<script setup>
import { fieldTypeLabelKey } from '@/sections/project/config/fieldTypes'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const { t } = useI18n()
const store = useProjectsStore()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()
const configureResource = inject('configureResource', null)
const createResource = inject('createResource', null)
const editField = inject('editField', null)
const createField = inject('createField', null)

const cfg = kindConfig('module')
const modules = computed(() => store.resourcesFor(props.project.id).filter(r => r.kind === 'module'))

// Per-module collapse state (expanded by default); keyed by module id.
const collapsedModules = ref({})
const isCollapsed = id => !!collapsedModules.value[id]
const toggle = id => {
  collapsedModules.value[id] = !collapsedModules.value[id]
}

// Localized field-type label, falling back to the raw kind for unknown types.
const fieldTypeText = type => {
  const key = fieldTypeLabelKey(type)
  return key ? t(key) : type
}

const fieldSummary = m => {
  const n = (m.fields || []).length
  if (!n) return t('project.field.noFields')
  return n === 1 ? t('project.field.one') : t('project.field.many', { count: n })
}

function removeModule(m) {
  confirmDelete({
    header: t('project.dataModel.removeModule.header'),
    message: t('project.dataModel.removeModule.message', { name: m.name }),
    onConfirm: () => handleRemoveModule(m),
  })
}

async function handleRemoveModule(m) {
  try {
    await store.removeResource(props.project.id, m.id)
    $toast.toastSuccess(m.name, t('project.dataModel.toast.removed'))
  } catch (err) {
    $toast.toastErrorHandler(t('project.dataModel.toast.removeFailed'))(err)
  }
}

function removeField(m, f) {
  confirmDelete({
    header: t('project.dataModel.removeField.header'),
    message: t('project.dataModel.removeField.message', {
      name: f.name || t('project.dataModel.untitledField'),
    }),
    onConfirm: () => handleRemoveField(m, f),
  })
}

async function handleRemoveField(m, f) {
  try {
    await store.removeField(props.project.id, m.id, f.id)
    $toast.toastSuccess(f.name, t('project.dataModel.toast.fieldRemoved'))
  } catch (err) {
    $toast.toastErrorHandler(t('project.dataModel.toast.fieldRemoveFailed'))(err)
  }
}
</script>
