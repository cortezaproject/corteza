<template>
  <Panel :header="$t('system.roles.editor.clone.title')" toggleable class="shadow">
    <div class="flex flex-col gap-4">
      <p class="text-sm text-muted-color">{{ $t('system.roles.editor.clone.description') }}</p>
      <div class="flex items-center gap-2">
        <CInputRole
          class="flex-1"
          v-model="sourceRoleID"
          :placeholder="$t('system.roles.editor.clone.placeholder')"
          @select="onRoleSelect"
          clear-on-select
        />
        <Button
          :label="$t('system.roles.editor.clone.button')"
          icon="pi pi-copy"
          :disabled="!selectedRole || cloning"
          :loading="cloning"
          @click="handleClone"
        />
      </div>
      <p v-if="selectedRole" class="text-sm">
        {{ $t('system.roles.editor.clone.selected') }}: {{ selectedRole.name || selectedRole.handle }}
      </p>
    </div>
  </Panel>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputRole } = components

const props = defineProps({
  roleID: {
    type: String,
    required: true,
  },
})

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const sourceRoleID = ref('')
const selectedRole = ref(null)
const cloning = ref(false)

function onRoleSelect(role) {
  if (!role) return
  selectedRole.value = role
  sourceRoleID.value = role.roleID
}

async function handleClone() {
  if (!selectedRole.value) return

  cloning.value = true
  try {
    if (typeof $SystemAPI.roleClone === 'function') {
      await $SystemAPI.roleClone({
        roleID: sourceRoleID.value,
        cloneToRoleID: props.roleID,
      })
    } else if (typeof $SystemAPI.permissionsClone === 'function') {
      await $SystemAPI.permissionsClone({
        roleID: sourceRoleID.value,
        cloneToRoleID: props.roleID,
      })
    } else {
      throw new Error('not-available')
    }
    $toast.toastSuccess(t('notification.role.clone.success'))
    selectedRole.value = null
    sourceRoleID.value = ''
  } catch (e) {
    if (e?.message === 'not-available') {
      $toast.toastError(t('notification.role.clone.notAvailable'))
    } else {
      $toast.toastErrorHandler(t('notification.role.clone.error'))(e)
    }
  } finally {
    cloning.value = false
  }
}
</script>
