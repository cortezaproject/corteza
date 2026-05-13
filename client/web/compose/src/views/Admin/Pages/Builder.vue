<template>
  <Teleport to="#topbar-title" defer>
    <span v-if="page">{{ $t('page.edit.pageBuilder') }}: {{ page.title }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <div v-if="page" class="flex items-center gap-2">
      <Select
        v-if="layouts.length > 1"
        :model-value="pageLayout?.pageLayoutID"
        :options="layoutOptions"
        option-label="label"
        option-value="value"
        size="small"
        style="min-width: 200px"
        @update:model-value="setLayout"
      />
      <ButtonGroup class="gap-1">
        <Button
          v-if="page.isRecordPage"
          :label="$t('page.moduleEdit')"
          icon="pi pi-database"
          size="small"
          @click="goToModuleEdit"
        />
        <Button
          :label="$t('page.edit.viewPage')"
          icon="pi pi-eye"
          size="small"
          @click="goToViewPage"
        />
        <Button
          v-tooltip.bottom="$t('navigation.editPage')"
          icon="pi pi-pencil"
          size="small"
          @click="goToEditPage"
        />
        <PageTranslator
          v-if="page.pageID && page.pageID !== '0' && namespace"
          :page="page"
          :namespace="namespace"
          :layouts="layouts.filter(l => l.pageLayoutID !== '0')"
          @update:page="page = $event"
          @update:layouts="layouts = $event"
        />
      </ButtonGroup>
    </div>
  </Teleport>

  <!-- Loading -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner />
  </div>

  <div v-else-if="page" class="flex flex-col h-full">
    <!-- Builder Grid -->
    <div class="flex-1 overflow-auto">
      <Grid ref="gridRef" :blocks="blocks" :namespace="namespace" :page="page" editable>
        <template #item-overlay="{ item }">
          <div class="block-toolbox bg-emphasis flex items-center">
            <Button
              v-tooltip.top="$t('page.tooltip.drag.block')"
              icon="pi pi-arrows-alt"
              text
              size="small"
              severity="secondary"
              class="block-drag-handle !cursor-move"
            />
            <Divider layout="vertical" class="!mx-1 !my-0" />
            <ButtonGroup>
              <Button
                :title="$t('page.tooltip.edit.block')"
                icon="pi pi-pencil"
                text
                size="small"
                severity="secondary"
                @click="editBlock(item.i)"
              />
              <Button
                :title="$t('page.tooltip.clone.block')"
                icon="pi pi-copy"
                text
                size="small"
                severity="secondary"
                @click="cloneBlock(item.i)"
              />
              <Button
                v-if="showTranslatorButton && blocks[item.i]?.blockID && blocks[item.i]?.blockID !== '0'"
                v-tooltip.top="$t('translator.button.tooltip')"
                icon="pi pi-language"
                text
                size="small"
                severity="secondary"
                @click.stop="openBlockTranslation(blocks[item.i])"
              />
              <Button
                :title="$t('page.tooltip.delete.block')"
                icon="pi pi-trash"
                text
                size="small"
                severity="danger"
                @click="deleteBlock(item.i)"
              />
            </ButtonGroup>
          </div>
        </template>
      </Grid>
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

        <Button
          :label="$t('page.build.addBlock')"
          icon="pi pi-plus"
          severity="secondary"
          @click="showAddBlock = true"
        />

        <div class="flex gap-2">
          <Button
            v-if="pageLayout"
            :label="$t('page.build.saveAsCopy')"
            icon="pi pi-copy"
            severity="secondary"
            :loading="saving"
            @click="handleSaveAsCopy"
          />
          <Button
            v-if="pageLayout && layouts.length > 1"
            v-tooltip.top="$t('page.build.layout.delete')"
            icon="pi pi-trash"
            severity="danger"
            outlined
            :loading="saving"
            @click="showDeleteLayoutConfirm = true"
          />
          <Button
            :label="$t('general.label.save')"
            icon="pi pi-save"
            :loading="saving"
            @click="handleSave"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Add Block Dialog -->
  <Dialog
    v-model:visible="showAddBlock"
    :header="$t('page.build.selectBlockTitle')"
    modal
    :style="{ width: '600px' }"
  >
    <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
      <template v-for="bt in addableBlockTypes" :key="bt.kind">
        <Divider v-if="bt.kind === 'divider'" class="col-span-2 md:col-span-3 my-0" />
        <div
          v-else
          class="flex flex-col items-center gap-2 p-4 border border-surface rounded cursor-pointer hover:bg-highlight transition-colors"
          @click="addBlock(bt.kind)"
        >
          <i :class="bt.icon" class="text-2xl text-primary" />
          <span class="text-sm font-medium text-center">{{ bt.label }}</span>
        </div>
      </template>
    </div>
  </Dialog>

  <!-- Delete Layout Confirmation Dialog -->
  <Dialog
    v-model:visible="showDeleteLayoutConfirm"
    :header="$t('page.build.layout.delete')"
    modal
    :style="{ width: '400px' }"
  >
    <p>{{ $t('page.build.layout.deleteConfirm') }}</p>
    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          size="small"
          outlined
          @click="showDeleteLayoutConfirm = false"
        />
        <Button
          :label="$t('general.label.delete')"
          severity="danger"
          size="small"
          :loading="saving"
          @click="handleDeleteLayout"
        />
      </div>
    </template>
  </Dialog>

  <!-- Block Configurator Dialog -->
  <Dialog
    v-model:visible="showConfigurator"
    :header="$t('page.changeBlock')"
    modal
    position="top"
    :style="{ width: '50vw' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :pt="{
      content: { class: 'p-0 flex flex-col !overflow-hidden' },
      footer: { class: 'border-t border-surface p-3' },
    }"
  >
    <Tabs v-if="editingBlock" v-model:value="configuratorTab" class="flex flex-col flex-1 min-h-0">
      <TabList class="shrink-0 z-10">
        <Tab value="general">{{ $t('block.general.label.general') }}</Tab>
        <Tab value="block">{{ editingBlockTypeLabel }}</Tab>
        <Tab v-if="hasAutomationTab" value="automation">{{ $t('block.automation.label') }}</Tab>
      </TabList>

      <TabPanels class="flex-1 overflow-y-auto">
        <!-- General Tab -->
        <TabPanel value="general">
          <div class="flex flex-col gap-3">
            <!-- Title -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('block.general.titleLabel') }}
              </label>
              <InputText
                v-model="editingBlock.title"
                :placeholder="$t('block.general.titlePlaceholder')"
                class="w-full"
              />
            </div>

            <!-- Description -->
            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('block.general.descriptionLabel') }}
              </label>
              <Textarea
                v-model="editingBlock.description"
                :placeholder="$t('block.general.descriptionPlaceholder')"
                rows="2"
                class="w-full"
              />
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <!-- Custom ID -->
              <div class="flex flex-col gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.general.customID.label') }}
                </label>
                <InputText
                  v-model="editingBlock.meta.customID"
                  :placeholder="$t('block.general.customID.placeholder')"
                  :invalid="customIDInvalid"
                  class="w-full"
                />
                <small v-if="customIDInvalid" class="text-red-500">
                  {{ $t('block.general.customID.invalid-state') }}
                </small>
                <small v-else class="text-muted-color">
                  {{ $t('block.general.customID.description') }}
                </small>
              </div>

              <!-- Custom CSS Class -->
              <div class="flex flex-col gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.general.customCSSClass.label') }}
                </label>
                <InputText
                  v-model="editingBlock.meta.customCSSClass"
                  :placeholder="$t('block.general.customCSSClass.placeholder')"
                  :invalid="customCSSClassInvalid"
                  class="w-full"
                />
                <small v-if="customCSSClassInvalid" class="text-red-500">
                  {{ $t('block.general.customCSSClass.invalid-state') }}
                </small>
                <small v-else class="text-muted-color">
                  {{ $t('block.general.customCSSClass.description') }}
                </small>
              </div>
            </div>

            <Divider />

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <!-- Header style -->
              <div class="flex flex-col gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.general.headerStyle') }}
                </label>
                <Select
                  v-model="editingBlock.style.variants.headerText"
                  :options="headerTextVariantOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                />
              </div>

              <!-- Magnify option -->
              <div class="flex flex-col gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.general.magnifyLabel') }}
                </label>
                <Select
                  v-model="editingBlock.options.magnifyOption"
                  :options="magnifyOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                />
              </div>
            </div>

            <!-- Wrap and border -->
            <div class="flex items-center gap-4">
              <div class="flex items-center gap-2">
                <Checkbox
                  :model-value="editingBlock.style.wrap.kind === 'card'"
                  binary
                  input-id="wrapCard"
                  @update:model-value="editingBlock.style.wrap.kind = $event ? 'card' : 'plain'"
                />
                <label for="wrapCard" class="text-sm">{{ $t('block.general.wrap') }}</label>
              </div>
              <div class="flex items-center gap-2">
                <Checkbox
                  v-model="editingBlock.style.border.enabled"
                  binary
                  input-id="borderEnabled"
                />
                <label for="borderEnabled" class="text-sm">
                  {{ $t('block.general.border.show') }}
                </label>
              </div>
            </div>

            <!-- Refresh rate -->
            <template v-if="editingBlock.options?.showRefresh !== undefined">
              <Divider />
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <div class="flex flex-col gap-1">
                  <label class="text-primary font-medium text-sm">
                    {{ $t('block.general.refresh.auto') }}
                  </label>
                  <div class="flex items-center gap-2">
                    <InputNumber
                      v-model="editingBlock.options.refreshRate"
                      :min="0"
                      suffix=" s"
                      class="w-full"
                      @blur="clampRefreshRate"
                    />
                  </div>
                  <small class="text-muted-color">
                    {{ $t('block.general.refresh.description') }}
                  </small>
                </div>
                <div class="flex gap-2">
                  <CInputSwitch
                    v-model="editingBlock.options.showRefresh"
                    :label="$t('block.general.refresh.show')"
                  />
                </div>
              </div>
            </template>

            <Divider />

            <!-- Visibility -->
            <h5 class="text-lg font-semibold text-primary m-0">
              {{ $t('block.general.visibility.label') }}
            </h5>

            <div class="flex flex-col gap-1">
              <div class="flex items-center gap-1">
                <label class="text-primary font-medium text-sm">
                  {{ $t('block.general.visibility.condition.label') }}
                </label>
                <i
                  class="pi pi-exclamation-triangle text-orange-500 text-xs"
                  v-tooltip="$t('block.general.visibility.tooltip.performance.condition')"
                />
              </div>
              <Textarea
                v-model="editingBlock.meta.visibility.expression"
                :placeholder="$t('block.general.visibility.condition.placeholder')"
                rows="2"
                class="w-full"
              />
              <small class="text-muted-color">
                {{
                  isRecordPage
                    ? $t('block.general.visibility.condition.description.record-page', [
                        'record.values.fieldName',
                        'user.(userID/email...)',
                        'screen.(width/height)',
                        'isView/isCreate/isEdit',
                        'user.userID == record.createdBy',
                        'screen.width < 1024',
                      ])
                    : $t('block.general.visibility.condition.description.non-record-page', [
                        'user.(userID/email...)',
                        'screen.(width/height)',
                        'user.email == "test@mail.com"',
                        'screen.width < 1024',
                      ])
                }}
              </small>
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-primary font-medium text-sm">
                {{ $t('block.general.visibility.roles.label') }}
              </label>
              <CInputRole
                v-model="editingBlock.meta.visibility.roles"
                multiple
                :placeholder="$t('block.general.visibility.roles.placeholder')"
                class="w-full"
              />
            </div>
          </div>
        </TabPanel>

        <!-- Block-specific Tab -->
        <TabPanel value="block">
          <component
            :is="blockConfigurator"
            v-if="blockConfigurator"
            :namespace="namespace"
            :page="page"
            :blocks="blocks"
            @edit-tab-block="onEditTabBlock"
            @create-tab-block="onCreateTabBlock"
          />

          <!-- Fallback when no configurator exists -->
          <div v-else class="p-3 text-muted-color italic">
            {{ editingBlock.kind }} {{ $t('block.general.label.noHandle') }}
          </div>
        </TabPanel>

        <!-- Automation Tab -->
        <TabPanel v-if="hasAutomationTab" value="automation">
          <AutomationButtonsEditor
            :buttons="editingBlock.options.selectionButtons || []"
            @update:buttons="onSelectionButtonsUpdate"
          />
        </TabPanel>
      </TabPanels>
    </Tabs>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          size="small"
          outlined
          @click="showConfigurator = false"
        />
        <Button :label="$t('general.label.save')" size="small" @click="saveBlockConfig" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, computed, watch, inject, provide, markRaw, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { compose } from '@planetcrust/human-js'
import { usePageStore } from '@/stores/page'
import { usePageLayoutStore } from '@/stores/page-layout'
import { useModuleStore } from '@/stores/module'
import Grid from '@/components/PageBlocks/Grid.vue'

// Block configurators — lazy imported
import ChartConfigurator from '@/components/PageBlocks/Configurators/ChartConfigurator.vue'
import ContentConfigurator from '@/components/PageBlocks/Configurators/ContentConfigurator.vue'
import RecordListConfigurator from '@/components/PageBlocks/Configurators/RecordListConfigurator.vue'
import RecordConfigurator from '@/components/PageBlocks/Configurators/RecordConfigurator.vue'
import MetricConfigurator from '@/components/PageBlocks/Configurators/MetricConfigurator.vue'
import IFrameConfigurator from '@/components/PageBlocks/Configurators/IFrameConfigurator.vue'
import FileConfigurator from '@/components/PageBlocks/Configurators/FileConfigurator.vue'
import CalendarConfigurator from '@/components/PageBlocks/Configurators/CalendarConfigurator.vue'
import CommentConfigurator from '@/components/PageBlocks/Configurators/CommentConfigurator.vue'
import AutomationConfigurator from '@/components/PageBlocks/Configurators/AutomationConfigurator.vue'
import NavigationConfigurator from '@/components/PageBlocks/Configurators/NavigationConfigurator.vue'
import TabsConfigurator from '@/components/PageBlocks/Configurators/TabsConfigurator.vue'
import ProgressConfigurator from '@/components/PageBlocks/Configurators/ProgressConfigurator.vue'
import RecordOrganizerConfigurator from '@/components/PageBlocks/Configurators/RecordOrganizerConfigurator.vue'
import RecordRevisionsConfigurator from '@/components/PageBlocks/Configurators/RecordRevisionsConfigurator.vue'
import GeometryConfigurator from '@/components/PageBlocks/Configurators/GeometryConfigurator.vue'
import ChatbotConfigurator from '@/components/PageBlocks/Configurators/ChatbotConfigurator.vue'
import AutomationButtonsEditor from '@/components/PageBlocks/Shared/AutomationButtonsEditor.vue'
import PageTranslator from '@/components/Admin/Page/PageTranslator.vue'
import { useResourceTranslations } from '@/composables/useResourceTranslations'
import { useTranslatorStore } from '@/stores/translator'
import { applyPageTranslations } from '@/lib/resource-translations'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')
const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()
const moduleStore = useModuleStore()
const { showTranslatorButton, currentLanguage } = useResourceTranslations()
const translatorStore = useTranslatorStore()

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

// State
const loading = ref(false)
const saving = ref(false)
const page = ref(null)
const layouts = ref([])
const pageLayout = ref(null)
const showDeleteLayoutConfirm = ref(false)
const blocks = ref([])
const showAddBlock = ref(false)
const showConfigurator = ref(false)
const editingBlock = ref(null)
const editingBlockIndex = ref(-1)

provide('blockDraft', editingBlock)
provide('$pageBuilder', {
  editTabbedBlock: (blockID) => editBlock(blockID),
  removeTabEntry: (tabsBlockID, tabIndex) => removeTabEntry(tabsBlockID, tabIndex),
  cloneTabbedBlock: (tabsBlockID, tabIndex) => cloneTabbedBlock(tabsBlockID, tabIndex),
})
const isNewBlock = ref(false)
const gridRef = ref(null)
const configuratorTab = ref('block')
const pendingTabBlockIndex = ref(null)

const headerTextVariantOptions = computed(() => [
  { value: 'dark', label: t('block.general.style.default') },
  { value: 'primary', label: t('block.general.style.primary') },
  { value: 'secondary', label: t('block.general.style.secondary') },
  { value: 'success', label: t('block.general.style.success') },
  { value: 'warning', label: t('block.general.style.warning') },
  { value: 'danger', label: t('block.general.style.danger') },
])

const magnifyOptions = computed(() => [
  { value: 'disabled', label: t('block.general.magnifyOptions.disabled') },
  { value: 'modal', label: t('block.general.magnifyOptions.modal') },
  { value: 'fullscreen', label: t('block.general.magnifyOptions.fullscreen') },
])

const isRecordPage = computed(() => !!page.value?.isRecordPage)

const layoutOptions = computed(() =>
  layouts.value.map(l => ({
    value: l.pageLayoutID,
    label: l.meta?.title || l.handle || l.pageLayoutID,
  })),
)

// Custom ID validation: must be at least 2 chars, alphanumeric/underscore/dash, end with letter/number
const customIDInvalid = computed(() => {
  const v = editingBlock.value?.meta?.customID
  if (!v) return false
  return !/^[a-zA-Z0-9_-]{2,}$/.test(v) || !/[a-zA-Z0-9]$/.test(v)
})

// CSS class validation: each class alphanumeric/underscore/dash, must end with letter/number
const customCSSClassInvalid = computed(() => {
  const v = editingBlock.value?.meta?.customCSSClass
  if (!v) return false
  const classes = v.trim().split(/\s+/)
  return classes.some(c => !/^[a-zA-Z0-9_-]+$/.test(c) || !/[a-zA-Z0-9]$/.test(c))
})

// Clamp refresh rate: if between 1-4, force to 5
function clampRefreshRate() {
  const val = editingBlock.value?.options?.refreshRate
  if (val > 0 && val < 5) {
    editingBlock.value.options.refreshRate = 5
  }
}

// Ensure nested objects exist when editingBlock is set
watch(
  editingBlock,
  block => {
    if (!block) return
    if (!block.meta) block.meta = {}
    if (!block.meta.visibility) block.meta.visibility = { expression: '', roles: [] }
    if (block.meta.visibility.expression === undefined) block.meta.visibility.expression = ''
    if (!block.meta.visibility.roles) block.meta.visibility.roles = []
    if (!block.style) block.style = {}
    if (!block.style.variants) block.style.variants = { headerText: '' }
    if (!block.style.wrap) block.style.wrap = { kind: 'card' }
    if (!block.style.border) block.style.border = { enabled: false }
    // Ensure options defaults
    if (!block.options) block.options = {}
    if (
      block.options.magnifyOption === undefined ||
      block.options.magnifyOption === null ||
      block.options.magnifyOption === ''
    )
      block.options.magnifyOption = 'disabled'
  },
  { immediate: true },
)

watch(showAddBlock, visible => {
  if (!visible) {
    pendingTabBlockIndex.value = null
  }
})

watch(showConfigurator, visible => {
  if (!visible && isNewBlock.value) {
    blocks.value.splice(editingBlockIndex.value, 1)
    syncTabbedBlockVisibility()
    gridRef.value?.rebuildLayout()
    isNewBlock.value = false
  }
})

// Label for the block-specific tab
const editingBlockTypeLabel = computed(() => {
  if (!editingBlock.value) return ''
  const bt = availableBlockTypes.value.find(b => b.kind === editingBlock.value.kind)
  return bt?.label || editingBlock.value.kind
})

const hasAutomationTab = computed(() => editingBlock.value?.kind === 'RecordList')

function onSelectionButtonsUpdate(next) {
  if (!editingBlock.value) return
  if (!editingBlock.value.options) editingBlock.value.options = {}
  editingBlock.value.options.selectionButtons = next
}

// Available block types for the "add block" dialog
const availableBlockTypes = computed(() => {
  const isRecordPage = page.value?.isRecordPage

  const recordBlocks = [
    { kind: 'RecordList', label: t('block.recordList.label'), icon: 'pi pi-list' },
    { kind: 'RecordOrganizer', label: t('block.recordOrganizer.label'), icon: 'pi pi-th-large' },
    ...(isRecordPage
      ? [
          { kind: 'Record', label: t('block.record.label'), icon: 'pi pi-objects-column' },
          {
            kind: 'RecordRevisions',
            label: t('block.recordRevisions.label'),
            icon: 'pi pi-history',
          },
        ]
      : []),
  ].sort((a, b) => a.label.localeCompare(b.label))

  const otherBlocks = [
    { kind: 'Automation', label: t('block.automation.label'), icon: 'pi pi-bolt' },
    { kind: 'Calendar', label: t('block.calendar.label'), icon: 'pi pi-calendar' },
    { kind: 'Chart', label: t('block.chart.label'), icon: 'pi pi-chart-bar' },
    { kind: 'Comment', label: t('block.comment.label'), icon: 'pi pi-comments' },
    { kind: 'Content', label: t('block.content.label'), icon: 'pi pi-align-left' },
    { kind: 'File', label: t('block.file.label'), icon: 'pi pi-paperclip' },
    { kind: 'Geometry', label: t('block.geometry.label'), icon: 'pi pi-map' },
    { kind: 'IFrame', label: t('block.iframe.label'), icon: 'pi pi-globe' },
    { kind: 'Metric', label: t('block.metric.label'), icon: 'pi pi-hashtag' },
    { kind: 'Navigation', label: t('block.navigation.label'), icon: 'pi pi-link' },
    { kind: 'Progress', label: t('block.progress.label'), icon: 'pi pi-percentage' },
    { kind: 'Tabs', label: t('block.tabs.label'), icon: 'pi pi-credit-card' },
    { kind: 'ChatbotSessions', label: t('block.chatbotSessions.label'), icon: 'pi pi-inbox' },
  ].sort((a, b) => a.label.localeCompare(b.label))

  return [...recordBlocks, { kind: 'divider' }, ...otherBlocks]
})

const addableBlockTypes = computed(() =>
  pendingTabBlockIndex.value !== null
    ? availableBlockTypes.value.filter(({ kind }) => kind !== 'Tabs' && kind !== 'divider')
    : availableBlockTypes.value,
)

// Configurator registry
const configurators = {
  Chart: markRaw(ChartConfigurator),
  Content: markRaw(ContentConfigurator),
  RecordList: markRaw(RecordListConfigurator),
  Record: markRaw(RecordConfigurator),
  Metric: markRaw(MetricConfigurator),
  IFrame: markRaw(IFrameConfigurator),
  File: markRaw(FileConfigurator),
  Calendar: markRaw(CalendarConfigurator),
  Comment: markRaw(CommentConfigurator),
  Automation: markRaw(AutomationConfigurator),
  Navigation: markRaw(NavigationConfigurator),
  Tabs: markRaw(TabsConfigurator),
  Progress: markRaw(ProgressConfigurator),
  RecordOrganizer: markRaw(RecordOrganizerConfigurator),
  RecordRevisions: markRaw(RecordRevisionsConfigurator),
  Geometry: markRaw(GeometryConfigurator),
  ChatbotSessions: markRaw(ChatbotConfigurator),
}

const blockConfigurator = computed(() => {
  if (!editingBlock.value) return null
  return configurators[editingBlock.value.kind] || null
})

// Unique ID for each block — blockID '0' is NoID (unsaved), so fall back to tempID
function getBlockId(block) {
  const bid = block.blockID
  if (bid && bid !== '0') return bid
  return block.meta?.tempID || ''
}

function syncTabbedBlockVisibility() {
  const tabbedBlockIds = new Set()

  for (const block of blocks.value) {
    if (block.kind !== 'Tabs') continue

    for (const tab of block.options?.tabs || []) {
      if (tab.blockID) {
        tabbedBlockIds.add(String(tab.blockID))
      }
    }
  }

  for (const block of blocks.value) {
    if (block.kind === 'Tabs') continue

    const blockID = String(getBlockId(block))
    if (!blockID) continue

    block.meta = {
      ...(block.meta || {}),
      hidden: tabbedBlockIds.has(blockID),
    }
  }
}

function commitEditingBlock() {
  if (editingBlockIndex.value < 0 || !editingBlock.value) return null

  const original = blocks.value[editingBlockIndex.value]
  if (!original) return null

  if (editingBlock.value.options?.magnifyOption === 'disabled') {
    editingBlock.value.options.magnifyOption = ''
  }

  const updated = compose.PageBlockMaker({
    ...editingBlock.value,
    xywh: original.xywh,
  })

  blocks.value.splice(editingBlockIndex.value, 1, updated)
  syncTabbedBlockVisibility()
  gridRef.value?.rebuildLayout()
  editingBlock.value = JSON.parse(JSON.stringify(updated))

  return updated
}

function addBlock(kind) {
  try {
    const block = compose.PageBlockMaker({ kind })
    const maxY = blocks.value.reduce((max, b) => {
      const [, y, , h] = b.xywh || [0, 0, 24, 18]
      return Math.max(max, y + h)
    }, 0)
    block.xywh = [0, maxY, 24, 18]

    if (pendingTabBlockIndex.value !== null && editingBlock.value?.kind === 'Tabs') {
      block.meta = { ...(block.meta || {}), hidden: true }
      blocks.value.push(block)
      const newBlockId = getBlockId(block)

      const items = [...(editingBlock.value.options?.tabs || [])]
      if (items[pendingTabBlockIndex.value]) {
        items[pendingTabBlockIndex.value] = {
          lazy: true,
          title: '',
          ...items[pendingTabBlockIndex.value],
          blockID: newBlockId,
        }

        editingBlock.value = {
          ...editingBlock.value,
          options: { ...editingBlock.value.options, tabs: items },
        }

        commitEditingBlock()
      }

      pendingTabBlockIndex.value = null
      showAddBlock.value = false
      editBlock(newBlockId)
      return
    }

    blocks.value.push(block)
    syncTabbedBlockVisibility()
    gridRef.value?.rebuildLayout()
    showAddBlock.value = false

    // Open the editor for the newly added block
    const index = blocks.value.length - 1
    editingBlock.value = JSON.parse(JSON.stringify(block))
    editingBlockIndex.value = index
    configuratorTab.value = 'block'
    isNewBlock.value = true
    showConfigurator.value = true
  } catch (e) {
    console.error('Failed to create block:', e)
  }
}

function cloneBlock(blockId) {
  const block = blocks.value.find(b => String(getBlockId(b)) === blockId)
  if (!block) return

  try {
    const clonedRaw = JSON.parse(JSON.stringify(block))
    clonedRaw.blockID = '0'
    if (clonedRaw.title) {
      clonedRaw.title = t('page.copyOf', { title: clonedRaw.title })
    }

    // Reset tempID
    if (clonedRaw.meta) {
      clonedRaw.meta.tempID = ''
    }

    const cloned = compose.PageBlockMaker(clonedRaw)
    // Offset the clone slightly
    const [, y, w, h] = cloned.xywh || [0, 0, 24, 18]
    cloned.xywh = [0, y + h, w, h]
    blocks.value.push(cloned)
    gridRef.value?.rebuildLayout()
  } catch (e) {
    console.error('Failed to clone block:', e)
  }
}

function editBlock(blockId) {
  const block = blocks.value.find(b => String(getBlockId(b)) === blockId)
  if (!block) return

  const index = blocks.value.findIndex(b => String(getBlockId(b)) === blockId)

  // Deep clone the block for editing
  editingBlock.value = JSON.parse(JSON.stringify(block))
  editingBlockIndex.value = index
  configuratorTab.value = 'block'
  showConfigurator.value = true
}

function onCreateTabBlock(index) {
  pendingTabBlockIndex.value = index
  showAddBlock.value = true
}

function onEditTabBlock(blockId) {
  commitEditingBlock()
  editBlock(blockId)
}

function replaceBlockOptions(blockIdx, optionsPatch) {
  const original = blocks.value[blockIdx]
  if (!original) return null
  const draft = JSON.parse(JSON.stringify(original))
  draft.options = { ...(draft.options || {}), ...optionsPatch }
  const updated = compose.PageBlockMaker(draft)
  blocks.value.splice(blockIdx, 1, updated)
  return updated
}

function removeTabEntry(tabsBlockID, tabIndex) {
  const idx = blocks.value.findIndex(b => String(getBlockId(b)) === String(tabsBlockID))
  if (idx < 0 || blocks.value[idx].kind !== 'Tabs') return

  const tabs = [...(blocks.value[idx].options?.tabs || [])]
  if (tabIndex < 0 || tabIndex >= tabs.length) return
  tabs.splice(tabIndex, 1)

  replaceBlockOptions(idx, { tabs })
  syncTabbedBlockVisibility()
  gridRef.value?.rebuildLayout()
}

function cloneTabbedBlock(tabsBlockID, tabIndex) {
  const idx = blocks.value.findIndex(b => String(getBlockId(b)) === String(tabsBlockID))
  if (idx < 0 || blocks.value[idx].kind !== 'Tabs') return

  const tabs = [...(blocks.value[idx].options?.tabs || [])]
  const tab = tabs[tabIndex]
  if (!tab?.blockID) return

  const source = blocks.value.find(b => String(getBlockId(b)) === String(tab.blockID))
  if (!source) return

  try {
    const clonedRaw = JSON.parse(JSON.stringify(source))
    clonedRaw.blockID = '0'
    if (clonedRaw.meta) clonedRaw.meta.tempID = ''
    if (clonedRaw.title) clonedRaw.title = t('page.copyOf', { title: clonedRaw.title })
    clonedRaw.meta = { ...(clonedRaw.meta || {}), hidden: true }

    const cloned = compose.PageBlockMaker(clonedRaw)
    blocks.value.push(cloned)
    const clonedID = getBlockId(cloned)

    tabs.splice(tabIndex + 1, 0, {
      lazy: tab.lazy !== false,
      title: tab.title ? t('page.copyOf', { title: tab.title }) : '',
      blockID: clonedID,
    })

    replaceBlockOptions(idx, { tabs })
    syncTabbedBlockVisibility()
    gridRef.value?.rebuildLayout()
  } catch (e) {
    console.error('Failed to clone tabbed block:', e)
  }
}

function deleteBlock(blockId) {
  const index = blocks.value.findIndex(b => String(getBlockId(b)) === blockId)
  if (index > -1) {
    blocks.value.splice(index, 1)
    syncTabbedBlockVisibility()
    gridRef.value?.rebuildLayout()
  }
}

function openBlockTranslation(block) {
  const { namespaceID, pageID } = page.value
  const blockID = block.blockID
  const res = `compose:page/${namespaceID}/${pageID}`
  translatorStore.open({
    resource: res,
    titles: {
      [res]: t('translator.resources.page.block.title', { title: block.title || blockID }),
    },
    fetcher: () =>
      $ComposeAPI
        .pageListTranslations({ namespaceID, pageID })
        .then((set) => set.filter((tr) => tr.key.startsWith(`pageBlock.${blockID}.`))),
    updater: async (changes) => {
      await $ComposeAPI.pageUpdateTranslations({ namespaceID, pageID, translations: changes })
      const fresh = await $ComposeAPI.pageListTranslations({ namespaceID, pageID })
      const updatedPage = JSON.parse(JSON.stringify(page.value))
      const updatedLayouts = JSON.parse(JSON.stringify(layouts.value))
      applyPageTranslations(updatedPage, updatedLayouts, fresh, currentLanguage.value)
      page.value = updatedPage
      layouts.value = updatedLayouts
    },
  })
}

function saveBlockConfig() {
  isNewBlock.value = false
  commitEditingBlock()
  showConfigurator.value = false
  editingBlock.value = null
  editingBlockIndex.value = -1
  pendingTabBlockIndex.value = null
}

function setLayout(layoutID) {
  const layout = layouts.value.find(l => l.pageLayoutID === layoutID)
  if (!layout) return

  pageLayout.value = layout

  // Update each block's xywh from the new layout (preserve all other block data)
  for (const block of blocks.value) {
    const blockID = getBlockId(block)
    const layoutBlock = layout.blocks?.find(lb => lb.blockID === blockID)
    if (layoutBlock?.xywh) {
      block.xywh = [...layoutBlock.xywh]
    }
  }
  gridRef.value?.rebuildLayout()
}

async function handleDeleteLayout() {
  if (!pageLayout.value) return

  showDeleteLayoutConfirm.value = false
  saving.value = true

  try {
    await pageLayoutStore.delete(toRaw(pageLayout.value))
    layouts.value = pageLayoutStore.getByPageID(page.value.pageID)

    if (layouts.value.length > 0) {
      pageLayout.value = layouts.value[0]
      setLayout(layouts.value[0].pageLayoutID)
    } else {
      pageLayout.value = null
    }

    $toast.toastSuccess(t('notification.page.page-layout.delete.success'))
  } catch (e) {
    console.error('Failed to delete layout:', e)
    $toast.toastDanger(t('notification.page.page-layout.delete.failed'))
  } finally {
    saving.value = false
  }
}

async function handleSaveAsCopy() {
  if (!pageLayout.value || !page.value) return

  saving.value = true

  try {
    const currentTitle = pageLayout.value.meta?.title || ''
    const layoutBlocks = blocks.value.map(b => ({
      blockID: getBlockId(b),
      xywh: b.xywh,
    }))

    const copy = {
      ...toRaw(pageLayout.value),
      pageLayoutID: '0',
      namespaceID: page.value.namespaceID || props.namespace?.namespaceID,
      pageID: page.value.pageID,
      handle: '',
      meta: {
        ...pageLayout.value.meta,
        title: currentTitle ? t('page.copyOf', { title: currentTitle }) : '',
      },
      blocks: layoutBlocks,
    }

    const created = await pageLayoutStore.create(copy)
    layouts.value = pageLayoutStore.getByPageID(page.value.pageID)
    setLayout(created.pageLayoutID)

    $toast.toastSuccess(t('notification.page.page-layout.saveAsCopy.success'))
  } catch (e) {
    console.error('Failed to save layout as copy:', e)
    $toast.toastDanger(t('notification.page.page-layout.saveAsCopy.failed'))
  } finally {
    saving.value = false
  }
}

function validateRequiredFields() {
  if (!page.value?.isRecordPage) return true

  const mod = moduleStore.getByID(page.value.moduleID)
  if (!mod?.fields?.length) return true

  const required = new Set(mod.fields.filter(f => f.isRequired).map(f => f.name))
  if (required.size === 0) return true

  for (const block of blocks.value) {
    if (block.kind !== 'Record') continue
    const fields = block.options?.fields || []
    // No field filter means all fields are shown
    if (!fields.length) return true
    for (const f of fields) {
      required.delete(f.name)
    }
  }

  return required.size === 0
}

async function loadPage() {
  const pageID = route.params.pageID
  if (!pageID) return

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

    // Load all layouts for this page
    layouts.value = pageLayoutStore.getByPageID(pageID)
    // Keep current layout if still in the list, otherwise pick the first
    const currentID = pageLayout.value?.pageLayoutID
    const keepLayout = currentID && layouts.value.find(l => l.pageLayoutID === currentID)
    pageLayout.value = keepLayout || layouts.value[0] || null

    // Initialize blocks — merge layout positions with page blocks
    if (page.value?.blocks) {
      blocks.value = page.value.blocks.map(b => {
        const layoutBlock = pageLayout.value?.blocks?.find(lb => lb.blockID === b.blockID)
        const merged = compose.PageBlockMaker({
          ...b,
          xywh: layoutBlock?.xywh || b.xywh || [0, 0, 48, 15],
        })
        return merged
      })
    }
  } catch (e) {
    console.error('Failed to load page:', e)
    $toast.toastDanger(t('notification.page.loadFailed'))
    router.push({ name: 'admin.pages' })
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!page.value) return

  // Warn if required module fields are not covered by any Record block
  if (!validateRequiredFields()) {
    $toast.toastWarning(t('notification.page.requiredFields.missing'))
  }

  saving.value = true

  try {
    const rawBlocks = blocks.value.map(b => toRaw(b))

    // Update page with blocks
    let updatedPage = await pageStore.update({
      ...toRaw(page.value),
      blocks: rawBlocks,
    })

    const savedBlocks = updatedPage?.blocks || []
    const blockIdMap = new Map(
      rawBlocks.map((block, index) => [
        String(getBlockId(block)),
        String(savedBlocks[index]?.blockID || block.blockID),
      ]),
    )

    const remappedBlocks = savedBlocks.map(block => {
      if (block.kind !== 'Tabs') return block

      let changed = false
      const tabs = (block.options?.tabs || []).map(tab => {
        const mappedBlockID = blockIdMap.get(String(tab.blockID))
        if (mappedBlockID && mappedBlockID !== String(tab.blockID)) {
          changed = true
          return { ...tab, blockID: mappedBlockID }
        }

        return tab
      })

      return changed ? { ...toRaw(block), options: { ...block.options, tabs } } : block
    })

    if (remappedBlocks.some((block, index) => block !== savedBlocks[index])) {
      updatedPage = await pageStore.update({
        ...toRaw(updatedPage),
        blocks: remappedBlocks,
      })
    }

    // If we have a layout, update it too with block positions
    if (pageLayout.value) {
      const persistedBlocks = updatedPage?.blocks || []
      const layoutBlocks = rawBlocks.map((b, index) => ({
        blockID: persistedBlocks[index]?.blockID || b.blockID,
        xywh: b.xywh,
      }))

      await pageLayoutStore.update({
        ...toRaw(pageLayout.value),
        blocks: layoutBlocks,
      })
    }

    $toast.toastSuccess(t('notification.page.saved'))

    // Reload to get updated blockIDs for new blocks
    await loadPage()
  } catch (e) {
    console.error('Failed to save page:', e)
    $toast.toastDanger(t('notification.page.saveFailed'))
  } finally {
    saving.value = false
  }
}

function goToViewPage() {
  if (page.value) {
    if (page.value.isRecordPage) {
      // Record page — open create record view
      router.push({
        name: 'page.record',
        params: {
          slug: route.params.slug,
          pageID: page.value.pageID,
          recordID: '0',
        },
      })
    } else {
      router.push({
        name: 'page',
        params: { pageID: page.value.pageID },
      })
    }
  }
}

function goToEditPage() {
  if (page.value) {
    router.push({
      name: 'admin.pages.edit',
      params: { pageID: page.value.pageID },
    })
  }
}

function goToModuleEdit() {
  if (page.value?.isRecordPage) {
    router.push({
      name: 'admin.modules.edit',
      params: { moduleID: page.value.moduleID },
    })
  }
}

// Load on mount
watch(
  () => route.params.pageID,
  () => loadPage(),
  { immediate: true },
)
</script>

<style scoped>
.block-toolbox {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 5;
  display: flex;
  justify-content: center;
  padding: 4px;
  border-bottom-left-radius: var(--p-card-border-radius);
  border-top-right-radius: var(--p-card-border-radius);
}
</style>
