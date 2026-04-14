<template>
  <div class="trigger-schema-editor">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-sm font-medium text-color">Agent Input Schema</h3>
      <button type="button" class="p-button p-component p-button-sm p-button-text" @click="addProperty">
        <span class="p-button-icon p-button-icon-left pi pi-plus" data-pc-section="icon"></span>
        <span class="p-button-label" data-pc-section="label">Add Property</span>
      </button>
    </div>

    <!-- Empty state -->
    <div v-if="!properties.length" class="text-sm text-color decoration-clone mb-4">
      No schema properties defined for this agent trigger. Add properties to allow the agent to pass parameters dynamically.
    </div>

    <!-- Schema properties list -->
    <div v-else class="flex flex-col gap-3 mb-4">
      <div v-for="(prop, index) in properties" :key="'prop-' + index" class="p-3 border border-surface rounded-border flex flex-col gap-2">
        <div class="flex justify-between items-center w-full">
          <input type="text" v-model="prop.name" placeholder="Property name (e.g. recordID)" class="p-inputtext p-component w-full text-sm mr-2" />
          <button type="button" class="p-button p-component p-button-danger p-button-text p-button-sm ml-auto" @click="removeProperty(index)">
            <span class="p-button-icon pi pi-trash" data-pc-section="icon"></span>
          </button>
        </div>

        <div class="w-full flex gap-2">
           <input type="text" v-model="prop.description" placeholder="Description for the LLM..." class="p-inputtext p-component w-full text-sm" />
        </div>

        <div class="flex items-center gap-4 mt-1">
          <div class="flex items-center gap-2">
             <label class="text-xs text-color">Type</label>
             <select v-model="prop.type" class="p-inputtext p-component p-inputtext-sm w-24">
               <option value="String">String</option>
               <option value="Boolean">Boolean</option>
               <option value="Number">Number</option>
               <option value="ID">ID</option>
             </select>
          </div>
          <div class="flex items-center gap-2">
            <input type="checkbox" v-model="prop.required" :id="'req-' + index" class="cursor-pointer" />
            <label :for="'req-' + index" class="text-xs text-color cursor-pointer">Required</label>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-2 mt-4">
      <button type="button" class="p-button p-component w-full" :class="{'p-disabled': isSaving}" :disabled="isSaving" @click="saveSchema">
        <span v-if="isSaving" class="p-button-icon p-button-icon-left pi pi-spinner pi-spin"></span>
        <span v-else class="p-button-icon p-button-icon-left pi pi-sync"></span>
        <span class="p-button-label">Sync Schema</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, getCurrentInstance } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'

const props = defineProps({
  triggerDef: {
    type: Object,
    required: true,
  },
})

const { t } = useI18n()
const toast = useToast()

const { proxy } = getCurrentInstance()

const definition = ref(null)
const properties = ref([])
const isSaving = ref(false)

// Extract handle from the resource type (e.g., automation:trigger-definition:agent_12345)
const getHandle = () => {
  const parts = props.triggerDef.resourceType.split(':')
  return parts[parts.length - 1]
}

const loadDefinition = async () => {
  const handle = getHandle()
  if (!handle) return

  try {
    const rawResult = await proxy.$AutomationAPI.api()
      .get('/ng-automation/trigger-definitions/', { params: { query: handle } })
      .then(r => r.data)

    // stdResolve equivalent or raw Axios unwrap
    const payload = rawResult.response || rawResult
    const match = payload.set ? payload.set.find(d => d.resourceType === props.triggerDef.resourceType) : null

    if (match) {
      definition.value = match
      properties.value = (match.properties || []).map(p => ({
        name: p.name || '',
        type: p.type || 'String',
        description: p.meta?.short || '',
        required: !!p.required
      }))
    } else {
      // Stub
      definition.value = { handle, skipEventBus: true, properties: [] }
      properties.value = []
    }
  } catch (err) {
    console.error('Failed to load trigger definition', err)
  }
}

onMounted(() => loadDefinition())

watch(() => props.triggerDef, () => {
  loadDefinition()
}, { deep: true })

const addProperty = () => {
  properties.value.push({
    name: '',
    type: 'String',
    description: '',
    required: false
  })
}

const removeProperty = (index) => {
  properties.value.splice(index, 1)
}

const saveSchema = async () => {
  isSaving.value = true
  
  // Clean empty properties
  const schema = properties.value.filter(p => p.name.trim() !== '').map(p => ({
    name: p.name,
    type: p.type,
    required: p.required,
    meta: { short: p.description },
  }))

  const payload = {
    ...definition.value,
    properties: schema
  }

  try {
    if (payload.id || payload.triggerDefinitionID) {
      payload.triggerDefinitionID = payload.triggerDefinitionID || payload.id
      await proxy.$AutomationAPI.api().put(`/ng-automation/trigger-definitions/${payload.triggerDefinitionID}`, payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Agent schema updated', life: 3000 })
    } else {
      await proxy.$AutomationAPI.api().post('/ng-automation/trigger-definitions/', payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Agent schema established', life: 3000 })
    }

    // Refresh cleanly
    await loadDefinition()
  } catch (err) {
    console.error('Failed to save schema', err)
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to update schema', life: 5000 })
  } finally {
    isSaving.value = false
  }
}
</script>

<style scoped>
.trigger-schema-editor {
  padding-top: 1rem;
}
</style>
