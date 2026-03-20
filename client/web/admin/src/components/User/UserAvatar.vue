<template>
  <div class="flex flex-col gap-4">
    <!-- Avatar Preview -->
    <div class="flex items-center gap-4">
      <div class="relative">
        <img
          v-if="user.meta?.avatarID && user.meta.avatarID !== '0'"
          :src="avatarURL"
          alt="Avatar"
          class="w-24 h-24 rounded-full object-cover border border-surface"
        />
        <div
          v-else
          class="w-24 h-24 rounded-full flex items-center justify-center text-2xl font-bold border border-surface"
          :style="{ backgroundColor: user.meta?.avatarBgColor || '#e5e7eb', color: user.meta?.avatarColor || '#374151' }"
        >
          {{ initials }}
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <span class="font-medium text-primary">{{ $t('system.users.editor.avatar.preview') }}</span>
        <span class="text-sm text-muted-color">
          {{ user.meta?.avatarID && user.meta.avatarID !== '0'
            ? $t('system.users.editor.avatar.hasAvatar')
            : $t('system.users.editor.avatar.noAvatar') }}
        </span>
      </div>
    </div>

    <!-- Upload -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-primary">{{ $t('system.users.editor.avatar.upload') }}</label>
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        class="block w-full text-sm text-muted-color file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-medium file:bg-primary file:text-primary-contrast hover:file:opacity-90 cursor-pointer"
        @change="handleFileChange"
      />
    </div>

    <!-- Color pickers -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-2">
        <label for="avatarColor" class="font-medium text-primary">
          {{ $t('system.users.editor.avatar.color') }}
        </label>
        <div class="flex items-center gap-2">
          <input
            id="avatarColor"
            type="color"
            :value="user.meta?.avatarColor || '#374151'"
            class="w-10 h-10 rounded cursor-pointer border border-surface p-0.5"
            @input="updateColor('avatarColor', $event.target.value)"
          />
          <span class="text-sm text-muted-color">{{ user.meta?.avatarColor || '#374151' }}</span>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <label for="avatarBgColor" class="font-medium text-primary">
          {{ $t('system.users.editor.avatar.bgColor') }}
        </label>
        <div class="flex items-center gap-2">
          <input
            id="avatarBgColor"
            type="color"
            :value="user.meta?.avatarBgColor || '#e5e7eb'"
            class="w-10 h-10 rounded cursor-pointer border border-surface p-0.5"
            @input="updateColor('avatarBgColor', $event.target.value)"
          />
          <span class="text-sm text-muted-color">{{ user.meta?.avatarBgColor || '#e5e7eb' }}</span>
        </div>
      </div>
    </div>

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
  const name = props.user?.name || props.user?.email || props.user?.handle || ''
  return name.charAt(0).toUpperCase() || '?'
})

function updateColor(field, value) {
  const updatedMeta = { ...props.user.meta, [field]: value }
  emit('update:user', { ...props.user, meta: updatedMeta })
}

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
