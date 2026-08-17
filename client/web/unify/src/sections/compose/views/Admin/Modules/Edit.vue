<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <ButtonGroup v-if="isEdit" class="gap-1">
      <CRouterLinkButton
        :to="{ name: 'admin.modules.record.list', params: { moduleID: module?.moduleID } }"
        :label="$t('module.allRecords.label')"
        icon="pi pi-table"
        size="small"
      />
      <ModuleTranslator
        v-if="module"
        :module="module"
        :namespace="namespace"
        @update:module="module = $event"
      />
    </ButtonGroup>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <!-- Form -->
  <Form
    v-else-if="module"
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
          v-if="discoveryEnabled"
          :label="$t('module.edit.discoverySettings.title')"
          icon="pi pi-globe"
          size="small"
          severity="secondary"
          outlined
          @click="discoveryModal = true"
        />
        <Button
          v-if="federationEnabled"
          :label="$t('module.edit.federationSettings.title')"
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
              <Tab value="dal">{{ $t('module.edit.config.dal.title') }}</Tab>
              <Tab value="unique">
                {{ $t('module.edit.config.uniqueValues.title') }}
              </Tab>
              <Tab value="revisions">
                {{ $t('module.edit.config.record-revisions.title') }}
              </Tab>
              <Tab v-if="hasIssues" value="issues" @click="onIssuesTabClick">
                <span class="text-red-500">
                  {{ $t('module.edit.issues.label', { count: module.issues.length }) }}
                </span>
              </Tab>
            </TabList>

            <TabPanels>
              <TabPanel value="fields">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                  <CFormGroup name="name" :label="$t('module.general.label.name')" required>
                    <!-- No :invalid — the resolver's required-only rule is the whole
                         rule, and the form drives the state from it. -->
                    <InputText id="name" name="name" v-model="module.name" />
                  </CFormGroup>

                  <CFormGroup name="handle" :label="$t('module.general.label.handle')">
                    <InputText
                      id="handle"
                      name="handle"
                      v-model="module.handle"
                      :invalid="!!module.handle && !isValidHandle(module.handle)"
                    />
                  </CFormGroup>
                </div>

                <Divider />

                <!-- Module Fields -->
                <div class="flex items-center mb-4">
                  <Button
                    :label="$t('module.edit.newField')"
                    icon="pi pi-plus"
                    severity="secondary"
                    size="small"
                    @click="addField"
                  />
                </div>

                <CFormList
                  v-model="module.fields"
                  draggable
                  :empty-message="$t('module.edit.fields.empty')"
                  :columns="fieldFormListColumns"
                >
                  <template #row="{ item: field, index }">
                    <div class="flex flex-col gap-1">
                      <InputText
                        v-model="field.name"
                        class="w-full"
                        size="small"
                        :invalid="validationTriggered && fieldNameError(field) !== ''"
                      />
                      <Message
                        v-if="validationTriggered && fieldNameError(field)"
                        severity="error"
                        size="small"
                        variant="simple"
                      >
                        {{ fieldNameError(field) }}
                      </Message>
                    </div>

                    <div class="flex flex-col gap-1">
                      <InputGroup>
                        <InputText
                          v-model="field.label"
                          class="w-full"
                          size="small"
                          :invalid="validationTriggered && (!field.label || !field.label.trim())"
                        />
                        <InputGroupAddon
                          v-if="
                            showTranslatorButton && isEdit && field.fieldID && field.fieldID !== '0'
                          "
                        >
                          <Button
                            icon="pi pi-language"
                            severity="secondary"
                            size="small"
                            class="w-full border-none"
                            v-tooltip.top="$t('field.translate.label')"
                            @click="openFieldTranslation(field)"
                          />
                        </InputGroupAddon>
                      </InputGroup>
                      <Message
                        v-if="validationTriggered && (!field.label || !field.label.trim())"
                        severity="error"
                        size="small"
                        variant="simple"
                      >
                        {{ $t('general.label.required') }}
                      </Message>
                    </div>

                    <InputGroup>
                      <Select
                        v-model="field.kind"
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
                          @click="openFieldConfigurator(field, index)"
                        />
                      </InputGroupAddon>
                    </InputGroup>

                    <div class="flex justify-center">
                      <Checkbox v-model="field.isRequired" :binary="true" />
                    </div>

                    <div class="flex justify-center">
                      <Checkbox v-model="field.isMulti" :binary="true" />
                    </div>

                    <div class="flex items-center justify-end gap-1">
                      <Button
                        v-if="canGrantField(field)"
                        icon="pi pi-lock"
                        text
                        severity="secondary"
                        size="small"
                        v-tooltip.top="$t('general.label.permissions')"
                        @click="openFieldPermissions(field)"
                      />
                      <Button
                        v-if="fieldActionsMenuItems(field, index).length"
                        icon="pi pi-ellipsis-v"
                        text
                        severity="secondary"
                        size="small"
                        @click="showFieldActionsMenu($event, field, index)"
                      />
                    </div>
                  </template>

                  <template #footer="{ gridStyle }">
                    <div
                      v-for="f in systemFieldsForDisplay"
                      :key="f.name"
                      class="border-t border-surface p-3 bg-highlight text-muted-color"
                      v-tooltip.left="$t('module.edit.systemField')"
                    >
                      <div :style="gridStyle" class="grid gap-2 items-center">
                        <i class="pi pi-lock text-muted-color text-center" />
                        <InputText :model-value="f.name" class="w-full" size="small" disabled />
                        <InputText :model-value="f.label" class="w-full" size="small" disabled />
                        <InputText :model-value="f.kind" class="w-full" size="small" disabled />
                        <div class="flex justify-center">
                          <i
                            :class="[
                              'pi',
                              f.isRequired ? 'pi-check text-primary' : 'pi-minus text-muted-color',
                            ]"
                          />
                        </div>
                        <div class="flex justify-center">
                          <i
                            :class="[
                              'pi',
                              f.isMulti ? 'pi-check text-primary' : 'pi-minus text-muted-color',
                            ]"
                          />
                        </div>
                        <span />
                        <span />
                      </div>
                    </div>
                  </template>
                </CFormList>

                <!-- One popup shared by every field row's actions button -->
                <Menu ref="fieldActionsMenuRef" :model="currentFieldMenuItems" popup />
              </TabPanel>

              <TabPanel value="dal">
                <DalSettings />
              </TabPanel>

              <TabPanel value="unique">
                <UniqueValues />
              </TabPanel>

              <TabPanel value="revisions">
                <RecordRevisionsSettings />
              </TabPanel>

              <TabPanel value="issues">
                <ModuleIssues :module="module" />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </template>
      </Card>
    </div>

    <CEditorActions
      :back-to="true"
      @back="goBack({ name: 'admin.modules', params: { slug: route.params.slug } })"
    >
      <CInputDelete
        v-if="isEdit && module.canDeleteModule"
        :label="$t('general.label.delete')"
        :message="$t('module.edit.deleteConfirm')"
        :header="module.name"
        :disabled="deleting"
        @confirm="handleDelete"
      />
      <Button
        v-if="isEdit"
        :label="$t('general.label.saveAsCopy')"
        icon="pi pi-copy"
        severity="secondary"
        :loading="cloning"
        :disabled="!canSave"
        @click="handleClone"
      />
      <Button
        type="submit"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
        :disabled="!canSave"
      />
    </CEditorActions>

    <!-- Field Configurator Modal -->
    <CFieldConfigurator
      v-model:visible="configuratorVisible"
      :field="activeConfiguratorField"
      :namespace="namespace"
      :module="module"
      @save="onFieldSave"
    />

    <!-- Config Modals -->
    <FederationSettings v-if="federationEnabled" v-model:modal="federationModal" :module="module" />
    <DiscoverySettings
      v-if="discoveryEnabled"
      v-model:modal="discoveryModal"
      :module="module"
      @save="onDiscoverySave"
    />
    <DalSchemaAlterations v-model:modal="schemaModal" :module="module" :batch="schemaBatch" />
  </Form>
</template>

<script setup>
import { useModuleStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { compose } from '@planetcrust/human-js'
import {
  components,
  useConfirmDelete,
  useHistoryBack,
  usePermissions,
  useDraftGuard,
} from '@planetcrust/human-vue'
import { computed, inject, onMounted, provide, ref, watch, nextTick } from 'vue'
import CFieldConfigurator from '@/sections/compose/components/ModuleFields/Configurator/index.vue'
import DalSettings from '@/sections/compose/components/Admin/Module/DalSettings.vue'
import UniqueValues from '@/sections/compose/components/Admin/Module/UniqueValues.vue'
import RecordRevisionsSettings from '@/sections/compose/components/Admin/Module/RecordRevisionsSettings.vue'
import FederationSettings from '@/sections/compose/components/Admin/Module/FederationSettings.vue'
import DiscoverySettings from '@/sections/compose/components/Admin/Module/DiscoverySettings.vue'
import DalSchemaAlterations from '@/sections/compose/components/Admin/Module/DalSchemaAlterations.vue'
import ModuleIssues from '@/sections/compose/components/Admin/Module/ModuleIssues.vue'
import ModuleTranslator from '@/sections/compose/components/Admin/Module/ModuleTranslator.vue'
import { useTranslatorStore } from '@/sections/compose/stores/translator'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'
import {
  applyFieldTranslations,
  applySelectTranslations,
  applyBoolTranslations,
} from '@/sections/compose/lib/resource-translations'

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
const goBack = useHistoryBack()
const { t } = useI18n()
const { confirmDelete } = useConfirmDelete()
const $toast = inject('$toast')
const moduleStore = useModuleStore()
const pageStore = usePageStore()

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const cloning = ref(false)
const module = ref(null)
const { capture, markSaved } = useDraftGuard({
  draft: module,
  busy: () => saving.value || deleting.value || cloning.value,
})

provide('moduleDraft', module)
const activeTab = ref('fields')

const creatingRecordPage = ref(false)
const creatingRecordListPage = ref(false)

// Configurator State
const configuratorVisible = ref(false)
const activeConfiguratorField = ref(null)
const activeConfiguratorFieldIndex = ref(-1)

// Field row action menu
const fieldActionsMenuRef = ref(null)
const currentFieldMenuItems = ref([])

// App State Defaults
const $SystemAPI = inject('$SystemAPI')
const $Settings = inject('$Settings')
const federationEnabled = computed(() => $Settings?.get('federation.enabled', false))
const discoveryEnabled = computed(() => $Settings?.get('discovery.enabled', false))
const federationModal = ref(false)
const discoveryModal = ref(false)
const schemaModal = ref(false)
const schemaBatch = ref(undefined)

// Permissions State
const permissionsMenu = ref(null)
const { open: openPermissions } = usePermissions()
const $ComposeAPI = inject('$ComposeAPI')
const translatorStore = useTranslatorStore()
const { showTranslatorButton, currentLanguage } = useResourceTranslations()

const togglePermissionsMenu = event => {
  permissionsMenu.value?.toggle(event)
}

const permissionsMenuItems = computed(() => {
  if (!module.value) return []
  return [
    {
      label: t('module.tooltip.permissions'),
      command: () => {
        openPermissions({
          resource: `corteza::compose:module/${module.value.namespaceID}/${module.value.moduleID}`,
          title: module.value.name || module.value.handle || module.value.moduleID,
          target: module.value.name || module.value.handle || module.value.moduleID,
        })
      },
    },
    {
      label: t('module.fieldPermissions'),
      command: () => {
        openPermissions({
          resource: `corteza::compose:module-field/${module.value.namespaceID}/${module.value.moduleID}/*`,
          title: module.value.name || module.value.handle || module.value.moduleID,
          target: module.value.name || module.value.handle || module.value.moduleID,
          allSpecific: true,
        })
      },
    },
    {
      label: t('module.recordPermissions'),
      command: () => {
        openPermissions({
          resource: `corteza::compose:record/${module.value.namespaceID}/${module.value.moduleID}/*`,
          title: module.value.name || module.value.handle || module.value.moduleID,
          target: module.value.name || module.value.handle || module.value.moduleID,
          allSpecific: true,
        })
      },
    },
  ]
})

function onDiscoverySave(mod) {
  module.value.config = mod.config
  discoveryModal.value = false
  // Flag to maybe save right after discovery save? It's fine to require top level form save.
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

// Field list column definitions (CFormList)
const fieldFormListColumns = computed(() => [
  {
    label: t('module.edit.fields.columns.name.label'),
    tooltip: t('module.edit.tooltip.name'),
    width: 'minmax(180px, 1.2fr)',
  },
  {
    label: t('module.edit.fields.columns.title.label'),
    tooltip: t('module.edit.tooltip.title'),
    width: 'minmax(180px, 1.2fr)',
  },
  {
    label: t('module.edit.fields.columns.type.label'),
    width: 'minmax(200px, 1.4fr)',
  },
  {
    label: t('module.edit.fields.columns.required.label'),
    width: '6rem',
    headerClass: 'text-center',
  },
  {
    label: t('module.edit.fields.columns.multi.label'),
    width: '6rem',
    headerClass: 'text-center',
  },
  { width: '5rem' },
])

const systemFieldsForDisplay = computed(() => {
  if (!module.value) return []
  const systemFieldEncoding = module.value.config?.dal?.systemFieldEncoding || {}
  return (module.value.systemFields() || []).map(sf => ({
    name: sf.name,
    label: sf.label || sf.name,
    kind: sf.kind,
    isRequired: !!sf.isRequired,
    isMulti: !!sf.isMulti,
    ...(systemFieldEncoding[sf.name] || {}),
  }))
})

// Computed
const isEdit = computed(() => !!route.params.moduleID)

const hasIssues = computed(() => {
  return (module.value?.issues || []).length > 0
})

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

// Valid field name: starts with a letter, then letters/numbers/underscores.
// Field names are identifiers (the server runs handle.IsValid over them);
// module names are labels and are not checked against this.
function isValidFieldName(name) {
  return /^[A-Za-z][A-Za-z0-9_]*$/.test(name)
}

// Valid handle: starts with letter, letters/numbers/underscores/dashes/dots, ends with letter/number
function isValidHandle(handle) {
  return /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$|^[A-Za-z]$/.test(handle)
}

const duplicateFieldNames = computed(() => {
  const seen = new Set()
  const dups = new Set()
  for (const f of module.value?.fields || []) {
    if (!f.name) continue
    if (seen.has(f.name)) dups.add(f.name)
    seen.add(f.name)
  }
  return dups
})

function isFieldNameDuplicate(name) {
  return !!name && duplicateFieldNames.value.has(name)
}

function fieldNameError(field) {
  if (!field || field.isSystem) return ''
  if (!field.name) return t('general.label.required')
  if (!isValidFieldName(field.name)) return t('module.edit.tooltip.name')
  if (isFieldNameDuplicate(field.name)) return t('module.edit.fields.duplicateName')
  return ''
}

const fieldsValid = computed(() => {
  const fields = module.value?.fields || []
  return fields.every(f => {
    if (!f.name || !isValidFieldName(f.name)) return false
    if (isFieldNameDuplicate(f.name)) return false
    if (!f.label || !f.label.trim()) return false
    return true
  })
})

const canSave = computed(() => {
  if (isEdit.value && !module.value?.canUpdateModule) return false
  return true
})

const validationTriggered = ref(false)

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
    capture()
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
    capture()

    // Auto-trigger schema alterations check if module has issues (matching Human behavior)
    if ((module.value.issues || []).length > 0) {
      checkSchemaAlterations()
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

function onIssuesTabClick() {
  checkSchemaAlterations()
}

function openFieldTranslation(field) {
  const nsID = props.namespace.namespaceID
  const modID = module.value.moduleID
  const fieldRes = `compose:module-field/${nsID}/${modID}/${field.fieldID}`
  translatorStore.open({
    resource: fieldRes,
    titles: {
      [fieldRes]: t('translator.resources.module.field.title', { name: field.label || field.name }),
    },
    fetcher: () =>
      $ComposeAPI
        .moduleListTranslations({ namespaceID: nsID, moduleID: modID })
        .then(set =>
          set
            .filter(tr => tr.resource === fieldRes)
            .filter(tr => !tr.key.startsWith('meta.options') && !tr.key.startsWith('meta.bool')),
        ),
    updater: async changes => {
      await $ComposeAPI.moduleUpdateTranslations({
        namespaceID: nsID,
        moduleID: modID,
        translations: changes,
      })
      const fresh = await $ComposeAPI.moduleListTranslations({ namespaceID: nsID, moduleID: modID })
      applyFieldTranslations(field, fresh, currentLanguage.value, fieldRes)
    },
  })
}

function canGrantField(field) {
  return !!(
    field &&
    !field.isSystem &&
    isEdit.value &&
    field.fieldID &&
    field.fieldID !== '0' &&
    module.value?.canGrant
  )
}

function openFieldPermissions(field) {
  openPermissions({
    resource: `corteza::compose:module-field/${module.value.namespaceID}/${module.value.moduleID}/${field.fieldID}`,
    title: field.label || field.name || field.fieldID,
  })
}

function fieldActionsMenuItems(field, index) {
  if (!field || field.isSystem) return []

  const items = []

  if (showTranslatorButton.value && isEdit.value && field.fieldID && field.fieldID !== '0') {
    const nsID = props.namespace.namespaceID
    const modID = module.value.moduleID
    const fieldRes = `compose:module-field/${nsID}/${modID}/${field.fieldID}`

    if (field.kind === 'Select') {
      items.push({
        label: t('field.translate.selectOptions'),
        icon: 'pi pi-language',
        command: () => {
          translatorStore.open({
            resource: fieldRes,
            titles: {
              [fieldRes]: t('translator.resources.module.field.selectOptions', {
                name: field.label || field.name,
              }),
            },
            keyPrettifier: key => {
              const match = key.match(/^meta\.options\.(.+)\.text$/)
              return match ? match[1] : key
            },
            fetcher: () =>
              $ComposeAPI
                .moduleListTranslations({ namespaceID: nsID, moduleID: modID })
                .then(set =>
                  set
                    .filter(tr => tr.resource === fieldRes)
                    .filter(tr => tr.key.startsWith('meta.options') && tr.key.endsWith('.text')),
                ),
            updater: async changes => {
              await $ComposeAPI.moduleUpdateTranslations({
                namespaceID: nsID,
                moduleID: modID,
                translations: changes,
              })
              const fresh = await $ComposeAPI.moduleListTranslations({
                namespaceID: nsID,
                moduleID: modID,
              })
              applySelectTranslations(field, fresh, currentLanguage.value, fieldRes)
            },
          })
        },
      })
    }

    if (field.kind === 'Bool') {
      items.push({
        label: t('field.translate.boolLabels'),
        icon: 'pi pi-language',
        command: () => {
          translatorStore.open({
            resource: fieldRes,
            titles: {
              [fieldRes]: t('translator.resources.module.field.boolLabels', {
                name: field.label || field.name,
              }),
            },
            fetcher: () =>
              $ComposeAPI
                .moduleListTranslations({ namespaceID: nsID, moduleID: modID })
                .then(set =>
                  set
                    .filter(tr => tr.resource === fieldRes)
                    .filter(tr => tr.key.startsWith('meta.bool') && tr.key.endsWith('.label')),
                ),
            updater: async changes => {
              await $ComposeAPI.moduleUpdateTranslations({
                namespaceID: nsID,
                moduleID: modID,
                translations: changes,
              })
              const fresh = await $ComposeAPI.moduleListTranslations({
                namespaceID: nsID,
                moduleID: modID,
              })
              applyBoolTranslations(field, fresh, currentLanguage.value, fieldRes)
            },
          })
        },
      })
    }
  }

  return items
}

// show() rather than toggle(): one popup serves every field row, so clicking a
// second row's button must re-anchor and open there rather than close the first.
function showFieldActionsMenu(event, field, index) {
  currentFieldMenuItems.value = fieldActionsMenuItems(field, index)
  fieldActionsMenuRef.value?.show(event, event.currentTarget)
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
  if (!valid || !fieldsValid.value) {
    validationTriggered.value = true
    activeTab.value = 'fields'
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error, .p-inputtext.p-invalid')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }
  if (!canSave.value) return
  validationTriggered.value = false

  saving.value = true
  try {
    const payload = {
      namespaceID: props.namespace.namespaceID,
      name: module.value.name,
      handle: module.value.handle,
      fields: module.value.fields,
      config: module.value.config,
      meta: module.value.meta,
    }

    if (isEdit.value) {
      payload.moduleID = module.value.moduleID
      const updated = await moduleStore.update(payload)
      module.value = new compose.Module({ ...updated })
      capture()
      $toast.toastSuccess(t('notification.module.saved'))
    } else {
      const created = await moduleStore.create(payload)
      $toast.toastSuccess(t('notification.module.created'))
      markSaved()
      router.push({
        name: 'admin.modules.edit',
        params: { moduleID: created.moduleID },
      })
    }
  } catch (e) {
    console.error('Failed to save module:', e)
    $toast.toastErrorHandler(t('notification.module.saveFailed'))(e)
  } finally {
    saving.value = false
  }
}

async function handleClone() {
  cloning.value = true
  try {
    const payload = {
      namespaceID: props.namespace.namespaceID,
      name: `${module.value.name} (${t('general.label.clone').toLowerCase()})`,
      handle: '',
      fields: module.value.fields,
      config: module.value.config,
      meta: module.value.meta,
    }

    const created = await moduleStore.create(payload)
    $toast.toastSuccess(t('notification.module.created'))
    router.push({
      name: 'admin.modules.edit',
      params: { moduleID: created.moduleID },
    })
  } catch (e) {
    console.error('Failed to clone module:', e)
    $toast.toastErrorHandler(t('notification.module.createFailed'))(e)
  } finally {
    cloning.value = false
  }
}

async function handleDelete() {
  deleting.value = true
  try {
    // namespaceID is required: the API client refuses the call before it
    // reaches the network without it, which is how deleting from this view
    // failed with "field namespaceID is empty" while the module list — which
    // passes both — worked fine.
    await moduleStore.delete({
      namespaceID: module.value.namespaceID,
      moduleID: module.value.moduleID,
    })
    capture()
    $toast.toastSuccess(t('notification.module.deleted'))
    router.push({ name: 'admin.modules' })
  } catch (e) {
    console.error('Failed to delete module:', e)
    $toast.toastErrorHandler(t('notification.module.deleteFailed'))(e)
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
    $toast.toastErrorHandler(t('notification.page.createFailed'))(e)
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
    $toast.toastErrorHandler(t('notification.page.createFailed'))(e)
  } finally {
    creatingRecordListPage.value = false
  }
}

onMounted(() => {
  loadModule()
})

watch(
  () => route.params.moduleID,
  () => {
    loadModule()
  },
)

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
