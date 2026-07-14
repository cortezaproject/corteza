<template>
  <div class="flex flex-col gap-4">
    <!-- Pick any platform user; selecting them adds their membership to this role
         at once. Mirrors the admin role editor's member picker. -->
    <CInputUser
      v-if="!disabled"
      class="w-full"
      clear-on-select
      size="small"
      :placeholder="$t('project.roleMembers.addPlaceholder')"
      :exclude-users="memberIds"
      @select="onSelect"
    />

    <CFormItemList
      :items="members"
      :empty-message="$t('project.roleMembers.empty')"
      :remove-label="$t('project.roleMembers.remove')"
      :hide-remove="disabled"
      :loading-key="busyId"
      item-key="userId"
      @remove="remove"
    >
      <template #default="{ item }">
        <CFormItemContent :title="item.name" :subtitle="item.email" />
      </template>
    </CFormItemList>
  </div>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { CInputUser } = components

const props = defineProps({
  project: { type: Object, required: true },
  role: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const usersStore = useProjectUsersStore()
const { t } = useI18n()
const $toast = inject('$toast')

const busyId = ref(null)

// Names/emails for members picked via the system-user search that aren't in the
// project directory (limited to 500) — so their card renders immediately.
const picked = reactive({})

// Members of THIS role, derived from the project-user list (grouped by user with
// the role IDs they hold), resolved against the directory or the pick cache.
const members = computed(() =>
  store
    .projectUsersFor(props.project.projectID)
    .filter(u => u.roleIds.includes(props.role.id))
    .map(u => {
      const dir = usersStore.findUser(u.userId) || picked[u.userId]
      return { userId: u.userId, name: dir?.name || u.userId, email: dir?.email || '' }
    }),
)

const memberIds = computed(() => members.value.map(m => m.userId))

async function onSelect(user) {
  if (!user) return
  const userId = String(user.userID)
  if (memberIds.value.includes(userId)) return
  picked[userId] = {
    name: user.name || user.label || user.email || userId,
    email: user.email || '',
  }
  try {
    await store.setProjectUserRole(props.project.projectID, userId, props.role.id, true)
  } catch (err) {
    $toast.toastErrorHandler(t('project.roleMembers.toastAddFailed'))(err)
  }
}

async function remove(item) {
  busyId.value = item.userId
  try {
    await store.setProjectUserRole(props.project.projectID, item.userId, props.role.id, false)
  } catch (err) {
    $toast.toastErrorHandler(t('project.roleMembers.toastRemoveFailed'))(err)
  } finally {
    busyId.value = null
  }
}
</script>
