import { CortezaID, NoID, ISO8601Date, Apply } from '../../cast'
import { IsOf } from '../../guards'

interface Meta {
  [key: string]: unknown;
}

interface PartialAttachment extends Partial<Omit<Attachment, 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

export class Attachment {
  public attachmentID = NoID;
  public ownerID = NoID;
  public name = '';
  public url = '';
  public previewUrl = '';
  public download = '';
  public meta: Meta = {};

  public createdAt?: Date = undefined;
  public updatedAt?: Date = undefined;
  public deletedAt?: Date = undefined;

  constructor (i?: PartialAttachment, baseURL?: string) {
    this.apply(i)
    this.setBaseURL(baseURL || '')
  }

  apply (i?: PartialAttachment): void {
    Apply(this, i, CortezaID, 'attachmentID', 'ownerID')
    Apply(this, i, String, 'name', 'url', 'previewUrl')

    if (IsOf(i, 'meta')) {
      this.meta = { ...i.meta }
    }

    Apply(this, i, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
  }

  setBaseURL (baseURL: string): void {
    this.url = baseURL + this.url
    this.previewUrl = baseURL + this.previewUrl
    this.download = this.url + '&download=1'
  }
}
