<template>
  <Panel :header="$t('chatbot.editor.panels.scenarios')" toggleable>
    <p class="text-xs text-muted-color mb-3">
      {{ $t('chatbot.editor.scenarios.description') }}
    </p>
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-4 min-h-[320px]">
      <div class="lg:col-span-1 flex flex-col gap-2 border border-surface rounded-lg p-2">
        <div class="flex items-center justify-between">
          <span class="font-medium text-sm">
            {{ $t('chatbot.editor.scenarios.list') }}
          </span>
          <Button
            :label="$t('chatbot.editor.scenarios.add')"
            icon="pi pi-plus"
            severity="secondary"
            size="small"
            @click="addScenario"
          />
        </div>
        <div v-if="!scenarios.length" class="text-xs text-muted-color p-2">
          {{ $t('chatbot.editor.scenarios.empty') }}
        </div>
        <TransitionGroup
          tag="ul"
          class="flex flex-col gap-1"
          move-class="transition-transform duration-300 ease-in-out"
        >
          <li
            v-for="(s, idx) in scenarios"
            :key="s.id || idx"
            class="group flex items-center gap-2 p-2 rounded cursor-grab"
            :class="[
              selectedIdx === idx ? 'bg-primary/10 border border-primary/40' : 'hover:bg-emphasis',
              dropTargetIdx === idx && draggedIdx !== idx ? 'border-t-2 !border-t-primary' : '',
            ]"
            draggable="true"
            @click="selectedIdx = idx"
            @dragstart="onDragStart(idx)"
            @dragover.prevent="onDragOver(idx)"
            @dragleave="onDragLeave"
            @drop.prevent="onDrop(idx)"
          >
            <span class="flex-1 text-sm truncate">
              {{ s.name || s.id || $t('chatbot.editor.scenarios.unnamed') }}
            </span>
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              :class="[
                'transition-opacity',
                selectedIdx === idx ? 'opacity-100' : 'opacity-0 group-hover:opacity-100',
              ]"
              @click.stop="removeScenario(idx)"
            />
          </li>
        </TransitionGroup>
      </div>

      <div class="xl:col-span-2 flex flex-col gap-3 border border-surface rounded-lg p-3">
        <template v-if="current">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
            <CFormGroup :label="$t('chatbot.editor.scenarios.name')">
              <InputText v-model="current.name" />
            </CFormGroup>
            <CFormGroup :label="$t('chatbot.editor.scenarios.handle')">
              <InputText v-model="current.id" :invalid="handleIsDuplicate" />
              <small v-if="handleIsDuplicate" class="text-red-500">
                {{ $t('chatbot.editor.scenarios.handleDuplicate') }}
              </small>
            </CFormGroup>
          </div>

          <div class="grid grid-cols-1 xl:grid-cols-2 gap-3">
            <CFormGroup :label="$t('chatbot.editor.scenarios.type')">
              <Select
                v-model="current.type"
                :options="scenarioTypes"
                optionLabel="label"
                optionValue="value"
                @change="ensureScenarioConfig"
              />
            </CFormGroup>
          </div>

          <Fieldset :legend="$t(`chatbot.editor.scenarios.types.${current.type}`)">
            <div class="grid grid-cols-1 xl:grid-cols-2 gap-3">
              <template v-if="current.type === 'static_message'">
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.static.message')"
                  class="xl:col-span-2"
                >
                  <Textarea v-model="current.config.message" :rows="3" />
                </CFormGroup>
                <CInputToggleCard
                  v-model="current.config.isMarkdown"
                  :label="$t('chatbot.editor.scenarios.static.isMarkdown')"
                  :description="$t('chatbot.editor.scenarios.static.isMarkdownDescription')"
                />
                <CFormGroup :label="$t('chatbot.editor.scenarios.static.autoAdvance')">
                  <InputNumber v-model="current.config.autoAdvanceMs" :min="0" />
                </CFormGroup>
              </template>

              <template v-else-if="current.type === 'conversation'">
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.agent.label')"
                  :description="$t('chatbot.editor.scenarios.agent.help')"
                  required
                >
                  <div class="flex gap-2 items-start">
                    <Select
                      v-model="current.agentID"
                      :options="agentOptions"
                      optionLabel="label"
                      optionValue="value"
                      :placeholder="$t('chatbot.editor.scenarios.agent.placeholder')"
                      filter
                      class="flex-1"
                    />
                    <Button
                      v-if="current.agentID"
                      v-tooltip.bottom="$t('chatbot.editor.scenarios.agent.open')"
                      icon="pi pi-external-link"
                      severity="secondary"
                      @click="openAgent(current.agentID)"
                    />
                    <Button
                      v-else
                      v-tooltip.bottom="$t('chatbot.editor.scenarios.agent.create')"
                      icon="pi pi-plus"
                      severity="secondary"
                      @click="createAgent"
                    />
                  </div>
                </CFormGroup>
                <CFormGroup :label="$t('chatbot.editor.scenarios.conversation.placeholder')">
                  <InputText v-model="current.config.placeholder" />
                </CFormGroup>
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.conversation.initialPrompt')"
                  class="xl:col-span-2"
                >
                  <Textarea v-model="current.config.initialPrompt" :rows="2" />
                </CFormGroup>
                <CInputToggleCard
                  v-model="current.config.typingIndicator"
                  :label="$t('chatbot.editor.scenarios.conversation.typingIndicator')"
                  :description="
                    $t('chatbot.editor.scenarios.conversation.typingIndicatorDescription')
                  "
                />
              </template>

              <template v-else-if="current.type === 'form'">
                <CFormGroup :label="$t('chatbot.editor.scenarios.form.fields')" class="xl:col-span-2">
                  <template #actions>
                    <Button
                      :label="$t('general.label.add')"
                      icon="pi pi-plus"
                      severity="secondary"
                      size="small"
                      @click="
                        current.config.fields.push({
                          name: '',
                          label: '',
                          type: 'text',
                          required: false,
                        })
                      "
                    />
                  </template>
                  <CFormList
                    v-model="current.config.fields"
                    :columns="[
                      { label: $t('chatbot.editor.scenarios.form.fieldName'), width: '3fr' },
                      { label: $t('chatbot.editor.scenarios.form.fieldLabel'), width: '4fr' },
                      { label: $t('chatbot.editor.scenarios.form.fieldType'), width: '3fr' },
                      {
                        label: $t('chatbot.editor.scenarios.form.fieldRequired'),
                        width: '6rem',
                        headerClass: 'text-center',
                      },
                    ]"
                    :empty-message="$t('chatbot.editor.scenarios.form.empty')"
                  >
                    <template #row="{ item }">
                      <InputText v-model="item.name" size="small" class="w-full" />
                      <InputText v-model="item.label" size="small" class="w-full" />
                      <Select
                        v-model="item.type"
                        :options="fieldTypeOptions"
                        optionLabel="label"
                        optionValue="value"
                        size="small"
                        class="w-full"
                      />
                      <div class="flex items-center justify-center">
                        <ToggleSwitch v-model="item.required" />
                      </div>
                    </template>
                  </CFormList>
                </CFormGroup>
                <CFormGroup :label="$t('chatbot.editor.scenarios.form.submitLabel')">
                  <InputText v-model="current.config.submitLabel" />
                </CFormGroup>
              </template>
            </div>
          </Fieldset>

          <Divider />

          <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
            <CFormGroup
              :label="$t('chatbot.editor.scenarios.automation.before')"
              :description="$t('chatbot.editor.scenarios.automation.beforeHelp')"
            >
              <CInputTAQ
                :model-value="automationID(current.automation?.before?.automation)"
                @update:model-value="setAutomation('before', $event)"
              />
            </CFormGroup>
            <CFormGroup
              :label="$t('chatbot.editor.scenarios.automation.after')"
              :description="$t('chatbot.editor.scenarios.automation.afterHelp')"
            >
              <CInputTAQ
                :model-value="automationID(current.automation?.after?.automation)"
                @update:model-value="setAutomation('after', $event)"
              />
            </CFormGroup>
          </div>
        </template>
        <div v-else class="text-sm text-muted-color p-4 text-center">
          {{ $t('chatbot.editor.scenarios.selectHint') }}
        </div>
      </div>
    </div>

    <Divider class="!my-4" />

    <div class="flex flex-col gap-4">
      <h3 class="text-primary font-medium">
        {{ $t('chatbot.editor.handoff.sectionTitle') }}
      </h3>

      <CInputToggleCard
        v-model="handoff.enabled"
        :label="$t('chatbot.editor.handoff.enabled.label')"
        :description="$t('chatbot.editor.handoff.enabled.help')"
        class="self-start"
      />

      <template v-if="handoff.enabled">
        <CFormGroup
          :label="$t('chatbot.editor.handoff.targetRoles.label')"
          :description="$t('chatbot.editor.handoff.targetRoles.help')"
        >
          <CInputRole
            v-model="handoff.targetRoles"
            multiple
            :placeholder="$t('chatbot.editor.handoff.targetRoles.placeholder')"
          />
        </CFormGroup>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
          <CFormGroup
            :label="$t('chatbot.editor.handoff.automation.onRequested.label')"
            :description="$t('chatbot.editor.handoff.automation.onRequested.help')"
          >
            <CInputTAQ
              :model-value="automationID(handoff.automation?.onRequested?.automation)"
              @update:model-value="setHandoffAutomation('onRequested', $event)"
            />
          </CFormGroup>
          <CFormGroup
            :label="$t('chatbot.editor.handoff.automation.onAccepted.label')"
            :description="$t('chatbot.editor.handoff.automation.onAccepted.help')"
          >
            <CInputTAQ
              :model-value="automationID(handoff.automation?.onAccepted?.automation)"
              @update:model-value="setHandoffAutomation('onAccepted', $event)"
            />
          </CFormGroup>
        </div>
      </template>
    </div>
  </Panel>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CInputTAQ, CInputRole, CInputToggleCard } = components
const { t } = useI18n()

const AUTOMATION_PREFIX = 'corteza::automation:ng-automation/'

function automationID(resourceStr) {
  if (!resourceStr) return null
  return resourceStr.startsWith(AUTOMATION_PREFIX)
    ? resourceStr.slice(AUTOMATION_PREFIX.length)
    : resourceStr
}

function setAutomation(field, id) {
  if (!current.value) return
  if (!current.value.automation) current.value.automation = emptyAutomation()
  const hook = current.value.automation[field] || { automation: '', async: false }
  current.value.automation[field] = {
    automation: id ? `${AUTOMATION_PREFIX}${id}` : '',
    async: !!hook.async,
  }
}

function emptyAutomation() {
  return {
    before: { automation: '', async: false },
    after: { automation: '', async: false },
  }
}

const props = defineProps({
  scenarios: { type: Array, required: true },
  agents: { type: Array, default: () => [] },
  handoff: { type: Object, required: true },
})

function setHandoffAutomation(phase, id) {
  if (!props.handoff.automation) {
    props.handoff.automation = {
      onRequested: { automation: '', async: false },
      onAccepted: { automation: '', async: false },
    }
  }
  const prev = props.handoff.automation[phase] || { automation: '', async: false }
  props.handoff.automation[phase] = {
    automation: id ? `${AUTOMATION_PREFIX}${id}` : '',
    async: !!prev.async,
  }
}

const agentOptions = computed(() =>
  props.agents.map(a => ({
    value: a.agentID,
    label: a.meta?.short || a.handle || a.agentID,
  })),
)

const selectedIdx = ref(props.scenarios.length ? 0 : -1)

const current = computed(() => (selectedIdx.value >= 0 ? props.scenarios[selectedIdx.value] : null))

const scenarioTypes = computed(() => [
  { label: t('chatbot.editor.scenarios.types.static_message'), value: 'static_message' },
  { label: t('chatbot.editor.scenarios.types.conversation'), value: 'conversation' },
  { label: t('chatbot.editor.scenarios.types.form'), value: 'form' },
])

const fieldTypeOptions = computed(() => [
  { label: t('chatbot.editor.scenarios.form.fieldTypes.text'), value: 'text' },
  { label: t('chatbot.editor.scenarios.form.fieldTypes.email'), value: 'email' },
  { label: t('chatbot.editor.scenarios.form.fieldTypes.number'), value: 'number' },
  { label: t('chatbot.editor.scenarios.form.fieldTypes.textarea'), value: 'textarea' },
])

function defaultScenarioConfig(type) {
  switch (type) {
    case 'static_message':
      return { message: '', isMarkdown: false, autoAdvanceMs: 500 }
    case 'conversation':
      return { placeholder: '', initialPrompt: '', typingIndicator: true }
    case 'form':
      return { fields: [], submitLabel: '' }
    default:
      return {}
  }
}

function ensureScenarioConfig() {
  if (!current.value) return
  current.value.config = {
    ...defaultScenarioConfig(current.value.type),
    ...(current.value.config || {}),
  }
}

function normalizeCurrent() {
  if (!current.value) return
  if (!current.value.type) current.value.type = 'static_message'
  ensureScenarioConfig()
}

const AUTO_NAME_PATTERN = /^(step|scenario)[_ ]\d+$/i
const typeLabelMap = computed(() =>
  Object.fromEntries(scenarioTypes.value.map(st => [st.value, st.label])),
)
const typeLabelSet = computed(() => new Set(Object.values(typeLabelMap.value)))

function isAutoName(name) {
  return !name || AUTO_NAME_PATTERN.test(name) || typeLabelSet.value.has(name)
}

watch(() => selectedIdx.value, normalizeCurrent, { immediate: true })
watch(
  () => current.value?.type,
  (newType, oldType) => {
    if (!current.value) return
    ensureScenarioConfig()
    if (newType && newType !== oldType && isAutoName(current.value.name)) {
      current.value.name = typeLabelMap.value[newType] || current.value.name
    }
  },
)

function nextStepNumber() {
  const taken = new Set(props.scenarios.map(s => s.id))
  let n = props.scenarios.length + 1
  while (taken.has(`step_${n}`)) n++
  return n
}

function addScenario() {
  const n = nextStepNumber()
  props.scenarios.push({
    id: `step_${n}`,
    name: t('chatbot.editor.scenarios.stepName', { n }),
    type: 'static_message',
    config: defaultScenarioConfig('static_message'),
    automation: emptyAutomation(),
  })
  selectedIdx.value = props.scenarios.length - 1
}

const handleIsDuplicate = computed(() => {
  if (!current.value?.id) return false
  return props.scenarios.some((s, i) => i !== selectedIdx.value && s.id === current.value.id)
})

const draggedIdx = ref(null)
const dropTargetIdx = ref(null)

function onDragStart(idx) {
  draggedIdx.value = idx
}

function onDragOver(idx) {
  dropTargetIdx.value = idx
}

function onDragLeave() {
  dropTargetIdx.value = null
}

function onDrop(idx) {
  const from = draggedIdx.value
  draggedIdx.value = null
  dropTargetIdx.value = null
  if (from === null || from === idx) return
  const selectedID = current.value?.id
  const [moved] = props.scenarios.splice(from, 1)
  props.scenarios.splice(idx, 0, moved)
  if (selectedID) {
    const newIdx = props.scenarios.findIndex(s => s.id === selectedID)
    if (newIdx >= 0) selectedIdx.value = newIdx
  }
}

function removeScenario(idx) {
  props.scenarios.splice(idx, 1)
  if (selectedIdx.value >= props.scenarios.length) {
    selectedIdx.value = props.scenarios.length - 1
  }
}

function openAgent(agentID) {
  window.open(`/agentic/${agentID}/edit`, '_blank', 'noopener')
}

function createAgent() {
  window.open('/agentic/create', '_blank', 'noopener')
}
</script>
