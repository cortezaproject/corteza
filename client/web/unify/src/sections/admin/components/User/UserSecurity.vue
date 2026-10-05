<template>
  <div class="flex flex-col gap-4">
    <div class="grid grid-cols-1 gap-4">
      <CInputToggleCard
        :modelValue="user.meta.securityPolicy.mfa.enforcedEmailOTP"
        :label="$t('system.users.editor.mfa.emailOTP.label')"
        :description="$t('system.users.editor.mfa.emailOTP.description')"
        @update:modelValue="val => $emit('update:mfa', 'enforcedEmailOTP', val)"
        :disabled="disabled"
      />

      <CInputToggleCard
        :modelValue="user.meta.securityPolicy.mfa.enforcedTOTP"
        :disabled="disabled || !user.meta.securityPolicy.mfa.enforcedTOTP"
        :label="$t('system.users.editor.mfa.TOTP.label')"
        :description="$t('system.users.editor.mfa.TOTP.description')"
        dim-when-off
        @update:modelValue="val => $emit('update:mfa', 'enforcedTOTP', val)"
      />
    </div>

    <CManualScriptButtons
      resource-type="system:user"
      ui-page="user/editor"
      ui-slot="passwordFooter"
      container-class="flex flex-wrap justify-end gap-2"
      @click="$emit('script', $event)"
    />
  </div>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'

const { CInputToggleCard, CManualScriptButtons } = components

defineProps({
  disabled: { type: Boolean, default: false },
  user: {
    type: Object,
    required: true,
  },
})

defineEmits(['update:mfa', 'script'])
</script>
