import { merge } from 'lodash'
import { Apply, CortezaID, ISO8601Date, NoID } from '../../cast'
import { PageBlock } from './page-block/base'
import { Button } from './page-block/types'

export type PageLayoutInput = PageLayout | Partial<PageLayout>

interface PageLayoutConfig {
  visibility: Visibility;
  actions: Action[];
}

interface Action {
  actionID: string;
  kind: string;
  placement: string;
  params: unknown;
  meta: unknown;
}

interface Visibility {
  expression: string;
  roles: string[];
}

interface Meta {
  title: string;
  description: string;
}

export class PageLayout {
  public pageLayoutID = NoID;
  public namespaceID = NoID;
  public pageID = NoID
  public handle = '';

  public weight = 0;

  public blocks: (Partial<PageBlock>)[] = [];

  public config: PageLayoutConfig = {
    visibility: {
      expression: '',
      roles: [],
    },
    actions: [],
  }

  public meta: Meta = {
    title: '',
    description: ''
  };

  public createdAt?: Date = undefined;
  public updatedAt?: Date = undefined;
  public deletedAt?: Date = undefined;

  public ownedBy = NoID;

  constructor (pl?: PageLayoutInput) {
    this.apply(pl)
  }

  apply (pl?: PageLayoutInput): void {
    if (!pl) return

    Apply(this, pl, CortezaID, 'pageLayoutID', 'namespaceID', 'pageID', 'ownedBy')
    Apply(this, pl, String, 'handle')
    Apply(this, pl, ISO8601Date, 'createdAt', 'updatedAt', 'deletedAt')

    this.blocks = (pl.blocks || []).map(({ blockID, xywh, meta }) => ({ blockID, xywh, meta }))

    if (pl.meta) {
      this.meta = { ...this.meta, ...pl.meta }
    }

    if (pl.config) {
      this.config = merge({}, this.config, pl.config)
    }
  }

  clone (): PageLayout {
    return new PageLayout(JSON.parse(JSON.stringify(this)))
  }

  /**
   * Returns resource ID
   */
  get resourceID (): string {
    return `${this.resourceType}:${this.pageLayoutID}`
  }

  /**
   * Resource type
   */
  get resourceType (): string {
    return 'compose:page-layout'
  }

  export (): PageLayoutInput {
    return {
      blocks: this.blocks,
      meta: this.meta,
    }
  }
}
