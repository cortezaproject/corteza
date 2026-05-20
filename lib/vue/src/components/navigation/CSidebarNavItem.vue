<template>
  <div class="w-full">
    <Divider v-if="showDivider" class="my-1" />

    <!-- Root level items styled as PrimeVue text buttons -->
    <Button
      severity="secondary"
      text
      size="small"
      class="w-full justify-between"
      :class="{
        'bg-highlight !text-primary': itemIsActive,
        '!w-[calc(100%-0.5rem)] ml-2': depth > 0,
      }"
      @click="handleAction($event)"
    >
      <template #default>
        <div class="flex items-center gap-2 flex-1 cursor-pointer overflow-hidden">
          <img
            v-if="itemIconIsImage"
            :src="itemIcon"
            class="h-4 w-4 object-contain"
            :alt="label"
          />
          <i v-else-if="itemIcon" :class="itemIcon" />
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
          :match-type="matchType"
          class="my-1"
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
import { useRoute, useRouter } from 'vue-router'

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
  matchType: {
    type: String,
    default: 'exact', // 'exact' or 'prefix'
  },
})

const emit = defineEmits(['select', 'toggle'])

const hasChildren = computed(() => props.node._children?.length > 0)
const isExpanded = computed(() => props.expandedIds.has(props.node[props.idKey]))

const nodeRoute = computed(() => {
  return props.routeKey ? props.node[props.routeKey] : null
})

const router = useRouter()
const currentRoute = useRoute()

// Safely resolve the item's route, returning null if params are missing
const resolvedRoute = computed(() => {
  if (!nodeRoute.value) return null
  try {
    return router.resolve(nodeRoute.value)
  } catch {
    return null
  }
})

// Unified active state: route-based when routeKey is set, otherwise prop-based
const itemIsActive = computed(() => {
  if (nodeRoute.value) {
    if (props.matchType === 'prefix') {
      if (nodeRoute.value.name && currentRoute.name?.startsWith(nodeRoute.value.name)) {
        return true
      }
    }

    const resolved = resolvedRoute.value
    if (!resolved) return false

    if (hasChildren.value) {
      // Root items with children: highlight when current path starts with this item's path
      return currentRoute.path.startsWith(resolved.path)
    }
    // Leaf items: highlight only on exact match
    return currentRoute.path === resolved.path
  }
  return props.activeId === props.node[props.idKey]
})

// Unified click action: navigate (route) or emit select (non-route)
function handleAction() {
  if (nodeRoute.value) {
    // Items with route: auto-expand but never collapse
    if (hasChildren.value && !isExpanded.value) {
      emit('toggle', props.node[props.idKey])
    }
    router.push(nodeRoute.value)
  } else {
    // Items without route: toggle collapse/expand
    if (hasChildren.value) {
      emit('toggle', props.node[props.idKey])
    }
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

const itemIconIsImage = computed(() => {
  const v = itemIcon.value
  return typeof v === 'string' && (v.startsWith('http') || v.startsWith('/'))
})
</script>
