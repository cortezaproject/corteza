import { FilterDefinition } from './filter'

interface AggregateColumn {
  name: string;
  expr: string;
  aggregate: string;
  // Kind specifies what kind the result will be.
  // This is lame and will change, but for now, bare with me.
  kind: string;
}

export interface StepLoad {
  name: string;
  source?: string;
  definition?: { [key: string]: unknown};
  filter?: FilterDefinition;
  sort?: string;
}

export interface StepLink {
  name: string;
  localSource: string;
  localColumn: string;
  foreignSource: string;
  foreignColumn: string;
}

export interface StepJoin {
  name: string;
  localSource: string;
  localColumn: string;
  foreignSource: string;
  foreignColumn: string;
}

export interface StepAggregate {
  name: string;
  source: string;
  filter?: FilterDefinition;
  keys?: Array<AggregateColumn>;
  columns?: Array<AggregateColumn>;
  sort?: string;
}

export interface Step {
  aggregate?: StepAggregate;
  load?: StepLoad;
  link?: StepLink;
}

export function StepFactory (step: Partial<Step>): Step {
  const k = Object.keys(step)[0]
  switch (k) {
    case 'load':
      return step as Step
    case 'link':
      return step as Step
    case 'join':
      return step as Step
    case 'aggregate':
      return step as Step
    default:
      throw new Error('unknown step: ' + k)
  }
}
