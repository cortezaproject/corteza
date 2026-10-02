export const userFields = [
  'createdBy',
  'updatedBy',
  'deletedBy',
  'ownedBy',
  'runAs',
]

const fields = [
  'createdAt',
  'createdBy',
  'updatedAt',
  'updatedBy',
  'deletedAt',
  'deletedBy',
  'archivedAt',
  'suspendedAt',
  'lastUsedAt',
  'completedAt',
  'ownedBy',
  'runAs',
]

export function getSystemFields (r) {
  return fields.filter(f => r[f])
}

export const kebabize = (str) => str.replace(/[A-Z]+(?![a-z])|[A-Z]/g, ($, ofs) => (ofs ? '-' : '') + $.toLowerCase())
