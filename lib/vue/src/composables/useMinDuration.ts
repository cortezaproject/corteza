/**
 * Ensures an async operation takes at least `minMs` milliseconds.
 *
 * Uses Promise.all so the delay runs in parallel with the actual work:
 * - If the work finishes in 50ms and minMs is 300ms, waits the remaining 250ms.
 * - If the work takes 500ms, no extra delay is added.
 *
 * This prevents loading spinners from flashing too quickly.
 *
 * @example
 * // As a standalone utility
 * const data = await withMinDuration(api.fetchItems(), 300)
 *
 * @example
 * // As a composable with reactive loading state
 * const { loading, run } = useMinDuration(300)
 * const data = await run(() => api.fetchItems())
 */
import { ref } from 'vue'

/**
 * Run a promise with a minimum duration guarantee.
 */
export async function withMinDuration<T>(promise: Promise<T>, minMs = 300): Promise<T> {
  const [result] = await Promise.all([promise, new Promise(resolve => setTimeout(resolve, minMs))])
  return result as T
}

/**
 * Composable that provides reactive loading state with minimum duration.
 */
export function useMinDuration(minMs = 300) {
  const loading = ref(false)

  async function run<T>(fn: () => Promise<T>): Promise<T> {
    loading.value = true
    try {
      return await withMinDuration(fn(), minMs)
    } finally {
      loading.value = false
    }
  }

  return { loading, run }
}
