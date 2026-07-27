<template>
  <!-- Boolean: checkbox sits inline with its label -->
  <div v-if="field.widget === 'boolean'" class="flex items-center gap-2">
    <Checkbox
      :model-value="field.value"
      :binary="true"
      :inputId="inputId"
      @update:model-value="patch({ value: $event })"
    />
    <label :for="inputId" class="font-medium text-sm">
      {{ field.label }}
      <span v-if="field.required" class="text-red-500">*</span>
    </label>
  </div>

  <CFormGroup
    v-else
    :label="field.label"
    :description="field.description"
    :required="field.required"
    :input-id="inputId"
  >
    <component
      :is="widget.component"
      :id="widget.component === InputText ? inputId : undefined"
      :model-value="field.value"
      v-bind="widget.props"
      class="w-full"
      @update:model-value="patch({ value: $event })"
    />
  </CFormGroup>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import { components } from '@planetcrust/human-vue'
const { CInputNamespace, CInputModule, CInputRecord, CInputUser, CInputRole, CInputDateTime } =
  components

const { t } = useI18n()

const props = defineProps({
  // Field descriptor built by lib/dry-run.js
  field: { type: Object, required: true },
  // name → picked value of the scope fields, for cascading dependencies
  scopeValues: { type: Object, default: () => ({}) },
})

// The owner of the fields array applies the patch onto the field descriptor
const emit = defineEmits(['update'])
const patch = p => emit('update', p)

const inputId = computed(() => `wf-input-${props.field.name}`)

const hasDep = name => (props.field.deps || []).includes(name)
const dep = name => (hasDep(name) ? props.scopeValues[name] : undefined)

// Record picker cascades off the sibling namespace/module scope fields
const namespaceID = computed(() => dep('namespace'))
const moduleID = computed(() => dep('module'))

// A picked value is stale once any sibling field it depends on changes
watch(
  () => (props.field.deps || []).map(name => props.scopeValues[name]),
  (next, prev) => {
    if (next.some((v, i) => v !== prev[i])) patch({ value: undefined })
  },
)

const widget = computed(() => {
  switch (props.field.widget) {
    case 'namespace':
      return { component: CInputNamespace }
    case 'module':
      return {
        component: CInputModule,
        props: {
          namespaceID: dep('namespace') || undefined,
          disabled: !dep('namespace'),
          placeholder: dep('namespace') ? undefined : t('editor.select-namespace-first'),
        },
      }
    case 'record':
      return {
        component: CInputRecord,
        props: {
          namespaceID: namespaceID.value || undefined,
          moduleID: moduleID.value || undefined,
          disabled: !namespaceID.value || !moduleID.value,
        },
      }
    case 'user':
      return { component: CInputUser }
    case 'role':
      return { component: CInputRole }
    case 'datetime':
      return { component: CInputDateTime }
    case 'number':
      return {
        component: InputNumber,
        props: { useGrouping: false, maxFractionDigits: 10 },
      }
    default:
      return { component: InputText, props: { placeholder: props.field.type } }
  }
})
</script>
