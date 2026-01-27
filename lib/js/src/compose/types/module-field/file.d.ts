import { ModuleField, Options } from './base'
export declare const modes: string[]
interface FileOptions extends Options {
    allowImages: boolean;
    allowDocuments: boolean;
    maxSize: number;
    mode: string;
    inline: boolean;
    hideFileName: boolean;
    mimetypes?: string;
    height?: string;
    width?: string;
    maxHeight?: string;
    maxWidth?: string;
    borderRadius?: string;
    margin?: string;
    backgroundColor?: string;
    clickToView?: boolean;
    enableDownload?: boolean;
    multiDelimiter: string;
    enableWebcam?: boolean;
}
export declare class ModuleFieldFile extends ModuleField {
  readonly kind = 'File'
  options: FileOptions
  constructor(i?: Partial<ModuleFieldFile>);
  applyOptions(o?: Partial<FileOptions>): void;
}
export {}
