<script setup lang="ts">
import type { App as McpApp } from '@modelcontextprotocol/ext-apps'
import CFieldBoolEditor from '@planetcrust/human-vue/src/components/field/editors/CFieldBoolEditor.vue'
import CFieldDateTimeEditor from '@planetcrust/human-vue/src/components/field/editors/CFieldDateTimeEditor.vue'
import CFieldEmailEditor from '@planetcrust/human-vue/src/components/field/editors/CFieldEmailEditor.vue'
import CFieldNumberEditor from '@planetcrust/human-vue/src/components/field/editors/CFieldNumberEditor.vue'
import CFieldSelectEditor from '@planetcrust/human-vue/src/components/field/editors/CFieldSelectEditor.vue'
import CFieldUrlEditor from '@planetcrust/human-vue/src/components/field/editors/CFieldUrlEditor.vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import MultiSelect from 'primevue/multiselect'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { mcpAppKey } from '../shared/mount'
import {
  type Choice,
  type Draft,
  type FormField,
  type FormValue,
  fromLines,
  initialValues,
  isReadOnly,
  linesOf,
  missingRequired,
  parseDraft,
  savedContext,
  valuesToSave,
} from './recordForm'

const { t } = useI18n()
const bridge = inject<McpApp>(mcpAppKey)!

// The webapp's editors for the kinds that edit without reaching Human's API.
// String is edited here in plain inputs; see isReadOnly for rich text.
const editors: Record<string, unknown> = {
  Bool: CFieldBoolEditor,
  DateTime: CFieldDateTimeEditor,
  Email: CFieldEmailEditor,
  Number: CFieldNumberEditor,
  Select: CFieldSelectEditor,
  Url: CFieldUrlEditor,
}

const draft = ref<Draft>()
const form = ref<Record<string, FormValue>>({})
const saving = ref(false)
const error = ref('')
const saved = ref<{ recordID: string; url?: string }>()

bridge.ontoolresult = r => {
  draft.value = parseDraft(r)
  form.value = initialValues(draft.value)
  saved.value = undefined
  error.value = ''
}

const fields = computed(() => draft.value?.view?.fields ?? [])

// The shape the webapp's editors read a field in.
function editorField(f: FormField) {
  return {
    name: f.name,
    kind: f.kind,
    isMulti: !!f.multi,
    isRequired: !!f.required,
    options: f.options ?? {},
  }
}

// A reference field's choices, with whatever it already holds kept in the list
// even when it is not among the ones the draft listed.
function choicesOf(f: FormField): Choice[] {
  const listed = draft.value?.view?.choices?.[f.name]?.items ?? []
  const held = [form.value[f.name]].flat().filter(Boolean) as string[]
  const extra = held
    .filter(id => !listed.some(c => c.id === id))
    .map(id => ({ id, label: draft.value!.refs[id] ?? id }))
  return [...extra, ...listed]
}

function hasMore(f: FormField) {
  return !!draft.value?.view?.choices?.[f.name]?.more
}

// A read-only field is shown by how many values it holds, or by its plain text
// when it is rich text.
function readOnlyText(f: FormField) {
  const v = [draft.value?.values[f.name] ?? []].flat().filter(Boolean)
  if (f.kind === 'String') {
    const doc = new DOMParser().parseFromString(v.join(' '), 'text/html')
    return doc.body.textContent?.trim() || '—'
  }
  return v.length ? String(v.length) : '—'
}

function resultText(r: { content?: { type: string; text?: string }[] }) {
  return (r.content ?? []).map(c => c.text ?? '').join(' ')
}

async function save() {
  if (!draft.value || saving.value) return

  const missing = missingRequired(draft.value, form.value)
  if (missing.length) {
    error.value = t('mcpApp.recordForm.missing', { fields: missing.join(', ') })
    return
  }

  const values = valuesToSave(draft.value, form.value)
  const args: Record<string, unknown> = {
    namespace: draft.value.namespaceID,
    module: draft.value.moduleID,
    values: JSON.stringify(values),
  }
  if (draft.value.recordID) args.recordID = draft.value.recordID

  saving.value = true
  error.value = ''
  try {
    const r = await bridge.callServerTool({
      name: draft.value.recordID ? 'compose_record_update' : 'compose_record_create',
      arguments: args,
    })
    if (r.isError) throw new Error(resultText(r) || 'error')

    const rec = JSON.parse(resultText(r) || '{}')
    saved.value = {
      recordID: rec.recordID ?? draft.value.recordID,
      url: rec.url ?? draft.value.view?.link,
    }

    await bridge.updateModelContext({
      content: [{ type: 'text', text: savedContext(draft.value, saved.value.recordID, values) }],
    })
  } catch (e) {
    error.value = t('mcpApp.recordForm.saveFailed', {
      error: e instanceof Error ? e.message : String(e),
    })
  } finally {
    saving.value = false
  }
}

function open(url?: string) {
  if (url) bridge.openLink({ url })
}
</script>

<template>
  <div class="p-3 text-sm">
    <div v-if="!draft" id="status" class="text-muted-color">
      {{ t('mcpApp.recordForm.waiting') }}
    </div>

    <template v-else>
      <div class="flex items-baseline justify-between gap-3 mb-3">
        <div class="font-semibold truncate">{{ draft.view?.module.name }}</div>
        <div id="status" class="text-muted-color text-xs shrink-0">
          {{ draft.recordID ? t('mcpApp.recordForm.edit') : t('mcpApp.recordForm.new') }}
        </div>
      </div>

      <form class="flex flex-col gap-3" @submit.prevent="save">
        <div v-for="f in fields" :key="f.name" class="flex flex-col gap-1">
          <label :for="`field-${f.name}`" class="text-xs font-medium">
            {{ f.label || f.name }}
            <span v-if="f.required" class="text-red-500">*</span>
          </label>

          <div v-if="isReadOnly(f)" class="text-muted-color">
            {{ readOnlyText(f) }} · {{ t('mcpApp.recordForm.readOnly') }}
          </div>

          <template v-else-if="f.kind === 'Record' || f.kind === 'User'">
            <MultiSelect
              v-if="f.multi"
              :id="`field-${f.name}`"
              v-model="form[f.name] as string[]"
              :options="choicesOf(f)"
              option-label="label"
              option-value="id"
              filter
              :disabled="!!saved"
              class="w-full"
            />
            <Select
              v-else
              :id="`field-${f.name}`"
              v-model="form[f.name] as string"
              :options="choicesOf(f)"
              option-label="label"
              option-value="id"
              filter
              show-clear
              :disabled="!!saved"
              class="w-full"
            />
            <div v-if="hasMore(f)" class="text-muted-color text-xs">
              {{ t('mcpApp.recordForm.moreChoices') }}
            </div>
          </template>

          <template v-else-if="f.kind === 'String'">
            <Textarea
              v-if="f.multi"
              :id="`field-${f.name}`"
              :model-value="linesOf(form[f.name])"
              :disabled="!!saved"
              auto-resize
              rows="2"
              class="w-full"
              @update:model-value="form[f.name] = fromLines($event ?? '')"
            />
            <Textarea
              v-else-if="f.options?.multiLine"
              :id="`field-${f.name}`"
              v-model="form[f.name] as string"
              :disabled="!!saved"
              auto-resize
              rows="3"
              class="w-full"
            />
            <InputText
              v-else
              :id="`field-${f.name}`"
              v-model="form[f.name] as string"
              :disabled="!!saved"
              class="w-full"
            />
          </template>

          <component
            :is="editors[f.kind]"
            v-else-if="editors[f.kind]"
            :id="`field-${f.name}`"
            v-model="form[f.name]"
            :field="editorField(f)"
            :disabled="!!saved"
          />
        </div>

        <div class="flex items-center gap-3 mt-1">
          <Button
            v-if="!saved"
            type="submit"
            :label="saving ? t('mcpApp.recordForm.saving') : t('mcpApp.recordForm.save')"
            :loading="saving"
            size="small"
          />
          <template v-else>
            <span class="text-green-600 font-medium">{{ t('mcpApp.recordForm.saved') }}</span>
            <Button
              v-if="saved.url"
              :label="t('mcpApp.recordForm.open')"
              size="small"
              severity="secondary"
              text
              @click="open(saved.url)"
            />
          </template>
          <span v-if="!saved && !error" class="text-muted-color text-xs">
            {{ t('mcpApp.recordForm.notSaved') }}
          </span>
          <span v-if="error" class="text-red-500 text-xs">{{ error }}</span>
        </div>
      </form>
    </template>
  </div>
</template>
