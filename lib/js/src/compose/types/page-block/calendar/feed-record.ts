import { Record } from '../../record'
import { Namespace } from '../../namespace'
import { Module } from '../../module'
import { Compose as ComposeAPI } from '../../../../api-clients'
import { makeColors, Event } from './shared'

// import variables from 'corteza-webapp-compose/src/themes/corteza-base/variables.scss'
// const defaultColor = variables.primary
// @todo fix this!
const defaultColor = '#568ba2'

interface FeedOptions {
  color: string;
  prefilter: string;
}

interface Feed {
  startField: string;
  endField: string;
  titleField: string;
  options: FeedOptions;
  allDay: boolean;
}

interface Range {
  end: Date;
  start: Date;
}

function getRecordValue (record: Readonly<Record>, field: string): (string|undefined)[] {
  const ef = record.module.fields.find(({ name }) => name === field)
  if (ef) {
    return ef.isMulti ? record.values[field] as string[] : [(record.values[field] as string) || undefined]
  } else {
    switch (field) {
      case 'recordID':
      case 'moduleID':
      case 'namespaceID':
        return [record[field]]
      case 'createdAt':
      case 'updatedAt':
      case 'deletedAt':
        if (record[field] !== undefined) {
          return [(record[field] as Date).toISOString()]
        }
        break
    }
  }

  return [undefined]
}

/**
 * Method expands the given record in a (set) of FC event objects.
 * Handles basic recurrence -- multiple date fields.
 * @param {Record} record Record to expand
 * @param {Feed} feed Feed, this record belongs to
 * @returns {Array} A set of expanded events
 */
function expandRecord (record: Readonly<Record>, feed: Feed): Event[] {
  const events: Event[] = []

  const starts = getRecordValue(record, feed.startField)
  const ends = getRecordValue(record, feed.endField)
  const title = getRecordValue(record, feed.titleField).shift() || record.recordID

  // Make sure ends is at least as long as starts, to avoid length checks
  ends.push(...(new Array(Math.max(starts.length - ends.length, 0)).fill(undefined)))

  const classNames = ['event', 'event-record']
  const { backgroundColor, borderColor, isLight } = makeColors(feed.options.color || defaultColor)
  if (isLight) {
    classNames.push('text-dark')
  } else {
    classNames.push('text-white')
  }

  starts.forEach((start, i) => {
    events.push({
      // So FC knows how to group these expanded events
      groupId: record.recordID,
      id: record.recordID,
      title,
      start: start,
      end: ends[i],
      allDay: feed.allDay,
      backgroundColor,
      borderColor,
      classNames,

      extendedProps: {
        moduleID: record.module.moduleID,
        recordID: record.recordID,
      },
    })
  })

  return events
}

/**
 * Checks if the given field can be used with the given record.
 * A field can be used if it's either defined as a record value OR it's a system field.
 *
 * @param r The record to check
 * @param field The field we wish to use
 */
function recordFeedFilter (r: Readonly<Record>, field: string): boolean {
  if (r.values[field]) {
    return true
  }

  return ['createdAt', 'updatedAt', 'deletedAt'].includes(field)
}

/**
 * Loads & converts module resource into FC events
 * @param {ComposeAPI} $ComposeAPI ComposeAPI provider
 * @param {Module} module Current module
 * @param {Namespace} namespace Current namespace
 * @param {Feed} feed Current feed
 * @param {Object} range Current date range
 * @returns {Promise<Array>} Resolves to a set of FC events to display
 */
export async function RecordFeed ($ComposeAPI: ComposeAPI, module: Module, namespace: Namespace, feed: Feed, range: Range): Promise<Event[]> {
  // Params for record fetching
  const params = {
    namespaceID: namespace.namespaceID,
    moduleID: module.moduleID,
    query: `(date(${feed.startField}) <= '${range.end.toISOString()}' AND '${range.start.toISOString()}' <= date(${feed.endField || feed.startField}))`,
  }

  if (feed.options.prefilter) {
    params.query += ` AND (${feed.options.prefilter})`
  }

  const events: Array<Event> = []
  return $ComposeAPI.recordList(params).then(({ set }) => {
    (set as Array<{ recordID: string }>)

      // Removes all duplicates
      .filter(({ recordID }, index, set) => set.findIndex((r) => recordID === r.recordID) === index)

      // cast & freeze
      .map(r => Object.freeze(new Record(module, r)))

      // drop record w/o proper values
      .filter(r => recordFeedFilter(r, feed.startField))
      // eslint-disable-next-line @typescript-eslint/no-use-before-define
      .forEach(r => events.push(...expandRecord(r, feed)))
    return events
  })
}
