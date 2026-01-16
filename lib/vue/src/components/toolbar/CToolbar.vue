<template>
  <div class="c-toolbar d-flex flex-column flex-sm-row bg-white gap-1 flex-wrap">
    <div
      class="d-flex align-items-center flex-sm-grow-0 flex-sm-shrink-0 flex-fill fill-width gap-1 flex-wrap"
    >
      <slot name="start" />
    </div>

    <div
      class="d-flex align-items-center flex-sm-grow-0 flex-sm-shrink-0 flex-fill fill-width gap-1 flex-wrap ml-sm-auto"
    >
      <slot
        name="center"
        class="d-flex align-items-center flex-sm-grow-0 flex-sm-shrink-0 flex-fill fill-width gap-1 flex-wrap ml-sm-auto"
      />
    </div>

    <div
      class="d-flex align-items-center flex-sm-grow-0 flex-sm-shrink-0 flex-fill fill-width gap-1 flex-wrap ml-sm-auto"
    >
      <slot name="end" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.c-toolbar {
  z-index: 1036;
  position: relative;
  min-height: 42px;
  padding: 0.75rem;
}

.fill-width > * {
  flex: 1;
}
</style>
