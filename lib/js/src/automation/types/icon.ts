/**
 * Typed icon definition supporting icon names, URLs, attachment IDs, and brand
 * slugs. A 'brand' value is a Simple Icons slug (e.g. "google-sheets") resolved
 * to a logo URL by brandIconUrl.
 */
export interface IconDef {
  type: 'name' | 'url' | 'attachment' | 'brand'
  value: string
}

/**
 * Type guard for IconDef
 */
export function isIconDef(value: unknown): value is IconDef {
  return (
    typeof value === 'object' &&
    value !== null &&
    'type' in value &&
    'value' in value &&
    typeof (value as IconDef).value === 'string' &&
    ['name', 'url', 'attachment', 'brand'].includes((value as IconDef).type)
  )
}

/**
 * Resolve a brand slug to its Simple Icons logo URL.
 */
export function brandIconUrl(slug: string): string {
  return `https://cdn.simpleicons.org/${(slug || '').trim().toLowerCase()}`
}

/**
 * Normalize a plain icon string or IconDef into a consistent IconDef.
 * Plain strings are treated as PrimeIcon names.
 */
export function normalizeIcon(icon: string | IconDef | undefined): IconDef | undefined {
  if (!icon) return undefined
  if (isIconDef(icon)) return icon
  if (typeof icon === 'string') {
    return { type: 'name', value: icon }
  }
  return undefined
}

/**
 * Default icons used across the TAQ builder
 */
export const DEFAULT_ICONS = {
  TRIGGER: { type: 'name', value: 'bolt' } as IconDef,
  ACTION: { type: 'name', value: 'cog' } as IconDef,
  BRANCH: { type: 'name', value: 'sitemap' } as IconDef,
  ITERATOR: { type: 'name', value: 'refresh' } as IconDef,
  END: { type: 'name', value: 'stop-circle' } as IconDef,
}
