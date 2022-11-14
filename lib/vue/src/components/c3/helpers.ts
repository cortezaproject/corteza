export interface Catalogue {
  [_: string]: CatalogueItem;
}

export interface CatalogueItem {
  group?: Array<string>;
  name?: string;
}

const del = '|||'

/**
 * Extract subgroups from the given groups
 * @param groups
 * @param path Array of groups that represent the path
 */
export function extractGroups (groups: Array<Array<string>>, ...path: Array<string>): Array<string> {
  const jpath = path.join(del)
  return [...(new Set(
    groups
      .filter(g => g.slice(0, path.length).join(del) === jpath)
      .map(g => g.slice(path.length, path.length + 1)[0])
      .filter(g => !!g),
  )).values()]
}

/**
 * Extract subgroups from the catalogue that match the given path
 *
 * @param cat Catalogue
 * @param path Array of groups that represent the path
 * @constructor
 */
export function ExtractSubgroups (cat: Catalogue, ...path: Array<string>): Array<string> {
  const groups = Object.values(cat)
    .map(({ group }) => group)
    .filter(g => g !== undefined) as Array<Array<string>>

  return extractGroups(groups, ...path)
}

/**
 * Returns all components that match the given path
 *
 * @param cat Catalogue
 * @param path Array of groups that represent the path

 */
export function ExtractComponents (cat: Catalogue, ...path: Array<string>): Array<CatalogueItem> {
  const jpath = path.join(del)

  return Object.values(cat)
    .filter(({ group }) => (group || []).join(del) === jpath)
}
