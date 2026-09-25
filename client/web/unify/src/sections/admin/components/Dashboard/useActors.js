import { ref } from 'vue'

// Resolves actor IDs to user names on demand. The admin app has no user
// store, so each unknown ID is read once and kept for the page's lifetime.
export function useActors($SystemAPI) {
  const cache = ref(new Map())

  async function resolve(ids) {
    const missing = [...new Set(ids.filter(id => id && id !== '0' && !cache.value.has(id)))]
    if (!missing.length) return

    const users = await Promise.all(
      missing.map(id => $SystemAPI.userRead({ userID: id }).catch(() => null)),
    )
    const next = new Map(cache.value)
    users.forEach((u, i) => next.set(missing[i], u))
    cache.value = next
  }

  function label(id) {
    if (!id || id === '0') return ''
    const u = cache.value.get(id)
    return u ? u.name || u.handle || u.email || id : id
  }

  return { resolve, label }
}
