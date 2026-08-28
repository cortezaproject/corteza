<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('agent.editor.tools.dialog.title')"
    :style="{ width: '72rem', height: '86vh' }"
    :contentStyle="{
      display: 'flex',
      flexDirection: 'column',
      minHeight: 0,
      flex: '1 1 auto',
      paddingLeft: 0,
      paddingRight: 0,
    }"
    :breakpoints="{ '1200px': '94vw' }"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-4 min-h-0 flex-1">
      <IconField class="mx-4">
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="search"
          :placeholder="$t('agent.editor.tools.dialog.search')"
          class="w-full"
          data-testid="tool-dialog-search"
        />
      </IconField>

      <!-- Access is deny-by-default, so an agent granted nothing does nothing.
           The other half is the part that is never otherwise on screen: a tool
           runs as the person invoking the agent, so it is a ceiling, not a
           widening. -->
      <Message v-if="!chosenCount" severity="secondary" :closable="false" class="!my-0 mx-4">
        {{ $t('agent.editor.tools.dialog.inherits') }}
      </Message>

      <div class="overflow-y-auto flex-1 min-h-0">
        <section v-for="d in domains" :key="d.key">
          <div
            class="tool-section-head sticky top-0 z-10 flex items-center gap-4 border-b border-surface px-4 py-2"
          >
            <button
              type="button"
              class="flex items-start gap-2 flex-1 min-w-0 text-left"
              :data-testid="`collapse-${d.key}`"
              @click="toggleCollapsed(d.key)"
            >
              <i
                class="pi pi-chevron-right tool-section-chevron text-xs text-muted-color w-3 mt-1"
                :data-open="isOpen(d.key)"
              />
              <span class="min-w-0">
                <span class="block text-sm font-semibold text-color">{{ d.label }}</span>

                <!-- What the subject is. The tools inside say what they do; a
                     section name alone does not say what it covers, and
                     "People and access" is the one nobody guesses. -->
                <span class="text-xs text-muted-color block">{{ d.description }}</span>
              </span>
            </button>

            <!-- The split, not the total: a section reads as what it lets the
                 agent do, and the total is the sum of what is shown. It sits
                 with the control that sets it rather than against the title,
                 where a name and a run of glyphs ran together. -->
            <span class="flex shrink-0 items-center gap-2.5 text-sm text-muted-color">
              <span
                v-for="c in sectionCounts(d)"
                :key="c.mode"
                v-tooltip.top="modeLabel(c.mode)"
                class="flex items-center gap-1"
                :data-testid="`section-count-${d.key}-${c.mode}`"
              >
                <i :class="[modeIcon(c.mode), modeColour(c.mode)]" />
                <span class="tabular-nums">{{ c.n }}</span>
              </span>
            </span>

            <!-- One statement about the whole section: what every tool in it
                 should do. It reads back as Custom when they disagree, which is
                 the normal state once anything has been set by hand. -->
            <Select
              :model-value="sectionMode(d)"
              :options="modeOptions"
              option-label="label"
              option-value="value"
              size="small"
              class="w-48 shrink-0"
              :disabled="disabled"
              :aria-label="$t('agent.editor.tools.mode.section')"
              :data-testid="`section-mode-${d.key}`"
              @update:model-value="v => setSectionMode(d, v)"
            >
              <!-- Custom carries no glyph: the other three each have one, so
                   the word standing alone is what says the section is not one
                   thing. An icon there read as a menu. -->
              <template #value="{ value }">
                <span class="flex items-center gap-2">
                  <i
                    v-if="modeIcon(value)"
                    :class="[modeIcon(value), modeColour(value)]"
                    class="text-xs"
                  />
                  {{ modeLabel(value) }}
                </span>
              </template>
              <!-- Every entry keeps its colour here. The list is the choices,
                   not the state, and the tick already says which is current.
                   Custom is not among them: it is what the control reads as
                   once its tools disagree, and picking it says nothing. -->
              <template #option="{ option }">
                <span class="flex items-center gap-2">
                  <i :class="[option.icon, option.colour]" class="text-xs" />
                  {{ option.label }}
                </span>
              </template>
            </Select>
          </div>

          <!-- PrimeVue's own collapse, the one Panel and Fieldset use, so a
               section folds at the same rate as the panels behind the dialog.
               Accordion is the component for this shape, but its header is a
               button and ours holds a Select. -->
          <transition name="p-collapsible">
            <div v-if="isOpen(d.key)" class="tool-section-body">
              <div class="pb-3">
                <div v-for="area in d.areas" :key="area.key">
                  <!-- A section holding one subject has already named it. An
                       area is the smaller grouping, and reads as one: at the
                       section header's size it claimed to be a section. -->
                  <div
                    v-if="d.areas.length > 1"
                    class="mb-1 mt-3 px-4 text-xs font-medium uppercase tracking-wide text-muted-color"
                  >
                    {{ area.label }}
                  </div>

                  <div
                    v-for="tool in area.tools"
                    :key="tool.name"
                    class="border-b border-surface last:border-b-0"
                  >
                    <div class="tool-row group flex items-start gap-3 px-4 py-2">
                      <div
                        class="flex-1 min-w-0 flex items-start gap-2"
                        :class="{ 'tool-row-blocked': !toolOn(tool) }"
                      >
                        <!-- Risk is a column, not a prefix. Inline, it pushed
                             every title off the left edge its own description
                             sat on. -->
                        <i
                          v-tooltip.top="riskLabel(tool.risk)"
                          :class="[riskIcon(tool.risk), riskColour(tool.risk)]"
                          class="text-xs shrink-0 w-3 mt-1"
                        />

                        <div class="min-w-0 flex-1">
                          <!-- Where a family covers this tool, the title says
                               so. On the row it wrapped every control in it,
                               and the risk glyph inside answered with the row's
                               tooltip instead of its own. -->
                          <span
                            v-tooltip.top="
                              covered(tool) ? $t('agent.editor.tools.dialog.fromFamily') : ''
                            "
                            class="text-sm text-color block"
                          >
                            {{ tool.title || tool.name }}
                          </span>
                          <!-- Under a search the row shows the sentence the
                               search hit, not the first one: matching runs over
                               the whole description, so a row whose opening
                               sentence lacks the term still earns its place and
                               has to be able to show why. -->
                          <span class="text-xs text-muted-color block line-clamp-2">
                            {{ describe(tool) }}
                          </span>

                          <!-- What the row's settings hold, so a configured tool
                               reads as one without being opened. -->
                          <span
                            v-if="toolOn(tool) && settingsSummary(tool)"
                            class="text-xs text-primary block mt-0.5"
                            :data-testid="`settings-summary-${tool.name}`"
                          >
                            {{ settingsSummary(tool) }}
                          </span>

                          <!-- A tool's own settings: the note the model reads
                               before calling it, and how far it may reach.
                               Narrowing shows only where the policy check can
                               act on it. It opens inside the row's text column
                               so it is the width of what it belongs to, and the
                               row's own controls stay where they were. -->
                          <transition name="p-collapsible">
                            <div
                              v-if="isConfiguring(tool)"
                              class="tool-section-body"
                              :data-testid="`tool-settings-${tool.name}`"
                            >
                              <div
                                class="min-h-0 mt-2 mb-1 p-3 rounded-border border border-surface bg-emphasis flex flex-col gap-3"
                              >
                                <!-- A named exception to the stacked-label rule:
                                     CFormGroup's label is uppercase at the size
                                     of a tool's title, which in a list of ninety
                                     rows makes a field label the loudest thing
                                     on screen. -->
                                <div class="flex flex-col gap-1.5">
                                  <label class="text-xs font-medium text-muted-color">
                                    {{ $t('agent.editor.tools.dialog.settings.note') }}
                                  </label>
                                  <InputText
                                    :model-value="entryOf(tool).description"
                                    :placeholder="
                                      $t('agent.editor.tools.dialog.settings.notePlaceholder')
                                    "
                                    :disabled="disabled"
                                    size="small"
                                    class="w-full"
                                    @update:model-value="v => patchEntry(tool, { description: v })"
                                  />
                                  <small class="text-muted-color">
                                    {{ $t('agent.editor.tools.dialog.settings.noteHelp') }}
                                  </small>
                                </div>

                                <div
                                  v-if="scopesModules(tool.name) && namespaces.length"
                                  class="flex flex-col gap-1.5"
                                >
                                  <label class="text-xs font-medium text-muted-color">
                                    {{ $t('agent.editor.tools.dialog.settings.modules') }}
                                  </label>
                                  <div v-for="ns in namespaces" :key="ns.namespaceID" class="mb-2">
                                    <CInputModule
                                      :model-value="modulesFor(tool, ns.namespaceID)"
                                      :namespace-i-d="ns.namespaceID"
                                      :multiple="true"
                                      :placeholder="
                                        $t('agent.editor.tools.dialog.settings.allModules')
                                      "
                                      :disabled="disabled"
                                      @update:model-value="
                                        v => setModulesFor(tool, ns.namespaceID, v)
                                      "
                                    />
                                  </div>
                                  <small class="text-muted-color">
                                    {{ $t('agent.editor.tools.dialog.settings.modulesHelp') }}
                                  </small>
                                </div>

                                <small
                                  v-else-if="canScope(tool.name) && !namespaces.length"
                                  class="text-muted-color"
                                >
                                  {{ $t('agent.editor.tools.dialog.settings.needsWorksIn') }}
                                </small>

                                <small v-else-if="scopeBlocked(tool.name)" class="text-muted-color">
                                  {{ $t('agent.editor.tools.dialog.settings.namedElsewhere') }}
                                </small>
                              </div>
                            </div>
                          </transition>
                        </div>
                      </div>

                      <!-- Blocked keeps the button's place rather than its
                           function: the permission controls stay in one column
                           however many rows above are off. -->
                      <Button
                        icon="pi pi-cog"
                        :severity="hasSettings(entryOf(tool)) ? 'primary' : 'secondary'"
                        text
                        rounded
                        size="small"
                        class="shrink-0"
                        :class="{ invisible: !toolOn(tool), 'tool-gear-open': isConfiguring(tool) }"
                        :disabled="disabled || !toolOn(tool)"
                        :aria-label="$t('agent.editor.tools.dialog.settings.label')"
                        :data-testid="`configure-${tool.name}`"
                        @click="toggleConfiguring(tool.name)"
                      />

                      <!-- Blocked is the off state, so the row shows nothing until it
                     is reached for. Three segments carry the whole decision:
                     a second control for on and off would only disagree. -->
                      <div class="tool-mode shrink-0" :data-blocked="!toolOn(tool)">
                        <SelectButton
                          :model-value="rowMode(tool)"
                          :options="modeOptions"
                          option-value="value"
                          :allow-empty="false"
                          size="small"
                          :disabled="disabled"
                          @update:model-value="v => setMode(tool, v)"
                        >
                          <!-- Only the chosen segment is coloured. Colouring all three
                         made every row shout in red about the state it was not
                         in, and left the control unable to show its own. -->
                          <template #option="{ option }">
                            <i
                              v-tooltip.top="option.label"
                              :class="[
                                option.icon,
                                rowMode(tool) === option.value ? option.colour : '',
                              ]"
                            />
                          </template>
                        </SelectButton>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </transition>
        </section>

        <div v-if="!domains.length" class="text-sm text-muted-color px-4 py-6 text-center">
          {{
            search.trim()
              ? $t('agent.editor.tools.dialog.noMatches')
              : $t('agent.editor.tools.dialog.empty')
          }}
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex flex-col items-end gap-2">
        <div class="flex gap-2">
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

        <!-- What Apply would grant, counted the way the panel behind this
             dialog counts it, so leaving here never changes the number. -->
        <span class="text-sm text-muted-color" data-testid="tool-dialog-total">
          {{ $t('agent.editor.tools.dialog.chosen', { n: heldCount }) }}
        </span>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { components } from '@planetcrust/human-vue'

import {
  DOMAINS,
  HIDDEN_AREAS,
  MODE_COLOURS,
  MODE_ICONS,
  MODES,
  PLACED_AREAS,
  RISK_ORDER,
  areaOf,
  canScope,
  coveredByFamily,
  defaultModeFor,
  hasSettings,
  modeOf,
  scopeBlocked,
  SUMMARY_MODES,
  scopesModules,
  summaryOf,
  tally,
} from '../toolAccess'

const { CInputModule } = components

const props = defineProps({
  visible: { type: Boolean, default: false },
  tools: { type: Array, default: () => [] },
  // The agent's access.tools, as stored.
  grants: { type: Array, default: () => [] },
  // The agent's own scope. A tool narrows within these and never past them.
  namespaces: { type: Array, default: () => [] },
  // The modules those namespaces hold, so a narrowing reads back as the names
  // it was chosen by rather than as IDs.
  modules: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:visible', 'apply'])

const { t } = useI18n()

const search = ref('')

// Editing happens on copies: closing with Cancel has to leave the agent
// exactly as it was. An entry carries everything the grant holds — its mode,
// the note the model reads, and any narrowing — because the row edits all
// three.
const draft = ref(new Map()) // tool name -> grant entry

// What the permission rule reads: it asks about modes, not whole entries.
const permissions = computed(
  () => new Map([...draft.value].map(([name, entry]) => [name, entry.permission || ''])),
)

// Which rows are open on their own settings.
const configuring = ref(new Set())

function isConfiguring(tool) {
  return configuring.value.has(tool.name) && toolOn(tool)
}

function toggleConfiguring(name) {
  const next = new Set(configuring.value)
  if (next.has(name)) next.delete(name)
  else next.add(name)
  configuring.value = next
}

function entryOf(tool) {
  return (
    draft.value.get(tool.name) || { name: tool.name, permission: '', description: '', allow: [] }
  )
}

function patchEntry(tool, patch) {
  const next = new Map(draft.value)
  next.set(tool.name, { ...entryOf(tool), ...patch })
  draft.value = next
}

// A tool narrows within the agent's own namespaces; it can never reach past
// them. No modules chosen for a namespace means every module in it.
function modulesFor(tool, namespaceID) {
  return entryOf(tool).allow?.find(a => a.namespaceID === namespaceID)?.moduleIDs || []
}

function setModulesFor(tool, namespaceID, moduleIDs) {
  const rest = (entryOf(tool).allow || []).filter(a => a.namespaceID !== namespaceID)
  const allow = moduleIDs?.length ? [...rest, { namespaceID, moduleIDs }] : rest
  patchEntry(tool, { allow })
}

// A grant may name a whole family instead of a tool. This dialog chooses tools
// and no longer writes families, but it still has to know which tools an
// existing family already covers, and hand every family back untouched.
const families = ref([])

// The tools the agent was opened with, so an override added in this session can
// be told from an entry that was already there.
const named = ref(new Set())
const chosenCount = computed(() => draft.value.size + families.value.length)

// How many tools the agent would end up holding. Counted through summaryOf so
// the dialog and the readout behind it can never state different totals; the
// label it takes is only used for the rows, which this ignores.
const heldCount = computed(
  () => summaryOf(props.tools, permissions.value, families.value, key => key).total,
)

// Which domains are folded shut, or null while the default still stands. A
// search folds separately: it opens everything so no match hides in a shut
// section, and what was folded before it is still folded when it clears.
const collapsed = ref(null)
const searchCollapsed = ref(null)

watch(
  () => props.visible,
  open => {
    if (!open) return
    search.value = ''
    collapsed.value = null
    searchCollapsed.value = null
    configuring.value = new Set()
    draft.value = new Map(
      (props.grants || [])
        .filter(g => g.name)
        .map(g => [g.name, { description: '', allow: [], ...g, permission: g.permission || '' }]),
    )
    families.value = (props.grants || []).filter(g => g.group)
    named.value = new Set(draft.value.keys())
  },
  { immediate: true },
)

// Starting or clearing a search hands back an unfolded list; refining one
// mid-search leaves whatever has been folded within it alone.
watch(
  () => Boolean(search.value.trim()),
  () => {
    searchCollapsed.value = null
  },
)

function modeIcon(mode) {
  return MODE_ICONS[mode] || MODE_ICONS.custom
}

function modeColour(mode) {
  return MODE_COLOURS[mode] || ''
}

function modeLabel(mode) {
  return t(`agent.editor.tools.mode.${mode || 'custom'}`)
}

function modeOption(value) {
  return { value, icon: modeIcon(value), colour: modeColour(value), label: modeLabel(value) }
}

// The three a row and a section are both set to. Custom is never among them:
// it is what a section reads as once its tools disagree, and the value slot
// shows it back without it having to be selectable.
const modeOptions = computed(() => MODES.map(modeOption))

// What the three segments read as. There is no fourth "default" segment — the
// default IS one of the three, shown as the chosen one.
function rowMode(tool) {
  return modeOf(tool, permissions.value, families.value)
}

// Whether the agent has the tool at all.
function toolOn(tool) {
  return rowMode(tool) !== 'deny'
}

function covered(tool) {
  return coveredByFamily(tool, families.value)
}

const moduleNames = computed(
  () => new Map((props.modules || []).map(m => [String(m.moduleID), m.name || m.handle])),
)

// The modules a row is narrowed to, by the names they were chosen by. One the
// agent's scope no longer holds falls back to its ID rather than vanishing from
// the summary.
function settingsModules(tool) {
  return (entryOf(tool).allow || [])
    .flatMap(a => a.moduleIDs || [])
    .map(id => moduleNames.value.get(String(id)) || String(id))
}

// What a row's settings amount to, in one line.
function settingsSummary(tool) {
  const parts = []
  const modules = settingsModules(tool)

  if (modules.length) {
    parts.push(
      t('agent.editor.tools.dialog.settings.summaryModules', { modules: modules.join(', ') }),
    )
  }

  if (entryOf(tool).description) parts.push(t('agent.editor.tools.dialog.settings.summaryNote'))

  return parts.join(' · ')
}

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
        description: t(`agent.editor.tools.dialog.domainHelp.${d.key}`),
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

// What the agent has, split the way the readout outside splits it. Blocked is
// not a third thing a section holds: a tool the agent may not use and one it was
// never given come to the same thing, and counting them says nothing about what
// the section lets the agent do.
function sectionCounts(d) {
  return tally(toolsIn(d), permissions.value, families.value).filter(c =>
    SUMMARY_MODES.includes(c.mode),
  )
}

// One click for a whole subject, which is the unit an agent is actually given:
// every record tool, or none of them.
function setSectionMode(d, mode) {
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
  const carry = permission => ({
    ...(next.get(tool.name) || { name: tool.name, description: '', allow: [] }),
    name: tool.name,
    permission,
  })

  // Blocking a tool nothing else grants is saying nothing about it, so the
  // entry goes rather than staying behind as a grant that grants nothing.
  if (mode === 'deny') {
    if (covered(tool)) next.set(tool.name, carry('deny'))
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

  next.set(tool.name, carry(isDefault ? '' : mode))
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
// showing a shut section holding the match, so it starts with all of them open.
function isOpen(key) {
  if (search.value.trim()) return !searchCollapsed.value?.has(key)
  if (collapsed.value) return !collapsed.value.has(key)
  return defaultOpen(key)
}

function toggleCollapsed(key) {
  if (search.value.trim()) {
    searchCollapsed.value = toggled(searchCollapsed.value ?? [], key)
    return
  }

  collapsed.value = toggled(
    collapsed.value ?? domains.value.filter(d => !defaultOpen(d.key)).map(d => d.key),
    key,
  )
}

function toggled(keys, key) {
  const next = new Set(keys)
  if (next.has(key)) next.delete(key)
  else next.add(key)
  return next
}

function byRiskThenName(a, b) {
  const ra = RISK_ORDER[a.risk] ?? RISK_ORDER.read
  const rb = RISK_ORDER[b.risk] ?? RISK_ORDER.read
  return ra - rb || (a.title || a.name).localeCompare(b.title || b.name)
}

function matches(tool, q) {
  return `${tool.name} ${tool.title || ''} ${tool.description || ''}`.toLowerCase().includes(q)
}

// What the row shows of a description: the sentence carrying the search term
// while one is active, and otherwise the first.
function describe(tool) {
  const q = search.value.trim().toLowerCase()
  if (!q) return summarise(tool.description)

  const text = String(tool.description || '').trim()
  const hit = text.split(/(?<=\.)\s+/).find(sentence => sentence.toLowerCase().includes(q))

  return hit || summarise(tool.description)
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
  emit('apply', [...families.value, ...draft.value.values()])
  emit('update:visible', false)
}
</script>

<style scoped>
/* The row the collapse animates from 0fr to 1fr. Its child needs a floor of
   zero to shrink past its own content. */
.tool-section-body {
  display: grid;
  grid-template-rows: 1fr;
}

.tool-section-body > * {
  min-height: 0;
}

/* A blocked row is legible, not erased: it is still a row you can read and turn
   back on. Only what the row says recedes — the control that turns it back on
   is the one thing that must not. */
.tool-row-blocked {
  opacity: 0.7;
}

/* One chevron that turns, at the rate the section folds. Two icons swapped for
   each other snap while everything around them eases. */
.tool-section-chevron {
  transition: transform 0.2s ease-out;
}

.tool-section-chevron[data-open='true'] {
  transform: rotate(90deg);
}

/* A pinned header keeps the subject and its permission control in reach while
   thirty rows go past. The dialog's own background rather than a surface class,
   so rows scrolling under it stay hidden in either theme. */
.tool-section-head {
  background: var(--p-dialog-background, var(--p-content-background, #fff));
}

/* The gear stays lit while the settings it opened are on screen, so the row
   says which button is holding them open. */
.tool-gear-open {
  background: var(--p-content-hover-background);
}

/* The row under the pointer. Ninety rows of the same shape, each with a control
   hard against the right edge, and nothing said which one a click would land
   on. The dialog does not pad its content sideways — each part pads itself — so
   the band is the width of the dialog rather than a box floating inside it. */
.tool-row {
  transition: background-color 0.1s ease;
}

.tool-row:hover {
  background: var(--p-content-hover-background);
}
</style>
