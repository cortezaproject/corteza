import { merge } from 'lodash'
import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'

interface KV {
  [_: string]: unknown;
}

interface PartialReminder extends Partial<Omit<Reminder, 'assignedAt' | 'dismissedAt' | 'remindAt' | 'createdAt'>> {
  assignedAt?: string|number|Date;
  dismissedAt?: string|number|Date;
  remindAt?: string|number|Date;
  createdAt?: string|number|Date;
}

export class Reminder {
  public reminderID = NoID
  public resource = NoID
  public payload: KV = {}
  public snoozeCount = 0
  public assignedTo = NoID
  public assignedBy = NoID
  public assignedAt?: Date = undefined
  public dismissedBy = NoID
  public dismissedAt?: Date = undefined
  public remindAt?: Date = undefined
  public createdAt?: Date = undefined
  public processed = false
  public actions: KV = {}
  public options: KV = {}

  constructor (r?: PartialReminder) {
    this.apply(r)
  }

  apply (r?: PartialReminder): void {
    if (!r) return

    Apply(this, r, CortezaID, 'reminderID')
    Apply(this, r, Number, 'snoozeCount')
    Apply(this, r, CortezaID, 'assignedTo', 'assignedBy', 'dismissedBy')

    // @todo actions, options, payload... all 3?
    this.payload = merge({}, this.payload, r.payload)
    this.actions = merge({}, this.actions, r.actions)
    this.options = merge({}, this.options, r.options)

    Apply(this, r, ISO8601Date, 'assignedAt', 'dismissedAt', 'remindAt', 'createdAt')
    Apply(this, r, Boolean, 'processed')
  }
}
