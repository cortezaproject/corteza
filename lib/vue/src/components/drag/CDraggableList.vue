<template>
  <VueDraggable
    v-model="items"
    :data-drag-key="dragKey"
    :tag="tag"
    :group="group"
    :disabled="disabled"
    :handle="handle"
    :draggable="itemSelector"
    :filter="filter"
    :prevent-on-filter="false"
    :clone="clone"
    :animation="150"
    :fallback-tolerance="5"
    :scroll-sensitivity="60"
    force-fallback
    fallback-on-body
    scroll
    bubble-scroll
    ghost-class="c-drag-ghost"
    chosen-class="c-drag-chosen"
    drag-class="c-drag-image"
    @add="e => emit('add', normalize(e))"
    @update="e => emit('update', normalize(e))"
    @remove="e => emit('remove', normalize(e))"
    @end="e => emit('end', normalize(e))"
  >
    <slot />
  </VueDraggable>
</template>

<script setup>
import { VueDraggable } from 'vue-draggable-plus'

// The one drag surface. Every list in the app that reorders — a board column, a
// form's rows, the app list — mounts this, so the grab threshold, the drop
// placeholder, the autoscroll and the lifted card look the same everywhere.
//
// Items are addressed by `itemSelector`, not by position, so a container may
// also hold an empty-state message or a header without them becoming
// draggable.
//
// Render the v-for unconditionally and make the empty state a sibling. Wrapping
// it in a v-if on the list's length tears the branch down as the last item is
// dragged out, while the sortable is still putting the dragged node back, and
// the node is left behind — visible in a list the data says is empty. It only
// bites lists an item can leave, but it costs nothing to write it the safe way.
//
// force-fallback replaces the browser's native drag image with a clone
// SortableJS positions itself. That is what makes touch work at all — native
// HTML5 drag never fires on a touchscreen — and what keeps the lifted card
// looking identical across browsers.

const items = defineModel({ type: Array, required: true })

defineProps({
  // Identifies this list to whoever handles a cross-list move: it arrives as
  // `fromKey`/`toKey` on the event, so a handler never reads the DOM.
  dragKey: { type: [String, Number], default: undefined },
  // A plain name, or SortableJS's { name, pull, put } for cross-list rules.
  group: { type: [String, Object], default: undefined },
  disabled: { type: Boolean, default: false },
  // A selector inside an item; unset means the whole item is the grab target.
  handle: { type: String, default: undefined },
  itemSelector: { type: String, default: '[data-drag-item]' },
  // Where a drag may NOT start, even on a whole-item surface. A form control
  // owns its own pointer gesture — dragging across an input is how you select
  // its text — and a row that is mostly inputs would otherwise be impossible to
  // type in. `prevent-on-filter` stays off so the control still takes the
  // click and places its caret; filtering only declines to start a drag.
  //
  // No `a` here on purpose: CAppList's rows ARE anchors, and SortableJS matches
  // the filter against the item root as well as its descendants, so listing it
  // would make that list undraggable rather than protecting anything.
  filter: {
    type: String,
    default: 'input, textarea, select, button, [contenteditable="true"]',
  },
  tag: { type: String, default: 'div' },
  // What the receiving list is handed when an item crosses into it. The item
  // itself, because these lists hold live model objects: the library's own
  // default is JSON.parse(JSON.stringify(item)), which flattens a class to a
  // plain object and rewrites anything with a toJSON — a compose.Record's
  // values come back as the raw [{name, value}] array no viewer reads, and the
  // card lands blank. A list that pulls copies rather than moving items is the
  // one that wants a real clone here.
  clone: { type: Function, default: item => item },
})

const emit = defineEmits(['add', 'update', 'remove', 'end'])

function normalize(e) {
  return {
    item: e.data,
    index: e.newIndex,
    oldIndex: e.oldIndex,
    fromKey: e.from?.dataset?.dragKey,
    toKey: e.to?.dataset?.dragKey,
  }
}
</script>

<style>
/* Deliberately global: the ghost and the lifted clone are moved out of the
   component's subtree (and onto <body>) while a drag is in flight, so scoped
   styles never reach them. */
.c-drag-ghost {
  opacity: 1;
  background: transparent;
  outline: 2px dashed var(--p-primary-color);
  outline-offset: -2px;
  border-radius: var(--p-content-border-radius);
}

.c-drag-ghost > * {
  visibility: hidden;
}

.c-drag-chosen {
  cursor: grabbing;
}

.c-drag-image {
  cursor: grabbing;
  opacity: 0.9;
  background: var(--p-content-background);
  border-radius: var(--p-content-border-radius);
  box-shadow: var(--p-overlay-popover-shadow, 0 8px 24px rgb(0 0 0 / 25%));
}
</style>
