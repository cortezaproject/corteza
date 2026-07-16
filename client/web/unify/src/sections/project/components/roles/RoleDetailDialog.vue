<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '52rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <KindIcon kind="role" size="lg" plain-icon />
        <div class="min-w-0">
          <DialogEyebrow>{{ $t('project.kinds.role.single') }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">
            {{ draft.name || $t('project.roleDetail.unnamed') }}
          </div>
        </div>
      </div>
    </template>

    <div v-if="role" class="flex flex-col gap-6">
      <!-- Role meta (staged; committed on Save) -->
      <div class="flex flex-col gap-5">
        <CFormGroup :label="$t('general.label.name')" required>
          <div>
            <InputText
              v-model="draft.name"
              size="small"
              fluid
              :invalid="submitted && !!nameError"
              autofocus
            />
            <ValidationMessage :message="submitted ? nameError : ''" />
          </div>
        </CFormGroup>
        <CFormGroup :label="$t('general.label.description')">
          <Textarea v-model="draft.description" rows="3" auto-resize fluid />
        </CFormGroup>
      </div>

      <!-- Members (persist immediately) -->
      <CFormGroup :label="$t('project.roleDetail.members')">
        <RoleMemberList :project="project" :role="role" class="mt-1" />
      </CFormGroup>

      <!-- Permissions overview — the full resource matrix for this one role. -->
      <CFormGroup :label="$t('project.roleDetail.permissions')">
        <p class="text-sm text-muted-color mb-2">{{ $t('project.roleDetail.permissionsHint') }}</p>
        <ProjectPermissionMatrix :project="project" :roles="[role]" hide-role-header />
      </CFormGroup>
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="resourceId"
        :to="{ name: 'system.roles.edit', params: { roleID: resourceId } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.roleDetail.openEditor')"
        icon="pi pi-external-link"
        severity="secondary"
        text
        size="small"
      />
      <span v-else />

      <div class="flex gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="visible = false"
        />
        <Button :label="$t('general.label.save')" size="small" :loading="saving" @click="onSave" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import ProjectPermissionMatrix from '@/sections/project/components/permissions/ProjectPermissionMatrix.vue'
import RoleMemberList from '@/sections/project/components/roles/RoleMemberList.vue'
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Existing access role to edit; always set when the dialog is opened.
  resourceId: { type: String, default: null },
})

const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})


const role = computed(() =>
  props.resourceId ? store.rolesFor(props.project?.projectID).find(r => r.id === props.resourceId) : null,
)

// --- Role meta draft (staged; nothing persists until Save) --------------------
const draft = reactive({ name: '', description: '' })

function initDraft() {
  submitted.value = false
  const r = role.value
  draft.name = r?.name || ''
  draft.description = r?.description || ''
}

// --- Validation ---------------------------------------------------------------
const submitted = ref(false)

const nameError = computed(() => (draft.name.trim() ? '' : t('project.roleDetail.nameRequired')))

const isValid = computed(() => !nameError.value)

// Load everything the members list and permission matrix need for this project
// (resources across kinds + role membership + the user directory). Best-effort;
// the matrix and member list read straight from the store caches these fill.
async function loadContext(id) {
  if (!id) return
  try {
    await Promise.all([
      store.loadResources(id),
      store.loadPages(id),
      store.loadAutomations(id),
      store.loadAgents(id),
      store.loadChatbots(id),
      store.loadConnections(id),
      store.loadProjectUsers(id),
      usersStore.load(),
    ])
  } catch (err) {
    $toast.toastErrorHandler(t('project.roleDetail.toastLoadFailed'))(err)
  }
}

watch(
  () => [props.modelValue, props.resourceId],
  async () => {
    if (!props.modelValue) return
    const openedFor = props.resourceId
    initDraft() // instant seed from the (possibly stale) cache
    const seeded = { name: draft.name, description: draft.description }
    await loadContext(props.project?.projectID)
    // loadContext refetches the roles; re-seed from the fresh copy, but only if
    // the dialog is still open for the same role and the user hasn't started
    // editing — never clobber in-progress input.
    if (
      props.modelValue &&
      props.resourceId === openedFor &&
      draft.name === seeded.name &&
      draft.description === seeded.description
    ) {
      initDraft()
    }
  },
  { immediate: true },
)

// --- Commit role meta ---------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    await store.updateRole(props.project.projectID, props.resourceId, {
      name: draft.name.trim(),
      description: draft.description,
    })
    emit('saved')
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.roleDetail.toastSaveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
