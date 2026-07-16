<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div v-if="!disabled && roles.length">
        <Button
          icon="pi pi-plus"
          :label="$t('project.accessUsers.add')"
          size="small"
          @click="createResource?.('user')"
        />
      </div>

      <!-- Need roles before users can be assigned to them. -->
      <CEmptyState v-if="!roles.length">
        {{ $t('project.accessUsers.noRoles') }}
      </CEmptyState>

      <CFormItemList
        v-else
        :items="rows"
        :empty-message="$t('project.accessUsers.empty')"
        :remove-label="$t('project.accessUsers.remove')"
        :hide-remove="disabled"
        item-key="userId"
        reveal-on-hover
        @select="row => inspectResource?.('user', row.userId)"
        @remove="remove"
      >
        <template #default="{ item: row }">
          <div class="flex items-center gap-4 min-w-0">
            <CFormItemContent :title="row.name" :subtitle="row.email" />
            <!-- Read-only role summary using the same role badge as the user
                 detail dialog; editing happens in the detail dialog. -->
            <div v-if="row.roleNames.length" class="flex flex-wrap gap-1.5">
              <KindBadge v-for="name in row.roleNames" :key="name" kind="role" :label="name" />
            </div>
            <span v-else class="text-xs text-muted-color italic">
              {{ $t('project.accessUsers.noRolesAssigned') }}
            </span>
          </div>
        </template>
        <template #hover-actions="{ item: row }">
          <CRouterLinkButton
            :to="{ name: 'system.users.edit', params: { userID: row.userId } }"
            icon="pi pi-external-link"
            severity="secondary"
            text
            rounded
            size="small"
            :aria-label="$t('project.accessUsers.openEditor')"
            :title="$t('project.accessUsers.openEditor')"
            @click.stop
          />
        </template>
      </CFormItemList>
    </div>
  </div>
</template>

<script setup>
import KindBadge from '@/sections/project/components/KindBadge.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, watch } from 'vue'
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

// Wizard mounts the user create/detail dialogs once and provides these openers.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

const roles = computed(() => store.rolesFor(props.project.projectID))
const projectUsers = computed(() => store.projectUsersFor(props.project.projectID))

const rows = computed(() =>
  projectUsers.value.map(u => {
    const dir = usersStore.findUser(u.userId)
    return {
      userId: u.userId,
      roleIds: u.roleIds,
      roleNames: u.roleIds
        .map(id => roles.value.find(r => r.id === id)?.name)
        .filter(Boolean),
      name: dir?.name || u.userId,
      email: dir?.email || '',
    }
  }),
)

function remove(row) {
  confirmDelete({
    header: t('project.accessUsers.removeConfirm.header'),
    message: t('project.accessUsers.removeConfirm.message', { name: row.name }),
    onConfirm: async () => {
      try {
        await store.removeProjectUser(props.project.projectID, row.userId, row.roleIds)
      } catch (err) {
        $toast.toastErrorHandler(t('project.accessUsers.toastRemoveFailed'))(err)
      }
    },
  })
}

async function refresh(id) {
  if (!id) return
  try {
    await Promise.all([usersStore.load(), store.loadProjectUsers(id)])
  } catch (err) {
    $toast.toastErrorHandler(t('project.accessUsers.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.projectID))
watch(() => props.project.projectID, refresh)
</script>
