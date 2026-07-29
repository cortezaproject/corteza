<template>
  <div class="flex flex-col gap-4">
    <!-- Leads with the reassurance, because the anxiety a first publish
         actually carries is "what am I about to break" — and the answer is
         nothing. -->
    <div class="rounded-xl border border-primary/40 bg-primary/5 p-4 flex gap-3.5 items-start">
      <span
        class="inline-flex items-center justify-center w-9 h-9 rounded-lg shrink-0 bg-primary text-primary-contrast"
      >
        <i class="pi pi-cloud-upload" />
      </span>
      <div class="min-w-0">
        <h3 class="font-medium">{{ $t('project.publish.firstRun.heading') }}</h3>
        <p class="text-sm text-muted-color mt-1">
          {{ $t('project.publish.firstRun.blurb', { name: projectName }) }}
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <section class="rounded-xl border border-surface bg-surface p-4">
        <h4 class="text-xs font-medium uppercase tracking-wide text-muted-color mb-3">
          {{ $t('project.publish.firstRun.whatGoesLive') }}
        </h4>
        <div v-if="inventory.length" class="flex flex-wrap gap-2">
          <ResourceCount
            v-for="entry in inventory"
            :key="entry.kind"
            :kind="entry.kind"
            :count="entry.count"
          />
        </div>
        <p v-else class="text-sm text-muted-color">
          {{ $t('project.publish.firstRun.nothingBuilt') }}
        </p>
      </section>

      <section class="rounded-xl border border-surface bg-surface p-4">
        <h4 class="text-xs font-medium uppercase tracking-wide text-muted-color mb-3">
          {{ $t('project.publish.firstRun.whoGetsAccess') }}
        </h4>
        <ul class="flex flex-col gap-2">
          <li v-for="member in members" :key="member.userId" class="flex items-center gap-2.5">
            <span
              class="inline-flex items-center justify-center w-6 h-6 rounded-full shrink-0 text-[10px] font-medium bg-primary/10 text-primary ring-1 ring-primary/25"
            >
              {{ member.initials }}
            </span>
            <span class="flex-1 min-w-0 text-sm truncate">{{ member.name }}</span>
            <span class="text-xs text-muted-color shrink-0">{{ member.role }}</span>
          </li>
        </ul>
        <p class="text-xs text-muted-color mt-3">
          {{ $t('project.publish.firstRun.accessNote') }}
        </p>
      </section>
    </div>

    <section class="rounded-xl border border-surface bg-surface p-4">
      <h4 class="text-xs font-medium uppercase tracking-wide text-muted-color mb-3">
        {{ $t('project.publish.firstRun.whatItDoes') }}
      </h4>
      <ul class="flex flex-col gap-2.5">
        <li class="flex gap-2.5 text-sm">
          <i class="pi pi-cloud-upload text-muted-color mt-0.5" />
          <span>{{ $t('project.publish.firstRun.becomesActive', { name: projectName }) }}</span>
        </li>
        <li class="flex gap-2.5 text-sm">
          <i class="pi pi-chart-line text-muted-color mt-0.5" />
          <span>{{ $t('project.publish.firstRun.getsDashboard') }}</span>
        </li>
        <!-- The one thing nobody knows the first time: a live project is not
             edited in place. Said here, where it is about to become true. -->
        <li class="flex gap-2.5 text-sm">
          <i class="pi pi-inbox text-muted-color mt-0.5" />
          <span>{{ $t('project.publish.firstRun.laterChanges') }}</span>
        </li>
      </ul>
    </section>
  </div>
</template>

<script setup>
// The first publish — its own screen, not the revision flow with its middle
// emptied out (ruled 2026-07-29). Nothing is being replaced and nothing can be
// lost, so there is no diff to read and no migration to decide; the questions
// worth answering are what this project contains and who is about to get it.
// It reports what exists and never judges readiness: everything shown is real
// persisted data, so it cannot be wrong.
import ResourceCount from './ResourceCount.vue'
import { rolePreset } from '@/sections/project/config/roles'
import { useProjectUsersStore } from '@/sections/project/stores/users'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  projectName: { type: String, default: '' },
  // [{ kind, count }] — same graph the Build canvas draws.
  inventory: { type: Array, default: () => [] },
  // Raw project members ({ userId, role }).
  projectMembers: { type: Array, default: () => [] },
})

const { t } = useI18n()
const usersStore = useProjectUsersStore()

const members = computed(() =>
  props.projectMembers.map(member => {
    const user = usersStore.findUser(member.userId)
    const name = user?.name || user?.email || member.userId
    return {
      userId: member.userId,
      name,
      role: t(rolePreset(member.role).labelKey),
      initials: initialsOf(name),
    }
  }),
)

// Falls back to the leading character of whatever we have — a userID still
// yields something stable rather than an empty circle.
function initialsOf(name) {
  return String(name)
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map(part => part[0])
    .join('')
    .toUpperCase()
}
</script>
