<template>
  <div class="flex h-full flex-col">
    <ReminderList
      v-if="!store.editing"
      :reminders="store.reminders"
      class="flex-1"
      @create="startCreate"
      @edit="store.startEdit($event)"
      @dismiss="store.setDismissed($event.reminder, $event.value)"
      @delete="store.deleteReminder($event)"
    />

    <ReminderEdit
      v-else
      :reminder="store.editing"
      :processing="store.processing"
      class="flex-1"
      @back="store.clearEdit()"
      @dismiss="store.setDismissed($event.reminder, $event.value)"
      @save="store.saveReminder($event)"
    />
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useRoute } from 'vue-router'
import { useNamespaceStore } from '@planetcrust/human-vue'
import { useReminderStore } from '@/sections/compose/stores/reminder'
import ReminderEdit from './ReminderEdit.vue'
import ReminderList from './ReminderList.vue'

const route = useRoute()
const namespaceStore = useNamespaceStore()
const store = useReminderStore()
const $Auth = inject('$Auth', {})

const currentNamespace = computed(() => {
  const slug = route.params.slug?.toString()
  if (!slug) return null

  return namespaceStore.getByUrlPart(slug) || null
})

function startCreate (seed = {}) {
  store.startCreate({
    assignedTo: $Auth?.user?.userID,
    resource: seed.resource || (currentNamespace.value ? `namespace:${currentNamespace.value.namespaceID}` : `system:user:${$Auth?.user?.userID || '0'}`),
    payload: seed.payload || {},
    remindAt: seed.remindAt,
  })
}
</script>
