<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="modelValue.enabled"
      :label="$t('system.settings.editor.external.oidc.enabled', 'Enabled')"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.handle', 'Handle') }}
        </label>
        <InputText v-model="modelValue.handle" :disabled="!fresh" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.issuer', 'OIDC Issuer URL') }}
        </label>
        <span class="text-xs text-surface-500">
          {{
            $t(
              'system.settings.editor.external.oidc.issuerHint',
              'Where to find the /.well-known/openid-configuration (without the /.well-known/openid-configuration part)',
            )
          }}
        </span>
        <InputText v-model="modelValue.issuer" placeholder="https://issuer.tld" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.clientKey', 'Client key') }}
        </label>
        <InputText v-model="modelValue.key" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.clientSecret', 'Secret') }}
        </label>
        <InputText v-model="modelValue.secret" class="w-full" />
      </div>

      <div class="flex flex-col gap-1 md:col-span-2">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.scope', 'Scope') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.oidc.scopeHint', 'Use space delimited string') }}
        </span>
        <InputText
          v-model="modelValue.scope"
          :placeholder="
            $t(
              'system.settings.editor.external.oidc.scopePlaceholder',
              'List out supported OAuth scope values',
            )
          "
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
