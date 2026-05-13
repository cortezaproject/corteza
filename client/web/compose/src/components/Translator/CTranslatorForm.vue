<!-- client/web/compose/src/components/Translator/CTranslatorForm.vue -->
<template>
  <div class="flex flex-col gap-0 overflow-hidden">
    <!-- Language column manager -->
    <div class="flex items-center gap-2 px-3 py-2 border-b border-surface bg-surface flex-wrap">
      <span class="text-sm font-medium text-muted-color shrink-0">{{ $t('translator.languages') }}</span>
      <div class="flex flex-wrap items-center gap-1 flex-1">
        <span
          v-for="lang in visibleLanguages"
          :key="lang.tag"
          class="inline-flex items-center gap-1 rounded-full bg-primary text-primary-contrast px-2.5 h-7 text-xs"
        >
          {{ lang.localizedName }}
          <button
            v-if="!lang.isDefault"
            class="opacity-70 hover:opacity-100 leading-none focus:outline-none"
            @click="setVisible(lang.tag, false)"
          >
            <i class="pi pi-times text-[9px]" />
          </button>
        </span>
        <Select
          v-if="hiddenLanguages.length"
          v-model="pendingAdd"
          :options="hiddenLanguages"
          option-label="localizedName"
          option-value="tag"
          :placeholder="$t('translator.add-language')"
          size="small"
          class="text-xs h-7"
          @change="onAddLanguage"
        />
      </div>
    </div>

    <!-- Translation table -->
    <div class="overflow-auto flex-1">
      <table class="w-full text-sm border-collapse">
        <thead class="bg-surface sticky top-0 z-10">
          <tr>
            <th class="text-left px-3 py-2 font-medium text-muted-color whitespace-nowrap border-b border-surface">
              {{ $t('translator.key') }}
            </th>
            <th
              v-for="lang in visibleLanguages"
              :key="lang.tag"
              class="text-left px-3 py-2 font-medium border-b border-surface"
            >
              {{ lang.localizedName }}
            </th>
          </tr>
        </thead>
        <tbody>
          <template v-for="section in sections" :key="section.resource">
            <!-- Section header for non-primary resources -->
            <tr class="bg-surface/60" v-if="!section.isPrimary">
              <td
                :colspan="visibleLanguages.length + 1"
                class="px-3 py-1 text-xs font-semibold text-muted-color border-b border-surface"
              >
                {{ section.title || section.resource }}
              </td>
            </tr>
            <!-- Translation rows -->
            <tr
              v-for="key in section.keys"
              :key="`${section.resource}:${key}`"
              :class="{ 'bg-amber-50 dark:bg-amber-950': isRowDirty(section.resource, key) }"
            >
              <td class="px-3 py-2 border-b border-surface text-muted-color font-mono text-xs align-top leading-relaxed whitespace-nowrap">
                {{ prettyKey(key) }}
              </td>
              <td
                v-for="lang in visibleLanguages"
                :key="lang.tag"
                :class="[
                  'border-b border-surface border-l border-surface p-0',
                  isDirty(section.resource, key, lang.tag) ? 'bg-amber-100 dark:bg-amber-900' : '',
                ]"
              >
                <textarea
                  :value="msg(section.resource, key, lang.tag)"
                  :placeholder="$t('translator.missing')"
                  rows="1"
                  class="w-full bg-transparent text-sm px-3 py-2 resize-none focus:outline-none placeholder:text-muted-color/50 min-h-[2.25rem]"
                  style="field-sizing: content"
                  @input="onUpdate(section.resource, key, lang.tag, ($event.target as HTMLTextAreaElement).value)"
                />
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import type { Language } from '@/stores/languages'
import type { ResourceTranslation } from '@/stores/translator'

const LS_KEY = 'compose:resource-translator:languages'

interface InternalTranslation extends ResourceTranslation {
  org: string
  dirty: boolean
}

interface InternalLanguage extends Language {
  isDefault: boolean
  visible: boolean
}

const props = defineProps<{
  languages: Language[]
  primaryResource: string
  translations: ResourceTranslation[]
  titles?: Record<string, string>
  highlightKey?: string
  keyPrettifier?: (key: string) => string
}>()

const emit = defineEmits<{
  change: [changes: ResourceTranslation[]]
}>()

// Restore previously visible languages from localStorage
const savedTags = new Set((localStorage.getItem(LS_KEY) || '').split(',').filter(Boolean))

const intLanguages = reactive<InternalLanguage[]>(
  props.languages.map((lang, i) => ({
    ...lang,
    isDefault: i === 0,
    visible: i === 0 || savedTags.has(lang.tag),
  })),
)

const intTranslations = reactive<InternalTranslation[]>(
  props.translations.map(t => ({ ...t, org: t.message, dirty: false })),
)

// Persist visible language selection
watch(
  intLanguages,
  () => {
    const tags = intLanguages.filter(l => l.visible).map(l => l.tag).join(',')
    localStorage.setItem(LS_KEY, tags)
  },
  { deep: true },
)

const visibleLanguages = computed(() => intLanguages.filter(l => l.visible))
const hiddenLanguages = computed(() => intLanguages.filter(l => !l.visible))

const pendingAdd = ref<string>('')

function setVisible(tag: string, visible: boolean): void {
  const lang = intLanguages.find(l => l.tag === tag)
  if (lang) lang.visible = visible
}

function onAddLanguage(): void {
  if (pendingAdd.value) {
    setVisible(pendingAdd.value, true)
    pendingAdd.value = ''
  }
}

const sections = computed(() => {
  const resources = [...new Set(intTranslations.map(t => t.resource))]
  return resources.map(resource => ({
    resource,
    isPrimary: resource === props.primaryResource,
    title: props.titles?.[resource] || '',
    keys: [...new Set(intTranslations.filter(t => t.resource === resource).map(t => t.key))],
  }))
})

function find(resource: string, key: string, lang: string): InternalTranslation | undefined {
  return intTranslations.find(t => t.resource === resource && t.key === key && t.lang === lang)
}

function isDirty(resource: string, key: string, lang: string): boolean {
  return find(resource, key, lang)?.dirty ?? false
}

function isRowDirty(resource: string, key: string): boolean {
  return intLanguages.some(l => isDirty(resource, key, l.tag))
}

function msg(resource: string, key: string, lang: string): string {
  return find(resource, key, lang)?.message ?? ''
}

function prettyKey(key: string): string {
  if (props.keyPrettifier) return props.keyPrettifier(key)
  return key
    .replace(/([A-Z])/g, ' $1')
    .toLowerCase()
    .replace(/(\d+)/g, '#$1')
    .split('.')
    .map(s => s.charAt(0).toUpperCase() + s.slice(1))
    .join(' › ')
}

function onUpdate(resource: string, key: string, lang: string, message: string): void {
  const existing = find(resource, key, lang)
  if (existing) {
    existing.dirty = existing.org !== message
    existing.message = message
  } else {
    intTranslations.push({ resource, key, lang, message, org: '', dirty: true })
  }

  const dirty = intTranslations
    .filter(t => t.dirty)
    .map(({ resource, key, lang, message }) => ({ resource, key, lang, message }))
  emit('change', dirty)
}
</script>
