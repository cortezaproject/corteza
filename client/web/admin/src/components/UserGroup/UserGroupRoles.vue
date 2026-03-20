<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <CInputRole class="flex-1" :placeholder="$t('system.user-groups.editor.roles.placeholder')" clear-on-select @select="addRole" />
    </div>
    <!-- Roles list -->
    <div v-if="roles.length === 0" class="text-muted-color p-4 border rounded-lg bg-highlight text-center">
      {{ $t('system.user-groups.editor.roles.empty') }}
    </div>
    <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
      <div v-for="r in roles" :key="r.roleID" class="flex items-center justify-between p-3">
        <span class="font-medium">{{ r.name || r.handle || r.roleID }}</span>
        <Button icon="pi pi-trash" severity="danger" text rounded size="small" @click="removeRole(r)" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputRole } = components

const props = defineProps({
  userGroupID: {
    type: String,
    required: true,
  },
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const roles = ref([])

async function loadRoles() {
  try {
    const result = await $SystemAPI.userGroupRoleList({ userGroupID: props.userGroupID })
    roles.value = result?.set || result || []
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.roles.fetchError', 'Failed to load roles'))(e)
  }
}

async function addRole(role) {
  if (!role) return
  try {
    await $SystemAPI.userGroupRoleAdd({ userGroupID: props.userGroupID, roleID: role.roleID })
    if (!roles.value.find(r => r.roleID === role.roleID)) {
      roles.value.push(role)
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.roles.addError', 'Failed to add role'))(e)
  }
}

async function removeRole(role) {
  try {
    await $SystemAPI.userGroupRoleRemove({ userGroupID: props.userGroupID, roleID: role.roleID })
    const idx = roles.value.findIndex(r => r.roleID === role.roleID)
    if (idx !== -1) roles.value.splice(idx, 1)
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.roles.removeError', 'Failed to remove role'))(e)
  }
}

onMounted(() => loadRoles())
</script>
