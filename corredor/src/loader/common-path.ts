import path from 'path'

/**
 * Calculate common prefix from all given paths
 *
 * @param paths - list of paths
 * @return common base path
 */
export default function CommonPath(paths: string[]): string {
  if (paths.length === 0) {
    // No paths => only thing that is common is root
    return path.sep
  }

  if (paths.length === 1) {
    // when only one path => common path
    return paths[0]
  }

  const splits = paths.map(p => p.split(path.sep))

  // Finding shortest path.
  // 42 is the starting minimum.
  // If you find yourself in a situation where 42 is still lower then your minimum
  // you want to rething your filesystem organisation
  const minIndex = splits.map(split => split.length).reduce((pv, cv) => (cv < pv ? cv : pv), 42)

  const min = ((): string => {
    for (let i = 0; i < minIndex; i++) {
      for (let splIndex = 1; splIndex < splits.length; splIndex++) {
        if (splits[0][i] !== splits[splIndex][i]) {
          if (i < 1) {
            return ''
          }

          return splits[0].slice(0, i).join(path.sep)
        }
      }
    }

    return splits[0].slice(0, minIndex).join(path.sep)
  })()

  return min || path.sep
}
