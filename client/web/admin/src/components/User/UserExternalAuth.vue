<template>
  <CFormItemList
    :items="credentials"
    :loading="loading"
    :empty-message="$t('system.users.editor.externalAuth.empty')"
    :remove-label="$t('system.users.editor.externalAuth.remove')"
    :loading-key="removingID"
    item-key="credentialsID"
    @remove="handleRemove"
  >
    <template #default="{ item }">
      <span class="font-medium">{{ formatKind(item.label || item.kind) }}</span>
      <span v-if="item.label && item.kind" class="text-xs text-muted-color block">
        {{ formatKind(item.kind) }}
      </span>
    </template>
  </CFormItemList>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  userID: {
    type: String,
    required: true,
  },
})

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(true)
const credentials = ref([])
const removingID = ref(null)

function formatKind(value) {
  if (!value) return ''
  return value.charAt(0).toUpperCase() + value.slice(1)
}

async function loadCredentials() {
  loading.value = true
  try {
    const result = await $SystemAPI.userListCredentials({ userID: props.userID })
    credentials.value = result?.set || result || []
  } catch (e) {
    $toast.toastErrorHandler(t('notification.user.credentials.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleRemove(item) {
  removingID.value = item.credentialsID
  try {
    await $SystemAPI.userDeleteCredentials({
      userID: props.userID,
      credentialsID: item.credentialsID,
    })
    credentials.value = credentials.value.filter(c => c.credentialsID !== item.credentialsID)
    $toast.toastSuccess(t('notification.user.credentials.remove.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.user.credentials.remove.error'))(e)
  } finally {
    removingID.value = null
  }
}

onMounted(() => loadCredentials())

defineExpose({ loadCredentials })
</script>
