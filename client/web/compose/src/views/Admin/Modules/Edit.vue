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
    <div class="container mx-auto p-4 flex-1 overflow-y-auto min-h-0 flex flex-col">
      <!-- Related Pages Actions -->
      <div
        v-if="isEdit && namespace?.canManageNamespace"
        class="flex justify-end gap-2 mb-4 shrink-0"
      >
        <!-- Discovery & Federation Buttons -->
        <Button
          v-if="module?.moduleID"
          :label="$t('module.edit.schemaAlterations.title', 'Schema Alterations')"
          icon="pi pi-database"
          size="small"
          severity="secondary"
          outlined
          @click="checkSchemaAlterations"
        />
        <Button
          :label="$t('module.edit.discoverySettings.title', 'Discovery')"
          icon="pi pi-globe"
          size="small"
          severity="secondary"
          outlined
          @click="discoveryModal = true"
        />
        <Button
          :label="$t('module.edit.federationSettings.title', 'Federation')"
          icon="pi pi-share-alt"
          size="small"
          severity="secondary"
          outlined
          @click="federationModal = true"
        />

        <!-- Export Button -->
        <Button
          v-if="namespace?.canExportModules"
          :label="$t('general.label.export')"
          icon="pi pi-download"
          size="small"
          severity="secondary"
          outlined
          @click="exportModule"
        />

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

        <Button
          v-if="isEdit && module.canGrant"
          type="button"
          icon="pi pi-lock"
          v-tooltip.bottom="$t('general.label.permissions')"
          size="small"
          severity="secondary"
          @click="togglePermissionsMenu"
          aria-haspopup="true"
          aria-controls="permissions_menu"
          outlined
        />
        <Menu
          ref="permissionsMenu"
          id="permissions_menu"
          :model="permissionsMenuItems"
          :popup="true"
        />
      </div>

      <Card
        :pt="{ body: { class: 'p-0 h-full' }, content: { class: 'p-0 h-full' } }"
        class="flex-1"
      >
        <template #content>
          <Tabs v-model:value="activeTab">
            <TabList class="rounded-t-lg overflow-x-auto whitespace-nowrap">
              <Tab value="fields">{{ $t('module.edit.fields.label') }}</Tab>
              <Tab value="dal">{{ $t('module.edit.config.dal.title', 'Data Store') }}</Tab>
              <Tab value="unique">
                {{ $t('module.edit.config.uniqueValues.title', 'Unique Values') }}
              </Tab>
              <Tab value="revisions">
                {{ $t('module.edit.config.record-revisions.title', 'Record Revisions') }}
              </Tab>
              <Tab value="privacy">{{ $t('module.edit.config.privacy.title', 'Privacy') }}</Tab>
              <Tab value="issues">{{ $t('module.edit.issues.tab', 'Issues') }}</Tab>
            </TabList>

            <TabPanels>
              <TabPanel value="fields">
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
                  <Button
                    :label="$t('module.edit.newField')"
                    icon="pi pi-plus"
                    size="small"
                    @click="addField"
                  />
                </div>

                <CResourceTable
                  ref="fieldTableRef"
                  :items="allFieldsForTable"
                  :fields="fieldTableColumns"
                  :action-items="getFieldActionsMenuItems"
                  :scroll-height="tableScrollHeight"
                  empty-message="—"
                >
                  <template #body-name="{ data }">
                    <span v-if="data.isSystem" class="text-muted-color">{{ data.name }}</span>
                    <InputText v-else v-model="data.name" class="w-full" size="small" />
                  </template>

                  <template #body-label="{ data }">
                    <span v-if="data.isSystem" class="text-muted-color">{{ data.label }}</span>
                    <InputText v-else v-model="data.label" class="w-full" size="small" />
                  </template>

                  <template #body-kind="{ data, index }">
                    <span v-if="data.isSystem" class="text-muted-color">{{ data.kind }}</span>
                    <InputGroup v-else>
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

                  <template #body-isRequired="{ data }">
                    <div v-if="!data.isSystem" class="flex justify-center">
                      <Checkbox v-model="data.isRequired" :binary="true" />
                    </div>
                    <span v-else />
                  </template>

                  <template #body-isMulti="{ data }">
                    <div v-if="!data.isSystem" class="flex justify-center">
                      <Checkbox v-model="data.isMulti" :binary="true" />
                    </div>
                    <span v-else />
                  </template>
                </CResourceTable>
              </TabPanel>

              <TabPanel value="dal">
                <DalSettings :module="module" />
              </TabPanel>

              <TabPanel value="unique">
                <UniqueValues :module="module" />
              </TabPanel>

              <TabPanel value="revisions">
                <RecordRevisionsSettings :module="module" />
              </TabPanel>

              <TabPanel value="privacy">
                <DataPrivacySettings
                  :resource="module"
                  :connection="{}"
                  :sensitivityLevels="sensitivityLevels"
                  :translations="privacyTranslations"
                />
              </TabPanel>

              <TabPanel value="issues">
                <ModuleIssues :module="module" />
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

    <!-- Config Modals -->
    <FederationSettings v-model:modal="federationModal" :module="module" />
    <DiscoverySettings v-model:modal="discoveryModal" :module="module" @save="onDiscoverySave" />
    <DalSchemaAlterations v-model:modal="schemaModal" :module="module" :batch="schemaBatch" />
  </Form>
</template>

<script setup>
import { useModuleStore } from '@/stores/module'
import { usePageStore } from '@/stores/page'
import { compose } from '@cortezaproject/corteza-js-next'
import { components, useConfirmDelete, usePermissions } from '@cortezaproject/corteza-vue-next'
import { computed, inject, onBeforeUnmount, onMounted, ref, watch, nextTick } from 'vue'
import CFieldConfigurator from '@/components/ModuleFields/Configurator/index.vue'
import DalSettings from '@/components/Admin/Module/DalSettings.vue'
import UniqueValues from '@/components/Admin/Module/UniqueValues.vue'
import RecordRevisionsSettings from '@/components/Admin/Module/RecordRevisionsSettings.vue'
import DataPrivacySettings from '@/components/Admin/Module/DataPrivacySettings.vue'
import FederationSettings from '@/components/Admin/Module/FederationSettings.vue'
import DiscoverySettings from '@/components/Admin/Module/DiscoverySettings.vue'
import DalSchemaAlterations from '@/components/Admin/Module/DalSchemaAlterations.vue'
import ModuleIssues from '@/components/Admin/Module/ModuleIssues.vue'

const { CInputDelete, CRouterLinkButton, CResourceTable } = components
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

const tableScrollHeight = ref('50vh')
let resizeObserver = null

// Configurator State
const configuratorVisible = ref(false)
const activeConfiguratorField = ref(null)
const activeConfiguratorFieldIndex = ref(-1)

// Field table ref
const fieldTableRef = ref()

// App State Defaults
const $SystemAPI = inject('$SystemAPI')
const sensitivityLevels = ref([])
const federationModal = ref(false)
const discoveryModal = ref(false)
const schemaModal = ref(false)
const schemaBatch = ref(undefined)

// Permissions State
const permissionsMenu = ref(null)
const { open: openPermissions } = usePermissions()

const togglePermissionsMenu = event => {
  permissionsMenu.value?.toggle(event)
}

const permissionsMenuItems = computed(() => {
  if (!module.value) return []
  return [
    {
      label: t('module.tooltip.permissions', 'Module Permissions'),
      command: () => {
        openPermissions({
          resource: `corteza::compose:module/${module.value.namespaceID}/${module.value.moduleID}`,
          title: module.value.name || module.value.handle || module.value.moduleID,
          target: module.value.name || module.value.handle || module.value.moduleID,
        })
      },
    },
    {
      label: t('module.fieldPermissions', 'Field Permissions'),
      command: () => {
        openPermissions({
          resource: `corteza::compose:module-field/${module.value.namespaceID}/${module.value.moduleID}/*`,
          title: module.value.name || module.value.handle || module.value.moduleID,
          target: module.value.name || module.value.handle || module.value.moduleID,
        })
      },
    },
    {
      label: t('module.recordPermissions', 'Record Permissions'),
      command: () => {
        openPermissions({
          resource: `corteza::compose:record/${module.value.namespaceID}/${module.value.moduleID}/*`,
          title: module.value.name || module.value.handle || module.value.moduleID,
          target: module.value.name || module.value.handle || module.value.moduleID,
        })
      },
    },
  ]
})

const privacyTranslations = computed(() => ({
  sensitivity: {
    label: t('module.edit.config.privacy.sensitivity-level.label', 'Sensitivity'),
    description: t('module.edit.config.privacy.sensitivity-level.description', 'Data access sensitivity'),
    placeholder: t('module.edit.config.privacy.sensitivity-level.placeholder', 'Select Sensitivity'),
  },
  usage: {
    label: t('module.edit.config.privacy.usage-disclosure.label', 'Usage Disclosure'),
  },
}))

function onDiscoverySave(mod) {
  module.value.config = mod.config
  discoveryModal.value = false
  // Flag to maybe save right after discovery save? It's fine to require top level form save.
}

async function fetchSensitivityLevels() {
  try {
    const { set } = await $SystemAPI.dalSensitivityLevelList()
    sensitivityLevels.value = set || []
  } catch (e) {}
}

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

// Field table column definitions
const fieldTableColumns = [
  { key: 'name', header: t('module.edit.fields.columns.name.label') },
  { key: 'label', header: t('module.edit.fields.columns.title.label') },
  { key: 'kind', header: t('module.edit.fields.columns.type.label') },
  {
    key: 'isRequired',
    header: t('module.edit.fields.columns.required.label'),
    headerStyle: 'width: 5rem',
    headerClass: 'text-center',
    bodyClass: 'text-center',
  },
  {
    key: 'isMulti',
    header: t('module.edit.fields.columns.multi.label'),
    headerStyle: 'width: 5rem',
    headerClass: 'text-center',
    bodyClass: 'text-center',
  },
]

// Stable key counter for new fields (fieldID is '0' for unsaved fields, so not usable as key)
let _fieldKeyCounter = 0

function ensureFieldKey(field) {
  if (!field._dataKey) {
    field._dataKey =
      field.fieldID && field.fieldID !== '0' ? field.fieldID : `new_${++_fieldKeyCounter}`
  }
}

// Fields with stable _dataKey so DataTable doesn't re-mount rows on name changes
const fieldsForTable = computed(() => {
  if (!module.value?.fields) return []
  module.value.fields.forEach(ensureFieldKey)
  return module.value.fields
})

// Combined regular + system fields for a single table
const allFieldsForTable = computed(() => {
  const regular = fieldsForTable.value
  if (!module.value) return regular
  const system = module.value.systemFields().map(f => ({
    ...f,
    _dataKey: `sys_${f.name}`,
  }))
  return [...regular, ...system]
})

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
  const field = new compose.ModuleFieldString()
  ensureFieldKey(field)
  module.value.fields.push(field)
}

function removeField(index) {
  module.value.fields.splice(index, 1)
}

function getFieldActionsMenuItems(field, index) {
  if (field.isSystem) return []

  const items = []

  if (isEdit.value && field.fieldID && field.fieldID !== '0' && module.value?.canGrant) {
    items.push({
      label: t('general.label.permissions'),
      icon: 'pi pi-lock',
      command: () => {
        openPermissions({
          resource: `corteza::compose:module-field/${module.value.namespaceID}/${module.value.moduleID}/${field.fieldID}`,
          title: field.label || field.name || field.fieldID,
        })
      },
    })
  }

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

async function checkSchemaAlterations() {
  try {
    const { set } = await $SystemAPI.dalSchemaAlterationList({
      resourceType: 'compose:module',
      resourceID: module.value.moduleID,
      completedAt: null,
    })
    schemaBatch.value = (set || []).filter(s => !!s.batchID).map(s => s.batchID)
    schemaModal.value = false
    setTimeout(() => {
      schemaModal.value = true
    }, 10)
  } catch (e) {
    if ($toast && $toast.toastErrorHandler)
      $toast.toastErrorHandler(t('module.edit.schemaAlterations.notification.load.error'))(e)
  }
}

async function handleRecordPageCreation() {
  creatingRecordPage.value = true
  try {
    const { name, moduleID } = module.value
    const { namespaceID } = props.namespace

    // Create a simple record page with a Record block
    const blocks = [new compose.PageBlockRecord({ xywh: [0, 0, 48, 36] })]
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
        xywh: [0, 0, 48, 36],
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
function updateTableScrollHeight() {
  nextTick(() => {
    const el = fieldTableRef.value?.dataTableRef?.$el
    if (!el) return

    // Find the header row inside DataTable to measure its height
    const header = el.querySelector('.p-datatable-header-cell')?.closest('thead')
    const headerHeight = header?.offsetHeight || 40

    const rect = el.getBoundingClientRect()
    const footerOffset = 70 // footer bar height + padding
    const remaining = window.innerHeight - rect.top - headerHeight - footerOffset
    const minHeight = window.innerHeight * 0.5 // 50vh

    tableScrollHeight.value = `${Math.max(remaining, minHeight)}px`
  })
}

onMounted(() => {
  loadModule()
  fetchSensitivityLevels()

  // Observe layout changes to recalculate scroll height
  updateTableScrollHeight()
  window.addEventListener('resize', updateTableScrollHeight)

  resizeObserver = new ResizeObserver(updateTableScrollHeight)
  const formEl = document.querySelector('form')
  if (formEl) resizeObserver.observe(formEl)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateTableScrollHeight)
  resizeObserver?.disconnect()
})

watch(
  () => route.params.moduleID,
  () => {
    loadModule()
  },
)

// Recalculate scroll height when loading finishes and DataTable renders
watch(loading, val => {
  if (!val) updateTableScrollHeight()
})

function exportModule() {
  if (!module.value) return
  const blob = new Blob([JSON.stringify({ type: 'module', list: [module.value] }, null, 2)], {
    type: 'application/json',
  })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${module.value.handle || module.value.name || 'module'}-export.json`
  a.click()
  URL.revokeObjectURL(url)
}
</script>
