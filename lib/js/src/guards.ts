export const IsOf = <T>(v: unknown, ...props: (keyof T)[]): v is T => {
  if (!v || typeof v !== 'object') {
    return false
  }

  for (const prop of props) {
    if (!Object.prototype.hasOwnProperty.call(v, prop)) {
      return false
    }
  }

  return true
}

// eslint-disable-next-line valid-typeof
const every = (a: unknown|unknown[], t: string): boolean => Array.isArray(a) && a.every(i => typeof i === t)

export const AreStrings = (a: unknown|unknown[]): a is string[] => every(a, 'string')
export const AreBooleans = (a: unknown|unknown[]): a is boolean[] => every(a, 'boolean')
export const AreNumbers = (a: unknown|unknown[]): a is number[] => every(a, 'number')
export const AreObjects = (a: unknown|unknown[]): a is object[] => every(a, 'object')

export function AreObjectsOf<T> (a: unknown|unknown[], ...props: (keyof T)[]): a is T[] {
  if (!a || !Array.isArray(a)) {
    return false
  }

  if (a.length === 0) {
    return true
  }

  return AreObjects(a) && a.every(i => IsOf<T>(i, ...props))
}
