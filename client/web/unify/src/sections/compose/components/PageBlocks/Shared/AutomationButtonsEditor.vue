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
          </TabPanels>
        </Tabs>
      </div>
    </CFormGroup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  buttons: { type: Array, default: () => [] },
  // ScopeEntry[] for the ${} variables a label may reference.
  scope: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:buttons'])

const $AutomationAPI = inject('$AutomationAPI', null)

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

const normalizedButtons = computed(() =>
  (props.buttons || []).map(b => ({
    ...b,
    scriptType: b.scriptType || (b.automationID ? 'taq' : 'workflow'),
  })),
)

// CFormList does in-place splice for remove and reassignment for reorder.
// Get returns the prop array directly so in-place mutations reach the parent;
// set emits an update for full reassignments (reorder).
const buttonsModel = computed({
  get: () => props.buttons || [],
  set: next => emit('update:buttons', next),
})

const availableTriggers = computed(() => {
  const existingKeys = normalizedButtons.value.map(b => {
    if (b.automationID) return `taq-${b.automationID}-${b.triggerHandle}`
    return b.workflowID ? `${b.workflowID}-${b.stepID}` : b.script
  })
  return triggerButtons.value.filter(t => {
    const key = t.isTAQ ? `taq-${t.automationID}-${t.triggerHandle}` : `${t.workflowID}-${t.stepID}`
    return !existingKeys.includes(key)
  })
})

const filteredWorkflows = computed(() => {
  const triggers = availableTriggers.value.filter(t => !t.isTAQ && !t.script)
  if (!searchQuery.value) return triggers
  const q = searchQuery.value.toLowerCase()
  return triggers.filter(t => `${t.label} ${t.description || ''}`.toLowerCase().includes(q))
})

const filteredScripts = computed(() => {
  const triggers = availableTriggers.value.filter(t => t.script)
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
  } else {
    newButton.workflowID = trigger.workflowID
    newButton.stepID = trigger.stepID
    newButton.scriptType = 'workflow'
  }

  emitButtons([...normalizedButtons.value, newButton])
}

function selectButton(index) {
  selectedIndex.value = selectedIndex.value === index ? -1 : index
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

    triggerButtons.value = [...(triggerButtons.value || []), ...taqButtons]
  } catch (e) {
    console.error('Failed to fetch triggers:', e)
  } finally {
    loadingTriggers.value = false
  }
}

onMounted(() => fetchTriggers())
</script>
