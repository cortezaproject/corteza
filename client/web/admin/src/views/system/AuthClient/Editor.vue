<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <Form
    v-else-if="authClient"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4">
      <Card
        :pt="{
          body: { class: 'p-0 flex flex-col h-full min-h-0' },
          content: { class: 'p-0 flex flex-col h-full min-h-0' },
        }"
        class="overflow-hidden flex-1 min-h-0 flex flex-col"
      >
        <template #content>
          <Tabs v-model:value="activeTab" class="flex flex-col h-full min-h-0">
            <TabList class="rounded-t-lg shrink-0">
              <Tab value="basic">{{ $t('system.authclients.editor.tabs.basic') }}</Tab>
              <Tab value="security">{{ $t('system.authclients.editor.tabs.security') }}</Tab>
            </TabList>

            <TabPanels class="flex-1 overflow-y-auto min-h-0">
              <TabPanel value="basic">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField name="name" class="flex flex-col gap-2">
                    <label for="name" class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.name') }} *
                    </label>
                    <InputText id="name" name="name" v-model="authClient.meta.name" />
                    <Message
                      v-if="$form.name?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.name.error?.message }}
                    </Message>
                  </FormField>

                  <FormField name="handle" class="flex flex-col gap-2">
                    <label for="handle" class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.handle') }}
                    </label>
                    <InputText id="handle" name="handle" v-model="authClient.handle" />
                    <Message
                      v-if="$form.handle?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.handle.error?.message }}
                    </Message>
                  </FormField>

                  <div class="flex flex-col gap-2">
                    <label for="validGrant" class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.validGrant') }}
                    </label>
                    <Select
                      id="validGrant"
                      v-model="authClient.validGrant"
                      :options="grantOptions"
                      option-label="label"
                      option-value="value"
                    />
                  </div>

                  <div class="flex flex-col gap-2">
                    <label for="scope" class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.scope') }}
                    </label>
                    <InputText id="scope" v-model="authClient.scope" />
                  </div>

                  <div class="flex flex-col gap-2 md:col-span-2">
                    <label for="redirectURI" class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.redirectURI') }}
                    </label>
                    <InputText id="redirectURI" v-model="authClient.redirectURI" />
                  </div>

                  <div class="flex items-center gap-3">
                    <ToggleSwitch id="enabled" v-model="authClient.enabled" />
                    <label for="enabled" class="font-medium text-primary cursor-pointer">
                      {{ $t('system.authclients.editor.info.enabled.label') }}
                    </label>
                  </div>

                  <div class="flex items-center gap-3">
                    <ToggleSwitch id="trusted" v-model="authClient.trusted" />
                    <label for="trusted" class="font-medium text-primary cursor-pointer">
                      {{ $t('system.authclients.editor.info.trusted.label') }}
                    </label>
                  </div>
                </div>
              </TabPanel>

              <TabPanel value="security">
                <div class="flex flex-col gap-6">
                  <!-- Permitted Roles -->
                  <div class="flex flex-col gap-3">
                    <h3 class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.security.permittedRoles.label') }}
                    </h3>
                    <div class="flex items-center gap-2">
                      <CInputRole
                        class="flex-1"
                        :placeholder="$t('system.authclients.editor.info.add')"
                        clear-on-select
                        filter-context-roles
                        @select="role => addRoleToList('permittedRoles', role)"
                      />
                    </div>
                    <RoleList
                      :roles="permittedRoles"
                      @remove="role => removeRoleFromList('permittedRoles', role)"
                    />
                  </div>

                  <!-- Prohibited Roles -->
                  <div class="flex flex-col gap-3">
                    <h3 class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.security.prohibitedRoles.label') }}
                    </h3>
                    <div class="flex items-center gap-2">
                      <CInputRole
                        class="flex-1"
                        :placeholder="$t('system.authclients.editor.info.add')"
                        clear-on-select
                        filter-context-roles
                        @select="role => addRoleToList('prohibitedRoles', role)"
                      />
                    </div>
                    <RoleList
                      :roles="prohibitedRoles"
                      @remove="role => removeRoleFromList('prohibitedRoles', role)"
                    />
                  </div>

                  <!-- Forced Roles -->
                  <div class="flex flex-col gap-3">
                    <h3 class="font-medium text-primary">
                      {{ $t('system.authclients.editor.info.security.forcedRoles.label') }}
                    </h3>
                    <div class="flex items-center gap-2">
                      <CInputRole
                        class="flex-1"
                        :placeholder="$t('system.authclients.editor.info.add')"
                        clear-on-select
                        filter-context-roles
                        @select="role => addRoleToList('forcedRoles', role)"
                      />
                    </div>
                    <RoleList
                      :roles="forcedRoles"
                      @remove="role => removeRoleFromList('forcedRoles', role)"
                    />
                  </div>
                </div>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.authClients' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && authClient.canDeleteAuthClient"
            :label="$t('system.authclients.editor.info.delete')"
            :message="$t('general.confirm.delete')"
            :header="authClient.meta?.name || authClient.handle"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputDelete, CInputRole } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const authClient = ref(null)
const activeTab = ref('basic')

// Local role lists with full role objects for display
const permittedRoles = ref([])
const prohibitedRoles = ref([])
const forcedRoles = ref([])

const isEdit = computed(() => !!route.params.authClientID)

const pageTitle = computed(() =>
  isEdit.value
    ? t('system.authclients.editor.title.edit')
    : t('system.authclients.editor.title.create'),
)

const grantOptions = computed(() => [
  {
    label: t('system.authclients.editor.info.validGrant.authorization_code'),
    value: 'authorization_code',
  },
  {
    label: t('system.authclients.editor.info.validGrant.client_credentials'),
    value: 'client_credentials',
  },
  { label: t('system.authclients.editor.info.validGrant.password'), value: 'password' },
])

const initialValues = computed(() => ({
  name: authClient.value?.meta?.name || '',
  handle: authClient.value?.handle || '',
}))

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [
      { message: t('system.authclients.editor.info.handle.invalid-handle-characters') },
    ]
  }

  return { errors }
})

// Inline role list component
const RoleList = {
  props: {
    roles: { type: Array, required: true },
  },
  emits: ['remove'],
  template: `
    <div v-if="roles.length === 0" class="text-muted-color p-3 border rounded-lg bg-highlight text-center text-sm">
      {{ $t('system.authclients.editor.info.security.noRoles') }}
    </div>
    <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
      <div v-for="role in roles" :key="role.roleID" class="flex items-center justify-between p-2 px-3">
        <span class="font-medium">{{ role.name || role.handle || role.roleID }}</span>
        <Button icon="pi pi-trash" severity="danger" text rounded size="small" @click="$emit('remove', role)" />
      </div>
    </div>
  `,
}

function addRoleToList(listName, role) {
  if (!role) return
  const list = { permittedRoles, prohibitedRoles, forcedRoles }[listName]
  if (!list.value.find(r => r.roleID === role.roleID)) {
    list.value.push(role)
    authClient.value.security[listName] = list.value.map(r => r.roleID)
  }
}

function removeRoleFromList(listName, role) {
  const list = { permittedRoles, prohibitedRoles, forcedRoles }[listName]
  const idx = list.value.findIndex(r => r.roleID === role.roleID)
  if (idx !== -1) {
    list.value.splice(idx, 1)
    authClient.value.security[listName] = list.value.map(r => r.roleID)
  }
}

async function loadRolesForList(ids, targetRef) {
  if (!ids || ids.length === 0) return
  const roles = await Promise.all(
    ids.map(id => $SystemAPI.roleRead({ roleID: id }).catch(() => null)),
  )
  targetRef.value = roles.filter(Boolean).map(r => new system.Role(r))
}

async function loadAuthClient() {
  const authClientID = route.params.authClientID
  if (!authClientID) {
    authClient.value = new system.AuthClient({ enabled: true })
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.authClientRead({ authClientID })
    authClient.value = new system.AuthClient(raw)

    // Load role objects for security lists
    await Promise.all([
      loadRolesForList(authClient.value.security.permittedRoles, permittedRoles),
      loadRolesForList(authClient.value.security.prohibitedRoles, prohibitedRoles),
      loadRolesForList(authClient.value.security.forcedRoles, forcedRoles),
    ])
  } catch (e) {
    $toast.toastErrorHandler(t('notification.authclient.fetch.error'))(e)
    router.push({ name: 'system.authClients' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  saving.value = true
  try {
    const payload = {
      handle: authClient.value.handle,
      meta: authClient.value.meta,
      scope: authClient.value.scope,
      redirectURI: authClient.value.redirectURI,
      validGrant: authClient.value.validGrant,
      enabled: authClient.value.enabled,
      trusted: authClient.value.trusted,
      security: authClient.value.security,
    }

    if (isEdit.value) {
      payload.authClientID = authClient.value.authClientID
      const raw = await $SystemAPI.authClientUpdate(payload)
      authClient.value = new system.AuthClient(raw)
      $toast.toastSuccess(t('notification.authclient.update.success'))
    } else {
      const created = await $SystemAPI.authClientCreate(payload)
      $toast.toastSuccess(t('notification.authclient.create.success'))
      router.push({
        name: 'system.authClients.edit',
        params: { authClientID: created.authClientID },
      })
    }
  } catch (e) {
    $toast.toastErrorHandler(
      t(
        `notification.authclient.${isEdit.value ? 'update' : 'create'}.error`,
        'Failed to save auth client',
      ),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.authClientDelete({ authClientID: authClient.value.authClientID })
    $toast.toastSuccess(t('notification.authclient.delete.success'))
    router.push({ name: 'system.authClients' })
  } catch (e) {
    $toast.toastErrorHandler(t('notification.authclient.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

onMounted(() => loadAuthClient())
watch(
  () => route.params.authClientID,
  () => loadAuthClient(),
)
</script>
