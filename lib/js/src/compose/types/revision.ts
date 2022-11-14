import { User } from '../../system'

export interface RevisionChange {
  key: string;
  old: Array<unknown>;
  new: Array<unknown>;
}

export interface Revision {
  changeID: string;
  timestamp: Date;
  resourceID: string;
  revision: number;
  operation: string;
  userID: string;
  user: User | null;
  changes: Array<RevisionChange>;
  comment: string;
}

export interface RawRevisionPayload {
  set: Array<{
    changeID: string;
    timestamp: string;
    resourceID: string;
    revision: number;
    operation: string;
    userID: string;
    changes: Array<RevisionChange>;
    comment: string;
  }>;
}

function isRawRevisionPayload (raw: unknown): raw is RawRevisionPayload {
  if (!raw || typeof raw !== 'object') {
    console.warn('not  an object', raw)
    return false
  }

  if (!Object.getOwnPropertyNames(raw).includes('set')) {
    console.warn('no set prop', raw)
    return false
  }

  if (!Array.isArray((raw as { set: unknown }).set)) {
    console.warn('set prop not array', raw)
    return false
  }

  return true
}

export function convertRevisionPayloadToRevision (payload: unknown, validChangeKeys: string[]): Array<Revision> {
  if (!isRawRevisionPayload(payload)) {
    throw new Error('Invalid revision payload')
  }

  let filterChanges = (cc: Array<RevisionChange>): Array<RevisionChange> => cc

  if (validChangeKeys.length > 0) {
    // filter out changes that don't have valid keys
    filterChanges = (cc: Array<RevisionChange>): Array<RevisionChange> => cc.filter(c => validChangeKeys.includes(c.key))
  }

  return payload.set.map(raw => ({
    changeID: raw.changeID,
    timestamp: new Date(raw.timestamp),
    resourceID: raw.resourceID,
    revision: raw.revision,
    operation: raw.operation,
    userID: raw.userID,
    user: null,
    comment: raw.comment,
    changes: filterChanges(raw.changes),
  }))
}
