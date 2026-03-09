<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="modelValue.enabled"
      :label="$t('system.settings.editor.external.saml.enabled', 'Enabled')"
    />

    <div class="flex flex-col gap-1">
      <label class="font-medium text-sm">
        {{ $t('system.settings.editor.external.saml.name', 'Name') }}
      </label>
      <span class="text-xs text-surface-500">
        {{ $t('system.settings.editor.external.saml.desc.name', 'IdP name used on the login screen (Login with <name>)') }}
      </span>
      <InputText v-model="modelValue.name" class="w-full" />
    </div>

    <Divider />

    <!-- Certificate -->
    <h5 class="text-sm font-semibold">
      {{ $t('system.settings.editor.external.saml.certificate', 'Certificate') }}
    </h5>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.cert.public', 'Public key') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.cert.public', 'Content will be minimized') }}
        </span>
        <Textarea v-model="modelValue.cert" rows="4" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.cert.private', 'Private key') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.cert.private', 'Content will be minimized') }}
        </span>
        <Textarea v-model="modelValue.key" rows="4" class="w-full" />
      </div>
    </div>

    <Divider />

    <!-- Requests -->
    <h5 class="text-sm font-semibold">
      {{ $t('system.settings.editor.external.saml.requests.title', 'Requests') }}
    </h5>

    <CInputSwitch
      v-model="modelValue['sign-requests']"
      :label="$t('system.settings.editor.external.saml.requests.sign-requests', 'Sign requests')"
      :description="$t('system.settings.editor.external.saml.desc.requests.sign-requests', 'Sign AuthNRequest and assertion')"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.requests.sign-method', 'Signature method') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.requests.sign-method', 'Method to use on signed requests') }}
        </span>
        <Select
          v-model="modelValue['sign-method']"
          :options="signMethods"
          option-label="text"
          option-value="value"
          :placeholder="$t('general.label.selectOption', 'Select an option')"
          class="w-full"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.requests.binding', 'Binding') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.requests.binding', 'The type of HTTP binding to use on AuthNRequest, defaults to HTTP Redirect (GET)') }}
        </span>
        <Select
          v-model="modelValue['binding']"
          :options="httpBindings"
          option-label="text"
          option-value="value"
          :placeholder="$t('general.label.selectOption', 'Select an option')"
          class="w-full"
        />
      </div>
    </div>

    <Divider />

    <!-- Identity Provider -->
    <h5 class="text-sm font-semibold">
      {{ $t('system.settings.editor.external.saml.idp.title', 'Identity provider') }}
    </h5>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.idp.url', 'URL') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.idp.url', 'Location of IdP metadata') }}
        </span>
        <InputText v-model="modelValue.idp.url" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.idp.ident-name', 'Name field') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.idp.ident-name', 'Name of the IdP field used for filling Corteza user full name') }}
        </span>
        <InputText v-model="modelValue.idp['ident-name']" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.idp.ident-handle', 'Handle field') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.idp.ident-handle', 'Name of the IdP field used for filling Corteza user handle or nickname') }}
        </span>
        <InputText v-model="modelValue.idp['ident-handle']" class="w-full" />
      </div>

      <div class="flex flex-col gap-1">
        <label class="font-medium text-sm">
          {{ $t('system.settings.editor.external.saml.idp.ident-identifier', 'Identifier field') }}
        </label>
        <span class="text-xs text-surface-500">
          {{ $t('system.settings.editor.external.saml.desc.idp.ident-identifier', 'Name of the IdP field used for filling and matching Corteza user email') }}
        </span>
        <InputText v-model="modelValue.idp['ident-identifier']" class="w-full" />
      </div>
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
    text: t('system.settings.editor.external.saml.requests.binding-post', 'HTTP POST'),
  },
  {
    value: 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect',
    text: t('system.settings.editor.external.saml.requests.binding-redirect', 'HTTP Redirect'),
  },
]
</script>
