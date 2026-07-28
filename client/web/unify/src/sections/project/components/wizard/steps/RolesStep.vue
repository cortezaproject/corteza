<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div v-if="!disabled">
        <Button
          icon="pi pi-plus"
          :label="$t('project.accessRoles.add')"
          size="small"
          @click="createResource?.('role')"
        />
      </div>

      <CEmptyState v-if="!roles.length">
        {{ $t('project.accessRoles.empty') }}
      </CEmptyState>

      <!-- One card per role — click the header to configure it; each card lists
           the role's members. Mirrors the data-model step's module cards. -->
      <div v-else class="flex flex-col gap-4">
        <div
          v-for="r in roles"
          :key="r.id"
          class="bg-surface border border-surface rounded-border shadow-sm overflow-hidden"
        >
          <div class="group flex items-center hover:bg-emphasis transition-colors">
            <button
              type="button"
              class="self-stretch flex items-center px-4 shrink-0 cursor-pointer text-muted-color hover:bg-emphasis transition-colors"
              :aria-label="
                $t(isCollapsed(r.id) ? 'general.label.expand' : 'general.label.collapse')
              "
              :title="$t(isCollapsed(r.id) ? 'general.label.expand' : 'general.label.collapse')"
              @click="toggle(r.id)"
            >
              <i
                class="pi pi-chevron-down text-xs transition-transform duration-200"
                :class="{ '-rotate-90': isCollapsed(r.id) }"
              />
            </button>
            <button
              type="button"
              class="flex-1 min-w-0 flex items-center gap-3 py-3 pr-3 pl-1 text-left cursor-pointer"
              :aria-label="$t('general.label.edit')"
              @click="inspectResource?.('role', r.id)"
            >
              <CFormItemContent :title="r.name">
                <template #subtitle>
                  <div class="text-xs text-muted-color truncate">
                    <template v-if="r.description">{{ r.description }} ·</template>
                    {{ memberSummary(r) }}
                  </div>
                </template>
              </CFormItemContent>
            </button>
            <CRouterLinkButton
              :to="{ name: 'system.roles.edit', params: { roleID: r.id } }"
              icon="pi pi-external-link"
              severity="secondary"
              text
              size="small"
              class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity mr-1"
              :aria-label="$t('project.accessRoles.openEditor')"
              :title="$t('project.accessRoles.openEditor')"
              @click.stop
            />
            <Button
              v-if="!disabled"
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity mr-2"
              :aria-label="$t('project.accessRoles.remove')"
              :title="$t('project.accessRoles.remove')"
              @click="onRemove(r)"
            />
          </div>
          <!-- Animated collapse via grid-template-rows 0fr <-> 1fr -->
          <div
            class="grid transition-[grid-template-rows] duration-200 ease-in-out"
            :class="isCollapsed(r.id) ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'"
          >
            <div class="overflow-hidden min-h-0">
              <div class="border-t border-surface p-3">
                <RoleMemberList :project="project" :role="r" :disabled="disabled" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import RoleMemberList from '@/sections/project/components/roles/RoleMemberList.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CEmptyState, CRouterLinkButton } = components

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

// The Wizard mounts the create/detail dialogs once and hands these openers down.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

const roles = computed(() => store.rolesFor(props.project.projectID))
const projectUsers = computed(() => store.projectUsersFor(props.project.projectID))

// Roles collapsed by default; the set holds the expanded ones.
const expandedIds = ref(new Set())
const isCollapsed = id => !expandedIds.value.has(id)
const toggle = id => {
  const next = new Set(expandedIds.value)
  next.has(id) ? next.delete(id) : next.add(id)
  expandedIds.value = next
}

const memberCount = role => projectUsers.value.filter(u => u.roleIds.includes(role.id)).length
const memberSummary = role => {
  const n = memberCount(role)
  if (!n) return t('project.roleMembers.none')
  return n === 1 ? t('project.roleMembers.one') : t('project.roleMembers.many', { count: n })
}

// A freshly created role (diffed against the previous list) starts expanded,
// matching the data-model step's create behaviour. The create dialog is owned by
// the Wizard now; we watch the store list for the new entry rather than owning it.
watch(roles, (list, prev) => {
  if (!prev || !prev.length) return // ignore the initial hydration burst
  const known = new Set(prev.map(r => r.id))
  const fresh = list.find(r => !known.has(r.id))
  if (fresh) expandedIds.value = new Set(expandedIds.value).add(fresh.id)
})

// --- remove ------------------------------------------------------------------
function onRemove(r) {
  confirmDelete({
    header: t('project.accessRoles.removeConfirm.header'),
    message: t('project.accessRoles.removeConfirm.message', { name: r.name }),
    onConfirm: async () => {
      try {
        await store.removeRole(props.project.projectID, r.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.accessRoles.toastRemoveFailed'))(err)
      }
    },
  })
}

// loadProjectUsers loads the roles too (member resolution needs them), so the
// member summaries and lists have data. The directory resolves names/emails.
async function refresh(id) {
  if (!id) return
  try {
    await Promise.all([store.loadProjectUsers(id), usersStore.load()])
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessRoles.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.projectID))
watch(
  () => props.project.projectID,
  id => id && refresh(id),
)
</script>
