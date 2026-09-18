// Sort expression for the record API: the columns the user picked, in order,
// or the block's presort when none are picked.
export function sortExpression(multiSortMeta, presort = '') {
  if (!multiSortMeta?.length) return presort || ''

  return multiSortMeta
    .map(({ field, order }) => `${field} ${order === 1 ? 'ASC' : 'DESC'}`)
    .join(', ')
}

// A sort expression (`name, other DESC`) as [{ field, order }], order 1 or -1.
export function parseSortExpression(raw) {
  return String(raw || '')
    .split(',')
    .map(s => s.trim())
    .filter(Boolean)
    .map(part => {
      const [field, dir] = part.split(/\s+/)
      return { field, order: (dir || '').toUpperCase() === 'DESC' ? -1 : 1 }
    })
}
