<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('agent.editor.tools.dialog.title')"
    :style="{ width: '54rem' }"
    :breakpoints="{ '960px': '90vw' }"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-3 min-h-0">
      <div class="flex items-center gap-2">
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
          {{ $t('agent.editor.tools.dialog.chosen', { count: chosen.size }) }}
        </span>
      </div>

      <!-- Nothing chosen is not "no access": the agent inherits what the
           person invoking it can already do. Saying so here is the difference
           between an empty panel that reads as broken and one that reads as a
           default. -->
      <Message v-if="!chosen.size" severity="secondary" :closable="false" class="!my-0">
        {{ $t('agent.editor.tools.dialog.inherits') }}
      </Message>

      <div class="overflow-y-auto max-h-[26rem] flex flex-col gap-4 pr-1">
        <div v-for="section in sections" :key="section.key">
          <div class="flex items-baseline gap-2 mb-1">
            <span class="text-sm font-semibold text-color">{{ section.label }}</span>
            <span class="text-xs text-muted-color">{{ section.group }}</span>
          </div>

          <div
            v-for="tool in section.tools"
            :key="tool.name"
            class="flex items-center gap-3 py-1.5 border-b border-surface last:border-0"
          >
            <Checkbox
              :model-value="chosen.has(tool.name)"
              binary
              :disabled="disabled"
              :inputId="`tool-${tool.name}`"
              @update:model-value="v => toggle(tool, v)"
            />
            <label :for="`tool-${tool.name}`" class="flex-1 min-w-0 cursor-pointer">
              <span class="text-sm text-color block truncate">{{ tool.title || tool.name }}</span>
              <span class="text-xs text-muted-color block truncate">{{ tool.name }}</span>
            </label>

            <Tag :severity="riskSeverity(tool.risk)" rounded>
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
              class="w-40"
              :disabled="disabled || !chosen.has(tool.name)"
              @update:model-value="v => setMode(tool, v)"
            />
          </div>
        </div>

        <div v-if="!sections.length" class="text-sm text-muted-color py-4 text-center">
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
// Editing happens on a copy: closing the dialog with Cancel has to leave the
// agent exactly as it was.
const draft = ref(new Map())

const chosen = computed(() => new Set(draft.value.keys()))

watch(
  () => props.visible,
  open => {
    if (!open) return
    search.value = ''
    draft.value = new Map(
      (props.grants || []).filter(g => g.name).map(g => [g.name, g.permission || '']),
    )
  },
  { immediate: true },
)

const modeOptions = computed(() => [
  { value: '', label: t('agent.editor.tools.dialog.mode.default') },
  { value: 'always', label: t('agent.editor.tools.dialog.mode.always') },
  { value: 'ask', label: t('agent.editor.tools.dialog.mode.ask') },
  { value: 'deny', label: t('agent.editor.tools.dialog.mode.deny') },
])

// Sections are the tool's own areas — the prefix its name already carries —
// rather than the two coarse groups, which put ninety tools under one heading.
const sections = computed(() => {
  const q = search.value.trim().toLowerCase()
  const byArea = new Map()

  for (const tool of props.tools) {
    if (q && !`${tool.name} ${tool.title || ''}`.toLowerCase().includes(q)) continue

    const key = areaOf(tool.name)
    if (!byArea.has(key)) {
      byArea.set(key, {
        key,
        label: t(`agent.editor.tools.dialog.area.${key}`, key),
        group: (tool.groups || []).join(' · '),
        tools: [],
      })
    }
    byArea.get(key).tools.push(tool)
  }

  return [...byArea.values()].sort((a, b) => a.label.localeCompare(b.label))
})

// compose_record_lookup -> compose_record; system_user_create -> system_user.
function areaOf(name) {
  const parts = String(name).split('_')
  return parts.length > 2 ? `${parts[0]}_${parts[1]}` : parts[0]
}

function modeOf(tool) {
  return draft.value.get(tool.name) ?? ''
}

function toggle(tool, on) {
  const next = new Map(draft.value)
  if (on) next.set(tool.name, '')
  else next.delete(tool.name)
  draft.value = next
}

function setMode(tool, mode) {
  if (!draft.value.has(tool.name)) return
  const next = new Map(draft.value)
  next.set(tool.name, mode || '')
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
  const existing = new Map((props.grants || []).filter(g => g.name).map(g => [g.name, g]))
  const groupGrants = (props.grants || []).filter(g => !g.name)

  const tools = [...draft.value.entries()].map(([name, permission]) => ({
    ...(existing.get(name) || { name, description: '', allow: [] }),
    name,
    permission,
  }))

  emit('apply', [...groupGrants, ...tools])
  emit('update:visible', false)
}
</script>
