import { FilterDefinition } from './filter'
interface AggregateColumn {
    name: string;
    expr: string;
    aggregate: string;
    kind: string;
}
export interface StepLoad {
    name: string;
    source?: string;
    definition?: {
        [key: string]: unknown;
    };
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
export declare function StepFactory(step: Partial<Step>): Step;
export {}
