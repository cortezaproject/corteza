<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="modelValue.enabled"
      :label="$t('system.settings.editor.external.oidc.enabled')"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.handle') }}
        </label>
        <InputText v-model="modelValue.handle" :disabled="!fresh" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.issuer') }}
        </label>
        <span class="text-xs text-muted-color">
          {{ $t('system.settings.editor.external.oidc.issuerHint') }}
        </span>
        <InputText v-model="modelValue.issuer" placeholder="https://issuer.tld" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.clientKey') }}
        </label>
        <InputText v-model="modelValue.key" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.clientSecret') }}
        </label>
        <InputText v-model="modelValue.secret" class="w-full" />
      </div>

      <div class="flex flex-col gap-1 md:col-span-2">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.scope') }}
        </label>
        <span class="text-xs text-muted-color">
          {{ $t('system.settings.editor.external.oidc.scopeHint') }}
        </span>
        <InputText
          v-model="modelValue.scope"
          :placeholder="$t('system.settings.editor.external.oidc.scopePlaceholder')"
          class="w-full"
        />
      </div>
    </div>

    <ExternalSecurity v-model="modelValue.security" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ExternalSecurity from './ExternalSecurity.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
})

const fresh = computed(() => {
  return Object.prototype.hasOwnProperty.call(props.modelValue, 'fresh') && props.modelValue.fresh
})
</script>
