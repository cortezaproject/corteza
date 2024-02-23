import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { IsOf, AreObjectsOf } from '../../guards'
import { PageBlock, PageBlockMaker } from './page-block'
import { merge } from 'lodash'

interface PartialPage extends Partial<Omit<Page, 'children' | 'meta' | 'blocks' |'createdAt' | 'updatedAt' | 'deletedAt'>> {
  children?: Array<PartialPage>;

  blocks?: PageBlock[];

  meta?: PageMeta;

  createdAt?: string|number|Date;
  updatedAt?: string|number|Date;
  deletedAt?: string|number|Date;
}

interface PageMeta {
  notifications: {
    enabled: boolean;
  };
}

interface PageConfig {
  navItem: {
    icon: {
      type: string;
      src: string;
    };
    expanded: false;
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

  public children?: Page[];

  public blocks: PageBlock[] = [];

  public config: PageConfig = {
    navItem: {
      icon: {
        type: '',
        src: '',
      },
      expanded: false,
    },
  }

  public meta: PageMeta = {
    notifications: {
      enabled: true,
    },
  };

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
      this.blocks = i.blocks.map(block => PageBlockMaker(block))
    }

    if (i.children) {
      this.children = []
      if (AreObjectsOf<Page>(i.children, 'pageID')) {
        this.children = i.children.map(c => new Page(c))
      }
    }

    if (IsOf(i, 'config')) {
      this.config = merge({}, this.config, i.config)
    }

    if (IsOf(i, 'meta')) {
      this.meta = merge({}, this.meta, i.meta)
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
