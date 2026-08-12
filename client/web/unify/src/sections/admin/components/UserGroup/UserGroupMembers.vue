<template>
  <div class="flex flex-col gap-4">
    <CInputUser
      class="w-full"
      :placeholder="$t('system.user-groups.editor.members.placeholder')"
      clear-on-select
      @select="addMember"
    />

    <CFormItemList
      :items="members"
      :loading="loading"
      :empty-message="$t('system.user-groups.editor.members.empty')"
      item-key="userID"
      hide-remove
    >
      <template #default="{ item }">
        <span class="font-medium">{{ item.name || item.email || item.handle || item.userID }}</span>
        <span
          v-if="(item.name || item.handle) && item.email"
          class="text-xs text-muted-color block"
        >
          {{ item.email }}
        </span>
      </template>
    </CFormItemList>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'

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
