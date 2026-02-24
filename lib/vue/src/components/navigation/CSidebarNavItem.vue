<template>
  <div>
    <Divider v-if="showDivider" class="my-1" />

    <!-- Root level items styled as PrimeVue text buttons -->
    <Button
      v-if="depth === 0"
      severity="secondary"
      text
      class="w-full justify-between font-medium"
      :class="{ 'bg-highlight !text-primary': itemIsActive }"
      @click="handleAction($event)"
    >
      <template #default>
        <div class="flex items-center gap-2 flex-1 cursor-pointer overflow-hidden">
          <i v-if="itemIcon" :class="itemIcon" />
          <span class="truncate">{{ label }}</span>
        </div>
        <i
          v-if="hasChildren"
          class="text-xs cursor-pointer transition-transform ml-auto p-2 -m-2"
          :class="isExpanded ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"
          @click.stop="$emit('toggle', node[idKey])"
        />
      </template>
    </Button>

    <!-- Child items -->
    <button
      v-else
      class="flex items-center gap-2 w-full text-left px-3 py-1.5 rounded-md transition-colors hover:bg-emphasis ml-2 mb-1"
      :class="itemIsActive ? 'bg-highlight !text-primary' : 'text-color-muted'"
      @click="handleAction($event)"
    >
      <slot name="item" :item="node" :depth="depth" :active="itemIsActive">
        <i v-if="itemIcon" :class="[itemIcon, 'text-xs opacity-60']" />
        <span class="truncate flex-1">{{ label }}</span>
        <i
          v-if="hasChildren"
          class="text-xs opacity-60 cursor-pointer transition-transform ml-auto p-2 -m-2"
          :class="isExpanded ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"
          @click.stop="$emit('toggle', node[idKey])"
        />
      </slot>
    </button>

    <Transition
      v-if="hasChildren"
      @before-enter="onBeforeEnter"
      @enter="onEnter"
      @after-enter="onAfterEnter"
      @before-leave="onBeforeLeave"
      @leave="onLeave"
    >
      <div v-if="isExpanded" class="overflow-hidden ml-3 border-l">
        <CSidebarNavItem
          v-for="child in node._children"
          :key="child[idKey]"
          :node="child"
          :id-key="idKey"
          :label-key="labelKey"
          :icon-key="iconKey"
          :icon="icon"
          :divider-key="dividerKey"
          :route-key="routeKey"
          :active-id="activeId"
          :expanded-ids="expandedIds"
          :depth="depth + 1"
          @select="$emit('select', $event)"
          @toggle="$emit('toggle', $event)"
        >
          <template v-if="$slots.item" #item="slotProps">
            <slot name="item" v-bind="slotProps" />
          </template>
        </CSidebarNavItem>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useLink } from 'vue-router'

const props = defineProps({
  node: {
    type: Object,
    required: true,
  },
  idKey: {
    type: String,
    required: true,
  },
  labelKey: {
    type: String,
    required: true,
  },
  iconKey: {
    type: String,
    default: undefined,
  },
  icon: {
    type: String,
    default: undefined,
  },
  dividerKey: {
    type: String,
    default: undefined,
  },
  routeKey: {
    type: String,
    default: undefined,
  },
  activeId: {
    type: String,
    default: null,
  },
  expandedIds: {
    type: Object,
    required: true,
  },
  depth: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['select', 'toggle'])

const hasChildren = computed(() => props.node._children?.length > 0)
const isExpanded = computed(() => props.expandedIds.has(props.node[props.idKey]))

const nodeRoute = computed(() => {
  return props.routeKey ? props.node[props.routeKey] : null
})

// Route-aware active state via vue-router's useLink
const routeTo = computed(() => nodeRoute.value || '/')
const {
  isActive: routeIsActive,
  isExactActive: routeIsExactActive,
  navigate,
} = useLink({ to: routeTo })

// Unified active state: route-based when routeKey is set, otherwise prop-based
const itemIsActive = computed(() => {
  if (nodeRoute.value) {
    // Root items with children: highlight when any child route is active
    // Leaf items: highlight only on exact match
    return hasChildren.value ? routeIsActive.value : routeIsExactActive.value
  }
  return props.activeId === props.node[props.idKey]
})

// Unified click action: navigate (route) or emit select (non-route)
function handleAction(e) {
  // Auto-expand children when clicking a parent
  if (hasChildren.value && !isExpanded.value) {
    emit('toggle', props.node[props.idKey])
  }

  if (nodeRoute.value) {
    navigate(e)
  } else {
    emit('select', props.node)
  }
}

// Expand/collapse transition hooks
function onBeforeEnter(el) {
  el.style.height = '0'
  el.style.opacity = '0'
}
function onEnter(el) {
  el.style.transition = 'height 150ms ease-out, opacity 150ms ease-out'
  el.style.height = el.scrollHeight + 'px'
  el.style.opacity = '1'
}
function onAfterEnter(el) {
  el.style.height = ''
  el.style.transition = ''
}
function onBeforeLeave(el) {
  el.style.height = el.scrollHeight + 'px'
  el.style.opacity = '1'
}
function onLeave(el) {
  el.style.transition = 'height 150ms ease-in, opacity 150ms ease-in'
  // Force reflow
  el.offsetHeight
  el.style.height = '0'
  el.style.opacity = '0'
}

const showDivider = computed(() => {
  return props.dividerKey && props.node[props.dividerKey]
})

const label = computed(() => {
  return props.node[props.labelKey] || props.node.handle || props.node[props.idKey]
})

const itemIcon = computed(() => {
  if (props.iconKey && props.node[props.iconKey]) {
    return props.node[props.iconKey]
  }
  return props.icon
})
</script>
