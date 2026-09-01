export const C311_DATA_STATES = [
  'loading',
  'empty',
  'populated',
  'forbidden',
  'not-found',
  'validation-error',
  'retryable-error',
  'terminal-error',
] as const

export type C311DataState = typeof C311_DATA_STATES[number]

export interface C311Session {
  authenticated: boolean
  actor?: {
    actor_id: string
    capabilities?: string[]
    scopes?: string[]
    available_routes?: string[]
  } | null
  expires_at?: string | null
}

export interface C311Provider {
  getSession (): Promise<C311Session>
}

export interface C311RouteMeta {
  c311?: {
    requiresAuth?: boolean
    route?: string
    capabilities?: string[]
    scopes?: string[]
    public?: boolean
  }
}

export interface C311ErrorLike {
  code?: string
  error?: string
  message?: string
  retryable?: boolean
  status?: number
  current_version?: number
  currentVersion?: number
  fieldErrors?: Array<{ field: string; code: string }>
  errors?: Array<{ field: string; code: string }>
}

export function hasC311Capability (session: C311Session | null | undefined, capability: string): boolean {
  return !!session?.actor?.capabilities?.includes(capability)
}

export function hasC311Scope (session: C311Session | null | undefined, scope: string): boolean {
  return !!session?.actor?.scopes?.includes(scope)
}

export function hasC311Route (session: C311Session | null | undefined, route: string): boolean {
  return !!session?.actor?.available_routes?.includes(route)
}

export function canAccessC311Route (session: C311Session | null | undefined, meta: C311RouteMeta['c311'] = {}): boolean {
  if (meta.requiresAuth && !session?.authenticated) return false
  if (meta.route && (!meta.public || meta.requiresAuth) && !hasC311Route(session, meta.route)) return false
  if (meta.capabilities?.some(capability => !hasC311Capability(session, capability))) return false
  if (meta.scopes?.some(scope => !hasC311Scope(session, scope))) return false
  return true
}

export function c311StateForError (error: C311ErrorLike, hasData = false): C311DataState {
  const code = error.code || error.error
  if (error.status === 401 || code === 'UNAUTHENTICATED' || code === 'INVALID_TOKEN') return 'forbidden'
  if (error.status === 403 || code === 'FORBIDDEN' || code === 'INSUFFICIENT_SCOPE') return 'forbidden'
  if (error.status === 404 || code === 'NOT_FOUND') return 'not-found'
  if (code === 'VALIDATION_ERROR' || error.status === 422) return 'validation-error'
  // Only declared 503 responses are stable server retry cases. Undeclared
  // 5xx responses remain terminal transport failures per contract-v1.
  if (error.status === 503 || error.status === 429 || (error.retryable && (error.status === undefined || error.status < 500))) return 'retryable-error'
  if (hasData && code === 'VERSION_CONFLICT') return 'terminal-error'
  return 'terminal-error'
}

export function c311DataState (options: { loading?: boolean; items?: unknown[] | null; error?: C311ErrorLike | null }): C311DataState {
  if (options.loading) return 'loading'
  if (options.error) return c311StateForError(options.error, !!options.items?.length)
  if (options.items && options.items.length === 0) return 'empty'
  return 'populated'
}

export function c311RouteStorageKey (route: string, actorID?: string): string {
  return `c311.route.${actorID || 'anonymous'}.${route}`
}
