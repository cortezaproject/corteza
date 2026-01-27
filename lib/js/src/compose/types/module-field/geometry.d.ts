import { Capabilities, ModuleField, Options } from './base'
interface GeometryOptions extends Options {
    center: number[];
    zoom: number;
    multiDelimiter: string;
    prefillWithCurrentLocation: boolean;
    hideCurrentLocationButton: boolean;
    hideGeoSearch: boolean;
}
export declare class ModuleFieldGeometry extends ModuleField {
  readonly kind = 'Geometry'
  options: GeometryOptions
  constructor(i?: Partial<ModuleFieldGeometry>);
  applyOptions(o?: Partial<GeometryOptions>): void;
  /**
     * Per module field type capabilities
     */
  get cap(): Readonly<Capabilities>;
}
export {}
