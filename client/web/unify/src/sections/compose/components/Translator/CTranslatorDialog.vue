<!-- client/web/compose/src/components/Translator/CTranslatorDialog.vue -->
<template>
  <Dialog
    v-model:visible="translatorStore.visible"
    modal
    :header="title"
    :style="{ width: '80vw', maxWidth: '1200px' }"
    :pt="{
      content: { class: 'p-0 flex flex-col overflow-hidden', style: 'height: 70vh' },
      footer: { class: 'border-t border-surface p-3' },
    }"
    @hide="translatorStore.close()"
  >
    <div v-if="fetchingTranslations" class="flex items-center justify-center h-32">
      <ProgressSpinner style="width: 24px; height: 24px" />
    </div>

    <CTranslatorForm
      v-else-if="languages.length"
      :languages="languages"
      :primary-resource="translatorStore.config?.resource || ''"
      :translations="translations"
      :titles="translatorStore.config?.titles"
      :highlight-key="translatorStore.config?.highlightKey"
      :key-prettifier="translatorStore.config?.keyPrettifier"
      class="flex-1 overflow-hidden flex flex-col"
      @change="pendingChanges = $event"
    />

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          outlined
          size="small"
          @click="translatorStore.close()"
        />
        <Button
          :label="$t('general.label.save')"
          size="small"
          :loading="saving"
          :disabled="pendingChanges.length === 0"
          @click="handleSave"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTranslatorStore, type ResourceTranslation } from '@/sections/compose/stores/translator'
import { useLanguagesStore } from '@/sections/compose/stores/languages'
import CTranslatorForm from './CTranslatorForm.vue'

const translatorStore = useTranslatorStore()
const languagesStore = useLanguagesStore()
const { t } = useI18n()
const $toast = inject('$toast') as any

const translations = ref<ResourceTranslation[]>([])
const pendingChanges = ref<ResourceTranslation[]>([])
const fetchingTranslations = ref(false)
const saving = ref(false)

const languages = computed(() => languagesStore.set)

const title = computed(() => {
  const cfg = translatorStore.config
  if (!cfg) return t('translator.title')
  if (cfg.titles?.[cfg.resource]) return cfg.titles[cfg.resource]
  // Extract last path segment (the ID) as fallback, e.g. "compose:module/nsID/modID" → "modID"
  const parts = cfg.resource.split('/')
  return parts[parts.length - 1] || cfg.resource
})

watch(
  () => translatorStore.visible,
  async (open) => {
    if (!open) {
      translations.value = []
      pendingChanges.value = []
      return
    }

    // Load languages if not yet loaded
    await languagesStore.load()

    // Fetch translations for this resource
    const cfg = translatorStore.config
    if (!cfg) return
    fetchingTranslations.value = true
    try {
      translations.value = await cfg.fetcher()
    } catch {
      $toast.toastDanger(t('notification.translations.loadFailed'))
      translatorStore.close()
    } finally {
      fetchingTranslations.value = false
    }
  },
)

async function handleSave(): Promise<void> {
  const cfg = translatorStore.config
  if (!cfg || pendingChanges.value.length === 0) return
  saving.value = true
  try {
    await cfg.updater(pendingChanges.value)
    $toast.toastSuccess(t('notification.translations.saved'))
    translatorStore.close()
  } catch {
    $toast.toastDanger(t('notification.translations.saveFailed'))
  } finally {
    saving.value = false
  }
}
</script>
