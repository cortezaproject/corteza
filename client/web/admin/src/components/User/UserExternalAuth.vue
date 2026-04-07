<template>
  <div class="flex flex-col gap-4">
    <div v-if="loading" class="flex justify-center p-4">
      <ProgressSpinner />
    </div>

    <div
      v-else-if="credentials.length === 0"
      class="text-muted-color p-4 border rounded-lg bg-surface text-center"
    >
      {{ $t('system.users.editor.externalAuth.empty') }}
    </div>

    <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
      <div
        v-for="item in credentials"
        :key="item.credentialsID"
        class="flex items-center justify-between p-3"
      >
        <div class="flex flex-col gap-1">
          <span class="font-medium">{{ item.label || item.kind }}</span>
          <span v-if="item.label && item.kind" class="text-xs text-muted-color">
            {{ item.kind }}
          </span>
        </div>
        <Button
          icon="pi pi-trash"
          severity="danger"
          text
          size="small"
          :aria-label="$t('system.users.editor.externalAuth.remove')"
          :title="$t('system.users.editor.externalAuth.remove')"
          :loading="removingID === item.credentialsID"
          @click="handleRemove(item)"
        />
      </div>
    </div>
  </div>
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
</script>
