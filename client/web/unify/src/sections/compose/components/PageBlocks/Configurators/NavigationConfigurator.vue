<template>
  <div class="flex flex-col gap-3">
    <Fieldset :legend="$t('block.navigation.appearance')">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <CFormGroup :label="$t('block.navigation.appearance')">
          <Select
            v-model="appearance"
            :options="appearanceOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('block.navigation.alignment')">
          <Select
            v-model="alignment"
            :options="alignmentOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup :label="$t('block.navigation.justify')">
          <Select
            v-model="justify"
            :options="justifyOptions"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>
      </div>
    </Fieldset>

    <Divider />

    <Fieldset :legend="$t('block.navigation.navigationItems')">
      <div class="flex flex-col gap-3">
        <div class="flex justify-end">
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
          :header="item.options?.item?.label || $t('block.navigation.newItem')"
          toggleable
        >
          <template #icons>
            <Button
              icon="pi pi-chevron-up"
              text
              rounded
              size="small"
              severity="secondary"
              :disabled="index === 0"
              @click="moveItem(index, -1)"
            />
            <Button
              icon="pi pi-chevron-down"
              text
              rounded
              size="small"
              severity="secondary"
              :disabled="index === navItems.length - 1"
              @click="moveItem(index, 1)"
            />
            <Button
              icon="pi pi-trash"
              text
              rounded
              size="small"
              severity="danger"
              @click="removeNavItem(index)"
            />
          </template>

          <div class="flex flex-col gap-2">
            <CFormGroup :label="$t('block.navigation.fieldLabel')">
              <InputText
                :model-value="item.options?.item?.label || ''"
                class="w-full"
                @update:model-value="updateNavItemOption(index, 'label', $event)"
              />
            </CFormGroup>

            <CFormGroup :label="$t('block.navigation.type')">
              <Select
                :model-value="item.type"
                :options="typeOptions"
                option-label="label"
                option-value="value"
                class="w-full"
                @update:model-value="updateNavItem(index, 'type', $event)"
              />
            </CFormGroup>

            <template v-if="item.type === 'url'">
              <CFormGroup :label="$t('block.navigation.url')">
                <CInputExpression
                  :ref="el => (urlInputs[index] = el)"
                  :model-value="item.options?.item?.url || ''"
                  dialect="interpolation"
                  :scope="scope"
                  placeholder="https://"
                  @update:model-value="updateNavItemOption(index, 'url', $event)"
                />
                <CExpressionHint :scope="scope" @insert="urlInputs[index]?.insert($event)" />
              </CFormGroup>
              <CFormGroup :label="$t('block.navigation.target')">
                <Select
                  :model-value="item.options?.item?.target || 'sameTab'"
                  :options="targetOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                  @update:model-value="updateNavItemOption(index, 'target', $event)"
                />
              </CFormGroup>
            </template>

            <template v-else-if="item.type === 'compose'">
              <CFormGroup :label="$t('block.navigation.composePage')">
                <Select
                  :model-value="item.options?.item?.pageID || null"
                  :options="pageOptions"
                  option-label="label"
                  option-value="value"
                  :placeholder="$t('block.navigation.selectPage')"
                  filter
                  show-clear
                  class="w-full"
                  @update:model-value="updateNavItemOption(index, 'pageID', $event)"
                />
              </CFormGroup>
              <CFormGroup :label="$t('block.navigation.target')">
                <Select
                  :model-value="item.options?.item?.target || 'sameTab'"
                  :options="targetOptions"
                  option-label="label"
                  option-value="value"
                  class="w-full"
                  @update:model-value="updateNavItemOption(index, 'target', $event)"
                />
              </CFormGroup>
              <CInputToggleCard
                :model-value="!!item.options?.item?.displaySubPages"
                :label="$t('block.navigation.displaySubPages')"
                :description="$t('block.navigation.displaySubPagesDescription')"
                @update:model-value="updateNavItemOption(index, 'displaySubPages', $event)"
              />
            </template>

            <template v-else-if="item.type === 'dropdown'">
              <CFormGroup :label="$t('block.navigation.dropdownItems')">
                <template #actions>
                  <Button
                    :label="$t('general.label.add')"
                    icon="pi pi-plus"
                    size="small"
                    severity="secondary"
                    text
                    @click="addDropdownItem(index)"
                  />
                </template>

                <div
                  v-for="(dItem, dIndex) in item.options?.item?.dropdown?.items || []"
                  :key="dIndex"
                  class="border border-surface rounded-border p-2 flex flex-col gap-2"
                >
                  <div class="flex items-center justify-between">
                    <CInputToggleCard
                      :model-value="!!dItem.delimiter"
                      :label="$t('block.navigation.delimiter')"
                      :description="$t('block.navigation.delimiterDescription')"
                      @update:model-value="updateDropdownItem(index, dIndex, 'delimiter', $event)"
                    />
                    <Button
                      icon="pi pi-trash"
                      text
                      rounded
                      size="small"
                      severity="danger"
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
              </CFormGroup>
            </template>

            <div class="grid grid-cols-2 gap-2">
              <CFormGroup :label="$t('block.navigation.textColor')">
                <CInputColorPicker
                  :model-value="item.options?.item?.textColor || ''"
                  show-text
                  @update:model-value="updateNavItemOption(index, 'textColor', $event)"
                />
              </CFormGroup>
              <CFormGroup :label="$t('block.navigation.backgroundColor')">
                <CInputColorPicker
                  :model-value="item.options?.item?.backgroundColor || ''"
                  show-text
                  @update:model-value="updateNavItemOption(index, 'backgroundColor', $event)"
                />
              </CFormGroup>
            </div>
          </div>
        </Panel>
      </div>
    </Fieldset>
  </div>
</template>

<script setup>
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { CInputColorPicker, CInputToggleCard } = components

const { t } = useI18n()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
})

const urlInputs = ref([])
const { scope } = useExpressionScope({ page: computed(() => props.page) })

const block = inject('blockDraft')

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
  { value: 'sameTab', label: t('block.navigation.sameTab') },
  { value: 'newTab', label: t('block.navigation.newTab') },
]

const pageOptions = computed(() => {
  return (pageStore.set || []).map(p => ({
    value: p.pageID,
    label: p.title || p.handle || p.pageID,
  }))
})

const navItems = computed(() => block.value.options?.navigationItems || [])

const appearance = computed({
  get: () => block.value.options?.display?.appearance || 'tabs',
  set: v => updateOptions('display', { ...block.value.options?.display, appearance: v }),
})

const alignment = computed({
  get: () => block.value.options?.display?.alignment || 'left',
  set: v => updateOptions('display', { ...block.value.options?.display, alignment: v }),
})

const justify = computed({
  get: () => block.value.options?.display?.justify || 'none',
  set: v => updateOptions('display', { ...block.value.options?.display, justify: v }),
})

function updateOptions(key, value) {
  if (!block.value.options) block.value.options = {}
  block.value.options[key] = value
}

function addNavItem() {
  const items = [
    ...navItems.value,
    {
      type: 'url',
      options: { enabled: true, item: { label: '', url: '', textColor: '', backgroundColor: '' } },
    },
  ]
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

function addDropdownItem(index) {
  const items = [...navItems.value]
  const dropdown = items[index].options?.item?.dropdown || { items: [] }
  const dItems = [
    ...(dropdown.items || []),
    { label: '', url: '', target: 'sameTab', delimiter: false },
  ]
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
