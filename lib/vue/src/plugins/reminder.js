/* eslint-disable @typescript-eslint/explicit-function-return-type */

import moment from 'moment'
import { system } from '@cortezaproject/corteza-js'

function intervalToMS (from, to) {
  if (!from || !to) {
    throw new Error('intervalToMS.invalidArgs')
  }
  return to - from
}

export class ReminderService {
  constructor ({ api, fetchOffset = 1000 * 60 * 5, resource = null, emitter } = {}) {
    if (!api) {
      throw new Error('reminderService.invalidParams')
    }

    this.emitter = emitter
    this.api = api
    this.fetchOffset = fetchOffset
    this.resource = resource

    this.set = []
    this.nextRemindAt = null
    this.tHandle = null
  }

  init (emitter, { filter = {} }) {
    if (emitter) {
      this.emitter = emitter
    }

    this.filter = {
      scheduledOnly: true,
      excludeDismissed: true,
      ...filter,
    }

    this.prefetch().then(rr => {
      this.enqueue(...rr)
      this.emitter.$emit('reminders.pull')
    })
  }

  /**
   * Fetches all reminders that are supposed to go off to date (time)
   *
   * @returns {Promise<system.Reminder>}
   */
  async prefetch () {
    return this.api.reminderList({
      limit: 0,
      resource: this.resource,
      toTime: moment().add(this.fetchOffset, 'min').toISOString(),
      ...this.filter,
    }).then(({ set }) => {
      return (set || []).map(r => new system.Reminder(r))
    })
  }

  enqueueRaw (raw) {
    this.enqueue(new system.Reminder(raw))
  }

  /**
   * Enqueue a given set of reminders
   * @param {Array<Reminder>} set Set of reminderIDs to enqueue
   */
  enqueue (...set) {
    set.forEach(r => {
      // New or replace
      const i = this.set.findIndex(({ reminderID }) => reminderID === r.reminderID)
      if (i > -1) {
        this.set.splice(i, 1, r)
      } else {
        this.set.push(r)
      }
    })

    // Should watcher restart
    const { changed, time } = this.findNextProcessTime(this.set, this.nextRemindAt)
    if (changed) {
      this.nextRemindAt = time
      this.scheduleReminderProcess(this.nextRemindAt)
    }
  }

  /**
   * Dequeue a given set of reminders
   * @param {Array} set Set of reminderIDs to remove
   */
  dequeue (IDs = []) {
    this.set = this.set.filter(({ reminderID }) => !IDs.includes(reminderID))

    // don't reuse time, since it could have been removed
    const { changed, time } = this.findNextProcessTime(this.set, null)
    if (changed) {
      this.nextRemindAt = time
      this.scheduleReminderProcess(this.nextRemindAt)
    }
  }

  /**
   * Determines we should use a new time
   * @param {Array} set Set of reminders to chose from
   * @param {Object|null} time Reference point
   * @private
   */
  findNextProcessTime (set = [], time = null) {
    let changed = false
    set.forEach(r => {
      if (!time || r.remindAt < time) {
        time = r.remindAt
        changed = true
      }
    })

    return { changed, time }
  }

  /**
   * Schedules processor ro run at the given time
   * @param {Moment} at When it should be ran
   * @param {Moment} now Ref to now; used for tests
   * @private
   */
  scheduleReminderProcess (at, now = new Date()) {
    // Determine ms until next reminder should be processed
    const t = intervalToMS(now, at)

    if (this.tHandle != null) {
      window.clearTimeout(this.tHandle)
    }
    this.tHandle = window.setTimeout(this.processQueue.bind(this), t)
  }

  /**
   * Processes the reminder queue. Emits due reminders &
   * removes them from state
   * @param {Moment} now Ref to now; used for tests
   * @private
   */
  processQueue (now = new Date()) {
    let nextRemindAt = null

    this.set.forEach(r => {
      if (now >= r.remindAt) {
        if (this.emitter) {
          this.emitter.$emit('reminder.show', r)
        } else {
          throw new Error('pool.noEmitter')
        }
        r.processed = true
      } else if (now < r.remindAt && (!nextRemindAt || r.remindAt < nextRemindAt)) {
        nextRemindAt = r.remindAt
      }
    })

    this.nextRemindAt = nextRemindAt
    this.set = this.set.filter(({ processed }) => !processed)
    if (this.nextRemindAt === null) {
      this.tHandle = null
    } else {
      this.scheduleReminderProcess(this.nextRemindAt)
    }
  }
}

export default {
  install (Vue, opts) {
    Vue.prototype.$Reminder = new ReminderService(opts)
  },
}
