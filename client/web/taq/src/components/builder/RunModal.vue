<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('builder.runModal.header', 'Execution Scope')"
    :style="{ width: '50rem' }"
    @update:visible="emit('update:visible', $event)"
  >
    <div class="mb-4 text-muted-color">
      {{
        $t(
          'builder.runModal.description',
          'Provide initial variables to execute triggers that expect a specific scope.',
        )
      }}
    </div>

    <div class="flex flex-col gap-4">
      <div v-for="prop in filteredProperties" :key="prop.name" class="flex flex-col gap-2">
        <label class="font-medium text-sm text-primary">
          {{ capitalize(prop.meta?.short || prop.name) }}
        </label>

        <CInputRecord
          v-if="prop.type === 'ComposeRecord'"
          v-model="scope[prop.name]"
          :namespaceID="prop.namespaceID || scope['namespace'] || undefined"
          :moduleID="prop.moduleID || scope['module'] || undefined"
          :disabled="!(prop.namespaceID || scope['namespace']) || !(prop.moduleID || scope['module'])"
          class="w-full"
        />
        <CInputModule
          v-else-if="prop.type === 'ComposeModule'"
          v-model="scope[prop.name]"
          :namespaceID="scope['namespace'] || undefined"
          :disabled="!scope['namespace']"
          :placeholder="
            !scope['namespace']
              ? $t('builder.configSidebar.selectNamespaceFirst', 'Select a namespace first')
              : undefined
          "
          class="w-full"
        />
        <CInputNamespace
          v-else-if="prop.type === 'ComposeNamespace'"
          v-model="scope[prop.name]"
          class="w-full"
        />
        <CInputUser
          v-else-if="prop.type === 'SystemUser'"
          v-model="scope[prop.name]"
          class="w-full"
        />
        <!-- Fallback to plain text for strings and other unsupported types right now -->
        <InputText v-else v-model="scope[prop.name]" class="w-full" />

        <span v-if="prop.meta?.description" class="text-xs text-muted-color">
          {{ prop.meta.description }}
        </span>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end w-full h-full items-center gap-2">
        <Button
          :label="$t('general.label.cancel', 'Cancel')"
          text
          size="small"
          severity="secondary"
          @click="close"
        />
        <Button
          :label="$t('builder.run', 'Run')"
          icon="pi pi-play"
          severity="success"
          size="small"
          @click="run"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import {
  CInputRecord,
  CInputModule,
  CInputNamespace,
  CInputUser,
} from '@cortezaproject/corteza-vue-next/src/components/input'

import { components } from '@cortezaproject/corteza-vue-next'
// Temporary fallback if InputText is not globally registered
const InputText = components.InputText || Object

const props = defineProps({
  visible: { type: Boolean, default: false },
  properties: { type: Array, default: () => [] },
})

const emit = defineEmits(['update:visible', 'run'])

const scope = ref({})

const supportedTypes = ['ComposeRecord', 'ComposeModule', 'ComposeNamespace', 'SystemUser']

import { computed } from 'vue'

const filteredProperties = computed(() => {
  return props.properties.filter(p => supportedTypes.includes(p.type))
})

function capitalize(str) {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// Reset scope when modal opens
watch(
  () => props.visible,
  isVisible => {
    if (isVisible) {
      scope.value = {}
      props.properties.forEach(p => {
        if (p.defaultValue) {
          scope.value[p.name] = p.defaultValue
        }
      })
    }
  },
)

function close() {
  emit('update:visible', false)
}

function run() {
  const formattedScope = {}
  Object.entries(scope.value).forEach(([key, val]) => {
    if (val !== undefined && val !== null && val !== '') {
      const propDef = props.properties.find(p => p.name === key)
      // If it is a string type or unknown, just pass as is
      if (!propDef || propDef.type === 'String' || propDef.type === 'Any') {
        formattedScope[key] = val
      } else {
        // Pack into typed value so the expr engine knows what it is
        let typedVal = val
        if (propDef.type === 'ComposeNamespace') {
          typedVal = { namespaceID: String(val) }
        } else if (propDef.type === 'ComposeModule') {
          typedVal = { 
            namespaceID: propDef.namespaceID || scope.value['namespace'] || undefined, 
            moduleID: String(val) 
          }
        } else if (propDef.type === 'ComposeRecord') {
          typedVal = { 
            namespaceID: propDef.namespaceID || scope.value['namespace'] || undefined, 
            moduleID: propDef.moduleID || scope.value['module'] || undefined, 
            recordID: String(val) 
          }
        } else if (propDef.type === 'SystemUser') {
          typedVal = { userID: String(val) }
        }
        
        formattedScope[key] = { '@type': propDef.type, '@value': typedVal }
      }
    }
  })

  emit('run', formattedScope)
  close()
}
</script>
