<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('agent.editor.tools.dialog.title')"
    :style="{ width: '72rem', height: '86vh' }"
    :contentStyle="{ display: 'flex', flexDirection: 'column', minHeight: 0, flex: '1 1 auto' }"
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

      <!-- Nothing chosen is not "no access": the agent inherits what the person
           invoking it can already do. Saying so here is the difference between
           an empty panel that reads as broken and one that reads as a default. -->
      <Message v-if="!chosenCount" severity="secondary" :closable="false" class="!my-0">
        {{ $t('agent.editor.tools.dialog.inherits') }}
      </Message>

      <!-- The scroller is pulled into the dialog's own padding so rows line up
           with the search field while the scrollbar sits in the gutter. -->
      <div class="overflow-y-auto flex-1 min-h-0 -mr-2 pr-2">
        <section v-for="d in domains" :key="d.key">
          <div class="flex items-center gap-2 py-2 border-b border-surface">
            <button
              type="button"
              class="flex items-center gap-2 flex-1 min-w-0 text-left"
              :data-testid="`collapse-${d.key}`"
              @click="toggleCollapsed(d.key)"
            >
              <i
                :class="isOpen(d.key) ? 'pi pi-chevron-down' : 'pi pi-chevron-right'"
                class="text-xs text-muted-color w-3"
              />
              <span class="text-sm font-semibold text-color">{{ d.label }}</span>

              <!-- The split, not the total: a section reads as what it lets the
                   agent do, and the total is the sum of what is shown. -->
              <span class="flex items-center gap-2.5 text-xs text-muted-color">
                <span
                  v-for="c in sectionCounts(d)"
                  :key="c.mode"
                  class="flex items-center gap-1"
                  :title="modeLabel(c.mode)"
                >
                  <i :class="modeIcon(c.mode)" />
                  <span class="tabular-nums">{{ c.n }}</span>
                </span>
              </span>
            </button>

            <!-- One statement about the whole section: what every tool in it
                 should do. It reads back as Custom when they disagree, which is
                 the normal state once anything has been set by hand. -->
            <Select
              :model-value="sectionMode(d)"
              :options="sectionModes"
              option-label="label"
              option-value="value"
              size="small"
              class="w-56 shrink-0"
              :disabled="disabled"
              :aria-label="$t('agent.editor.tools.mode.section')"
              :data-testid="`section-mode-${d.key}`"
              @update:model-value="v => setSectionMode(d, v)"
            >
              <template #value="{ value }">
                <span class="flex items-center gap-2">
                  <i :class="modeIcon(value)" class="text-xs" />
                  {{ modeLabel(value) }}
                </span>
              </template>
              <template #option="{ option }">
                <span class="flex items-center gap-2">
                  <i :class="option.icon" class="text-xs" />
                  {{ option.label }}
                </span>
              </template>
            </Select>
          </div>

          <div v-if="isOpen(d.key)" class="pb-3">
            <div v-for="area in d.areas" :key="area.key">
              <div class="font-medium text-muted-color text-sm uppercase tracking-wide mt-3 mb-1">
                {{ area.label }}
              </div>

              <div
                v-for="tool in area.tools"
                :key="tool.name"
                class="group flex items-start gap-3 py-2 border-b border-surface last:border-0"
              >
                <!-- A rail only where the choice loosens what the risk would
                     have done. The list stays quiet until there is something to
                     notice, and the thing to notice is an allow set by hand on
                     a tool that writes. -->
                <span
                  class="w-0.5 self-stretch rounded-full shrink-0"
                  :class="loosened(tool) ? 'bg-amber-500' : 'bg-transparent'"
                />

                <div
                  class="flex-1 min-w-0"
                  :class="{ 'opacity-50': !toolOn(tool) }"
                  :title="covered(tool) ? $t('agent.editor.tools.dialog.fromFamily') : ''"
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
                </div>

                <!-- Blocked is the off state, so the row shows nothing until it
                     is reached for. Three segments carry the whole decision:
                     a second control for on and off would only disagree. -->
                <div class="tool-mode shrink-0" :data-blocked="!toolOn(tool)">
                  <SelectButton
                    :model-value="rowMode(tool)"
                    :options="toolModes"
                    option-value="value"
                    :allow-empty="false"
                    size="small"
                    :disabled="disabled"
                    @update:model-value="v => setMode(tool, v)"
                  >
                    <template #option="{ option }">
                      <i :class="option.icon" :title="option.label" />
                    </template>
                  </SelectButton>
                </div>
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
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="$emit('update:visible', false)"
        />
        <Button
          :label="$t('general.label.apply')"
          size="small"
          :disabled="disabled"
          @click="apply"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import {
  MODE_ICONS,
  RISK_ORDER,
  coveredByFamily,
  defaultModeFor,
  modeOf,
  tally,
} from '../toolAccess'

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

// A grant may name a whole family instead of a tool. This dialog chooses tools
// and no longer writes families, but it still has to know which tools an
// existing family already covers, and hand every family back untouched.
const families = ref([])

const chosen = computed(() => new Set(draft.value.keys()))

// The tools the agent was opened with, so an override added in this session can
// be told from an entry that was already there.
const named = ref(new Set())
const chosenCount = computed(() => draft.value.size + families.value.length)

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
    families.value = (props.grants || []).filter(g => g.group)
    named.value = new Set(draft.value.keys())
  },
  { immediate: true },
)

function modeIcon(mode) {
  return MODE_ICONS[mode] || MODE_ICONS.custom
}

function modeLabel(mode) {
  return t(`agent.editor.tools.mode.${mode || 'custom'}`)
}

function modeOption(value) {
  return { value, icon: modeIcon(value), label: modeLabel(value) }
}

const toolModes = computed(() => ['always', 'ask', 'deny'].map(modeOption))

// Custom is not a choice, it is what the section reads as once its tools
// disagree — offered so the control can show it back rather than lie.
const sectionModes = computed(() => ['always', 'ask', 'deny', 'custom'].map(modeOption))

// What the three segments read as. There is no fourth "default" segment — the
// default IS one of the three, shown as the chosen one.
function rowMode(tool) {
  return modeOf(tool, draft.value, families.value)
}

// Whether the agent has the tool at all.
function toolOn(tool) {
  return rowMode(tool) !== 'deny'
}

function covered(tool) {
  return coveredByFamily(tool, families.value)
}

// An allow set by hand on a tool that writes: the one state where the choice
// permits more than the risk rule would have.
function loosened(tool) {
  return rowMode(tool) === 'always' && !!tool.risk && tool.risk !== 'read'
}

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
      }
    })
    .filter(d => d.areas.length)
})

function toolsIn(d) {
  return d.areas.flatMap(a => a.tools)
}

// What the section reads as: the mode its tools agree on, or Custom.
function sectionMode(d) {
  const tools = toolsIn(d)
  if (!tools.length) return 'custom'

  const first = rowMode(tools[0])
  return tools.every(tool => rowMode(tool) === first) ? first : 'custom'
}

function sectionCounts(d) {
  return tally(toolsIn(d), draft.value, families.value)
}

// One click for a whole subject, which is the unit an agent is actually given:
// every record tool, or none of them.
function setSectionMode(d, mode) {
  if (mode === 'custom') return

  const next = new Map(draft.value)
  for (const tool of toolsIn(d)) applyMode(next, tool, mode)
  draft.value = next
}

// Setting a mode on a tool a family already covers writes a named entry beside
// the family grant, which is how "all data tools, but never delete" is said.
// The named entry wins: the runtime expands families after them.
function setMode(tool, mode) {
  const next = new Map(draft.value)
  applyMode(next, tool, mode)
  draft.value = next
}

function applyMode(next, tool, mode) {
  // Blocking a tool nothing else grants is saying nothing about it, so the
  // entry goes rather than staying behind as a grant that grants nothing.
  if (mode === 'deny') {
    if (covered(tool)) next.set(tool.name, 'deny')
    else next.delete(tool.name)
    return
  }

  // Choosing the mode the risk would have given anyway clears the override
  // rather than pinning it, so a tool that should simply follow the rule keeps
  // doing so if the rule ever changes.
  const isDefault = mode === defaultModeFor(tool.risk)
  if (isDefault && covered(tool) && !named.value.has(tool.name)) {
    next.delete(tool.name)
    return
  }

  next.set(tool.name, isDefault ? '' : mode)
}

// Open where there is already something to see, so a configured agent shows its
// own grants and a fresh one still opens on something it can act on.
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

  const toolEntries = [...draft.value.entries()].map(([name, permission]) => ({
    ...(byName.get(name) || { name, description: '', allow: [] }),
    name,
    permission,
  }))

  emit('apply', [...families.value, ...toolEntries])
  emit('update:visible', false)
}
</script>

<style scoped>
/* A blocked row shows nothing until it is reached for. Written here rather than
   as utilities so the rule survives whatever else styles the row. */
.tool-mode {
  transition: opacity 0.15s ease;
}

.tool-mode[data-blocked='true']:not(:focus-within) {
  opacity: 0;
}

.group:hover .tool-mode[data-blocked='true'] {
  opacity: 1;
}
</style>
