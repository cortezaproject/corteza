interface KV {
  [_: string]: unknown
}
interface PartialReminder extends Partial<
  Omit<Reminder, 'assignedAt' | 'dismissedAt' | 'remindAt' | 'createdAt'>
> {
  assignedAt?: string | number | Date
  dismissedAt?: string | number | Date
  remindAt?: string | number | Date
  createdAt?: string | number | Date
}
export declare class Reminder {
  reminderID: string
  resource: string
  payload: KV
  snoozeCount: number
  assignedTo: string
  assignedBy: string
  assignedAt?: Date
  dismissedBy: string
  dismissedAt?: Date
  remindAt?: Date
  createdAt?: Date
  processed: boolean
  actions: KV
  options: KV
  constructor(r?: PartialReminder)
  apply(r?: PartialReminder): void
}
export {}
