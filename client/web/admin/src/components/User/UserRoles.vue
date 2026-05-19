<template>
  <div class="flex flex-col gap-4">
    <CInputRole
      class="w-full"
      :placeholder="$t('system.users.editor.roles.placeholder')"
      clear-on-select
      filter-context-roles
      :exclude-roles="Array.from(membershipIDs)"
      @select="onRoleSelect"
    />

    <CFormItemList
      :items="currentRoles"
      :loading="loading"
      :empty-message="$t('system.users.editor.roles.empty')"
      :remove-label="$t('system.users.editor.roles.remove')"
      item-key="roleID"
      @remove="removeRole"
    >
      <template #default="{ item }">
        <span class="font-medium">{{ item.name || item.handle || item.roleID }}</span>
        <span v-if="item.name && item.handle" class="text-xs text-muted-color block">
          {{ item.handle }}
        </span>
      </template>
    </CFormItemList>
  </div>
</template>

<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'

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
      const systemHandles = ['authenticated', 'anonymous', 'everyone']
      currentRoles.value = roles
        .filter(Boolean)
        .filter(r => !systemHandles.includes((r.handle || '').toLowerCase()))
        .filter(r => !(r.meta?.context?.expr || r.meta?.context?.resourceTypes?.length > 0))
        .map(r => new system.Role(r))
    }
  } catch (e) {
    console.error('Failed to load roles:', e)
    $toast.toastErrorHandler(t('notification.user.roles.error'))(e)
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
