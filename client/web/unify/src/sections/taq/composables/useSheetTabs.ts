import { computed, inject, ref, watch, type Ref } from 'vue'

// Fetches the worksheets of the spreadsheet selected in the current step as
// { label, value } items (value = numeric sheet id). Re-fetches on change/refresh.
export function useSheetTabs() {
  const $SystemAPI = inject<any>('$SystemAPI')
  const argumentsRef = inject<Ref<any[]>>('taq-arguments', ref([]))
  const refreshNonce = inject<Ref<number>>('taq-refresh-nonce', ref(0))

  const tabs = ref<Array<{ label: string; value: string }>>([])
  const loading = ref(false)
  const error = ref('')

  function argValue(name: string) {
    const arg = (argumentsRef.value || []).find((a: any) => a.argumentName === name)
    return arg?.value ?? null
  }

  const configID = computed(() => argValue('configurationID'))
  const spreadsheetId = computed(() => argValue('spreadsheetId'))

  async function load() {
    error.value = ''
    if (!configID.value || !spreadsheetId.value) {
      tabs.value = []
      return
    }
    loading.value = true
    try {
      const { data } = await $SystemAPI.api().get(
        `/configured-connections/${configID.value}/sheet-tabs`,
        { params: { spreadsheetId: spreadsheetId.value } },
      )
      tabs.value = Array.isArray(data?.response) ? data.response : []
    } catch (e: any) {
      tabs.value = []
      error.value = e?.response?.data?.error?.message || e?.message || 'Failed to load worksheets'
    } finally {
      loading.value = false
    }
  }

  watch([configID, spreadsheetId, refreshNonce], load, { immediate: true })

  return { tabs, loading, error }
}
