import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf } from '../../guards'

export enum NotificationKind {
  Simple = 'simple',
  Record = 'record',
}

export interface SimpleNotificationConfig {
  title: string;
  description: string;
}

export interface RecordNotificationConfig {
  title: string;
  description: string;
  moduleID: string;
  namespaceID: string;
  recordID: string;
  openMode: string; // "modal", "newTab", or "sameTab" (default)
  edit: boolean; // Whether to open the record in edit mode
}

export interface NotificationConfig {
  simple?: SimpleNotificationConfig;
  record?: RecordNotificationConfig;
}

interface PartialNotification extends Partial<Omit<Notification, 'createdAt' | 'updatedAt' | 'deletedAt' | 'readAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
  readAt?: string|number|Date;
}

export class Notification {
  public notificationID = NoID
  public kind: NotificationKind = NotificationKind.Simple
  public config: NotificationConfig = { simple: { title: '', description: '' } }
  public recipient = NoID
  public createdBy = NoID
  public createdAt?: Date = undefined
  public readAt?: Date = undefined
  public updatedAt?: Date = undefined
  public deletedAt?: Date = undefined

  constructor (n?: PartialNotification) {
    this.apply(n)
  }

  apply (n?: PartialNotification): void {
    if (!n) return

    Apply(this, n, CortezaID, 'notificationID', 'recipient', 'createdBy')
    Apply(this, n, String, 'kind')
    Apply(this, n, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt', 'readAt')

    if (IsOf(n, 'config')) {
      this.config = n.config[this.kind] || {}
    }
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.notificationID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'system:notification'
  }

  clone (): Notification {
    return new Notification(JSON.parse(JSON.stringify(this)))
  }
}
