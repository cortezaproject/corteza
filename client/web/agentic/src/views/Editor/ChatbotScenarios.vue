<template>
  <Panel :header="$t('agent.editor.chatbot.panels.scenarios')" toggleable>
    <p class="text-xs text-muted-color mb-3">
      {{ $t('agent.editor.chatbot.scenarios.description') }}
    </p>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 min-h-[320px]">
      <!-- list -->
      <div class="md:col-span-1 flex flex-col gap-2 border border-surface rounded-lg p-2">
        <div class="flex items-center justify-between">
          <span class="font-medium text-sm">
            {{ $t('agent.editor.chatbot.scenarios.list') }}
          </span>
          <Button
            icon="pi pi-plus"
            text
            size="small"
            @click="addScenario"
            v-tooltip.bottom="$t('agent.editor.chatbot.scenarios.add')"
          />
        </div>
        <div v-if="!scenarios.length" class="text-xs text-muted-color p-2">
          {{ $t('agent.editor.chatbot.scenarios.empty') }}
        </div>
        <TransitionGroup
          tag="ul"
          class="flex flex-col gap-1"
          move-class="transition-transform duration-300 ease-in-out"
        >
          <li
            v-for="(s, idx) in scenarios"
            :key="s.id || idx"
            class="flex items-center gap-2 p-2 rounded cursor-grab"
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
            <i class="pi pi-bars text-xs text-muted-color" />
            <span class="flex-1 text-sm truncate">
              {{ s.name || s.id || $t('agent.editor.chatbot.scenarios.unnamed') }}
            </span>
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              size="small"
              @click.stop="removeScenario(idx)"
            />
          </li>
        </TransitionGroup>
      </div>

      <!-- detail -->
      <div class="md:col-span-2 flex flex-col gap-3 border border-surface rounded-lg p-3">
        <template v-if="current">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.name') }}
              </label>
              <InputText v-model="current.name" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.handle') }}
              </label>
              <InputText v-model="current.id" />
            </div>
          </div>

          <div class="flex flex-col gap-1">
            <label class="font-medium text-primary">
              {{ $t('agent.editor.chatbot.scenarios.type') }}
            </label>
            <Select
              v-model="current.type"
              :options="scenarioTypes"
              optionLabel="label"
              optionValue="value"
              @change="ensureScenarioConfig"
            />
          </div>

          <!-- type: static_message -->
          <template v-if="current.type === 'static_message'">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.static.message') }}
              </label>
              <Textarea v-model="current.config.message" :rows="3" />
            </div>
            <div class="flex items-center gap-2">
              <ToggleSwitch v-model="current.config.isMarkdown" />
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.static.isMarkdown') }}
              </label>
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.static.autoAdvance') }}
              </label>
              <InputNumber v-model="current.config.autoAdvanceMs" :min="0" />
            </div>
          </template>

          <!-- type: conversation -->
          <template v-else-if="current.type === 'conversation'">
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.conversation.placeholder') }}
              </label>
              <InputText v-model="current.config.placeholder" />
            </div>
            <div class="flex flex-col gap-1">
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.conversation.initialPrompt') }}
              </label>
              <Textarea v-model="current.config.initialPrompt" :rows="2" />
            </div>
            <div class="flex items-center gap-2">
              <ToggleSwitch v-model="current.config.typingIndicator" />
              <label class="font-medium text-primary">
                {{ $t('agent.editor.chatbot.scenarios.conversation.typingIndicator') }}
              </label>
            </div>
          </template>

          <!-- type: form -->
          <template v-else-if="current.type === 'form'">
            <div class="flex flex-col gap-2">
              <label class="font-medium text-primary flex items-center justify-between">
                <span>{{ $t('agent.editor.chatbot.scenarios.form.fields') }}</span>
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
                    {{ $t('agent.editor.chatbot.scenarios.form.fieldName') }}
                  </span>
                  <span class="col-span-4 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('agent.editor.chatbot.scenarios.form.fieldLabel') }}
                  </span>
                  <span class="col-span-3 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('agent.editor.chatbot.scenarios.form.fieldType') }}
                  </span>
                  <span class="col-span-1 text-xs uppercase tracking-wide text-muted-color">
                    {{ $t('agent.editor.chatbot.scenarios.form.fieldRequired') }}
                  </span>
                  <span class="col-span-1" />
                </div>
                <div
                  v-if="!current.config.fields.length"
                  class="text-xs text-muted-color text-center py-2"
                >
                  {{ $t('agent.editor.chatbot.scenarios.form.empty') }}
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
                {{ $t('agent.editor.chatbot.scenarios.form.submitLabel') }}
              </label>
              <InputText v-model="current.config.submitLabel" />
            </div>
          </template>
        </template>
        <div v-else class="text-sm text-muted-color p-4 text-center">
          {{ $t('agent.editor.chatbot.scenarios.selectHint') }}
        </div>
      </div>
    </div>
  </Panel>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  scenarios: { type: Array, required: true },
})

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
    case 'static_message': return { message: '', isMarkdown: false, autoAdvanceMs: 1500 }
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

const AUTO_NAME_PATTERN = /^scenario_\d+$/i
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

function addScenario() {
  const id = `scenario_${props.scenarios.length + 1}`
  props.scenarios.push({
    id,
    name: id,
    type: 'static_message',
    config: defaultScenarioConfig('static_message'),
  })
  selectedIdx.value = props.scenarios.length - 1
}

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
</script>
