<template>
  <div
    class="grid grid-cols-[1fr_auto_1fr] items-center p-3 gap-2 bg-[--p-content-background] border-t border-surface"
  >
    <div class="flex items-center justify-start gap-2">
      <slot name="start" />
    </div>
    <div class="flex items-center justify-center gap-1">
      <slot name="center" />
    </div>
    <div class="flex items-center justify-end gap-2">
      <slot name="end" />
    </div>
  </div>
</template>
