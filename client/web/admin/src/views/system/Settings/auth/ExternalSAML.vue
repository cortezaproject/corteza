<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="modelValue.enabled"
      :label="$t('system.settings.editor.external.saml.enabled')"
    />

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.name') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.name') }}
      </span>
      <InputText v-model="modelValue.name" class="w-full" />
    </div>

    <Divider />

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.cert.public') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.cert.public') }}
      </span>
      <Textarea v-model="modelValue.cert" rows="4" class="w-full" />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.cert.private') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.cert.private') }}
      </span>
      <Textarea v-model="modelValue.key" rows="4" class="w-full" />
    </div>

    <Divider />

    <CInputSwitch
      v-model="modelValue['sign-requests']"
      :label="$t('system.settings.editor.external.saml.requests.sign-requests')"
      :description="$t('system.settings.editor.external.saml.desc.requests.sign-requests')"
    />

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.requests.sign-method') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.requests.sign-method') }}
      </span>
      <Select
        v-model="modelValue['sign-method']"
        :options="signMethods"
        option-label="text"
        option-value="value"
        :placeholder="$t('general.label.selectOption')"
        class="w-full"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.requests.binding') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.requests.binding') }}
      </span>
      <Select
        v-model="modelValue['binding']"
        :options="httpBindings"
        option-label="text"
        option-value="value"
        :placeholder="$t('general.label.selectOption')"
        class="w-full"
      />
    </div>

    <Divider />

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.idp.url') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.idp.url') }}
      </span>
      <InputText v-model="modelValue.idp.url" class="w-full" />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.idp.ident-name') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.idp.ident-name') }}
      </span>
      <InputText v-model="modelValue.idp['ident-name']" class="w-full" />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.idp.ident-handle') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.idp.ident-handle') }}
      </span>
      <InputText v-model="modelValue.idp['ident-handle']" class="w-full" />
    </div>

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.idp.ident-identifier') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.idp.ident-identifier') }}
      </span>
      <InputText v-model="modelValue.idp['ident-identifier']" class="w-full" />
    </div>

    <ExternalSecurity v-model="modelValue.security" />
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import ExternalSecurity from './ExternalSecurity.vue'

const { t } = useI18n()

defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
})

const signMethods = [
  { value: 'http://www.w3.org/2000/09/xmldsig#rsa-sha1', text: 'SHA1' },
  { value: 'http://www.w3.org/2001/04/xmldsig-more#rsa-sha256', text: 'SHA256' },
  { value: 'http://www.w3.org/2001/04/xmldsig-more#rsa-sha512', text: 'SHA512' },
]

const httpBindings = [
  {
    value: 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST',
    text: t('system.settings.editor.external.saml.requests.binding-post'),
  },
  {
    value: 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect',
    text: t('system.settings.editor.external.saml.requests.binding-redirect'),
  },
]
</script>
