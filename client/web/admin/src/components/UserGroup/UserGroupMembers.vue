<template>
  <div class="flex flex-col gap-4">
    <!-- User search -->
    <div class="flex items-center gap-2">
      <CInputUser class="flex-1" :placeholder="$t('system.user-groups.editor.members.placeholder')" clear-on-select @select="addMember" />
    </div>
    <!-- Members list -->
    <div v-if="loading" class="flex justify-center p-4">
      <ProgressSpinner />
    </div>
    <div v-else-if="members.length === 0" class="text-muted-color p-4 border rounded-lg bg-highlight text-center">
      {{ $t('system.user-groups.editor.members.empty') }}
    </div>
    <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
      <div v-for="m in members" :key="m.userID" class="flex items-center justify-between p-3">
        <div class="flex flex-col">
          <span class="font-medium">{{ m.name || m.email || m.handle || m.userID }}</span>
          <span v-if="(m.name || m.handle) && m.email" class="text-xs text-muted-color">{{ m.email }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
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

const loading = ref(false)
const members = ref([])

async function loadMembers() {
  loading.value = true
  try {
    const result = await $SystemAPI.userGroupMemberList({ userGroupID: props.userGroupID })
    const ids = result?.set || result || []
    if (ids.length > 0) {
      const users = await Promise.all(
        ids.map(id => $SystemAPI.userRead({ userID: id }).catch(() => null)),
      )
      members.value = users.filter(Boolean).map(u => new system.User(u))
    } else {
      members.value = []
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.fetchError'))(e)
  } finally {
    loading.value = false
  }
}

async function addMember(user) {
  if (!user) return
  try {
    await $SystemAPI.userGroupMemberAdd({ userGroupID: props.userGroupID, userID: user.userID })
    if (!members.value.find(m => m.userID === user.userID)) {
      members.value.push(new system.User(user))
    }
  } catch (e) {
    $toast.toastErrorHandler(t('system.user-groups.editor.members.addError'))(e)
  }
}

onMounted(() => loadMembers())
</script>
