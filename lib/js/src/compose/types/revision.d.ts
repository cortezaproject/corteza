import { User } from '../../system'
export interface RevisionChange {
  key: string
  old: Array<unknown>
  new: Array<unknown>
}
export interface Revision {
  changeID: string
  timestamp: Date
  resourceID: string
  revision: number
  operation: string
  userID: string
  user: User | null
  changes: Array<RevisionChange>
  comment: string
}
export interface RawRevisionPayload {
  set: Array<{
    changeID: string
    timestamp: string
    resourceID: string
    revision: number
    operation: string
    userID: string
    changes: Array<RevisionChange>
    comment: string
  }>
}
export declare function convertRevisionPayloadToRevision(
  payload: unknown,
  validChangeKeys: string[],
): Array<Revision>
