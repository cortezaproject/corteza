import { System as SystemAPI } from '../../../../api-clients'
import { User, Reminder } from '../../../../system'
import { makeColors, Event } from './shared'
import { AreObjectsOf } from '../../../../guards'

interface FeedOptions {
  color: string;
}

interface Feed {
  options: FeedOptions;
}

interface Range {
  end: Date;
  start: Date;
}

/**
 * Loads & converts reminder resource into FC events
 * @param {SystemAPI} $SystemAPI SystemAPI provider
 * @param {User} user Current user
 * @param {Feed} feed Current feed
 * @param {Object} range Current date range
 * @returns {Promise<Array>} Resolves to a set of FC events to display
 */
export async function ReminderFeed($SystemAPI: SystemAPI, user: User, feed: Feed, range: Range, options = {}): Promise<Event[]> {
  return $SystemAPI.reminderList({
    scheduledFrom: range.start.toISOString(),
    scheduledUntil: range.end.toISOString(),
    scheduledOnly: true,
    excludeDismissed: true,
    assignedTo: user.userID,
  }, options).then(({ set }) => {
    const { backgroundColor, borderColor, textColor } = makeColors(feed.options.color)

    if (!AreObjectsOf<Reminder>(set, 'reminderID', 'assignedTo', 'remindAt', 'payload')) {
      return []
    }

    return set.map(r => {
      r = new Reminder(r)

      const classNames = ['event', 'event-reminder']
      if (r.assignedTo !== user.userID) {
        classNames.push('event-not-owner')
      }

      const e: Event = {
        id: r.reminderID,
        title: (r.payload.title as string) ?? r.reminderID ?? '-',
        start: r.remindAt?.toISOString(),
        backgroundColor,
        borderColor,
        textColor,
        classNames,
        allDay: false,
        extendedProps: {
          reminderID: r.reminderID,
        },
      }

      return e
    })
  })
}
