<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ $t('system.email.editor.title') }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else class="flex flex-col h-full">
    <CViewContainer scroll gap="5">
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('system.email.editor.server.testSmtpConfigs.button')"
          icon="pi pi-bolt"
          severity="secondary"
          size="small"
          outlined
          :loading="testing"
          @click="handleTestSmtp"
        />
      </div>

      <Panel
        :header="$t('system.email.editor.server.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="flex flex-col gap-4">
          <!-- Host : Port -->
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.email.editor.server.host.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('system.email.editor.server.host.description') }}
            </span>
            <InputGroup>
              <InputText v-model="server.host" placeholder="host.domain.tld" class="flex-1" />
              <InputGroupAddon>:</InputGroupAddon>
              <InputNumber v-model="server.port" :use-grouping="false" class="w-24" />
            </InputGroup>
          </div>

          <Divider />

          <!-- User / Password -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-sm">
                {{ $t('system.email.editor.server.user.label') }}
              </label>
              <span class="text-xs text-muted-color">
                {{ $t('system.email.editor.server.user.description') }}
              </span>
              <InputText v-model="server.user" autocomplete="off" class="w-full" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-sm">
                {{ $t('system.email.editor.server.password.label') }}
              </label>
              <span class="text-xs text-muted-color">
                {{ $t('system.email.editor.server.password.description') }}
              </span>
              <InputText v-model="server.pass" type="password" autocomplete="off" class="w-full" />
            </div>
          </div>

          <Divider />

          <!-- From address -->
          <div class="flex flex-col gap-1">
            <label class="font-medium text-sm">
              {{ $t('system.email.editor.server.from.label') }}
            </label>
            <span class="text-xs text-muted-color">
              {{ $t('system.email.editor.server.from.description') }}
            </span>
            <InputText v-model="server.from" type="email" class="w-full" />
          </div>

          <Divider />

          <!-- TLS -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-sm">
                {{ $t('system.email.editor.server.tlsServerName.label') }}
              </label>
              <span class="text-xs text-muted-color">
                {{ $t('system.email.editor.server.tlsServerName.description') }}
              </span>
              <InputText v-model="server.tlsServerName" class="w-full" />
            </div>
            <CInputSwitch
              v-model="server.tlsInsecure"
              :label="$t('system.email.editor.server.tlsInsecure.label')"
              :description="$t('system.email.editor.server.tlsInsecure.description')"
            />
          </div>
        </div>
      </Panel>
    </CViewContainer>

    <CEditorActions>
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="handleSave"
      />
    </CEditorActions>
  </div>
</template>

<script setup>
import { inject, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { CViewContainer } = components
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const testing = ref(false)

const server = reactive({
  host: '',
  port: 25,
  user: '',
  pass: '',
  from: '',
  tlsInsecure: false,
  tlsServerName: '',
})

async function loadSettings() {
  loading.value = true
  try {
    const result = await $SystemAPI.settingsList()
    const smtpSetting = (result || []).find(s => s.name === 'smtp.servers')

    if (
      smtpSetting &&
      smtpSetting.value &&
      smtpSetting.value.length > 0 &&
      typeof smtpSetting.value[0] === 'object'
    ) {
      Object.assign(server, smtpSetting.value[0])
    }
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.fetch.error'))(e)
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await $SystemAPI.settingsUpdate({
      values: [{ name: 'smtp.servers', value: [{ ...server }] }],
    })
    $toast.toastSuccess(t('notification.settings.update.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.update.error'))(e)
  } finally {
    saving.value = false
  }
}

async function handleTestSmtp() {
  testing.value = true
  try {
    const response = await $SystemAPI.smtpConfigurationCheckerCheck({
      host: server.host,
      port: parseInt(server.port),
      recipients: [server.from],
      username: server.user,
      password: server.pass,
      tlsInsecure: server.tlsInsecure,
      tlsServerName: server.tlsServerName,
    })

    if (Object.values(response).every(resp => resp === '')) {
      $toast.toastSuccess(t('notification.settings.system.smtpCheck.success'))
    }

    Object.keys(response).forEach(key => {
      if (response[key]) {
        $toast.toastWarning(`${key}: ${response[key]}`)
      }
    })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.settings.system.smtpCheck.error'))(e)
  } finally {
    testing.value = false
  }
}

onMounted(() => loadSettings())
</script>
