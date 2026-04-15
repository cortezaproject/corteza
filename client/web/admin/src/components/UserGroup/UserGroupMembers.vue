<template>
  <div class="flex flex-col gap-4">
    <!-- User search -->
    <div class="flex items-center gap-2">
      <CInputUser class="flex-1" :placeholder="$t('system.user-groups.editor.members.placeholder')" clear-on-select @select="addMember" />
    </div>
    <!-- Members list -->
    <div v-if="members.length === 0" class="text-muted-color p-4 border rounded-lg bg-highlight text-center">
      {{ $t('system.user-groups.editor.members.empty') }}
    </div>
    <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
      <div v-for="m in members" :key="m.userID" class="flex items-center justify-between p-3">
        <div class="flex flex-col">
          <span class="font-medium">{{ m.name || m.email || m.handle || m.userID }}</span>
          <span v-if="(m.name || m.handle) && m.email" class="text-xs text-muted-color">{{ m.email }}</span>
        </div>
        <Button icon="pi pi-trash" severity="danger" text rounded size="small" @click="removeMember(m)" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputUser } = components

const props = defineProps({
  userGroupID: {
    type: String,
    required: true,
  },
})

const { t } = useI18n()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const members = ref([])

async function loadMembers() {
  try {
    const result = await $SystemAPI.userGroupMemberList({ userGroupID: props.userGroupID })
    members.value = result?.set || result || []
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.fetchError'))(e)
  }
}

async function addMember(user) {
  if (!user) return
  try {
    await $SystemAPI.userGroupMemberAdd({ userGroupID: props.userGroupID, memberID: user.userID })
    if (!members.value.find(m => m.userID === user.userID)) {
      members.value.push(user)
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.addError'))(e)
  }
}

async function removeMember(member) {
  try {
    await $SystemAPI.userGroupMemberRemove({ userGroupID: props.userGroupID, memberID: member.userID })
    const idx = members.value.findIndex(m => m.userID === member.userID)
    if (idx !== -1) members.value.splice(idx, 1)
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.removeError'))(e)
  }
}

onMounted(() => loadMembers())
</script>
