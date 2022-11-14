import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'
import { Compose as ComposeAPI, System as SystemAPI } from '../../../api-clients'
import { Record, Module } from '../../'
import { User } from '../../../system'
import { convertRevisionPayloadToRevision, RawRevisionPayload, Revision } from '../revision'

const kind = 'RecordRevisions'
interface Options {
  // do we preload changes or not
  preload: boolean;

  // what fields do we want to display
  // empty array means all fields
  displayedFields: string[];

  // referenced fields (records, users) we want to expand
  expRefFields: string[];
}

const defaults: Readonly<Options> = Object.freeze({
  preload: false,
  displayedFields: [],
  expRefFields: [],
})

export class PageBlockRecordRevisions extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    Apply(this.options, o, Boolean, 'preload')

    // set new values to displayed fields
    if (Array.isArray(o?.displayedFields)) {
      this.options.displayedFields = o.displayedFields.map(String)
    }

    // set new values to expanded reference fields
    if (Array.isArray(o?.expRefFields)) {
      this.options.expRefFields = o.expRefFields.map(String)
    }
  }

  /**
   * fetch is a utility method on record revision page block
   * that fetches revisions for a record and converts them to RevisionPayload class
   *
   * this function also strips out all fields that should not be dispalyed
   * (as per displayedFields option)
   *
   * @param api Compose API to be used
   * @param record Record to fetch revisions for
   */
  async fetch (api: ComposeAPI, record: Record): Promise<Array<Revision>> {
    const { namespaceID, moduleID, recordID } = record

    return api
      .recordRevisions({ namespaceID, moduleID, recordID })
      .then(payload => convertRevisionPayloadToRevision(
        (payload as unknown) as RawRevisionPayload,
        this.options.displayedFields,
      ))
  }

  /**
   * iterates over all expandable reference fields and fetches values and updates original value
   *
   * @todo ref-record resolving
   *
   * @param $SystemAPI
   * @param $ComposeAPI
   * @param m
   * @param rr
   */
  async expandReferences ({ $SystemAPI, $ComposeAPI }: { $SystemAPI: SystemAPI; $ComposeAPI: ComposeAPI }, m: Module, rr: Array<Revision>): Promise<unknown> {
    // collect user IDs to fetch
    const userIDs = new Set<string>()
    // const recordIDs = new Set<string>()
    //
    // const userRefFields: string[] = ['ownedBy']
    // const recRefFields: string[] = []

    // const numeric = /^\d+$/
    //
    // const addIDvalues = (set: Set<string>, values: unknown[]): void => {
    //   values.forEach(v => {
    //     if (typeof v === 'string' && numeric.test(v)) {
    //       set.add(v)
    //     }
    //   })
    // }

    rr.forEach(r => {
      // collecting revision authors
      userIDs.add(r.userID)

      // r.changes.forEach(c => {
      //   switch (true) {
      //     case userRefFields.includes(c.key):
      //       addIDvalues(userIDs, [...c.old, ...c.new])
      //       break
      //     case recRefFields.includes(c.key):
      //       addIDvalues(recordIDs, [...c.old, ...c.new])
      //       break
      //   }
      // })
    })

    return Promise.all([
      $SystemAPI.userList({ userID: [...userIDs] })
        .then((uu): void => {
          const map = new Map<string, User>()
          if (uu.set && Array.isArray(uu.set)) {
            uu.set.forEach(u => {
              map.set(u.userID, new User(u))
            })
          }

          // map user to revision author
          rr.filter(r => map.has(r.userID))
            .forEach(r => {
              r.user = map.get(r.userID) as User
              console.log('resolving revision author', r.user)
            })
        }),

      // @todo if the final decision on implementation of record-reference resolving is
      //       that it needs to be on the frontend we should fetch all records that are referenced
      //       and append the resolved versions back in the similar way as we are handling users.
    ])
  }
}

Registry.set(kind, PageBlockRecordRevisions)
