<template>
  <Teleport to="#topbar-title" defer>
    <span>{{ pageTitle }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <ButtonGroup v-if="isEdit" class="gap-1">
      <Button
        v-if="page?.isRecordPage"
        :label="$t('page.moduleEdit')"
        icon="pi pi-database"
        size="small"
        @click="goToModuleEdit"
      />
      <Button
        :label="$t('page.edit.pageBuilder')"
        icon="pi pi-wrench"
        size="small"
        @click="goToBuilder"
      />
      <Button
        :label="$t('page.edit.viewPage')"
        icon="pi pi-eye"
        size="small"
        @click="goToViewPage"
      />
    </ButtonGroup>
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
    <div class="container mx-auto p-4 flex-1 overflow-auto flex flex-col gap-4">
      <!-- General Panel -->
      <Panel :header="$t('general.label.general')" toggleable>
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

        <div class="flex flex-col gap-2 mb-6">
          <label for="description" class="font-medium text-primary">
            {{ $t('page.label.description') }}
          </label>
          <Textarea id="description" v-model="page.description" rows="4" auto-resize />
        </div>

        <!-- Page Icon + Other Options (side by side like Corteza) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <!-- Page Icon -->
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-2">
              <label class="font-medium text-primary">{{ $t('page.icon.page') }}</label>
              <Button
                v-tooltip.top="$t('page.icon.configure')"
                icon="pi pi-pencil"
                text
                severity="secondary"
                size="small"
                @click="openIconModal"
              />
            </div>

            <img
              v-if="pageIconSrc"
              :src="pageIconSrc"
              width="auto"
              height="50"
            >
            <span v-else class="text-muted-color">
              {{ $t('page.icon.noIcon') }}
            </span>
          </div>

          <!-- Other Options -->
          <div class="flex flex-col gap-2">
            <label class="font-medium text-primary">{{ $t('page.edit.otherOptions') }}</label>

            <div v-if="!isRecordPage" class="flex items-center gap-3">
              <Checkbox v-model="page.visible" :binary="true" input-id="visible" />
              <label for="visible" class="cursor-pointer">
                {{ $t('page.edit.visible') }}
              </label>
            </div>

            <div class="flex items-center gap-3">
              <Checkbox v-model="showSubPages" :binary="true" input-id="showSubPages" />
              <label for="showSubPages" class="cursor-pointer">
                {{ $t('page.showSubPages') }}
              </label>
            </div>

            <div v-if="isRecordPage" class="flex items-center gap-3">
              <Checkbox v-model="notificationsEnabled" :binary="true" input-id="notifications" />
              <label for="notifications" class="cursor-pointer">
                {{ $t('page.edit.notifications.enabled') }}
              </label>
            </div>
          </div>
        </div>
      </Panel>

      <!-- Layouts Panel -->
      <Panel v-if="isEdit" :header="$t('page.page-layout.layouts')" toggleable>
        <div class="flex items-center justify-end mb-4">
          <Button
            :label="$t('page.page-layout.add')"
            icon="pi pi-plus"
            size="small"
            @click="addLayout"
          />
        </div>

        <DataTable
          v-if="layouts.length > 0"
          :value="layouts"
          striped-rows
          data-key="_key"
          class="border border-b-0 border-surface rounded-border overflow-auto"
          :pt="{ headerCell: { class: 'bg-highlight-emphasis' } }"
        >
          <Column header-style="width: 5rem">
            <template #body="{ index }">
              <div class="flex gap-1">
                <Button
                  icon="pi pi-arrow-up"
                  text
                  severity="secondary"
                  size="small"
                  :disabled="index === 0"
                  @click="moveLayout(index, -1)"
                />
                <Button
                  icon="pi pi-arrow-down"
                  text
                  severity="secondary"
                  size="small"
                  :disabled="index === layouts.length - 1"
                  @click="moveLayout(index, 1)"
                />
              </div>
            </template>
          </Column>

          <Column :header="$t('page.page-layout.title')" style="min-width: 250px">
            <template #body="{ data }">
              <InputText
                v-model="data.meta.title"
                class="w-full"
                size="small"
                @input="data._updated = true"
              />
            </template>
          </Column>

          <Column :header="$t('page.page-layout.handle')" style="min-width: 250px">
            <template #body="{ data }">
              <InputGroup>
                <InputText
                  v-model="data.handle"
                  class="w-full"
                  size="small"
                  @input="data._updated = true"
                />
                <InputGroupAddon>
                  <Button
                    v-tooltip.top="$t('page.page-layout.tooltip.configure')"
                    icon="pi pi-cog"
                    severity="secondary"
                    size="small"
                    class="w-full border-none"
                    @click="openLayoutConfig(data)"
                  />
                </InputGroupAddon>
                <InputGroupAddon>
                  <Button
                    v-tooltip.top="$t('page.page-layout.tooltip.builder')"
                    icon="pi pi-wrench"
                    severity="secondary"
                    size="small"
                    class="w-full border-none"
                    :disabled="data.pageLayoutID === NoID"
                    @click="goToLayoutBuilder(data)"
                  />
                </InputGroupAddon>
              </InputGroup>
            </template>
          </Column>

          <Column header-style="width: 4rem">
            <template #body="{ data, index }">
              <CInputDelete
                text
                size="small"
                :message="$t('page.edit.deleteConfirm')"
                :header="data.meta.title || data.handle"
                @confirm="removeLayout(index)"
              />
            </template>
          </Column>
        </DataTable>

        <div v-else class="text-center py-4 text-muted-color border border-surface rounded-border">
          {{ $t('page.noBlock') }}
        </div>
      </Panel>
    </div>

    <!-- Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.push({ name: 'admin.pages' })"
        />
        <div class="flex gap-2">
          <!-- Delete with strategy for pages with children -->
          <template v-if="isEdit && page.canDeletePage">
            <template v-if="hasChildren">
              <Button
                :label="$t('general.label.delete')"
                icon="pi pi-trash"
                severity="danger"
                outlined
                size="small"
                @click="toggleDeleteMenu"
              />
              <TieredMenu ref="deleteMenu" :model="deleteMenuItems" popup />
            </template>

            <CInputDelete
              v-else
              :label="$t('general.label.delete')"
              :message="$t('page.edit.deleteConfirm')"
              :header="page.title"
              :disabled="deleting"
              @confirm="handleDelete('abort')"
            />
          </template>

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

  <!-- Layout Configuration Dialog -->
  <Dialog
    v-model:visible="layoutConfigVisible"
    :header="layoutConfigTitle"
    modal
    :style="{ width: '50vw' }"
    :breakpoints="{ '768px': '90vw' }"
    :closable="true"
    @hide="onLayoutConfigClose"
  >
    <template v-if="configLayout">
      <!-- General -->
      <h5 class="font-semibold mb-3">{{ $t('page.page-layout.general') }}</h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary">{{ $t('page.page-layout.title') }}</label>
          <InputText v-model="configLayout.meta.title" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-primary">{{ $t('page.page-layout.handle') }}</label>
          <InputText v-model="configLayout.handle" />
        </div>
      </div>

      <!-- Use Title (record pages only) -->
      <div v-if="isRecordPage" class="flex items-center gap-3 mb-4">
        <ToggleSwitch v-model="configLayout.config.useTitle" />
        <label class="font-medium text-primary">{{ $t('page.page-layout.useTitle') }}</label>
      </div>

      <Divider />

      <!-- Visibility -->
      <h5 class="font-semibold mb-3">{{ $t('page.page-layout.visibility') }}</h5>

      <div class="flex flex-col gap-2 mb-4">
        <label class="font-medium text-primary">{{ $t('page.page-layout.condition.label') }}</label>
        <InputGroup>
          <InputGroupAddon>ƒ</InputGroupAddon>
          <InputText
            v-model="configLayout.config.visibility.expression"
            :placeholder="$t('page.page-layout.condition.placeholder')"
          />
        </InputGroup>
        <small class="text-muted-color" v-if="isRecordPage">
          {{ $t('page.page-layout.condition.description.record-page', {
            0: 'record.values.fieldName',
            1: 'user.(userID/email...)',
            2: 'screen.(width/height)',
            3: 'isView/isCreate/isEdit',
            4: 'user.userID == record.createdBy',
            5: 'screen.width < 1024',
          }) }}
        </small>
        <small class="text-muted-color" v-else>
          {{ $t('page.page-layout.condition.description.non-record-page', {
            0: 'user.(userID/email...)',
            1: 'screen.(width/height)',
            2: 'user.email == "test@mail.com"',
            3: 'screen.width < 1024',
          }) }}
        </small>
      </div>

      <div class="flex flex-col gap-2 mb-4">
        <label class="font-medium text-primary">{{ $t('page.page-layout.roles.label') }}</label>
        <CInputRole
          :value="configLayoutRoles"
          :placeholder="$t('page.page-layout.roles.placeholder')"
          multiple
          @input="onConfigLayoutRoleChange"
        />
      </div>

      <!-- Record Toolbar (record pages only) -->
      <template v-if="isRecordPage">
        <Divider />

        <h5 class="font-semibold mb-3">{{ $t('page.page-layout.recordToolbar.label') }}</h5>

        <div class="flex flex-col gap-2 mb-4">
          <label class="font-medium text-primary">{{ $t('page.page-layout.recordToolbar.buttons.label') }}</label>

          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-3">
              <Checkbox v-model="configLayout.config.buttons.back.enabled" :binary="true" input-id="btn-back" />
              <label for="btn-back">{{ $t('page.page-layout.recordToolbar.buttons.showBack') }}</label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox v-model="configLayout.config.buttons.delete.enabled" :binary="true" input-id="btn-delete" />
              <label for="btn-delete">{{ $t('page.page-layout.recordToolbar.buttons.showDelete') }}</label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox v-model="configLayout.config.buttons.clone.enabled" :binary="true" input-id="btn-clone" />
              <label for="btn-clone">{{ $t('page.page-layout.recordToolbar.buttons.showClone') }}</label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox v-model="configLayout.config.buttons.new.enabled" :binary="true" input-id="btn-new" />
              <label for="btn-new">{{ $t('page.page-layout.recordToolbar.buttons.showNew') }}</label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox v-model="configLayout.config.buttons.edit.enabled" :binary="true" input-id="btn-edit" />
              <label for="btn-edit">{{ $t('page.page-layout.recordToolbar.buttons.showEdit') }}</label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox v-model="configLayout.config.buttons.submit.enabled" :binary="true" input-id="btn-submit" />
              <label for="btn-submit">{{ $t('page.page-layout.recordToolbar.buttons.showSave') }}</label>
            </div>
          </div>
        </div>

        <Divider />

        <!-- Custom Actions -->
        <div class="flex items-center justify-between mb-3">
          <h5 class="font-semibold">{{ $t('page.page-layout.recordToolbar.actions.label') }}</h5>
          <Button
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            size="small"
            @click="addLayoutAction"
          />
        </div>

        <DataTable
          v-if="configLayout.config.actions.length > 0"
          :value="configLayout.config.actions"
          striped-rows
          class="border border-surface rounded-border mb-4"
        >
          <Column :header="$t('page.page-layout.recordToolbar.actions.buttonLabel')" style="min-width: 200px">
            <template #body="{ data }">
              <InputText v-model="data.meta.label" class="w-full" size="small" />
            </template>
          </Column>

          <Column :header="$t('page.page-layout.recordToolbar.actions.kind.label')" style="min-width: 180px">
            <template #body="{ data }">
              <Select
                v-model="data.kind"
                :options="actionKindOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                size="small"
                @change="onActionKindChange(data)"
              />
            </template>
          </Column>

          <Column :header="$t('page.page-layout.recordToolbar.actions.variant')" style="min-width: 140px">
            <template #body="{ data }">
              <Select
                v-model="data.meta.style.variant"
                :options="actionVariantOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                size="small"
              />
            </template>
          </Column>

          <Column :header="$t('page.page-layout.recordToolbar.actions.placement.label')" style="min-width: 120px">
            <template #body="{ data }">
              <Select
                v-model="data.placement"
                :options="actionPlacementOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                size="small"
              />
            </template>
          </Column>

          <Column :header="$t('page.page-layout.recordToolbar.actions.visible')" header-style="width: 5rem" header-class="text-center" body-class="text-center">
            <template #body="{ data }">
              <Checkbox v-model="data.enabled" :binary="true" />
            </template>
          </Column>

          <Column header-style="width: 3rem">
            <template #body="{ index }">
              <Button
                icon="pi pi-trash"
                text
                severity="danger"
                size="small"
                @click="removeLayoutAction(index)"
              />
            </template>
          </Column>
        </DataTable>

        <!-- Action-specific config (shown per action in the table) -->
        <div
          v-for="(action, aIdx) in configLayout.config.actions"
          :key="aIdx"
          class="mb-2"
        >
          <div v-if="action.kind === 'toLayout'" class="flex flex-col gap-2">
            <label class="text-sm text-muted-color">
              {{ $t('page.page-layout.recordToolbar.actions.toLayout.label') }} — {{ action.meta.label || `#${aIdx + 1}` }}
            </label>
            <Select
              v-model="action.params.pageLayoutID"
              :options="actionLayoutOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              size="small"
            />
          </div>

          <div v-if="action.kind === 'toURL'" class="flex flex-col gap-2">
            <label class="text-sm text-muted-color">
              {{ $t('page.page-layout.recordToolbar.actions.toURL.label') }} — {{ action.meta.label || `#${aIdx + 1}` }}
            </label>
            <InputText
              v-model="action.params.url"
              :placeholder="$t('page.page-layout.recordToolbar.actions.toURL.placeholder')"
              size="small"
            />
            <Select
              v-model="action.params.openIn"
              :options="actionOpenInOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              size="small"
            />
          </div>
        </div>
      </template>
    </template>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        @click="layoutConfigVisible = false"
      />
      <Button
        :label="$t('general.label.save')"
        :disabled="!configLayout?.meta?.title"
        @click="saveLayoutConfig"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { usePageStore } from '@/stores/page'
import { usePageLayoutStore } from '@/stores/page-layout'
import { compose, NoID } from '@cortezaproject/corteza-js-next'
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
const $ComposeAPI = inject('$ComposeAPI')
const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()

// State
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const page = ref(null)

// Page layouts state
const layouts = ref([])
const removedLayouts = ref([])
let layoutKeyCounter = 0

// Layout config dialog state
const layoutConfigVisible = ref(false)
const configLayout = ref(null)
const configLayoutIndex = ref(-1)
const configLayoutRoles = ref([])

// Delete menu ref
const deleteMenu = ref()

// Computed
const isEdit = computed(() => !!route.params.pageID)

const isRecordPage = computed(() => {
  return page.value && page.value.moduleID && page.value.moduleID !== NoID
})

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

const hasChildren = computed(() => {
  if (!page.value) return false
  return pageStore.set.some(p => p.selfID === page.value.pageID)
})

// Page icon source
const pageIconSrc = computed(() => {
  const icon = page.value?.config?.navItem?.icon
  if (!icon?.src) return ''
  return icon.type === 'link' ? icon.src : icon.src
})

// Show sub-pages computed property with getter/setter
const showSubPages = computed({
  get () {
    return page.value?.config?.navItem?.expanded ?? false
  },
  set (val) {
    if (page.value) {
      if (!page.value.config) page.value.config = {}
      if (!page.value.config.navItem) page.value.config.navItem = {}
      page.value.config.navItem.expanded = val
    }
  },
})

// Notifications enabled computed property
const notificationsEnabled = computed({
  get () {
    return page.value?.meta?.notifications?.enabled ?? false
  },
  set (val) {
    if (page.value) {
      if (!page.value.meta) page.value.meta = {}
      if (!page.value.meta.notifications) page.value.meta.notifications = {}
      page.value.meta.notifications.enabled = val
    }
  },
})

// Delete menu items for pages with children
const deleteMenuItems = computed(() => [
  {
    label: t('page.delete.rebase'),
    icon: 'pi pi-arrow-up',
    command: () => handleDelete('rebase'),
  },
  {
    label: t('page.delete.cascade'),
    icon: 'pi pi-trash',
    command: () => handleDelete('cascade'),
  },
])

// Layout config dialog title
const layoutConfigTitle = computed(() => {
  if (configLayout.value?.meta?.title) {
    return t('page.page-layout.configure', { title: configLayout.value.meta.title })
  }
  return t('page.page-layout.configure', { title: '' })
})

// Action options for layout config
const actionKindOptions = computed(() => [
  { value: 'toLayout', label: t('page.page-layout.recordToolbar.actions.kind.toLayout') },
  { value: 'toURL', label: t('page.page-layout.recordToolbar.actions.kind.toURL') },
])

const actionLayoutOptions = computed(() => {
  const options = [{ value: '', label: t('page.page-layout.recordToolbar.actions.toLayout.placeholder') }]
  layouts.value
    .filter(l => l.pageLayoutID !== NoID)
    .forEach(l => {
      options.push({
        value: l.pageLayoutID,
        label: l.meta.title || l.handle || l.pageLayoutID,
      })
    })
  return options
})

const actionOpenInOptions = computed(() => [
  { value: 'sameTab', label: t('page.page-layout.recordToolbar.actions.openIn.sameTab') },
  { value: 'newTab', label: t('page.page-layout.recordToolbar.actions.openIn.newTab') },
])

const actionVariantOptions = computed(() => [
  { value: 'primary', label: 'Primary' },
  { value: 'secondary', label: 'Secondary' },
  { value: 'success', label: 'Success' },
  { value: 'warning', label: 'Warning' },
  { value: 'danger', label: 'Danger' },
  { value: 'info', label: 'Info' },
])

const actionPlacementOptions = computed(() => [
  { value: 'start', label: t('page.page-layout.recordToolbar.actions.placement.start') },
  { value: 'center', label: t('page.page-layout.recordToolbar.actions.placement.center') },
  { value: 'end', label: t('page.page-layout.recordToolbar.actions.placement.end') },
])

// ─── Methods ────────────────────────────────────────────────────────────────

function ensureLayoutKey (layout) {
  if (!layout._key) {
    layout._key = layout.pageLayoutID && layout.pageLayoutID !== NoID
      ? layout.pageLayoutID
      : `new_${++layoutKeyCounter}`
  }
}

async function loadPage () {
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

    // Load layouts for this page
    await loadLayouts()
  } catch (e) {
    console.error('Failed to load page:', e)
    $toast.toastDanger(t('notification.page.loadFailed'))
    router.push({ name: 'admin.pages' })
  } finally {
    loading.value = false
  }
}

async function loadLayouts () {
  if (!page.value?.pageID || page.value.pageID === NoID) return

  try {
    const storeLayouts = await pageLayoutStore.findByPageID({
      namespaceID: props.namespace.namespaceID,
      pageID: page.value.pageID,
      force: true,
    })

    layouts.value = storeLayouts.map(l => {
      const layout = new compose.PageLayout(l)
      ensureLayoutKey(layout)
      return layout
    })
    removedLayouts.value = []
  } catch (e) {
    console.error('Failed to load page layouts:', e)
  }
}

// ─── Layout CRUD ────────────────────────────────────────────────────────────

function addLayout () {
  const layout = new compose.PageLayout({
    namespaceID: props.namespace.namespaceID,
    pageID: page.value.pageID,
  })
  ensureLayoutKey(layout)
  layout._updated = true
  layouts.value.push(layout)
}

function removeLayout (index) {
  const layout = layouts.value[index]
  if (layout && layout.pageLayoutID !== NoID) {
    removedLayouts.value.push(layout)
  }
  layouts.value.splice(index, 1)
}

function moveLayout (index, direction) {
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= layouts.value.length) return

  const item = layouts.value.splice(index, 1)[0]
  layouts.value.splice(newIndex, 0, item)
}

// ─── Layout Config Dialog ───────────────────────────────────────────────────

function openLayoutConfig (layout) {
  const idx = layouts.value.indexOf(layout)
  configLayoutIndex.value = idx
  // Deep clone the layout for editing
  configLayout.value = JSON.parse(JSON.stringify(layout))

  // Resolve roles
  configLayoutRoles.value = (configLayout.value.config?.visibility?.roles || [])
    .map(roleID => ({ roleID }))

  layoutConfigVisible.value = true
}

function onLayoutConfigClose () {
  configLayout.value = null
  configLayoutIndex.value = -1
  configLayoutRoles.value = []
}

function saveLayoutConfig () {
  if (configLayoutIndex.value >= 0 && configLayout.value) {
    // Apply roles back
    configLayout.value.config.visibility.roles = configLayoutRoles.value.map(r => r.roleID || r)
    configLayout.value._updated = true

    // Preserve the _key
    configLayout.value._key = layouts.value[configLayoutIndex.value]._key
    layouts.value.splice(configLayoutIndex.value, 1, configLayout.value)
  }
  layoutConfigVisible.value = false
}

function onConfigLayoutRoleChange (roles) {
  configLayoutRoles.value = roles
}

// ─── Layout Actions ─────────────────────────────────────────────────────────

function addLayoutAction () {
  if (!configLayout.value) return
  if (!configLayout.value.config.actions) {
    configLayout.value.config.actions = []
  }
  configLayout.value.config.actions.push({
    kind: 'toLayout',
    placement: 'end',
    enabled: true,
    params: {
      pageLayoutID: '',
    },
    meta: {
      label: '',
      style: {
        variant: 'primary',
      },
    },
  })
}

function removeLayoutAction (index) {
  configLayout.value.config.actions.splice(index, 1)
}

function onActionKindChange (action) {
  if (action.kind === 'toURL' && !action.params.openIn) {
    action.params.openIn = 'sameTab'
  }
}

// ─── Save / Delete ──────────────────────────────────────────────────────────

async function handleSubmit ({ valid }) {
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
      config: page.value.config || {},
      meta: page.value.meta || {},
    }

    if (isEdit.value) {
      payload.pageID = page.value.pageID
      await pageStore.update(payload)

      // Save layouts
      await saveLayouts()

      $toast.toastSuccess(t('notification.page.saved'))
    } else {
      const created = await pageStore.create(payload)

      // Auto-create a primary layout for the new page (matching Corteza)
      await pageLayoutStore.create({
        namespaceID: props.namespace.namespaceID,
        pageID: created.pageID,
        handle: 'primary',
        meta: { title: created.title || payload.title },
      })

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

async function saveLayouts () {
  const { namespaceID } = props.namespace
  const pageID = page.value.pageID

  // Delete removed layouts first (so old handles don't interfere)
  await Promise.all(
    removedLayouts.value.map(layout =>
      pageLayoutStore.delete({
        namespaceID,
        pageID,
        pageLayoutID: layout.pageLayoutID,
      }).catch(e => console.error('Failed to delete layout:', e)),
    ),
  )

  // Create or update layouts
  await Promise.all(
    layouts.value.map(layout => {
      if (layout.pageLayoutID === NoID) {
        return pageLayoutStore.create({
          namespaceID,
          pageID,
          handle: layout.handle,
          meta: layout.meta,
          config: layout.config,
          blocks: layout.blocks || [],
        })
      } else if (layout._updated) {
        return pageLayoutStore.update({
          namespaceID,
          pageID,
          pageLayoutID: layout.pageLayoutID,
          handle: layout.handle,
          meta: layout.meta,
          config: layout.config,
          blocks: layout.blocks || [],
        })
      }
      return Promise.resolve()
    }),
  )

  // Reorder
  const pageIDs = layouts.value.map(l => l.pageLayoutID).filter(id => id !== NoID)
  if (pageIDs.length > 0) {
    try {
      await $ComposeAPI.pageLayoutReorder({ namespaceID, pageID, pageIDs })
      // Reload the store to reflect new order
      await pageLayoutStore.load({ namespaceID, clear: true, force: true })
    } catch (e) {
      console.error('Failed to reorder layouts:', e)
    }
  }

  // Reload layouts from API
  removedLayouts.value = []
  await loadLayouts()
}

function toggleDeleteMenu (event) {
  deleteMenu.value.toggle(event)
}

async function handleDelete (strategy = 'abort') {
  deleting.value = true
  try {
    await pageStore.delete({
      pageID: page.value.pageID,
      namespaceID: props.namespace.namespaceID,
      strategy,
    })
    $toast.toastSuccess(t('notification.page.deleted'))
    router.push({ name: 'admin.pages' })
  } catch (e) {
    console.error('Failed to delete page:', e)
    $toast.toastDanger(t('notification.page.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

// ─── Icon ───────────────────────────────────────────────────────────────────

function openIconModal () {
  // TODO: Implement full icon management dialog (upload, select, URL link)
  console.warn('Icon management dialog not yet implemented')
}

// ─── Navigation ─────────────────────────────────────────────────────────────

function goToBuilder () {
  router.push({
    name: 'admin.pages.builder',
    params: { pageID: page.value.pageID },
  })
}

function goToViewPage () {
  router.push({
    name: 'page',
    params: { pageID: page.value.pageID },
  })
}

function goToModuleEdit () {
  if (page.value?.moduleID) {
    router.push({
      name: 'admin.modules.edit',
      params: { moduleID: page.value.moduleID },
    })
  }
}

function goToLayoutBuilder (layout) {
  router.push({
    name: 'admin.pages.builder',
    params: { pageID: page.value.pageID },
    query: { layoutID: layout.pageLayoutID },
  })
}

// ─── Lifecycle ──────────────────────────────────────────────────────────────

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
