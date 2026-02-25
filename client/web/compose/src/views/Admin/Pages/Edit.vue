<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <div v-if="isEdit" class="flex gap-2">
      <Button
        :label="$t('page.edit.pageBuilder')"
        icon="pi pi-pencil"
        size="small"
        severity="secondary"
        @click="goToBuilder"
      />
      <Button
        :label="$t('page.edit.viewPage')"
        icon="pi pi-eye"
        size="small"
        severity="secondary"
        @click="goToViewPage"
      />
    </div>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <Form
    v-else-if="page"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1">
      <Card :pt="{ body: { class: 'p-0' }, content: { class: 'p-0' } }" class="overflow-hidden">
        <template #content>
          <Tabs v-model:value="activeTab">
            <TabList class="rounded-t-lg">
              <Tab value="general">{{ $t('general.label.general') }}</Tab>
              <Tab value="settings">{{ $t('page.edit.settings') }}</Tab>
            </TabList>

            <TabPanels>
              <!-- General Tab -->
              <TabPanel value="general">
                <h4 class="font-semibold text-lg mb-4 mt-2">
                  {{ $t('page.edit.pageInfo') }}
                </h4>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                  <FormField name="title" class="flex flex-col gap-2">
                    <label for="title" class="font-medium text-primary">
                      {{ $t('page.label.title') }}
                    </label>
                    <InputText id="title" name="title" v-model="page.title" />
                    <Message
                      v-if="$form.title?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.title.error?.message }}
                    </Message>
                  </FormField>

                  <FormField name="handle" class="flex flex-col gap-2">
                    <label for="handle" class="font-medium text-primary">
                      {{ $t('page.label.handle') }}
                    </label>
                    <InputText id="handle" name="handle" v-model="page.handle" />
                    <Message
                      v-if="$form.handle?.invalid"
                      severity="error"
                      size="small"
                      variant="simple"
                    >
                      {{ $form.handle.error?.message }}
                    </Message>
                  </FormField>
                </div>

                <div class="flex flex-col gap-2">
                  <label for="description" class="font-medium text-primary">
                    {{ $t('page.label.description') }}
                  </label>
                  <Textarea id="description" v-model="page.description" rows="4" auto-resize />
                </div>
              </TabPanel>

              <!-- Settings Tab -->
              <TabPanel value="settings">
                <h4 class="font-semibold text-lg mb-4 mt-2">
                  {{ $t('page.edit.settings') }}
                </h4>

                <div class="flex flex-col gap-6">
                  <!-- Visible -->
                  <div class="flex items-center gap-3">
                    <ToggleSwitch v-model="page.visible" input-id="visible" />
                    <label for="visible" class="font-medium text-primary cursor-pointer">
                      {{ $t('page.edit.visible') }}
                    </label>
                  </div>

                  <!-- Parent Page -->
                  <div class="flex flex-col gap-2">
                    <label for="parent-page" class="font-medium text-primary">
                      {{ $t('page.edit.parentPage') }}
                    </label>
                    <Select
                      id="parent-page"
                      v-model="page.selfID"
                      :options="parentPageOptions"
                      option-label="label"
                      option-value="value"
                      :placeholder="$t('page.edit.noParent')"
                      show-clear
                    />
                  </div>
                </div>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <!-- Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.back()"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="isEdit && page.canDeletePage"
            :label="$t('general.label.delete')"
            :message="$t('page.edit.deleteConfirm')"
            :header="page.title"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            :disabled="!canSave"
          />
        </div>
      </div>
    </div>
  </Form>
</template>

<script setup>
import { usePageStore } from '@/stores/page'
import { compose } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, ref, watch } from 'vue'

const { CInputDelete } = components
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const pageStore = usePageStore()

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const page = ref(null)
const activeTab = ref('general')

// Computed
const isEdit = computed(() => !!route.params.pageID)

const pageTitle = computed(() => {
  return isEdit.value ? t('page.edit.edit') : t('page.edit.create')
})

const initialValues = computed(() => {
  return {
    title: page.value?.title || '',
    handle: page.value?.handle || '',
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.title || values.title.trim().length === 0) {
    errors.title = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('page.block.general.invalid-handle-characters') }]
  }

  return { errors }
})

const canSave = computed(() => {
  if (isEdit.value && !page.value?.canUpdatePage) return false
  return true
})

// Parent page options — exclude the current page and its descendants
const parentPageOptions = computed(() => {
  const currentID = page.value?.pageID
  const options = [{ label: t('page.edit.noParent'), value: '0' }]

  pageStore.set.forEach(p => {
    // Don't allow setting self as parent
    if (p.pageID === currentID) return
    options.push({
      label: p.title || p.handle || p.pageID,
      value: p.pageID,
    })
  })

  return options
})

// Methods
async function loadPage() {
  const pageID = route.params.pageID
  if (!pageID) {
    // Create new
    page.value = new compose.Page({
      namespaceID: props.namespace?.namespaceID,
      visible: true,
    })
    return
  }

  loading.value = true
  try {
    const found = pageStore.getByID(pageID)
    if (found) {
      page.value = new compose.Page({ ...found })
    } else {
      const p = await pageStore.findByID({
        namespaceID: props.namespace?.namespaceID,
        pageID,
      })
      page.value = new compose.Page({ ...p })
    }
  } catch (e) {
    console.error('Failed to load page:', e)
    $toast.toastDanger(t('notification.page.loadFailed'))
    router.push({ name: 'admin.pages' })
  } finally {
    loading.value = false
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return
  if (!canSave.value) return

  saving.value = true
  try {
    const payload = {
      namespaceID: props.namespace.namespaceID,
      title: page.value.title,
      handle: page.value.handle,
      description: page.value.description,
      visible: page.value.visible,
      selfID: page.value.selfID || '0',
      blocks: page.value.blocks || [],
    }

    if (isEdit.value) {
      payload.pageID = page.value.pageID
      await pageStore.update(payload)
      $toast.toastSuccess(t('notification.page.saved'))
    } else {
      const created = await pageStore.create(payload)
      $toast.toastSuccess(t('notification.page.created'))
      router.push({
        name: 'admin.pages.edit',
        params: { pageID: created.pageID },
      })
    }
  } catch (e) {
    console.error('Failed to save page:', e)
    $toast.toastDanger(t('notification.page.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await pageStore.delete({ pageID: page.value.pageID, namespaceID: props.namespace.namespaceID })
    $toast.toastSuccess(t('notification.page.deleted'))
    router.push({ name: 'admin.pages' })
  } catch (e) {
    console.error('Failed to delete page:', e)
    $toast.toastDanger(t('notification.page.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function goToBuilder() {
  router.push({
    name: 'admin.pages.builder',
    params: { pageID: page.value.pageID },
  })
}

function goToViewPage() {
  router.push({
    name: 'page',
    params: { pageID: page.value.pageID },
  })
}

// Lifecycle
onMounted(() => {
  loadPage()
})

watch(
  () => route.params.pageID,
  () => {
    loadPage()
  },
)
</script>
