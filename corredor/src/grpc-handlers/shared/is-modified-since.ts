import * as grpc from '@grpc/grpc-js'

/**
 * Compares if-modified-since value from grpc headers (metadata) with the given date
 *
 * It silently ignores invalid dates from metadata
 *
 * @param {Date} lastModified
 * @param {grpc.Metadata} metadata
 * @return {boolean}
 */
export default function IsModifiedSince(lastModified: Date, metadata: grpc.Metadata): boolean {
  const imsMD = metadata.get('if-modified-since')
  if (imsMD.length === 0) {
    return true
  }

  const imd = Date.parse(imsMD[0].toString())
  return isNaN(imd) || lastModified > new Date(imd)
}
