<template>
  <div class="h-full overflow-auto p-4">
    <div>
      <div class="flex flex-col gap-2">
        <div v-if="!disabled">
          <Button
            icon="pi pi-plus"
            :label="$t('project.dataModel.addModule')"
            size="small"
            @click="openCreate"
          />
        </div>

        <CEmptyState v-if="!modules.length">
          {{ $t('project.dataModel.empty') }}
        </CEmptyState>

        <!-- One card per module — click the header to configure it; each card
             lists the module's fields with their type. -->
        <div v-else class="flex flex-col gap-4">
          <div
            v-for="m in modules"
            :key="m.id"
            class="bg-surface border border-surface rounded-border shadow-sm overflow-hidden"
          >
            <div class="group flex items-center hover:bg-emphasis transition-colors">
              <button
                type="button"
                class="self-stretch flex items-center px-4 shrink-0 cursor-pointer text-muted-color hover:bg-emphasis transition-colors"
                :aria-label="$t(isCollapsed(m.id) ? 'general.label.expand' : 'general.label.collapse')"
                :title="$t(isCollapsed(m.id) ? 'general.label.expand' : 'general.label.collapse')"
                @click="toggle(m.id)"
              >
                <i
                  class="pi pi-chevron-down text-xs transition-transform duration-200"
                  :class="{ '-rotate-90': isCollapsed(m.id) }"
                />
              </button>
              <button
                type="button"
                class="flex-1 min-w-0 flex items-center gap-3 py-3 pr-3 pl-1 text-left cursor-pointer"
                :aria-label="$t('general.label.edit')"
                @click="inspectResource?.('module', m.id)"
              >
                <CFormItemContent :title="m.name" :subtitle="fieldSummary(m)" />
              </button>
              <CRouterLinkButton
                :to="{ name: 'admin.modules.edit', params: { slug: project.namespaceID, moduleID: m.id } }"
                icon="pi pi-external-link"
                severity="secondary"
                text
                size="small"
                class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity mr-1"
                :aria-label="$t('project.dataModel.openEditor')"
                :title="$t('project.dataModel.openEditor')"
                @click.stop
              />
              <Button
                v-if="!disabled"
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity mr-2"
                :aria-label="$t('general.label.remove')"
                :title="$t('general.label.remove')"
                @click="removeModule(m)"
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
                    class="group flex items-center gap-3 px-3 min-h-11 hover:bg-emphasis transition-colors cursor-pointer"
                    @click="editField?.(m.id, f.id)"
                  >
                    <div class="min-w-0 flex-1 flex items-center gap-2">
                      <span class="text-sm truncate">
                        {{ f.name || $t('project.dataModel.untitledField') }}
                      </span>
                      <FieldKindTag :type="f.type" />
                      <template v-if="f.required">
                        <span class="shrink-0 text-muted-color opacity-50 text-[11px]">·</span>
                        <span
                          class="shrink-0 text-[11px] font-medium text-amber-600 dark:text-amber-400"
                        >
                          {{ $t('project.dataModel.required') }}
                        </span>
                      </template>
                      <template v-if="f.multi">
                        <span class="shrink-0 text-muted-color opacity-50 text-[11px]">·</span>
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
      </div>
    </div>
  </div>
</template>

<script setup>
import FieldKindTag from '@/sections/project/components/datamodel/FieldKindTag.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CEmptyState, CRouterLinkButton } = components

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const { t } = useI18n()
const store = useProjectsStore()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)
const editField = inject('editField', null)
const createField = inject('createField', null)

const modules = computed(() => store.resourcesFor(props.project.projectID).filter(r => r.kind === 'module'))

// Modules collapsed by default; the set holds the expanded ones.
const expandedIds = ref(new Set())
const isCollapsed = id => !expandedIds.value.has(id)
const toggle = id => {
  const next = new Set(expandedIds.value)
  next.has(id) ? next.delete(id) : next.add(id)
  expandedIds.value = next
}

// "Add Module" arms a flag; the module that then appears (diffed against the
// previous list) starts expanded. Without the flag the initial load — which
// also fills an empty list — would look like a create.
const expandNext = ref(false)
const openCreate = () => {
  expandNext.value = true
  createResource?.('module')
}
watch(modules, (list, prev) => {
  if (!expandNext.value) return
  const known = new Set((prev || []).map(m => m.id))
  const fresh = list.find(m => !known.has(m.id))
  if (fresh) {
    expandedIds.value = new Set(expandedIds.value).add(fresh.id)
    expandNext.value = false
  }
})

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
    await store.removeResource(props.project.projectID, m.id)
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
    await store.removeField(props.project.projectID, m.id, f.id)
    $toast.toastSuccess(f.name, t('project.dataModel.toast.fieldRemoved'))
  } catch (err) {
    $toast.toastErrorHandler(t('project.dataModel.toast.fieldRemoveFailed'))(err)
  }
}
</script>
