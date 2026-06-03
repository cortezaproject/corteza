import { useRouter } from 'vue-router'
import { handleInternalAnchorClick, resolveInternalPath } from '../utils/internalNav'

/**
 * Router-aware navigation helpers for components that link to other "apps"
 * which may actually be internal sections of the unified shell. Internal
 * targets navigate client-side (no reload / no splash); everything else falls
 * back to native navigation.
 */
export function useInternalLink() {
  const router = useRouter()

  return {
    resolveInternalPath: (target: string) => resolveInternalPath(router, target),
    onAnchorClick: (event: MouseEvent, target: string) =>
      handleInternalAnchorClick(router, event, target),
  }
}
