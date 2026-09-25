<template>
  <Dialog
    :visible="visible"
    modal
    :header="item ? item.title : ''"
    class="w-full max-w-2xl"
    data-testid="dashboard-event-dialog"
    @update:visible="$emit('update:visible', $event)"
  >
    <div v-if="item" class="flex flex-col gap-4">
      <p v-if="item.detail" class="text-base text-color break-words">{{ item.detail }}</p>

      <table class="w-full text-base">
        <tbody class="divide-y divide-surface">
          <tr v-for="row in rows" :key="row.label">
            <th
              class="py-1.5 pr-4 text-left font-medium text-muted-color whitespace-nowrap align-top w-40"
            >
              {{ row.label }}
            </th>
            <td class="py-1.5 break-all text-color">{{ row.value }}</td>
          </tr>
        </tbody>
      </table>

      <section v-if="metaRows.length" class="flex flex-col gap-1">
        <h3 class="text-base font-semibold text-color">{{ $t('dashboard.event.meta') }}</h3>
        <table class="w-full text-sm">
          <tbody class="divide-y divide-surface">
            <tr v-for="row in metaRows" :key="row.label">
              <th
                class="py-1 pr-4 text-left font-medium text-muted-color whitespace-nowrap align-top w-40"
              >
                {{ row.label }}
              </th>
              <td class="py-1 break-all text-color">{{ row.value }}</td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>

    <template #footer>
      <Button v-if="item && item.to" :label="openLabel" severity="secondary" @click="open" />
      <Button :label="$t('general.label.close')" @click="$emit('update:visible', false)" />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { filters } from '@planetcrust/human-vue'
import { resourceID, resourceTypeLabel } from './eventLabel'
import { useActors } from './useActors'

const { locFullDateTime } = filters

const props = defineProps({
  visible: { type: Boolean, default: false },
  // an attention/bucket item: { kind, title, detail, at, to, raw }
  item: { type: Object, default: null },
})
const emit = defineEmits(['update:visible'])

const { t } = useI18n()
const router = useRouter()
const actors = useActors(inject('$SystemAPI'))

watch(
  () => props.item,
  item => {
    const id = item?.raw?.actorID || item?.raw?.createdBy
    if (id) actors.resolve([id])
  },
  { immediate: true },
)

const rows = computed(() => {
  const it = props.item
  if (!it) return []
  const raw = it.raw || {}
  const out = [{ label: t('dashboard.event.when'), value: locFullDateTime(it.at) }]

  if (it.kind === 'session') {
    out.push(
      {
        label: t('dashboard.event.workflow'),
        value: raw.workflowName || raw.workflowHandle || raw.workflowID,
      },
      { label: t('dashboard.event.status'), value: raw.status || t('dashboard.runs.failed') },
      {
        label: t('dashboard.event.trigger'),
        value: [raw.eventType, raw.resourceType].filter(Boolean).join(' · '),
      },
      { label: t('dashboard.event.session'), value: raw.sessionID },
    )
    if (raw.createdBy)
      out.push({ label: t('dashboard.event.actor'), value: actors.label(raw.createdBy) })
  } else {
    out.push(
      { label: t('dashboard.event.resource'), value: resourceTypeLabel(raw.resource) },
      { label: t('dashboard.event.action'), value: raw.action },
    )
    const id = resourceID(raw.resource)
    if (id) out.push({ label: t('dashboard.event.resourceID'), value: id })
    if (raw.actorID && raw.actorID !== '0')
      out.push({ label: t('dashboard.event.actor'), value: actors.label(raw.actorID) })
    if (raw.description)
      out.push({ label: t('dashboard.event.description'), value: raw.description })
    if (raw.error) out.push({ label: t('dashboard.event.error'), value: raw.error })
  }

  return out.filter(r => r.value !== undefined && r.value !== null && r.value !== '')
})

const metaRows = computed(() => {
  const meta = props.item?.raw?.meta
  if (!meta || typeof meta !== 'object') return []
  return Object.entries(meta).map(([k, v]) => ({
    label: k,
    value: typeof v === 'object' ? JSON.stringify(v) : String(v),
  }))
})

const openLabel = computed(() => {
  const name = props.item?.to?.name || ''
  if (name.startsWith('automation.sessions')) return t('dashboard.event.openSession')
  if (name.startsWith('automation.taq')) return t('dashboard.event.openTaq')
  return t('dashboard.event.openLog')
})

function open() {
  emit('update:visible', false)
  router.push(props.item.to)
}
</script>
