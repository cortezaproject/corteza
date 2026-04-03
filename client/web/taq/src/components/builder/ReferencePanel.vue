<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between px-3 py-2 border-b border-surface">
      <h4 class="text-sm font-semibold text-color">{{ $t('builder.referencePanel.title') }}</h4>
      <Button icon="pi pi-times" text rounded size="small" @click="emit('close')" />
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-auto">
      <div v-if="filteredResults.length === 0" class="text-sm text-muted-color text-center py-8">
        {{ $t('builder.referencePanel.noResults') }}
      </div>

      <Accordion v-else v-model:value="expandedPanels" multiple class="flex flex-col gap-3 p-3">
        <AccordionPanel
          v-for="step in filteredResults"
          :key="step.handle"
          :value="step.handle"
          class="border border-surface rounded-border overflow-hidden bg-surface shadow-sm"
        >
          <AccordionHeader>
            <div class="flex items-center gap-2 p-1">
              <i v-if="step.icon?.iconClass" :class="step.icon.iconClass" class="text-sm" />
              <span class="text-sm font-medium">{{ step.label }}</span>
            </div>
          </AccordionHeader>
          <AccordionContent>
            <div class="flex flex-col gap-1">
              <template v-for="result in step.properties || step.results" :key="result.sourceName">
                <!-- Expandable result (e.g. ComposeRecord) — nested accordion -->
                <div v-if="result.expandable" class="mt-1">
                  <Accordion
                    :value="expandedSubPanels[step.handle + ':' + result.sourceName] || []"
                    @update:value="
                      v => (expandedSubPanels[step.handle + ':' + result.sourceName] = v)
                    "
                    multiple
                  >
                    <AccordionPanel
                      :value="result.sourceName"
                      class="border border-surface rounded-border overflow-hidden"
                    >
                      <AccordionHeader>
                        <div class="flex items-center justify-between w-full gap-2">
                          <div class="flex items-center gap-2">
                            <i class="pi pi-database text-xs text-[--p-text-muted-color]" />
                            <span class="text-sm text-color">{{ result.name }}</span>
                          </div>
                        </div>
                      </AccordionHeader>
                      <AccordionContent>
                        <!-- Top-level: select the whole record -->
                        <button
                          class="flex items-center justify-between px-2 py-1.5 rounded-md hover:bg-emphasis cursor-pointer text-left transition-colors w-full mb-1"
                          :class="
                            isActive(step.handle, result.sourceName)
                              ? 'bg-highlight !text-primary'
                              : ''
                          "
                          @click="emit('select', { scope: step.handle, source: result.sourceName })"
                        >
                          <span
                            class="text-sm italic"
                            :class="isActive(step.handle, result.sourceName) ? '' : 'text-color'"
                          >
                            {{ result.name }} (whole)
                          </span>
                        </button>

                        <!-- Loading state -->
                        <div
                          v-if="loadingFields[fieldKey(step, result)]"
                          class="flex items-center gap-2 px-3 py-4 justify-center"
                        >
                          <i class="pi pi-spin pi-spinner text-[--p-text-muted-color]" />
                          <span class="text-xs text-muted">Loading fields...</span>
                        </div>

                        <!-- Sub-fields -->
                        <template v-else-if="recordFields[fieldKey(step, result)]?.length">
                          <button
                            v-for="field in recordFields[fieldKey(step, result)]"
                            :key="field.name"
                            class="flex items-center justify-between px-2 py-1.5 rounded-md hover:bg-emphasis cursor-pointer text-left transition-colors w-full"
                            :class="
                              isActive(
                                step.handle,
                                field.isProperty
                                  ? result.sourceName + '.' + field.name
                                  : result.sourceName + '.values.' + field.name,
                              )
                                ? 'bg-highlight !text-primary'
                                : ''
                            "
                            @click="
                              emit('select', {
                                scope: step.handle,
                                source: field.isProperty
                                  ? result.sourceName + '.' + field.name
                                  : result.sourceName + '.values.' + field.name,
                              })
                            "
                          >
                            <span
                              class="flex items-center justify-between w-full text-sm italic"
                              :class="
                                isActive(
                                  step.handle,
                                  field.isProperty
                                    ? result.sourceName + '.' + field.name
                                    : result.sourceName + '.values.' + field.name,
                                )
                                  ? ''
                                  : 'text-color'
                              "
                            >
                              <span>{{ field.label || field.name }}</span>
                            </span>
                          </button>
                        </template>

                        <!-- No module configured for custom fields -->
                        <div
                          v-else-if="result.types?.includes('ComposeRecord') && !result.moduleID"
                          class="text-xs text-[--p-text-muted-color] px-2 py-1.5"
                        >
                          No module configured
                        </div>

                        <!-- No fields found -->
                        <div v-else class="text-xs text-[--p-text-muted-color] px-1.5 py-1">
                          No fields found
                        </div>
                      </AccordionContent>
                    </AccordionPanel>
                  </Accordion>
                </div>

                <!-- Simple result — clickable button -->
                <button
                  v-else
                  class="flex items-center justify-between px-2 py-1.5 rounded-md hover:bg-emphasis cursor-pointer text-left transition-colors"
                  :class="
                    isActive(step.handle, result.sourceName) ? 'bg-highlight !text-primary' : ''
                  "
                  @click="emit('select', { scope: step.handle, source: result.sourceName })"
                >
                  <span
                    class="text-sm"
                    :class="isActive(step.handle, result.sourceName) ? '' : 'text-color'"
                  >
                    {{ result.name }}
                  </span>
                </button>
              </template>
            </div>
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, watch, ref } from 'vue'
import { useComposeResourceStore } from '@cortezaproject/corteza-vue-next'

const STRUCT_FIELDS = {
  ComposeRecord: [
    { name: 'recordID', label: 'Record ID', kind: 'ID', isProperty: true },
    { name: 'moduleID', label: 'Module ID', kind: 'ID', isProperty: true },
    { name: 'namespaceID', label: 'Namespace ID', kind: 'ID', isProperty: true },
    { name: 'values', label: 'Values', kind: 'Object', isProperty: true },
    { name: 'ownedBy', label: 'Owned By', kind: 'ID', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
  ],
  SystemUser: [
    { name: 'userID', label: 'User ID', kind: 'ID', isProperty: true },
    { name: 'email', label: 'Email', kind: 'String', isProperty: true },
    { name: 'name', label: 'Name', kind: 'String', isProperty: true },
    { name: 'username', label: 'Username', kind: 'String', isProperty: true },
    { name: 'handle', label: 'Handle', kind: 'String', isProperty: true },
    { name: 'emailConfirmed', label: 'Email Confirmed', kind: 'Boolean', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
    { name: 'suspendedAt', label: 'Suspended At', kind: 'DateTime', isProperty: true },
  ],
  SystemRole: [
    { name: 'roleID', label: 'Role ID', kind: 'ID', isProperty: true },
    { name: 'name', label: 'Name', kind: 'String', isProperty: true },
    { name: 'handle', label: 'Handle', kind: 'String', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
    { name: 'archivedAt', label: 'Archived At', kind: 'DateTime', isProperty: true },
  ],
  SystemApplication: [
    { name: 'applicationID', label: 'Application ID', kind: 'ID', isProperty: true },
    { name: 'name', label: 'Name', kind: 'String', isProperty: true },
    { name: 'enabled', label: 'Enabled', kind: 'Boolean', isProperty: true },
    { name: 'unlisted', label: 'Unlisted', kind: 'Boolean', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
  ],
  ComposeModule: [
    { name: 'moduleID', label: 'Module ID', kind: 'ID', isProperty: true },
    { name: 'namespaceID', label: 'Namespace ID', kind: 'ID', isProperty: true },
    { name: 'name', label: 'Name', kind: 'String', isProperty: true },
    { name: 'handle', label: 'Handle', kind: 'String', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
  ],
  ComposeNamespace: [
    { name: 'namespaceID', label: 'Namespace ID', kind: 'ID', isProperty: true },
    { name: 'name', label: 'Name', kind: 'String', isProperty: true },
    { name: 'slug', label: 'Slug', kind: 'String', isProperty: true },
    { name: 'enabled', label: 'Enabled', kind: 'Boolean', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
  ],
  ComposePage: [
    { name: 'pageID', label: 'Page ID', kind: 'ID', isProperty: true },
    { name: 'moduleID', label: 'Module ID', kind: 'ID', isProperty: true },
    { name: 'namespaceID', label: 'Namespace ID', kind: 'ID', isProperty: true },
    { name: 'title', label: 'Title', kind: 'String', isProperty: true },
    { name: 'handle', label: 'Handle', kind: 'String', isProperty: true },
    { name: 'description', label: 'Description', kind: 'String', isProperty: true },
    { name: 'createdAt', label: 'Created At', kind: 'DateTime', isProperty: true },
    { name: 'updatedAt', label: 'Updated At', kind: 'DateTime', isProperty: true },
    { name: 'deletedAt', label: 'Deleted At', kind: 'DateTime', isProperty: true },
  ],
  HttpRequest: [
    { name: 'method', label: 'Method', kind: 'String', isProperty: true },
    { name: 'url', label: 'URL', kind: 'String', isProperty: true },
    { name: 'header', label: 'Headers', kind: 'Object', isProperty: true },
    { name: 'cookie', label: 'Cookies', kind: 'Object', isProperty: true },
    { name: 'form', label: 'Form', kind: 'Object', isProperty: true },
    { name: 'query', label: 'Query', kind: 'Object', isProperty: true },
    { name: 'body', label: 'Body', kind: 'Reader', isProperty: true },
    { name: 'postForm', label: 'Post Form', kind: 'Object', isProperty: true },
    { name: 'remoteAddr', label: 'Remote Address', kind: 'String', isProperty: true },
  ],
}

const props = defineProps({
  upstreamResults: {
    type: Array,
    default: () => [],
  },
  activeArgument: {
    type: Object,
    default: null,
  },
  currentReference: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['select', 'close'])

const store = useComposeResourceStore()

// Check if a reference item matches the currently active reference
function isActive(scope, source) {
  if (!props.currentReference) return false
  return props.currentReference.scope === scope && props.currentReference.source === source
}

// Track expanded sub-panels
const expandedSubPanels = reactive({})

// Expand all top-level panels by default but allow collapsing
const expandedPanels = ref([])

// Cache fetched module fields per step+result key
const recordFields = reactive({})
const loadingFields = reactive({})

function fieldKey(step, result) {
  return `${step.handle}:${result.sourceName}`
}

// Type compatibility check
function typesOverlap(acceptedTypes, resultTypes) {
  if (!acceptedTypes?.length || acceptedTypes.includes('Any')) return true
  if (!resultTypes?.length) return true
  return acceptedTypes.some(t => resultTypes.includes(t))
}

// Filtered results based on active argument's accepted types
const filteredResults = computed(() => {
  const accepted = props.activeArgument?.types

  // No active argument (browse mode) → show everything
  if (!accepted?.length) return props.upstreamResults

  const isAny = accepted.includes('Any')
  if (isAny) return props.upstreamResults

  return props.upstreamResults
    .map(step => {
      const items = step.properties || step.results || []
      const filtered = items.filter(result => {
        if (result.expandable) {
          // Show expandable if param accepts the parent type (e.g. ComposeRecord)
          if (typesOverlap(accepted, result.types)) return true

          // Or if any sub-field kind matches the accepted types
          const fields = recordFields[fieldKey(step, result)] || []
          return fields.some(f => accepted.includes(f.kind))
        }

        // Flat result: check type overlap
        return typesOverlap(accepted, result.types)
      })

      if (filtered.length === 0) return null

      const newStep = { ...step }
      if (step.properties) {
        newStep.properties = filtered
      } else {
        newStep.results = filtered
      }
      return newStep
    })
    .filter(Boolean)
})

watch(
  filteredResults,
  newResults => {
    const current = new Set(expandedPanels.value)
    const newHandles = newResults.map(s => s.handle).filter(h => !current.has(h))
    if (newHandles.length > 0) {
      expandedPanels.value = [...expandedPanels.value, ...newHandles]
    }
  },
  { immediate: true },
)

// Fetch module fields for expandable results
async function fetchFields(step, result) {
  const key = fieldKey(step, result)
  if (recordFields[key] || loadingFields[key]) return

  let resolvedFields = []

  // Add standard struct fields based on result type
  for (const t of result.types || []) {
    if (STRUCT_FIELDS[t]) {
      resolvedFields = [...resolvedFields, ...STRUCT_FIELDS[t]]
    }
  }

  // If it's a ComposeRecord with module configuration, fetch custom fields
  if (result.types?.includes('ComposeRecord') && result.namespaceID && result.moduleID) {
    loadingFields[key] = true
    try {
      const mod = await store.resolveModule(result.namespaceID, result.moduleID)
      if (mod?.fields) {
        const moduleFields = mod.fields
          .filter(f => !f.isSystem)
          .map(f => ({
            name: f.name,
            label: f.label || f.name,
            kind: f.kind,
            isProperty: false,
          }))
        resolvedFields = [...resolvedFields, ...moduleFields]
      }
    } catch {
      // Gracefully continue with standard properties
    } finally {
      loadingFields[key] = false
    }
  }

  recordFields[key] = resolvedFields
}

// Auto-fetch fields for all expandable results when upstream data changes
watch(
  () => props.upstreamResults,
  results => {
    for (const step of results) {
      const items = step.properties || step.results || []
      for (const result of items) {
        if (result.expandable) {
          // Auto-expand sub-panels
          const subKey = step.handle + ':' + result.sourceName
          if (!expandedSubPanels[subKey]) {
            expandedSubPanels[subKey] = [result.sourceName]
          }
          fetchFields(step, result)
        }
      }
    }
  },
  { immediate: true },
)
</script>

<style scoped>
:deep(.p-accordionheader) {
  padding: 0.375rem 0.5rem;
  font-size: 0.8125rem;
}
:deep(.p-accordioncontent-content) {
  padding: 0.25rem 0.25rem 0.25rem 0.5rem;
}
:deep(.p-accordion) {
  /* Let tailwind flex handle gap */
}
:deep(.p-accordionpanel) {
  /* Allowed standard borders to apply through tailwind */
}
</style>
