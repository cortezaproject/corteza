<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="model.enabled"
      :label="$t('system.settings.editor.external.standard.enabled')"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <CFormGroup :label="$t('system.settings.editor.external.standard.clientKey')">
        <InputText v-model="model.key" :required="model.enabled" class="w-full" />
      </CFormGroup>

      <CFormGroup :label="$t('system.settings.editor.external.standard.clientSecret')">
        <InputText v-model="model.secret" :required="model.enabled" class="w-full" />
      </CFormGroup>
    </div>

    <ExternalSecurity v-model="model.security" />
  </div>
</template>

<script setup>
import ExternalSecurity from './ExternalSecurity.vue'

const model = defineModel({ type: Object, required: true })
</script>
