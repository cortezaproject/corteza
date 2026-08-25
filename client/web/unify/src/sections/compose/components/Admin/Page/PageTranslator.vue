<!-- client/web/compose/src/components/Admin/Page/PageTranslator.vue -->
<template>
  <CTranslatorButton
    v-if="isEdit"
    :resource="resource"
    :titles="titles"
    :fetcher="fetcher"
    :updater="updater"
    :highlight="highlight"
    :key-prettifier="keyPrettifier"
    :disabled="disabled"
  />
</template>

<script setup lang="ts">
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import CTranslatorButton from '@/sections/compose/components/Translator/CTranslatorButton.vue'
import { useTranslatorStore, type TranslationTarget } from '@/sections/compose/stores/translator'
import { compose } from '@planetcrust/human-js'
import { applyPageTranslations, pageKeyLabel } from '@/sections/compose/lib/resource-translations'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'

const props = defineProps<{
  page: any
  namespace: any
  /** All page layouts — pass to also include layout keys in the translation set */
  layouts?: any[]
  /** When set, restrict fetcher to only this block's keys */
  block?: any
  highlight?: TranslationTarget
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:page': [page: any]
  'update:layouts': [layouts: any[]]
}>()

const $ComposeAPI = inject('$ComposeAPI') as any
const { t } = useI18n()
const { currentLanguage } = useResourceTranslations()
const translatorStore = useTranslatorStore()

const isEdit = computed(() => props.page?.pageID && props.page.pageID !== '0')

const resource = computed(() => `compose:page/${props.namespace.namespaceID}/${props.page.pageID}`)

const titles = computed(() => {
  const tt: Record<string, string> = {}

  if (props.block) {
    tt[resource.value] = t('translator.resources.page.block.title', {
      title: props.block.title || props.block.blockID,
    })
  } else {
    tt[resource.value] = t('translator.resources.page.title', {
      handle: props.page.title || props.page.handle || props.page.pageID,
    })
    for (const layout of props.layouts || []) {
      const lRes = `compose:page-layout/${layout.namespaceID}/${layout.pageID}/${layout.pageLayoutID}`
      tt[lRes] = t('translator.resources.page.layout.title', {
        handle: layout.handle || layout.meta?.title || layout.pageLayoutID,
      })
    }
  }

  return tt
})

// The set spans the page, its blocks and every layout; each key reads as what
// it labels rather than as its dotted path.
function keyPrettifier(key: string): string {
  return pageKeyLabel(key, t, props.page?.blocks || [])
}

function fetcher() {
  const { namespaceID, pageID } = props.page
  return $ComposeAPI.pageListTranslations({ namespaceID, pageID }).then((set: any[]) => {
    if (props.block) {
      return set.filter((tr: any) => tr.key.startsWith(`pageBlock.${props.block.blockID}.`))
    }
    return set
  })
}

// The page's translations are one set spanning the page, its blocks and its
// layouts, so the wrapper that knows how to fetch and persist them is also the
// one that opens the dialog aimed at a particular row.
function open(highlight?: TranslationTarget): void {
  translatorStore.open({
    resource: resource.value,
    titles: titles.value,
    highlight: highlight ?? props.highlight,
    keyPrettifier,
    fetcher,
    updater,
  })
}

defineExpose({ open })

async function updater(changes: any[]) {
  const { namespaceID, pageID } = props.page
  await $ComposeAPI.pageUpdateTranslations({ namespaceID, pageID, translations: changes })

  const fresh = await $ComposeAPI.pageListTranslations({ namespaceID, pageID })
  const updatedPage = JSON.parse(JSON.stringify(props.page))
  const updatedLayouts = JSON.parse(JSON.stringify(props.layouts || []))
  applyPageTranslations(updatedPage, updatedLayouts, fresh, currentLanguage.value)
  // Reconstructed rather than handed back as the plain clones they were applied
  // to: an editor holds the typed resource, not a copy of its data.
  emit('update:page', new compose.Page(updatedPage))
  emit(
    'update:layouts',
    updatedLayouts.map((l: any) => new compose.PageLayout(l)),
  )
}
</script>
