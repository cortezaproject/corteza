import { DisplayElement, Registry } from './base'
export { DisplayElementChart, ChartOptionsMaker } from './chart'
export { DisplayElementTable } from './table'
export { DisplayElementText } from './text'
export { DisplayElementMetric } from './metric'

export function DisplayElementMaker<T extends DisplayElement> (i: { kind: string }): T {
  const DisplayElementTemp = Registry.get(i.kind)
  if (DisplayElementTemp === undefined) {
    throw new Error(`unknown display element kind '${i.kind}'`)
  }

  if (i instanceof DisplayElement) {
    // Get rid of the references
    i = JSON.parse(JSON.stringify(i))
  }

  return new DisplayElementTemp(i) as T
}

export {
  Registry as DisplayElementRegistry,
  DisplayElement,
}
