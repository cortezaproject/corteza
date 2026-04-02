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
          v-else-if="prop.type === 'SystemUser' || prop.type === 'User'"
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
          :loading="isFetchingContext"
          :disabled="isFetchingContext"
          @click="run"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, watch, computed, inject } from 'vue'
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

const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const scope = ref({})
const isFetchingContext = ref(false)

const supportedTypes = ['ComposeRecord', 'ComposeModule', 'ComposeNamespace', 'SystemUser', 'User']

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

async function run() {
  const formattedScope = {}
  isFetchingContext.value = true

  try {
    for (const [key, val] of Object.entries(scope.value)) {
      if (val !== undefined && val !== null && val !== '') {
        const propDef = props.properties.find(p => p.name === key)
        // If it is a string type or unknown, just pass as is
        if (!propDef || propDef.type === 'String' || propDef.type === 'Any') {
          formattedScope[key] = val
        } else {
          // Fetch full object so the expr engine receives the rich datatype as it expects
          let typedVal = val
          switch (propDef.type) {
            case 'ComposeNamespace':
              typedVal = await $ComposeAPI.namespaceRead({ namespaceID: String(val) })
              break
            case 'ComposeModule':
              typedVal = await $ComposeAPI.moduleRead({ 
                namespaceID: propDef.namespaceID || scope.value['namespace'] || undefined, 
                moduleID: String(val) 
              })
              break
            case 'ComposeRecord':
              typedVal = await $ComposeAPI.recordRead({ 
                namespaceID: propDef.namespaceID || scope.value['namespace'] || undefined, 
                moduleID: propDef.moduleID || scope.value['module'] || undefined, 
                recordID: String(val) 
              })
              break
            case 'User':
            case 'SystemUser':
              typedVal = await $SystemAPI.userRead({ userID: String(val) })
              break
          }
          
          formattedScope[key] = { '@type': propDef.type, '@value': typedVal }
        }
      }
    }
  } catch (err) {
    console.error('[RunModal] Error fetching context references:', err)
    $toast?.toastDanger('Failed to fetch full object scopes for execution. Ensure your selections are valid.')
    return
  } finally {
    isFetchingContext.value = false
  }

  emit('run', formattedScope)
  close()
}
</script>
