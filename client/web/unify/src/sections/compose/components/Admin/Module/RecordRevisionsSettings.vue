<template>
  <div v-if="module" class="flex flex-col gap-6">
    <div class="flex items-center gap-3">
      <ToggleSwitch
        v-model="module.config.recordRevisions.enabled"
        inputId="recordRevisionsEnabled"
      />
      <label
        for="recordRevisionsEnabled"
        class="font-medium text-primary cursor-pointer select-none"
      >
        {{ $t('module.edit.config.record-revisions.enabled') }}
      </label>
    </div>

    <CFormGroup
      name="recordRevisions.ident"
      :label="$t('module.edit.config.record-revisions.ident.label')"
      :description="$t('module.edit.config.record-revisions.ident.description')"
      input-id="ident"
      class="max-w-lg"
    >
      <InputText
        id="ident"
        v-model="module.config.recordRevisions.ident"
        :disabled="!module.config.recordRevisions.enabled"
        :placeholder="$t('module.edit.config.record-revisions.ident.placeholder')"
        class="w-full"
      />
    </CFormGroup>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const module = inject('moduleDraft')

// Ensure config.recordRevisions exists safely
if (!module.value.config) module.value.config = {}
if (!module.value.config.recordRevisions) {
  module.value.config.recordRevisions = { enabled: false, ident: '' }
}
</script>
