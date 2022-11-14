import { FilterDefinition } from './filter'

export interface FrameColumn {
  kind?: string;
  label?: string;
  name?: string;
  primary?: boolean;
  unique?: boolean;
}

export interface FramePaging {
  limit?: number;
  cursor?: string;
}

export class FrameDefinition {
  name?: string;
  source?: string;
  ref?: string;

  sort?: string;
  filter?: FilterDefinition;
  paging?: FramePaging

  refValue?: string;
  relColumn?: string;
  relSource?: string;

  columns?: Array<FrameColumn>;
  rows?: Array<string>;
}

export interface DefinitionOptions {
  [header: string]: FrameDefinition;
}
