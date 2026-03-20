/**
 * Typed icon definition supporting icon names, URLs, and attachment IDs.
 */
export interface IconDef {
  type: 'name' | 'url' | 'attachment'
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
    ['name', 'url', 'attachment'].includes((value as IconDef).type)
  )
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
