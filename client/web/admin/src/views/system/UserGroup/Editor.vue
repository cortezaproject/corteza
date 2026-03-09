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
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 flex flex-col min-h-0 gap-4 overflow-y-auto">
      <Panel :header="$t('system.user-groups.editor.info.title')" toggleable :collapsed="false">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField name="name" class="flex flex-col gap-2">
            <label for="name" class="font-medium text-primary">
              {{ $t('system.user-groups.editor.info.meta.short') }}
            </label>
            <InputText id="name" name="name" v-model="userGroup.meta.short" />
            <Message v-if="$form.name?.invalid" severity="error" size="small" variant="simple">
              {{ $form.name.error?.message }}
            </Message>
          </FormField>

          <FormField name="handle" class="flex flex-col gap-2">
            <label for="handle" class="font-medium text-primary">
              {{ $t('system.user-groups.editor.info.handle') }}
            </label>
            <InputText id="handle" name="handle" v-model="userGroup.handle" />
            <Message v-if="$form.handle?.invalid" severity="error" size="small" variant="simple">
              {{ $form.handle.error?.message }}
            </Message>
          </FormField>

          <FormField name="description" class="flex flex-col gap-2 md:col-span-2">
            <label for="description" class="font-medium text-primary">
              {{ $t('system.user-groups.editor.info.meta.description') }}
            </label>
            <Textarea
              id="description"
              name="description"
              v-model="userGroup.meta.description"
              rows="3"
            />
          </FormField>
        </div>
      </Panel>
    </div>

    <!-- Bottom Actions Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="p-3 flex items-center justify-between">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'system.userGroups' })"
        />
        <div class="flex gap-2">
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
const userGroup = ref(null)

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

// Methods
async function loadUserGroup() {
  const userGroupID = route.params.userGroupID
  if (!userGroupID) {
    // Create new
    userGroup.value = new system.UserGroup({})
    return
  }

  loading.value = true
  try {
    const raw = await $SystemAPI.userGroupRead({ userGroupID })
    userGroup.value = new system.UserGroup(raw)
  } catch (e) {
    console.error('Failed to load user group:', e)
    $toast.toastErrorHandler(t('notification.userGroup.fetch.error'))(e)
    router.push({ name: 'system.userGroups' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return

  if (isEdit.value && !userGroup.value?.canUpdateUserGroup) return

  saving.value = true
  try {
    const payload = {
      handle: userGroup.value.handle,
      meta: userGroup.value.meta,
    }

    if (isEdit.value) {
      payload.userGroupID = userGroup.value.userGroupID
      const raw = await $SystemAPI.userGroupUpdate(payload)
      userGroup.value = new system.UserGroup(raw)
      $toast.toastSuccess(t('notification.userGroup.update.success'))
    } else {
      const created = await $SystemAPI.userGroupCreate(payload)
      $toast.toastSuccess(t('notification.userGroup.create.success'))
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

onMounted(() => {
  loadUserGroup()
})
</script>
