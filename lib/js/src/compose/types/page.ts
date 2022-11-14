import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf, AreObjectsOf } from '../../guards'
import { PageBlock, PageBlockMaker } from './page-block'
import { Button } from './page-block/types'

interface PartialPage extends Partial<Omit<Page, 'children' | 'blocks' | 'createdAt' | 'updatedAt' | 'deletedAt'>> {
  children?: Array<PartialPage>;

  blocks?: (Partial<PageBlock>)[];

  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

interface PageConfig {
  buttons: {
    submit: Button;
    delete: Button;
    new: Button;
    edit: Button;
    clone: Button;
    back: Button;
  };
}

export class Page {
  public pageID = NoID;
  public selfID = NoID;
  public moduleID = NoID;
  public namespaceID = NoID;

  public title = '';
  public handle = '';
  public description = '';
  public weight = 0;

  public labels: object = {}

  public visible = false;

  public children?: Array<Page>

  public blocks: (InstanceType<typeof PageBlock>)[] = [];

  public config: PageConfig = {
    buttons: {
      submit: { enabled: true },
      delete: { enabled: true },
      new: { enabled: true },
      edit: { enabled: true },
      clone: { enabled: true },
      back: { enabled: true },
    },
  }

  public createdAt?: Date = undefined;
  public updatedAt?: Date = undefined;
  public deletedAt?: Date = undefined;

  public canUpdatePage = false;
  public canDeletePage = false;
  public canGrant = false;

  constructor (i?: PartialPage) {
    this.apply(i)
  }

  clone (): Page {
    return new Page(JSON.parse(JSON.stringify(this)))
  }

  apply (i?: PartialPage): void {
    if (!i) return

    Apply(this, i, CortezaID, 'pageID', 'selfID', 'moduleID', 'namespaceID')
    Apply(this, i, String, 'title', 'handle', 'description')
    Apply(this, i, Number, 'weight')
    Apply(this, i, Boolean, 'visible')

    if (i.blocks) {
      this.blocks = []
      if (AreObjectsOf<PageBlock>(i.blocks, 'kind') && AreObjectsOf<PageBlock>(i.blocks, 'xywh')) {
        this.blocks = i.blocks.map((b: { kind: string }) => PageBlockMaker(b))
      }
    }

    if (i.children) {
      this.children = []
      if (AreObjectsOf<Page>(i.children, 'pageID')) {
        this.children = i.children.map(c => new Page(c))
      }
    }

    if (i.config) {
      this.config = i.config
    }

    if (IsOf(i, 'labels')) {
      this.labels = { ...i.labels }
    }

    Apply(this, i, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')
    Apply(this, i, Boolean,
      'canUpdatePage',
      'canDeletePage',
      'canGrant',
    )
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.pageID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:page'
  }

  get isRecordPage (): boolean {
    return this.moduleID !== NoID
  }

  get firstLevel (): boolean {
    return this.selfID === NoID
  }

  /**
   * Validates page & it's blocks
   */
  validate (): Array<string> {
    const ee: Array<string> = []

    if (this.blocks.length === 0) {
      ee.push('blocks missing')
    } else {
      this.blocks.forEach(b => {
        ee.push(...b.validate())
      })
    }

    return ee
  }

  export (): PartialPage {
    return {
      title: this.title,
      handle: this.handle,
      description: this.description,
      visible: this.visible,
      blocks: this.blocks,
    }
  }
}
