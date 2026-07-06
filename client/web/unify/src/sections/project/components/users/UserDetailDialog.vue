<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '64rem', maxWidth: '96vw' }"
    :pt="{
      content: { class: '!pt-2', style: 'max-height: 75vh; overflow: auto' },
      footer: { class: 'flex justify-between gap-2 p-3' },
    }"
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
      <!-- Project role assignment. Toggling a chip applies immediately (adds or
           removes the user from that role); the permissions matrix below reacts
           live — the role's column appears/disappears and the evaluated column
           re-traces. -->
      <CFormGroup :label="$t('project.userDetail.roles')">
        <p class="text-sm text-muted-color mb-2">{{ $t('project.userDetail.rolesHint') }}</p>
        <div v-if="roles.length" class="flex flex-wrap gap-2">
          <button
            v-for="role in roles"
            :key="role.id"
            type="button"
            class="inline-flex items-center gap-2 pl-1.5 pr-3 py-1.5 rounded-lg border border-surface text-sm font-medium transition-all"
            :class="
              draft.roleIds.includes(role.id)
                ? 'text-color'
                : 'text-muted-color opacity-50 hover:opacity-100'
            "
            @click="toggleRole(role.id)"
          >
            <!-- Same role badge as the matrix column headers in both states; the
                 whole chip just dims when the role isn't held. -->
            <span
              class="inline-flex items-center justify-center w-5 h-5 rounded ring-1 shrink-0"
              :class="[roleCfg.bg, roleCfg.ring]"
            >
              <i :class="[roleCfg.icon, roleCfg.text, 'text-[10px]']" />
            </span>
            {{ role.name }}
          </button>
        </div>
        <p v-else class="text-sm text-muted-color italic">
          {{ $t('project.userDetail.noRoles') }}
        </p>
      </CFormGroup>

      <!-- Effective-permissions overview. The leading column is this user's
           resolved (read-only) access across every role they hold; the role
           columns are those same roles and toggle their permissions exactly
           like the permissions step. Columns track the assigned-role chips above
           live — toggling a chip shows/hides its column and re-evaluates. -->
      <CFormGroup :label="$t('project.userDetail.permissions.label')">
        <p class="text-sm text-muted-color mb-2">
          {{ $t('project.userDetail.permissions.hint') }}
        </p>
        <!-- Always rendered: even with no roles held, the evaluated column shows
             the user's resolved (denied) access. Role columns appear as roles
             are toggled on above. -->
        <ProjectPermissionMatrix
          :project="project"
          :roles="membershipRoles"
          :eval-user-id="resourceId"
          :eval-label="displayName"
        />
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

      <!-- Everything in this dialog (role membership + permissions) applies
           live, so the footer only closes. -->
      <Button
        :label="$t('general.label.close')"
        severity="secondary"
        size="small"
        @click="visible = false"
      />
    </template>
  </Dialog>
</template>

<script setup>
import ProjectPermissionMatrix from '@/sections/project/components/permissions/ProjectPermissionMatrix.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, watch } from 'vue'
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
// Role kind visuals (violet id-card) — used to style the role toggle chips so
// they read as roles, matching how roles look elsewhere.
const roleCfg = kindConfig('role')

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

// --- Role assignment (applies live) -------------------------------------------
// `draft` mirrors the user's current membership. It's seeded from the store on
// open and kept in sync by optimistic toggles (rolled back if a write fails).
const draft = reactive({ roleIds: [] })

// Permission-matrix columns track the assigned roles live, so toggling a chip
// shows/hides that role's column immediately. The evaluated column re-traces on
// its own once the membership write lands (loadProjectUsers → touch bumps
// graphVersion, which the matrix watches).
const membershipRoles = computed(() =>
  roles.value.filter(r => draft.roleIds.includes(r.id)),
)

function initDraft() {
  draft.roleIds = [...(entity.value?.roleIds || [])]
}

// Toggling adds/removes the user from the role right away. Optimistically flip
// the chip + column, then persist; roll back on failure.
async function toggleRole(roleId) {
  const on = !draft.roleIds.includes(roleId)
  if (on) draft.roleIds.push(roleId)
  else draft.roleIds.splice(draft.roleIds.indexOf(roleId), 1)
  try {
    await store.setProjectUserRole(props.project.id, props.resourceId, roleId, on)
    emit('saved')
  } catch (err) {
    if (on) draft.roleIds.splice(draft.roleIds.indexOf(roleId), 1)
    else draft.roleIds.push(roleId)
    $toast.toastErrorHandler(t('project.userDetail.toastSaveFailed'))(err)
  }
}

// Load the roles + membership + directory this dialog needs on open. Best-effort;
// the body reads straight from the store caches these fill. Re-seed from the
// fresh copy afterwards, but never clobber in-progress edits.
async function loadContext(id) {
  if (!id) return
  try {
    // Membership + directory for the header/roles, plus the resource caches the
    // permission matrix charts (modules, pages, agents, chatbots).
    await Promise.all([
      store.loadProjectUsers(id),
      store.loadResources(id),
      store.loadPages(id),
      store.loadAgents(id),
      store.loadChatbots(id),
      usersStore.load(),
    ])
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
</script>
