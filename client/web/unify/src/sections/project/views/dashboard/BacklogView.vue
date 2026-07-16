<template>
  <!-- Backlog: every category item carrying at least one backlog tag,
       aggregated across all five categories and grouped by tag. Title-bar
       idiom matches CategoryView (leading badge + title/description); the
       body below is plain sections rather than a CResourceList since rows
       here span multiple categories and only need to be scanned, not sorted
       or filtered (v1 keeps this simple). -->
  <div class="flex flex-col h-full min-w-0 overflow-y-auto">
    <header class="shrink-0 border-b border-surface px-4 py-3 flex items-center gap-3">
      <KindIcon :config="BADGE" size="xl" />
      <div class="min-w-0">
        <h2 class="text-xl font-semibold text-color truncate">
          {{ $t('project.dashboard.backlog.title') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('project.dashboard.backlog.desc') }}</p>
      </div>
    </header>

    <div class="flex-1 min-h-0 overflow-y-auto p-4 flex flex-col gap-6">
      <!-- First load (or a project switch): skeleton sections rather than a
           zeroed-out page. -->
      <div v-if="store.loading" class="flex flex-col gap-6">
        <div v-for="n in 3" :key="n" class="flex flex-col gap-2">
          <span class="h-5 w-24 rounded bg-emphasis block animate-pulse motion-reduce:animate-none" />
          <span
            v-for="m in 3"
            :key="m"
            class="h-10 rounded-lg bg-emphasis block animate-pulse motion-reduce:animate-none"
          />
        </div>
      </div>

      <!-- Nothing tagged yet across any category. -->
      <CEmptyState v-else-if="!backlogGroups.length">
        <p class="font-medium text-color">{{ $t('project.dashboard.backlog.empty') }}</p>
        <p class="text-xs mt-1">{{ $t('project.dashboard.backlog.emptyHint') }}</p>
      </CEmptyState>

      <!-- One section per tag, alphabetical; items within a tag lead with the
           soonest-due first (blank dates sort last). An item with several
           tags appears once per tag it carries. -->
      <template v-else>
        <section v-for="group in backlogGroups" :key="group.tag" class="flex flex-col gap-2">
          <header class="flex items-center gap-2">
            <span
              class="font-mono text-xs px-2 py-0.5 rounded border border-surface bg-surface text-color"
            >
              {{ group.tag }}
            </span>
            <span class="text-xs text-muted-color">
              {{ $t('project.dashboard.backlog.itemCount', group.items.length) }}
            </span>
          </header>

          <div class="rounded-lg border border-surface divide-y divide-surface overflow-hidden bg-surface">
            <RouterLink
              v-for="item in group.items"
              :key="`${item.category}-${item.id}`"
              :to="{ name: 'project.overview.category', params: { projectId, category: item.category } }"
              class="flex items-center gap-3 px-3 py-2 hover:bg-emphasis transition-colors"
            >
              <KindIcon :config="CATEGORY_CONFIG[item.category].badge" size="md" />
              <span class="flex-1 min-w-0 font-medium text-color truncate">{{ item.title || '—' }}</span>
              <EventBadge :value="item.status" variant="status" />
              <span class="text-sm text-muted-color w-24 text-right shrink-0">
                {{ formatDate(item.dateDue) }}
              </span>
              <UserCell v-if="ownerName(item)" :name="ownerName(item)" />
            </RouterLink>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup>
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import UserCell from '@/sections/project/components/dashboard/UserCell.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { CATEGORY_CONFIG } from '@/sections/project/config/categories'
import { useEventsStore } from '@/sections/project/stores/events'
import { components } from '@planetcrust/human-vue'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const { CEmptyState } = components

defineOptions({ name: 'BacklogView' })

const route = useRoute()
const store = useEventsStore()

const projectId = computed(() => route.params.projectId)

// Title-bar badge — neutral (this page spans every category, so it doesn't
// borrow any single category's colour), same icon-square shape as
// CategoryView's per-category badge.
const BADGE = { icon: 'pi pi-th-large', bg: 'bg-emphasis', ring: 'ring-surface', text: 'text-color' }

// Which resolved field on a mapped event (see stores/events.js#mapRow) holds
// its primary owner's display name, per category — mirrors FIELDS.ownerKey in
// config/categories.js (kept local since that map isn't exported and this is
// the only other call site that needs it).
const OWNER_FIELD = {
  incident: 'issueOwner',
  feature: 'featureOwner',
  privacy: 'requestOwner',
  task: 'owner',
  review: 'reviewer',
}
const ownerName = item => item[OWNER_FIELD[item.category]] || ''

// Group every loaded category's items by backlog tag. `store.events` is
// already flattened across all five categories and kept live by
// DashboardLayout's load() on entry; note it caps each category's list at 200
// rows (stores/events.js#load), so a category with more open items than that
// is under-represented here — acceptable for v1, no fix needed.
const backlogGroups = computed(() => {
  const byTag = new Map()
  for (const item of store.events) {
    for (const tag of item.backlog || []) {
      if (!byTag.has(tag)) byTag.set(tag, [])
      byTag.get(tag).push(item)
    }
  }
  return [...byTag.entries()]
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([tag, items]) => ({
      tag,
      items: [...items].sort((a, b) => {
        const ad = a.dateDue ? new Date(a.dateDue).getTime() : Number.MAX_SAFE_INTEGER
        const bd = b.dateDue ? new Date(b.dateDue).getTime() : Number.MAX_SAFE_INTEGER
        return ad - bd
      }),
    }))
})

const formatDate = v => {
  if (!v) return '—'
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleDateString()
}
</script>
