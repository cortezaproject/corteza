import { PageBlock, PageBlockInput, Registry } from './base'
import { Apply } from '../../../cast'

const kind = 'File'

interface Options {
  mode: string;
  attachments: string[];
  hideFileName: boolean;
  height: string;
  width: string;
  maxHeight: string;
  maxWidth: string;
  borderRadius: string;
  margin: string;
  backgroundColor: string;
  magnifyOption: string;
  clickToView?: boolean;
  enableDownload?: boolean;
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
  height: '',
  width: '',
  maxHeight: '',
  maxWidth: '',
  borderRadius: '',
  margin: 'auto',
  backgroundColor: '',
  magnifyOption: '',
  clickToView: true,
  enableDownload: true
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

    Apply(this.options, o, Boolean, 'hideFileName', 'clickToView', 'enableDownload')
    Apply(this.options, o, String, 'height', 'width', 'maxHeight', 'maxWidth', 'borderRadius', 'margin', 'backgroundColor', 'magnifyOption')

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
