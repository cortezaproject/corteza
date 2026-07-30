<template>
  <div v-if="system" class="flex flex-col gap-6 max-w-4xl">
    <!-- Header + back, mirroring FriaScenarioEditor's action bar. -->
    <div class="flex items-start gap-3">
      <Button
        icon="pi pi-arrow-left"
        severity="secondary"
        text
        rounded
        :aria-label="$t('general.label.back')"
        @click="emit('closed')"
      />
      <div class="flex-1 min-w-0">
        <h2 class="text-lg font-medium truncate">
          {{ system.name || $t('fria.aiSystems.untitled') }}
        </h2>
        <p class="text-sm text-muted-color">{{ $t('fria.aiSystems.editorHint') }}</p>
      </div>
      <Button
        v-if="!disabled"
        icon="pi pi-trash"
        severity="danger"
        text
        rounded
        :aria-label="$t('general.label.delete')"
        @click="remove"
      />
    </div>

    <!-- 1. Identity + intended purpose -->
    <section class="rounded-xl border border-surface bg-surface p-4 flex flex-col gap-4">
      <h3 class="font-medium">{{ $t('fria.aiSystems.sections.identity') }}</h3>

      <div class="flex flex-col gap-1">
        <label class="text-sm text-muted-color" for="ais-name">
          {{ $t('fria.aiSystems.fields.name') }}
        </label>
        <InputText id="ais-name" v-model="draft.name" :disabled="disabled" fluid />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm text-muted-color" for="ais-purpose">
          {{ $t('fria.aiSystems.fields.intendedPurpose') }}
        </label>
        <!-- Intended purpose is not decoration: Art. 3(12) makes it the thing
             a system's whole classification turns on. -->
        <Textarea
          id="ais-purpose"
          v-model="draft.intendedPurpose"
          :disabled="disabled"
          rows="3"
          auto-resize
          fluid
          :placeholder="$t('fria.aiSystems.fields.intendedPurposePlaceholder')"
        />
      </div>
    </section>

    <!-- 2. Risk classification -->
    <section class="rounded-xl border border-surface bg-surface p-4 flex flex-col gap-4">
      <div>
        <h3 class="font-medium">{{ $t('fria.aiSystems.sections.classification') }}</h3>
        <p class="text-sm text-muted-color mt-1">
          {{ $t('fria.aiSystems.sections.classificationHint') }}
        </p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <button
          v-for="rc in AI_ACT_RISK_CLASSES"
          :key="rc.key"
          type="button"
          :disabled="disabled"
          class="text-left rounded-lg border p-3 transition-colors disabled:opacity-60"
          :class="
            draft.riskClass === rc.key
              ? RISK_CLASS_CARD_CLASSES[rc.key]
              : 'border-surface hover:bg-emphasis'
          "
          @click="draft.riskClass = rc.key"
        >
          <div class="flex items-center gap-2">
            <span class="font-medium text-sm">{{ $t(rc.labelKey) }}</span>
            <span
              v-if="rc.friaRequired"
              class="text-xs px-1.5 py-0.5 rounded bg-primary/10 text-primary"
            >
              {{ $t('fria.aiSystems.friaRequired') }}
            </span>
          </div>
          <p class="text-xs text-muted-color mt-1">{{ $t(rc.descriptionKey) }}</p>
        </button>
      </div>

      <!-- Annex III only becomes a question once the system is high-risk;
           asking it otherwise implies a classification nobody made. -->
      <div v-if="draft.riskClass === 'high'" class="flex flex-col gap-1">
        <label class="text-sm text-muted-color" for="ais-annex">
          {{ $t('fria.aiSystems.fields.annexIIIPoint') }}
        </label>
        <Select
          id="ais-annex"
          v-model="draft.annexIIIPoint"
          :options="annexOptions"
          option-label="label"
          option-value="value"
          :disabled="disabled"
          show-clear
          fluid
        />
        <!-- Honest about a real gap rather than pretending it saved. -->
        <small class="text-muted-color">
          {{ $t('fria.aiSystems.fields.annexIIINotPersisted') }}
        </small>
      </div>
    </section>

    <!-- 3. Membership -->
    <section class="rounded-xl border border-surface bg-surface p-4 flex flex-col gap-4">
      <div>
        <h3 class="font-medium">{{ $t('fria.aiSystems.sections.members') }}</h3>
        <p class="text-sm text-muted-color mt-1">
          {{ $t('fria.aiSystems.sections.membersHint') }}
        </p>
      </div>

      <div v-for="kind in MEMBER_KINDS" :key="kind" class="flex flex-col gap-1.5">
        <div class="text-xs uppercase tracking-wide text-muted-color">
          {{ $t(`project.kinds.${kind}.plural`) }}
        </div>
        <div v-if="candidatesFor(kind).length" class="flex flex-wrap gap-2">
          <button
            v-for="c in candidatesFor(kind)"
            :key="c.ref"
            type="button"
            :disabled="disabled || busyRef === c.ref"
            class="text-xs px-2 py-1 rounded-full border transition-colors disabled:opacity-50"
            :class="
              isMember(c.ref)
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-surface hover:bg-emphasis'
            "
            @click="toggle(c.ref)"
          >
            <i :class="isMember(c.ref) ? 'pi pi-check' : 'pi pi-plus'" class="mr-1 text-[10px]" />
            {{ c.name }}
          </button>
        </div>
        <div v-else class="text-xs text-muted-color italic">
          {{ $t('fria.aiSystems.noneOfKind') }}
        </div>
      </div>

      <!-- Refs whose resource no longer exists. Kept, never swept: the
           assessment claimed to cover them, so their disappearance is
           material and a reassessment trigger. -->
      <div v-if="orphanRefs.length" class="rounded-lg border border-surface p-3 bg-emphasis">
        <div class="text-xs uppercase tracking-wide text-muted-color mb-2">
          {{ $t('fria.aiSystems.removedResources') }}
        </div>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="ref in orphanRefs"
            :key="ref"
            class="text-xs px-2 py-1 rounded-full border border-surface text-muted-color line-through"
          >
            {{ ref }}
          </span>
        </div>
        <p class="text-xs text-muted-color mt-2">
          {{ $t('fria.aiSystems.removedResourcesHint') }}
        </p>
      </div>
    </section>

    <!-- 4. Human oversight (Art. 14) -->
    <section class="rounded-xl border border-surface bg-surface p-4 flex flex-col gap-3">
      <div>
        <h3 class="font-medium">{{ $t('fria.aiSystems.sections.oversight') }}</h3>
        <p class="text-sm text-muted-color mt-1">
          {{ $t('fria.aiSystems.sections.oversightHint') }}
        </p>
      </div>
      <div v-if="candidatesFor('role').length" class="flex flex-wrap gap-2">
        <button
          v-for="c in candidatesFor('role')"
          :key="c.ref"
          type="button"
          :disabled="disabled || busyRef === c.ref"
          class="text-xs px-2 py-1 rounded-full border transition-colors disabled:opacity-50"
          :class="
            isMember(c.ref)
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-surface hover:bg-emphasis'
          "
          @click="toggle(c.ref)"
        >
          <i :class="isMember(c.ref) ? 'pi pi-check' : 'pi pi-plus'" class="mr-1 text-[10px]" />
          {{ c.name }}
        </button>
      </div>
      <div v-else class="text-xs text-muted-color italic">
        {{ $t('fria.aiSystems.noneOfKind') }}
      </div>
    </section>

    <div v-if="!disabled" class="flex justify-end gap-2">
      <Button
        :label="$t('general.label.save')"
        icon="pi pi-check"
        :loading="saving"
        @click="save"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import { useI18n } from 'vue-i18n'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { AI_ACT_RISK_CLASSES, ANNEX_III_POINTS } from '@/sections/project/config/aiSystemClasses'
import {
  MEMBER_KINDS,
  buildResourceRef,
  parseResourceRef,
} from '@/sections/project/config/resourceRefs'

const props = defineProps({
  project: { type: Object, required: true },
  aiSystemId: { type: String, required: true },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['closed'])

const { t } = useI18n()
const store = useProjectsStore()
const projectId = computed(() => props.project?.projectID)

const system = computed(() => store.aiSystem(projectId.value, props.aiSystemId))

// Same worst→least family as RiskClassBadge, for the selected card's fill.
const RISK_CLASS_CARD_CLASSES = {
  prohibited: 'border-red-500 bg-red-50 dark:bg-red-500/10',
  high: 'border-orange-500 bg-orange-50 dark:bg-orange-500/10',
  limited: 'border-amber-500 bg-amber-50 dark:bg-amber-500/10',
  minimal: 'border-yellow-500 bg-yellow-50 dark:bg-yellow-500/10',
}

// Local draft for the free-text + classification fields, committed on Save —
// same explicit-save contract the FRIA scenario editor uses. Membership is
// deliberately NOT part of it: entries are their own rows and each toggle is
// its own request, so there is no half-saved membership to reconcile.
const draft = reactive({ name: '', intendedPurpose: '', riskClass: null, annexIIIPoint: null })
const saving = ref(false)
const busyRef = ref(null)

watch(
  system,
  s => {
    if (!s) return
    draft.name = s.name
    draft.intendedPurpose = s.intendedPurpose
    draft.riskClass = s.riskClass
  },
  { immediate: true },
)

onMounted(() => {
  // The membership lists are drawn from the per-kind caches, which the Build
  // steps populate lazily — a Govern step may be the first thing opened.
  store.loadResources?.(projectId.value)
  store.loadAgents?.(projectId.value)
  store.loadChatbots?.(projectId.value)
  store.loadAutomations?.(projectId.value)
  store.loadConnections?.(projectId.value)
  store.loadRoles?.(projectId.value)
  store.loadPages?.(projectId.value)
})

const annexOptions = computed(() =>
  ANNEX_III_POINTS.map(p => ({ value: p.key, label: `${p.point} · ${t(p.labelKey)}` })),
)

const memberRefs = computed(() => system.value?.resourceRefs || [])

function isMember(ref) {
  return memberRefs.value.includes(ref)
}

// Candidate resources per kind, as {ref, name}. Reads the store's existing
// per-kind caches rather than the graph payload, so the checklist works on the
// Govern tab where no graph is mounted.
const KIND_SOURCES = {
  module: () => store.resourcesFor(projectId.value),
  page: () => store.pagesFor(projectId.value),
  automation: () => store.automationsFor(projectId.value),
  agent: () => store.agentsFor(projectId.value),
  chatbot: () => store.chatbotsFor(projectId.value),
  connection: () => store.connectionsFor(projectId.value),
  role: () => store.rolesFor(projectId.value),
}

function candidatesFor(kind) {
  const rows = KIND_SOURCES[kind]?.() || []
  return rows
    .map(r => ({
      ref: buildResourceRef(kind, r.id ?? r.resourceID),
      name: r.name || r.handle || r.title || '',
    }))
    .filter(c => c.ref)
}

// Stored refs that match no live resource. Rendered as tombstones, never
// silently dropped.
const orphanRefs = computed(() => {
  const known = new Set()
  for (const kind of Object.keys(KIND_SOURCES)) {
    for (const c of candidatesFor(kind)) known.add(c.ref)
  }
  return memberRefs.value.filter(r => !known.has(r))
})

async function toggle(ref) {
  if (props.disabled || busyRef.value) return
  busyRef.value = ref
  try {
    if (isMember(ref)) {
      await store.removeAiSystemResource(projectId.value, props.aiSystemId, ref)
    } else {
      await store.addAiSystemResource(projectId.value, props.aiSystemId, ref)
    }
  } finally {
    busyRef.value = null
  }
}

async function save() {
  saving.value = true
  try {
    await store.updateAiSystem(projectId.value, props.aiSystemId, {
      name: draft.name,
      intendedPurpose: draft.intendedPurpose,
      riskClass: draft.riskClass,
    })
  } finally {
    saving.value = false
  }
}

async function remove() {
  await store.removeAiSystem(projectId.value, props.aiSystemId)
  emit('closed')
}

// Referenced by the template's tombstone list only; kept exported-by-use so
// the parse helper's intent is visible where refs are rendered raw.
void parseResourceRef
</script>
