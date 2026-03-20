<template>
  <div class="flex flex-col gap-4">
    <Message severity="info" variant="simple" class="mb-0">
      {{ $t('block.tabs.alertTitle') }}
    </Message>

    <!-- Style section -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.tabs.style.label') }}</h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.tabs.style.appearance') }}</label>
          <Select
            v-model="appearance"
            :options="appearanceOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.tabs.style.alignment') }}</label>
          <Select
            v-model="alignment"
            :options="alignmentOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.tabs.style.justify') }}</label>
          <Select
            v-model="justify"
            :options="justifyOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.tabs.style.orientation') }}</label>
          <Select
            v-model="orientation"
            :options="orientationOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.tabs.style.position') }}</label>
          <Select
            v-model="position"
            :options="positionOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>
      </div>
    </div>

    <Divider />

    <!-- Tab list -->
    <div class="flex items-center justify-between">
      <label class="text-primary font-medium text-sm">{{ $t('block.tabs.title') }}</label>
      <Button
        :label="$t('general.label.add')"
        icon="pi pi-plus"
        size="small"
        severity="secondary"
        @click="addTab"
      />
    </div>

    <div v-for="(tab, index) in tabs" :key="index" class="border border-surface rounded-border p-3">
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <div class="flex flex-col gap-1">
            <Button
              icon="pi pi-chevron-up"
              text
              rounded
              size="small"
              severity="secondary"
              :disabled="index === 0"
              @click="moveTab(index, -1)"
            />
            <Button
              icon="pi pi-chevron-down"
              text
              rounded
              size="small"
              severity="secondary"
              :disabled="index === tabs.length - 1"
              @click="moveTab(index, 1)"
            />
          </div>

          <InputText
            :model-value="tab.title || ''"
            :placeholder="`${$t('block.tabs.tab')} ${index + 1}`"
            class="flex-1"
            @update:model-value="updateTab(index, 'title', $event)"
          />

          <Button
            icon="pi pi-trash"
            text
            rounded
            size="small"
            severity="danger"
            @click="removeTab(index)"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('block.tabs.blockAssociation') }}</label>
          <Select
            :model-value="tab.blockID || ''"
            :options="blockOptions"
            option-label="label"
            option-value="value"
            :placeholder="$t('block.tabs.selectBlock')"
            class="w-full"
            show-clear
            @update:model-value="updateTab(index, 'blockID', $event)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  blocks: { type: Array, default: () => [] },
})

const emit = defineEmits(['update:block'])

const appearanceOptions = [
  { value: 'tabs', label: t('block.tabs.style.appearanceTabs') },
  { value: 'pills', label: t('block.tabs.style.appearancePills') },
  { value: 'small', label: t('block.tabs.style.appearanceSmall') },
]

const alignmentOptions = [
  { value: 'left', label: t('block.tabs.style.alignmentLeft') },
  { value: 'center', label: t('block.tabs.style.alignmentCenter') },
  { value: 'right', label: t('block.tabs.style.alignmentRight') },
]

const justifyOptions = [
  { value: 'justify', label: t('block.tabs.style.justifyJustify') },
  { value: 'none', label: t('block.tabs.style.justifyNone') },
]

const orientationOptions = [
  { value: 'horizontal', label: t('block.tabs.style.orientationHorizontal') },
  { value: 'vertical', label: t('block.tabs.style.orientationVertical') },
]

const positionOptions = [
  { value: 'start', label: t('block.tabs.style.positionStart') },
  { value: 'end', label: t('block.tabs.style.positionEnd') },
]

const tabs = computed(() => props.block.options?.tabs || [])

const blockOptions = computed(() => {
  return (props.blocks || []).map(b => {
    const bid = (b.blockID && b.blockID !== '0') ? b.blockID : b.meta?.tempID || ''
    return {
      value: bid,
      label: b.title || b.kind || bid,
    }
  }).filter(o => o.value)
})

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

function updateStyle(key, value) {
  const style = { ...(props.block.options?.style || {}), [key]: value }
  updateOptions('style', style)
}

const appearance = computed({
  get: () => props.block.options?.style?.appearance || 'tabs',
  set: v => updateStyle('appearance', v),
})

const alignment = computed({
  get: () => props.block.options?.style?.alignment || 'left',
  set: v => updateStyle('alignment', v),
})

const justify = computed({
  get: () => props.block.options?.style?.justify || 'none',
  set: v => updateStyle('justify', v),
})

const orientation = computed({
  get: () => props.block.options?.style?.orientation || 'horizontal',
  set: v => updateStyle('orientation', v),
})

const position = computed({
  get: () => props.block.options?.style?.position || 'start',
  set: v => updateStyle('position', v),
})

function addTab() {
  updateOptions('tabs', [...tabs.value, { title: '', blockID: '' }])
}

function removeTab(index) {
  const items = [...tabs.value]
  items.splice(index, 1)
  updateOptions('tabs', items)
}

function updateTab(index, key, value) {
  const items = [...tabs.value]
  items[index] = { ...items[index], [key]: value }
  updateOptions('tabs', items)
}

function moveTab(index, direction) {
  const items = [...tabs.value]
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= items.length) return
  const temp = items[index]
  items[index] = items[newIndex]
  items[newIndex] = temp
  updateOptions('tabs', items)
}
</script>
