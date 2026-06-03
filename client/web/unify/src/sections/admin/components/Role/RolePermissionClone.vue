<template>
  <Dialog
    :visible="visible"
    @update:visible="$emit('update:visible', $event)"
    :header="$t('system.roles.editor.clone.title')"
    modal
    class="w-[30rem]"
  >
    <div class="flex flex-col gap-4 py-2">
      <p class="text-sm text-muted-color">{{ $t('system.roles.editor.clone.description') }}</p>
      <CInputRole
        class="w-full"
        v-model="sourceRoleID"
        :placeholder="$t('system.roles.editor.clone.placeholder')"
        @select="onRoleSelect"
      />
    </div>
    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          outlined
          size="small"
          :disabled="cloning"
          @click="$emit('update:visible', false)"
        />
        <Button
          :label="$t('system.roles.editor.clone.button')"
          size="small"
          :disabled="!selectedRole || cloning"
          :loading="cloning"
          @click="handleClone"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CInputRole } = components

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  roleID: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['update:visible', 'cloned'])

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const sourceRoleID = ref('')
const selectedRole = ref(null)
const cloning = ref(false)

function onRoleSelect(role) {
  if (!role) {
    selectedRole.value = null
    sourceRoleID.value = ''
    return
  }
  selectedRole.value = role
  sourceRoleID.value = role.roleID
}

async function handleClone() {
  if (!selectedRole.value) return

  cloning.value = true
  try {
    if (typeof $SystemAPI.roleCloneRules === 'function') {
      await $SystemAPI.roleCloneRules({
        roleID: props.roleID,
        cloneToRoleID: [sourceRoleID.value],
      })
    } else {
      throw new Error('not-available')
    }
    $toast.toastSuccess(t('notification.role.clone.success'))
    selectedRole.value = null
    sourceRoleID.value = ''
    emit('cloned')
    emit('update:visible', false)
  } catch (e) {
    if (e?.message === 'not-available') {
      $toast.toastWarning(t('notification.role.clone.notAvailable'))
    } else {
      $toast.toastErrorHandler(t('notification.role.clone.error'))(e)
    }
  } finally {
    cloning.value = false
  }
}
</script>
