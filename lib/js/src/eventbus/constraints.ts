import minimatch from 'minimatch'
import { IsOf } from '../guards'

interface Constraint {
  name?:
      string;
  op?:
      string;
  value:
      string[];
}

export interface ConstraintMatcher {
  Name (): string|undefined;
  Values (): string[];
  Match(value: string): boolean;
}

export class Equal {
  readonly name?: string
  readonly values: string[]
  protected not: boolean

  constructor (name: string|undefined, vv: string[], not = false) {
    this.name = name
    this.values = vv
    this.not = not
  }

  Name (): string|undefined {
    return this.name
  }

  Values (): string[] {
    return this.values
  }

  Match (value: string): boolean {
    for (const v of this.values) {
      if (value === v) {
        return !this.not
      }
    }

    return this.not
  }
}

/**
 * Handle glob-like pattern matching
 *
 * See: https://github.com/isaacs/minimatch
 */
export class Like extends Equal {
  constructor (name: string|undefined, vv: string[], not = false) {
    super(
      name,
      vv.map(v => v.replace('%', '*').replace('_', '?')),
      not,
    )
  }

  Match (value: string): boolean {
    for (const v of this.values) {
      if (minimatch(value, v)) {
        return !this.not
      }
    }

    return this.not
  }
}

/**
 * Regex matcher
 */
export class Match extends Equal {
  protected re: RegExp[]
  constructor (name: string|undefined, vv: string[], not = false) {
    super(name, vv, not)
    this.re = vv.map(v => new RegExp(v))
  }

  Match (value: string): boolean {
    for (const re of this.re) {
      if (re.test(value)) {
        return !this.not
      }
    }

    return this.not
  }
}

export function ConstraintMaker (c: Constraint|unknown): ConstraintMatcher {
  if (!IsOf<Constraint>(c, 'value')) {
    throw new Error('invalid constraint input')
  }

  const { name = '', op = '', value } = c

  switch (op.toLowerCase()) {
    case '':
    case 'eq':
    case '=':
    case '==':
    case '===':
      return new Equal(name, value)
    case 'not eq':
    case 'ne':
    case '!=':
    case '!==':
      return new Equal(name, value, true)
    case 'like':
      return new Like(name, value)
    case 'not like':
      return new Like(name, value, true)
    case '~':
      return new Match(name, value)
    case '!~':
      return new Match(name, value, true)
    default:
      throw new Error('unsupported constraint operator')
  }
}
