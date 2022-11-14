import { FilterDefinition } from './filter'

interface GroupColumn {
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

export interface StepJoin {
  name: string;
  localSource: string;
  localColumn: string;
  foreignSource: string;
  foreignColumn: string;
}

export interface StepGroup {
  name: string;
  source: string;
  filter?: FilterDefinition;
  keys?: Array<GroupColumn>;
  columns?: Array<GroupColumn>;
  sort?: string;
}

interface TransformColumn {
  name: string;
  expr: string;
  rows: unknown;
}

export interface StepTransform {
  name?: string;
  dimension?: string;
  columns?: Array<TransformColumn>;
}

export interface Step {
  group?: StepGroup;
  load?: StepLoad;
  join?: StepJoin;
  transform?: StepTransform;
}

export function StepFactory (step: Partial<Step>): Step {
  const k = Object.keys(step)[0]
  switch (k) {
    case 'load':
      return step as Step
    case 'join':
      return step as Step
    case 'group':
      return step as Step
    case 'transform':
      return step as Step
    default:
      throw new Error('unknown step: ' + k)
  }
}
