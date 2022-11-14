import { feedResources } from './resources'
import { Apply, NoID } from '../../../../cast'
import { IsOf } from '../../../../guards'

interface FeedOptions {
  color: string;
  prefilter: string;
  moduleID: string;
}

interface LegacyFeed {
  moduleID?: string;
  startField?: string;
  endField?: string;
  titleField?: string;
  allDay?: boolean;
}

export type FeedInput = Partial<Feed> | Feed | LegacyFeed

const defOptions = {
  moduleID: NoID,
  color: '#ffffff',
  prefilter: '',
}

/**
 * Feed class represents an event feed for the given calendar
 */
export default class Feed {
  public resource = ''
  public startField = ''
  public endField = ''
  public titleField = ''
  public options: FeedOptions = defOptions

  public allDay = false

  constructor (i?: FeedInput) {
    this.apply(i)
  }

  apply (i?: FeedInput): void {
    if (!i) return

    if (!IsOf<Feed>(i, 'resource') && IsOf<LegacyFeed>(i, 'moduleID')) {
      i = Feed.fromLegacy(i)
    }

    if (IsOf<Feed>(i, 'resource')) {
      Apply(this, i, String, 'resource', 'startField', 'endField', 'titleField')
      Apply(this, i, Boolean, 'allDay')

      if (i.options) {
        this.options = { ...this.options, ...i.options }
      }
    }
  }

  static fromLegacy (legacy: LegacyFeed): Partial<Feed> {
    const p: Partial<Feed> = {
      // legacy does not have resource,
      // we've used it with records only
      resource: feedResources.record,

      ...legacy,
    }

    if (legacy.moduleID) {
      if (!p.options) {
        p.options = { ...defOptions }
      }

      // module was moved under options
      p.options.moduleID = legacy.moduleID
    }

    return p
  }
}
