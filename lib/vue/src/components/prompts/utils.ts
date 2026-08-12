import { automation } from '@planetcrust/human-js'

export function pVal<T = unknown>(vars: automation.Vars, k: string, def?: T): T | undefined {
  if (vars && vars[k] && vars[k]['@value'] !== undefined) {
    return vars[k]['@value']
  }

  return def
}

export function pType(vars: automation.Vars, k: string, def?: string): string | undefined {
  if (vars && vars[k] && vars[k]['@type'] !== undefined) {
    return vars[k]['@type']
  }

  return def
}

export interface ButtonStyle {
  severity?: string
  variant?: string
}

export function variantToSeverity(name: string | undefined): ButtonStyle {
  switch (name) {
    case 'secondary':
      return { severity: 'secondary' }
    case 'success':
      return { severity: 'success' }
    case 'warning':
      return { severity: 'warn' }
    case 'danger':
      return { severity: 'danger' }
    case 'info':
      return { severity: 'info' }
    case 'light':
      return { severity: 'secondary', variant: 'text' }
    case 'dark':
      return { severity: 'contrast' }
    case 'primary':
    default:
      return { severity: 'primary' }
  }
}
