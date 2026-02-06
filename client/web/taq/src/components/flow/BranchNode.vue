<template>
  <div class="flow-node">
    <Handle type="target" :position="Position.Top" />

    <Card
      class="rounded-border overflow-hidden border border-surface hover:scale-[1.02] hover:shadow transition-all cursor-pointer"
      :class="{ '!border-primary scale-[1.02] shadow': selected }"
      :style="{ width: `${NODE_DIMENSIONS.WIDTH}px` }"
      :pt="{ body: { class: 'p-3' } }"
    >
      <template #content>
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-border flex items-center justify-center">
            <i :class="data?.icon || 'pi pi-sitemap'" class="text-lg text-primary" />
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-medium text-color">{{ data?.label || $t('builder.nodes.branch') }}</div>
            <div class="text-sm text-muted-color truncate">
              {{ data?.description || $t('builder.nodes.splitFlow') }}
            </div>
          </div>
        </div>
      </template>
    </Card>

    <!-- Single source handle for N-way branching -->
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<script setup>
import { NODE_DIMENSIONS } from '@/utils/flow-constants';
import { Handle, Position } from '@vue-flow/core';
defineProps({
  id: { type: String, required: true },
  data: { type: Object, required: true },
  selected: { type: Boolean, default: false },
})
</script>
