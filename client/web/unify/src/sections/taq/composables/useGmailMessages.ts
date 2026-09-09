import { computed, inject, ref, watch, type Ref } from 'vue'

// Fetches recent Gmail messages for the configured connection selected in the
// current step as { label, value } items (value = message id).
export function useGmailMessages() {
  const $SystemAPI = inject<any>('$SystemAPI')
  const argumentsRef = inject<Ref<any[]>>('taq-arguments', ref([]))
  const refreshNonce = inject<Ref<number>>('taq-refresh-nonce', ref(0))

  const messages = ref<Array<{ label: string; value: string }>>([])
  const loading = ref(false)
  const error = ref('')

  function argValue(name: string) {
    const arg = (argumentsRef.value || []).find((a: any) => a.argumentName === name)
    return arg?.value ?? null
  }

  const configID = computed(() => argValue('configurationID'))

  async function load() {
    error.value = ''
    if (!configID.value) {
      messages.value = []
      return
    }
    loading.value = true
    try {
      const { data } = await $SystemAPI
        .api()
        .get(`/configured-connections/${configID.value}/gmail-messages`)
      messages.value = Array.isArray(data?.response) ? data.response : []
    } catch (e: any) {
      messages.value = []
      error.value = e?.response?.data?.error?.message || e?.message || 'Failed to load messages'
    } finally {
      loading.value = false
    }
  }

  watch([configID, refreshNonce], load, { immediate: true })

  return { messages, loading, error }
}
