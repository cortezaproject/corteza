<template>
  <div class="flow-node group">
    <Handle type="target" :position="Position.Top" />

    <Card
      class="rounded-border overflow-hidden border border-surface hover:scale-[1.02] hover:shadow transition-all cursor-pointer"
      :class="{ '!border-primary scale-[1.02] shadow': selected }"
      :style="{ width: `${NODE_DIMENSIONS.WIDTH}px` }"
      :pt="{ body: { class: 'p-2' } }"
    >
      <template #content>
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-border flex items-center justify-center shrink-0">
            <TaqIcon
              :icon="data?.icon"
              :fallback="DEFAULT_ICONS.ACTION"
              class="text-lg text-primary"
            />
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-1">
              <div class="font-medium text-color flex-1 truncate">
                {{ data?.label || $t('builder.nodes.step') }}
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
              {{ data?.description || '' }}
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

    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<script setup>
import TaqIcon from '@/components/common/TaqIcon.vue'
import { NODE_DIMENSIONS } from '@/utils/flow-constants'
import { DEFAULT_ICONS } from '@cortezaproject/corteza-js-next/src/automation/types/icon'
import { Handle, Position } from '@vue-flow/core'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
})

const emit = defineEmits(['delete'])
const { t } = useI18n()
const menuRef = ref()
const menuOpen = ref(false)

const menuItems = [
  {
    label: t('builder.nodes.menu.delete'),
    icon: 'pi pi-trash',
    command: () => emit('delete', props.id),
  },
]

function toggle(event) {
  menuRef.value.toggle(event)
}
</script>
