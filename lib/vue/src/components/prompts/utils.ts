import { automation } from '@cortezaproject/corteza-js'

export function pVal<T = unknown> (vars: automation.Vars, k: string, def?: T): T | undefined {
  if (vars && vars[k] && vars[k]['@value'] !== undefined) {
    return vars[k]['@value']
  }

  return def
}

export function pType (vars: automation.Vars, k: string, def?: string): string | undefined {
  if (vars && vars[k] && vars[k]['@type'] !== undefined) {
    return vars[k]['@type']
  }

  return def
}
