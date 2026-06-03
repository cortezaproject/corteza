<template>
  <Teleport to="body">
    <div v-if="style" class="fixed z-[10000] w-[360px]" :style="style">
      <div class="bg-surface border border-surface rounded shadow-md overflow-hidden">
        <div class="p-2 border-b border-surface">
          <div class="font-medium text-sm text-color break-words">{{ title }}</div>
          <div
            v-if="description"
            class="mt-0.5 text-xs text-muted-color whitespace-pre-line break-words"
          >
            {{ description }}
          </div>
        </div>
        <table v-if="rows.length > 0" class="w-full table-fixed">
          <tr v-for="(row, idx) in rows" :key="idx" :class="row.class">
            <td
              v-for="(cell, ci) in row.cells"
              :key="ci"
              :colspan="row.cells.length === 1 ? 2 : 1"
              class="align-top px-2 py-1.5 text-xs break-words"
              v-html="cell"
            />
          </tr>
        </table>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
defineProps({
  style: { type: Object, default: null },
  title: { type: String, default: '' },
  description: { type: String, default: '' },
  rows: { type: Array, default: () => [] },
})
</script>
