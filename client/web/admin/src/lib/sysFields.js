export function getSystemFields (r) {
  const fields = [
    'createdAt',
    'updatedAt',
    'deletedAt',
    'archivedAt',
    'suspendedAt',
    'lastUsedAt',
    'completedAt',
    'createdBy',
  ]
  const viableFields = Object.keys(r).filter(f => fields.includes(f) && r[f])
  const systemFields = viableFields.map(f => f)
  return systemFields
}

export const kebabize = (str) => str.replace(/[A-Z]+(?![a-z])|[A-Z]/g, ($, ofs) => (ofs ? '-' : '') + $.toLowerCase())
