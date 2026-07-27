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
      <PageTranslator
        v-if="page && namespace"
        :page="page"
        :namespace="namespace"
        :layouts="layouts.filter(l => l.pageLayoutID !== '0')"
        @update:page="page = $event"
        @update:layouts="layouts = $event"
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
    :resolver="resolver"
    :initialValues="initialValues"
    @submit="handleSubmit"
    class="flex flex-col h-full"
  >
    <div class="container mx-auto p-4 flex-1 overflow-auto flex flex-col gap-4">
      <div v-if="isEdit && page?.canGrant" class="flex justify-end">
        <CPermissionsButton
          :resource="`corteza::compose:page/${page.namespaceID}/${page.pageID}`"
          :title="page.title || page.handle || page.pageID"
          :target="page.title || page.handle || page.pageID"
          v-tooltip.bottom="$t('general.label.permissions')"
          severity="secondary"
          size="small"
        />
      </div>

      <!-- General Panel -->
      <Panel :header="$t('general.label.general')" toggleable>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <CFormGroup name="title" :label="$t('page.label.title')" required>
            <InputText id="title" name="title" v-model="page.title" />
          </CFormGroup>

          <CFormGroup name="handle" :label="$t('page.label.handle')">
            <InputText id="handle" name="handle" v-model="page.handle" />
          </CFormGroup>
        </div>

        <CFormGroup :label="$t('page.label.description')" input-id="description" class="mb-6">
          <Textarea id="description" v-model="page.description" rows="4" auto-resize />
        </CFormGroup>

        <!-- Page Icon + Other Options (side by side like Human) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <!-- Page Icon -->
          <CFormGroup :label="$t('page.icon.page')">
            <template #actions>
              <Button
                v-tooltip.top="$t('page.icon.configure')"
                icon="pi pi-pencil"
                text
                severity="secondary"
                size="small"
                @click="openIconModal"
              />
            </template>

            <div class="inline-flex">
              <img v-if="pageIconSrc" :src="pageIconSrc" class="h-10 w-auto" />
              <span v-else class="text-muted-color">
                {{ $t('page.icon.noIcon') }}
              </span>
            </div>
          </CFormGroup>

          <!-- Other Options -->
          <CFormGroup>
            <CInputToggleCard
              v-model="page.visible"
              :label="$t('page.edit.visible')"
              :description="$t('page.edit.visibleDescription')"
            />

            <CInputToggleCard
              v-model="showSubPages"
              :label="$t('page.showSubPages')"
              :description="$t('page.showSubPagesDescription')"
            />

            <CInputToggleCard
              v-if="isRecordPage"
              v-model="notificationsEnabled"
              :label="$t('page.edit.notifications.enabled')"
              :description="$t('page.edit.notifications.description')"
            />
          </CFormGroup>
        </div>
      </Panel>

      <!-- Layouts Panel -->
      <Panel v-if="isEdit" :header="$t('page.page-layout.layouts')" toggleable>
        <div class="flex items-center mb-4">
          <Button
            :label="$t('page.page-layout.add')"
            icon="pi pi-plus"
            severity="secondary"
            size="small"
            @click="addLayout"
          />
        </div>

        <CFormList
          v-model="layouts"
          hide-remove
          draggable
          :empty-message="$t('page.noLayouts')"
          :columns="[
            {
              label: $t('page.page-layout.title'),
              width: '1fr',
              tooltip: $t('page.page-layout.tooltip.title'),
            },
            {
              label: $t('page.page-layout.handle'),
              width: '1fr',
              tooltip: $t('page.page-layout.tooltip.handle'),
            },
            { width: '2.5rem' },
          ]"
        >
          <template #row="{ item, index }">
            <div class="flex flex-col gap-1">
              <InputText
                v-model="item.meta.title"
                class="w-full"
                size="small"
                :invalid="validationTriggered && (!item.meta?.title || !item.meta.title.trim())"
                @input="item._updated = true"
              />
              <Message
                v-if="validationTriggered && (!item.meta?.title || !item.meta.title.trim())"
                severity="error"
                size="small"
                variant="simple"
              >
                {{ $t('general.label.required') }}
              </Message>
            </div>

            <InputGroup>
              <InputText
                v-model="item.handle"
                class="w-full"
                size="small"
                @input="item._updated = true"
              />
              <InputGroupAddon>
                <Button
                  v-tooltip.top="$t('page.page-layout.tooltip.configure')"
                  icon="pi pi-cog"
                  severity="secondary"
                  size="small"
                  class="w-full border-none"
                  @click="openLayoutConfig(item)"
                />
              </InputGroupAddon>
              <InputGroupAddon>
                <Button
                  v-tooltip.top="$t('page.page-layout.tooltip.builder')"
                  icon="pi pi-wrench"
                  severity="secondary"
                  size="small"
                  class="w-full border-none"
                  :disabled="item.pageLayoutID === NoID"
                  @click="goToLayoutBuilder(item)"
                />
              </InputGroupAddon>
            </InputGroup>

            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click="removeLayout(index)"
            />
          </template>
        </CFormList>
      </Panel>
    </div>

    <CEditorActions :back-to="{ name: 'admin.pages' }">
      <!-- Delete with strategy for pages with children -->
      <template v-if="isEdit && page.canDeletePage">
        <template v-if="hasChildren">
          <Button
            :label="$t('general.label.delete')"
            icon="pi pi-trash"
            severity="danger"
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
        v-if="isEdit && !isRecordPage"
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
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <CFormGroup :label="$t('page.page-layout.title')">
          <InputGroup v-if="isRecordPage && configLayout.config.useTitle">
            <InputGroupAddon>ƒ</InputGroupAddon>
            <InputText
              v-model="configLayout.meta.title"
              :placeholder="$t('page.page-layout.titleExpressionPlaceholder')"
            />
          </InputGroup>
          <InputText v-else v-model="configLayout.meta.title" />
        </CFormGroup>
        <CFormGroup :label="$t('page.page-layout.handle')">
          <InputText v-model="configLayout.handle" />
        </CFormGroup>
      </div>

      <!-- Use Title (record pages only) -->
      <CInputToggleCard
        v-if="isRecordPage"
        v-model="configLayout.config.useTitle"
        :label="$t('page.page-layout.useTitle')"
        :description="$t('page.page-layout.useTitleDescription')"
        class="mb-4"
      />

      <Divider />

      <CFormGroup :label="$t('page.page-layout.condition.label')" class="mb-4">
        <InputGroup>
          <InputGroupAddon>ƒ</InputGroupAddon>
          <InputText
            v-model="configLayout.config.visibility.expression"
            :placeholder="$t('page.page-layout.condition.placeholder')"
          />
        </InputGroup>
        <template #description>
          <template v-if="isRecordPage">
            {{
              $t('page.page-layout.condition.description.record-page', [
                'record.values.fieldName',
                'user.(userID/email...)',
                'screen.(width/height)',
                'isView/isCreate/isEdit',
                'user.userID == record.createdBy',
                'screen.width < 1024',
              ])
            }}
          </template>
          <template v-else>
            {{
              $t('page.page-layout.condition.description.non-record-page', [
                'user.(userID/email...)',
                'screen.(width/height)',
                'user.email == "test@mail.com"',
                'screen.width < 1024',
              ])
            }}
          </template>
        </template>
      </CFormGroup>

      <CFormGroup :label="$t('page.page-layout.roles.label')" class="mb-4">
        <CInputRole
          :value="configLayoutRoles"
          :placeholder="$t('page.page-layout.roles.placeholder')"
          multiple
          @input="onConfigLayoutRoleChange"
        />
      </CFormGroup>

      <!-- Record Toolbar (record pages only) -->
      <template v-if="isRecordPage">
        <Divider />

        <CFormGroup :label="$t('page.page-layout.recordToolbar.buttons.label')" class="mb-4">
          <div class="flex flex-col gap-2">
            <div class="flex items-center gap-3">
              <Checkbox
                v-model="configLayout.config.buttons.back.enabled"
                :binary="true"
                input-id="btn-back"
              />
              <label for="btn-back">
                {{ $t('page.page-layout.recordToolbar.buttons.showBack') }}
              </label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox
                v-model="configLayout.config.buttons.delete.enabled"
                :binary="true"
                input-id="btn-delete"
              />
              <label for="btn-delete">
                {{ $t('page.page-layout.recordToolbar.buttons.showDelete') }}
              </label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox
                v-model="configLayout.config.buttons.clone.enabled"
                :binary="true"
                input-id="btn-clone"
              />
              <label for="btn-clone">
                {{ $t('page.page-layout.recordToolbar.buttons.showClone') }}
              </label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox
                v-model="configLayout.config.buttons.new.enabled"
                :binary="true"
                input-id="btn-new"
              />
              <label for="btn-new">
                {{ $t('page.page-layout.recordToolbar.buttons.showNew') }}
              </label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox
                v-model="configLayout.config.buttons.edit.enabled"
                :binary="true"
                input-id="btn-edit"
              />
              <label for="btn-edit">
                {{ $t('page.page-layout.recordToolbar.buttons.showEdit') }}
              </label>
            </div>
            <div class="flex items-center gap-3">
              <Checkbox
                v-model="configLayout.config.buttons.submit.enabled"
                :binary="true"
                input-id="btn-submit"
              />
              <label for="btn-submit">
                {{ $t('page.page-layout.recordToolbar.buttons.showSave') }}
              </label>
            </div>
          </div>
        </CFormGroup>

        <Divider />

        <!-- Custom Actions -->
        <CFormGroup :label="$t('page.page-layout.recordToolbar.actions.label')">
          <template #actions>
            <Button
              :label="$t('general.label.add')"
              icon="pi pi-plus"
              severity="secondary"
              size="small"
              @click="addLayoutAction"
            />
          </template>

          <CFormList
            v-model="configLayout.config.actions"
            draggable
            :empty-message="$t('page.page-layout.recordToolbar.actions.empty')"
            :columns="[
              { label: $t('page.page-layout.recordToolbar.actions.buttonLabel'), width: '1fr' },
              { label: $t('page.page-layout.recordToolbar.actions.kind.label'), width: '180px' },
              { label: $t('page.page-layout.recordToolbar.actions.variant'), width: '140px' },
              {
                label: $t('page.page-layout.recordToolbar.actions.placement.label'),
                width: '120px',
              },
              {
                label: $t('page.page-layout.recordToolbar.actions.visible'),
                width: '5rem',
                headerClass: 'text-center',
              },
            ]"
          >
            <template #row="{ item }">
              <InputText v-model="item.meta.label" class="w-full" size="small" />

              <Select
                v-model="item.kind"
                :options="actionKindOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                size="small"
                @change="onActionKindChange(item)"
              />

              <Select
                v-model="item.meta.style.variant"
                :options="actionVariantOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                size="small"
              />

              <Select
                v-model="item.placement"
                :options="actionPlacementOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                size="small"
              />

              <div class="flex justify-center">
                <Checkbox v-model="item.enabled" :binary="true" />
              </div>
            </template>

            <template #extra="{ item }">
              <div
                v-if="item.kind === 'toLayout' || item.kind === 'toURL'"
                class="border-t border-surface pt-3 mt-1 grid grid-cols-1 md:grid-cols-2 gap-3"
              >
                <CFormGroup
                  v-if="item.kind === 'toLayout'"
                  :label="$t('page.page-layout.recordToolbar.actions.toLayout.label')"
                >
                  <Select
                    v-model="item.params.pageLayoutID"
                    :options="actionLayoutOptions"
                    option-label="label"
                    option-value="value"
                    class="w-full"
                    size="small"
                  />
                </CFormGroup>

                <template v-else-if="item.kind === 'toURL'">
                  <CFormGroup :label="$t('page.page-layout.recordToolbar.actions.toURL.label')">
                    <InputText
                      v-model="item.params.url"
                      :placeholder="$t('page.page-layout.recordToolbar.actions.toURL.placeholder')"
                      size="small"
                    />
                  </CFormGroup>
                  <CFormGroup :label="$t('page.page-layout.recordToolbar.actions.openIn.label')">
                    <Select
                      v-model="item.params.openIn"
                      :options="actionOpenInOptions"
                      option-label="label"
                      option-value="value"
                      class="w-full"
                      size="small"
                    />
                  </CFormGroup>
                </template>
              </div>
            </template>
          </CFormList>
        </CFormGroup>
      </template>
    </template>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="layoutConfigVisible = false"
        />
        <Button
          :label="$t('general.label.save')"
          size="small"
          :disabled="!configLayout?.meta?.title"
          @click="saveLayoutConfig"
        />
      </div>
    </template>
  </Dialog>

  <!-- Icon Configuration Dialog -->
  <Dialog
    v-model:visible="showIconModal"
    :header="$t('page.icon.configure')"
    modal
    :style="{ width: '40rem' }"
    :breakpoints="{ '768px': '90vw' }"
    :closable="true"
    @hide="closeIconModal"
  >
    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('page.icon.upload')">
        <CFileDropZone
          accept="image/*"
          :uploading="iconUploading"
          :error="iconUploadError"
          :drop-label="$t('general.label.dropFiles')"
          compact
          @select="onIconFileSelect"
        />
      </CFormGroup>

      <CFormGroup :label="$t('page.url.label')">
        <InputGroup>
          <InputText v-model="linkUrl" :disabled="isIconSet" />
          <InputGroupAddon>
            <Button
              v-tooltip.top="$t('page.tooltip.preview-link')"
              icon="pi pi-external-link"
              severity="secondary"
              :disabled="!linkUrl"
              @click="showLinkPreview = true"
            />
          </InputGroupAddon>
        </InputGroup>
      </CFormGroup>

      <template v-if="attachments.length > 0">
        <Divider />

        <CFormGroup :label="$t('page.icon.list')" class="mb-4">
          <div v-if="processingIcon" class="flex items-center justify-center h-24">
            <ProgressSpinner style="width: 2rem; height: 2rem" />
          </div>
          <div v-else class="flex flex-wrap gap-2">
            <img
              v-for="a in attachments"
              :key="a.attachmentID"
              :src="a.src"
              :alt="a.name"
              class="h-12 w-auto rounded cursor-pointer p-1 border-2"
              :class="
                selectedAttachmentID === a.attachmentID ? 'border-primary' : 'border-transparent'
              "
              @click="toggleSelectedIcon(a.attachmentID)"
            />
          </div>
        </CFormGroup>
      </template>
    </div>

    <template #footer>
      <div class="flex items-center w-full">
        <CInputDelete
          v-if="selectedAttachmentID"
          :label="$t('page.icon.delete')"
          :message="$t('page.icon.delete')"
          :header="$t('page.icon.delete')"
          :disabled="processingIcon"
          size="small"
          @confirm="deleteIcon"
        />
        <div class="ml-auto flex gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="closeIconModal"
          />
          <Button :label="$t('general.label.saveAndClose')" size="small" @click="saveIconModal" />
        </div>
      </div>
    </template>
  </Dialog>

  <Dialog
    v-model:visible="showLinkPreview"
    modal
    dismissable-mask
    :closable="true"
    :show-header="false"
    :style="{ maxWidth: '90vw' }"
  >
    <img :src="linkUrl" class="max-w-full h-auto" />
  </Dialog>
</template>

<script setup>
import { usePageStore } from '@planetcrust/human-vue'
import { usePageLayoutStore } from '@planetcrust/human-vue'
import { compose, NoID } from '@planetcrust/human-js'
import { components, useFileUpload, useUnsavedGuard } from '@planetcrust/human-vue'
import { cloneDeep, isEqual } from 'lodash-es'
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'

const { CInputDelete, CInputToggleCard, CFileDropZone } = components
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import PageTranslator from '@/sections/compose/components/Admin/Page/PageTranslator.vue'

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
const cloning = ref(false)
const page = ref(null)
const initialPage = ref(null)

// Page layouts state
const layouts = ref([])
const initialLayouts = ref([])
const removedLayouts = ref([])

const { markSaved } = useUnsavedGuard({
  isDirty: () => {
    if (saving.value || deleting.value || cloning.value) return false
    if (!page.value || !initialPage.value) return false
    if (!isEqual(page.value, initialPage.value)) return true
    if (!isEqual(layouts.value, initialLayouts.value)) return true
    return false
  },
  messageKey: 'general.editor.unsavedChanges',
})
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

const layoutsValid = computed(() => {
  return layouts.value.every(l => !!l.meta?.title && l.meta.title.trim().length > 0)
})

const canSave = computed(() => {
  if (isEdit.value && !page.value?.canUpdatePage) return false
  return true
})

const validationTriggered = ref(false)

const hasChildren = computed(() => {
  if (!page.value) return false
  return pageStore.set.some(p => p.selfID === page.value.pageID)
})

// ─── Icon state ─────────────────────────────────────────────────────────────
const showIconModal = ref(false)
const showLinkPreview = ref(false)
const attachments = ref([])
const selectedAttachmentID = ref('')
const linkUrl = ref('')
const processingIcon = ref(false)

const {
  uploading: iconUploading,
  uploadError: iconUploadError,
  uploadFileRaw: uploadIconRaw,
  reset: resetIconUpload,
} = useFileUpload()

const pageIcon = computed({
  get() {
    return page.value?.config?.navItem?.icon || {}
  },
  set(icon) {
    if (!page.value) return
    if (!page.value.config) page.value.config = {}
    if (!page.value.config.navItem) page.value.config.navItem = {}
    page.value.config.navItem.icon = icon
  },
})

function makeAttachmentUrl(src) {
  return `${$ComposeAPI.baseURL}${src}`
}

// Page icon source
const pageIconSrc = computed(() => {
  const icon = pageIcon.value
  if (!icon?.src) return ''
  return icon.type === 'link' ? icon.src : makeAttachmentUrl(icon.src)
})

const isIconSet = computed(() => !!selectedAttachmentID.value)

// Show sub-pages computed property with getter/setter
const showSubPages = computed({
  get() {
    return page.value?.config?.navItem?.expanded ?? false
  },
  set(val) {
    if (page.value) {
      if (!page.value.config) page.value.config = {}
      if (!page.value.config.navItem) page.value.config.navItem = {}
      page.value.config.navItem.expanded = val
    }
  },
})

// Notifications enabled computed property
const notificationsEnabled = computed({
  get() {
    return page.value?.meta?.notifications?.enabled ?? false
  },
  set(val) {
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
  const options = [
    { value: '', label: t('page.page-layout.recordToolbar.actions.toLayout.placeholder') },
  ]
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

function ensureLayoutKey(layout) {
  if (!layout._key) {
    layout._key =
      layout.pageLayoutID && layout.pageLayoutID !== NoID
        ? layout.pageLayoutID
        : `new_${++layoutKeyCounter}`
  }
}

async function loadPage() {
  const pageID = route.params.pageID
  if (!pageID) {
    // Create new
    page.value = new compose.Page({
      namespaceID: props.namespace?.namespaceID,
      visible: true,
    })
    initialPage.value = cloneDeep(page.value)
    initialLayouts.value = []
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
    await fetchAttachments()
    initialPage.value = cloneDeep(page.value)
    initialLayouts.value = cloneDeep(layouts.value)
  } catch (e) {
    console.error('Failed to load page:', e)
    $toast.toastDanger(t('notification.page.loadFailed'))
    router.push({ name: 'admin.pages' })
  } finally {
    loading.value = false
  }
}

async function loadLayouts() {
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

function addLayout() {
  const layout = new compose.PageLayout({
    namespaceID: props.namespace.namespaceID,
    pageID: page.value.pageID,
  })
  ensureLayoutKey(layout)
  layout._updated = true
  layouts.value.push(layout)
}

function removeLayout(index) {
  const layout = layouts.value[index]
  if (layout && layout.pageLayoutID !== NoID) {
    removedLayouts.value.push(layout)
  }
  layouts.value.splice(index, 1)
}

// ─── Layout Config Dialog ───────────────────────────────────────────────────

function openLayoutConfig(layout) {
  const idx = layouts.value.indexOf(layout)
  configLayoutIndex.value = idx
  // Deep clone the layout for editing
  configLayout.value = JSON.parse(JSON.stringify(layout))

  // Resolve roles
  configLayoutRoles.value = (configLayout.value.config?.visibility?.roles || []).map(roleID => ({
    roleID,
  }))

  layoutConfigVisible.value = true
}

function onLayoutConfigClose() {
  configLayout.value = null
  configLayoutIndex.value = -1
  configLayoutRoles.value = []
}

function saveLayoutConfig() {
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

function onConfigLayoutRoleChange(roles) {
  configLayoutRoles.value = roles
}

// ─── Layout Actions ─────────────────────────────────────────────────────────

function addLayoutAction() {
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

function onActionKindChange(action) {
  if (action.kind === 'toURL' && !action.params.openIn) {
    action.params.openIn = 'sameTab'
  }
}

// ─── Save / Delete ──────────────────────────────────────────────────────────

async function handleSubmit({ valid }) {
  if (!valid || !layoutsValid.value) {
    validationTriggered.value = true
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

      // Save icon first; merge the returned icon into the page payload
      const savedIcon = await saveIcon()
      if (savedIcon) {
        if (!payload.config.navItem) payload.config.navItem = {}
        payload.config.navItem.icon = savedIcon
      }

      const updated = await pageStore.update(payload)
      page.value = new compose.Page({ ...updated })

      // Save layouts
      await saveLayouts()

      initialPage.value = cloneDeep(page.value)
      initialLayouts.value = cloneDeep(layouts.value)

      $toast.toastSuccess(t('notification.page.saved'))
    } else {
      const created = await pageStore.create(payload)

      // Auto-create a primary layout for the new page (matching Human).
      // Use the PageLayout type so the default config (all record toolbar
      // buttons enabled) is persisted instead of zero-value disabled buttons.
      await pageLayoutStore.create(
        new compose.PageLayout({
          namespaceID: props.namespace.namespaceID,
          pageID: created.pageID,
          handle: 'primary',
          meta: { title: created.title || payload.title },
        }),
      )

      $toast.toastSuccess(t('notification.page.created'))
      markSaved()
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

async function saveLayouts() {
  const { namespaceID } = props.namespace
  const pageID = page.value.pageID

  // Delete removed layouts first (so old handles don't interfere)
  await Promise.all(
    removedLayouts.value.map(layout =>
      pageLayoutStore
        .delete({
          namespaceID,
          pageID,
          pageLayoutID: layout.pageLayoutID,
        })
        .catch(e => console.error('Failed to delete layout:', e)),
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
      await pageLayoutStore.load({ namespaceID, clear: true })
    } catch (e) {
      console.error('Failed to reorder layouts:', e)
    }
  }

  // Reload layouts from API
  removedLayouts.value = []
  await loadLayouts()
}

function toggleDeleteMenu(event) {
  deleteMenu.value.toggle(event)
}

async function handleClone() {
  cloning.value = true
  try {
    const payload = {
      namespaceID: props.namespace.namespaceID,
      title: `${page.value.title} (${t('general.label.clone').toLowerCase()})`,
      handle: '',
      description: page.value.description,
      visible: page.value.visible,
      selfID: page.value.selfID || '0',
      blocks: page.value.blocks || [],
      config: page.value.config || {},
      meta: page.value.meta || {},
    }

    const created = await pageStore.create(payload)
    // Seed layout with the cloned blocks — View.vue intersects layout.blocks with
    // page.blocks, so an empty layout would hide everything until the next save.
    const layoutBlocks = (created.blocks || []).map(b => ({
      blockID: b.blockID,
      xywh: b.xywh,
    }))
    // Use the PageLayout type so the default config (all record toolbar
    // buttons enabled) is persisted instead of zero-value disabled buttons.
    await pageLayoutStore.create(
      new compose.PageLayout({
        namespaceID: props.namespace.namespaceID,
        pageID: created.pageID,
        handle: 'primary',
        meta: { title: created.title || payload.title },
        blocks: layoutBlocks,
      }),
    )
    $toast.toastSuccess(t('notification.page.created'))
    markSaved()
    router.push({
      name: 'admin.pages.edit',
      params: { pageID: created.pageID },
    })
  } catch (e) {
    console.error('Failed to clone page:', e)
    $toast.toastDanger(t('notification.page.cloneFailed'))
  } finally {
    cloning.value = false
  }
}

async function handleDelete(strategy = 'abort') {
  deleting.value = true
  try {
    await pageStore.delete({
      pageID: page.value.pageID,
      namespaceID: props.namespace.namespaceID,
      strategy,
    })
    initialPage.value = cloneDeep(page.value)
    initialLayouts.value = cloneDeep(layouts.value)
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

async function fetchAttachments() {
  processingIcon.value = true
  try {
    const { set = [] } = await $ComposeAPI.iconList({ sort: 'id DESC' })
    const baseURL = $ComposeAPI.baseURL
    if (set.length === 0) {
      attachments.value = []
      if (page.value && pageIcon.value?.src && pageIcon.value.type !== 'link') {
        pageIcon.value = {}
      }
      if (initialPage.value?.config?.navItem) {
        initialPage.value.config.navItem.icon = cloneDeep(pageIcon.value)
      }
    } else {
      attachments.value = set.map(a => ({
        ...a,
        src: a.url && !a.url.includes(baseURL) ? makeAttachmentUrl(a.url) : a.url,
      }))
    }
  } catch (e) {
    console.error('Failed to fetch icons:', e)
    $toast.toastDanger(t('notification.page.iconFetchFailed'))
  } finally {
    processingIcon.value = false
  }
}

function setCurrentIcon() {
  const match = attachments.value.find(a => a.url === pageIcon.value?.src)
  selectedAttachmentID.value = match?.attachmentID || ''
  if (!selectedAttachmentID.value && pageIcon.value?.type !== 'link') {
    pageIcon.value = {}
  }
}

function toggleSelectedIcon(attachmentID = '') {
  selectedAttachmentID.value = selectedAttachmentID.value === attachmentID ? '' : attachmentID
}

function openIconModal() {
  linkUrl.value = pageIcon.value?.type === 'link' ? pageIcon.value.src : ''
  setCurrentIcon()
  resetIconUpload()
  showIconModal.value = true
}

function closeIconModal() {
  linkUrl.value = pageIcon.value?.type === 'link' ? pageIcon.value.src : ''
  setCurrentIcon()
  resetIconUpload()
  showIconModal.value = false
}

function saveIconModal() {
  const type = selectedAttachmentID.value ? 'attachment' : 'link'
  let src = linkUrl.value
  if (selectedAttachmentID.value) {
    const att = attachments.value.find(
      ({ attachmentID }) => attachmentID === selectedAttachmentID.value,
    )
    src = att?.url || ''
  }

  pageIcon.value = type === 'link' && !src ? {} : { type, src }
  showIconModal.value = false
}

async function onIconFileSelect(files) {
  const file = files?.[0]
  if (!file) return
  try {
    const endpoint = $ComposeAPI.baseURL + $ComposeAPI.iconUploadEndpoint()
    const token = $ComposeAPI.accessTokenFn ? $ComposeAPI.accessTokenFn() : ''
    const att = await uploadIconRaw(file, {
      url: endpoint,
      token,
      fieldName: 'icon',
    })
    await fetchAttachments()
    if (att?.attachmentID) {
      toggleSelectedIcon(att.attachmentID)
    }
  } catch (e) {
    console.error('Failed to upload icon:', e)
  }
}

async function deleteIcon() {
  if (!selectedAttachmentID.value) return
  processingIcon.value = true
  try {
    await $ComposeAPI.iconDelete({ iconID: selectedAttachmentID.value })
    await fetchAttachments()
    setCurrentIcon()
    $toast.toastSuccess(t('notification.page.iconDeleteSuccess'))
  } catch (e) {
    console.error('Failed to delete icon:', e)
    $toast.toastDanger(t('notification.page.iconDeleteFailed'))
  } finally {
    processingIcon.value = false
  }
}

async function saveIcon() {
  if (!page.value?.pageID || page.value.pageID === NoID) return null
  const icon = pageIcon.value || {}
  return $ComposeAPI.pageUpdateIcon({
    namespaceID: props.namespace.namespaceID,
    pageID: page.value.pageID,
    type: icon.type || 'link',
    source: icon.src || '',
  })
}

// ─── Navigation ─────────────────────────────────────────────────────────────

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

function goToModuleEdit() {
  if (page.value?.moduleID) {
    router.push({
      name: 'admin.modules.edit',
      params: { moduleID: page.value.moduleID },
    })
  }
}

function goToLayoutBuilder(layout) {
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
