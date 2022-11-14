import fzy from 'fuzzysort'

// A helper function to remove diacritics from the given string.
// This allows us so search over characters like čšž with csz.
//
// stackoverflow.com/questions/990904/remove-accents-diacritics-in-a-string-in-javascript/37511463#37511463
function toNFD (s: string): string {
  return s.normalize('NFD').replace(/[\u0300-\u036f]/g, '')
}

// Here we try to calculate the quality of the match using relative numbers.
// By default, fuzzy sort uses score that is not capped.
const goodMatchThreshold = 0.35
function fuzzyMatch ({ indexes }: { indexes: Array<number> }): number {
  return indexes.length / (Math.max(...indexes) + 1)
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function assert (match: (v: string, query: string) => boolean, obj: any, query: string, field?: string, ...fields: Array<string>): boolean {
  if (!obj) {
    return false
  }

  query = toNFD(query.trim().toLowerCase())

  if (typeof obj === 'string') {
    return match(toNFD(obj.trim().toLowerCase()), query)
  } else {
    if (!field) {
      throw new Error('field must be provided when searching over objects')
    }
    return fields.concat([field])
      .some(f => match(toNFD(obj[f].trim().toLowerCase()), query))
  }
}

// AssertStrict checks for a substring and doesn't do any additional guessing
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function AssertStrict (obj: any, query: string, field?: string, ...fields: Array<string>): boolean {
  return assert((v: string, query: string) => v.includes(query), obj, query, field, ...fields)
}

// AssertFuzzy uses some heuristic to determine if the query maybe probably perhaps matches
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function AssertFuzzy (obj: any, query: string, field?: string, ...fields: Array<string>): boolean {
  return assert(
    (v: string, query: string) => {
      const rs = fzy.single(query, v)
      if (!rs) {
        return false
      }
      return fuzzyMatch(rs) >= goodMatchThreshold
    },
    obj, query, field, ...fields,
  )
}

// Assert uses both AssertStrict and AssertFuzzy to check if the query matches,
// If the strict assertion passes that is used, else it tries to do a AssertFuzzy just to make sure
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function Assert (obj: any, query: string, field?: string, ...fields: Array<string>): boolean {
  return AssertStrict(obj, query, field, ...fields) || AssertFuzzy(obj, query, field, ...fields)
}
