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
    v-else-if="role"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <div v-if="isEdit" class="flex justify-end gap-2">
        <Button
          v-if="role.canGrant"
          :label="$t('system.roles.editor.clone.title')"
          icon="pi pi-copy"
          severity="secondary"
          size="small"
          outlined
          @click="showCloneDialog = true"
        />
        <CPermissionsButton
          v-if="role.canGrant"
          v-tooltip.bottom="$t('general.label.permissions')"
          :resource="`corteza::system:role/${role.roleID}`"
          :title="role.name || role.handle || role.roleID"
          :target="role.name || role.handle || role.roleID"
        />
      </div>

      <Panel :header="$t('system.roles.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="name" class="flex flex-col gap-2">
            <label for="name" class="font-medium text-primary">
              {{ $t('system.roles.editor.info.name') }}
              <span class="text-red-500">*</span>
            </label>
            <InputText id="name" name="name" v-model="role.name" />
            <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
              {{ $form.name.error?.message }}
            </Message>
          </FormField>

          <FormField name="handle" class="flex flex-col gap-2">
            <label for="handle" class="font-medium text-primary">
              {{ $t('system.roles.editor.info.handle') }}
            </label>
            <InputText id="handle" name="handle" v-model="role.handle" />
            <Message v-if="$form.handle?.invalid" severity="error" size="small" variant="simple">
              {{ $form.handle.error?.message }}
            </Message>
          </FormField>

          <FormField name="description" class="flex flex-col gap-2 md:col-span-2">
            <label for="description" class="font-medium text-primary">
              {{ $t('system.roles.editor.info.description') }}
            </label>
            <Textarea
              id="description"
              name="description"
              v-model="role.meta.description"
              rows="3"
            />
          </FormField>

          <!-- isContextual toggle -->
          <div class="flex items-center gap-3">
            <ToggleSwitch id="isContextual" v-model="isContextual" />
            <label for="isContextual" class="font-medium text-primary cursor-pointer">
              {{ $t('system.roles.editor.info.context.label') }}
            </label>
          </div>

          <!-- Contextual section (shown when isContextual) -->
          <div
            v-if="isContextual"
            class="md:col-span-2 flex flex-col gap-4 p-4 border rounded-lg bg-surface"
          >
            <div class="flex flex-col gap-2">
              <label class="font-medium text-primary">
                {{ $t('system.roles.editor.info.context.expression-label') }}
              </label>
              <InputText v-model="role.meta.context.expr" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="font-medium text-primary">
                {{ $t('system.roles.editor.info.context.resource-types-label') }}
              </label>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-2">
                <div v-for="rt in resourceTypes" :key="rt.value" class="flex items-center gap-2">
                  <Checkbox
                    :inputId="rt.value"
                    v-model="role.meta.context.resourceTypes"
                    :value="rt.value"
                  />
                  <label :for="rt.value" class="cursor-pointer text-sm">{{ rt.label }}</label>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Panel>

      <Panel
        v-if="isEdit && !isContextual"
        :header="$t('system.roles.editor.members.title')"
        toggleable
        :collapsed="false"
      >
        <RoleMembers
          :role="role"
          :initialMemberIDs="initialMemberIDs"
          v-model:memberIDs="memberIDs"
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
          @click="$router.push({ name: 'system.roles' })"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && role.canDeleteRole && !role.deletedAt && !role.isSystem"
            :label="$t('system.roles.editor.info.delete')"
            :message="$t('general.confirm.delete')"
            :header="role.name || role.handle || role.roleID"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            v-if="isEdit && role.deletedAt && !role.isSystem"
            :label="$t('system.roles.editor.info.undelete')"
            icon="pi pi-refresh"
            severity="success"
            :disabled="saving"
            @click="handleUndelete"
          />
          <Button
            v-if="isEdit && !role.archivedAt && !role.isSystem"
            :label="$t('system.roles.editor.info.archive')"
            icon="pi pi-box"
            severity="warn"
            outlined
            :disabled="saving"
            @click="handleArchive"
          />
          <Button
            v-if="isEdit && role.archivedAt && !role.isSystem"
            :label="$t('system.roles.editor.info.unarchive')"
            icon="pi pi-box"
            severity="success"
            outlined
            :disabled="saving"
            @click="handleUnarchive"
          />
          <Button
            v-if="!role.isSystem || role.canUpdateRole"
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
          />
        </div>
      </div>
    </div>
  </Form>

  <!-- Clone Permissions Dialog -->
  <RolePermissionClone
    v-if="isEdit && role"
    v-model:visible="showCloneDialog"
    :roleID="role.roleID"
  />
</template>

<script setup>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { system } from '@cortezaproject/corteza-js-next'
import { components, usePermissions, useUnsavedGuard } from '@cortezaproject/corteza-vue-next'
import { cloneDeep, isEqual } from 'lodash-es'
import RoleMembers from '@/components/Role/RoleMembers.vue'
import RolePermissionClone from '@/components/Role/RolePermissionClone.vue'

const { CInputDelete } = components

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const role = ref(null)
const initialRole = ref(null)
const memberIDs = ref(new Set())
const initialMemberIDs = ref(new Set())
const showCloneDialog = ref(false)

// Computed
const isEdit = computed(() => !!route.params.roleID)

const isContextual = ref(false)

watch(isContextual, val => {
  if (val) {
    memberIDs.value = new Set()
  } else if (role.value?.meta?.context) {
    role.value.meta.context.expr = ''
    role.value.meta.context.resourceTypes = []
  }
})

const resourceTypes = [
  { value: 'corteza::system:auth-client', label: 'system:auth-client' },
  { value: 'corteza::system:role', label: 'system:role' },
  { value: 'corteza::system:user', label: 'system:user' },
  { value: 'corteza::compose:module', label: 'compose:module' },
  { value: 'corteza::compose:namespace', label: 'compose:namespace' },
  { value: 'corteza::compose:page', label: 'compose:page' },
  { value: 'corteza::compose:record', label: 'compose:record' },
  { value: 'corteza::automation:workflow', label: 'automation:workflow' },
]

const pageTitle = computed(() => {
  return isEdit.value ? t('system.roles.editor.title.edit') : t('system.roles.editor.title.create')
})

const initialValues = computed(() => {
  return {
    name: role.value?.name || '',
    handle: role.value?.handle || '',
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('system.roles.editor.info.invalid-handle-characters') }]
  }

  return { errors }
})

// Methods
async function loadRole() {
  const roleID = route.params.roleID
  if (!roleID) {
    // Create new
    role.value = new system.Role({ meta: { context: { expr: '', resourceTypes: [] } } })
    initialRole.value = cloneDeep(role.value)
    memberIDs.value = new Set()
    initialMemberIDs.value = new Set()
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.roleRead({ roleID })
    role.value = new system.Role(raw)

    if (!role.value.meta.context) {
      role.value.meta.context = { expr: '', resourceTypes: [] }
    }

    isContextual.value = !!(
      role.value.meta.context.expr || role.value.meta.context.resourceTypes?.length
    )

    if (!isContextual.value) {
      // Load member IDs
      const membersResult = await $SystemAPI.roleMemberList({ roleID })
      const ids = new Set((membersResult?.set || membersResult || []).map(u => u.userID || u))
      memberIDs.value = ids
      initialMemberIDs.value = new Set(ids)
    } else {
      memberIDs.value = new Set()
      initialMemberIDs.value = new Set()
    }
    initialRole.value = cloneDeep(role.value)
  } catch (e) {
    console.error('Failed to load role:', e)
    $toast.toastErrorHandler(t('notification.role.fetch.error'))(e)
    router.push({ name: 'system.roles' })
  } finally {
    loading.value = false
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

  if (isEdit.value && !role.value?.canUpdateRole) return

  saving.value = true
  try {
    const payload = {
      name: role.value.name,
      handle: role.value.handle,
      meta: role.value.meta,
    }

    if (isEdit.value) {
      payload.roleID = role.value.roleID
      const raw = await $SystemAPI.roleUpdate(payload)
      role.value = new system.Role(raw)
      initialRole.value = cloneDeep(role.value)

      // Sync member changes
      if (!isContextual.value) {
        const added = [...memberIDs.value].filter(id => !initialMemberIDs.value.has(id))
        const removed = [...initialMemberIDs.value].filter(id => !memberIDs.value.has(id))
        await Promise.all([
          ...added.map(userID => $SystemAPI.roleMemberAdd({ roleID: role.value.roleID, userID })),
          ...removed.map(userID =>
            $SystemAPI.roleMemberRemove({ roleID: role.value.roleID, userID }),
          ),
        ])
        initialMemberIDs.value = new Set(memberIDs.value)
      } else if (initialMemberIDs.value.size > 0) {
        // Roles that are contextual shouldn't retain explicit members
        const removed = [...initialMemberIDs.value]
        await Promise.all(
          removed.map(userID => $SystemAPI.roleMemberRemove({ roleID: role.value.roleID, userID })),
        )
        initialMemberIDs.value = new Set()
      }

      $toast.toastSuccess(t('notification.role.update.success'))
    } else {
      const created = await $SystemAPI.roleCreate(payload)
      $toast.toastSuccess(t('notification.role.create.success'))
      router.push({
        name: 'system.roles.edit',
        params: { roleID: created.roleID },
      })
    }
  } catch (e) {
    console.error('Failed to save role:', e)
    $toast.toastErrorHandler(
      isEdit.value ? t('notification.role.update.error') : t('notification.role.create.error'),
    )(e)
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await $SystemAPI.roleDelete({ roleID: role.value.roleID })
    $toast.toastSuccess(t('notification.role.delete.success'))
    router.push({ name: 'system.roles' })
  } catch (e) {
    console.error('Failed to delete role:', e)
    $toast.toastErrorHandler(t('notification.role.delete.error'))(e)
  } finally {
    deleting.value = false
  }
}

async function handleUndelete() {
  saving.value = true
  try {
    await $SystemAPI.roleUndelete({ roleID: role.value.roleID })
    const raw = await $SystemAPI.roleRead({ roleID: role.value.roleID })
    role.value = new system.Role(raw)
    $toast.toastSuccess(t('notification.role.undelete.success'))
  } catch (e) {
    console.error('Failed to undelete role:', e)
    $toast.toastErrorHandler(t('notification.role.undelete.error'))(e)
  } finally {
    saving.value = false
  }
}

async function handleArchive() {
  saving.value = true
  try {
    await $SystemAPI.roleArchive({ roleID: role.value.roleID })
    const raw = await $SystemAPI.roleRead({ roleID: role.value.roleID })
    role.value = new system.Role(raw)
    $toast.toastSuccess(t('notification.role.archive.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.role.archive.error'))(e)
  } finally {
    saving.value = false
  }
}

async function handleUnarchive() {
  saving.value = true
  try {
    await $SystemAPI.roleUnarchive({ roleID: role.value.roleID })
    const raw = await $SystemAPI.roleRead({ roleID: role.value.roleID })
    role.value = new system.Role(raw)
    $toast.toastSuccess(t('notification.role.unarchive.success'))
  } catch (e) {
    $toast.toastErrorHandler(t('notification.role.unarchive.error'))(e)
  } finally {
    saving.value = false
  }
}

useUnsavedGuard({
  isDirty: () => !saving.value && !deleting.value && !!role.value && !!initialRole.value && (!isEqual(role.value, initialRole.value) || !isEqual([...memberIDs.value], [...initialMemberIDs.value])),
  messageKey: 'general.editor.unsavedChanges',
})

watch(
  () => route.params.roleID,
  (newID, oldID) => {
    if (newID !== oldID) {
      loadRole()
    }
  },
)

onMounted(() => {
  loadRole()
})
</script>
