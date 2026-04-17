<template>
  <div class="flex flex-col gap-4">
    <!-- Style section -->
    <div class="flex flex-col gap-3">
      <h5 class="text-lg font-semibold text-primary m-0">{{ $t('block.tabs.style.label') }}</h5>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.tabs.style.appearance') }}
          </label>
          <Select
            v-model="appearance"
            :options="appearanceOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.tabs.style.justify') }}
          </label>
          <Select
            v-model="justify"
            :options="justifyOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.tabs.style.orientation') }}
          </label>
          <Select
            v-model="orientation"
            :options="orientationOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.tabs.style.position') }}
          </label>
          <Select
            v-model="position"
            :options="positionOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-primary font-medium text-sm">
            {{ $t('block.tabs.style.alignment') }}
          </label>
          <Select
            v-model="alignment"
            :options="alignmentOptions"
            :disabled="justify === 'justify'"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </div>
      </div>
    </div>

    <Divider />

    <!-- Tab list -->
    <CResourceTable
      :items="tabsTableRows"
      :fields="tableFields"
      primary-key="_rowKey"
      empty-message="—"
    >
      <template #header>
        <div class="flex items-center justify-between w-full">
          <label class="text-primary font-medium text-sm">{{ $t('block.tabs.title') }}</label>
          <Button
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            size="small"
            severity="secondary"
            @click="addTab"
          />
        </div>
      </template>

      <template #empty>
        {{ $t('block.tabs.noTabs') }}
      </template>

      <template #body-_order="{ data }">
        <div class="flex gap-1">
          <Button
            icon="pi pi-chevron-up"
            text
            severity="secondary"
            size="small"
            :disabled="data._index === 0"
            @click="moveTab(data._index, -1)"
          />
          <Button
            icon="pi pi-chevron-down"
            text
            severity="secondary"
            size="small"
            :disabled="data._index === tabs.length - 1"
            @click="moveTab(data._index, 1)"
          />
        </div>
      </template>

      <template #body-title="{ data }">
        <InputText
          :model-value="data.title || ''"
          :placeholder="`${$t('block.tabs.tab')} ${data._index + 1}`"
          class="w-full"
          size="small"
          @update:model-value="updateTab(data._index, 'title', $event)"
        />
      </template>

      <template #body-blockID="{ data }">
        <div class="flex items-center gap-1 min-w-0">
          <Select
            :model-value="data.blockID || ''"
            :options="getBlockOptions(data.blockID)"
            option-label="label"
            option-value="value"
            :placeholder="$t('block.tabs.selectBlock')"
            class="w-full min-w-0"
            size="small"
            show-clear
            @update:model-value="updateTab(data._index, 'blockID', $event)"
          />

          <Button
            v-if="data.blockID"
            v-tooltip.top="$t('block.tabs.tooltip.edit')"
            icon="pi pi-pencil"
            text
            severity="secondary"
            size="small"
            class="shrink-0"
            @click="editTabBlock(data._index)"
          />

          <Button
            v-else
            v-tooltip.top="$t('block.tabs.tooltip.addBlock')"
            icon="pi pi-plus"
            text
            severity="secondary"
            size="small"
            class="shrink-0"
            @click="createTabBlock(data._index)"
          />
        </div>
      </template>

      <Column
        field="lazy"
        :header-style="'width: 8rem'"
        header-class="text-center whitespace-nowrap"
        body-class="text-center whitespace-nowrap"
      >
        <template #header>
          <div class="flex items-center justify-center gap-1 whitespace-nowrap">
            <span class="p-datatable-column-title">
              {{ $t('block.tabs.table.columns.lazy.label') }}
            </span>
            <i
              class="pi pi-info-circle text-muted-color text-sm"
              v-tooltip.top="$t('block.tabs.table.columns.lazy.tooltip')"
            />
          </div>
        </template>

        <template #body="{ data }">
          <div class="flex justify-center">
            <Checkbox
              :model-value="data.lazy !== false"
              binary
              :input-id="`tabs-lazy-${data._index}`"
              @update:model-value="updateTab(data._index, 'lazy', $event)"
            />
          </div>
        </template>
      </Column>

      <Column
        field="_delete"
        :header-style="'width: 4rem'"
        header-class="whitespace-nowrap"
        body-class="whitespace-nowrap"
      >
        <template #body="{ data }">
          <div class="flex justify-center">
            <Button
              icon="pi pi-trash"
              text
              rounded
              size="small"
              severity="danger"
              @click="removeTab(data._index)"
            />
          </div>
        </template>
      </Column>
    </CResourceTable>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'

const { t } = useI18n()
const { CResourceTable } = components

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  blocks: { type: Array, default: () => [] },
})

const block = inject('blockDraft')

const emit = defineEmits(['edit-tab-block', 'create-tab-block'])

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

const tabs = computed(() => block.value.options?.tabs || [])

const tabsTableRows = computed(() =>
  tabs.value.map((tab, index) => ({
    ...tab,
    _index: index,
    _rowKey: `${index}-${tab.blockID || ''}-${tab.title || ''}`,
  })),
)

const tableFields = computed(() => [
  {
    key: '_order',
    header: '',
    style: 'width: 5rem',
    headerClass: 'whitespace-nowrap',
    bodyClass: 'whitespace-nowrap',
  },
  {
    key: 'title',
    header: t('block.tabs.table.columns.title.label'),
    style: 'min-width: 14rem',
    headerClass: 'whitespace-nowrap',
    bodyClass: 'whitespace-nowrap',
  },
  {
    key: 'blockID',
    header: t('block.tabs.table.columns.block.label'),
    style: 'min-width: 16rem',
    headerClass: 'whitespace-nowrap',
    bodyClass: 'whitespace-nowrap',
  },
])

function getBlockId(b) {
  return b.blockID && b.blockID !== '0' ? b.blockID : b.meta?.tempID || ''
}

function getBlockOptions(currentBlockID = '') {
  return (props.blocks || [])
    .filter(b => {
      const blockID = getBlockId(b)
      if (!blockID) return false
      if (b.kind === 'Tabs') return false
      if (blockID === getBlockId(block.value)) return false

      // Keep the current selection available, but prevent picking blocks already used by other tabs.
      return blockID === currentBlockID || !tabs.value.some(tab => tab.blockID === blockID)
    })
    .map(b => {
      const blockID = getBlockId(b)
      return {
        value: blockID,
        label: block.title || block.kind || blockID,
      }
    })
}

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

function updateStyle(key, value) {
  const style = { ...(block.value.options?.style || {}), [key]: value }
  updateOptions('style', style)
}

const appearance = computed({
  get: () => block.value.options?.style?.appearance || 'tabs',
  set: v => updateStyle('appearance', v),
})

const alignment = computed({
  get: () => block.value.options?.style?.alignment || 'center',
  set: v => updateStyle('alignment', v),
})

const justify = computed({
  get: () => block.value.options?.style?.justify || 'justify',
  set: v => updateStyle('justify', v),
})

const orientation = computed({
  get: () => block.value.options?.style?.orientation || 'horizontal',
  set: v => updateStyle('orientation', v),
})

const position = computed({
  get: () => block.value.options?.style?.position || 'start',
  set: v => updateStyle('position', v),
})

function addTab() {
  updateOptions('tabs', [...tabs.value, { title: '', blockID: '', lazy: true }])
}

function removeTab(index) {
  const items = [...tabs.value]
  items.splice(index, 1)
  updateOptions('tabs', items)
}

function updateTab(index, key, value) {
  const items = [...tabs.value]
  items[index] = {
    lazy: true,
    title: '',
    blockID: '',
    ...items[index],
    [key]: value,
  }
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

function editTabBlock(index) {
  const blockID = tabs.value[index]?.blockID
  if (!blockID) return
  emit('edit-tab-block', blockID)
}

function createTabBlock(index) {
  emit('create-tab-block', index)
}
</script>
