<template>
  <div class="flex flex-col gap-4">
    <!-- Avatar Preview -->
    <CFormGroup :label="$t('system.users.editor.avatar.preview')">
      <div class="relative pointer-events-none">
        <Button
          rounded
          variant="outlined"
          severity="secondary"
          size="large"
          class="text-color !p-0 !w-[2.5rem] !h-[2.5rem]"
        >
          <template #default>
            <Avatar
              :image="avatarURL || undefined"
              :label="!avatarURL ? initials : undefined"
              :icon="!avatarURL && !initials ? 'pi pi-user' : undefined"
              shape="circle"
              class="!w-full !h-full !text-sm"
            />
          </template>
        </Button>
      </div>
    </CFormGroup>

    <!-- Upload -->
    <CFormGroup :label="$t('system.users.editor.avatar.upload')">
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        class="block w-full text-sm text-muted-color file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-medium file:bg-primary file:text-primary-contrast hover:file:opacity-90 cursor-pointer"
        @change="handleFileChange"
      />
    </CFormGroup>

    <!-- Remove avatar button -->
    <div v-if="user.meta?.avatarID && user.meta.avatarID !== '0'">
      <Button
        :label="$t('system.users.editor.avatar.remove')"
        icon="pi pi-trash"
        severity="danger"
        outlined
        :loading="removing"
        @click="handleRemoveAvatar"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:user'])

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const fileInput = ref(null)
const removing = ref(false)
const uploading = ref(false)

const avatarURL = computed(() => {
  if (!props.user?.userID || !props.user?.meta?.avatarID || props.user.meta.avatarID === '0') {
    return ''
  }
  return $SystemAPI.baseURL + '/system/users/' + props.user.userID + '/avatar'
})

const initials = computed(() => {
  const name = props.user?.name
  if (name) {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }
  const email = props.user?.email
  if (email) return email[0].toUpperCase()
  const handle = props.user?.handle
  if (handle) return handle[0].toUpperCase()
  return ''
})

async function handleFileChange(event) {
  const file = event.target.files?.[0]
  if (!file) return

  uploading.value = true
  try {
    const raw = await $SystemAPI.userProfileAvatarUpload({
      userID: props.user.userID,
      upload: file,
    })
    emit('update:user', { ...props.user, meta: { ...props.user.meta, ...raw?.meta } })
    $toast.toastSuccess(t('notification.user.avatar.upload.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.user.avatar.upload.error'))(e)
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

async function handleRemoveAvatar() {
  removing.value = true
  try {
    await $SystemAPI.userProfileAvatarDelete({ userID: props.user.userID })
    const updatedMeta = { ...props.user.meta, avatarID: '0' }
    emit('update:user', { ...props.user, meta: updatedMeta })
    $toast.toastSuccess(t('notification.user.avatar.remove.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.user.avatar.remove.error'))(e)
  } finally {
    removing.value = false
  }
}
</script>
