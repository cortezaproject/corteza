<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="model.enabled"
      :label="$t('system.settings.editor.external.oidc.enabled')"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.handle') }}
        </label>
        <InputText v-model="model.handle" :disabled="!fresh" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.issuer') }}
        </label>
        <span class="text-xs text-muted-color">
          {{ $t('system.settings.editor.external.oidc.issuerHint') }}
        </span>
        <InputText v-model="model.issuer" placeholder="https://issuer.tld" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.clientKey') }}
        </label>
        <InputText v-model="model.key" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.clientSecret') }}
        </label>
        <InputText v-model="model.secret" class="w-full" />
      </div>

      <div class="flex flex-col gap-1 md:col-span-2">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.oidc.scope') }}
        </label>
        <span class="text-xs text-muted-color">
          {{ $t('system.settings.editor.external.oidc.scopeHint') }}
        </span>
        <InputText
          v-model="model.scope"
          :placeholder="$t('system.settings.editor.external.oidc.scopePlaceholder')"
          class="w-full"
        />
      </div>
    </div>

    <ExternalSecurity v-model="model.security" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ExternalSecurity from './ExternalSecurity.vue'

const model = defineModel({ type: Object, required: true })

const fresh = computed(() => {
  return Object.prototype.hasOwnProperty.call(model.value, 'fresh') && model.value.fresh
})
</script>
