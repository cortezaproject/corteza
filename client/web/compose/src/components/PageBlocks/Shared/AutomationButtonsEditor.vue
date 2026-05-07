<template>
  <div class="flex flex-col gap-3">
    <!-- Configured buttons -->
    <div class="flex flex-col gap-1">
      <div class="flex items-center justify-between">
        <label class="text-primary font-medium text-sm">
          {{ $t('block.automation.configuredButtons') }}
        </label>
        <Button
          :label="$t('block.automation.addPlaceholderLabel')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          @click="addPlaceholder"
        />
      </div>

      <div v-if="!normalizedButtons.length" class="text-muted-color text-sm italic p-2">
        {{ $t('block.automation.noScripts') }}
      </div>

      <div v-else class="flex flex-col gap-2">
        <div
          v-for="(btn, i) in normalizedButtons"
          :key="i"
          class="border border-surface rounded-border p-3"
          :class="{ 'border-primary': selectedIndex === i }"
        >
          <div class="flex items-center gap-2">
            <div class="flex gap-1">
              <Button
                icon="pi pi-chevron-up"
                text
                rounded
                size="small"
                severity="secondary"
                :disabled="i === 0"
                @click="moveButton(i, -1)"
              />
              <Button
                icon="pi pi-chevron-down"
                text
                rounded
                size="small"
                severity="secondary"
                :disabled="i === normalizedButtons.length - 1"
                @click="moveButton(i, 1)"
              />
            </div>
            <Tag
              v-if="(btn.variant || 'primary') !== 'primary'"
              :severity="mapSeverity(btn.variant)"
              :value="btn.variant"
            />
            <span class="flex-1 font-medium text-sm truncate">{{ btn.label || '-' }}</span>
            <Button icon="pi pi-pencil" text rounded size="small" @click="selectButton(i)" />
            <Button
              icon="pi pi-trash"
              text
              rounded
              size="small"
              severity="danger"
              @click="removeButton(i)"
            />
          </div>

          <!-- Inline editor when selected -->
          <div
            v-if="selectedIndex === i"
            class="mt-3 flex flex-col gap-2 border-t border-surface pt-3"
          >
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">
                {{ $t('block.automation.buttonLabel') }}
              </label>
              <InputText
                :model-value="btn.label"
                class="w-full"
                @update:model-value="updateField(i, 'label', $event)"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">
                {{ $t('block.automation.buttonVariant') }}
              </label>
              <Select
                :model-value="btn.variant"
                :options="variantOptions"
                class="w-full"
                @update:model-value="updateField(i, 'variant', $event)"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">
                {{ $t('block.automation.buttonResourceType') }}
              </label>
              <Select
                :model-value="btn.resourceType"
                :options="resourceTypeOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                @update:model-value="updateField(i, 'resourceType', $event)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <Divider />

    <!-- Available workflows/triggers -->
    <div class="flex flex-col gap-1">
      <label class="text-primary font-medium text-sm">
        {{
          $t('block.automation.availableScriptsAndWorkflow', { count: availableTriggers.length })
        }}
      </label>

      <div v-if="loadingTriggers" class="flex justify-center p-3">
        <ProgressSpinner style="width: 24px; height: 24px" />
      </div>

      <div v-else class="flex flex-col gap-1 mt-1">
        <InputText
          v-model="searchQuery"
          :placeholder="$t('block.automation.searchPlaceholder')"
          class="w-full"
        />

        <Tabs v-model:value="activeTab">
          <TabList>
            <Tab value="taqs">{{ $t('block.automation.tabs.taqs') }}</Tab>
            <Tab value="workflows">{{ $t('block.automation.tabs.workflows') }}</Tab>
            <Tab value="scripts">{{ $t('block.automation.tabs.scripts') }}</Tab>
          </TabList>
          <TabPanels class="px-0 pb-0">
            <TabPanel value="taqs">
              <div v-if="filteredTaqs.length" class="flex flex-col gap-1">
                <div
                  v-for="(trigger, i) in filteredTaqs"
                  :key="i"
                  class="p-3 border border-surface rounded-border cursor-pointer hover:bg-highlight transition-colors"
                  @click="addTriggerButton(trigger)"
                >
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-sm">{{ trigger.label }}</span>
                  </div>
                  <p v-if="trigger.description" class="text-sm text-muted-color mt-1 mb-0">
                    {{ trigger.description }}
                  </p>
                </div>
              </div>
              <div v-else class="text-muted-color text-sm italic p-2">
                {{ $t('block.automation.noScripts') }}
              </div>
            </TabPanel>

            <TabPanel value="workflows">
              <div v-if="filteredWorkflows.length" class="flex flex-col gap-1">
                <div
                  v-for="(trigger, i) in filteredWorkflows"
                  :key="i"
                  class="p-3 border border-surface rounded-border cursor-pointer hover:bg-highlight transition-colors"
                  @click="addTriggerButton(trigger)"
                >
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-sm">{{ trigger.label }}</span>
                  </div>
                  <p v-if="trigger.description" class="text-sm text-muted-color mt-1 mb-0">
                    {{ trigger.description }}
                  </p>
                </div>
              </div>
              <div v-else class="text-muted-color text-sm italic p-2">
                {{ $t('block.automation.noScripts') }}
              </div>
            </TabPanel>

            <TabPanel value="scripts">
              <div v-if="filteredScripts.length" class="flex flex-col gap-1">
                <div
                  v-for="(trigger, i) in filteredScripts"
                  :key="i"
                  class="p-3 border border-surface rounded-border cursor-pointer hover:bg-highlight transition-colors"
                  @click="addTriggerButton(trigger)"
                >
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-sm">{{ trigger.label }}</span>
                  </div>
                  <p v-if="trigger.description" class="text-sm text-muted-color mt-1 mb-0">
                    {{ trigger.description }}
                  </p>
                </div>
              </div>
              <div v-else class="text-muted-color text-sm italic p-2">
                {{ $t('block.automation.noScripts') }}
              </div>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  buttons: { type: Array, default: () => [] },
})

const emit = defineEmits(['update:buttons'])

const $AutomationAPI = inject('$AutomationAPI', null)

const activeTab = ref('taqs')
const selectedIndex = ref(-1)
const searchQuery = ref('')
const loadingTriggers = ref(false)
const triggerButtons = ref([])
const workflowData = ref([])
const taqData = ref([])

const variantOptions = ['primary', 'secondary', 'success', 'danger', 'warning', 'info']

const resourceTypeOptions = [
  { value: 'compose', label: 'Compose' },
  { value: 'compose:namespace', label: 'Namespace' },
  { value: 'compose:page', label: 'Page' },
  { value: 'compose:module', label: 'Module' },
  { value: 'compose:record', label: 'Record' },
]

const normalizedButtons = computed(() =>
  (props.buttons || []).map(b => ({
    ...b,
    scriptType: b.scriptType || (b.automationID ? 'taq' : 'workflow'),
  })),
)

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

function mapSeverity(variant) {
  const map = {
    primary: undefined,
    secondary: 'secondary',
    success: 'success',
    danger: 'danger',
    warning: 'warn',
    info: 'info',
  }
  return map[variant] || undefined
}

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

function removeButton(index) {
  const next = [...normalizedButtons.value]
  next.splice(index, 1)
  emitButtons(next)
  if (selectedIndex.value === index) selectedIndex.value = -1
}

function selectButton(index) {
  selectedIndex.value = selectedIndex.value === index ? -1 : index
}

function moveButton(index, direction) {
  const next = [...normalizedButtons.value]
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= next.length) return
  const temp = next[index]
  next[index] = next[newIndex]
  next[newIndex] = temp
  emitButtons(next)
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
