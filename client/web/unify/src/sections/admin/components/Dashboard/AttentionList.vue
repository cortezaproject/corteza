<template>
  <div class="flex flex-col min-w-0">
    <div v-if="!items.length" class="text-base text-muted-color py-2">{{ emptyText }}</div>
    <ul v-else class="flex flex-col divide-y divide-surface">
      <li v-for="item in items" :key="item.id">
        <button
          type="button"
          class="flex w-full items-start gap-3 py-2 min-w-0 text-left rounded hover:bg-emphasis transition-colors"
          @click="$emit('select', item)"
        >
          <span
            class="mt-2 inline-block h-2 w-2 shrink-0 rounded-full"
            :style="{ backgroundColor: item.color }"
          />
          <span class="flex-1 min-w-0">
            <span class="flex items-baseline justify-between gap-3">
              <span class="text-base font-medium text-color break-words min-w-0">
                {{ item.title }}
              </span>
              <span class="text-sm text-muted-color whitespace-nowrap shrink-0">
                {{ item.time }}
              </span>
            </span>
            <span v-if="item.detail" class="block text-sm text-muted-color break-words">
              {{ item.detail }}
            </span>
          </span>
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup>
defineProps({
  // [{ id, kind, title, detail, time, at, to, color, raw }]
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
})
defineEmits(['select'])
</script>
