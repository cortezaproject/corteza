<template>
  <div class="flex flex-col gap-4">
    <CInputSwitch
      v-model="model.enabled"
      :label="$t('system.settings.editor.external.saml.enabled')"
    />

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.name')"
      :description="$t('system.settings.editor.external.saml.desc.name')"
    >
      <InputText v-model="model.name" class="w-full" />
    </CFormGroup>

    <Divider />

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.cert.public')"
      :description="$t('system.settings.editor.external.saml.desc.cert.public')"
    >
      <Textarea v-model="model.cert" rows="4" class="w-full" />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.cert.private')"
      :description="$t('system.settings.editor.external.saml.desc.cert.private')"
    >
      <Textarea v-model="model.key" rows="4" class="w-full" />
    </CFormGroup>

    <Divider />

    <CInputSwitch
      v-model="model['sign-requests']"
      :label="$t('system.settings.editor.external.saml.requests.sign-requests')"
      :description="$t('system.settings.editor.external.saml.desc.requests.sign-requests')"
    />

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.requests.sign-method')"
      :description="$t('system.settings.editor.external.saml.desc.requests.sign-method')"
    >
      <Select
        v-model="model['sign-method']"
        :options="signMethods"
        option-label="text"
        option-value="value"
        :placeholder="$t('general.label.selectOption')"
        class="w-full"
      />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.requests.binding')"
      :description="$t('system.settings.editor.external.saml.desc.requests.binding')"
    >
      <Select
        v-model="model['binding']"
        :options="httpBindings"
        option-label="text"
        option-value="value"
        :placeholder="$t('general.label.selectOption')"
        class="w-full"
      />
    </CFormGroup>

    <Divider />

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.idp.url')"
      :description="$t('system.settings.editor.external.saml.desc.idp.url')"
    >
      <InputText v-model="model.idp.url" class="w-full" />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.idp.ident-name')"
      :description="$t('system.settings.editor.external.saml.desc.idp.ident-name')"
    >
      <InputText v-model="model.idp['ident-name']" class="w-full" />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.idp.ident-handle')"
      :description="$t('system.settings.editor.external.saml.desc.idp.ident-handle')"
    >
      <InputText v-model="model.idp['ident-handle']" class="w-full" />
    </CFormGroup>

    <CFormGroup
      :label="$t('system.settings.editor.external.saml.idp.ident-identifier')"
      :description="$t('system.settings.editor.external.saml.desc.idp.ident-identifier')"
    >
      <InputText v-model="model.idp['ident-identifier']" class="w-full" />
    </CFormGroup>

    <ExternalSecurity v-model="model.security" />
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import ExternalSecurity from './ExternalSecurity.vue'

const { t } = useI18n()

const model = defineModel({ type: Object, required: true })

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
