<template>
  <div class="flex flex-col gap-3">
    <!-- Configured buttons -->
    <div class="flex flex-col gap-1">
      <div class="flex items-center justify-between">
        <label class="text-primary font-medium text-sm">{{ $t('block.automation.configuredButtons') }}</label>
        <Button
          :label="$t('block.automation.addPlaceholderLabel')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          @click="addPlaceholder"
        />
      </div>

      <div v-if="!buttons.length" class="text-muted-color text-sm italic p-2">
        {{ $t('block.automation.noScripts') }}
      </div>

      <div v-else class="flex flex-col gap-2">
        <div
          v-for="(btn, i) in buttons"
          :key="i"
          class="border border-surface rounded-border p-3"
          :class="{ 'border-primary': selectedIndex === i }"
        >
          <div class="flex items-center gap-2">
            <div class="flex flex-col gap-1">
              <Button
                icon="pi pi-chevron-up"
                text rounded size="small" severity="secondary"
                :disabled="i === 0"
                @click="moveButton(i, -1)"
              />
              <Button
                icon="pi pi-chevron-down"
                text rounded size="small" severity="secondary"
                :disabled="i === buttons.length - 1"
                @click="moveButton(i, 1)"
              />
            </div>
            <Tag :severity="mapSeverity(btn.variant)" :value="btn.variant || 'primary'" />
            <span class="flex-1 font-medium text-sm truncate">{{ btn.label || '-' }}</span>
            <Button
              icon="pi pi-pencil"
              text rounded size="small"
              @click="selectButton(i)"
            />
            <Button
              icon="pi pi-trash"
              text rounded size="small" severity="danger"
              @click="removeButton(i)"
            />
          </div>

          <!-- Inline editor when selected -->
          <div v-if="selectedIndex === i" class="mt-3 flex flex-col gap-2 border-t border-surface pt-3">
            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.automation.buttonLabel') }}</label>
              <InputText v-model="btn.label" class="w-full" />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.automation.buttonVariant') }}</label>
              <Select
                v-model="btn.variant"
                :options="variantOptions"
                class="w-full"
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.automation.buttonWorkflow') }}</label>
              <Select
                v-model="btn.workflowID"
                :options="workflows"
                option-label="label"
                option-value="workflowID"
                :placeholder="$t('block.automation.pickWorkflow')"
                class="w-full"
                filter
                show-clear
              />
            </div>

            <div v-if="btn.workflowID" class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.automation.buttonStep') }}</label>
              <Select
                v-model="btn.stepID"
                :options="getSteps(btn.workflowID)"
                option-label="label"
                option-value="stepID"
                :placeholder="$t('block.automation.pickStep')"
                class="w-full"
                filter
                show-clear
              />
            </div>

            <div class="flex flex-col gap-1">
              <label class="text-sm text-muted-color">{{ $t('block.automation.buttonResourceType') }}</label>
              <Select
                v-model="btn.resourceType"
                :options="resourceTypeOptions"
                option-label="label"
                option-value="value"
                class="w-full"
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
        {{ $t('block.automation.availableScriptsAndWorkflow', { count: availableTriggers.length }) }}
      </label>

      <div v-if="loadingTriggers" class="flex justify-center p-3">
        <ProgressSpinner style="width: 24px; height: 24px" />
      </div>

      <div v-else-if="availableTriggers.length" class="flex flex-col gap-1">
        <InputText
          v-model="searchQuery"
          :placeholder="$t('block.automation.searchPlaceholder')"
          class="w-full"
        />

        <div
          v-for="(trigger, i) in filteredTriggers"
          :key="i"
          class="p-3 border border-surface rounded-border cursor-pointer hover:bg-highlight transition-colors"
          @click="addTriggerButton(trigger)"
        >
          <div class="flex items-center gap-2">
            <span class="font-medium text-sm">{{ trigger.label }}</span>
            <Tag severity="info" value="workflow" class="text-xs" />
          </div>
          <p v-if="trigger.description" class="text-sm text-muted-color mt-1 mb-0">{{ trigger.description }}</p>
        </div>
      </div>

      <div v-else class="text-muted-color text-sm italic p-2">
        {{ $t('block.automation.noScripts') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const $AutomationAPI = inject('$AutomationAPI', null)

const selectedIndex = ref(-1)
const searchQuery = ref('')
const loadingTriggers = ref(false)
const triggerButtons = ref([])
const workflowData = ref([])

const variantOptions = ['primary', 'secondary', 'success', 'danger', 'warning', 'info']

const resourceTypeOptions = [
  { value: 'compose', label: 'Compose' },
  { value: 'compose:namespace', label: 'Namespace' },
  { value: 'compose:page', label: 'Page' },
  { value: 'compose:module', label: 'Module' },
  { value: 'compose:record', label: 'Record' },
]

const buttons = computed(() => props.block.options?.buttons || [])

const workflows = computed(() => {
  return workflowData.value.map(wf => ({
    workflowID: wf.workflowID,
    label: wf.meta?.name || wf.handle || wf.workflowID,
    steps: wf.steps || [],
  }))
})

function getSteps(workflowID) {
  const wf = workflowData.value.find(w => w.workflowID === workflowID)
  if (!wf) return []
  return (wf.steps || []).map(s => ({
    stepID: s.stepID,
    label: s.meta?.label || s.stepID,
  }))
}

const availableTriggers = computed(() => {
  const existingKeys = buttons.value.map(b => b.workflowID ? `${b.workflowID}-${b.stepID}` : b.script)
  return triggerButtons.value.filter(t => !existingKeys.includes(`${t.workflowID}-${t.stepID}`))
})

const filteredTriggers = computed(() => {
  if (!searchQuery.value) return availableTriggers.value
  const q = searchQuery.value.toLowerCase()
  return availableTriggers.value.filter(t =>
    `${t.label} ${t.description || ''}`.toLowerCase().includes(q),
  )
})

function mapSeverity(variant) {
  const map = { primary: undefined, secondary: 'secondary', success: 'success', danger: 'danger', warning: 'warn', info: 'info' }
  return map[variant] || undefined
}

function updateButtons(newButtons) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, buttons: newButtons },
  })
}

function addPlaceholder() {
  const newButtons = [...buttons.value, {
    label: t('block.automation.dummyButtonLabel'),
    variant: 'primary',
    resourceType: 'compose',
  }]
  updateButtons(newButtons)
  selectedIndex.value = newButtons.length - 1
}

function addTriggerButton(trigger) {
  const newButtons = [...buttons.value, {
    label: trigger.label,
    variant: 'primary',
    workflowID: trigger.workflowID,
    stepID: trigger.stepID,
    resourceType: trigger.resourceType || 'compose',
  }]
  updateButtons(newButtons)
}

function removeButton(index) {
  const newButtons = [...buttons.value]
  newButtons.splice(index, 1)
  updateButtons(newButtons)
  if (selectedIndex.value === index) selectedIndex.value = -1
}

function selectButton(index) {
  selectedIndex.value = selectedIndex.value === index ? -1 : index
}

function moveButton(index, direction) {
  const newButtons = [...buttons.value]
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= newButtons.length) return
  const temp = newButtons[index]
  newButtons[index] = newButtons[newIndex]
  newButtons[newIndex] = temp
  updateButtons(newButtons)
}

async function fetchTriggers() {
  if (!$AutomationAPI) return

  loadingTriggers.value = true

  try {
    const { set: triggers = [] } = await $AutomationAPI.triggerList({ eventType: 'onManual' })

    const triggerData = triggers.map(({ triggerID, workflowID, resourceType, stepID }) => ({
      triggerID, workflowID, resourceType, stepID,
    }))

    const workflowIDs = [...new Set(triggers.map(t => t.workflowID))]

    if (workflowIDs.length) {
      const { set: wfSet = [] } = await $AutomationAPI.workflowList({ workflowID: workflowIDs })
      workflowData.value = wfSet

      triggerButtons.value = triggerData.map(trigger => {
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
      }).filter(Boolean)
    }
  } catch (e) {
    console.error('Failed to fetch triggers:', e)
  } finally {
    loadingTriggers.value = false
  }
}

onMounted(() => fetchTriggers())
</script>
