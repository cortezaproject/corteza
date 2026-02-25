import { DisplayElement, Registry } from './base'
export { DisplayElementChart, ChartOptionsMaker } from './chart'
export { DisplayElementTable } from './table'
export { DisplayElementText } from './text'
export { DisplayElementMetric } from './metric'
export declare function DisplayElementMaker<T extends DisplayElement>(i: { kind: string }): T
export { Registry as DisplayElementRegistry, DisplayElement }
