interface Meta {
  [key: string]: unknown
}
interface PartialAttachment extends Partial<
  Omit<Attachment, 'createdAt' | 'updatedAt' | 'deletedAt'>
> {
  createdAt?: string | number | Date
  updatedAt?: string | number | Date
  deletedAt?: string | number | Date
}
export declare class Attachment {
  attachmentID: string
  ownerID: string
  name: string
  url: string
  previewUrl: string
  download: string
  meta: Meta
  createdAt?: Date
  updatedAt?: Date
  deletedAt?: Date
  constructor(i?: PartialAttachment, baseURL?: string)
  apply(i?: PartialAttachment): void
  setBaseURL(baseURL: string): void
}
export {}
