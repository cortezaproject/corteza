<template>
  <div class="flex flex-col gap-6">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center p-4">
      <ProgressSpinner />
    </div>

    <div v-else class="flex flex-col gap-4">
      <!-- Add Role -->
      <div class="flex items-center gap-2">
        <CInputRole
          class="flex-1"
          :placeholder="$t('system.users.editor.roles.placeholder', 'Select a role to add')"
          clear-on-select
          filter-context-roles
          @select="onRoleSelect"
        />
      </div>

      <!-- Current Roles -->
      <div
        v-if="currentRoles.length === 0"
        class="text-surface-500 p-4 border rounded-lg bg-surface-50 dark:bg-surface-900/50 text-center"
      >
        {{ $t('system.users.editor.roles.empty', 'No roles assigned to this user yet.') }}
      </div>

      <div v-else class="flex flex-col border rounded-lg divide-y bg-surface-0 dark:bg-surface-900">
        <div
          v-for="role in currentRoles"
          :key="role.roleID"
          class="flex items-center justify-between p-3"
        >
          <div class="flex flex-col">
            <span class="font-medium">{{ role.name || role.handle || role.roleID }}</span>
            <span v-if="role.name && role.handle" class="text-xs text-surface-500">
              {{ role.handle }}
            </span>
          </div>
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            :aria-label="$t('system.users.editor.roles.remove', 'Remove')"
            :title="$t('system.users.editor.roles.remove', 'Remove')"
            @click="removeRole(role)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputRole } = components

const { t } = useI18n()

defineProps({
  user: {
    type: Object,
    required: true,
  },
  initialMembershipIDs: {
    type: Object, // Set
    required: true,
  },
})

const membershipIDs = defineModel('membershipIDs', {
  type: Object, // Set
  required: true,
})

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(true)
const currentRoles = ref([])

// Methods
async function loadData() {
  loading.value = true
  try {
    const ids = Array.from(membershipIDs.value)
    if (ids.length > 0) {
      const roles = await Promise.all(
        ids.map(id => $SystemAPI.roleRead({ roleID: id }).catch(() => null)),
      )
      currentRoles.value = roles.filter(Boolean).map(r => new system.Role(r))
    }
  } catch (e) {
    console.error('Failed to load roles:', e)
    $toast.toastErrorHandler(t('notification.user.roles.error', 'Failed to load roles.'))(e)
  } finally {
    loading.value = false
  }
}

function onRoleSelect(role) {
  if (!role) return
  if (!membershipIDs.value.has(role.roleID)) {
    const newSet = new Set(membershipIDs.value)
    newSet.add(role.roleID)
    membershipIDs.value = newSet
    currentRoles.value.push(new system.Role(role))
  }
}

function removeRole(role) {
  const newSet = new Set(membershipIDs.value)
  newSet.delete(role.roleID)
  membershipIDs.value = newSet
  currentRoles.value = currentRoles.value.filter(r => r.roleID !== role.roleID)
}

onMounted(() => {
  loadData()
})
</script>
