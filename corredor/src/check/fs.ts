import fs from 'fs'

export function canAccessPath(path: string, mode = fs.constants.R_OK): boolean {
  try {
    fs.accessSync(path, mode)
    return true
  } catch {
    return false
  }
}
