<template>
  <nav class="w-full h-full overflow-y-auto rounded-xl border border-surface bg-surface-50 dark:bg-surface-950">
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
              : 'hover:bg-surface-100 dark:hover:bg-surface-800'
          "
          @click="go(item)"
        >
          <span
            class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
            :class="isActive(item) ? 'bg-primary/10 ring-primary/30' : 'bg-emphasis ring-surface'"
          >
            <i :class="['pi', item.icon, 'text-xs', isActive(item) ? 'text-primary' : 'text-muted-color']" />
          </span>
          <span class="flex-1 min-w-0 truncate">{{ $t(item.labelKey) }}</span>
          <span
            v-if="item.badge != null"
            class="text-[10px] leading-none rounded-full px-1.5 py-1 font-medium"
            :class="badgeClass(item)"
          >
            {{ item.badge }}
          </span>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { DASHBOARD_NAV } from '@/sections/project/config/dashboard'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const nav = DASHBOARD_NAV

// An item is active when its route matches and — for the shared Events view —
// its category query matches (so the plain "All Events" item and the Category
// filter items highlight independently).
function isActive(item) {
  if (route.name !== item.route) return false
  return (item.query?.category ?? null) === (route.query.category ?? null)
}

function go(item) {
  router.push({
    name: item.route,
    params: { projectId: route.params.projectId },
    query: item.query || {},
  })
}

function badgeClass(item) {
  if (item.badgeTone === 'warning') return 'bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-300'
  if (item.badgeTone === 'muted' || item.badge === 0) return 'bg-surface-200 text-muted-color dark:bg-surface-700'
  return 'bg-primary/15 text-primary'
}
</script>
