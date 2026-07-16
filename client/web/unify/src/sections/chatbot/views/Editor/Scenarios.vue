<template>
  <Panel :header="$t('chatbot.editor.panels.scenarios')" toggleable>
    <p class="text-xs text-muted-color mb-3">
      {{ $t('chatbot.editor.scenarios.description') }}
    </p>

    <div class="grid grid-cols-1 xl:grid-cols-3 gap-4 min-h-80 pb-3">
      <Fieldset class="lg:col-span-1">
        <template #legend>
          <div class="flex items-center gap-2">
            <span>{{ $t('chatbot.editor.scenarios.list') }}</span>
            <Button
              :label="$t('chatbot.editor.scenarios.add')"
              icon="pi pi-plus"
              severity="secondary"
              size="small"
              @click="addScenario"
            />
          </div>
        </template>
        <CFormItemList
          :items="props.scenarios"
          :empty-message="$t('chatbot.editor.scenarios.empty')"
          item-key="id"
          :selected-key="current?.id"
          draggable
          @select="(_item, index) => (selectedIdx = index)"
          @remove="onScenarioRemove"
          @reorder="onScenariosReorder"
        >
          <template #default="{ item, index }">
            <span
              class="block w-full text-sm truncate"
              :class="selectedIdx === index ? 'text-primary font-medium' : ''"
            >
              {{ item.name || item.id || $t('chatbot.editor.scenarios.unnamed') }}
            </span>
          </template>
        </CFormItemList>
      </Fieldset>

      <div class="xl:col-span-2 flex flex-col gap-3 border border-surface rounded-lg p-3 mt-6">
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

          <Fieldset :legend="$t(`chatbot.editor.scenarios.types.${current.type}`)" class="mb-3">
            <div class="grid grid-cols-1 xl:grid-cols-2 gap-3">
              <template v-if="current.type === 'static_message'">
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.static.message')"
                  class="xl:col-span-2"
                  required
                >
                  <CRichTextInput
                    v-model="current.config.message"
                    min-body-height="6rem"
                    max-body-height="14rem"
                  />
                </CFormGroup>
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
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.conversation.placeholder')"
                  :description="$t('chatbot.editor.scenarios.conversation.placeholderHelp')"
                  class="xl:col-span-2"
                >
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
                  class="col-span-2"
                />
              </template>

              <template v-else-if="current.type === 'form'">
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.form.fields')"
                  class="xl:col-span-2"
                >
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
                      { label: $t('chatbot.editor.scenarios.form.fieldLabel'), width: '3fr' },
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

              <template v-else-if="current.type === 'consent'">
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.consent.body.label')"
                  :description="$t('chatbot.editor.scenarios.consent.body.help')"
                  class="xl:col-span-2"
                  required
                >
                  <CRichTextInput
                    v-model="current.config.body"
                    :placeholder="$t('chatbot.editor.scenarios.consent.body.placeholder')"
                    min-body-height="6rem"
                    max-body-height="14rem"
                  />
                </CFormGroup>
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.consent.acceptLabel.label')"
                  :description="$t('chatbot.editor.scenarios.consent.acceptLabel.help')"
                >
                  <InputText
                    v-model="current.config.acceptLabel"
                    :placeholder="$t('chatbot.editor.scenarios.consent.acceptLabel.placeholder')"
                  />
                </CFormGroup>
                <CFormGroup
                  :label="$t('chatbot.editor.scenarios.consent.rejectLabel.label')"
                  :description="$t('chatbot.editor.scenarios.consent.rejectLabel.help')"
                >
                  <InputText
                    v-model="current.config.rejectLabel"
                    :placeholder="$t('chatbot.editor.scenarios.consent.rejectLabel.placeholder')"
                  />
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
                :placeholder="$t('chatbot.editor.scenarios.automation.placeholder')"
                @update:model-value="setAutomation('before', $event)"
              />
            </CFormGroup>
            <CFormGroup
              :label="$t('chatbot.editor.scenarios.automation.after')"
              :description="$t('chatbot.editor.scenarios.automation.afterHelp')"
            >
              <CInputTAQ
                :model-value="automationID(current.automation?.after?.automation)"
                :placeholder="$t('chatbot.editor.scenarios.automation.placeholder')"
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

    <Divider />

    <Fieldset :legend="$t('chatbot.editor.handoff.sectionTitle')">
      <div class="flex flex-col gap-4">
        <CInputToggleCard
          v-model="handoff.enabled"
          :label="$t('chatbot.editor.handoff.enabled.label')"
          :description="$t('chatbot.editor.handoff.enabled.help')"
        />

        <template v-if="handoff.enabled">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
            <CFormGroup
              :label="$t('chatbot.editor.handoff.automation.onRequested.label')"
              :description="$t('chatbot.editor.handoff.automation.onRequested.help')"
            >
              <CInputTAQ
                :model-value="automationID(handoff.automation?.onRequested?.automation)"
                :placeholder="$t('chatbot.editor.handoff.automation.placeholder')"
                @update:model-value="setHandoffAutomation('onRequested', $event)"
              />
            </CFormGroup>
            <CFormGroup
              :label="$t('chatbot.editor.handoff.automation.onAccepted.label')"
              :description="$t('chatbot.editor.handoff.automation.onAccepted.help')"
            >
              <CInputTAQ
                :model-value="automationID(handoff.automation?.onAccepted?.automation)"
                :placeholder="$t('chatbot.editor.handoff.automation.placeholder')"
                @update:model-value="setHandoffAutomation('onAccepted', $event)"
              />
            </CFormGroup>
          </div>
        </template>
      </div>
    </Fieldset>
  </Panel>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CInputTAQ, CInputToggleCard, CRichTextInput } = components
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
  { label: t('chatbot.editor.scenarios.types.consent'), value: 'consent' },
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
      return { message: '', autoAdvanceMs: 500 }
    case 'conversation':
      return { placeholder: '', initialPrompt: '', typingIndicator: true }
    case 'form':
      return { fields: [], submitLabel: '' }
    case 'consent':
      return { body: '', acceptLabel: '', rejectLabel: '' }
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

// CFormItemList passes items as a one-way prop, so we mutate the parent's
// scenarios array in place on reorder/remove (keeping the reference stable
// for the Editor's chatbot.scenarios binding). Selection is re-mapped by
// scenario ID so the editor stays on the same step after reorder.
function onScenariosReorder(reordered) {
  const selectedID = current.value?.id
  props.scenarios.splice(0, props.scenarios.length, ...reordered)
  if (!selectedID) return
  const newIdx = props.scenarios.findIndex(s => s.id === selectedID)
  if (newIdx >= 0) selectedIdx.value = newIdx
}

function onScenarioRemove(_item, index) {
  props.scenarios.splice(index, 1)
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
