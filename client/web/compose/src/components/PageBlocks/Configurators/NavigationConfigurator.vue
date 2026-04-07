<template>
  <div class="flex flex-col gap-3">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.navigation.appearance') }}</label>
        <Select
          v-model="appearance"
          :options="appearanceOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.navigation.alignment') }}</label>
        <Select
          v-model="alignment"
          :options="alignmentOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-primary font-medium text-sm">{{ $t('block.navigation.justify') }}</label>
        <Select
          v-model="justify"
          :options="justifyOptions"
          option-label="label"
          option-value="value"
          class="w-full"
        />
      </div>
    </div>

    <Divider />

    <div class="flex items-center justify-between">
      <label class="text-primary font-medium text-sm">{{ $t('block.navigation.navigationItems') }}</label>
      <Button
        :label="$t('block.navigation.add')"
        icon="pi pi-plus"
        size="small"
        severity="secondary"
        @click="addNavItem"
      />
    </div>

    <div v-for="(item, index) in navItems" :key="index" class="border border-surface rounded-border p-3">
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <div class="flex flex-col gap-1">
            <Button
              icon="pi pi-chevron-up"
              text rounded size="small" severity="secondary"
              :disabled="index === 0"
              @click="moveItem(index, -1)"
            />
            <Button
              icon="pi pi-chevron-down"
              text rounded size="small" severity="secondary"
              :disabled="index === navItems.length - 1"
              @click="moveItem(index, 1)"
            />
          </div>

          <span class="text-sm font-medium flex-1">#{{ index + 1 }}</span>
          <Button
            icon="pi pi-trash"
            text rounded size="small" severity="danger"
            @click="removeNavItem(index)"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('block.navigation.type') }}</label>
          <Select
            :model-value="item.type"
            :options="typeOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="updateNavItem(index, 'type', $event)"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('block.navigation.fieldLabel') }}</label>
          <InputText
            :model-value="item.options?.item?.label || ''"
            class="w-full"
            @update:model-value="updateNavItemLabel(index, $event)"
          />
        </div>

        <div v-if="item.type === 'url'" class="flex flex-col gap-1">
          <label class="text-sm text-muted-color">{{ $t('block.navigation.url') }}</label>
          <InputText
            :model-value="item.options?.item?.url || ''"
            placeholder="https://"
            class="w-full"
            @update:model-value="updateNavItemUrl(index, $event)"
          />
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.navigation.textColor') }}</label>
            <CInputColorPicker
              :model-value="item.options?.item?.textColor || ''"
              show-text
              @update:model-value="updateNavItemStyle(index, 'textColor', $event)"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-sm text-muted-color">{{ $t('block.navigation.backgroundColor') }}</label>
            <CInputColorPicker
              :model-value="item.options?.item?.backgroundColor || ''"
              show-text
              @update:model-value="updateNavItemStyle(index, 'backgroundColor', $event)"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'

const { CInputColorPicker } = components

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const typeOptions = [
  { value: 'url', label: t('block.navigation.url') },
  { value: 'compose', label: t('block.navigation.composePage') },
  { value: 'dropdown', label: t('block.navigation.dropdown') },
  { value: 'text-section', label: t('block.navigation.textSection') },
]

const appearanceOptions = [
  { value: 'tabs', label: t('block.navigation.tabs') },
  { value: 'pills', label: t('block.navigation.pills') },
  { value: 'small', label: t('block.navigation.small') },
]

const alignmentOptions = [
  { value: 'left', label: t('block.navigation.left') },
  { value: 'center', label: t('block.navigation.center') },
  { value: 'right', label: t('block.navigation.right') },
]

const justifyOptions = [
  { value: 'justify', label: t('block.navigation.justify') },
  { value: 'none', label: t('block.navigation.none') },
]

const navItems = computed(() => props.block.options?.navigationItems || [])

const appearance = computed({
  get: () => props.block.options?.display?.appearance || 'tabs',
  set: v => updateOptions('display', { ...props.block.options?.display, appearance: v }),
})

const alignment = computed({
  get: () => props.block.options?.display?.alignment || 'left',
  set: v => updateOptions('display', { ...props.block.options?.display, alignment: v }),
})

const justify = computed({
  get: () => props.block.options?.display?.justify || 'none',
  set: v => updateOptions('display', { ...props.block.options?.display, justify: v }),
})

function updateOptions(key, value) {
  emit('update:block', {
    ...props.block,
    options: { ...props.block.options, [key]: value },
  })
}

function addNavItem() {
  const items = [...navItems.value, {
    type: 'url',
    options: { enabled: true, item: { label: '', url: '', textColor: '', backgroundColor: '' } },
  }]
  updateOptions('navigationItems', items)
}

function removeNavItem(index) {
  const items = [...navItems.value]
  items.splice(index, 1)
  updateOptions('navigationItems', items)
}

function updateNavItem(index, key, value) {
  const items = [...navItems.value]
  items[index] = { ...items[index], [key]: value }
  updateOptions('navigationItems', items)
}

function updateNavItemLabel(index, label) {
  const items = [...navItems.value]
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: { ...items[index].options?.item, label },
    },
  }
  updateOptions('navigationItems', items)
}

function updateNavItemUrl(index, url) {
  const items = [...navItems.value]
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: { ...items[index].options?.item, url },
    },
  }
  updateOptions('navigationItems', items)
}

function updateNavItemStyle(index, styleKey, value) {
  const items = [...navItems.value]
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: { ...items[index].options?.item, [styleKey]: value },
    },
  }
  updateOptions('navigationItems', items)
}

function moveItem(index, direction) {
  const items = [...navItems.value]
  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= items.length) return
  const temp = items[index]
  items[index] = items[newIndex]
  items[newIndex] = temp
  updateOptions('navigationItems', items)
}
</script>
