<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="model.enabled"
      :label="$t('system.settings.editor.external.oidc.enabled')"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <CFormGroup :label="$t('system.settings.editor.external.oidc.handle')">
        <InputText v-model="model.handle" :disabled="!fresh" class="w-full" />
      </CFormGroup>

      <CFormGroup
        :label="$t('system.settings.editor.external.oidc.issuer')"
        :description="$t('system.settings.editor.external.oidc.issuerHint')"
      >
        <InputText v-model="model.issuer" placeholder="https://issuer.tld" class="w-full" />
      </CFormGroup>

      <CFormGroup :label="$t('system.settings.editor.external.oidc.clientKey')">
        <InputText v-model="model.key" class="w-full" />
      </CFormGroup>

      <CFormGroup :label="$t('system.settings.editor.external.oidc.clientSecret')">
        <InputText v-model="model.secret" class="w-full" />
      </CFormGroup>

      <CFormGroup
        :label="$t('system.settings.editor.external.oidc.scope')"
        :description="$t('system.settings.editor.external.oidc.scopeHint')"
        class="md:col-span-2"
      >
        <InputText
          v-model="model.scope"
          :placeholder="$t('system.settings.editor.external.oidc.scopePlaceholder')"
          class="w-full"
        />
      </CFormGroup>
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
