<template>
  <div v-if="filter.params && filter.params.length" class="flex flex-col gap-4">
    <CFormGroup
      v-for="(param, index) in filter.params"
      :key="index"
    >
      <template #label>
        {{ getParamLabel(param.label) }}
        <a
          v-if="param.label === 'expr'"
          href="https://docs.planetcrust.io/human-docs/latest/integrator-guide/expr/index.html"
          target="_blank"
          class="ml-1 text-muted-color"
        >
          <i class="pi pi-question-circle text-xs" />
        </a>
      </template>
      <template v-if="filter.ref === 'response' && param.type === 'header'" #actions>
        <Button
          :label="$t('system.apigw.editor.filters.addHeader')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          @click="getHeaderValue(param).push({ name: '', expr: '' })"
        />
      </template>

      <!-- Boolean -->
      <ToggleSwitch
        v-if="param.type === 'bool'"
        v-model="param.value"
      />

      <!-- Workflow picker -->
      <Select
        v-else-if="param.label === 'workflow'"
        v-model="param.value"
        :options="workflows"
        option-label="label"
        option-value="workflowID"
        :placeholder="$t('system.apigw.editor.filters.placeholders.workflow')"
        :loading="loadingWorkflows"
        showClear
        filter
      />

      <!-- HTTP Status -->
      <Select
        v-else-if="param.label === 'status'"
        v-model="param.value"
        :options="httpStatusOptions"
        option-label="label"
        option-value="value"
        showClear
      />

      <!-- Response filter: input type -->
      <template v-else-if="filter.ref === 'response' && param.type === 'input'">
        <Select
          v-model="getInputValue(param).type"
          :options="inputTypeOptions"
          class="mb-2"
        />
        <InputGroup>
          <InputGroupAddon>ƒ</InputGroupAddon>
          <InputText
            v-model="getInputValue(param).expr"
            :placeholder="$t('system.apigw.editor.filters.help.expression.example')"
          />
        </InputGroup>
      </template>

      <!-- Response filter: header type -->
      <CFormList
        v-else-if="filter.ref === 'response' && param.type === 'header'"
        :model-value="getHeaderValue(param)"
        :columns="[
          { label: $t('system.apigw.editor.filters.labels.name'), width: '1fr' },
          { label: $t('system.apigw.editor.filters.labels.value'), width: '1fr' },
        ]"
      >
        <template #row="{ item }">
          <InputText
            v-model="item.name"
            :placeholder="$t('system.apigw.editor.filters.labels.name')"
            class="w-full"
          />
          <InputText
            v-model="item.expr"
            :placeholder="$t('system.apigw.editor.filters.labels.value')"
            class="w-full"
          />
        </template>
      </CFormList>

      <!-- JS function -->
      <Textarea
        v-else-if="param.label === 'jsfunc'"
        v-model="param.value"
        rows="6"
        autoResize
        class="font-mono text-sm"
      />

      <!-- Expression -->
      <InputGroup v-else-if="param.label === 'expr'">
        <InputGroupAddon>ƒ</InputGroupAddon>
        <InputText
          v-model="param.value"
          :placeholder="$t('system.apigw.editor.filters.help.expression.example')"
        />
      </InputGroup>

      <!-- Default text input -->
      <InputText
        v-else
        v-model="param.value"
      />
    </CFormGroup>
  </div>
</template>

<script setup>
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, te } = useI18n()
const $AutomationAPI = inject('$AutomationAPI')

const props = defineProps({
  filter: {
    type: Object,
    required: true,
  },
})

const loadingWorkflows = ref(false)
const workflows = ref([])

const httpStatusOptions = [
  { value: 300, label: t('system.apigw.editor.filters.httpStatus.300') },
  { value: 301, label: t('system.apigw.editor.filters.httpStatus.301') },
  { value: 302, label: t('system.apigw.editor.filters.httpStatus.302') },
  { value: 303, label: t('system.apigw.editor.filters.httpStatus.303') },
  { value: 304, label: t('system.apigw.editor.filters.httpStatus.304') },
  { value: 307, label: t('system.apigw.editor.filters.httpStatus.307') },
  { value: 308, label: t('system.apigw.editor.filters.httpStatus.308') },
]

const inputTypeOptions = [
  'String',
  'Any',
  'Array',
  'KV',
  'DateTime',
  'Float',
  'Integer',
  'Reader',
  'Vars',
]

// Safe translation fallback — uses the translation key if it exists,
// otherwise humanizes the raw label (e.g. "jsfunc" → "Jsfunc")
function getParamLabel(label) {
  const key = `system.apigw.editor.filters.labels.${label}`
  return te(key) ? t(key) : (label || '').replace(/([A-Z])/g, ' $1').replace(/^./, s => s.toUpperCase())
}

// Null-safe getter for response input param value
function getInputValue(param) {
  if (!param.value || typeof param.value !== 'object') {
    param.value = { type: 'Any', expr: '' }
  }
  return param.value
}

// Null-safe getter for response header param value
function getHeaderValue(param) {
  if (!Array.isArray(param.value)) {
    param.value = []
  }
  return param.value
}

async function loadWorkflows() {
  const needsWorkflows = (props.filter.params || []).some(p => p.label === 'workflow')
  if (!needsWorkflows) return

  loadingWorkflows.value = true
  try {
    const result = await $AutomationAPI.workflowList()
    workflows.value = (result?.set || []).map(({ workflowID, handle, meta }) => ({
      label: meta?.name || handle || workflowID,
      workflowID,
    }))
  } catch (e) {
    console.error('Failed to load workflows:', e)
  } finally {
    loadingWorkflows.value = false
  }
}

watch(() => props.filter, () => loadWorkflows(), { immediate: true })
</script>
