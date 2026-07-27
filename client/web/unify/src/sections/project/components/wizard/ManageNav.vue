<template>
  <nav class="w-full h-full overflow-y-auto rounded-xl border border-surface bg-surface">
    <div class="p-3 flex flex-col gap-4">
      <div v-for="section in nav" :key="section.key" class="flex flex-col gap-0.5">
        <div
          class="px-3 pt-1 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-color"
        >
          {{ $t(section.labelKey) }}
        </div>
        <button
          v-for="item in section.items"
          :key="item.key"
          type="button"
          class="w-full text-left rounded-md px-3 py-2 flex items-center gap-2 text-sm transition-colors"
          :class="isActive(item) ? 'bg-primary/10 text-primary font-medium' : 'hover:bg-emphasis'"
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
import { MANAGE_NAV } from '@/sections/project/config/manageNav'
import { useEventsStore } from '@/sections/project/stores/events'

// Prop-driven, not route-driven: Manage & Monitor lives inside one Wizard
// tab and switches content via a `section` query param the Wizard view owns
// (see config/manageNav.js for why), so this rail takes the active key as a
// prop and emits the selection back up instead of pushing routes itself —
// unlike DashboardNav, which the live (routed) dashboard uses.
const props = defineProps({
  activeKey: { type: String, default: '' },
})
const emit = defineEmits(['select'])

const eventsStore = useEventsStore()

const nav = MANAGE_NAV

// Item keys are unique across the whole nav (unlike DashboardNav, which
// disambiguates category items on a shared route's `:category` param), so a
// plain key match is enough here.
function isActive(item) {
  return item.key === props.activeKey
}

// Icon badge classes — same idiom as DashboardNav: category items reuse the
// SAME colored badge as their screen title (CATEGORY_CONFIG.badge —
// bg/ring/icon/text); the rest use the neutral badge that tints to primary
// when active (the row background, not the icon, marks the active item).
function iconBox(item) {
  if (item.category) {
    const { badge } = CATEGORY_CONFIG[item.category]
    return { wrap: [badge.bg, badge.ring], icon: [badge.icon, badge.text] }
  }
  return {
    wrap: ['bg-emphasis', 'ring-surface'],
    icon: ['pi', item.icon, 'text-muted-color'],
  }
}

function go(item) {
  emit('select', item.key)
}

// Category items show a LIVE count from the events store, same as
// DashboardNav; everything else keeps its static config badge (may be null
// → no pill).
function badgeValue(item) {
  if (item.category) return eventsStore.countByCategory(item.category)
  return item.badge ?? null
}

function badgeClass(item) {
  if (item.badgeTone === 'muted' || badgeValue(item) === 0) return 'bg-emphasis text-muted-color'
  return 'bg-primary/15 text-primary'
}
</script>
