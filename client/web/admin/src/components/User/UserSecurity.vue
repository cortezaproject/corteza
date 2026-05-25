<template>
  <div class="flex flex-col gap-8">
    <!-- Password Section -->
    <div class="flex flex-col gap-4">
      <div class="flex flex-col gap-4">
        <CFormGroup name="password" :label="$t('system.users.editor.password.new')">
          <Password
            id="password"
            name="password"
            v-model="passwords.password"
            toggleMask
            :feedback="false"
            :inputProps="{ autocomplete: 'new-password' }"
            inputClass="w-full"
            class="w-full relative"
          />
        </CFormGroup>

        <CFormGroup name="confirmPassword" :label="$t('system.users.editor.password.confirm')">
          <Password
            id="confirmPassword"
            name="confirmPassword"
            v-model="passwords.confirmPassword"
            toggleMask
            :feedback="false"
            :inputProps="{ autocomplete: 'new-password' }"
            inputClass="w-full"
            class="w-full"
          />
        </CFormGroup>
      </div>
    </div>

    <!-- MFA Section -->
    <div class="grid grid-cols-1 gap-4">
      <CInputToggleCard
        :modelValue="user.meta.securityPolicy.mfa.enforcedEmailOTP"
        :label="$t('system.users.editor.mfa.emailOTP.label')"
        :description="$t('system.users.editor.mfa.emailOTP.description')"
        @update:modelValue="val => $emit('update:mfa', 'enforcedEmailOTP', val)"
      />

      <CInputToggleCard
        :modelValue="user.meta.securityPolicy.mfa.enforcedTOTP"
        :disabled="!user.meta.securityPolicy.mfa.enforcedTOTP"
        :label="$t('system.users.editor.mfa.TOTP.label')"
        :description="$t('system.users.editor.mfa.TOTP.description')"
        dim-when-off
        @update:modelValue="val => $emit('update:mfa', 'enforcedTOTP', val)"
      />
    </div>
  </div>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'

const { CInputToggleCard } = components

defineProps({
  user: {
    type: Object,
    required: true,
  },
})

defineEmits(['update:mfa'])

const passwords = defineModel('passwords', {
  type: Object,
  required: true,
})
</script>
