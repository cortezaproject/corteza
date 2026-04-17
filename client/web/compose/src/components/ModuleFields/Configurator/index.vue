<template>
  <Dialog
    :visible="visible"
    modal
    position="top"
    :header="header"
    :style="{ width: '50vw' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :pt="{
      content: { class: 'p-0 flex flex-col !overflow-hidden' },
      footer: { class: 'border-t border-surface p-3' },
    }"
    @update:visible="emit('update:visible', $event)"
  >
    <Tabs v-if="mockField" v-model:value="activeTab" class="flex flex-col flex-1 min-h-0">
      <TabList class="shrink-0 z-10 w-full overflow-x-auto whitespace-nowrap">
        <Tab value="basic">{{ $t('general.label.general') }}</Tab>
        <Tab v-if="hasKindSettings" value="kind">
          {{ $t(`general.fieldKinds.${mockField.kind}.label`) }}
        </Tab>
        <Tab v-if="mockField.cap?.multi" :disabled="!mockField.isMulti" value="multi">
          {{ $t('field.label.multi') }}
        </Tab>
        <Tab value="validation">{{ $t('field.validators.label') }}</Tab>
      </TabList>

      <TabPanels class="flex-1 overflow-y-auto w-full">
        <!-- General Settings -->
        <TabPanel value="basic" class="px-0">
          <CConfiguratorBasic :namespace="namespace" />
        </TabPanel>

        <!-- Field-Specific Settings -->
        <TabPanel v-if="hasKindSettings" value="kind" class="px-0">
          <component :is="kindComponent" />
        </TabPanel>

        <!-- Multi-Value -->
        <TabPanel v-if="mockField.cap?.multi" value="multi" class="px-0">
          <CConfiguratorMultiDelimiter />
        </TabPanel>

        <!-- Validation -->
        <TabPanel value="validation" class="px-0">
          <CConfiguratorValidation />
        </TabPanel>
      </TabPanels>
    </Tabs>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          size="small"
          outlined
          @click="handleCancel"
        />
        <Button :label="$t('general.label.save')" size="small" @click="handleSave" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { compose } from '@planetcrust/human-js'
import { computed, provide, ref, shallowRef, watch } from 'vue'
import CConfiguratorBasic from './CConfiguratorBasic.vue'
import CConfiguratorValidation from './CConfiguratorValidation.vue'
import CConfiguratorMultiDelimiter from './CConfiguratorMultiDelimiter.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  field: {
    type: Object,
    default: null,
  },
  namespace: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:visible', 'save'])

const activeTab = ref('basic')
const mockField = ref(null)
const kindComponent = shallowRef(null)

provide('fieldDraft', mockField)

watch(
  () => props.visible,
  val => {
    if (val && props.field) {
      activeTab.value = 'basic'
      // Deep clone field data and create an actual module field class instance
      const raw = JSON.parse(JSON.stringify(props.field))
      mockField.value = compose.ModuleFieldMaker(raw)
    } else {
      mockField.value = null
    }
  },
  { immediate: true },
)

// Try to import the kind-specific configurator; set to null if none exists
watch(
  () => mockField.value?.kind,
  async kind => {
    if (!kind) {
      kindComponent.value = null
      return
    }
    try {
      const mod = await import(`./kinds/${kind}.vue`)
      kindComponent.value = mod.default
    } catch {
      kindComponent.value = null
    }
  },
  { immediate: true },
)

const hasKindSettings = computed(() => !!kindComponent.value)

const header = computed(() => {
  if (!mockField.value) return ''
  return mockField.value.label || mockField.value.name || mockField.value.kind
})

function handleSave() {
  emit('save', mockField.value)
  emit('update:visible', false)
}

function handleCancel() {
  emit('update:visible', false)
}
</script>
