<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('agent.editor.tools.dialog.title')"
    :style="{ width: '72rem', height: '86vh' }"
    :contentStyle="{ display: 'flex', flexDirection: 'column', minHeight: 0 }"
    :breakpoints="{ '1200px': '94vw' }"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-3 min-h-0 flex-1">
      <div class="flex items-center gap-3">
        <IconField class="flex-1">
          <InputIcon class="pi pi-search" />
          <InputText
            v-model="search"
            :placeholder="$t('agent.editor.tools.dialog.search')"
            class="w-full"
            data-testid="tool-dialog-search"
          />
        </IconField>
        <span class="text-sm text-muted-color whitespace-nowrap">
          {{ $t('agent.editor.tools.dialog.chosen', { count: chosenCount }) }}
        </span>
      </div>

      <!-- Nothing chosen is not "no access": the agent inherits what the
           person invoking it can already do. Saying so here is the difference
           between an empty panel that reads as broken and one that reads as a
           default. -->
      <Message v-if="!chosenCount" severity="secondary" :closable="false" class="!my-0">
        {{ $t('agent.editor.tools.dialog.inherits') }}
      </Message>

      <div class="overflow-y-auto flex-1 min-h-0 flex flex-col gap-5 pr-1">
        <section v-for="grp in groups" :key="grp.key">
          <!-- The whole group in one grant. It keeps meaning "every data tool"
               as tools are added, which naming them one by one does not. -->
          <div
            class="flex items-center gap-3 pb-2 mb-2 border-b border-surface sticky top-0 bg-surface-0 dark:bg-surface-900 z-10"
          >
            <Checkbox
              :model-value="!!groupGrant(grp.key)"
              binary
              :disabled="disabled"
              :inputId="`grp-${grp.key}`"
              @update:model-value="v => toggleGroup(grp.key, v)"
            />
            <label :for="`grp-${grp.key}`" class="flex-1 min-w-0 cursor-pointer">
              <span class="text-base font-semibold text-color block">{{ grp.label }}</span>
              <span class="text-xs text-muted-color block">{{ grp.help }}</span>
            </label>
            <Select
              :model-value="groupGrant(grp.key)?.maxRisk || 'read'"
              :options="riskCeilings"
              option-label="label"
              option-value="value"
              size="small"
              class="w-52"
              :disabled="disabled || !groupGrant(grp.key)"
              @update:model-value="v => setGroupRisk(grp.key, v)"
            />
          </div>

          <div v-for="area in grp.areas" :key="area.key" class="mb-3">
            <div class="text-sm font-medium text-muted-color mb-1 pl-1">{{ area.label }}</div>

            <div
              v-for="tool in area.tools"
              :key="tool.name"
              class="flex items-start gap-3 py-2 border-b border-surface last:border-0"
            >
              <Checkbox
                :model-value="chosen.has(tool.name)"
                binary
                :disabled="disabled || !!groupGrant(grp.key)"
                :inputId="`tool-${tool.name}`"
                class="mt-0.5"
                @update:model-value="v => toggle(tool, v)"
              />
              <label :for="`tool-${tool.name}`" class="flex-1 min-w-0 cursor-pointer">
                <span class="text-sm text-color block">{{ tool.title || tool.name }}</span>
                <span class="text-xs text-muted-color block line-clamp-2">
                  {{ summarise(tool.description) }}
                </span>
              </label>

              <Tag :severity="riskSeverity(tool.risk)" rounded class="mt-0.5 shrink-0">
                <span class="text-xs">{{ riskLabel(tool.risk) }}</span>
              </Tag>

              <!-- The default mode is stored as an empty string, which PrimeVue
                   reads as "no value" and renders blank; the placeholder is what
                   makes the row say what it will actually do. -->
              <Select
                :model-value="modeOf(tool)"
                :options="modeOptions"
                option-label="label"
                option-value="value"
                :placeholder="$t('agent.editor.tools.dialog.mode.default')"
                size="small"
                class="w-40 shrink-0"
                :disabled="disabled || !(chosen.has(tool.name) || groupGrant(grp.key))"
                @update:model-value="v => setMode(tool, v)"
              />
            </div>
          </div>
        </section>

        <div v-if="!groups.length" class="text-sm text-muted-color py-6 text-center">
          {{ $t('agent.editor.tools.dialog.noMatches') }}
        </div>
      </div>
    </div>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        @click="$emit('update:visible', false)"
      />
      <Button :label="$t('general.label.apply')" :disabled="disabled" @click="apply" />
    </template>
  </Dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  visible: { type: Boolean, default: false },
  tools: { type: Array, default: () => [] },
  // The agent's access.tools, as stored.
  grants: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:visible', 'apply'])

const { t } = useI18n()

const search = ref('')

// Editing happens on copies: closing with Cancel has to leave the agent
// exactly as it was.
const draft = ref(new Map()) // tool name -> permission
const groupDraft = ref(new Map()) // group key -> maxRisk

const chosen = computed(() => new Set(draft.value.keys()))

// The tools the agent was opened with, so an override added in this session can
// be told from an entry that was already there.
const named = ref(new Set())
const chosenCount = computed(() => draft.value.size + groupDraft.value.size)

watch(
  () => props.visible,
  open => {
    if (!open) return
    search.value = ''
    draft.value = new Map(
      (props.grants || []).filter(g => g.name).map(g => [g.name, g.permission || '']),
    )
    groupDraft.value = new Map(
      (props.grants || []).filter(g => g.group).map(g => [g.group, g.maxRisk || 'read']),
    )
    named.value = new Set(draft.value.keys())
  },
  { immediate: true },
)

const modeOptions = computed(() => [
  { value: '', label: t('agent.editor.tools.dialog.mode.default') },
  { value: 'always', label: t('agent.editor.tools.dialog.mode.always') },
  { value: 'ask', label: t('agent.editor.tools.dialog.mode.ask') },
  { value: 'deny', label: t('agent.editor.tools.dialog.mode.deny') },
])

const riskCeilings = computed(() => [
  { value: 'read', label: t('agent.editor.tools.dialog.ceiling.read') },
  { value: 'write', label: t('agent.editor.tools.dialog.ceiling.write') },
  { value: 'destructive', label: t('agent.editor.tools.dialog.ceiling.destructive') },
])

const GROUP_KEYS = ['usage', 'configuring']

// Two levels: the group is what a grant can name, and the areas inside it are
// what a person scans by. One level of either alone is a list of ninety.
const groups = computed(() => {
  const q = search.value.trim().toLowerCase()

  return GROUP_KEYS.map(key => {
    const byArea = new Map()

    for (const tool of props.tools) {
      if (!(tool.groups || []).includes(key)) continue
      if (q && !matches(tool, q)) continue

      const areaKey = areaOf(tool.name)
      if (!byArea.has(areaKey)) {
        byArea.set(areaKey, {
          key: areaKey,
          label: t(`agent.editor.tools.dialog.area.${areaKey}`, areaKey),
          tools: [],
        })
      }
      byArea.get(areaKey).tools.push(tool)
    }

    return {
      key,
      label: t(`agent.editor.tools.dialog.group.${key}`),
      help: t(`agent.editor.tools.dialog.group.${key}_help`),
      areas: [...byArea.values()].sort((a, b) => a.label.localeCompare(b.label)),
    }
  }).filter(g => g.areas.length)
})

function matches(tool, q) {
  return `${tool.name} ${tool.title || ''} ${tool.description || ''}`.toLowerCase().includes(q)
}

// compose_record_lookup -> compose_record; discovery_search -> discovery.
function areaOf(name) {
  const parts = String(name).split('_')
  return parts.length > 2 ? `${parts[0]}_${parts[1]}` : parts[0]
}

// The first sentence, which is what a tool description leads with. The rest is
// written for the model that has already decided to call it.
function summarise(description) {
  const text = String(description || '').trim()
  if (!text) return ''
  const end = text.search(/\.\s/)
  return end > 0 ? text.slice(0, end + 1) : text
}

function modeOf(tool) {
  return draft.value.get(tool.name) ?? ''
}

function groupGrant(key) {
  const maxRisk = groupDraft.value.get(key)
  return maxRisk ? { group: key, maxRisk } : null
}

function toggle(tool, on) {
  const next = new Map(draft.value)
  if (on) next.set(tool.name, '')
  else next.delete(tool.name)
  draft.value = next
}

function toggleGroup(key, on) {
  const next = new Map(groupDraft.value)
  if (on) next.set(key, next.get(key) || 'read')
  else next.delete(key)
  groupDraft.value = next
}

function setGroupRisk(key, maxRisk) {
  if (!groupDraft.value.has(key)) return
  const next = new Map(groupDraft.value)
  next.set(key, maxRisk)
  groupDraft.value = next
}

// Setting a mode on a tool the group already covers writes a named entry beside
// the group grant, which is how "all data tools, but ask before deleting" is
// said. The named entry wins: the runtime expands groups after them.
function setMode(tool, mode) {
  const next = new Map(draft.value)
  if (!mode && !named.value.has(tool.name)) next.delete(tool.name)
  else next.set(tool.name, mode || '')
  draft.value = next
}

function riskLabel(risk) {
  return t(`agent.editor.tools.dialog.risk.${risk || 'read'}`)
}

function riskSeverity(risk) {
  if (risk === 'destructive') return 'danger'
  if (risk === 'write') return 'warn'
  return 'secondary'
}

// Applying keeps whatever scope a grant already carried: the dialog chooses
// tools and modes, and namespace scoping is set elsewhere.
function apply() {
  const named = new Map((props.grants || []).filter(g => g.name).map(g => [g.name, g]))
  const grouped = new Map((props.grants || []).filter(g => g.group).map(g => [g.group, g]))

  const groupEntries = [...groupDraft.value.entries()].map(([group, maxRisk]) => ({
    ...(grouped.get(group) || { description: '', allow: [] }),
    group,
    maxRisk,
  }))

  const toolEntries = [...draft.value.entries()].map(([name, permission]) => ({
    ...(named.get(name) || { name, description: '', allow: [] }),
    name,
    permission,
  }))

  emit('apply', [...groupEntries, ...toolEntries])
  emit('update:visible', false)
}
</script>
