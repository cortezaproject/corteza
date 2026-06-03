<template>
  <div class="flex flex-col gap-4">
    <CInputRole
      class="w-full"
      :placeholder="$t('system.user-groups.editor.roles.placeholder')"
      :exclude-roles="roles.map(r => r.roleID)"
      clear-on-select
      filter-context-roles
      @select="addRole"
    />

    <CFormItemList
      :items="roles"
      :empty-message="$t('system.user-groups.editor.roles.empty')"
      item-key="roleID"
      @remove="removeRole"
    >
      <template #default="{ item }">
        <span class="font-medium">{{ item.name || item.handle || item.roleID }}</span>
      </template>
    </CFormItemList>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

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
    const result = await $SystemAPI.roleList({ userGroupID: props.userGroupID })
    roles.value = result?.set || result || []
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.roles.fetchError'))(e)
  }
}

async function addRole(role) {
  if (!role) return
  try {
    await $SystemAPI.roleMemberAddGroup({ userGroupID: props.userGroupID, roleID: role.roleID })
    if (!roles.value.find(r => r.roleID === role.roleID)) {
      roles.value.push(role)
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.roles.addError'))(e)
  }
}

async function removeRole(role) {
  try {
    await $SystemAPI.roleMemberRemoveGroup({ userGroupID: props.userGroupID, roleID: role.roleID })
    const idx = roles.value.findIndex(r => r.roleID === role.roleID)
    if (idx !== -1) roles.value.splice(idx, 1)
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.roles.removeError'))(e)
  }
}

onMounted(() => loadRoles())
</script>
