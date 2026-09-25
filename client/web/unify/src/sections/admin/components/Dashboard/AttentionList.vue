<template>
  <div class="flex flex-col min-w-0">
    <div v-if="!items.length" class="text-base text-muted-color py-2">{{ emptyText }}</div>
    <ul v-else class="flex flex-col divide-y divide-surface">
      <li v-for="item in items" :key="item.id" class="flex items-start gap-3 py-2 min-w-0">
        <span
          class="mt-1.5 inline-block h-2 w-2 shrink-0 rounded-full"
          :style="{ backgroundColor: item.color }"
        />
        <div class="flex-1 min-w-0">
          <div class="flex items-baseline justify-between gap-3">
            <component
              :is="item.to ? 'router-link' : 'span'"
              :to="item.to"
              class="text-base font-medium text-color break-words min-w-0"
              :class="item.to ? 'hover:underline' : ''"
            >
              {{ item.title }}
            </component>
            <span class="text-sm text-muted-color whitespace-nowrap shrink-0">{{ item.time }}</span>
          </div>
          <div v-if="item.detail" class="text-sm text-muted-color break-words">
            {{ item.detail }}
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
defineProps({
  // [{ id, title, detail, time, to, color }]
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
})
</script>
