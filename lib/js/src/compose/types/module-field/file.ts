import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply, ApplyWhitelisted } from '../../../cast'

const kind = 'File'

export const modes = [
  // list of attachments, no preview
  'list',
  // grid of icons
  'grid',
  // single (first) image/file, show preview
  'single',
  // list of all images/files, show preview
  'gallery',
]

interface FileOptions extends Options {
  allowImages: boolean;
  allowDocuments: boolean;
  maxSize: number;
  mode: string;
  inline: boolean;
  hideFileName: boolean;
  mimetypes?: string;
  height?: number;
  width?: number;
  maxHeight?: number;
  maxWidth?: number;
  borderRadius?: number;
  margin?: number;
  backgroundColor?: string;
}

const defaults = (): Readonly<FileOptions> => Object.freeze({
  ...defaultOptions(),
  allowImages: true,
  allowDocuments: true,
  maxSize: 0,
  mode: '\n',
  inline: true,
  hideFileName: false,
  mimetypes: '',
  height: undefined,
  width: undefined,
  maxHeight: undefined,
  maxWidth: undefined,
  borderRadius: undefined,
  margin: undefined,
  backgroundColor: undefined,
})

export class ModuleFieldFile extends ModuleField {
  readonly kind = kind

  options: FileOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldFile>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<FileOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, Number, 'maxSize', 'height', 'width', 'maxHeight', 'maxWidth', 'borderRadius', 'margin')
    Apply(this.options, o, Boolean, 'allowImages', 'allowDocuments', 'inline', 'hideFileName')
    Apply(this.options, o, String, 'mimetypes', 'backgroundColor')
    ApplyWhitelisted(this.options, o, modes, 'mode')
  }
}

Registry.set(kind, ModuleFieldFile)
