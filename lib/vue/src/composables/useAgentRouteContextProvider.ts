import { useRoute } from 'vue-router'

// Default route-only context provider for CAgentSidebar.
// Returns a function that, at agentExec time, snapshots what the router knows
// about the user's current location. Apps that want richer context (e.g.
// resolved namespace/page/record entities) build their own provider on top.
export function useAgentRouteContextProvider(appKind: string) {
  const route = useRoute()

  return function contextProvider(): Record<string, unknown> {
    const ctx: Record<string, unknown> = {
      app: appKind,
      routeName: route.name ? String(route.name) : '',
      routePath: route.path,
    }

    const params = route.params || {}
    const paramKeys = Object.keys(params)
    if (paramKeys.length > 0) {
      ctx.routeParams = paramKeys.reduce<Record<string, string>>((acc, key) => {
        const value = (params as Record<string, unknown>)[key]
        if (value !== undefined && value !== null && value !== '') {
          acc[key] = String(value)
        }
        return acc
      }, {})
    }

    const query = route.query || {}
    const queryKeys = Object.keys(query)
    if (queryKeys.length > 0) {
      ctx.routeQuery = queryKeys.reduce<Record<string, string>>((acc, key) => {
        const value = (query as Record<string, unknown>)[key]
        if (Array.isArray(value)) {
          acc[key] = value.map(v => String(v)).join(',')
        } else if (value !== undefined && value !== null && value !== '') {
          acc[key] = String(value)
        }
        return acc
      }, {})
    }

    return ctx
  }
}
