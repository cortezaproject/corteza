<template>
  <div class="flex flex-col gap-4">
    <template v-if="totalFields">
      <!-- View toggle: a grouped overview (by sensitivity level) vs. a matrix
           where you classify each field inline (grouped by module). -->
      <div class="flex items-center justify-end gap-2">
        <span class="text-sm text-muted-color">{{ $t('project.dataSensitivity.groupBy') }}</span>
        <div class="inline-flex rounded-md border border-surface overflow-hidden text-sm">
          <button
            type="button"
            class="px-3 py-1.5 transition-colors"
            :class="
              view === 'level'
                ? 'bg-primary text-primary-contrast'
                : 'text-muted-color hover:bg-emphasis'
            "
            @click="view = 'level'"
          >
            {{ $t('project.dataSensitivity.viewByLevel') }}
          </button>
          <button
            type="button"
            class="px-3 py-1.5 border-l border-surface transition-colors"
            :class="
              view === 'module'
                ? 'bg-primary text-primary-contrast'
                : 'text-muted-color hover:bg-emphasis'
            "
            @click="view = 'module'"
          >
            {{ $t('project.dataSensitivity.viewByModule') }}
          </button>
        </div>
      </div>

      <!-- By sensitivity: one collapsible card per (non-empty) level, fields
           grouped under their module. An overview — click a field to edit it. -->
      <div v-if="view === 'level'" class="flex flex-col gap-4">
        <div
          v-for="g in nonEmptyGroups"
          :key="g.id"
          class="border border-surface rounded-border shadow-sm overflow-hidden"
        >
          <button
            type="button"
            class="w-full flex items-center gap-2 p-3 text-left hover:bg-emphasis transition-colors"
            @click="toggle(g.id)"
          >
            <i
              class="pi pi-chevron-down text-xs text-muted-color shrink-0 transition-transform duration-200"
              :class="{ '-rotate-90': isCollapsed(g.id) }"
            />
            <Tag :value="g.label" :severity="g.severity" />
            <span class="text-xs text-muted-color">{{ countLabel(g.count) }}</span>
          </button>

          <!-- Animated collapse via grid-template-rows 0fr <-> 1fr -->
          <div
            class="grid transition-[grid-template-rows] duration-200 ease-in-out"
            :class="isCollapsed(g.id) ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'"
          >
            <div class="overflow-hidden min-h-0">
              <div class="flex flex-col gap-3 p-3 border-t border-surface">
                <div v-for="mod in g.modules" :key="mod.id">
                  <!-- Module label -->
                  <div class="flex items-center gap-2.5 mb-2">
                    <KindIcon kind="module" size="lg" />
                    <span class="font-medium truncate">
                      {{ mod.name }}
                    </span>
                  </div>

                  <!-- Fields in this module at this level, indented under a guide line -->
                  <div class="ml-4 pl-3 border-l border-surface flex flex-col">
                    <div
                      v-for="f in mod.fields"
                      :key="f.fieldId"
                      class="flex items-center gap-3 pl-2 pr-1 min-h-[2.5rem] rounded hover:bg-emphasis transition-colors"
                    >
                      <button
                        type="button"
                        class="min-w-0 flex-1 truncate text-left text-sm cursor-pointer"
                        @click="editField?.(f.moduleId, f.fieldId)"
                      >
                        {{ f.name }}
                      </button>
                      <Select
                        :model-value="f.level"
                        :options="sensitivityOptions"
                        option-label="label"
                        option-value="id"
                        size="small"
                        class="w-44 shrink-0"
                        :placeholder="$t('project.dataSensitivity.unclassified')"
                        :disabled="disabled"
                        @update:model-value="v => setLevel(f, v)"
                      >
                        <template #option="{ option }">
                          <span class="flex items-center gap-2">
                            <span
                              class="w-2.5 h-2.5 rounded-full shrink-0"
                              :style="{ backgroundColor: option.color }"
                            />
                            <span>{{ option.label }}</span>
                          </span>
                        </template>
                        <template #value="{ value, placeholder }">
                          <span class="flex items-center gap-2">
                            <span
                              class="w-2.5 h-2.5 rounded-full shrink-0"
                              :style="{ backgroundColor: sensitivityOption(value).color }"
                            />
                            <span>{{ sensitivityOption(value).label || placeholder }}</span>
                          </span>
                        </template>
                      </Select>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- By module: one collapsible card per module, fields as rows with a
           radio per level — click a radio to (re)classify. -->
      <div v-else class="flex flex-col gap-4">
        <div
          v-for="mod in moduleRows"
          :key="mod.id"
          class="border border-surface rounded-border shadow-sm overflow-hidden"
        >
          <button
            type="button"
            class="w-full flex items-center gap-3 p-3 text-left hover:bg-emphasis transition-colors"
            @click="toggleModule(mod.id)"
          >
            <i
              class="pi pi-chevron-down text-xs text-muted-color shrink-0 transition-transform duration-200"
              :class="{ '-rotate-90': isModuleCollapsed(mod.id) }"
            />
            <KindIcon kind="module" size="lg" />
            <div class="min-w-0 flex-1">
              <div class="font-medium truncate">{{ mod.name }}</div>
              <div class="text-xs text-muted-color">{{ countLabel(mod.fields.length) }}</div>
            </div>
          </button>

          <!-- Animated collapse via grid-template-rows 0fr <-> 1fr -->
          <div
            class="grid transition-[grid-template-rows] duration-200 ease-in-out"
            :class="isModuleCollapsed(mod.id) ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'"
          >
            <div class="overflow-hidden min-h-0">
              <div class="border-t border-surface">
                <!-- Level column headers -->
                <div
                  class="grid items-center bg-emphasis/40 border-b border-surface"
                  :style="matrixCols"
                >
                  <div class="px-3 py-2" />
                  <div
                    v-for="opt in levelChoices"
                    :key="opt.id || 'unclassified'"
                    class="px-1 py-2 text-center"
                  >
                    <span class="text-[11px] font-semibold" :style="{ color: opt.color }">
                      {{ opt.label }}
                    </span>
                  </div>
                </div>

                <div
                  v-for="f in mod.fields"
                  :key="f.fieldId"
                  class="grid items-center border-b border-surface last:border-b-0 hover:bg-emphasis transition-colors"
                  :style="matrixCols"
                >
                  <button
                    type="button"
                    class="px-3 py-3 text-left text-sm truncate cursor-pointer"
                    @click="editField?.(f.moduleId, f.fieldId)"
                  >
                    {{ f.name }}
                  </button>
                  <div
                    v-for="opt in levelChoices"
                    :key="opt.id || 'unclassified'"
                    class="flex justify-center"
                  >
                    <button
                      type="button"
                      class="w-5 h-5 rounded-full border flex items-center justify-center transition-colors disabled:cursor-default"
                      :class="
                        (f.level || null) === opt.id
                          ? ''
                          : 'border-surface enabled:hover:border-primary'
                      "
                      :style="
                        (f.level || null) === opt.id
                          ? { backgroundColor: opt.color, borderColor: opt.color }
                          : null
                      "
                      :disabled="disabled"
                      :aria-label="opt.label"
                      :title="opt.label"
                      @click="setLevel(f, opt.id)"
                    >
                      <span
                        v-if="(f.level || null) === opt.id"
                        class="w-1.5 h-1.5 rounded-full bg-white"
                      />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <div v-else class="text-muted-color text-sm p-4 border rounded-lg bg-emphasis text-center">
      {{ $t('project.dataSensitivity.empty') }}
    </div>
  </div>
</template>

<script setup>
import KindIcon from '@/sections/project/components/KindIcon.vue'
import { SENSITIVITY_OPTIONS, sensitivity } from '@/sections/project/config/sensitivity'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const { t } = useI18n()
const store = useProjectsStore()
const $toast = inject('$toast')
const editField = inject('editField', null)

// Grouped overview ('level') vs. classification matrix ('module').
const view = ref('level')

// Each level's accent colour, pulled from the very PrimeVue palette its Tag
// severity uses, so dots and labels match the level Tags exactly (and track the
// theme). Mapping: public=success→green, internal=info→sky,
// confidential=warn→orange, restricted=danger→red, unclassified=secondary→surface.
const LEVEL_COLOR = {
  unclassified: 'var(--p-surface-400)',
  public: 'var(--p-green-500)',
  internal: 'var(--p-sky-500)',
  confidential: 'var(--p-orange-500)',
  restricted: 'var(--p-red-500)',
}

// Localized classification choices (incl. Unclassified), each with its colour.
const levelChoices = computed(() =>
  SENSITIVITY_OPTIONS.map(o => ({
    id: o.id,
    label: t(o.labelKey),
    color: LEVEL_COLOR[o.id || 'unclassified'],
  })),
)

// Same set for the by-sensitivity row Select.
const sensitivityOptions = levelChoices
// Resolve an option by its level id (null = Unclassified) for the value slot.
const sensitivityOption = id =>
  sensitivityOptions.value.find(o => (o.id || null) === (id || null)) || {}

// Matrix grid: a flexible name column plus one fixed column per level.
const matrixCols = `grid-template-columns: minmax(0, 1fr) repeat(${SENSITIVITY_OPTIONS.length}, 5.25rem)`

// Most sensitive first; unclassified (level === null) is appended last.
const LEVEL_ORDER = ['restricted', 'confidential', 'internal', 'public']

const moduleResources = computed(() =>
  store.resourcesFor(props.project.projectID).filter(m => m.kind === 'module'),
)

const totalFields = computed(() =>
  moduleResources.value.reduce((n, m) => n + (m.fields || []).length, 0),
)

// Every module with its fields (for the matrix view), modules with no fields
// dropped. Each field carries its current level for the inline classifier.
const moduleRows = computed(() =>
  moduleResources.value
    .map(m => ({
      id: m.id,
      name: m.name,
      fields: (m.fields || []).map(f => ({
        fieldId: f.id,
        name: f.name || t('project.dataSensitivity.untitledField'),
        level: f.sensitivity || null,
        moduleId: m.id,
      })),
    }))
    .filter(m => m.fields.length),
)

// Per sensitivity level (most sensitive first, unclassified last), fields are
// sub-grouped by their owning module so each level card reads module-by-module.
const groups = computed(() => {
  const modulesAtLevel = levelId =>
    moduleResources.value
      .map(m => ({
        id: m.id,
        name: m.name,
        fields: (m.fields || [])
          .filter(f => (f.sensitivity || null) === levelId)
          .map(f => ({
            fieldId: f.id,
            name: f.name || t('project.dataSensitivity.untitledField'),
            level: f.sensitivity || null,
            moduleId: m.id,
          })),
      }))
      .filter(m => m.fields.length)

  const levels = [
    ...LEVEL_ORDER.map(id => ({
      id,
      label: t(sensitivity(id).labelKey),
      severity: sensitivity(id).severity,
      levelId: id,
    })),
    {
      id: 'unclassified',
      label: t('project.dataSensitivity.unclassified'),
      severity: 'secondary',
      levelId: null,
    },
  ]

  return levels.map(l => {
    const modules = modulesAtLevel(l.levelId)
    const count = modules.reduce((n, m) => n + m.fields.length, 0)
    return { ...l, modules, count }
  })
})

const nonEmptyGroups = computed(() => groups.value.filter(g => g.count))

// Per-level collapse state (expanded by default); keyed by group id.
const collapsedGroups = ref({})
const isCollapsed = id => !!collapsedGroups.value[id]
const toggle = id => {
  collapsedGroups.value[id] = !collapsedGroups.value[id]
}

// Per-module collapse state for the by-module view; keyed by module id.
const collapsedModules = ref({})
const isModuleCollapsed = id => !!collapsedModules.value[id]
const toggleModule = id => {
  collapsedModules.value[id] = !collapsedModules.value[id]
}

function countLabel(n) {
  return n === 1 ? t('project.field.one') : t('project.field.many', { count: n })
}

async function setLevel(f, level) {
  try {
    await store.updateField(props.project.projectID, f.moduleId, f.fieldId, { sensitivity: level || null })
  } catch (err) {
    $toast.toastErrorHandler(t('project.dataSensitivity.toastUpdateFailed'))(err)
  }
}
</script>
