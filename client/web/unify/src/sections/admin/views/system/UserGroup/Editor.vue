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
    v-else-if="userGroup"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2 shrink-0">
        <CPermissionsButton
          v-if="userGroup.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:user-group/${userGroup.userGroupID}`"
          :title="userGroup.meta?.short || userGroup.handle || userGroup.userGroupID"
          :target="userGroup.meta?.short || userGroup.handle || userGroup.userGroupID"
        />
      </div>
      <Panel :header="$t('system.user-groups.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <CFormGroup name="name" :label="$t('system.user-groups.editor.info.meta.short')" required>
            <InputText id="name" name="name" v-model="userGroup.meta.short" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('system.user-groups.editor.info.handle')">
            <InputText id="handle" name="handle" v-model="userGroup.handle" />
          </CFormGroup>

          <CFormGroup name="description" :label="$t('system.user-groups.editor.info.meta.description')" class="md:col-span-2">
            <Textarea
              id="description"
              name="description"
              v-model="userGroup.meta.description"
              rows="3"
            />
          </CFormGroup>

          <!-- Parent hierarchy (non-root groups) -->
          <CFormGroup
            v-if="!userGroup.isRoot"
            :label="$t('system.user-groups.editor.info.parents.title')"
            class="md:col-span-2"
          >
            <template #actions>
              <Button
                icon="pi pi-plus"
                :label="$t('general.label.add')"
                severity="secondary"
                size="small"
                @click="addParent"
              />
            </template>
            <CFormList
              v-if="userGroup.config?.path"
              v-model="userGroup.config.path"
              :min-items="1"
              :empty-message="$t('system.user-groups.editor.info.parents.empty')"
              :columns="[
                { label: $t('system.user-groups.editor.info.parents.parent.label'), width: '1fr' },
                { label: $t('system.user-groups.editor.info.parents.name.label'), width: '1fr' },
              ]"
            >
              <template #row="{ item }">
                <CInputUserGroup
                  v-model="item.selfID"
                  class="w-full"
                  :placeholder="$t('system.user-groups.editor.info.parents.parent.placeholder')"
                />
                <InputText
                  v-model="item.name"
                  :placeholder="$t('system.user-groups.editor.info.parents.name.placeholder')"
                  class="w-full"
                />
              </template>
            </CFormList>
          </CFormGroup>
        </div>
      </Panel>

      <Panel
        v-if="isEdit"
        :header="$t('system.user-groups.editor.members.title')"
        toggleable
        class="shadow"
      >
        <UserGroupMembers :userGroupID="userGroup.userGroupID" />
      </Panel>

      <Panel
        v-if="isEdit"
        :header="$t('system.user-groups.editor.roles.title')"
        toggleable
        class="shadow"
      >
        <UserGroupRoles :userGroupID="userGroup.userGroupID" />
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'system.userGroups' }">
      <CInputDelete
        v-if="isEdit && userGroup.canDeleteUserGroup && !userGroup.deletedAt"
        :label="$t('system.user-groups.editor.info.delete')"
        :message="$t('general.confirm.delete')"
        :header="userGroup.meta.short || userGroup.handle || userGroup.userGroupID"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        v-if="isEdit && userGroup.deletedAt"
        :label="$t('system.user-groups.editor.info.undelete')"
        icon="pi pi-refresh"
        severity="success"
        :disabled="saving"
        @click="handleUndelete"
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
import UserGroupMembers from '@/sections/admin/components/UserGroup/UserGroupMembers.vue'
import UserGroupRoles from '@/sections/admin/components/UserGroup/UserGroupRoles.vue'

const { CInputDelete, CInputUserGroup } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const userGroup = ref(null)
const initialUserGroup = ref(null)

// Computed
const isEdit = computed(() => !!route.params.userGroupID)

const pageTitle = computed(() => {
  return isEdit.value
    ? t('system.user-groups.editor.title.edit')
    : t('system.user-groups.editor.title.create')
})

const initialValues = computed(() => {
  return {
    name: userGroup.value?.meta?.short || '',
    handle: userGroup.value?.handle || '',
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('system.user-groups.editor.info.invalid-handle-characters') }]
  }

  return { errors }
})

// Parent hierarchy management
function addParent() {
  if (!userGroup.value.config) userGroup.value.config = { path: [] }
  if (!userGroup.value.config.path) userGroup.value.config.path = []
  userGroup.value.config.path.push({ selfID: '', name: '' })
}

// Methods
async function loadUserGroup() {
  const userGroupID = route.params.userGroupID
  if (!userGroupID) {
    userGroup.value = new system.UserGroup({})
    try {
      const result = await $SystemAPI.userGroupList({ limit: 100 })
      const groups = result?.set || []
      const defaultGroup = groups.find(g => g.handle === 'default-root' || g.isRoot) || groups[0]
      if (defaultGroup) {
        userGroup.value.config = { path: [{ selfID: defaultGroup.userGroupID, name: '' }] }
      }
    } catch {
      // silent — user can pick manually
    }
    initialUserGroup.value = cloneDeep(userGroup.value)
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.userGroupRead({ userGroupID })
    userGroup.value = new system.UserGroup(raw)
    initialUserGroup.value = cloneDeep(userGroup.value)
  } catch (e) {
    console.error('Failed to load user group:', e)
    $toast.toastErrorHandler(t('notification.userGroup.fetch.error'))(e)
    router.push({ name: 'system.userGroups' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }

  if (isEdit.value && !userGroup.value?.canUpdateUserGroup) return

  saving.value = true
  try {
    const payload = {
      handle: userGroup.value.handle,
      meta: userGroup.value.meta,
      config: userGroup.value.config,
    }

    if (isEdit.value) {
      payload.userGroupID = userGroup.value.userGroupID
      const raw = await $SystemAPI.userGroupUpdate(payload)
      userGroup.value = new system.UserGroup(raw)
      initialUserGroup.value = cloneDeep(userGroup.value)
      $toast.toastSuccess(t('notification.userGroup.update.success'))
    } else {
      const created = await $SystemAPI.userGroupCreate(payload)
      $toast.toastSuccess(t('notification.userGroup.create.success'))
      markSaved()
      router.push({
        name: 'system.userGroups.edit',
        params: { userGroupID: created.userGroupID },
      })
    }
  } catch (e) {
    console.error('Failed to save user group:', e)
    $toast.toastErrorHandler(
      isEdit.value
        ? t('notification.userGroup.update.error')
        : t('notification.userGroup.create.error'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.userGroupDelete({ userGroupID: userGroup.value.userGroupID })
    $toast.toastSuccess(t('notification.userGroup.delete.success'))
    router.push({ name: 'system.userGroups' })
  } catch (e) {
    console.error('Failed to delete user group:', e)
    $toast.toastErrorHandler(t('notification.userGroup.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

async function handleUndelete() {
  saving.value = true
  try {
    await $SystemAPI.userGroupUndelete({ userGroupID: userGroup.value.userGroupID })
    const raw = await $SystemAPI.userGroupRead({ userGroupID: userGroup.value.userGroupID })
    userGroup.value = new system.UserGroup(raw)
    $toast.toastSuccess(t('notification.userGroup.undelete.success'))
  } catch (e) {
    console.error('Failed to undelete user group:', e)
    $toast.toastErrorHandler(t('notification.userGroup.undelete.error'))(e)
  } finally {
    saving.value = false
  }
}

const { markSaved } = useUnsavedGuard({
  isDirty: () =>
    !saving.value &&
    !deleting.value &&
    !!userGroup.value &&
    !!initialUserGroup.value &&
    !isEqual(userGroup.value, initialUserGroup.value),
  messageKey: 'general.editor.unsavedChanges',
})

watch(
  () => route.params.userGroupID,
  (newID, oldID) => {
    if (newID !== oldID) {
      loadUserGroup()
    }
  },
)

onMounted(() => {
  loadUserGroup()
})
</script>
