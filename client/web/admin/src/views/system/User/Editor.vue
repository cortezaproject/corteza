<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <Form
    v-else-if="user"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto">
      <div v-if="isEdit && user.canUpdateUser" class="flex justify-end gap-2">
        <Button
          v-if="!user.suspendedAt"
          :label="$t('system.users.editor.info.suspend')"
          icon="pi pi-pause"
          severity="secondary"
          size="small"
          outlined
          :disabled="suspending"
          @click="confirmSuspend"
        />
        <Button
          v-else
          :label="$t('system.users.editor.info.unsuspend')"
          icon="pi pi-play"
          severity="secondary"
          size="small"
          outlined
          :disabled="suspending"
          @click="confirmUnsuspend"
        />
      </div>

      <Panel
        :header="$t('system.users.editor.info.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="email" class="flex flex-col gap-2">
            <label for="email" class="font-medium text-primary">
              {{ $t('system.users.editor.info.email') }}
            </label>
            <InputText id="email" name="email" v-model="user.email" type="email" />
            <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">
              {{ $form.email.error?.message }}
            </Message>
          </FormField>

          <FormField name="name" class="flex flex-col gap-2">
            <label for="name" class="font-medium text-primary">
              {{ $t('system.users.editor.info.name') }}
            </label>
            <InputText id="name" name="name" v-model="user.name" />
            <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
              {{ $form.name.error?.message }}
            </Message>
          </FormField>

          <FormField name="handle" class="flex flex-col gap-2">
            <label for="handle" class="font-medium text-primary">
              {{ $t('system.users.editor.info.handle') }}
            </label>
            <InputText id="handle" name="handle" v-model="user.handle" />
            <Message v-if="$form.handle?.invalid" severity="error" size="small" variant="simple">
              {{ $form.handle.error?.message }}
            </Message>
          </FormField>

          <FormField name="userGroupID" class="flex flex-col gap-2">
            <label for="userGroupID" class="font-medium text-primary">
              {{ $t('system.users.editor.info.userGroup.label') }}
            </label>
            <CInputUserGroup id="userGroupID" v-model="user.userGroupID" class="w-full" />
          </FormField>
        </div>
      </Panel>

      <Panel
        v-if="isEdit"
        :header="$t('system.users.editor.security.title')"
        toggleable
        class="shadow"
      >
        <UserSecurity
          :user="user"
          v-model:passwords="passwords"
          @update:mfa="(key, val) => (user.meta.securityPolicy.mfa[key] = val)"
        />
      </Panel>

      <Panel
        v-if="isEdit"
        :header="$t('system.users.editor.roles.title')"
        toggleable
        class="shadow"
      >
        <UserRoles
          :user="user"
          v-model:membershipIDs="membershipIDs"
          :initialMembershipIDs="initialMembershipIDs"
        />
      </Panel>
    </div>

    <!-- Bottom Actions Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.users' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && user.canDeleteUser"
            :label="$t('system.users.editor.info.delete')"
            :message="$t('system.users.editor.info.deleteConfirm')"
            :header="$t('system.users.editor.info.delete')"
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
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
import { useConfirm } from 'primevue/useconfirm'

const { CInputDelete, CInputUserGroup } = components

import UserSecurity from '@/components/User/UserSecurity.vue'
import UserRoles from '@/components/User/UserRoles.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const confirm = useConfirm()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const suspending = ref(false)
const user = ref(null)

// Lifted State for Tabs
const passwords = ref({
  password: '',
  confirmPassword: '',
})
const initialMembershipIDs = ref(new Set())
const membershipIDs = ref(new Set())

// Computed
const isEdit = computed(() => !!route.params.userID)

const pageTitle = computed(() => {
  return isEdit.value ? t('system.users.editor.title.edit') : t('system.users.editor.title.create')
})

const initialValues = computed(() => {
  return {
    email: user.value?.email || '',
    name: user.value?.name || '',
    handle: user.value?.handle || '',
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.email || values.email.trim().length === 0) {
    errors.email = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('system.users.editor.info.invalid-characters') }]
  }

  return { errors }
})

// Methods
async function loadUser() {
  const userID = route.params.userID
  if (!userID) {
    // Create new
    user.value = new system.User({})
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.userRead({ userID })
    user.value = new system.User(raw)

    // Load initial roles if editing
    const memRes = await $SystemAPI.userMembershipList({ userID })
    const ids = (memRes.set || []).map(m => m.roleID)
    initialMembershipIDs.value = new Set(ids)
    membershipIDs.value = new Set(ids)
  } catch (e) {
    console.error('Failed to load user:', e)
    $toast.toastErrorHandler(t('notification.user.fetch.error'))(e)
    router.push({ name: 'system.users' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  if (isEdit.value && !user.value?.canUpdateUser) return

  saving.value = true
  try {
    const payload = {
      email: user.value.email,
      name: user.value.name,
      handle: user.value.handle,
      userGroupID: user.value.userGroupID || '0',
      meta: user.value.meta, // To save MFA toggles!
    }

    if (isEdit.value) {
      payload.userID = user.value.userID
      const raw = await $SystemAPI.userUpdate(payload)
      user.value = new system.User(raw)

      // Handle Password if provided
      if (passwords.value.password) {
        await $SystemAPI.userSetPassword({
          userID: payload.userID,
          password: passwords.value.password,
        })
        passwords.value = { password: '', confirmPassword: '' }
      }

      // Handle Role updates
      const rolesToAdd = [...membershipIDs.value].filter(id => !initialMembershipIDs.value.has(id))
      const rolesToRemove = [...initialMembershipIDs.value].filter(
        id => !membershipIDs.value.has(id),
      )

      for (const roleID of rolesToAdd) {
        await $SystemAPI.roleMemberAdd({ roleID, userID: payload.userID })
      }
      for (const roleID of rolesToRemove) {
        await $SystemAPI.roleMemberRemove({ roleID, userID: payload.userID })
      }
      initialMembershipIDs.value = new Set(membershipIDs.value)

      $toast.toastSuccess(t('notification.user.update.success'))
    } else {
      const created = await $SystemAPI.userCreate(payload)
      $toast.toastSuccess(t('notification.user.create.success'))
      router.push({
        name: 'system.users.edit',
        params: { userID: created.userID },
      })
    }
  } catch (e) {
    console.error('Failed to save user:', e)
    $toast.toastErrorHandler(
      isEdit.value ? t('notification.user.update.error') : t('notification.user.create.error'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.userDelete({ userID: user.value.userID })
    $toast.toastSuccess(t('notification.user.delete.success'))
    router.push({ name: 'system.users' })
  } catch (e) {
    console.error('Failed to delete user:', e)
    $toast.toastErrorHandler(t('notification.user.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

async function handleSuspend() {
  suspending.value = true
  try {
    await $SystemAPI.userSuspend({ userID: user.value.userID })
    const raw = await $SystemAPI.userRead({ userID: user.value.userID })
    user.value = new system.User(raw)
    $toast.toastSuccess(t('notification.user.suspend.success'))
  } catch (e) {
    console.error('Failed to suspend user:', e)
    $toast.toastErrorHandler(t('notification.user.suspend.error'))(e)
  } finally {
    suspending.value = false
  }
}

function confirmSuspend(event) {
  confirm.require({
    target: event.currentTarget,
    message: t('system.users.editor.info.suspendConfirm'),
    header: t('system.users.editor.info.suspend'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      outlined: true,
      size: 'small',
    },
    acceptProps: {
      label: t('system.users.editor.info.suspend'),
      severity: 'warn',
      size: 'small',
    },
    accept: () => {
      handleSuspend()
    },
  })
}

async function handleUnsuspend() {
  suspending.value = true
  try {
    await $SystemAPI.userUnsuspend({ userID: user.value.userID })
    const raw = await $SystemAPI.userRead({ userID: user.value.userID })
    user.value = new system.User(raw)
    $toast.toastSuccess(t('notification.user.unsuspend.success'))
  } catch (e) {
    console.error('Failed to unsuspend user:', e)
    $toast.toastErrorHandler(t('notification.user.unsuspend.error'))(e)
  } finally {
    suspending.value = false
  }
}

function confirmUnsuspend(event) {
  confirm.require({
    target: event.currentTarget,
    message: t('system.users.editor.info.unsuspendConfirm'),
    header: t('system.users.editor.info.unsuspend'),
    icon: 'pi pi-info-circle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      outlined: true,
      size: 'small',
    },
    acceptProps: {
      label: t('system.users.editor.info.unsuspend'),
      severity: 'success',
      size: 'small',
    },
    accept: () => {
      handleUnsuspend()
    },
  })
}

onMounted(() => {
  loadUser()
})
</script>
