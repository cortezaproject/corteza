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
    <div class="flex flex-col gap-4 min-h-0 flex-1">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="search"
          :placeholder="$t('agent.editor.tools.dialog.search')"
          class="w-full"
          data-testid="tool-dialog-search"
        />
      </IconField>

      <!-- A family grant is a standing rule, not a selection: it keeps meaning
           "every data tool" as tools are added. That is a different kind of
           statement from picking tools, so it is stated in one place rather
           than repeated in every section header. -->
      <div class="rounded-border border border-surface p-3 flex flex-col gap-3">
        <div>
          <span class="font-medium text-color text-sm block">
            {{ $t('agent.editor.tools.dialog.broad.label') }}
          </span>
          <!-- Nothing chosen is not "no access": the agent inherits what the
               person invoking it can already do. Saying so here is the
               difference between an empty panel that reads as broken and one
               that reads as a default. -->
          <span class="text-xs text-muted-color block">
            {{
              chosenCount
                ? $t('agent.editor.tools.dialog.broad.help')
                : $t('agent.editor.tools.dialog.inherits')
            }}
          </span>
        </div>

        <div v-for="key in GROUP_KEYS" :key="key" class="flex items-center gap-3">
          <Checkbox
            :model-value="!!groupGrant(key)"
            binary
            :disabled="disabled"
            :inputId="`grp-${key}`"
            :data-testid="`group-${key}`"
            @update:model-value="v => toggleGroup(key, v)"
          />
          <label :for="`grp-${key}`" class="flex-1 min-w-0 cursor-pointer">
            <span class="text-sm text-color block">
              {{ $t(`agent.editor.tools.dialog.group.${key}`) }}
            </span>
            <span class="text-xs text-muted-color block">
              {{ $t(`agent.editor.tools.dialog.group.${key}_help`) }}
            </span>
          </label>
          <span class="text-xs text-muted-color shrink-0">
            {{ $t('agent.editor.tools.dialog.broad.upTo') }}
          </span>
          <Select
            :model-value="groupGrant(key)?.maxRisk || 'read'"
            :options="riskCeilings"
            option-label="label"
            option-value="value"
            size="small"
            class="w-52 shrink-0"
            :disabled="disabled || !groupGrant(key)"
            @update:model-value="v => setGroupRisk(key, v)"
          />
        </div>
      </div>

      <!-- The scroller is pulled into the dialog's own padding so rows line up
           with the search field while the scrollbar sits in the gutter. -->
      <div class="overflow-y-auto flex-1 min-h-0 -mr-2 pr-2">
        <section v-for="d in domains" :key="d.key">
          <button
            type="button"
            class="w-full flex items-center gap-2 py-2 text-left border-b border-surface"
            :data-testid="`collapse-${d.key}`"
            @click="toggleCollapsed(d.key)"
          >
            <i
              :class="isOpen(d.key) ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"
              class="text-xs text-muted-color w-3"
            />
            <span class="text-sm font-semibold text-color">{{ d.label }}</span>
            <span class="text-xs text-muted-color ml-auto tabular-nums">{{ d.count }}</span>
          </button>

          <div v-if="isOpen(d.key)" class="pb-3">
            <div v-for="area in d.areas" :key="area.key">
              <div class="font-medium text-muted-color text-sm uppercase tracking-wide mt-3 mb-1">
                {{ area.label }}
              </div>

              <div
                v-for="tool in area.tools"
                :key="tool.name"
                class="flex items-start gap-3 py-2 border-b border-surface last:border-0"
              >
                <!-- A rail only where the choice loosens what the risk would
                     have done. The list stays quiet until there is something
                     to notice, and the thing to notice is an allow you set on
                     a tool that writes, not a deny. -->
                <span
                  class="w-0.5 self-stretch rounded-full shrink-0"
                  :class="loosened(tool) ? 'bg-amber-500' : 'bg-transparent'"
                />

                <span
                  :title="coveredByGroup(tool) ? $t('agent.editor.tools.dialog.fromFamily') : ''"
                >
                  <Checkbox
                    :model-value="chosen.has(tool.name) || coveredByGroup(tool)"
                    binary
                    :disabled="disabled || coveredByGroup(tool)"
                    :inputId="`tool-${tool.name}`"
                    class="mt-0.5"
                    @update:model-value="v => toggle(tool, v)"
                  />
                </span>

                <label
                  :for="`tool-${tool.name}`"
                  class="flex-1 min-w-0 cursor-pointer"
                  :class="{ 'opacity-50': effectiveMode(tool) === 'deny' }"
                >
                  <span class="text-sm text-color block">
                    <i
                      :class="[riskIcon(tool.risk), riskColour(tool.risk), 'text-xs mr-1.5']"
                      :title="riskLabel(tool.risk)"
                    />
                    {{ tool.title || tool.name }}
                  </span>
                  <span class="text-xs text-muted-color block line-clamp-2">
                    {{ summarise(tool.description) }}
                  </span>
                </label>

                <!-- Three states, always all three visible. The one that reads
                     as chosen is the effective mode: what the operator set, or
                     what the tool's risk decides until they set something. -->
                <SelectButton
                  :model-value="effectiveMode(tool)"
                  :options="modeOptions"
                  option-value="value"
                  :allow-empty="false"
                  size="small"
                  class="shrink-0"
                  :disabled="disabled || !(chosen.has(tool.name) || coveredByGroup(tool))"
                  @update:model-value="v => setMode(tool, v)"
                >
                  <template #option="{ option }">
                    <i :class="option.icon" :title="option.label" />
                  </template>
                </SelectButton>
              </div>
            </div>
          </div>
        </section>

        <div v-if="!domains.length" class="text-sm text-muted-color py-6 text-center">
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

// Which domains are folded shut, or null while the default still stands.
const collapsed = ref(null)

watch(
  () => props.visible,
  open => {
    if (!open) return
    search.value = ''
    collapsed.value = null
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
  {
    value: 'always',
    icon: 'pi pi-check-circle',
    label: t('agent.editor.tools.dialog.mode.always'),
  },
  { value: 'ask', icon: 'pi pi-question-circle', label: t('agent.editor.tools.dialog.mode.ask') },
  { value: 'deny', icon: 'pi pi-ban', label: t('agent.editor.tools.dialog.mode.deny') },
])

// What the tool will actually do: the mode set on it, or the one its risk
// decides until someone sets another. There is no fourth "default" segment —
// the default IS one of the three, shown as the chosen one.
function effectiveMode(tool) {
  return draft.value.get(tool.name) || defaultModeFor(tool.risk)
}

function defaultModeFor(risk) {
  return risk && risk !== 'read' ? 'ask' : 'always'
}

// An allow set by hand on a tool that writes: the one state where the choice
// permits more than the risk rule would have.
function loosened(tool) {
  return effectiveMode(tool) === 'always' && !!tool.risk && tool.risk !== 'read'
}

const riskCeilings = computed(() => [
  { value: 'read', label: t('agent.editor.tools.dialog.ceiling.read') },
  { value: 'write', label: t('agent.editor.tools.dialog.ceiling.write') },
  { value: 'destructive', label: t('agent.editor.tools.dialog.ceiling.destructive') },
])

const GROUP_KEYS = ['usage', 'configuring']

// Skills are attached to a tool, not chosen: the runtime injects one when the
// tool that triggers it is used. Offering them here invites a choice that
// changes nothing.
const HIDDEN_AREAS = new Set(['system_skill'])

// Sections are subjects, in the order an agent meets them: what it works with,
// then what it runs, then what it builds on, then who it touches.
//
// Grouping by `usage` and `configuring` instead put 17 tools in one section and
// 94 in the other, and split five subjects across both — every TAQ tool but
// `exec` in one section, `exec` in the other. A subject now appears once, with
// all of its tools, whichever group each one belongs to.
const DOMAINS = [
  { key: 'data', areas: ['compose_record', 'system_reminder'] },
  {
    key: 'automation',
    areas: ['automation_taq', 'automation_workflow', 'automation_trigger', 'automation_event'],
  },
  {
    key: 'structure',
    areas: ['compose_namespace', 'compose_module', 'compose_page', 'compose_chart'],
  },
  { key: 'people', areas: ['system_user', 'system_role', 'system_auth'] },
  { key: 'ai', areas: ['system_agent', 'system_chatbot'] },
  { key: 'workspace', areas: ['system_application', 'system_theme'] },
]

const PLACED_AREAS = new Set(DOMAINS.flatMap(d => d.areas))

// Reads first, then writes, then deletes: the harmless surface leads, and the
// one that cannot be taken back is last.
const RISK_ORDER = { read: 0, write: 1, destructive: 2 }

const domains = computed(() => {
  const q = search.value.trim().toLowerCase()
  const byArea = new Map()

  for (const tool of props.tools) {
    const areaKey = areaOf(tool.name)
    if (HIDDEN_AREAS.has(areaKey)) continue
    if (q && !matches(tool, q)) continue
    if (!byArea.has(areaKey)) byArea.set(areaKey, [])
    byArea.get(areaKey).push(tool)
  }

  // A tool family nobody has placed yet gets a section of its own rather than
  // disappearing from a dialog that claims to list everything.
  const spec = [
    ...DOMAINS,
    { key: 'other', areas: [...byArea.keys()].filter(k => !PLACED_AREAS.has(k)).sort() },
  ]

  return spec
    .map(d => {
      const areas = d.areas
        .filter(k => byArea.has(k))
        .map(k => ({
          key: k,
          label: t(`agent.editor.tools.dialog.area.${k}`, k),
          tools: [...byArea.get(k)].sort(byRiskThenName),
        }))

      return {
        key: d.key,
        label: t(`agent.editor.tools.dialog.domain.${d.key}`),
        areas,
        count: areas.reduce((n, a) => n + a.tools.length, 0),
      }
    })
    .filter(d => d.areas.length)
})

// Open where there is already something to see, so a configured agent shows
// its own grants and a fresh one still opens on something it can act on.
function defaultOpen(key) {
  const configured = domains.value.filter(d =>
    d.areas.some(a => a.tools.some(tool => named.value.has(tool.name))),
  )
  if (configured.length) return configured.some(d => d.key === key)
  return domains.value[0]?.key === key
}

// A search that only looked inside open sections would report nothing while
// showing a shut section holding the match.
function isOpen(key) {
  if (search.value.trim()) return true
  if (collapsed.value) return !collapsed.value.has(key)
  return defaultOpen(key)
}

function toggleCollapsed(key) {
  const next = new Set(
    collapsed.value ?? domains.value.filter(d => !defaultOpen(d.key)).map(d => d.key),
  )
  if (next.has(key)) next.delete(key)
  else next.add(key)
  collapsed.value = next
}

function byRiskThenName(a, b) {
  const ra = RISK_ORDER[a.risk] ?? RISK_ORDER.read
  const rb = RISK_ORDER[b.risk] ?? RISK_ORDER.read
  return ra - rb || (a.title || a.name).localeCompare(b.title || b.name)
}

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

function groupGrant(key) {
  const maxRisk = groupDraft.value.get(key)
  return maxRisk ? { group: key, maxRisk } : null
}

// A ceiling admits its own level and everything below it, which is what the
// runtime expands. A destructive tool under a write ceiling is not covered and
// stays available to grant by name.
function coveredByGroup(tool) {
  return (tool.groups || []).some(g => {
    const ceiling = groupDraft.value.get(g)
    return !!ceiling && (RISK_ORDER[tool.risk] ?? 0) <= (RISK_ORDER[ceiling] ?? 0)
  })
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

// Setting a mode on a tool a family already covers writes a named entry beside
// the family grant, which is how "all data tools, but ask before deleting" is
// said. The named entry wins: the runtime expands families after them.
//
// Choosing the mode the risk would have given anyway clears the override rather
// than pinning it, so a tool that should simply follow the rule keeps doing so
// if the rule ever changes.
function setMode(tool, mode) {
  const next = new Map(draft.value)
  const isDefault = mode === defaultModeFor(tool.risk)

  if (isDefault && !named.value.has(tool.name)) next.delete(tool.name)
  else next.set(tool.name, isDefault ? '' : mode)

  draft.value = next
}

function riskLabel(risk) {
  return t(`agent.editor.tools.dialog.risk.${risk || 'read'}`)
}

// The word is a tooltip rather than a badge: it repeats on every row and the
// shape carries it faster than reading does.
function riskIcon(risk) {
  if (risk === 'destructive') return 'pi pi-trash'
  if (risk === 'write') return 'pi pi-pencil'
  return 'pi pi-eye'
}

function riskColour(risk) {
  if (risk === 'destructive') return 'text-red-500'
  if (risk === 'write') return 'text-amber-500'
  return 'text-muted-color'
}

// Applying keeps whatever scope a grant already carried: the dialog chooses
// tools and modes, and namespace scoping is set elsewhere.
function apply() {
  const byName = new Map((props.grants || []).filter(g => g.name).map(g => [g.name, g]))
  const byGroup = new Map((props.grants || []).filter(g => g.group).map(g => [g.group, g]))

  const groupEntries = [...groupDraft.value.entries()].map(([group, maxRisk]) => ({
    ...(byGroup.get(group) || { description: '', allow: [] }),
    group,
    maxRisk,
  }))

  const toolEntries = [...draft.value.entries()].map(([name, permission]) => ({
    ...(byName.get(name) || { name, description: '', allow: [] }),
    name,
    permission,
  }))

  emit('apply', [...groupEntries, ...toolEntries])
  emit('update:visible', false)
}
</script>
