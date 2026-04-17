<template>
  <div ref="nodeRef" class="flow-node group">
    <Card
      class="rounded-border overflow-hidden border hover:scale-[1.02] hover:shadow transition-all cursor-pointer"
      :class="[
        selected ? '!border-primary scale-[1.02] shadow' : '',
        !selected && traceActive ? '!border-green-500' : '',
        !selected && !traceActive ? 'border-surface' : '',
      ]"
      :style="{ width: `${NODE_DIMENSIONS.WIDTH}px` }"
      :pt="{ body: { class: 'p-2' } }"
      @mouseenter="onMouseEnter"
      @mouseleave="onMouseLeave"
    >
      <template #content>
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-border flex items-center justify-center shrink-0">
            <TaqIcon
              :icon="freshNode.data?.icon"
              :fallback="DEFAULT_ICONS.TRIGGER"
              class="text-lg text-primary"
            />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-1">
              <div class="font-medium text-color flex-1 truncate">
                {{ freshNode.data?.label || $t('builder.nodes.trigger') }}
              </div>
              <Button
                icon="pi pi-ellipsis-v"
                text
                size="small"
                severity="secondary"
                class="shrink-0 transition-opacity !w-7 !h-7"
                :class="menuOpen ? 'opacity-100' : 'opacity-0 group-hover:opacity-100'"
                @click.stop="toggle"
              />
            </div>
            <div class="text-sm text-muted-color truncate mb-1">
              {{ freshNode.data?.description || $t('builder.nodes.startAutomation') }}
            </div>
          </div>
        </div>
      </template>
    </Card>

    <Menu
      ref="menuRef"
      :model="menuItems"
      :popup="true"
      @show="menuOpen = true"
      @hide="menuOpen = false"
    />

    <StepPreviewPopover
      :visible="showPreview || alwaysShowPreview"
      :node="freshNode"
      :triggers="triggers"
      :nodes="nodes"
      :always-show="alwaysShowPreview"
      :node-el="nodeRef"
    />

    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<script setup>
import TaqIcon from '@/components/common/TaqIcon.vue'
import StepPreviewPopover from '@/components/flow/StepPreviewPopover.vue'
import { NODE_DIMENSIONS } from '@/utils/flow-constants'
import { DEFAULT_ICONS } from '@planetcrust/human-js/src/automation/types/icon'
import { Handle, Position } from '@vue-flow/core'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
  triggers: { type: Array, default: () => [] },
  nodes: { type: Array, default: () => [] },
  alwaysShowPreview: { type: Boolean, default: false },
  traceActive: { type: Boolean, default: false },
})

const emit = defineEmits(['delete', 'replace'])
const { t } = useI18n()
const nodeRef = ref(null)
const menuRef = ref()
const menuOpen = ref(false)
const showPreview = ref(false)
let hoverTimer = null

const freshNode = computed(() => {
  return props.nodes.find(n => n.id === props.id) || { type: 'trigger', data: props.data }
})

const menuItems = [
  {
    label: t('builder.nodes.menu.replace'),
    icon: 'pi pi-sync',
    command: () => emit('replace', props.id),
  },
  {
    label: t('builder.nodes.menu.delete'),
    icon: 'pi pi-trash',
    class: 'text-red-500',
    command: () => emit('delete', props.id),
  },
]

function toggle(event) {
  menuRef.value.toggle(event)
}

function onMouseEnter() {
  hoverTimer = setTimeout(() => {
    showPreview.value = true
  }, 150)
}

function onMouseLeave() {
  clearTimeout(hoverTimer)
  showPreview.value = false
}
</script>
