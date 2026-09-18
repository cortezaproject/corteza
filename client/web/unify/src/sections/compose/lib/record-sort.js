// Sort expression for the record API: the columns the user picked, in order,
// or the block's presort when none are picked.
export function sortExpression(multiSortMeta, presort = '') {
  if (!multiSortMeta?.length) return presort || ''

  return multiSortMeta
    .map(({ field, order }) => `${field} ${order === 1 ? 'ASC' : 'DESC'}`)
    .join(', ')
}
