<template>
  <nav class="w-full h-full overflow-y-auto rounded-xl border border-surface bg-surface">
    <div class="p-3 flex flex-col gap-4">
      <div v-for="section in nav" :key="section.key" class="flex flex-col gap-0.5">
        <div class="px-3 pt-1 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-color">
          {{ $t(section.labelKey) }}
        </div>
        <button
          v-for="item in section.items"
          :key="item.key"
          type="button"
          class="w-full text-left rounded-md px-3 py-2 flex items-center gap-2 text-sm transition-colors"
          :class="
            isActive(item)
              ? 'bg-primary/10 text-primary font-medium'
              : 'hover:bg-emphasis'
          "
          @click="go(item)"
        >
          <span
            class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
            :class="iconBox(item).wrap"
          >
            <i :class="[...iconBox(item).icon, 'text-xs']" />
          </span>
          <span class="flex-1 min-w-0 truncate">{{ $t(item.labelKey) }}</span>
          <span
            v-if="badgeValue(item) != null"
            class="text-[10px] leading-none rounded-full px-1.5 py-1 font-medium"
            :class="badgeClass(item)"
          >
            {{ badgeValue(item) }}
          </span>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { DASHBOARD_NAV } from '@/sections/project/config/dashboard'
import { useBacklogItemsStore } from '@/sections/project/stores/backlogItems'
import { useEventsStore } from '@/sections/project/stores/events'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const eventsStore = useEventsStore()
const backlogStore = useBacklogItemsStore()

const nav = DASHBOARD_NAV

// An item is active when its route matches. Category items own a dedicated
// route and disambiguate on the `:category` param (so each category rail entry
// highlights independently).
function isActive(item) {
  if (item.category) {
    return route.name === 'project.overview.category' && route.params.category === item.category
  }
  return route.name === item.route
}

// Icon badge classes. Category items reuse the SAME colored badge as their
// screen title (CATEGORY_CONFIG.badge — bg/ring/icon/text); items with a fixed
// `iconTone` mirror their screen's header badge (neutral box + tinted icon,
// regardless of active state); the rest use the neutral badge that tints to
// primary when active.
function iconBox(item) {
  if (item.category) {
    const { badge } = CATEGORY_CONFIG[item.category]
    return { wrap: [badge.bg, badge.ring], icon: [badge.icon, badge.text] }
  }
  if (item.iconTone === 'primary') {
    return { wrap: ['bg-emphasis', 'ring-surface'], icon: ['pi', item.icon, 'text-primary'] }
  }
  // Badge is constant regardless of active state — the row background (see the
  // button :class) is what marks the active item, not the icon.
  return {
    wrap: ['bg-emphasis', 'ring-surface'],
    icon: ['pi', item.icon, 'text-muted-color'],
  }
}

function go(item) {
  const params = { projectId: route.params.projectId }
  if (item.category) params.category = item.category
  router.push({ name: item.route, params })
}

// Category items show a LIVE count from the events store; the backlog item
// shows a live count of open backlog items (stores/backlogItems.js#openCount);
// everything else keeps its static config badge (may be null → no pill).
function badgeValue(item) {
  if (item.category) return eventsStore.countByCategory(item.category)
  if (item.key === 'backlog') return backlogStore.openCount
  return item.badge ?? null
}

function badgeClass(item) {
  if (item.badgeTone === 'muted' || badgeValue(item) === 0) return 'bg-emphasis text-muted-color'
  return 'bg-primary/15 text-primary'
}
</script>
