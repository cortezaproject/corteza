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

    <Panel
      v-for="(item, index) in navItems"
      :key="index"
      :header="item.options?.item?.label || $t('block.navigation.newItem', 'New item')"
      toggleable
    >
      <template #icons>
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
        <Button
          icon="pi pi-trash"
          text rounded size="small" severity="danger"
          @click="removeNavItem(index)"
        />
      </template>

      <div class="flex flex-col gap-2">
        <!-- Label -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.navigation.fieldLabel') }}</label>
          <InputText
            :model-value="item.options?.item?.label || ''"
            class="w-full"
            @update:model-value="updateNavItemOption(index, 'label', $event)"
          />
        </div>

        <!-- Type selector -->
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">{{ $t('block.navigation.type') }}</label>
          <Select
            :model-value="item.type"
            :options="typeOptions"
            option-label="label"
            option-value="value"
            class="w-full"
            @update:model-value="updateNavItem(index, 'type', $event)"
          />
        </div>

        <!-- URL type: URL input + target -->
        <template v-if="item.type === 'url'">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.url') }}</label>
            <InputText
              :model-value="item.options?.item?.url || ''"
              placeholder="https://"
              class="w-full"
              @update:model-value="updateNavItemOption(index, 'url', $event)"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.target', 'Open in') }}</label>
            <Select
              :model-value="item.options?.item?.target || 'sameTab'"
              :options="targetOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              @update:model-value="updateNavItemOption(index, 'target', $event)"
            />
          </div>
        </template>

        <!-- Compose Page type: page selector + options -->
        <template v-else-if="item.type === 'compose'">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.composePage', 'Page') }}</label>
            <Select
              :model-value="item.options?.item?.pageID || null"
              :options="pageOptions"
              option-label="label"
              option-value="value"
              :placeholder="$t('block.navigation.selectPage', 'Select page')"
              filter
              show-clear
              class="w-full"
              @update:model-value="updateNavItemOption(index, 'pageID', $event)"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.target', 'Open in') }}</label>
            <Select
              :model-value="item.options?.item?.target || 'sameTab'"
              :options="targetOptions"
              option-label="label"
              option-value="value"
              class="w-full"
              @update:model-value="updateNavItemOption(index, 'target', $event)"
            />
          </div>
          <CInputSwitch
            :model-value="!!item.options?.item?.displaySubPages"
            :label="$t('block.navigation.displaySubPages', 'Show as dropdown with sub-pages')"
            @update:model-value="updateNavItemOption(index, 'displaySubPages', $event)"
          />
        </template>

        <!-- Dropdown type: dropdown label + items -->
        <template v-else-if="item.type === 'dropdown'">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.dropdownLabel', 'Dropdown button label') }}</label>
            <InputText
              :model-value="item.options?.item?.dropdown?.label || ''"
              class="w-full"
              @update:model-value="updateDropdownLabel(index, $event)"
            />
          </div>

          <div class="flex flex-col gap-2">
            <div class="flex items-center justify-between">
              <label class="text-primary font-medium text-sm">{{ $t('block.navigation.dropdownItems', 'Items') }}</label>
              <Button
                :label="$t('general.label.add')"
                icon="pi pi-plus"
                size="small"
                severity="secondary"
                text
                @click="addDropdownItem(index)"
              />
            </div>

            <div
              v-for="(dItem, dIndex) in (item.options?.item?.dropdown?.items || [])"
              :key="dIndex"
              class="border border-surface rounded-border p-2 flex flex-col gap-2"
            >
              <div class="flex items-center justify-between">
                <CInputSwitch
                  :model-value="!!dItem.delimiter"
                  :label="$t('block.navigation.delimiter', 'Separator')"
                  @update:model-value="updateDropdownItem(index, dIndex, 'delimiter', $event)"
                />
                <Button
                  icon="pi pi-trash"
                  text rounded size="small" severity="danger"
                  @click="removeDropdownItem(index, dIndex)"
                />
              </div>

              <template v-if="!dItem.delimiter">
                <InputText
                  :model-value="dItem.label || ''"
                  :placeholder="$t('block.navigation.fieldLabel')"
                  class="w-full"
                  size="small"
                  @update:model-value="updateDropdownItem(index, dIndex, 'label', $event)"
                />
                <InputText
                  :model-value="dItem.url || ''"
                  placeholder="https://"
                  class="w-full"
                  size="small"
                  @update:model-value="updateDropdownItem(index, dIndex, 'url', $event)"
                />
                <Select
                  :model-value="dItem.target || 'sameTab'"
                  :options="targetOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                  size="small"
                  @update:model-value="updateDropdownItem(index, dIndex, 'target', $event)"
                />
              </template>
            </div>
          </div>
        </template>

        <!-- Text and Background Colours (all types) -->
        <div class="grid grid-cols-2 gap-2">
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.textColor') }}</label>
            <CInputColorPicker
              :model-value="item.options?.item?.textColor || ''"
              show-text
              @update:model-value="updateNavItemOption(index, 'textColor', $event)"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-primary font-medium text-sm">{{ $t('block.navigation.backgroundColor') }}</label>
            <CInputColorPicker
              :model-value="item.options?.item?.backgroundColor || ''"
              show-text
              @update:model-value="updateNavItemOption(index, 'backgroundColor', $event)"
            />
          </div>
        </div>
      </div>
    </Panel>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'
import { usePageStore } from '@/stores/page'

const { CInputColorPicker, CInputSwitch } = components

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const emit = defineEmits(['update:block'])

const pageStore = usePageStore()

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

const targetOptions = [
  { value: 'sameTab', label: t('block.navigation.sameTab', 'Same tab') },
  { value: 'newTab', label: t('block.navigation.newTab', 'New tab') },
]

const pageOptions = computed(() => {
  return (pageStore.set || []).map(p => ({
    value: p.pageID,
    label: p.title || p.handle || p.pageID,
  }))
})

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

// Update a field inside item.options.item
function updateNavItemOption(index, key, value) {
  const items = [...navItems.value]
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: { ...items[index].options?.item, [key]: value },
    },
  }
  updateOptions('navigationItems', items)
}

function updateDropdownLabel(index, label) {
  const items = [...navItems.value]
  const dropdown = items[index].options?.item?.dropdown || {}
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: {
        ...items[index].options?.item,
        dropdown: { ...dropdown, label },
      },
    },
  }
  updateOptions('navigationItems', items)
}

function addDropdownItem(index) {
  const items = [...navItems.value]
  const dropdown = items[index].options?.item?.dropdown || { label: '', items: [] }
  const dItems = [...(dropdown.items || []), { label: '', url: '', target: 'sameTab', delimiter: false }]
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: {
        ...items[index].options?.item,
        dropdown: { ...dropdown, items: dItems },
      },
    },
  }
  updateOptions('navigationItems', items)
}

function removeDropdownItem(index, dIndex) {
  const items = [...navItems.value]
  const dropdown = items[index].options?.item?.dropdown || {}
  const dItems = [...(dropdown.items || [])]
  dItems.splice(dIndex, 1)
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: {
        ...items[index].options?.item,
        dropdown: { ...dropdown, items: dItems },
      },
    },
  }
  updateOptions('navigationItems', items)
}

function updateDropdownItem(index, dIndex, key, value) {
  const items = [...navItems.value]
  const dropdown = items[index].options?.item?.dropdown || {}
  const dItems = [...(dropdown.items || [])]
  dItems[dIndex] = { ...dItems[dIndex], [key]: value }
  items[index] = {
    ...items[index],
    options: {
      ...items[index].options,
      item: {
        ...items[index].options?.item,
        dropdown: { ...dropdown, items: dItems },
      },
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
