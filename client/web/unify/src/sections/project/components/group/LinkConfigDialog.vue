<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '72rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-end gap-2' } }"
  >
    <template #header>
      <div v-if="kind" class="flex items-center gap-2.5 min-w-0">
        <span
          class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
          :class="[srcCfg.bg, srcCfg.ring]"
        >
          <i :class="[srcCfg.icon, srcCfg.text]" />
        </span>
        <div class="min-w-0">
          <div class="text-[10px] uppercase tracking-wider text-muted-color leading-none mb-0.5">
            {{ kindLabel }}
          </div>
          <div class="font-semibold truncate leading-tight">{{ draft.name || 'Unnamed' }}</div>
        </div>
      </div>
    </template>

    <div v-if="kind" class="flex flex-col gap-5">
      <!-- Module: name + data sensitivity side by side -->
      <div v-if="isModule" class="grid grid-cols-1 sm:grid-cols-2 gap-x-5 gap-y-4">
        <CFormGroup label="Name">
          <InputText v-model="draft.name" size="small" fluid :disabled="readonly" autofocus />
        </CFormGroup>
        <CFormGroup label="Data sensitivity" description="Default for new fields">
          <Select
            v-model="draft.sensitivity"
            :options="SENSITIVITY_OPTIONS"
            option-label="label"
            option-value="id"
            size="small"
            fluid
            :disabled="readonly"
          />
        </CFormGroup>
      </div>
      <!-- Connection: sensitivity only (name comes from the connector) -->
      <CFormGroup v-else-if="isConnection" label="Data sensitivity">
        <Select
          v-model="draft.sensitivity"
          :options="SENSITIVITY_OPTIONS"
          option-label="label"
          option-value="id"
          size="small"
          class="w-72"
          :disabled="readonly"
        />
      </CFormGroup>
      <!-- Other resources: name -->
      <CFormGroup v-else label="Name">
        <InputText v-model="draft.name" size="small" fluid :disabled="readonly" autofocus />
      </CFormGroup>

      <!-- Role: permissions across resources (existing roles only) -->
      <CFormGroup v-if="isRole && !isNew" label="Permissions">
        <PermList
          :rows="permTargets"
          :get-level="r => draftLevel(resourceId, r.id)"
          :baseline="r => isBaseline(project, resourceId, r.id)"
          :disabled="readonly"
          @set="(r, level) => setDraftBinding(resourceId, r.id, level)"
        />
      </CFormGroup>

      <!-- Non-role: fields (modules) or links (others), plus who can access it -->
      <template v-else-if="!isRole">
        <CFormGroup v-if="isModule" label="Fields">
          <template #actions>
            <Button
              v-if="!readonly"
              icon="pi pi-plus"
              label="Add field"
              severity="secondary"
              size="small"
              @click="addDraftField"
            />
          </template>
          <FieldsEditor v-model="draft.fields" :module-options="moduleOptions" :disabled="readonly" />
        </CFormGroup>

        <!-- Plain resources link to others; modules infer links from record fields,
             connections have no links of their own. -->
        <CFormGroup v-if="!isModule && !isConnection" label="Links to">
          <div class="rounded-lg border border-surface divide-y divide-surface py-1">
            <GroupKindRow
              v-for="k in LINK_KINDS"
              :key="k"
              :project="project"
              :kind="k"
              :linked-ids="draft.linkedIds"
              :exclude-id="resourceId"
              :allow-create="false"
              :disabled="readonly"
              @link="linkDraft"
              @unlink="unlinkDraft"
            />
          </div>
        </CFormGroup>

        <CFormGroup v-if="!isNew && !isConnection && roles.length" label="Roles & permissions">
          <PermList
            :rows="roles"
            :get-level="r => draftLevel(r.id, resourceId)"
            :baseline="r => isBaseline(project, r.id, resourceId)"
            :disabled="readonly"
            @set="(r, level) => setDraftBinding(r.id, resourceId, level)"
          />
        </CFormGroup>
      </template>
    </div>

    <template #footer>
      <Button v-if="readonly" label="Close" severity="secondary" text size="small" @click="visible = false" />
      <template v-else>
        <Button label="Cancel" severity="secondary" text size="small" @click="visible = false" />
        <Button label="Save" size="small" :disabled="!draft.name.trim()" @click="onSave" />
      </template>
    </template>
  </Dialog>
</template>

<script setup>
import FieldsEditor from '@/sections/project/components/datamodel/FieldsEditor.vue'
import GroupKindRow from '@/sections/project/components/group/GroupKindRow.vue'
import PermList from '@/sections/project/components/group/PermList.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { SENSITIVITY_OPTIONS } from '@/sections/project/config/sensitivity'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { linkedResourceIds } from '@/sections/project/utils/links'
import { effectiveLevel, isBaseline, isPermissionTarget } from '@/sections/project/utils/rbac'
import { computed, reactive, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Existing resource to edit…
  resourceId: { type: String, default: null },
  // …or a kind to create a brand-new resource (staged until Save).
  newKind: { type: String, default: null },
  // For a new connection: the catalog connector id and a default name.
  newConnector: { type: String, default: null },
  newName: { type: String, default: '' },
  // Open as read-only (e.g. approver view, or a submitted/approved step).
  readonly: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const store = useProjectsStore()

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const LINK_KINDS = ['module', 'page', 'chart', 'automation', 'agent', 'chatbot', 'connection']

const isNew = computed(() => !!props.newKind && !props.resourceId)
const resource = computed(() =>
  props.resourceId ? (props.project?.resources || []).find(r => r.id === props.resourceId) : null,
)
const kind = computed(() => (isNew.value ? props.newKind : resource.value?.kind))
const isRole = computed(() => kind.value === 'role')
const isModule = computed(() => kind.value === 'module')
const isConnection = computed(() => kind.value === 'connection')
const srcCfg = computed(() => kindConfig(kind.value))
const kindLabel = computed(() => srcCfg.value.label.replace(/s$/, ''))

const allResources = computed(() => props.project?.resources || [])
const roles = computed(() => allResources.value.filter(r => r.kind === 'role'))
const permTargets = computed(() => allResources.value.filter(r => isPermissionTarget(r.kind)))
const moduleOptions = computed(() =>
  allResources.value.filter(r => r.kind === 'module' && r.id !== props.resourceId),
)

// --- Draft (staged) edits; nothing touches the store until Save -------------
const draft = reactive({ name: '', fields: [], linkedIds: [], bindings: {}, sensitivity: null })
// Once a new resource is created on Save, remember its id so a re-save (before
// the dialog closes) updates it instead of creating another copy.
const createdId = ref(null)

function initDraft() {
  createdId.value = null
  if (isNew.value) {
    draft.name = props.newName || ''
    draft.fields = []
    draft.linkedIds = []
    draft.sensitivity = null
  } else {
    const r = resource.value
    draft.name = r?.name || ''
    draft.fields = (r?.fields || []).map(f => ({ ...f }))
    draft.linkedIds = linkedResourceIds(props.project, props.resourceId)
    draft.sensitivity = r?.sensitivity || null
  }
  draft.bindings = {}
}

watch(
  () => [props.modelValue, props.resourceId, props.newKind],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

// Permission draft helpers (key: `${roleId}|${resourceId}`).
const draftLevel = (roleId, resId) => {
  const k = `${roleId}|${resId}`
  return k in draft.bindings ? draft.bindings[k] : effectiveLevel(props.project, roleId, resId)
}
const setDraftBinding = (roleId, resId, level) => {
  draft.bindings[`${roleId}|${resId}`] = level
}

// Field draft helpers (modules).
const nid = () => `f-${Math.random().toString(36).slice(2, 9)}`
const addDraftField = () => {
  // New fields inherit the module's sensitivity as their default.
  draft.fields.push({
    id: nid(),
    name: '',
    type: 'String',
    required: false,
    targetModuleId: null,
    sensitivity: draft.sensitivity || null,
  })
}

// Link draft helpers.
const linkDraft = id => {
  if (!draft.linkedIds.includes(id)) draft.linkedIds.push(id)
}
const unlinkDraft = id => {
  draft.linkedIds = draft.linkedIds.filter(x => x !== id)
}

// --- Commit -----------------------------------------------------------------
function onSave() {
  if (!draft.name.trim()) return
  const pid = props.project.projectID
  const k = kind.value
  try {
    // Edit the existing resource, the already-created one, or create a new one.
    const patch = { name: draft.name }
    if (isConnection.value || isModule.value) patch.sensitivity = draft.sensitivity || ''
    let id = props.resourceId || createdId.value
    if (id) {
      store.updateResource(pid, id, patch)
    } else {
      id = store.addResource(pid, { kind: k, name: draft.name, connector: props.newConnector || resource.value?.connector })
      createdId.value = id
      if (id) store.updateResource(pid, id, patch)
    }
    if (!id) return

    if (k === 'module') {
      store.setFields(pid, id, draft.fields)
    } else if (!isRole.value && !isConnection.value) {
      const current = new Set(linkedResourceIds(props.project, id))
      const next = new Set(draft.linkedIds)
      for (const t of next) if (!current.has(t)) store.addLink(pid, id, t)
      for (const t of current) if (!next.has(t)) store.removeLink(pid, id, t)
    }

    for (const [key, level] of Object.entries(draft.bindings)) {
      const [roleId, resId] = key.split('|')
      store.setBinding(pid, roleId, resId, level)
    }
  } finally {
    // Always close, even if a write hiccups — prevents re-save duplicates.
    visible.value = false
  }
}
</script>
