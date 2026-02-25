interface FeedOptions {
  moduleID: string
  color: string
  prefilter: string
}
interface LegacyFeed {
  moduleID?: string
  startField?: string
  endField?: string
  titleField?: string
  allDay?: boolean
}
export type FeedInput = Partial<Feed> | Feed | LegacyFeed
/**
 * Feed class represents an event feed for the given calendar
 */
export default class Feed {
  resource: string
  startField: string
  endField: string
  titleField: string
  options: FeedOptions
  allDay: boolean
  constructor(i?: FeedInput)
  apply(i?: FeedInput): void
  static fromLegacy(legacy: LegacyFeed): Partial<Feed>
}
export {}
