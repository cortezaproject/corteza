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
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-5 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2">
        <template v-if="user.canUpdateUser">
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
          <Button
            :label="$t('system.users.editor.info.revokeAllSession')"
            icon="pi pi-sign-out"
            severity="secondary"
            size="small"
            outlined
            :disabled="revoking || isSelf"
            @click="confirmRevokeSessions"
          />
        </template>
        
        <CPermissionsButton
          v-if="user.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:user/${user.userID}`"
          :title="user.name || user.handle || user.email || user.userID"
          :target="user.name || user.handle || user.email || user.userID"
        />
      </div>

      <Panel
        :header="$t('system.users.editor.info.title')"
        toggleable
        :collapsed="false"
        class="shadow"
      >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="email" :label="$t('system.users.editor.info.email')" required>
            <InputText id="email" name="email" v-model="user.email" type="email" />
          </CFormGroup>

          <CFormGroup name="name" :label="$t('system.users.editor.info.name')">
            <InputText id="name" name="name" v-model="user.name" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('system.users.editor.info.handle')">
            <InputText id="handle" name="handle" v-model="user.handle" />
          </CFormGroup>

          <CFormGroup name="userGroupID" :label="$t('system.users.editor.info.userGroup.label')">
            <CInputUserGroup id="userGroupID" v-model="user.userGroupID" :clearable="false" class="w-full" />
          </CFormGroup>
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

      <Panel v-if="isEdit" :header="$t('system.users.editor.avatar.title')" toggleable class="shadow">
        <UserAvatar :user="user" @update:user="(u) => (user = u)" />
      </Panel>

      <Panel v-if="isEdit" :header="$t('system.users.editor.externalAuth.title')" toggleable class="shadow">
        <UserExternalAuth ref="externalAuthRef" :userID="user.userID" />
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'system.users' }">
      <CInputDelete
        v-if="isEdit && user.canDeleteUser"
        :label="$t('system.users.editor.info.delete')"
        :message="$t('system.users.editor.info.deleteConfirm')"
        :header="user.name || user.handle || user.email || user.userID"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
      />
    </CEditorActions>
  </Form>
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@planetcrust/human-js'
import { components, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'
import { useConfirm } from 'primevue/useconfirm'

const { CInputDelete, CInputUserGroup } = components

import UserSecurity from '@/sections/admin/components/User/UserSecurity.vue'
import UserRoles from '@/sections/admin/components/User/UserRoles.vue'
import UserAvatar from '@/sections/admin/components/User/UserAvatar.vue'
import UserExternalAuth from '@/sections/admin/components/User/UserExternalAuth.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const confirm = useConfirm()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const $Auth = inject('$Auth')
// Shared user cache (provided app-wide) — keep it fresh after edits here so
// other sections (record displays, user pickers) reflect changes immediately.
const $userStore = inject('$userStore', null)

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const suspending = ref(false)
const revoking = ref(false)
const user = ref(null)
const initialUser = ref(null)
const externalAuthRef = ref(null)

// Lifted State for Tabs
const passwords = ref({
  password: '',
  confirmPassword: '',
})
const initialMembershipIDs = ref(new Set())
const membershipIDs = ref(new Set())

// Computed
const isEdit = computed(() => !!route.params.userID)

const isSelf = computed(() => !!user.value && $Auth?.user?.userID === user.value.userID)

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
    errors.handle = [{ message: t('system.users.editor.info.invalid-handle-characters') }]
  }

  return { errors }
})

// Methods
async function loadUser() {
  const userID = route.params.userID
  if (!userID) {
    // Create new
    user.value = new system.User({})
    initialMembershipIDs.value = new Set()
    membershipIDs.value = new Set()
    passwords.value = { password: '', confirmPassword: '' }

    // Preselect the default user group
    await fetchDefaultUserGroup()

    initialUser.value = cloneDeep(user.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.userRead({ userID })
    user.value = new system.User(raw)

    // Preselect the default user group if none is currently selected
    if (!user.value.userGroupID || user.value.userGroupID === '0') {
      await fetchDefaultUserGroup()
    }

    // Load initial roles if editing
    const memRes = await $SystemAPI.userMembershipList({ userID })
    const ids = Array.isArray(memRes) ? memRes : (memRes.set || []).map(m => m.roleID)
    initialMembershipIDs.value = new Set(ids)
    membershipIDs.value = new Set(ids)
    initialUser.value = cloneDeep(user.value)
  } catch (e) {
    console.error('Failed to load user:', e)
    $toast.toastErrorHandler(t('notification.user.fetch.error'))(e)
    router.push({ name: 'system.users' })
  } finally {
    loading.value = false
  }
}

async function fetchDefaultUserGroup() {
  try {
    const result = await $SystemAPI.userGroupList({ limit: 100 })
    if (result?.set?.length > 0) {
      // Find 'default-root', 'users', or anything with 'Default' in name, otherwise fallback to the first available group
      const defaultGroup = result.set.find(g => 
        g.handle === 'default-root' || 
        g.handle === 'users' || 
        g.meta?.short?.includes('Default')
      ) || result.set[0]
      
      if (defaultGroup) {
        user.value.userGroupID = defaultGroup.userGroupID
      }
    }
  } catch (e) {
    console.warn('Silent fail fetching default user group', e)
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document.querySelector('.p-message-error')?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

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
      initialUser.value = cloneDeep(user.value)
      $userStore?.storeUsers([user.value])

      // Handle Password if provided
      if (passwords.value.password) {
        await $SystemAPI.userSetPassword({
          userID: payload.userID,
          password: passwords.value.password,
        })
        passwords.value = { password: '', confirmPassword: '' }
        await externalAuthRef.value?.loadCredentials()
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
      // Set the user instance to the newly created user to immediately populate the userID
      // preventing components like UserExternalAuth from fetching with an invalid ID
      // during the router transition.
      user.value = new system.User(created)
      $userStore?.storeUsers([user.value])
      $toast.toastSuccess(t('notification.user.create.success'))
      markSaved()
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
    $userStore?.removeUsers(user.value.userID)
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

async function handleRevokeSessions() {
  revoking.value = true
  try {
    await $SystemAPI.userSessionsRemove({ userID: user.value.userID })
    $toast.toastSuccess(t('notification.user.sessionsRevoke.success'))
  } catch (e) {
    console.error('Failed to revoke sessions:', e)
    $toast.toastErrorHandler(t('notification.user.sessionsRevoke.error'))(e)
  } finally {
    revoking.value = false
  }
}

function confirmRevokeSessions(event) {
  confirm.require({
    target: event.currentTarget,
    message: t('system.users.editor.info.revokeAllSession') + '?',
    header: t('system.users.editor.info.revokeAllSession'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      outlined: true,
      size: 'small',
    },
    acceptProps: {
      label: t('system.users.editor.info.revokeAllSession'),
      severity: 'warn',
      size: 'small',
    },
    accept: () => {
      handleRevokeSessions()
    },
  })
}

const { markSaved } = useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!user.value && !!initialUser.value && (!isEqual(user.value, initialUser.value) || !isEqual([...membershipIDs.value], [...initialMembershipIDs.value])),
  messageKey: 'general.editor.unsavedChanges',
})

watch(
  () => route.params.userID,
  (newID, oldID) => {
    if (newID !== oldID) {
      loadUser()
    }
  }
)

onMounted(() => {
  loadUser()
})
</script>
