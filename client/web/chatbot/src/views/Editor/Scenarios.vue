<template>
  <Panel :header="$t('chatbot.editor.panels.scenarios')" toggleable>
    <p class="text-xs text-muted-color mb-3">
      {{ $t('chatbot.editor.scenarios.description') }}
    </p>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 min-h-[320px]">
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

      <div class="lg:col-span-2 flex flex-col gap-3 border border-surface rounded-lg p-3">
        <template v-if="current">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.name') }}
              </label>
              <InputText v-model="current.name" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.handle') }}
              </label>
              <InputText v-model="current.id" :invalid="handleIsDuplicate" />
              <small v-if="handleIsDuplicate" class="text-red-500">
                {{ $t('chatbot.editor.scenarios.handleDuplicate') }}
              </small>
            </div>
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-primary">
              {{ $t('chatbot.editor.scenarios.type') }}
            </label>
            <Select
              v-model="current.type"
              :options="scenarioTypes"
              optionLabel="label"
              optionValue="value"
              @change="ensureScenarioConfig"
            />
          </div>

          <template v-if="current.type === 'static_message'">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.static.message') }}
              </label>
              <Textarea v-model="current.config.message" :rows="3" />
            </div>
            <div class="flex items-center gap-2">
              <ToggleSwitch v-model="current.config.isMarkdown" />
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.static.isMarkdown') }}
              </label>
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.static.autoAdvance') }}
              </label>
              <InputNumber v-model="current.config.autoAdvanceMs" :min="0" />
            </div>
          </template>

          <template v-else-if="current.type === 'conversation'">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.agent.label') }}
                <span class="text-red-500">*</span>
              </label>
              <small class="text-muted-color">
                {{ $t('chatbot.editor.scenarios.agent.help') }}
              </small>
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
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.conversation.placeholder') }}
              </label>
              <InputText v-model="current.config.placeholder" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.conversation.initialPrompt') }}
              </label>
              <Textarea v-model="current.config.initialPrompt" :rows="2" />
            </div>
            <div class="flex items-center gap-2">
              <ToggleSwitch v-model="current.config.typingIndicator" />
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.conversation.typingIndicator') }}
              </label>
            </div>
          </template>

          <template v-else-if="current.type === 'form'">
            <div class="flex flex-col gap-2">
              <label class="font-medium text-primary flex items-center justify-between">
                <span>{{ $t('chatbot.editor.scenarios.form.fields') }}</span>
                <Button
                  icon="pi pi-plus"
                  text
                  size="small"
                  @click="current.config.fields.push({ name:'', label:'', type:'text', required:false })"
                />
              </label>
              <div class="border border-surface rounded-lg p-3 flex flex-col gap-2">
                <div
                  v-if="current.config.fields.length"
                  class="grid grid-cols-12 gap-2 items-center"
                >
                  <span class="col-span-3 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('chatbot.editor.scenarios.form.fieldName') }}
                  </span>
                  <span class="col-span-4 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('chatbot.editor.scenarios.form.fieldLabel') }}
                  </span>
                  <span class="col-span-3 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('chatbot.editor.scenarios.form.fieldType') }}
                  </span>
                  <span class="col-span-1 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('chatbot.editor.scenarios.form.fieldRequired') }}
                  </span>
                  <span class="col-span-1" />
                </div>
                <div
                  v-if="!current.config.fields.length"
                  class="text-xs text-muted-color text-center py-2"
                >
                  {{ $t('chatbot.editor.scenarios.form.empty') }}
                </div>
                <div
                  v-for="(f, fIdx) in current.config.fields"
                  :key="fIdx"
                  class="grid grid-cols-12 gap-2 items-center"
                >
                  <InputText v-model="f.name" class="col-span-3" />
                  <InputText v-model="f.label" class="col-span-4" />
                  <Select
                    v-model="f.type"
                    :options="['text','email','number','textarea']"
                    class="col-span-3"
                  />
                  <ToggleSwitch v-model="f.required" class="col-span-1" />
                  <Button
                    icon="pi pi-trash"
                    severity="danger"
                    text
                    size="small"
                    class="col-span-1"
                    @click="current.config.fields.splice(fIdx, 1)"
                  />
                </div>
              </div>
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('chatbot.editor.scenarios.form.submitLabel') }}
              </label>
              <InputText v-model="current.config.submitLabel" />
            </div>
          </template>
        </template>
        <div v-else class="text-sm text-muted-color p-4 text-center">
          {{ $t('chatbot.editor.scenarios.selectHint') }}
        </div>
      </div>
    </div>
  </Panel>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  scenarios: { type: Array, required: true },
  agents: { type: Array, default: () => [] },
})

const agentOptions = computed(() =>
  props.agents.map(a => ({
    value: a.agentID,
    label: a.meta?.short || a.handle || a.agentID,
  })),
)

const selectedIdx = ref(props.scenarios.length ? 0 : -1)

const current = computed(() =>
  selectedIdx.value >= 0 ? props.scenarios[selectedIdx.value] : null,
)

const scenarioTypes = [
  { label: 'Static message', value: 'static_message' },
  { label: 'Conversation', value: 'conversation' },
  { label: 'Form', value: 'form' },
]

function defaultScenarioConfig(type) {
  switch (type) {
    case 'static_message': return { message: '', isMarkdown: false, autoAdvanceMs: 500 }
    case 'conversation':   return { placeholder: '', initialPrompt: '', typingIndicator: true }
    case 'form':           return { fields: [], submitLabel: '' }
    default:               return {}
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
const typeLabelMap = Object.fromEntries(scenarioTypes.map(t => [t.value, t.label]))
const typeLabelSet = new Set(Object.values(typeLabelMap))

function isAutoName(name) {
  return !name || AUTO_NAME_PATTERN.test(name) || typeLabelSet.has(name)
}

watch(() => selectedIdx.value, normalizeCurrent, { immediate: true })
watch(() => current.value?.type, (newType, oldType) => {
  if (!current.value) return
  ensureScenarioConfig()
  if (newType && newType !== oldType && isAutoName(current.value.name)) {
    current.value.name = typeLabelMap[newType] || current.value.name
  }
})

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
    name: `Step ${n}`,
    type: 'static_message',
    config: defaultScenarioConfig('static_message'),
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
