<template>
  <div class="flex flex-col gap-3">
    <CFormGroup :label="$t('block.custom.modeLabel')">
      <SelectButton
        v-model="mode"
        :options="modes"
        option-label="label"
        option-value="value"
        :allow-empty="false"
        data-test-id="custom-block-mode"
      />
    </CFormGroup>

    <CFormGroup
      v-if="mode === 'app'"
      :label="$t('block.custom.applicationLabel')"
      :description="$t('block.custom.applicationDesc')"
      required
    >
      <Select
        v-model="applicationID"
        :options="applications"
        option-label="label"
        option-value="applicationID"
        :placeholder="$t('block.custom.pickApplication')"
        :loading="loadingApps"
        filter
        class="w-full"
        data-test-id="custom-block-application"
      />
    </CFormGroup>

    <template v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <CustomAppDeclaration v-model="declaration" :namespace="namespace" />
      </div>

      <CFormGroup
        :label="$t('block.custom.sourceLabel')"
        :description="$t('block.custom.sourceDesc')"
        required
      >
        <template #actions>
          <CustomAppSnippets
            :modules="declaration.modules"
            @insert="sourceEditor?.insert($event)"
          />
        </template>
        <div data-test-id="custom-block-source">
          <CCodeEditor
            id="customBlockSource"
            ref="sourceEditor"
            v-model="source"
            language="html"
            min-height="360px"
            :assist="assist"
          />
        </div>
      </CFormGroup>
    </template>

    <CFormGroup
      :label="$t('block.custom.paramsLabel')"
      :description="$t('block.custom.paramsDesc')"
    >
      <template #actions>
        <Button
          :label="$t('general.label.add')"
          icon="pi pi-plus"
          severity="secondary"
          size="small"
          data-test-id="custom-block-params-add"
          @click="paramRows.push({ name: '', value: '' })"
        />
      </template>
      <CFormList
        v-model="paramRows"
        data-test-id="custom-block-params"
        :columns="[
          { label: $t('block.custom.paramName'), width: '1fr' },
          { label: $t('block.custom.paramValue'), width: '2fr' },
        ]"
        :empty-message="$t('block.custom.noParams')"
        fit-width
      >
        <template #row="{ index }">
          <InputText
            v-model="paramRows[index].name"
            size="small"
            class="w-full"
            placeholder="status"
          />
          <InputText
            v-model="paramRows[index].value"
            size="small"
            class="w-full"
            placeholder="doing"
          />
        </template>
      </CFormList>
    </CFormGroup>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import CustomAppDeclaration from '@/sections/app/components/CustomAppDeclaration.vue'
import CustomAppSnippets from '@/sections/app/components/CustomAppSnippets.vue'
import { bridgeAssist } from '@/sections/app/hints'

const { CCodeEditor } = components

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const { t } = useI18n()
const block = inject('blockDraft')
const $SystemAPI = inject('$SystemAPI')
const $ComposeAPI = inject('$ComposeAPI')

const sourceEditor = ref(null)
const applications = ref([])
const loadingApps = ref(false)

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

const modes = computed(() => [
  { value: 'inline', label: t('block.custom.modeInline') },
  { value: 'app', label: t('block.custom.modeApp') },
])

// Picked once and kept on the draft as which option is filled in; switching
// clears the other, so a block never carries both.
const chosen = ref(null)
const mode = computed({
  get: () =>
    chosen.value || (isApplicationID(block.value.options?.applicationID) ? 'app' : 'inline'),
  set: v => {
    chosen.value = v
    if (v === 'app') {
      updateOptions('source', '')
      declaration.value = {}
    } else {
      updateOptions('applicationID', '')
      declareWhereShown()
    }
  },
})

function isApplicationID(id) {
  return !!id && id !== '0'
}

// A page written into the block reads this page's namespace, and starts out
// declaring the module of the record page it sits on. Only an empty
// declaration is filled, so nothing an author chose is replaced.
async function declareWhereShown() {
  const moduleID = props.page?.moduleID
  if (block.value.options?.modules?.length || !moduleID || moduleID === '0') return

  const module = await $ComposeAPI
    .moduleRead({ namespaceID: props.namespace.namespaceID, moduleID })
    .catch(() => null)

  if (!module?.handle || block.value.options?.modules?.length) return
  declaration.value = { ...declaration.value, modules: [module.handle], writes: [] }
}

const applicationID = computed({
  get: () => block.value.options?.applicationID || '',
  set: v => updateOptions('applicationID', v || ''),
})

const source = computed({
  get: () => block.value.options?.source || '',
  set: v => updateOptions('source', v),
})

// The modules the page may read and change, by handle, in this page's
// namespace.
const declaration = computed({
  get: () => ({
    modules: block.value.options?.modules || [],
    writes: block.value.options?.writes || [],
    deletes: block.value.options?.deletes || [],
    origins: block.value.options?.origins || [],
    automations: block.value.options?.automations || [],
    chatbots: block.value.options?.chatbots || [],
  }),
  set: d => {
    updateOptions('chatbots', d.chatbots || [])
    updateOptions('automations', d.automations || [])
    updateOptions('origins', d.origins || [])
    updateOptions('modules', d.modules || [])
    updateOptions('deletes', d.deletes || [])
    updateOptions('writes', d.writes || [])
  },
})

// The parameters as rows of name and value. A value is read as JSON where it
// parses — a number, true or false, a list — and kept as text otherwise; a
// row with no name is one still being written and is left out. Nothing is
// written until the author changes a row.
function shown(value) {
  return typeof value === 'string' ? value : JSON.stringify(value)
}

function parsed(text) {
  try {
    return JSON.parse(text)
  } catch {
    return text
  }
}

const paramRows = ref(
  Object.entries(block.value.options?.params || {}).map(([name, value]) => ({
    name,
    value: shown(value),
  })),
)

watch(
  paramRows,
  rows => {
    const params = {}
    for (const { name, value } of rows) {
      if (name.trim()) params[name.trim()] = parsed(value)
    }
    updateOptions('params', params)
  },
  { deep: true },
)

// What the editor suggests: what this block declared and its parameters.
const assist = bridgeAssist(t, () => ({
  ...declaration.value,
  params: paramRows.value.map(r => r.name.trim()).filter(Boolean),
}))

onMounted(async () => {
  // A block opened for the first time has neither an application nor a page.
  if (mode.value === 'inline' && !block.value.options?.source) declareWhereShown()

  loadingApps.value = true
  try {
    const { set = [] } = await $SystemAPI.applicationList({ limit: 200 })
    applications.value = set
      .filter(app => app.unify?.kind === 'custom')
      .map(app => ({ applicationID: app.applicationID, label: app.unify?.name || app.name }))
  } finally {
    loadingApps.value = false
  }
})
</script>
