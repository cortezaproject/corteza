<template>
  <!-- One list per level. Nothing here moves while a page is carried: the
       list view draws the line where it will land and saves on release. -->
  <TransitionGroup tag="ul" :class="root ? 'tree-root' : 'tree-branch'" name="page-node">
    <li
      v-for="node in nodes"
      :key="node.key"
      class="page-node"
      :class="{ 'page-node-carried': tree.carriedKey.value === node.key }"
    >
      <div
        class="page-row group relative inline-flex items-center gap-3 w-fit min-w-56 max-w-full border rounded-md bg-[var(--p-content-background)] transition-colors hover:bg-emphasis cursor-pointer pl-3 pr-1 py-2"
        data-test-id="page-tree-node"
        :data-key="node.key"
        :data-parent="parentId"
        :data-depth="depth"
        :class="{ 'page-row-target': tree.targetKey.value === node.key }"
        @click="tree.onSelect(node)"
      >
        <!-- The single top-level page has nowhere to go; everything else does -->
        <span
          v-if="tree.canDrag.value && !(root && nodes.length === 1)"
          class="page-grip shrink-0"
          :title="$t('page.list.dragToMove')"
          @pointerdown="tree.onGripDown($event, node)"
          @click.stop
        />
        <span
          class="whitespace-normal break-words"
          :class="{ 'text-muted-color': !node.data.visible && !node.data.isRecordPage }"
          data-test-id="page-tree-title"
        >
          {{ node.label }}
        </span>
        <span
          v-if="node.data.description"
          class="text-xs text-muted-color truncate min-w-0"
          :title="node.data.description"
        >
          {{ node.data.description }}
        </span>

        <!-- A record page wears its module -->
        <Tag
          v-if="node.data.isRecordPage"
          v-tooltip.bottom="tree.recordPageLabel(node.data)"
          :value="tree.recordModuleName(node.data)"
          icon="pi pi-database"
          severity="secondary"
          rounded
          class="text-xs leading-none py-1 shrink-0"
        />
        <Tag
          v-else-if="!node.data.visible"
          v-tooltip.bottom="$t('page.list.hiddenTooltip')"
          :value="$t('page.notVisible')"
          icon="pi pi-eye-slash"
          severity="secondary"
          rounded
          class="text-xs leading-none py-1 shrink-0"
        />
        <Button
          v-if="tree.actionItems(node.data).length"
          v-tooltip.bottom="$t('general.label.actions')"
          icon="pi pi-ellipsis-v"
          text
          severity="secondary"
          size="small"
          class="shrink-0 ml-auto opacity-0 transition-opacity group-hover:opacity-100 focus-visible:opacity-100"
          :aria-label="$t('general.label.actions')"
          data-test-id="page-tree-actions"
          @click.stop="tree.showActionsMenu($event, node.data)"
        />
      </div>

      <PageTreeBranch
        v-if="node.children.length"
        :nodes="node.children"
        :parent-id="node.key"
        :depth="depth + 1"
      />
    </li>
  </TransitionGroup>
</template>

<script setup>
import { inject } from 'vue'

defineOptions({ name: 'PageTreeBranch' })

defineProps({
  nodes: { type: Array, required: true },
  // The page these nodes hang under; '0' at the top level
  parentId: { type: String, required: true },
  depth: { type: Number, default: 0 },
  root: { type: Boolean, default: false },
})

// Everything a level needs that is the same at every level — the drag, the
// row actions, the carried page — comes from the list view
const tree = inject('pageTree')
</script>

<style scoped>
/* The lists are flex columns of content-width rows. A branch is indented
 * under its parent's grip, where the parent's trunk drops from. */
.tree-root,
.tree-branch {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.75rem;
  margin: 0;
  padding: 0;
  list-style: none;
}

.tree-branch {
  margin-left: 1rem;
  padding-top: 0.75rem;
  padding-left: 15px;
}

.page-node {
  position: relative;
  transition: opacity 150ms ease;
}

/* After a drop the rows slide to their new places; a page arriving from
 * another level fades in there */
.page-node-move {
  transition: transform 220ms ease;
}

.page-node-enter-active {
  transition: opacity 220ms ease;
}

.page-node-enter-from {
  opacity: 0;
}

/* The carried page and everything beneath it stay in place, faded: the whole
 * block that moves, and the room it leaves. The card at the pointer is the
 * page alone with a count. */
.page-node-carried {
  opacity: 0.4;
}

.page-row {
  transition:
    border-color 120ms ease,
    box-shadow 120ms ease,
    background-color 120ms ease;
}

/* The row the carried page will go into */
.page-row-target {
  border-color: var(--p-primary-color);
  box-shadow: inset 0 0 0 1px var(--p-primary-color);
  background: color-mix(in srgb, var(--p-primary-color) 8%, var(--p-content-background));
}

/* Tree guides. A parent's trunk drops from under its grip to its children,
 * and the elbow hangs off each child row so it meets the row's middle
 * whatever the row holds. Below a top-level page the trunk runs on to the
 * next one, so the top-level pages hang off one spine; it passes behind the
 * opaque rows. */
.tree-branch > .page-node::before {
  content: '';
  position: absolute;
  left: -15px;
  top: -0.75rem;
  bottom: 0;
  border-left: 1px solid var(--page-tree-guide);
}

.tree-branch > .page-node:last-child::before {
  bottom: auto;
  height: calc(0.75rem + 24px);
}

.tree-root > .page-node:not(:last-child)::before {
  content: '';
  position: absolute;
  left: 1rem;
  top: 24px;
  bottom: -0.75rem;
  border-left: 1px solid var(--page-tree-guide);
}

.tree-branch > .page-node > .page-row::before {
  content: '';
  position: absolute;
  left: -15px;
  top: 50%;
  width: 11px;
  border-top: 1px solid var(--page-tree-guide);
}

/* The grab handle: two columns of dots, as on any sortable row. touch-action
 * none, so a finger on it drags the page rather than the screen. */
.page-grip {
  width: 10px;
  height: 16px;
  color: var(--p-text-muted-color);
  background-image: radial-gradient(circle, currentColor 1.2px, transparent 1.7px);
  background-size: 5px 5.5px;
  background-position: 0 0;
  cursor: grab;
  opacity: 0.6;
  touch-action: none;
}

.page-row:hover .page-grip {
  opacity: 1;
}
</style>
