// client/web/compose/src/composables/useResourceTranslations.ts
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRBACStore } from '@planetcrust/human-vue'

export function useResourceTranslations() {
  const $Settings = inject('$Settings') as any
  const { locale } = useI18n()
  const rbac = useRBACStore()

  const resourceTranslationLanguages = computed<string[]>(() => {
    const ll = $Settings.get('resourceTranslations.languages')
    if (!ll || !Array.isArray(ll) || ll.length === 0) return ['en']
    return ll
  })

  const resourceTranslationsEnabled = computed(() => resourceTranslationLanguages.value.length > 1)

  const canManageResourceTranslations = computed(() =>
    rbac.can('compose/', 'resource-translations.manage'),
  )

  const showTranslatorButton = computed(
    () => resourceTranslationsEnabled.value && canManageResourceTranslations.value,
  )

  const currentLanguage = computed(() => locale.value)

  return {
    showTranslatorButton,
    currentLanguage,
    resourceTranslationsEnabled,
    canManageResourceTranslations,
  }
}
