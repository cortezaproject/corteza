<template>
  <div class="flex flex-col gap-6">
    <div v-if="loading" class="flex justify-center p-4">
      <ProgressSpinner />
    </div>

    <div v-else class="flex flex-col gap-4">
      <!-- Add User -->
      <div class="flex items-center gap-2">
        <CInputUser
          class="flex-1"
          :placeholder="$t('system.roles.editor.members.placeholder')"
          clear-on-select
          @select="onUserSelect"
        />
      </div>

      <!-- Current Members -->
      <div
        v-if="currentMembers.length === 0"
        class="text-muted-color p-4 border rounded-lg bg-highlight text-center"
      >
        {{ $t('system.roles.editor.members.empty') }}
      </div>

      <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
        <div
          v-for="member in currentMembers"
          :key="member.userID"
          class="flex items-center justify-between p-3"
        >
          <div class="flex flex-col">
            <span class="font-medium">
              {{ member.name || member.email || member.handle || member.userID }}
            </span>
            <span
              v-if="(member.name || member.handle) && member.email"
              class="text-xs text-muted-color"
            >
              {{ member.email }}
            </span>
          </div>
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            :aria-label="$t('system.roles.editor.members.remove')"
            :title="$t('system.roles.editor.members.remove')"
            @click="removeMember(member)"
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
