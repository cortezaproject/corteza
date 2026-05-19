<template>
  <div class="flex flex-col gap-4">
    <CInputUser
      class="w-full"
      :placeholder="$t('system.roles.editor.members.placeholder')"
      clear-on-select
      :exclude-users="Array.from(memberIDs)"
      @select="onUserSelect"
    />

    <CFormItemList
      :items="currentMembers"
      :loading="loading"
      :empty-message="$t('system.roles.editor.members.empty')"
      :remove-label="$t('system.roles.editor.members.remove')"
      item-key="userID"
      @remove="removeMember"
    >
      <template #default="{ item }">
        <span class="font-medium">
          {{ item.name || item.email || item.handle || item.userID }}
        </span>
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
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components } from '@planetcrust/human-vue'

const { CInputUser } = components

const { t } = useI18n()

defineProps({
  role: {
    type: Object,
    required: true,
  },
  initialMemberIDs: {
    type: Object, // Set
    required: true,
  },
})

const memberIDs = defineModel('memberIDs', {
  type: Object, // Set
  required: true,
})

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(true)
const currentMembers = ref([])

async function loadData() {
  loading.value = true
  try {
    const ids = Array.from(memberIDs.value)
    if (ids.length > 0) {
      const users = await Promise.all(
        ids.map(id => $SystemAPI.userRead({ userID: id }).catch(() => null)),
      )
      currentMembers.value = users.filter(Boolean).map(u => new system.User(u))
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.role.members.error'))(e)
  } finally {
    loading.value = false
  }
}

function onUserSelect(user) {
  if (!user) return
  if (!memberIDs.value.has(user.userID)) {
    const newSet = new Set(memberIDs.value)
    newSet.add(user.userID)
    memberIDs.value = newSet
    currentMembers.value.push(new system.User(user))
  }
}

function removeMember(user) {
  const newSet = new Set(memberIDs.value)
  newSet.delete(user.userID)
  memberIDs.value = newSet
  currentMembers.value = currentMembers.value.filter(m => m.userID !== user.userID)
}

onMounted(() => loadData())
</script>
