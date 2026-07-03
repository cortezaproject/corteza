<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '40rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <span
          class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
          :class="[cfg.bg, cfg.ring]"
        >
          <i :class="[cfg.icon, cfg.text]" />
        </span>
        <div class="min-w-0">
          <div class="text-[10px] uppercase tracking-wider text-muted-color leading-none mb-0.5">
            {{ $t('project.kinds.user.single') }}
          </div>
          <div class="font-semibold truncate leading-tight">
            {{ displayName || $t('project.userDetail.unnamed') }}
          </div>
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-6">
      <!-- Read-only profile summary from the user directory. -->
      <div class="rounded-lg border border-surface divide-y divide-surface">
        <div class="flex items-center justify-between gap-4 px-4 py-2.5">
          <span class="text-xs uppercase tracking-wider text-muted-color">
            {{ $t('project.userDetail.name') }}
          </span>
          <span class="text-sm font-medium truncate">{{ displayName || '—' }}</span>
        </div>
        <div class="flex items-center justify-between gap-4 px-4 py-2.5">
          <span class="text-xs uppercase tracking-wider text-muted-color">
            {{ $t('project.userDetail.email') }}
          </span>
          <span class="text-sm text-muted-color truncate">{{ email || '—' }}</span>
        </div>
      </div>

      <!-- Project role assignment (staged; committed on Save). -->
      <CFormGroup :label="$t('project.userDetail.roles')">
        <p class="text-sm text-muted-color mb-2">{{ $t('project.userDetail.rolesHint') }}</p>
        <div v-if="roles.length" class="flex flex-wrap gap-1.5">
          <button
            v-for="role in roles"
            :key="role.id"
            type="button"
            class="px-2.5 py-1 rounded-md text-xs font-medium border transition-colors"
            :class="
              draft.roleIds.includes(role.id)
                ? 'bg-primary border-primary text-primary-contrast'
                : 'bg-transparent border-surface text-muted-color hover:border-primary'
            "
            @click="toggleRole(role.id)"
          >
            {{ role.name }}
          </button>
        </div>
        <p v-else class="text-sm text-muted-color italic">
          {{ $t('project.userDetail.noRoles') }}
        </p>
      </CFormGroup>
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="resourceId"
        :to="{ name: 'system.users.edit', params: { userID: resourceId } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.userDetail.openEditor')"
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
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Project user to inspect; this dialog only opens for existing project users.
  resourceId: { type: String, default: null },
})

const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const { t } = useI18n()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const cfg = kindConfig('user')

const roles = computed(() => store.rolesFor(props.project?.id))

// The project-user row carries the userId + the roles they currently hold.
const entity = computed(() =>
  props.resourceId
    ? store.projectUsersFor(props.project?.id).find(u => u.userId === props.resourceId)
    : null,
)

const directoryUser = computed(() =>
  props.resourceId ? usersStore.findUser(props.resourceId) : null,
)
const displayName = computed(() => directoryUser.value?.name || props.resourceId || '')
const email = computed(() => directoryUser.value?.email || '')

// --- Role assignment draft (staged; nothing persists until Save) --------------
const draft = reactive({ roleIds: [] })

function initDraft() {
  draft.roleIds = [...(entity.value?.roleIds || [])]
}

function toggleRole(roleId) {
  const i = draft.roleIds.indexOf(roleId)
  if (i === -1) draft.roleIds.push(roleId)
  else draft.roleIds.splice(i, 1)
}

// Load the roles + membership + directory this dialog needs on open. Best-effort;
// the body reads straight from the store caches these fill. Re-seed from the
// fresh copy afterwards, but never clobber in-progress edits.
async function loadContext(id) {
  if (!id) return
  try {
    await Promise.all([store.loadProjectUsers(id), usersStore.load()])
  } catch (err) {
    $toast.toastErrorHandler(t('project.userDetail.toastLoadFailed'))(err)
  }
}

watch(
  () => [props.modelValue, props.resourceId],
  async () => {
    if (!props.modelValue) return
    const openedFor = props.resourceId
    initDraft() // instant seed from the (possibly stale) cache
    const seeded = [...draft.roleIds]
    await loadContext(props.project?.id)
    if (
      props.modelValue &&
      props.resourceId === openedFor &&
      sameSet(draft.roleIds, seeded)
    ) {
      initDraft()
    }
  },
  { immediate: true },
)

function sameSet(a, b) {
  if (a.length !== b.length) return false
  const s = new Set(a)
  return b.every(x => s.has(x))
}

// --- Commit role changes ------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  saving.value = true
  try {
    const before = new Set(entity.value?.roleIds || [])
    const after = new Set(draft.roleIds)
    const toAdd = [...after].filter(id => !before.has(id))
    const toRemove = [...before].filter(id => !after.has(id))
    // setProjectUserRole refetches the list after each call; run sequentially so
    // the final store state reflects every toggle.
    for (const roleId of toAdd) {
      await store.setProjectUserRole(props.project.id, props.resourceId, roleId, true)
    }
    for (const roleId of toRemove) {
      await store.setProjectUserRole(props.project.id, props.resourceId, roleId, false)
    }
    emit('saved')
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.userDetail.toastSaveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
