<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <div v-if="isEdit" class="flex gap-2">
      <CRouterLinkButton
        :to="{ name: 'admin.modules.record.list', params: { moduleID: module?.moduleID } }"
        :label="$t('module.allRecords.label')"
        icon="pi pi-table"
        size="small"
        severity="secondary"
      />
    </div>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <Form
    v-else-if="module"
    v-slot="$form"
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1">
      <!-- Related Pages Actions -->
      <div v-if="isEdit && namespace?.canManageNamespace" class="flex justify-end gap-2 mb-4">
        <!-- Record Page Button -->
        <CRouterLinkButton
          v-if="recordPage"
          :to="{ name: 'admin.pages.builder', params: { pageID: recordPage.pageID } }"
          :label="$t('module.recordPage.edit')"
          icon="pi pi-file-edit"
          size="small"
          severity="secondary"
          outlined
        />
        <Button
          v-else
          :label="$t('module.recordPage.create')"
          icon="pi pi-file-plus"
          size="small"
          severity="secondary"
          outlined
          :loading="creatingRecordPage"
          @click="handleRecordPageCreation"
        />

        <!-- Record List Page Button -->
        <CRouterLinkButton
          v-if="recordListPage"
          :to="{ name: 'admin.pages.builder', params: { pageID: recordListPage.pageID } }"
          :label="$t('module.recordListPage.edit')"
          icon="pi pi-list"
          size="small"
          severity="secondary"
          outlined
        />
        <Button
          v-else
          :label="$t('module.recordListPage.create')"
          icon="pi pi-list"
          size="small"
          severity="secondary"
          outlined
          :loading="creatingRecordListPage"
          :disabled="!recordPage"
          @click="handleRecordListPageCreation"
        />
      </div>

      <Card :pt="{ body: { class: 'p-0' }, content: { class: 'p-0' } }" class="overflow-hidden">
        <template #content>
          <Tabs v-model:value="activeTab">
            <TabList class="rounded-t-lg">
              <Tab value="fields">{{ $t('module.edit.fields.label') }}</Tab>
            </TabList>

            <TabPanels>
              <TabPanel value="fields">
                <!-- Module Info -->
                <h4 class="font-semibold text-lg mb-4">
                  {{ $t('module.edit.moduleInfo') }}
                </h4>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                  <FormField name="name" class="flex flex-col gap-2">
                    <label for="name" class="font-medium text-primary">
                      {{ $t('module.general.label.name') }}
                    </label>
                    <InputText id="name" name="name" v-model="module.name" />
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
                      {{ $t('module.general.label.handle') }}
                    </label>
                    <InputText id="handle" name="handle" v-model="module.handle" />
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

                <Divider />

                <!-- Module Fields -->
                <div class="flex items-center justify-between mb-4">
                  <h4 class="font-semibold text-lg">
                    {{ $t('module.edit.manageRecordFields') }}
                  </h4>
                  <Button
                    :label="$t('module.edit.newField')"
                    icon="pi pi-plus"
                    size="small"
                    @click="addField"
                  />
                </div>

                <DataTable
                  :value="module.fields"
                  striped-rows
                  data-key="name"
                  class="border border-b-0 border-surface rounded-border overflow-auto"
                  :pt="{ headerCell: { class: 'bg-highlight-emphasis' } }"
                >
                  <Column field="name" :header="$t('module.edit.fields.columns.name.label')">
                    <template #body="{ data }">
                      <InputText v-model="data.name" class="w-full" size="small" />
                    </template>
                  </Column>

                  <Column field="label" :header="$t('module.edit.fields.columns.title.label')">
                    <template #body="{ data }">
                      <InputText v-model="data.label" class="w-full" size="small" />
                    </template>
                  </Column>

                  <Column field="kind" :header="$t('module.edit.fields.columns.type.label')">
                    <template #body="{ data, index }">
                      <InputGroup>
                        <Select
                          v-model="data.kind"
                          :options="fieldKinds"
                          option-label="label"
                          option-value="value"
                          size="small"
                        />
                        <InputGroupAddon>
                          <Button
                            icon="pi pi-cog"
                            severity="secondary"
                            size="small"
                            class="w-full border-none"
                            @click="openFieldConfigurator(data, index)"
                          />
                        </InputGroupAddon>
                      </InputGroup>
                    </template>
                  </Column>

                  <Column
                    field="isRequired"
                    :header="$t('module.edit.fields.columns.required.label')"
                    header-style="width: 5rem"
                    header-class="text-center"
                    body-class="text-center"
                  >
                    <template #body="{ data }">
                      <div class="flex justify-center">
                        <Checkbox v-model="data.isRequired" :binary="true" />
                      </div>
                    </template>
                  </Column>

                  <Column
                    field="isMulti"
                    :header="$t('module.edit.fields.columns.multi.label')"
                    header-style="width: 5rem"
                    header-class="text-center"
                    body-class="text-center"
                  >
                    <template #body="{ data }">
                      <div class="flex justify-center">
                        <Checkbox v-model="data.isMulti" :binary="true" />
                      </div>
                    </template>
                  </Column>

                  <Column header-style="width: 3rem">
                    <template #body="{ data, index }">
                      <div class="flex justify-end gap-1">
                        <Button
                          icon="pi pi-ellipsis-v"
                          text
                          severity="secondary"
                          size="small"
                          class="row-action-btn w-full"
                          @click.stop="toggleFieldActionsMenu($event, data, index)"
                        />
                      </div>
                    </template>
                  </Column>

                  <template #empty>
                    <div class="text-center py-4 text-muted-color">
                      {{ $t('module.noModule') }}
                    </div>
                  </template>
                </DataTable>

                <TieredMenu ref="fieldActionsMenu" :model="fieldActionsMenuItems" popup>
                  <template #item="{ item, props }">
                    <a v-ripple v-bind="props.action" :class="item.class">
                      <span :class="item.icon" />
                      <span class="ml-2">{{ item.label }}</span>
                    </a>
                  </template>
                </TieredMenu>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

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
            v-if="isEdit && module.canDeleteModule"
            :label="$t('general.label.delete')"
            :message="$t('module.edit.deleteConfirm')"
            :header="module.name"
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

    <!-- Field Configurator Modal -->
    <CFieldConfigurator
      v-model:visible="configuratorVisible"
      :field="activeConfiguratorField"
      @save="onFieldSave"
    />
  </Form>
</template>

<script setup>
import { useModuleStore } from '@/stores/module'
import { usePageStore } from '@/stores/page'
import { compose } from '@cortezaproject/corteza-js-next'
import { components, useConfirmDelete } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onMounted, ref, watch } from 'vue'
import CFieldConfigurator from '@/components/ModuleFields/Configurator/index.vue'

const { CInputDelete, CRouterLinkButton } = components
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
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const moduleStore = useModuleStore()
const pageStore = usePageStore()

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const module = ref(null)
const activeTab = ref('fields')
const creatingRecordPage = ref(false)
const creatingRecordListPage = ref(false)

// Configurator State
const configuratorVisible = ref(false)
const activeConfiguratorField = ref(null)
const activeConfiguratorFieldIndex = ref(-1)

// Field actions menu
const fieldActionsMenu = ref()
const fieldActionsMenuItems = ref([])

// Field type options
const fieldKinds = [
  { label: t('general.fieldKinds.String.label'), value: 'String' },
  { label: t('general.fieldKinds.Number.label'), value: 'Number' },
  { label: t('general.fieldKinds.Bool.label'), value: 'Bool' },
  { label: t('general.fieldKinds.DateTime.label'), value: 'DateTime' },
  { label: t('general.fieldKinds.Select.label'), value: 'Select' },
  { label: t('general.fieldKinds.Email.label'), value: 'Email' },
  { label: t('general.fieldKinds.Url.label'), value: 'Url' },
  { label: t('general.fieldKinds.File.label'), value: 'File' },
  { label: t('general.fieldKinds.User.label'), value: 'User' },
  { label: t('general.fieldKinds.Record.label'), value: 'Record' },
  { label: t('general.fieldKinds.Geometry.label'), value: 'Geometry' },
]

// Computed
const isEdit = computed(() => !!route.params.moduleID)

const pageTitle = computed(() => {
  return isEdit.value ? t('module.edit.edit') : t('module.edit.create')
})

const initialValues = computed(() => {
  return {
    name: module.value?.name || '',
    handle: module.value?.handle || '',
  }
})

const resolver = ref(({ values }) => {
  const errors = {}

  if (!values.name || values.name.trim().length === 0) {
    errors.name = [{ message: t('general.label.required') }]
  }

  if (values.handle && !/^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(values.handle)) {
    errors.handle = [{ message: t('module.general.placeholder.invalid-handle-characters') }]
  }

  return { errors }
})

const canSave = computed(() => {
  if (isEdit.value && !module.value?.canUpdateModule) return false
  return true
})

// Related Pages - find existing pages for this module
const recordPage = computed(() => {
  if (!module.value?.moduleID) return null
  return pageStore.set.find(p => p.moduleID === module.value.moduleID)
})

const recordListPage = computed(() => {
  if (!module.value?.moduleID) return null
  return pageStore.set.find(p => {
    return p.blocks?.some(
      b => b.kind === 'RecordList' && b.options?.moduleID === module.value.moduleID,
    )
  })
})

// Methods
async function loadModule() {
  const moduleID = route.params.moduleID
  if (!moduleID) {
    // Create new
    module.value = new compose.Module({
      namespaceID: props.namespace?.namespaceID,
      fields: [],
    })
    return
  }

  loading.value = true
  try {
    const found = moduleStore.getByID(moduleID)
    if (found) {
      module.value = new compose.Module({ ...found })
    } else {
      const m = await moduleStore.findByID({
        namespaceID: props.namespace?.namespaceID,
        moduleID,
      })
      module.value = new compose.Module({ ...m })
    }
  } catch (e) {
    console.error('Failed to load module:', e)
    $toast.toastDanger(t('notification.module.loadFailed'))
    router.push({ name: 'admin.modules' })
  } finally {
    loading.value = false
  }
}

function addField() {
  if (!module.value.fields) {
    module.value.fields = []
  }
  module.value.fields.push(new compose.ModuleFieldString())
}

function removeField(index) {
  module.value.fields.splice(index, 1)
}

function toggleFieldActionsMenu(event, field, index) {
  fieldActionsMenuItems.value = getFieldActionsMenuItems(field, index)
  fieldActionsMenu.value.toggle(event)
}

function getFieldActionsMenuItems(field, index) {
  const items = []

  items.push({
    label: t('general.label.delete'),
    icon: 'pi pi-trash',
    class: 'text-red-500',
    command: () => onConfirmFieldDelete(field, index),
  })

  return items
}

function onConfirmFieldDelete(field, index) {
  confirmDelete({
    message: t('module.edit.fields.deleteConfirm'),
    header: field.label || field.name || t('module.edit.fields.columns.name.label'),
    onConfirm: () => removeField(index),
  })
}

function openFieldConfigurator(field, index) {
  activeConfiguratorField.value = field
  activeConfiguratorFieldIndex.value = index
  configuratorVisible.value = true
}

function onFieldSave(updatedField) {
  if (activeConfiguratorFieldIndex.value > -1) {
    module.value.fields.splice(activeConfiguratorFieldIndex.value, 1, updatedField)
  }
}

async function handleSubmit({ valid }) {
  if (!valid) return
  if (!canSave.value) return

  saving.value = true
  try {
    const payload = {
      namespaceID: props.namespace.namespaceID,
      name: module.value.name,
      handle: module.value.handle,
      fields: module.value.fields,
    }

    if (isEdit.value) {
      payload.moduleID = module.value.moduleID
      await moduleStore.update(payload)
      $toast.toastSuccess(t('notification.module.saved'))
    } else {
      const created = await moduleStore.create(payload)
      $toast.toastSuccess(t('notification.module.created'))
      router.push({
        name: 'admin.modules.edit',
        params: { moduleID: created.moduleID },
      })
    }
  } catch (e) {
    console.error('Failed to save module:', e)
    $toast.toastDanger(t('notification.module.saveFailed'))
  } finally {
    saving.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    await moduleStore.delete({ moduleID: module.value.moduleID })
    $toast.toastSuccess(t('notification.module.deleted'))
    router.push({ name: 'admin.modules' })
  } catch (e) {
    console.error('Failed to delete module:', e)
    $toast.toastDanger(t('notification.module.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

async function handleRecordPageCreation() {
  creatingRecordPage.value = true
  try {
    const { name, moduleID } = module.value
    const { namespaceID } = props.namespace

    // Create a simple record page with a Record block
    const blocks = [new compose.PageBlockRecord({ xywh: [0, 0, 48, 82] })]
    const selfID = recordListPage.value?.pageID || '0'

    const page = new compose.Page({
      namespaceID,
      moduleID,
      selfID,
      title: t('module.forModule.recordPage', { name }),
      blocks,
    })

    await pageStore.create(page)
    $toast.toastSuccess(t('notification.page.created'))
  } catch (e) {
    console.error('Failed to create record page:', e)
    $toast.toastDanger(t('notification.page.createFailed'))
  } finally {
    creatingRecordPage.value = false
  }
}

async function handleRecordListPageCreation() {
  creatingRecordListPage.value = true
  try {
    const { name, moduleID } = module.value
    const { namespaceID } = props.namespace

    // Create a page with a RecordList block for this module
    const blocks = [
      new compose.PageBlockRecordList({
        xywh: [0, 0, 48, 82],
        options: {
          moduleID,
          fields: [],
        },
      }),
    ]

    const page = new compose.Page({
      title: name,
      namespaceID,
      blocks,
      visible: true,
    })

    const createdPage = await pageStore.create(page)

    // Update the record page to set this as its parent
    if (recordPage.value) {
      await pageStore.update({
        ...recordPage.value,
        selfID: createdPage.pageID,
      })
    }

    $toast.toastSuccess(t('notification.page.created'))
  } catch (e) {
    console.error('Failed to create record list page:', e)
    $toast.toastDanger(t('notification.page.createFailed'))
  } finally {
    creatingRecordListPage.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadModule()
})

watch(
  () => route.params.moduleID,
  () => {
    loadModule()
  },
)
</script>
