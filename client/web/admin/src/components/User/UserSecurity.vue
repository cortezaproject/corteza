<template>
  <div class="flex flex-col gap-8">
    <!-- Password Section -->
    <div class="flex flex-col gap-4">
      <div class="flex flex-col gap-4">
        <FormField name="password" class="flex flex-col gap-2">
          <label for="password" class="font-medium text-primary">
            {{ $t('system.users.editor.password.new') }}
          </label>
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
        </FormField>

        <FormField name="confirmPassword" class="flex flex-col gap-2">
          <label for="confirmPassword" class="font-medium text-primary">
            {{ $t('system.users.editor.password.confirm') }}
          </label>
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
        </FormField>
      </div>
    </div>

    <!-- MFA Section -->
    <div class="flex flex-col gap-4">
      <div class="flex flex-col md:flex-row flex-wrap gap-6">
        <!-- Email OTP -->
        <div
          class="flex-1 min-w-[300px] flex items-start justify-between gap-4 p-4 border rounded-lg bg-emphasis"
        >
          <div class="flex flex-col gap-1">
            <span class="font-medium">{{ $t('system.users.editor.mfa.emailOTP.label') }}</span>
            <span class="text-sm text-surface-500" style="white-space: pre-line">
              {{ $t('system.users.editor.mfa.emailOTP.description') }}
            </span>
          </div>
          <div class="flex items-center">
            <ToggleSwitch
              :modelValue="user.meta.securityPolicy.mfa.enforcedEmailOTP"
              @update:modelValue="val => $emit('update:mfa', 'enforcedEmailOTP', val)"
            />
          </div>
        </div>

        <!-- TOTP -->
        <div
          class="flex-1 min-w-[300px] flex items-start justify-between gap-4 p-4 border rounded-lg bg-emphasis"
        >
          <div class="flex flex-col gap-1">
            <span class="font-medium">{{ $t('system.users.editor.mfa.TOTP.label') }}</span>
            <span class="text-sm text-surface-500">
              {{ $t('system.users.editor.mfa.TOTP.description') }}
            </span>
          </div>
          <div class="flex items-center">
            <ToggleSwitch
              :modelValue="user.meta.securityPolicy.mfa.enforcedTOTP"
              :disabled="!user.meta.securityPolicy.mfa.enforcedTOTP"
              @update:modelValue="val => $emit('update:mfa', 'enforcedTOTP', val)"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
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
