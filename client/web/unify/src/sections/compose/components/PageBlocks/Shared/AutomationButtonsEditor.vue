<template>
  <div class="flex flex-col gap-3">
    <CFormGroup :label="$t('block.automation.configuredButtons')">
      <template #actions>
        <Button
          :label="$t('block.automation.addPlaceholderLabel')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          @click="addPlaceholder"
        />
      </template>

      <CFormList
        v-model="buttonsModel"
        draggable
        :empty-message="$t('block.automation.noScripts')"
        :columns="[{ width: '1fr' }, { width: '3rem' }]"
        @change="selectedIndex = -1"
        @reorder="selectedIndex = -1"
      >
        <template #row="{ item, index }">
          <Button
            :label="item.label || '-'"
            :severity="mapVariantSeverity(item.variant)"
            size="small"
            class="pointer-events-none truncate justify-self-start"
          />
          <Button icon="pi pi-pencil" text rounded size="small" @click="selectButton(index)" />
        </template>

        <template #extra="{ item, index }">
          <div
            v-if="selectedIndex === index"
            class="flex flex-col gap-2 border-t border-surface pt-3"
          >
            <div v-if="item.script" class="flex items-center gap-2 min-w-0">
              <Tag :value="$t('block.automation.badge.script')" severity="secondary" />
              <span class="text-sm text-muted-color truncate">{{ item.script }}</span>
            </div>

            <CFormGroup :label="$t('block.automation.buttonLabel')">
              <CInputExpression
                :ref="el => (labelInputs[index] = el)"
                :model-value="item.label"
                dialect="interpolation"
                :scope="scope"
                @update:model-value="updateField(index, 'label', $event)"
              />
              <CExpressionHint :scope="scope" @insert="labelInputs[index]?.insert($event)" />
            </CFormGroup>

            <CFormGroup :label="$t('block.automation.buttonVariant')">
              <Select
                :model-value="item.variant"
                :options="variantOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                @update:model-value="updateField(index, 'variant', $event)"
              >
                <template #value="{ value, placeholder }">
                  <Button
                    v-if="value"
                    :label="variantLabel(value)"
                    :severity="mapVariantSeverity(value)"
                    size="small"
                    class="pointer-events-none !py-0.5 !px-2 !text-xs"
                  />
                  <span v-else>{{ placeholder }}</span>
                </template>
                <template #option="{ option }">
                  <Button
                    :label="option.label"
                    :severity="mapVariantSeverity(option.value)"
                    size="small"
                    class="pointer-events-none !py-0.5 !px-2 !text-xs"
                  />
                </template>
              </Select>
            </CFormGroup>

            <CFormGroup :description="visibilityDescription">
              <template #label>
                <span class="flex items-center gap-1">
                  {{ $t('block.automation.buttonVisibility.label') }}
                  <i
                    class="pi pi-exclamation-triangle text-orange-500 text-xs"
                    v-tooltip="$t('block.automation.buttonVisibility.tooltip')"
                  />
                </span>
              </template>
              <CInputExpression
                :model-value="item.visibility?.expression || ''"
                dialect="expr"
                :scope="exprScope"
                :min-lines="1"
                :placeholder="$t('block.automation.buttonVisibility.placeholder')"
                @update:model-value="updateVisibility(index, $event)"
              />
            </CFormGroup>
          </div>
        </template>
      </CFormList>
    </CFormGroup>

    <Divider />

    <CFormGroup
      :label="
        $t('block.automation.availableScriptsAndWorkflow', { count: availableTriggers.length })
      "
    >
      <div v-if="loadingTriggers" class="flex justify-center p-3">
        <ProgressSpinner style="width: 24px; height: 24px" />
      </div>

      <div v-else class="flex flex-col gap-1 mt-1">
        <InputText
          v-model="searchQuery"
          :placeholder="$t('block.automation.searchPlaceholder')"
          class="w-full"
        />

        <Tabs
          v-model:value="activeTab"
          class="border border-surface rounded-border overflow-hidden"
        >
          <TabList>
            <Tab value="taqs">{{ $t('block.automation.tabs.taqs') }}</Tab>
            <Tab value="workflows">{{ $t('block.automation.tabs.workflows') }}</Tab>
            <Tab value="scripts">{{ $t('block.automation.tabs.scripts') }}</Tab>
          </TabList>
          <TabPanels>
            <TabPanel value="taqs">
              <CFormItemList
                :items="filteredTaqs"
                :empty-message="$t('block.automation.noScripts')"
                hide-remove
                @select="addTriggerButton"
              >
                <template #default="{ item }">
                  <span class="font-medium text-sm">{{ item.label }}</span>
                  <p v-if="item.description" class="text-sm text-muted-color mt-1 mb-0">
                    {{ item.description }}
                  </p>
                </template>
              </CFormItemList>
            </TabPanel>

            <TabPanel value="workflows">
              <CFormItemList
                :items="filteredWorkflows"
                :empty-message="$t('block.automation.noScripts')"
                hide-remove
                @select="addTriggerButton"
              >
                <template #default="{ item }">
                  <span class="font-medium text-sm">{{ item.label }}</span>
                  <p v-if="item.description" class="text-sm text-muted-color mt-1 mb-0">
                    {{ item.description }}
                  </p>
                </template>
              </CFormItemList>
            </TabPanel>

            <TabPanel value="scripts">
              <CFormItemList
                :items="filteredScripts"
                :empty-message="$t('block.automation.noScripts')"
                hide-remove
                @select="addScriptButton"
              >
                <template #default="{ item }">
                  <div :class="['flex flex-col', { 'opacity-50': !item.applies }]">
                    <div class="flex items-center gap-2 flex-wrap">
                      <span class="font-medium text-sm">{{ item.label }}</span>
                      <Tag
                        v-if="item.resourceType"
                        :value="item.resourceType"
                        severity="info"
                        class="text-xs"
                      />
                      <Tag
                        v-for="(chip, i) in item.constraintChips"
                        :key="i"
                        :value="chip"
                        severity="secondary"
                        class="text-xs"
                      />
                    </div>
                    <p v-if="item.description" class="text-sm text-muted-color mt-1 mb-0">
                      {{ item.description }}
                    </p>
                    <p v-if="!item.applies" class="text-sm text-muted-color italic mt-1 mb-0">
                      {{ $t('block.automation.scriptNeverAppliesHere') }}
                    </p>
                  </div>
                </template>
              </CFormItemList>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </div>
    </CFormGroup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { constraintChips } from '@planetcrust/human-vue'
import { pageFitsResourceType, triggerCanApply } from '@/sections/compose/lib/script-events'

const { t } = useI18n()

const props = defineProps({
  buttons: { type: Array, default: () => [] },
  // ScopeEntry[] for the ${} variables a label may reference.
  scope: {
    type: Array,
    default: () => [],
  },
  // ScopeEntry[] for the server-evaluated visibility condition.
  exprScope: {
    type: Array,
    default: () => [],
  },
  // Whether the block renders beside a record the condition may read.
  isRecordPage: { type: Boolean, default: false },
  // The page the block sits on and the resources it can hand a script. Left out,
  // no script is held back: a caller that names no page states no context.
  page: { type: Object, default: null },
  namespace: { type: Object, default: null },
  module: { type: Object, default: null },
  // Whether the block hands a script the page's record; a record list never does.
  canSupplyRecord: { type: Boolean, default: true },
})

const emit = defineEmits(['update:buttons'])

const $AutomationAPI = inject('$AutomationAPI', null)
const $ComposeAPI = inject('$ComposeAPI', null)

// Button labels are interpolated at render time (AutomationBlock.buttonLabel).
const labelInputs = ref([])

const activeTab = ref('taqs')
const selectedIndex = ref(-1)
const searchQuery = ref('')
const loadingTriggers = ref(false)
const triggerButtons = ref([])
const workflowData = ref([])
const taqData = ref([])

const variantKeys = ['primary', 'secondary', 'success', 'danger', 'warning', 'info']

const variantLabel = key => t(`block.automation.variants.${key}`)

const variantOptions = computed(() =>
  variantKeys.map(value => ({ value, label: variantLabel(value) })),
)

const variantSeverityMap = {
  primary: undefined,
  secondary: 'secondary',
  success: 'success',
  danger: 'danger',
  warning: 'warn',
  info: 'info',
}

const mapVariantSeverity = key => variantSeverityMap[key]

const visibilityDescription = computed(() =>
  props.isRecordPage
    ? t('block.general.visibility.condition.description.record-page', [
        'record.values.fieldName',
        'user.(userID/email...)',
        'screen.(width/height)',
        'isView/isCreate/isEdit',
        'record.values.status == "draft"',
        'user.userID == record.ownedBy',
      ])
    : t('block.general.visibility.condition.description.non-record-page', [
        'user.(userID/email...)',
        'screen.(width/height)',
        'user.email == "test@mail.com"',
        'screen.width < 1024',
      ]),
)

const scriptTypeOf = b => {
  if (b.automationID) return 'taq'
  if (b.script) return 'script'
  return 'workflow'
}

const normalizedButtons = computed(() =>
  (props.buttons || []).map(b => ({
    ...b,
    scriptType: b.scriptType || scriptTypeOf(b),
  })),
)

// CFormList does in-place splice for remove and reassignment for reorder.
// Get returns the prop array directly so in-place mutations reach the parent;
// set emits an update for full reassignments (reorder).
const buttonsModel = computed({
  get: () => props.buttons || [],
  set: next => emit('update:buttons', next),
})

const triggerKey = t => {
  if (t.automationID) return `taq-${t.automationID}-${t.triggerHandle}`
  if (t.script) return `script-${t.script}`
  return `${t.workflowID}-${t.stepID}`
}

const availableTriggers = computed(() => {
  const existingKeys = normalizedButtons.value.map(triggerKey)
  return triggerButtons.value.filter(t => !existingKeys.includes(triggerKey(t)))
})

const filteredWorkflows = computed(() => {
  const triggers = availableTriggers.value.filter(t => !t.isTAQ && !t.script)
  if (!searchQuery.value) return triggers
  const q = searchQuery.value.toLowerCase()
  return triggers.filter(t => `${t.label} ${t.description || ''}`.toLowerCase().includes(q))
})

// Scripts this block may offer, each told how it stands against the page: one
// whose resource the page cannot hand it is left out, one whose constraints the
// page contradicts is shown and refused.
const filteredScripts = computed(() => {
  const triggers = availableTriggers.value
    .filter(t => t.script)
    .filter(t => !props.page || pageFitsResourceType(props.page, t.resourceType))
    .filter(t => props.canSupplyRecord || t.resourceType !== 'compose:record')
    .map(t => ({
      ...t,
      constraintChips: constraintChips(t.constraints),
      applies: triggerCanApply(
        { constraints: t.constraints },
        { namespace: props.namespace, module: props.module },
      ),
    }))

  if (!searchQuery.value) return triggers
  const q = searchQuery.value.toLowerCase()
  return triggers.filter(t => `${t.label} ${t.description || ''}`.toLowerCase().includes(q))
})

const filteredTaqs = computed(() => {
  const triggers = availableTriggers.value.filter(t => t.isTAQ)
  if (!searchQuery.value) return triggers
  const q = searchQuery.value.toLowerCase()
  return triggers.filter(t => `${t.label} ${t.description || ''}`.toLowerCase().includes(q))
})

function emitButtons(next) {
  emit('update:buttons', next)
}

function updateField(index, key, value) {
  const next = [...normalizedButtons.value]
  next[index] = { ...next[index], [key]: value }
  emitButtons(next)
}

// An empty condition drops the key, so a button never carries a blank rule.
function updateVisibility(index, expression) {
  const next = [...normalizedButtons.value]
  const { visibility: _, ...rest } = next[index]
  next[index] = expression ? { ...rest, visibility: { expression } } : rest
  emitButtons(next)
}

function addPlaceholder() {
  const next = [
    ...normalizedButtons.value,
    {
      label: t('block.automation.dummyButtonLabel'),
      variant: 'primary',
      resourceType: 'compose',
      scriptType: 'workflow',
    },
  ]
  emitButtons(next)
  selectedIndex.value = next.length - 1
}

function addTriggerButton(trigger) {
  const newButton = {
    label: trigger.label,
    variant: 'primary',
    resourceType: trigger.resourceType || 'compose',
  }

  if (trigger.isTAQ) {
    newButton.automationID = trigger.automationID
    newButton.triggerHandle = trigger.triggerHandle
    newButton.scriptType = 'taq'
  } else if (trigger.isScript) {
    newButton.script = trigger.script
    newButton.scriptType = 'script'
  } else {
    newButton.workflowID = trigger.workflowID
    newButton.stepID = trigger.stepID
    newButton.scriptType = 'workflow'
  }

  emitButtons([...normalizedButtons.value, newButton])
}

// A script the page contradicts is listed so it can be recognised, not picked.
function addScriptButton(trigger) {
  if (trigger.applies === false) return
  addTriggerButton(trigger)
}

function selectButton(index) {
  selectedIndex.value = selectedIndex.value === index ? -1 : index
}

// Apps whose manual Corredor scripts a compose page button can trigger. A
// trigger naming no app at all is offered everywhere.
const scriptApps = ['compose', 'unify']

async function fetchScriptTriggers() {
  if (!$ComposeAPI) return []

  try {
    const { set = [] } = await $ComposeAPI.automationList({
      eventTypes: ['onManual'],
      excludeInvalid: true,
    })

    return set.flatMap(script =>
      (script.triggers || [])
        .filter(trigger => {
          const app = (trigger.uiProps || []).find(p => p.name === 'app')?.value
          return !app || scriptApps.includes(app)
        })
        .map(trigger => ({
          script: script.name,
          label: script.label || script.name,
          resourceType: (trigger.resourceTypes || [])[0],
          description: script.description,
          isScript: true,
          constraints: trigger.constraints || [],
        })),
    )
  } catch (e) {
    console.error('Failed to fetch automation scripts:', e)
    return []
  }
}

async function fetchTriggers() {
  if (!$AutomationAPI) return

  loadingTriggers.value = true

  try {
    const [workflowsResp, taqsResp] = await Promise.all([
      $AutomationAPI.triggerList({ eventType: 'onManual' }),
      $AutomationAPI.ngAutomationListCancellable({ limit: 100, sort: 'name ASC' }),
    ])

    const { set: triggers = [] } = workflowsResp
    const taqResult = await taqsResp.response()
    const taqsRaw = Array.isArray(taqResult) ? taqResult : taqResult.set || []
    taqData.value = taqsRaw

    const triggerData = triggers.map(({ triggerID, workflowID, resourceType, stepID }) => ({
      triggerID,
      workflowID,
      resourceType,
      stepID,
    }))

    const workflowIDs = [...new Set(triggers.map(t => t.workflowID))]

    if (workflowIDs.length) {
      const { set: wfSet = [] } = await $AutomationAPI.workflowList({ workflowID: workflowIDs })
      workflowData.value = wfSet

      triggerButtons.value = triggerData
        .map(trigger => {
          const wf = wfSet.find(w => w.workflowID === trigger.workflowID)
          if (!wf) return null

          let label = wf.meta?.name || wf.handle || wf.workflowID
          const step = (wf.steps || []).find(s => s.stepID === trigger.stepID)
          if (step?.meta?.label) label = `${label} (${step.meta.label})`

          return {
            label,
            workflowID: trigger.workflowID,
            stepID: trigger.stepID,
            resourceType: trigger.resourceType,
            description: wf.meta?.description,
          }
        })
        .filter(Boolean)
    }

    const taqButtons = taqData.value.flatMap(taq => {
      const triggers = taq.triggers || []
      const manualTriggers = triggers.filter(t => t.eventType === 'onManual')

      if (manualTriggers.length === 0) return []

      return manualTriggers.map(t => {
        let label = taq.meta?.short || taq.handle || taq.automationID
        if (t.meta?.short) label = `${label} (${t.meta.short})`

        return {
          label,
          automationID: taq.automationID,
          triggerHandle: t.handle,
          resourceType: t.resourceType || 'compose',
          description: taq.meta?.description,
          isTAQ: true,
        }
      })
    })

    const scriptButtons = await fetchScriptTriggers()

    triggerButtons.value = [...(triggerButtons.value || []), ...taqButtons, ...scriptButtons]
  } catch (e) {
    console.error('Failed to fetch triggers:', e)
  } finally {
    loadingTriggers.value = false
  }
}

onMounted(() => fetchTriggers())
</script>
