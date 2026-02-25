export declare enum NotificationKind {
  Simple = 'simple',
  Record = 'record',
}
export interface SimpleNotificationConfig {
  title: string
  description: string
}
export interface RecordNotificationConfig {
  title: string
  description: string
  moduleID: string
  namespaceID: string
  recordID: string
  openMode: string
  edit: boolean
}
export interface NotificationConfig {
  simple?: SimpleNotificationConfig
  record?: RecordNotificationConfig
}
interface PartialNotification extends Partial<
  Omit<Notification, 'createdAt' | 'updatedAt' | 'deletedAt' | 'readAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
  readAt?: string | number | Date
}
export declare class Notification {
  notificationID: string
  kind: NotificationKind
  config: NotificationConfig
  recipient: string
  createdBy: string
  createdAt?: Date
  readAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  constructor(n?: PartialNotification)
  apply(n?: PartialNotification): void
  /**
   * Returns resource ID
   */
  get resourceID(): string
  /**
   * Resource type
   */
  get resourceType(): string
  clone(): Notification
}
export {}
