import { Capabilities, ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'

const kind = 'Geometry'

interface GeometryOptions extends Options {
  center: number[];
  zoom: number;
  multiDelimiter: string;
  prefillWithCurrentLocation: boolean;
  hideCurrentLocationButton: boolean;
  hideGeoSearch: boolean,
}

const defaults = (): Readonly<GeometryOptions> => Object.freeze({
  ...defaultOptions(),
  center: [30, 30],
  zoom: 3,
  multiDelimiter: '\n',
  prefillWithCurrentLocation: false,
  hideCurrentLocationButton: false,
  hideGeoSearch: false,
})

export class ModuleFieldGeometry extends ModuleField {
  readonly kind = kind

  options: GeometryOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldGeometry>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<GeometryOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'multiDelimiter')
    Apply(this.options, o, Number, 'zoom')
    Apply(this.options, o, Boolean, 'prefillWithCurrentLocation', 'hideCurrentLocationButton', 'hideGeoSearch')

    if (o.center) {
      this.options.center = o.center
    }
  }

  /**
   * Per module field type capabilities
   */
  public get cap (): Readonly<Capabilities> {
    return {
      ...super.cap,
      multi: true,
    }
  }
}

Registry.set(kind, ModuleFieldGeometry)
