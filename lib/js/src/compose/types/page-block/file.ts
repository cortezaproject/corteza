import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'File'

interface Options {
  mode: string;
  attachments: string[];
  hideFileName: boolean;
  height?: number;
  width?: number;
  maxHeight?: number;
  maxWidth?: number;
  borderRadius?: number;
  margin?: number;
  backgroundColor?: string;
}

const PageBlockFileDefaultMode = 'list'
const PageBlockFileModes = [
  // list of attachments, no preview
  'list',
  // grid of icons
  'grid',
  // single (first) image/file, show preview
  'single',
  // list of all images/files, show preview
  'gallery',
]

const defaults: Readonly<Options> = Object.freeze({
  mode: PageBlockFileDefaultMode,
  attachments: [],
  hideFileName: false,
  height: undefined,
  width: undefined,
  maxHeight: undefined,
  maxWidth: undefined,
  borderRadius: undefined,
  margin: undefined,
  backgroundColor: undefined,
})

export class PageBlockFile extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    if (o.attachments) {
      this.options.attachments = o.attachments
    }

    Apply(this.options, o, Boolean, 'hideFileName')
    Apply(this.options, o, String, 'backgroundColor')
    Apply(this.options, o, Number, 'height', 'width', 'maxHeight', 'maxWidth', 'borderRadius', 'margin')

    if (o.mode) {
      if (PageBlockFileModes.includes(o.mode)) {
        this.options.mode = o.mode
      } else {
        o.mode = PageBlockFileDefaultMode
      }
    }
  }
}

Registry.set(kind, PageBlockFile)
