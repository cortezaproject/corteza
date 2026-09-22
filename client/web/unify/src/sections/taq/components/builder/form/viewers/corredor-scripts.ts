// Corredor script names resolved to the label the script declares.
//
// The list request is shared by every viewer instance: a preview carrying
// several Corredor steps, and every preview after it, read the same answer.

interface CorredorScript {
  name: string
  label?: string
}

// Server-side scripts with a manual trigger on the system resource — the set
// lib/vue `CInputCorredorScript` offers for the same argument.
const SCRIPT_FILTER = {
  eventTypes: ['onManual'],
  resourceTypes: ['system'],
  excludeInvalid: true,
  excludeClientScripts: true,
}

let pending: Promise<CorredorScript[]> | null = null

/** The deployed manual system scripts, fetched at most once. */
export function fetchCorredorScripts(api: any): Promise<CorredorScript[]> {
  if (!api) return Promise.resolve([])

  if (!pending) {
    pending = api.automationList(SCRIPT_FILTER).then(
      (result: any) => (Array.isArray(result) ? result : result?.set || []),
      () => {
        // Corredor unreachable or the request refused; the next caller asks again
        pending = null
        return []
      },
    )
  }

  return pending
}

/**
 * The label a script carries, or the name itself when the list has no such
 * script — one deployed later, or hidden from the caller, still reads as itself.
 */
export async function resolveCorredorScriptLabel(api: any, name: string): Promise<string> {
  if (!name) return ''

  const scripts = await fetchCorredorScripts(api)
  return scripts.find(script => script?.name === name)?.label || name
}

/** Drops the shared list, so a caller starts from no answer. */
export function resetCorredorScriptCache(): void {
  pending = null
}
