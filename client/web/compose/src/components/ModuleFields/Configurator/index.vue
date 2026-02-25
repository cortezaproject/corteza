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
      <TabList class="shrink-0 z-10">
        <Tab value="basic">{{ $t('general.label.general') }}</Tab>
        <Tab v-if="hasKindSettings" value="kind">
          {{ $t(`general.fieldKinds.${mockField.kind}.label`) }}
        </Tab>
      </TabList>

      <TabPanels class="flex-1 overflow-y-auto">
        <!-- General Settings -->
        <TabPanel value="basic" class="px-0 py-4">
          <CConfiguratorBasic :field="mockField" />
        </TabPanel>

        <!-- Field-Specific Settings -->
        <TabPanel v-if="hasKindSettings" value="kind" class="px-0 py-4">
          <component :is="kindComponent" :field="mockField" />
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
import { compose } from '@cortezaproject/corteza-js-next'
import { computed, defineAsyncComponent, ref, watch } from 'vue'
import CConfiguratorBasic from './CConfiguratorBasic.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  field: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:visible', 'save'])

const activeTab = ref('basic')
const mockField = ref(null)

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

const header = computed(() => {
  if (!mockField.value) return ''
  return mockField.value.label || mockField.value.name || mockField.value.kind
})

// Dynamically load the specific configurator for this field kind
const kindComponent = computed(() => {
  if (!mockField.value?.kind) return null
  return defineAsyncComponent(() => import(`./kinds/${mockField.value.kind}.vue`).catch(() => null))
})

// Hardcode which types have specific configurator settings for now
const hasKindSettings = computed(() => {
  if (!mockField.value?.kind) return false
  const kindsWithSettings = ['String', 'Number', 'DateTime', 'Bool', 'Select']
  return kindsWithSettings.includes(mockField.value.kind)
})

function handleSave() {
  emit('save', mockField.value)
  emit('update:visible', false)
}

function handleCancel() {
  emit('update:visible', false)
}
</script>
