import { computed, inject, ref, watch, type Ref } from 'vue'

// Fetches the header row of the sheet selected in the current step, reading the
// sibling spreadsheet/tab/config args. Re-fetches on change or on refresh.
export function useSheetColumns() {
  const $SystemAPI = inject<any>('$SystemAPI')
  const argumentsRef = inject<Ref<any[]>>('taq-arguments', ref([]))
  const refreshNonce = inject<Ref<number>>('taq-refresh-nonce', ref(0))

  const columns = ref<string[]>([])
  const loading = ref(false)
  const error = ref('')

  function argValue(name: string) {
    const arg = (argumentsRef.value || []).find((a: any) => a.argumentName === name)
    return arg?.value ?? null
  }

  const configID = computed(() => argValue('configurationID'))
  const spreadsheetId = computed(() => argValue('spreadsheetId'))
  const tab = computed(() => argValue('sheetName') ?? argValue('range'))

  async function load() {
    error.value = ''
    if (!configID.value || !spreadsheetId.value) {
      columns.value = []
      return
    }
    loading.value = true
    try {
      const { data } = await $SystemAPI.api().get(
        `/configured-connections/${configID.value}/sheet-columns`,
        { params: { spreadsheetId: spreadsheetId.value, tab: tab.value || '' } },
      )
      columns.value = Array.isArray(data?.response) ? data.response : []
    } catch (e: any) {
      columns.value = []
      error.value = e?.response?.data?.error?.message || e?.message || 'Failed to load columns'
    } finally {
      loading.value = false
    }
  }

  watch([configID, spreadsheetId, tab, refreshNonce], load, { immediate: true })

  return { columns, loading, error, spreadsheetId, tab }
}
