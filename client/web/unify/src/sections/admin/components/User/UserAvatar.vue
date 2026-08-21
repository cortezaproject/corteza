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
          class="text-color !p-0 !w-10 !h-10"
        >
          <template #default>
            <Avatar
              :image="avatarImage || undefined"
              :label="!avatarImage ? initials : undefined"
              :icon="!avatarImage && !initials ? 'pi pi-user' : undefined"
              shape="circle"
              class="!w-full !h-full !text-sm"
              @error="avatarBroken = true"
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
        accept="image/png,image/jpeg"
        :disabled="disabled || uploading"
        class="block w-full text-sm text-muted-color file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-medium file:bg-primary file:text-primary-contrast hover:file:opacity-90 cursor-pointer"
        @change="handleFileChange"
      />
    </CFormGroup>

    <!-- Removal is for an uploaded avatar only: every other user wears a
         generated initials avatar, which the server re-creates the moment it is
         deleted, so the button would promise something it cannot do. -->
    <div v-if="hasUploadedAvatar && !disabled">
      <Button
        :label="$t('system.users.editor.avatar.remove')"
        icon="pi pi-trash"
        severity="danger"
        outlined
        size="small"
        :loading="removing"
        @click="handleRemoveAvatar"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'

const props = defineProps({
  disabled: { type: Boolean, default: false },
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

// Both avatar kinds are attachments; the image is served from the attachment
// endpoint, not from the user's own /avatar path, which only accepts a POST.
const avatarURL = computed(() => {
  const attachmentID = props.user?.meta?.avatarID
  if (!attachmentID || attachmentID === '0') {
    return ''
  }
  return (
    $SystemAPI.baseURL +
    $SystemAPI.attachmentOriginalEndpoint({ attachmentID, kind: 'avatar', name: 'avatar' })
  )
})

// A user can carry an avatarID whose attachment is gone, and the broken image
// that draws says less about the account than the initials do.
const avatarBroken = ref(false)
watch(avatarURL, () => {
  avatarBroken.value = false
})

const avatarImage = computed(() => (avatarBroken.value ? '' : avatarURL.value))

const hasUploadedAvatar = computed(() => props.user?.meta?.avatarKind === 'avatar')

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

// Uploading and removing both leave the server holding an avatar the caller
// cannot predict — a delete regenerates initials under a fresh attachment ID —
// so the new state is read back rather than guessed. Only the avatar half of
// meta is taken: the rest of the draft is what the human has typed and not saved.
async function refreshAvatarMeta() {
  const raw = await $SystemAPI.userRead({ userID: props.user.userID })
  const { avatarID, avatarKind, avatarColor, avatarBgColor } = raw?.meta || {}

  emit(
    'update:user',
    new system.User({
      ...props.user,
      meta: { ...props.user.meta, avatarID, avatarKind, avatarColor, avatarBgColor },
    }),
  )
}

// The generated client posts its arguments as JSON, which turns a File into an
// empty object, so the upload is sent as multipart against the same endpoint.
// A rejected file comes back as a plain-text body on a 200.
async function uploadAvatar(file) {
  const formData = new FormData()
  formData.append('upload', file, file.name)

  const { data } = await $SystemAPI
    .api()
    .post($SystemAPI.userProfileAvatarEndpoint({ userID: props.user.userID }), formData, {
      headers: { 'Content-Type': undefined },
    })

  if (typeof data === 'string') {
    throw new Error(data.split('\n')[0].replace(/^Error:\s*/, ''))
  }

  if (data?.error) {
    throw new Error(data.error.message || data.error)
  }
}

async function handleFileChange(event) {
  const file = event.target.files?.[0]
  if (!file) return

  uploading.value = true
  try {
    await uploadAvatar(file)
    await refreshAvatarMeta()
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
    await $SystemAPI.userDeleteAvatar({ userID: props.user.userID })
    await refreshAvatarMeta()
    $toast.toastSuccess(t('notification.user.avatar.remove.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.user.avatar.remove.error'))(e)
  } finally {
    removing.value = false
  }
}
</script>
