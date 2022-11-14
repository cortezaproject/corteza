import { capitalize } from 'lodash'
export interface Typed {
  '@type': string;
  '@value': any;
}

export type Vars = { [_: string]: Typed }

export function IsTyped (a: unknown): a is Typed {
  return typeof a &&
    a === 'object' &&
    Object.prototype.hasOwnProperty.call(a, '@type') &&
    Object.prototype.hasOwnProperty.call(a, '@value')
}

function unwrap (v: unknown): any {
  return IsTyped(v) ? v['@value'] : v
}

function cast (v: unknown): Typed {
  return { '@value': unwrap(v), '@type': guessType(v) }
}

function guessType (v: any): string {
  switch (typeof v) {
    case 'boolean':
      return 'Boolean'
    case 'string':
      return 'String'
    case 'number':
      return Number(v) === v && v % 1 === 0 ? 'Float' : 'Integer'
    case 'object':
      if (v.resourceType) {
        // converts foo:bar into FooBar
        return v.resourceType.split(':').map(capitalize).join('')
      }
      return 'Any'
    default:
      return 'Any'
  }
}

/**
 *
 * @param any
 * @constructor
 */
export function Encode (input: {[_: string]: any}): Vars {
  const output: Vars = {}

  for (const key in input) {
    output[key] = IsTyped(input[key]) ? input[key] : cast(input[key])
  }

  return output
}
