<template>
  <Teleport to="#topbar-title" defer>
    <span v-if="page">{{ $t('page.edit.pageBuilder') }}: {{ page.title }}</span>
  </Teleport>

  <Teleport to="#topbar-tools" defer>
    <div v-if="page" class="flex items-center gap-2">
      <Select
        v-if="layouts.length > 1"
        :key="layoutSelectKey"
        :model-value="pageLayout?.pageLayoutID"
        :options="layoutOptions"
        option-label="label"
        option-value="value"
        size="small"
        style="min-width: 200px"
        @update:model-value="onLayoutSelect"
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
          @update:page="onPageTranslated"
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
    <!-- No layout: the builder can't render a page without one. Offer to add it. -->
    <div
      v-if="!pageLayout"
      class="flex-1 flex flex-col items-center justify-center gap-4 p-6 text-center"
    >
      <i class="pi pi-table text-4xl text-muted-color" />
      <p class="text-muted-color">{{ $t('page.build.noLayout') }}</p>
      <Button
        :label="$t('page.build.addLayout')"
        icon="pi pi-plus"
        :loading="saving"
        @click="handleCreateLayout"
      />
    </div>

    <!-- Builder Grid -->
    <div v-else class="flex-1 overflow-auto">
      <Grid ref="gridRef" :blocks="blocks" :namespace="namespace" :page="page" editable>
        <template #item-overlay="{ item, block }">
          <div
            v-tooltip.top="$t('page.tooltip.drag.block')"
            class="block-drag-handle block-drag-tab bg-emphasis"
          >
            <i class="pi pi-bars text-base" />
          </div>
          <div class="block-toolbox bg-emphasis flex items-center">
            <i
              v-if="issuesFor(item.i).length"
              v-tooltip.top="issueSummary(issuesFor(item.i))"
              class="pi pi-exclamation-triangle text-orange-500 text-sm px-2"
            />
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
                v-if="showTranslatorButton && block?.blockID && block?.blockID !== '0'"
                v-tooltip.top="$t('translator.button.tooltip')"
                icon="pi pi-language"
                text
                size="small"
                severity="secondary"
                @click.stop="openBlockTranslation(block)"
              />
              <Button
                :title="$t('page.tooltip.removeFromLayout')"
                icon="pi pi-times"
                text
                size="small"
                severity="secondary"
                @click="confirmRemoveBlock(item.i)"
              />
            </ButtonGroup>
          </div>
        </template>
      </Grid>
    </div>

    <CEditorActions
      :back-to="true"
      @back="goBack({ name: 'admin.pages', params: { slug: route.params.slug } })"
    >
      <template #center>
        <Button
          v-if="pageLayout"
          :label="$t('page.build.addBlock')"
          icon="pi pi-plus"
          severity="secondary"
          @click="showAddBlock = true"
        />
      </template>
      <Button
        v-if="pageLayout"
        :label="$t('general.label.delete')"
        icon="pi pi-trash"
        severity="danger"
        :loading="saving"
        @click="showDeleteLayoutConfirm = true"
      />
      <Button
        v-if="pageLayout"
        :label="$t('page.build.saveAsCopy')"
        icon="pi pi-copy"
        severity="secondary"
        :loading="saving"
        @click="handleSaveAsCopy"
      />
      <Button
        v-if="pageLayout"
        :label="$t('general.label.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="handleSave"
      />
    </CEditorActions>
  </div>

  <!-- Add Block Dialog -->
  <Dialog
    v-model:visible="showAddBlock"
    :header="$t('page.build.selectBlockTitle')"
    modal
    :style="{ width: '600px' }"
  >
    <div class="flex flex-col gap-5">
      <!-- New block -->
      <div>
        <h4 v-if="orphanBlocks.length" class="text-sm font-semibold text-muted-color mb-2">
          {{ $t('page.build.newBlock') }}
        </h4>
        <div class="grid grid-cols-3 md:grid-cols-4 gap-3">
          <template v-for="bt in addableBlockTypes" :key="bt.kind">
            <Divider v-if="bt.kind === 'divider'" class="col-span-3 md:col-span-4 my-0" />
            <div
              v-else
              class="flex flex-col items-center justify-center gap-1.5 px-2 py-5 border border-surface rounded cursor-pointer hover:bg-highlight transition-colors"
              @click="addBlock(bt.kind)"
            >
              <i :class="bt.icon" class="text-xl text-primary" />
              <span class="text-xs font-medium text-center">{{ bt.label }}</span>
            </div>
          </template>
        </div>
      </div>

      <!-- Existing page blocks not yet placed in this layout (as a list) -->
      <div v-if="orphanBlocks.length">
        <h4 class="text-sm font-semibold text-muted-color mb-2">
          {{ $t('page.build.existingBlocks') }}
        </h4>
        <div class="flex flex-col gap-1">
          <div
            v-for="ob in orphanBlocks"
            :key="ob.blockID"
            class="group flex items-center gap-2 pr-1 border border-surface rounded hover:bg-highlight transition-colors"
          >
            <button
              type="button"
              class="flex items-center gap-3 flex-1 min-w-0 p-2 text-left cursor-pointer"
              @click="addExistingBlock(ob)"
            >
              <i
                :class="blockTypeMeta[ob.kind]?.icon || 'pi pi-box'"
                class="text-lg text-primary shrink-0"
              />
              <span class="text-sm font-medium truncate">
                {{ ob.title || blockTypeMeta[ob.kind]?.label || ob.kind }}
              </span>
            </button>
            <Button
              icon="pi pi-trash"
              text
              rounded
              size="small"
              severity="danger"
              class="shrink-0 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity"
              :title="$t('general.label.delete')"
              @click="confirmDeletePageBlock(ob)"
            />
          </div>
        </div>
      </div>
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
          text
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
    :header="configuratorHeader"
    modal
    position="top"
    :style="{ width: '50vw' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :pt="{
      content: { class: 'p-0 flex flex-col !overflow-hidden' },
      footer: { class: 'border-t border-surface p-3' },
    }"
  >
    <Message
      v-if="editingBlockIssues.length"
      severity="warn"
      size="small"
      variant="simple"
      class="mx-3 mt-3 shrink-0"
    >
      {{ $t('block.issue.summary', { needs: issueSummary(editingBlockIssues) }) }}
    </Message>

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
            <CFormGroup :label="$t('block.general.titleLabel')">
              <InputText
                v-model="editingBlock.title"
                :placeholder="$t('block.general.titlePlaceholder')"
                class="w-full"
              />
            </CFormGroup>

            <!-- Description -->
            <CFormGroup :label="$t('block.general.descriptionLabel')">
              <Textarea
                v-model="editingBlock.description"
                :placeholder="$t('block.general.descriptionPlaceholder')"
                rows="2"
                class="w-full"
              />
            </CFormGroup>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <!-- Custom ID -->
              <CFormGroup :label="$t('block.general.customID.label')">
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
              </CFormGroup>

              <!-- Custom CSS Class -->
              <CFormGroup :label="$t('block.general.customCSSClass.label')">
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
              </CFormGroup>
            </div>

            <Divider />

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <!-- Header style -->
              <CFormGroup :label="$t('block.general.headerStyle')">
                <Select
                  v-model="editingBlock.style.variants.headerText"
                  :options="headerTextVariantOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                />
              </CFormGroup>

              <!-- Magnify option -->
              <CFormGroup :label="$t('block.general.magnifyLabel')">
                <Select
                  v-model="editingBlock.options.magnifyOption"
                  :options="magnifyOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                />
              </CFormGroup>
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
                <CFormGroup
                  :label="$t('block.general.refresh.auto')"
                  :description="$t('block.general.refresh.description')"
                >
                  <div class="flex items-center gap-2">
                    <InputNumber
                      v-model="editingBlock.options.refreshRate"
                      :min="0"
                      suffix=" s"
                      class="w-full"
                      @blur="clampRefreshRate"
                    />
                  </div>
                </CFormGroup>
                <div class="flex gap-2">
                  <CInputSwitch
                    v-model="editingBlock.options.showRefresh"
                    :label="$t('block.general.refresh.show')"
                  />
                </div>
              </div>
            </template>

            <Divider />

            <Fieldset>
              <template #legend>
                <span class="flex items-center gap-1">
                  {{ $t('block.general.visibility.label') }}
                  <i
                    class="pi pi-exclamation-triangle text-orange-500 text-xs"
                    v-tooltip="$t('block.general.visibility.tooltip.performance.condition')"
                  />
                </span>
              </template>
              <div class="flex flex-col gap-3">
                <CFormGroup
                  :label="$t('block.general.visibility.condition.label')"
                  :description="visibilityConditionDescription"
                >
                  <CInputExpression
                    ref="visibilityInput"
                    v-model="editingBlock.meta.visibility.expression"
                    dialect="expr"
                    :scope="exprScope"
                    :placeholder="$t('block.general.visibility.condition.placeholder')"
                  />
                </CFormGroup>

                <CFormGroup :label="$t('block.general.visibility.roles.label')">
                  <CInputRole
                    v-model="editingBlock.meta.visibility.roles"
                    multiple
                    :placeholder="$t('block.general.visibility.roles.placeholder')"
                    class="w-full"
                  />
                </CFormGroup>
              </div>
            </Fieldset>
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
            :scope="scope"
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
          text
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
import { useConfirm } from 'primevue/useconfirm'
import { useRoute, useRouter } from 'vue-router'
import { compose } from '@planetcrust/human-js'
import { usePageStore } from '@planetcrust/human-vue'
import { usePageLayoutStore } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { useHistoryBack } from '@planetcrust/human-vue'
import { useDraftGuard } from '@planetcrust/human-vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'

// Block configurators — lazy imported
import ChartConfigurator from '@/sections/compose/components/PageBlocks/Configurators/ChartConfigurator.vue'
import ContentConfigurator from '@/sections/compose/components/PageBlocks/Configurators/ContentConfigurator.vue'
import RecordListConfigurator from '@/sections/compose/components/PageBlocks/Configurators/RecordListConfigurator.vue'
import RecordConfigurator from '@/sections/compose/components/PageBlocks/Configurators/RecordConfigurator.vue'
import MetricConfigurator from '@/sections/compose/components/PageBlocks/Configurators/MetricConfigurator.vue'
import IFrameConfigurator from '@/sections/compose/components/PageBlocks/Configurators/IFrameConfigurator.vue'
import FileConfigurator from '@/sections/compose/components/PageBlocks/Configurators/FileConfigurator.vue'
import CalendarConfigurator from '@/sections/compose/components/PageBlocks/Configurators/CalendarConfigurator.vue'
import CommentConfigurator from '@/sections/compose/components/PageBlocks/Configurators/CommentConfigurator.vue'
import AutomationConfigurator from '@/sections/compose/components/PageBlocks/Configurators/AutomationConfigurator.vue'
import NavigationConfigurator from '@/sections/compose/components/PageBlocks/Configurators/NavigationConfigurator.vue'
import TabsConfigurator from '@/sections/compose/components/PageBlocks/Configurators/TabsConfigurator.vue'
import ProgressConfigurator from '@/sections/compose/components/PageBlocks/Configurators/ProgressConfigurator.vue'
import RecordOrganizerConfigurator from '@/sections/compose/components/PageBlocks/Configurators/RecordOrganizerConfigurator.vue'
import RecordRevisionsConfigurator from '@/sections/compose/components/PageBlocks/Configurators/RecordRevisionsConfigurator.vue'
import GeometryConfigurator from '@/sections/compose/components/PageBlocks/Configurators/GeometryConfigurator.vue'
import ChatbotInboxConfigurator from '@/sections/compose/components/PageBlocks/Configurators/ChatbotInboxConfigurator.vue'
import AgentChatConfigurator from '@/sections/compose/components/PageBlocks/Configurators/AgentChatConfigurator.vue'
import AutomationButtonsEditor from '@/sections/compose/components/PageBlocks/Shared/AutomationButtonsEditor.vue'
import PageTranslator from '@/sections/compose/components/Admin/Page/PageTranslator.vue'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'
import { useTranslatorStore } from '@/sections/compose/stores/translator'
import { applyPageTranslations } from '@/sections/compose/lib/resource-translations'

const { t } = useI18n()
const confirm = useConfirm()
const route = useRoute()
const router = useRouter()
const goBack = useHistoryBack()
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
  editTabbedBlock: blockID => editBlock(blockID),
  removeTabEntry: (tabsBlockID, tabIndex) => removeTabEntry(tabsBlockID, tabIndex),
  cloneTabbedBlock: (tabsBlockID, tabIndex) => cloneTabbedBlock(tabsBlockID, tabIndex),
})
const isNewBlock = ref(false)
const gridRef = ref(null)
const configuratorTab = ref('block')
const pendingTabBlockIndex = ref(null)

// Orphan-block deletions are staged like every other edit and applied on Save;
// leaving without saving discards them.
const stagedBlockDeletes = ref(new Set())

// Everything the builder does is staged until Save: block config, adds, removes,
// orphan deletes and drag/resize positions all live in memory until then.
const { isDirty, capture: captureBlocksBaseline } = useDraftGuard({
  draft: blocks,
  busy: saving,
  extra: () => [...stagedBlockDeletes.value],
})

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

const visibilityInput = ref(null)
const { scope, exprScope } = useExpressionScope({
  page: computed(() => page.value),
  hasRecord: isRecordPage,
})

const visibilityConditionDescription = computed(() =>
  isRecordPage.value
    ? t('block.general.visibility.condition.description.record-page', [
        'record.values.fieldName',
        'user.(userID/email...)',
        'screen.(width/height)',
        'isView/isCreate/isEdit',
        'user.userID == record.createdBy',
        'screen.width < 1024',
      ])
    : t('block.general.visibility.condition.description.non-record-page', [
        'user.(userID/email...)',
        'screen.(width/height)',
        'user.email == "test@mail.com"',
        'screen.width < 1024',
      ]),
)

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
  // Closing the configurator on a staged new block (never saved) discards it —
  // it was never added to the page, so there's nothing to remove.
  if (!visible && isNewBlock.value) {
    isNewBlock.value = false
  }
})

// Label for the block-specific tab
const editingBlockTypeLabel = computed(() => {
  if (!editingBlock.value) return ''
  const bt = availableBlockTypes.value.find(b => b.kind === editingBlock.value.kind)
  return bt?.label || editingBlock.value.kind
})

// The same dialog stages a new block and edits a placed one, so the header
// names which of the two is happening and to what.
const configuratorHeader = computed(() =>
  t(isNewBlock.value ? 'page.build.blockDialog.add' : 'page.build.blockDialog.edit', {
    block: editingBlockTypeLabel.value,
  }),
)

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
    { kind: 'ChatbotInbox', label: t('block.chatbotInbox.label'), icon: 'pi pi-headphones' },
    { kind: 'AgentChat', label: t('block.agentChat.label'), icon: 'pi pi-sparkles' },
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
  ChatbotInbox: markRaw(ChatbotInboxConfigurator),
  AgentChat: markRaw(AgentChatConfigurator),
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

// First save of a staged new block — add it to the page/layout now.
function addEditingAsNewBlock() {
  if (!editingBlock.value) return
  if (editingBlock.value.options?.magnifyOption === 'disabled') {
    editingBlock.value.options.magnifyOption = ''
  }
  const block = compose.PageBlockMaker(JSON.parse(JSON.stringify(toRaw(editingBlock.value))))
  blocks.value.push(block)
  syncTabbedBlockVisibility()
  gridRef.value?.rebuildLayout()
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
  editingBlock.value = compose.PageBlockMaker(JSON.parse(JSON.stringify(updated)))

  return updated
}

// Place an existing page block (orphan) into the current layout. It already
// lives on the page, so this only adds it to the working/layout set.
function addExistingBlock(pageBlock) {
  const maxY = blocks.value.reduce((max, b) => {
    const [, y, , h] = b.xywh || [0, 0, 24, 18]
    return Math.max(max, y + h)
  }, 0)
  const block = compose.PageBlockMaker({
    ...JSON.parse(JSON.stringify(toRaw(pageBlock))),
    // Same default footprint as a newly added block.
    xywh: [0, maxY, 24, 18],
  })
  blocks.value.push(block)
  syncTabbedBlockVisibility()
  gridRef.value?.rebuildLayout()
  showAddBlock.value = false
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

    showAddBlock.value = false

    // Stage the new block in the editor; it's added to the page only on the
    // first save of its config. Cancelling discards it (nothing was added).
    editingBlock.value = compose.PageBlockMaker(JSON.parse(JSON.stringify(block)))
    editingBlockIndex.value = -1
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

  // Deep clone the block for editing, preserving the class instance so prototype methods (e.g. fetch) survive
  editingBlock.value = compose.PageBlockMaker(JSON.parse(JSON.stringify(block)))
  editingBlockIndex.value = index
  configuratorTab.value = 'block'
  isNewBlock.value = false
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

// Remove from the current layout (the block stays on the page as an orphan).
function confirmRemoveBlock(blockId) {
  confirm.require({
    header: t('page.build.removeBlock.header'),
    message: t('page.build.removeBlock.message'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: { label: t('page.tooltip.removeFromLayout'), size: 'small' },
    accept: () => deleteBlock(blockId),
  })
}

function confirmDeletePageBlock(pageBlock) {
  confirm.require({
    header: t('page.build.deleteBlock.header'),
    message: t('page.build.deleteBlock.message'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: { label: t('general.label.delete'), severity: 'danger', size: 'small' },
    accept: () => {
      stagedBlockDeletes.value = new Set([...stagedBlockDeletes.value, String(pageBlock.blockID)])
    },
  })
}

function syncBlockTranslations(updatedPageBlocks) {
  for (const updatedBlock of updatedPageBlocks || []) {
    const block = blocks.value.find(b => String(getBlockId(b)) === String(updatedBlock.blockID))
    if (block) {
      block.title = updatedBlock.title
      block.description = updatedBlock.description
    }
  }
}

function onPageTranslated(updatedPage) {
  page.value = updatedPage
  syncBlockTranslations(updatedPage.blocks)
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
        .then(set => set.filter(tr => tr.key.startsWith(`pageBlock.${blockID}.`))),
    updater: async changes => {
      await $ComposeAPI.pageUpdateTranslations({ namespaceID, pageID, translations: changes })
      const fresh = await $ComposeAPI.pageListTranslations({ namespaceID, pageID })
      const updatedPage = JSON.parse(JSON.stringify(page.value))
      const updatedLayouts = JSON.parse(JSON.stringify(layouts.value))
      applyPageTranslations(updatedPage, updatedLayouts, fresh, currentLanguage.value)
      page.value = updatedPage
      layouts.value = updatedLayouts
      syncBlockTranslations(updatedPage.blocks)
    },
  })
}

function saveBlockConfig() {
  // The general tab already marks these red; committing them anyway would
  // persist an ID or class the page cannot use.
  if (customIDInvalid.value || customCSSClassInvalid.value) {
    configuratorTab.value = 'general'
    $toast.toastDanger(t('notification.page.blockSaveFailedInvalid'))
    return
  }

  if (isNewBlock.value) {
    addEditingAsNewBlock()
  } else {
    commitEditingBlock()
  }
  isNewBlock.value = false
  showConfigurator.value = false
  editingBlock.value = null
  editingBlockIndex.value = -1
  pendingTabBlockIndex.value = null
}

// The builder edits one layout at a time. A block is defined on the page; the
// layout positions a subset of those blocks. buildLayoutBlocks resolves the
// current layout's block list (with positions) against the page's definitions —
// mirroring the view, so the builder shows exactly what the view renders.
function buildLayoutBlocks(pg, layout) {
  if (!pg?.blocks?.length || !layout?.blocks?.length) return []
  return layout.blocks
    .map(lb => {
      const pb = pg.blocks.find(b => String(b.blockID) === String(lb.blockID))
      return pb
        ? compose.PageBlockMaker({ ...pb, xywh: lb.xywh || pb.xywh || [0, 0, 48, 15] })
        : null
    })
    .filter(Boolean)
}

// Page blocks not placed in the current layout — offered in the Add block dialog
// so the user can re-add an existing (orphaned) block instead of creating one.
const orphanBlocks = computed(() => {
  if (!page.value?.blocks?.length) return []
  const placed = new Set(blocks.value.map(b => String(getBlockId(b))))
  return page.value.blocks.filter(
    b =>
      b.blockID &&
      !placed.has(String(b.blockID)) &&
      !stagedBlockDeletes.value.has(String(b.blockID)),
  )
})

// kind -> { icon, label }, for rendering orphan tiles in the Add block dialog.
const blockTypeMeta = computed(() => {
  const map = {}
  for (const bt of availableBlockTypes.value) {
    if (bt.kind !== 'divider') map[bt.kind] = bt
  }
  return map
})

function setLayout(layoutID) {
  const layout = layouts.value.find(l => l.pageLayoutID === layoutID)
  if (!layout) return

  pageLayout.value = layout
  // Each layout positions its own subset of the page's blocks.
  blocks.value = buildLayoutBlocks(page.value, layout)
  stagedBlockDeletes.value = new Set()
  captureBlocksBaseline()
  gridRef.value?.rebuildLayout()
}

// Select keeps its own copy of the selection, so a switch the user backs out of
// leaves it displaying a layout the builder is not editing. Remounting it puts
// the label back on whatever pageLayout actually holds.
const layoutSelectKey = ref(0)

// The working set is rebuilt from the saved page, so switching layout discards
// the same staged edits that leaving the builder would.
function onLayoutSelect(layoutID) {
  if (!isDirty()) {
    setLayout(layoutID)
    return
  }

  confirm.require({
    // Both paths out of the dialog: the reject button, and dismissing it.
    reject: () => layoutSelectKey.value++,
    onHide: () => layoutSelectKey.value++,
    header: t('page.build.switchLayout.header'),
    message: t('page.build.switchLayout.message'),
    icon: 'pi pi-exclamation-triangle',
    rejectProps: {
      label: t('general.label.cancel'),
      severity: 'secondary',
      text: true,
      size: 'small',
    },
    acceptProps: { label: t('page.build.switchLayout.confirm'), size: 'small' },
    accept: () => setLayout(layoutID),
  })
}

async function handleCreateLayout() {
  if (!page.value) return

  saving.value = true
  try {
    // Empty 'primary' layout. PageLayout type carries the default config (record
    // toolbar buttons enabled); a plain object would persist them disabled.
    const created = await pageLayoutStore.create(
      new compose.PageLayout({
        namespaceID: page.value.namespaceID || props.namespace?.namespaceID,
        pageID: page.value.pageID,
        handle: 'primary',
        meta: { title: page.value.title },
        blocks: [],
      }),
    )
    layouts.value = pageLayoutStore.getByPageID(page.value.pageID)
    setLayout(created.pageLayoutID)
    $toast.toastSuccess(t('notification.page.page-layout.create.success'))
  } catch (e) {
    console.error('Failed to create layout:', e)
    $toast.toastDanger(t('notification.page.page-layout.create.failed'))
  } finally {
    saving.value = false
  }
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
      // No layouts left — the builder shows its "Add layout" empty state.
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

// What each block still needs, keyed by block id. The block classes own the
// rule; this is the one read of it, so the grid badge, the configurator and
// the save refusal cannot disagree.
const blockIssues = computed(() => {
  const byBlock = new Map()
  for (const block of blocks.value) {
    const issues = block.validate?.() || []
    if (issues.length) byBlock.set(String(getBlockId(block)), issues)
  }
  return byBlock
})

// The draft is not in `blocks` until it is saved, so the dialog asks it directly.
const editingBlockIssues = computed(() => editingBlock.value?.validate?.() || [])

function issueSummary(issues) {
  return issues.map(i => t(i.labelKey)).join(', ')
}

function issuesFor(blockId) {
  return blockIssues.value.get(String(blockId)) || []
}

function blockLabel(block) {
  return (
    block.title || availableBlockTypes.value.find(b => b.kind === block.kind)?.label || block.kind
  )
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
    // A configured field is a name or a {name} object
    for (const f of fields) {
      required.delete(f.name ?? f)
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
    // Deep link (?layoutID) wins, then the current layout if still in the
    // list, otherwise the first
    const requestedID = route.query.layoutID
    const requested = requestedID && layouts.value.find(l => l.pageLayoutID === requestedID)
    const currentID = pageLayout.value?.pageLayoutID
    const keepLayout = currentID && layouts.value.find(l => l.pageLayoutID === currentID)
    pageLayout.value = requested || keepLayout || layouts.value[0] || null

    // Show only the current layout's blocks (the view does the same); page
    // blocks not in the layout are orphans, addable via the Add block dialog.
    blocks.value = buildLayoutBlocks(page.value, pageLayout.value)
    captureBlocksBaseline()
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

  // A record page whose required module fields are not covered by any Record
  // block would render a form users can't complete — refuse to save it.
  if (!validateRequiredFields()) {
    $toast.toastDanger(t('notification.page.saveFailedRequired'))
    return
  }

  // A block missing what it cannot render without would ship to the live page
  // as empty chrome. Name the blocks rather than just refusing.
  if (blockIssues.value.size) {
    const named = blocks.value
      .filter(b => issuesFor(getBlockId(b)).length)
      .map(
        b =>
          `${blockLabel(b)} — ${issuesFor(getBlockId(b))
            .map(i => t(i.labelKey))
            .join(', ')}`,
      )
      .join('; ')
    $toast.toastDanger(t('notification.page.saveFailedBlockConfig', { blocks: named }))
    return
  }

  saving.value = true

  try {
    // Working set = the blocks placed in the current layout (edited/new/added).
    const working = blocks.value.map(b => toRaw(b))
    const workingIds = new Set(working.map(b => String(getBlockId(b))))

    // Preserve page blocks not in this layout (orphans + blocks owned by other
    // layouts) so saving a layout never deletes them from the page. They go
    // first; the working set is the tail, aligned by index with the response.
    // Staged orphan deletions are applied here: dropped from the page payload.
    const preserved = (page.value.blocks || [])
      .map(b => toRaw(b))
      .filter(
        b => !workingIds.has(String(b.blockID)) && !stagedBlockDeletes.value.has(String(b.blockID)),
      )

    const pageBlocks = [...preserved, ...working]

    let updatedPage = await pageStore.update({
      ...toRaw(page.value),
      blocks: pageBlocks,
    })

    let savedBlocks = updatedPage?.blocks || []
    const blockIdMap = new Map(
      pageBlocks.map((block, index) => [
        String(getBlockId(block)),
        String(savedBlocks[index]?.blockID || getBlockId(block)),
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
      savedBlocks = updatedPage?.blocks || savedBlocks
    }

    // The current layout references only its working blocks (by index against the
    // saved tail, so multiple new blocks get their right ids) with positions.
    if (pageLayout.value) {
      const layoutBlocks = working.map((b, i) => ({
        blockID: savedBlocks[preserved.length + i]?.blockID || String(getBlockId(b)),
        xywh: b.xywh,
      }))

      await pageLayoutStore.update({
        ...toRaw(pageLayout.value),
        blocks: layoutBlocks,
      })
    }

    // Staged block deletions: drop references from every other layout too
    if (stagedBlockDeletes.value.size) {
      for (const layout of layouts.value) {
        if (layout.pageLayoutID === pageLayout.value?.pageLayoutID) continue
        if ((layout.blocks || []).some(lb => stagedBlockDeletes.value.has(String(lb.blockID)))) {
          await pageLayoutStore.update({
            ...toRaw(layout),
            blocks: (layout.blocks || []).filter(
              lb => !stagedBlockDeletes.value.has(String(lb.blockID)),
            ),
          })
        }
      }
      stagedBlockDeletes.value = new Set()
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
/* Positions are relative to .grid-stack-item (the outer gridstack-positioned div).
 * .grid-stack-item-content is inset 6px on every side by gridstack's margin and
 * carries the dashed border. Anchoring the overlay at 6px (= the border-box outer
 * edge of .grid-stack-item-content) makes it land exactly on top of the border. */
.block-toolbox {
  position: absolute;
  bottom: 5px;
  left: 5px;
  z-index: 10;
  display: flex;
  justify-content: center;
  padding: 4px;
  border-top-right-radius: var(--p-card-border-radius);
  border-bottom-left-radius: var(--p-card-border-radius);
}

.block-drag-tab {
  position: absolute;
  top: 6px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  padding: 6px 44px;
  cursor: move;
  color: var(--p-text-muted-color);
  background-clip: padding-box;
  border: 1px solid var(--p-content-border-color);
  border-top: none;
  border-bottom-left-radius: var(--p-card-border-radius);
  border-bottom-right-radius: var(--p-card-border-radius);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.8;
  transition:
    opacity 0.15s ease,
    color 0.15s ease;
}

.block-drag-tab:hover {
  opacity: 1;
  color: var(--p-primary-color);
}
</style>
