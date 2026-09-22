/**
 * Filters and reduces list of script and returns yungest one
 *
 * @param {Array<{ updatedAt?: Date }>} scripts
 * @param {Date}                        fallback
 * @return {Date} latest updatedAt from the array of scripts
 */
export default function GetLastUpdated(
  scripts: Array<{ updatedAt?: Date }>,
  fallback = new Date('0000-01-01'),
): Date {
  return scripts
    .map(({ updatedAt }) => updatedAt)
    .filter((updatedAt): updatedAt is Date => !!updatedAt)
    .reduce((last: Date, updatedAt: Date) => {
      return last < updatedAt ? updatedAt : last
    }, fallback)
}
