<template>
  <div class="flow-node">
    <Card
      class="rounded-border overflow-hidden border border-surface hover:scale-[1.02] hover:shadow transition-all cursor-pointer"
      :class="{ '!border-primary scale-[1.02] shadow': selected }"
      :style="{ width: `${NODE_DIMENSIONS.WIDTH}px` }"
      :pt="{ body: { class: 'p-3' } }"
    >
      <template #content>
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-border flex items-center justify-center">
            <i :class="data?.icon || 'pi pi-bolt'" class="text-lg text-primary" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-medium text-color">{{ data?.label || $t('builder.nodes.trigger') }}</div>
            <div class="text-sm text-muted-color truncate">
              {{ data?.description || $t('builder.nodes.startAutomation') }}
            </div>
          </div>
        </div>
      </template>
    </Card>

    <!-- Source handle at bottom (invisible) -->
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<script setup>
import { NODE_DIMENSIONS } from '@/utils/flow-constants';
import { Handle, Position } from '@vue-flow/core';
import Card from 'primevue/card';

defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
})
</script>
