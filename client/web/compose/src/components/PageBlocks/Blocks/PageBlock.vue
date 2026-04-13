<template>
  <Card class="h-full overflow-hidden" :pt="cardPt">
    <template v-if="showHeader" #title>
      <div class="flex items-center gap-2">
        <div class="flex-1 truncate">
          <slot name="title">
            {{ block.title }}
          </slot>
        </div>

        <div v-if="showHeaderActions" class="flex items-center gap-1 ml-auto shrink-0">
          <Button
            v-if="block.options?.showRefresh"
            v-tooltip.bottom="$t('block.general.label.refresh')"
            icon="pi pi-refresh"
            text
            severity="secondary"
            size="small"
            class="p-1 border-0"
            @click="$emit('refreshBlock')"
          />
          <Button
            v-if="showMagnifyButton"
            v-tooltip.bottom="$t('block.general.label.magnify')"
            icon="pi pi-search-plus"
            text
            severity="secondary"
            size="small"
            class="p-1 border-0"
            @click="magnified = true"
          />
        </div>
      </div>
    </template>
    <template v-if="$slots.subtitle || block.description" #subtitle>
      <slot name="subtitle">
        {{ block.description }}
      </slot>
    </template>
    <template #content>
      <slot />
    </template>
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </Card>

  <!-- Magnify dialog -->
  <Dialog
    v-if="showMagnifyButton"
    v-model:visible="magnified"
    :modal="true"
    class="magnify-dialog"
    :style="{ width: '90vw', height: '90vh' }"
    :pt="{
      content: { class: 'flex-1 flex flex-col overflow-hidden p-0' },
      header: { class: 'pl-3 py-2 pr-2 border-b border-surface gap-1' },
      title: {
        style:
          'font-size: var(--p-card-title-font-size); font-weight: var(--p-card-title-font-weight)',
      },
      headerActions: { class: 'ml-auto' },
    }"
  >
    <template #header>
      <div class="flex items-center flex-1 min-w-0 pr-2">
        <span class="font-semibold text-lg truncate">{{ block.title }}</span>
        <Button
          v-if="block.options?.showRefresh"
          v-tooltip.bottom="$t('block.general.label.refresh')"
          icon="pi pi-refresh"
          text
          severity="secondary"
          size="small"
          class="p-1 border-0 ml-auto"
          @click="$emit('refreshBlock')"
        />
      </div>
    </template>

    <div class="flex flex-col flex-1 overflow-hidden h-full">
      <!-- Block content (toolbar + body, same as Card #content) -->
      <div class="p-0 flex-1 flex flex-col overflow-hidden">
        <slot />
      </div>

      <!-- Block footer (e.g. pagination) -->
      <div v-if="$slots.footer">
        <slot name="footer" />
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import { computed, ref, useSlots, watch, onBeforeUnmount } from 'vue'

const $slots = useSlots()

const emit = defineEmits(['refreshBlock'])

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
})

const magnified = ref(false)

const isPlain = computed(() => {
  return props.block.style?.wrap?.kind !== 'card'
})

// Map old Bootstrap variant names to PrimeVue CSS classes
const headerTextClass = computed(() => {
  const variant = props.block.style?.variants?.headerText
  const map = {
    dark: 'text-color',
    primary: 'text-primary',
    secondary: 'text-muted-color',
    success: 'text-green-500',
    warning: 'text-orange-500',
    danger: 'text-red-500',
  }
  return map[variant] || ''
})

const hasBorder = computed(() => {
  return props.block.style?.border?.enabled
})

const showHeader = computed(() => {
  return $slots.title || props.block.title || showHeaderActions.value
})

const magnifyOption = computed(() => props.block.options?.magnifyOption || '')
const showMagnifyButton = computed(() => !!magnifyOption.value)

const showHeaderActions = computed(() => {
  return props.block.options?.showRefresh || showMagnifyButton.value
})

const cardPt = computed(() => ({
  root: {
    class: [
      isPlain.value ? 'bg-transparent shadow-none' : '',
      hasBorder.value ? 'border border-surface' : '',
    ],
  },
  body: { class: 'p-0 flex-1 flex flex-col overflow-hidden gap-0' },
  caption: { class: ['pl-3 py-2 pr-2 border-b border-surface gap-1', headerTextClass.value] },
  content: { class: 'p-0 flex-1 flex flex-col overflow-hidden' },
}))

const refreshRate = computed(() => props.block.options?.refreshRate || 0)
let interval = null

watch(
  refreshRate,
  rate => {
    if (interval) clearInterval(interval)
    if (rate > 0) {
      interval = setInterval(() => {
        emit('refreshBlock')
      }, rate * 1000)
    }
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  if (interval) clearInterval(interval)
})
</script>

<style scoped>
.magnify-dialog :deep(.p-dialog-content) {
  min-height: 60vh;
}
</style>
